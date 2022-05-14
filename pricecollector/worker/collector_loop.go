package worker

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"time"

	"github.com/bombnp/cloud-final-services/lib/pubsub"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/pkg/errors"
)

func (c *collector) LoopCollectEvents(ctx context.Context) error {
	msgCh, err := c.sub.Subscribe(ctx, pubsub.SyncEventsTopic)
	if err != nil {
		return errors.Wrapf(err, "can't subscribe %s\n", pubsub.SyncEventsTopic)
	}
	log.Println("subscribed to", pubsub.SyncEventsTopic)
	log.Println("processing messages...")
	for msg := range msgCh {
		var events []pubsub.SyncEventMsg
		err = json.Unmarshal(msg.Payload, &events)
		if err != nil {
			log.Println("can't unmarshal sync event", err)
			continue
		}
		c.writePricePoints(events)
		msg.Ack()
	}
	return errors.New("channel unexpectedly closed")
}

func (c *collector) writePricePoints(events []pubsub.SyncEventMsg) {
	writeApi := c.repo.Influx.WriteAPI
	for _, event := range events {
		p, err := c.pointFromEvent(event)
		if err != nil {
			log.Println(errors.Wrap(err, "cant write price point").Error())
			continue
		}
		writeApi.WritePoint(p)
	}
	log.Println("Wrote points of block", events[0].Block)
}

func (c *collector) pointFromEvent(event pubsub.SyncEventMsg) (*write.Point, error) {
	pair, ok := c.pairMap[event.Address]
	if !ok {
		return nil, errors.New("unknown pair address")
	}
	baseReserve := new(big.Float)
	quoteReserve := new(big.Float)
	if pair.IsBaseToken0 {
		baseReserve.SetInt(event.Reserve0)
		quoteReserve.SetInt(event.Reserve1)
	} else {
		baseReserve.SetInt(event.Reserve1)
		quoteReserve.SetInt(event.Reserve0)
	}
	price, _ := new(big.Float).Quo(quoteReserve, baseReserve).Float64()
	tags := map[string]string{
		"address": event.Address.String(),
	}
	fields := map[string]any{
		"price": price,
	}
	return influxdb2.NewPoint("price", tags, fields, time.Unix(int64(event.Timestamp), 0)), nil
}
