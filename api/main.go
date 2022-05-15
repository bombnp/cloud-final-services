package main

import (
	"fmt"
	"log"

	"github.com/bombnp/cloud-final-services/api/config"
	"github.com/bombnp/cloud-final-services/api/repository"
	"github.com/bombnp/cloud-final-services/api/services"
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
		return
	})

	service := services.NewHandler(services.NewService(repository.New(pg, influx)))
	api_handler := router.Group("/api")
	{
		api_handler.GET("/pair", service.GetAllPairHandler)

		subscribe_handler := api_handler.Group("/subscribe")
		{
			subscribe_handler.POST("/alert", service.AlertSubscribeHandler)
		}
	}

	log.Printf("Server started on port %d", conf.Server.Port)
	err = router.Run(fmt.Sprintf(":%d", conf.Server.Port))
	if err != nil {
		log.Fatalln("Error running router", err.Error())
	}
}
