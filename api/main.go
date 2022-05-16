package main

import (
	"fmt"
	"log"

	"github.com/bombnp/cloud-final-services/api/config"
	"github.com/bombnp/cloud-final-services/api/packages/alert"
	"github.com/bombnp/cloud-final-services/api/packages/pair"
	"github.com/bombnp/cloud-final-services/api/packages/subscribe"
	"github.com/bombnp/cloud-final-services/api/repository"
	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"github.com/bombnp/cloud-final-services/lib/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.InitConfig()

	pg, err := postgres.New(&conf.Database.Postgres)
	if err != nil {
		log.Fatalln("Postgres are not connected")
		return
	}

	influx, err := influxdb.NewService(&conf.Database.InfluxDB)
	if err != nil {
		log.Fatalln("Postgres are not connected")
		return
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	repo := repository.New(pg, influx)

	pairHandler := pair.NewHandler(pair.NewService(repo))
	subscribeHandler := subscribe.NewHandler(subscribe.NewService(repo))
	alertHandler := alert.NewHandler(alert.NewService(repo))

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/pair", pairHandler.GetPairs)
		apiGroup.GET("/alert", alertHandler.GetTokenAlertSummaryHandler)

		subscribeGroup := apiGroup.Group("/subscribe")
		{
			subscribeGroup.POST("/alert", subscribeHandler.PostAlertSubscribe)
		}
	}

	log.Printf("Server started on port %d", conf.Server.Port)
	err = router.Run(fmt.Sprintf(":%d", conf.Server.Port))
	if err != nil {
		log.Fatalln("Error running router", err.Error())
	}
}
