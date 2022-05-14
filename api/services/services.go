package services

import (
	"fmt"

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

		base, err := s.Database.QueryToken(e.BaseAddress)

		if err != nil {
			return nil, err
		}

		quote, err := s.Database.QueryToken(e.QuoteAddress)

		if err != nil {
			return nil, err
		}

		pair_list = append(pair_list, PairResponse{
			PoolAddress: e.PoolAddress,
			PoolName:    fmt.Sprintf("%s/%s", base.Symbol, quote.Symbol),
		})

	}

	return pair_list, nil
}
