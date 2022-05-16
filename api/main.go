package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bombnp/cloud-final-services/api/config"
	"github.com/bombnp/cloud-final-services/api/packages/alert"
	"github.com/bombnp/cloud-final-services/api/packages/pair"
	"github.com/bombnp/cloud-final-services/api/packages/subscribe"
	"github.com/bombnp/cloud-final-services/api/repository"
	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"github.com/bombnp/cloud-final-services/lib/postgres"
	"github.com/bombnp/cloud-final-services/lib/pubsub"
	"github.com/bombnp/cloud-final-services/lib/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.InitConfig()

	pg, err := postgres.New(&conf.Database.Postgres)
	if err != nil {
		log.Fatalln("Postgres are not connected", err)
		return
	}

	influx, err := influxdb.NewService(&conf.Database.InfluxDB)
	if err != nil {
		log.Fatalln("Postgres are not connected", err)
		return
	}

	rd, err := redis.New(context.Background(), &conf.Database.Redis)
	if err != nil {
		log.Fatalln("Redis are not connected", err)
		return
	}

	pub, err := pubsub.NewPublisher(conf.Publisher)
	if err != nil {
		log.Fatalln("can't init google cloud publisher", err)
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	repo := repository.New(pg, influx, rd)

	pairHandler := pair.NewHandler(pair.NewService(repo))
	subscribeHandler := subscribe.NewHandler(subscribe.NewService(repo))
	alertHandler := alert.NewHandler(alert.NewService(repo, pub))

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/pair", pairHandler.GetPairs)

		subscribeGroup := apiGroup.Group("/subscribe")
		{
			subscribeGroup.GET("/alert", subscribeHandler.GetAlertSubscribe)
			subscribeGroup.POST("/alert", subscribeHandler.PostAlertSubscribe)
		}

		triggerGroup := apiGroup.Group("/trigger")
		{
			triggerGroup.POST("/alert", alertHandler.TriggerPriceAlert)
		}
	}

	log.Printf("Server started on port %d", conf.Server.Port)
	err = router.Run(fmt.Sprintf(":%d", conf.Server.Port))
	if err != nil {
		log.Fatalln("Error running router", err.Error())
	}
}
