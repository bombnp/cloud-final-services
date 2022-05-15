package repository

import (
	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"gorm.io/gorm"
)

type Databaser struct {
	Postgres *gorm.DB
	InfluxDB *influxdb.Service
}

func New(pg *gorm.DB, influx *influxdb.Service) *Databaser {
	return &Databaser{
		Postgres: pg,
		InfluxDB: influx,
	}
}

func (db *Databaser) InsertNewSubscribe(id, pool, t, channel string) error {

	query := `INSERT INTO pair_subscriptions (server_id,pool_address,type,channel_id) VALUE('?','?','?','?')`
	return db.Postgres.Exec(query, id, pool, t, channel).Error

}

func (db *Databaser) QueryToken(address string) (Token, error) {

	var token Token

	query := `SELECT * FROM tokens where address = ?`
	err := db.Postgres.Raw(query, address).First(&token).Error

	return token, err

}

func (db *Databaser) QueryAllPair() ([]Pair, error) {
	var pair_list []Pair

	query := `SELECT * FROM pairs`
	err := db.Postgres.Raw(query).Scan(pair_list).Error

	if err != nil {
		return nil, err
	} else {
		return pair_list, nil
	}
}
