package pair

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

func (s *Service) GetPairs() ([]GetPairsResponse, error) {

	query, err := s.repository.QueryAllPairs()

	if err != nil {
		return nil, err
	}

	var pairList []GetPairsResponse

	for _, e := range query {

		pairList = append(pairList, GetPairsResponse{
			PoolAddress: e.PoolAddress,
			PoolName:    e.Name,
		})

	}

	return pairList, nil
}
