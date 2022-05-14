package ethutils

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Fetcher interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethtypes.Header, error)
}

type BlockTimestampCache struct {
	rdb     *redis.Client
	fetcher Fetcher

	blockTimestamps sync.Map
}

func redisCacheKeyFrom(number uint64) string {
	return fmt.Sprintf("cache:blockTimestamp:%d", number)
}

func NewBlockTimestampCache(rdb *redis.Client, fetcher Fetcher) *BlockTimestampCache {
	return &BlockTimestampCache{
		rdb:     rdb,
		fetcher: fetcher,
	}
}

func (c *BlockTimestampCache) GetTimestamp(ctx context.Context, blockNumber uint64) (uint64, error) {
	// try fetching from memory first
	timestampFromMap, ok := c.blockTimestamps.Load(blockNumber)
	if ok {
		return timestampFromMap.(uint64), nil
	}

	timestamp, err := c.rdb.Get(ctx, redisCacheKeyFrom(blockNumber)).Uint64()
	if err != nil {
		// not in redis or redis encountered error, try fetching from node
		if !errors.Is(err, redis.Nil) {
			log.Println("error fetching block timestamp from redis", err)
		}
		timestamp, err := c.FetchAndCacheTimestamp(ctx, blockNumber)
		if err != nil {
			return 0, errors.Wrap(err, "error while fetching block timestamp from node")
		}
		return timestamp, nil
	}
	c.blockTimestamps.Store(blockNumber, timestamp)
	return timestamp, nil
}

func (c *BlockTimestampCache) FetchAndCacheTimestamp(ctx context.Context, number uint64) (uint64, error) {
	header, err := c.fetcher.HeaderByNumber(ctx, new(big.Int).SetUint64(number))
	if err != nil {
		return 0, errors.WithStack(err)
	}
	timestamp := header.Time
	c.blockTimestamps.Store(number, timestamp)
	err = c.rdb.Set(ctx, redisCacheKeyFrom(number), timestamp, 24*time.Hour).Err()
	if err != nil {
		log.Println("failed to set block timestamp cache", err)
	}
	return timestamp, nil
}
