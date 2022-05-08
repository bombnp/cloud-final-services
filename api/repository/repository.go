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

func (db *Databaser) QueryAllToken() ([]Token, error) {

	var token_list []Token

	query := `SELECT * FROM tokens`
	err := db.Postgres.Raw(query).Scan(&token_list).Error

	return token_list, err

}

func (db *Databaser) QueryToken(address string) (Token, error) {

	var token Token

	query := `SELECT * FROM tokens where address = ?`
	err := db.Postgres.Raw(query, address).First(&token).Error

	return token, err

}
