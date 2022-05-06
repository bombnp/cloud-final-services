package redis

import (
	"context"

	redis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// Config defines redis config
type Config struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"` // optional
}

// New creates new redis instance
func New(ctx context.Context, config *Config) (*redis.Client, error) {
	if config.Address == "" {
		config.Address = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Database,
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot connect to redis")
	}
	return redisClient, nil
}
