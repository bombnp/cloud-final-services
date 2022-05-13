package main

import (
	"context"
	"log"

	"github.com/bombnp/cloud-final-services/lib/postgres"
	"github.com/bombnp/cloud-final-services/lib/pubsub"
	"github.com/bombnp/cloud-final-services/lib/redis"
	"github.com/bombnp/cloud-final-services/txpublisher/config"
	"github.com/bombnp/cloud-final-services/txpublisher/repository"
	"github.com/bombnp/cloud-final-services/txpublisher/worker"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func main() {
	conf := config.InitConfig()

	// Context
	ctx := context.Background()

	// Databases
	pg, err := postgres.New(conf.Postgres)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't connect to postgres").Error())
	}
	log.Println("connected to postgres!")

	// ETH client
	ethClient, err := ethclient.Dial(conf.Chain.NodeURL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't init eth client").Error())
	}
	log.Println("initialized eth client!")

	// Redis
	redisClient, err := redis.New(ctx, conf.Redis)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't connect to redis").Error())
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Println(errors.Wrap(err, "can't close redis client").Error())
		} else {
			log.Println("gracefully stopped redis client")
		}
	}()

	// Pub/Sub
	pub, err := pubsub.NewPublisher(conf.Publisher)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't init google cloud publisher").Error())
	}

	repo := repository.NewRepository(pg, ethClient, redisClient)

	streamer, err := worker.NewStreamer(repo, pub)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't create streamer").Error())
	}
	err = streamer.PollPreviousLogs(ctx)
	if err != nil {
		log.Println(errors.Wrap(err, "error during logs polling").Error())
	}
	err = streamer.LoopConsumeLog(ctx)
	if err != nil {
		log.Println(errors.Wrap(err, "error during logs consumption").Error())
	}
}
