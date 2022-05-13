package worker

import (
	"context"
	"log"
	"math/big"

	"github.com/bombnp/cloud-final-services/lib/ethutils"
	"github.com/bombnp/cloud-final-services/txpublisher/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const BscBlockTime = 3

func (s *streamer) PollPreviousLogs(ctx context.Context) error {
	conf := config.InitConfig()
	startBlock := s.startBlock
	ethClient := s.repo.EthClient

	var maxBlocksPerQuery = conf.Chain.MaxBlocksPerQuery

	log.Printf("streamer_poll: Starting transactions poll. Target block: %d, maxBlocksPerQuery: %d\n", startBlock, maxBlocksPerQuery)

	// get last recorded blocks
	lastRecordedBlock, err := s.repo.GetLastRecordedBlock(ctx)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return errors.Wrap(err, "can't get last recorded block")
		}
		log.Println("streamer_poll: Last recorded block not found: starting from last 24h")
		lastRecordedBlock = startBlock - 24*60*60/BscBlockTime
	}

	// start polling
	for lastBlock := lastRecordedBlock; lastBlock < startBlock; {
		fromBlock := lastBlock + 1
		toBlock := lastBlock + maxBlocksPerQuery
		// bound toBlock so it doesn't exceed startBlock
		if toBlock > startBlock {
			toBlock = startBlock
		}
		logs, err := ethClient.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(toBlock)),
			Topics:    [][]common.Hash{ethutils.GetSyncTopics()},
			Addresses: s.pairAddresses,
		})
		log.Printf("streamer_poll: Polled %d blocks %d - %d. %d logs", toBlock-fromBlock+1, fromBlock, toBlock, len(logs))
		if err != nil {
			return errors.Wrapf(err, "can't poll logs from block %d to block %d", fromBlock, toBlock)
		}
		for _, eventLog := range logs {
			err := s.ProcessLog(ctx, eventLog)
			if err != nil {
				log.Println(errors.Wrap(err, "can't process log").Error())
			}
		}
		lastBlock = toBlock
	}

	return nil
}
