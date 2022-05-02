package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	PG *gorm.DB
}

func NewRepository(pg *gorm.DB) *Repository {
	return &Repository{
		PG: pg,
	}
}
