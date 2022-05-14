package worker

import (
	"context"

	"github.com/bombnp/cloud-final-services/lib/postgres/models"
	"github.com/bombnp/cloud-final-services/lib/pubsub"
	"github.com/bombnp/cloud-final-services/pricecollector/config"
	"github.com/bombnp/cloud-final-services/pricecollector/repository"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type collector struct {
	repo *repository.Repository
	sub  *pubsub.Subscriber

	pairMap map[common.Address]models.Pair
}

type Collector interface {
	LoopCollectEvents(ctx context.Context) error
}

func NewCollector(ctx context.Context) (Collector, error) {
	conf := config.InitConfig()
	repo, err := repository.NewRepository()
	if err != nil {
		return nil, errors.Wrap(err, "can't create repository")
	}
	sub, err := pubsub.NewSubscriber(conf.Subscriber)
	if err != nil {
		return nil, errors.Wrap(err, "can't init google cloud publisher")
	}

	pairs, err := repo.GetPairs()
	if err != nil {
		return nil, errors.Wrap(err, "can't get pairs during collector init")
	}
	pairMap := make(map[common.Address]models.Pair)
	for _, pair := range pairs {
		pairMap[common.HexToAddress(pair.PoolAddress)] = pair
	}

	c := &collector{
		repo:    repo,
		sub:     sub,
		pairMap: pairMap,
	}
	return c, nil
}
