package main

import (
	"fmt"
	"log"

	"github.com/bombnp/cloud-final-services/api/config"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.InitConfig()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
		return
	})

	log.Printf("Server started on port %d", conf.Server.Port)
	err := router.Run(fmt.Sprintf(":%d", conf.Server.Port))
	if err != nil {
		log.Fatalln("Error running router", err.Error())
	}
}
