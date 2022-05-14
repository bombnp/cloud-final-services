package repository

import (
	"log"

	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"github.com/bombnp/cloud-final-services/lib/postgres"
	"github.com/bombnp/cloud-final-services/pricecollector/config"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	PG     *gorm.DB
	Influx *influxdb.Service
}

func NewRepository() (*Repository, error) {
	conf := config.InitConfig()

	// Databases
	pg, err := postgres.New(conf.Postgres)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to postgres")
	}
	log.Println("connected to postgres!")

	influx, err := influxdb.NewService(conf.Influx)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to influxdb")
	}
	log.Println("connected to influxdb!")

	return &Repository{
		PG:     pg,
		Influx: influx,
	}, nil
}
