package main

import (
	"context"
	"log"

	"github.com/bombnp/cloud-final-services/txpublisher/worker"
	"github.com/pkg/errors"
)

func main() {
	// Context
	ctx := context.Background()

	streamer, err := worker.NewStreamer(ctx)
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
