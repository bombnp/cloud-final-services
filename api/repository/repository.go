package repository

import "gorm.io/gorm"

type Databaser struct {
	Postgres *gorm.DB
}

func New(pg *gorm.DB) *Databaser {
	return &Databaser{
		Postgres: pg,
	}
}

func (db *Databaser) InsertNewSubscribe(id, pool, t string) error {
	query := `INSERT INTO pair_subscriptions (server_id,pool_address,type) VALUE('?','?','?')`
	return db.Postgres.Exec(query, id, pool, t).Error
}
