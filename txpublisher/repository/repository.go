package repository

import (
	"context"
	"log"

	"github.com/bombnp/cloud-final-services/lib/postgres"
	libredis "github.com/bombnp/cloud-final-services/lib/redis"
	"github.com/bombnp/cloud-final-services/txpublisher/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	PG          *gorm.DB
	EthClient   *ethclient.Client
	RedisClient *redis.Client
}

func NewRepository(ctx context.Context) (*Repository, error) {
	conf := config.InitConfig()

	// Databases
	pg, err := postgres.New(conf.Postgres)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}
	log.Println("connected to postgres!")

	// ETH client
	ethClient, err := ethclient.Dial(conf.Chain.NodeURL)
	if err != nil {
		return nil, errors.Wrap(err, "can't init eth client")
	}
	log.Println("initialized eth client!")

	// Redis
	redisClient, err := libredis.New(ctx, conf.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to redis")
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Println(errors.Wrap(err, "can't close redis client").Error())
		} else {
			log.Println("gracefully stopped redis client")
		}
	}()
	log.Println("connected to redis!")
	return &Repository{
		PG:          pg,
		EthClient:   ethClient,
		RedisClient: redisClient,
	}, nil
}
