package pubsub

import (
	googlepubsub "cloud.google.com/go/pubsub"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

// UUIDHeaderKey is the key of the Pub/Sub attribute that carries Waterfall UUID.
const UUIDHeaderKey = "_watermill_message_uuid"

func Marshal(msg *message.Message, orderingKey string) (*googlepubsub.Message, error) {
	if value := msg.Metadata.Get(UUIDHeaderKey); value != "" {
		return nil, errors.Errorf("metadata %s is reserved by watermill for message UUID", UUIDHeaderKey)
	}

	attributes := map[string]string{
		UUIDHeaderKey: msg.UUID,
	}

	for k, v := range msg.Metadata {
		attributes[k] = v
	}

	marshaledMsg := &googlepubsub.Message{
		Data:        msg.Payload,
		Attributes:  attributes,
		OrderingKey: orderingKey,
	}

	return marshaledMsg, nil
}
