package subscribe

import (
	"github.com/bombnp/cloud-final-services/api/repository"
	"github.com/bombnp/cloud-final-services/lib/postgres/models"
)

type Service struct {
	repository *repository.Repository
}

func NewService(db *repository.Repository) *Service {
	return &Service{
		repository: db,
	}
}

func (s *Service) AlertSubscribe(id string, pool string, channel string) error {
	return s.repository.InsertNewSubscribe(id, pool, models.AlertSubscription, channel)
}

func (s *Service) GetAlert(address string) ([]models.PairSubscription, error) {
	return s.repository.QuerySubscribeByAddress(address)
}
