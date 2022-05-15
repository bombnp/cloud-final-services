package services

import (
	"github.com/bombnp/cloud-final-services/api/repository"
)

type Service struct {
	Database *repository.Repository
}

func NewService(db *repository.Repository) *Service {
	return &Service{
		Database: db,
	}
}

func (s *Service) AlertSubscribe(id string, pool string, channel string) error {
	return s.Database.InsertNewSubscribe(id, pool, "alert", channel)
}

func (s *Service) GetAllPair() ([]PairResponse, error) {

	query, err := s.Database.QueryAllPairs()

	if err != nil {
		return nil, err
	}

	var pairList []PairResponse

	for _, e := range query {

		pairList = append(pairList, PairResponse{
			PoolAddress: e.PoolAddress,
			PoolName:    e.Name,
		})

	}

	return pairList, nil
}
