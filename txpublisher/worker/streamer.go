package worker

import (
	"context"
	"log"

	"github.com/bombnp/cloud-final-services/lib/ethutils"
	"github.com/bombnp/cloud-final-services/lib/pubsub"
	"github.com/bombnp/cloud-final-services/txpublisher/config"
	"github.com/bombnp/cloud-final-services/txpublisher/repository"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type streamer struct {
	repo          *repository.Repository
	pub           *pubsub.Publisher
	pairAddresses []common.Address

	logCh      chan ethtypes.Log
	startBlock uint64
}

type Streamer interface {
	PollPreviousLogs(ctx context.Context) error
	LoopConsumeLog(ctx context.Context) error
}

func NewStreamer(ctx context.Context) (Streamer, error) {
	conf := config.InitConfig()
	repo, err := repository.NewRepository(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't create repository")
	}
	pub, err := pubsub.NewPublisher(conf.Publisher)
	if err != nil {
		return nil, errors.Wrap(err, "can't init google cloud publisher")
	}

	pairs, err := repo.GetPairs()
	if err != nil {
		return nil, errors.Wrap(err, "can't get pairs during streamer init")
	}
	var pairAddresses []common.Address
	for _, pair := range pairs {
		pairAddresses = append(pairAddresses, common.HexToAddress(pair.PoolAddress))
	}

	s := &streamer{
		repo:          repo,
		pub:           pub,
		pairAddresses: pairAddresses,
	}
	err = s.subscribeLogs()
	if err != nil {
		return nil, errors.Wrap(err, "can't subscribe logs")
	}
	return s, nil
}

func (s *streamer) subscribeLogs() error {
	ctx := context.Background()
	ethClient := s.repo.EthClient

	logCh := make(chan ethtypes.Log)
	startBlock, err := ethClient.BlockNumber(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get block number")
	}

	s.startBlock = startBlock
	topics := [][]common.Hash{ethutils.GetSyncTopics()}
	_, err = ethClient.SubscribeFilterLogs(ctx, ethereum.FilterQuery{
		Topics:    topics,
		Addresses: s.pairAddresses,
	}, logCh)
	if err != nil {
		return errors.Wrap(err, "can't subscribe filter logs")
	}
	log.Printf("Subscribed to chain at block number: %d\n", startBlock)
	log.Printf("Subscribed filter logs, with %d pairs and topics: %s", len(s.pairAddresses), topics)
	s.logCh = logCh
	return nil
}
