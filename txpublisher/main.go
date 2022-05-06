package main

import (
	"context"
	"log"

	"github.com/bombnp/cloud-final-services/lib/postgres"
	"github.com/bombnp/cloud-final-services/lib/redis"
	"github.com/bombnp/cloud-final-services/txpublisher/config"
	"github.com/bombnp/cloud-final-services/txpublisher/repository"
	"github.com/bombnp/cloud-final-services/txpublisher/worker"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	conf := config.InitConfig()

	// Context
	ctx := context.Background()

	// Databases
	pg, err := postgres.New(conf.Postgres)
	if err != nil {
		log.Fatal("can't connect to postgres", err.Error())
	}
	log.Println("connected to postgres!")

	// ETH client
	ethClient, err := ethclient.Dial(conf.Chain.NodeURL)
	if err != nil {
		log.Fatal("can't init eth client", err.Error())
	}
	log.Println("initialized eth client!")

	// Redis
	redisClient, err := redis.New(ctx, conf.Redis)
	if err != nil {
		log.Fatal("can't connect to redis", err.Error())
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Println("can't close redis client", err.Error())
		} else {
			log.Println("gracefully stopped redis client")
		}
	}()

	repo := repository.NewRepository(pg, ethClient, redisClient)

	streamer, err := worker.NewStreamer(repo)
	if err != nil {
		log.Fatal("can't create streamer", err.Error())
	}
}
