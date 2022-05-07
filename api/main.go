package main

import (
	"fmt"
	"log"

	"github.com/bombnp/cloud-final-services/api/config"
	"github.com/bombnp/cloud-final-services/api/services"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.InitConfig()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
		return
	})

	service := services.NewHandler(services.NewService())
	api_handler := router.Group("/api")
	subscribe_handler := api_handler.Group("/subscribe")
	{
		subscribe_handler.POST("/alert", service.AlertSubscribeHandler)
	}

	log.Printf("Server started on port %d", conf.Server.Port)
	err := router.Run(fmt.Sprintf(":%d", conf.Server.Port))
	if err != nil {
		log.Fatalln("Error running router", err.Error())
	}
}
