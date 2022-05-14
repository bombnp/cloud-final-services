package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bombnp/cloud-final-services/lib/ethutils"
	"github.com/bombnp/cloud-final-services/lib/pubsub"
	"github.com/bombnp/cloud-final-services/txpublisher/config"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var currentBlock uint64
var currentIsRemoved bool
var blockStartTime = time.Now()
var syncEvents []pubsub.SyncEventMsg

func (s *streamer) LoopConsumeLog(ctx context.Context) error {
	logCh := s.logCh

	flushTime := 50 * time.Millisecond
	flushTimer := time.NewTimer(flushTime)

	log.Println("streamer_loop: Started logs processor.")
	for {
		select {
		case eventLog, ok := <-logCh:
			// channel unexpectedly closed.
			if !ok {
				return errors.New("channel unexpectedly closed")
			}
			flushTimer.Stop()

			err := s.ProcessLog(ctx, eventLog)
			if err != nil {
				log.Println(errors.Wrap(err, "can't process log").Error())
			}

			flushTimer.Reset(flushTime)
		case err := <-s.logSub.Err():
			errClose := s.pub.Close()
			log.Println("error while closing publisher", errClose)
			return errors.Wrap(err, "error during subscription")
		case <-ctx.Done():
			return ctx.Err()
		case <-flushTimer.C:
			// flush messages if channel is empty AND 50ms has passed
			if len(syncEvents) != 0 {
				s.publishMessages(ctx, syncEvents, time.Since(blockStartTime))
				syncEvents = []pubsub.SyncEventMsg{}
			}
		}
	}
}

func (s *streamer) ProcessLog(ctx context.Context, log ethtypes.Log) error {
	// flush messages if new block comes
	if log.BlockNumber != currentBlock || log.Removed != currentIsRemoved {
		if len(syncEvents) != 0 {
			s.publishMessages(ctx, syncEvents, time.Since(blockStartTime))
			syncEvents = []pubsub.SyncEventMsg{}
		}
		currentBlock = log.BlockNumber
		currentIsRemoved = log.Removed
		blockStartTime = time.Now()
	}

	lp, _ := ethutils.GetLp(ethutils.DexTypeUniswapV2)
	topic := log.Topics[0]
	switch topic {
	case lp.SyncEventID():
		syncEvent, err := lp.UnmarshalSyncEvent(log.Data)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal sync event")
		}
		timestamp, err := s.btCache.GetTimestamp(ctx, log.BlockNumber)
		if err != nil {
			return errors.Wrap(err, "failed to get block timestamp")
		}
		syncEventMsg := pubsub.SyncEventMsg{
			Address:   log.Address,
			Block:     log.BlockNumber,
			Timestamp: timestamp,
			Reserve0:  syncEvent.Reserve0,
			Reserve1:  syncEvent.Reserve1,
		}
		syncEvents = append(syncEvents, syncEventMsg)
		break
	}
	return nil
}

func (s *streamer) publishMessages(ctx context.Context, syncEvents []pubsub.SyncEventMsg, timeTaken time.Duration) {
	conf := config.InitConfig()
	pub := s.pub
	block := syncEvents[0].Block
	if len(syncEvents) == 0 {
		log.Printf("no new sync events to be published. block: %d, timeTaken (ms): %d\n", block, timeTaken.Milliseconds())
		return
	}
	out, err := json.Marshal(syncEvents)
	if err != nil {
		log.Printf("can't marshal sync events to json. block: %d. %s\n", block, err.Error())
		return
	}

	msg := message.NewMessage(watermill.NewUUID(), out)
	var orderingKey string
	if conf.Publisher.EnableMessageOrdering {
		orderingKey = "sync-events"
	}
	if err = pub.Publish(ctx, pubsub.SyncEventsTopic, orderingKey, msg); err != nil {
		log.Printf("can't publish message to pubsub. block: %d. %s\n", block, err.Error())
	}
	log.Println("published block", block, "number of events", len(syncEvents))

	err = s.repo.SetLastRecordedBlock(ctx, block)
	if err != nil {
		log.Printf("can't set last recorded block after publishing messages. block: %d. %s\n", block, err.Error())
	}
}
