package subscribe

import (
	"github.com/bombnp/cloud-final-services/api/repository"
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
	return s.repository.InsertNewSubscribe(id, pool, "alert", channel)
}
