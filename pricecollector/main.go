package main

import (
	"context"
	"log"

	"github.com/bombnp/cloud-final-services/pricecollector/worker"
	"github.com/pkg/errors"
)

func main() {
	// Context
	ctx := context.Background()

	collector, err := worker.NewCollector(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't initialize collector").Error())
	}

	err = collector.LoopCollectEvents(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "collector stopped").Error())
	}
}
