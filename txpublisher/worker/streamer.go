package worker

import (
	"context"
	"log"

	"github.com/bombnp/cloud-final-services/lib/ethutils"
	"github.com/bombnp/cloud-final-services/lib/postgres/models"
	"github.com/bombnp/cloud-final-services/txpublisher/repository"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type streamer struct {
	repo *repository.Repository

	pairs      []models.Pair
	logCh      chan ethtypes.Log
	startBlock uint64
}

type Streamer interface {
	PollPreviousLogs() error
}

func NewStreamer(repo *repository.Repository) (Streamer, error) {
	pairs, err := repo.GetPairs()
	if err != nil {
		return nil, errors.Wrap(err, "can't get pairs during streamer init")
	}
	s := &streamer{
		repo:  repo,
		pairs: pairs,
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
	log.Println("Subscribed to chain at block number:", startBlock)
	_, err = ethClient.SubscribeFilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: nil,
		Topics:    [][]common.Hash{ethutils.GetSyncTopics()},
	}, logCh)
	if err != nil {
		return errors.Wrap(err, "can't subscribe filter logs")
	}
	log.Println("Subscribed filter logs")
	s.logCh = logCh
	return nil
}
