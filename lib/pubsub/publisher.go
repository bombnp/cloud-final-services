package pubsub

import (
	"context"
	"log"
	"sync"
	"time"

	googlepubsub "cloud.google.com/go/pubsub"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type PublisherConfig struct {
	ProjectID                 string `mapstructure:"project_id"`
	DoNotCreateTopicIfMissing bool   `mapstructure:"do_not_create_topic_if_missing"`
	EnableMessageOrdering     bool   `mapstructure:"enable_message_ordering"`
	DelayThreshold            *int   `mapstructure:"delay_threshold"` // seconds to wait before publishing the batch
}

type Publisher struct {
	config *PublisherConfig
	client *googlepubsub.Client

	topics     map[string]*googlepubsub.Topic
	topicsLock sync.RWMutex

	closed bool
}

func NewPublisher(config *PublisherConfig) (*Publisher, error) {
	ctx := context.Background()
	client, err := googlepubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "can't create google pubsub client")
	}
	return &Publisher{
		config: config,
		client: client,
		topics: map[string]*googlepubsub.Topic{},
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, topic string, orderingKey string, messages ...*message.Message) error {
	if p.closed {
		return errors.New("Publisher is closed")
	}

	t, err := p.getTopic(ctx, topic)
	if err != nil {
		return errors.Wrapf(err, "can't get topic %s", topic)
	}

	for _, msg := range messages {
		pubsubMsg, err := Marshal(msg, orderingKey)
		if err != nil {
			return errors.Wrapf(err, "cannot marshal message %s", msg.UUID)
		}

		result := t.Publish(ctx, pubsubMsg)
		go func(result *googlepubsub.PublishResult) {
			_, err := result.Get(ctx)
			if err != nil {
				t.ResumePublish(orderingKey)
				log.Println("publish message failed", err)
			}
		}(result)
	}

	return nil
}

func (p *Publisher) Close() error {
	if p.closed {
		return nil
	}
	p.closed = true

	p.topicsLock.Lock()
	for _, t := range p.topics {
		t.Stop()
	}
	p.topicsLock.Unlock()

	return p.client.Close()
}
func (p *Publisher) getTopic(ctx context.Context, topic string) (t *googlepubsub.Topic, err error) {
	p.topicsLock.RLock()
	t, ok := p.topics[topic]
	p.topicsLock.RUnlock()
	// if topic exists in map, return
	if ok {
		return t, nil
	}

	// if not, create a new topic instance
	p.topicsLock.Lock()
	defer func() {
		if err == nil {
			p.topics[topic] = t
		}
		p.topicsLock.Unlock()
	}()

	t = p.client.Topic(topic)
	t.EnableMessageOrdering = p.config.EnableMessageOrdering
	if p.config.DelayThreshold != nil {
		delayThreshold := time.Duration(*p.config.DelayThreshold) * time.Second
		t.PublishSettings.DelayThreshold = delayThreshold
		log.Printf("Topic %s delay threshold set to %s", topic, delayThreshold)
	}

	exists, err := t.Exists(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "could not check if topic %s exists", topic)
	}

	if exists {
		return t, nil
	}

	if p.config.DoNotCreateTopicIfMissing {
		return nil, errors.Errorf("topic does not exist: %s", topic)
	}

	t, err = p.client.CreateTopic(ctx, topic)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create topic %s", topic)
	}

	return t, nil
}
