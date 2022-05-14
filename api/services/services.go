package services

import (
	"github.com/bombnp/cloud-final-services/api/repository"
)

type Service struct {
	Database Databaser
}

type Databaser interface {
	InsertNewSubscribe(id, pool, t, channel string) error
	QueryToken(address string) (repository.Token, error)
	QueryAllPair() ([]repository.Pair, error)
}

func NewService(db Databaser) *Service {
	return &Service{
		Database: db,
	}
}

func (s *Service) AlertSubscribe(id string, pool string, channel string) error {
	return s.Database.InsertNewSubscribe(id, pool, "alert", channel)
}

func (s *Service) GetAllPair() ([]PairResponse, error) {

	query, err := s.Database.QueryAllPair()

	if err != nil {
		return nil, err
	}

	var pair_list []PairResponse

	for _, e := range query {

		pair_list = append(pair_list, PairResponse{
			PoolAddress: e.PoolAddress,
			PoolName:    e.Name,
		})

	}

	return pair_list, nil
}
