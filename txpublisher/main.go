package main

import (
	"log"

	"github.com/bombnp/cloud-final-services/lib/postgres"
	"github.com/bombnp/cloud-final-services/txpublisher/config"
	"github.com/bombnp/cloud-final-services/txpublisher/repository"
)

func main() {
	conf := config.InitConfig()

	// Databases
	pg, err := postgres.New(conf.Postgres)
	if err != nil {
		log.Fatal("can't connect to postgres", err.Error())
	}
	log.Println("connected to postgres!")

	repo := repository.NewRepository(pg)

	log.Println(repo)
}
