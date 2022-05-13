package pubsub

import (
	"context"

	googlepubsub "cloud.google.com/go/pubsub"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type SubscriberConfig struct {
	ProjectID                 string `mapstructure:"project_id"`
	DoNotCreateTopicIfMissing bool   `mapstructure:"do_not_create_topic_if_missing"`
	EnableMessageOrdering     bool   `mapstructure:"enable_message_ordering"`
	SubscriptionSuffix        string `mapstructure:"subscription_suffix"`
	LoggerDebug               bool   `mapstructure:"logger_debug"`
	LoggerTrace               bool   `mapstructure:"logger_trace"`
}

type Subscriber struct {
	config *SubscriberConfig
	sub    *googlecloud.Subscriber
}

func NewSubscriber(config *SubscriberConfig) (*Subscriber, error) {
	if config.SubscriptionSuffix == "" {
		return nil, errors.New("Subscription suffix is required")
	}
	logger := watermill.NewStdLogger(config.LoggerDebug, config.LoggerTrace)
	generateSubscriptionName := googlecloud.TopicSubscriptionNameWithSuffix("-" + config.SubscriptionSuffix)
	sub, err := googlecloud.NewSubscriber(googlecloud.SubscriberConfig{
		GenerateSubscriptionName:  generateSubscriptionName,
		ProjectID:                 config.ProjectID,
		DoNotCreateTopicIfMissing: config.DoNotCreateTopicIfMissing,
		SubscriptionConfig: googlepubsub.SubscriptionConfig{
			EnableMessageOrdering: config.EnableMessageOrdering,
		},
	}, logger)
	if err != nil {
		return nil, errors.Wrap(err, "can't create google cloud subscriber")
	}
	return &Subscriber{
		config: config,
		sub:    sub,
	}, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	messages, err := s.sub.Subscribe(ctx, topic)
	if err != nil {
		return nil, errors.Wrap(err, "can't subscribe")
	}
	return messages, nil
}

func (s *Subscriber) Close() error {
	err := s.sub.Close()
	if err != nil {
		return errors.Wrap(err, "can't close subscriber")
	}
	return err
}
