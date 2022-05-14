package repository

import (
	"github.com/bombnp/cloud-final-services/lib/postgres/models"
	"github.com/pkg/errors"
)

func (r *Repository) GetPairs() ([]models.Pair, error) {
	var pairs []models.Pair
	if err := r.PG.Find(&pairs).Error; err != nil {
		return nil, errors.Wrap(err, "can't get pairs from postgres")
	}
	return pairs, nil
}
