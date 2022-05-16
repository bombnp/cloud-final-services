package repository

import (
	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"github.com/bombnp/cloud-final-services/lib/postgres/models"
	"gorm.io/gorm"
)

type Repository struct {
	Postgres *gorm.DB
	InfluxDB *influxdb.Service
}

func New(pg *gorm.DB, influx *influxdb.Service) *Repository {
	return &Repository{
		Postgres: pg,
		InfluxDB: influx,
	}
}

func (db *Repository) InsertNewSubscribe(id, pool, t, channel string) error {

	query := `INSERT INTO pair_subscriptions (server_id,pool_address,type,channel_id) VALUE('?','?','?','?')`
	return db.Postgres.Exec(query, id, pool, t, channel).Error

}

func (db *Repository) QuerySubscribeByAddress(address string) ([]models.PairSubscription, error) {
	query := `SELECT * FROM pair_subscriptions WHERE pool_address = ?`

	var q []models.PairSubscription

	err := db.Postgres.Raw(query, address).Scan(&q).Error

	return q, err
}

func (db *Repository) QueryToken(address string) (models.Token, error) {

	var token models.Token

	query := `SELECT * FROM tokens where address = ?`
	err := db.Postgres.Raw(query, address).First(&token).Error

	return token, err

}

func (db *Repository) QueryAllPairs() ([]models.Pair, error) {
	var pairList []models.Pair

	query := `SELECT * FROM pairs`
	err := db.Postgres.Raw(query).Scan(&pairList).Error

	if err != nil {
		return nil, err
	} else {
		return pairList, nil
	}
}
