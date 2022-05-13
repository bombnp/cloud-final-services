package pubsub

import (
	"context"
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
}

type publisher struct {
	config *PublisherConfig
	client *googlepubsub.Client

	topics     map[string]*googlepubsub.Topic
	topicsLock sync.RWMutex

	closed bool
}

func NewPublisher(config *PublisherConfig) (*publisher, error) {
	ctx := context.Background()
	client, err := googlepubsub.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "can't create google pubsub client")
	}
	return &publisher{
		config: config,
		client: client,
		topics: map[string]*googlepubsub.Topic{},
	}, nil
}

func (p *publisher) Publish(ctx context.Context, topic string, orderingKey string, messages ...*message.Message) error {
	if p.closed {
		return errors.New("publisher is closed")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	t, err := p.getTopic(ctx, topic)
	if err != nil {
		return errors.Wrapf(err, "can't get topic %s", topic)
	}

	// intentionally waiting for each message to finish before publishing a new one
	for _, msg := range messages {
		pubsubMsg, err := Marshal(msg, orderingKey)
		if err != nil {
			return errors.Wrapf(err, "cannot marshal message %s", msg.UUID)
		}

		result := t.Publish(ctx, pubsubMsg)
		<-result.Ready()

		_, err = result.Get(ctx)
		if err != nil {
			t.ResumePublish(orderingKey)
			return errors.Wrapf(err, "publishing message %s failed", msg.UUID)
		}
	}

	return nil
}

func (p *publisher) Close() error {
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
func (p *publisher) getTopic(ctx context.Context, topic string) (t *googlepubsub.Topic, err error) {
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
