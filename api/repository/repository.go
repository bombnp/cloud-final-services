package repository

import (
	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"github.com/bombnp/cloud-final-services/lib/postgres/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	Postgres *gorm.DB
	InfluxDB *influxdb.Service
	Redis    *redis.Client
}

func New(pg *gorm.DB, influx *influxdb.Service, rd *redis.Client) *Repository {
	return &Repository{
		Postgres: pg,
		InfluxDB: influx,
		Redis:    rd,
	}
}

func (db *Repository) InsertNewSubscribe(id string, pool string, t models.SubscriptionType, channel string) error {
	pairSub := models.PairSubscription{
		ServerId:    id,
		PoolAddress: pool,
		Type:        t,
		ChannelId:   channel,
	}
	err := db.Postgres.Clauses(clause.OnConflict{UpdateAll: true}).Create(&pairSub).Error
	if err != nil {
		return errors.Wrap(err, "can't execute create query")
	}
	return nil

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

func (db *Repository) QueryPairNames() (map[common.Address]string, error) {
	pairList, err := db.QueryAllPairs()
	if err != nil {
		return nil, errors.Wrap(err, "can't get pairs from postgres")
	}
	pairNames := make(map[common.Address]string)
	for _, pair := range pairList {
		pairNames[common.HexToAddress(pair.PoolAddress)] = pair.Name
	}
	return pairNames, nil
}

func (db *Repository) QueryPairSubscriptionsMap() (map[common.Address][]models.PairSubscription, error) {
	var pairSubs []models.PairSubscription
	if err := db.Postgres.Find(&pairSubs).Error; err != nil {
		return nil, errors.Wrap(err, "can't get pair subscriptions from postgres")
	}

	pairSubMap := make(map[common.Address][]models.PairSubscription)
	for _, pairSub := range pairSubs {
		pairSubList, ok := pairSubMap[common.HexToAddress(pairSub.PoolAddress)]
		if !ok {
			pairSubList = []models.PairSubscription{}
		}
		pairSubList = append(pairSubList, pairSub)
		pairSubMap[common.HexToAddress(pairSub.PoolAddress)] = pairSubList
	}
	return pairSubMap, nil
}
