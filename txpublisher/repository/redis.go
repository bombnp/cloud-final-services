package repository

import (
	"context"

	"github.com/pkg/errors"
)

var lastRecordedBlockKey = "txpublisher:lastRecordedBlock"

func (r *Repository) GetLastRecordedBlock(ctx context.Context) (uint64, error) {
	redisClient := r.RedisClient
	block, err := redisClient.Get(ctx, lastRecordedBlockKey).Uint64()
	if err != nil {
		return 0, errors.Wrapf(err, "can't get %s", lastRecordedBlockKey)
	}
	return block, nil
}

func (r *Repository) SetLastRecordedBlock(ctx context.Context, block uint64) error {
	redisClient := r.RedisClient
	err := redisClient.Set(ctx, lastRecordedBlockKey, block, 0).Err()
	if err != nil {
		return errors.Wrapf(err, "can't set %s", lastRecordedBlockKey)
	}
	return nil
}
