package repository

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Repository struct {
	PG          *gorm.DB
	EthClient   *ethclient.Client
	RedisClient *redis.Client
}

func NewRepository(pg *gorm.DB, ethClient *ethclient.Client, redisClient *redis.Client) *Repository {
	return &Repository{
		PG:          pg,
		EthClient:   ethClient,
		RedisClient: redisClient,
	}
}
