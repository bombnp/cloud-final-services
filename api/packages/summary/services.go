package summary

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bombnp/cloud-final-services/api/config"
	"github.com/bombnp/cloud-final-services/api/repository"
	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"github.com/bombnp/cloud-final-services/lib/pubsub"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type Service struct {
	repository *repository.Repository
	pub        *pubsub.Publisher
}

func NewService(repository *repository.Repository, pub *pubsub.Publisher) *Service {
	return &Service{
		repository: repository,
		pub:        pub,
	}
}

type dailySummary struct {
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Change float64
}

func (s *Service) SendSummaryReports(ctx context.Context, summaryMap map[common.Address]dailySummary) error {
	pairSubMap, err := s.repository.QueryPairSubscriptionsMap()
	if err != nil {
		return errors.Wrap(err, "can't get pair subscriptions map")
	}
	pairNames, err := s.repository.QueryPairNames()
	if err != nil {
		return errors.Wrap(err, "can't get pair subscriptions map")
	}

	var summaryMessages []pubsub.PriceSummaryMsg
	for address, summary := range summaryMap {
		pairSubs, ok := pairSubMap[address]
		if !ok {
			continue
		}
		for _, pairSub := range pairSubs {
			summaryMsg := pubsub.PriceSummaryMsg{
				ServerId:    pairSub.ServerId,
				PoolAddress: pairSub.PoolAddress,
				ChannelId:   pairSub.ChannelId,
				PairName:    pairNames[address],
				Date:        time.Now().Add(-24 * time.Hour).Format("January 02, 2006"),
				Open:        summary.Open,
				Close:       summary.Close,
				High:        summary.High,
				Low:         summary.Low,
				Change:      summary.Change,
			}
			summaryMessages = append(summaryMessages, summaryMsg)
		}
	}
	err = s.publishSummaryMessages(ctx, summaryMessages)
	if err != nil {
		return errors.Wrap(err, "can't publish summary messages")
	}
	return nil

}

func (s *Service) publishSummaryMessages(ctx context.Context, summaryMessages []pubsub.PriceSummaryMsg) error {
	conf := config.InitConfig()
	pub := s.pub
	out, err := json.Marshal(summaryMessages)
	if err != nil {
		return errors.Wrap(err, "can't marshal summary reports to json")
	}

	msg := message.NewMessage(watermill.NewUUID(), out)
	var orderingKey string
	if conf.Publisher.EnableMessageOrdering {
		orderingKey = pubsub.PriceSummaryTopic
	}
	if err = pub.Publish(ctx, pubsub.PriceSummaryTopic, orderingKey, msg); err != nil {
		return errors.Wrap(err, "can't publish message to pubsub")
	}
	log.Printf("published %d summary reports\n", len(summaryMessages))
	return nil
}

func (s *Service) GetTokenDailySummary(ctx context.Context) (map[common.Address]dailySummary, error) {
	conf := config.InitConfig()
	currentTime := time.Now()
	args := &dailySummaryTemplateArgs{
		Bucket: conf.Database.InfluxDB.Bucket,
		Start:  currentTime.Add(-24 * time.Hour).Unix(),
		Stop:   currentTime.Unix(),
	}
	result, err := s.repository.InfluxDB.Query(ctx, dailySummaryTemplate, args)
	if err != nil {
		return nil, errors.Wrap(err, "error during influx query")
	}

	summaryMap := make(map[common.Address]dailySummary)
	var currentTable string
	for result.Next() {
		if result.TableChanged() {
			columns := result.TableMetadata().Columns()
			currentTable = columns[0].DefaultValue()
		}
		address := common.HexToAddress(result.Record().ValueByKey("address").(string))
		summary, ok := summaryMap[address]
		if !ok {
			summary = dailySummary{}
		}
		switch currentTable {
		case "open":
			summary.Open = result.Record().Value().(float64)
		case "close":
			summary.Close = result.Record().Value().(float64)
		case "high":
			summary.High = result.Record().Value().(float64)
		case "low":
			summary.Low = result.Record().Value().(float64)
		}
		summaryMap[address] = summary
	}
	for address, summary := range summaryMap {
		summary.Change = (summary.Close - summary.Open) / summary.Open
		summaryMap[address] = summary
	}
	return summaryMap, nil
}

type dailySummaryTemplateArgs struct {
	Bucket string
	Start  int64
	Stop   int64
}

var dailySummaryTemplate = influxdb.NewQueryTemplate("prices", `
from(bucket: "{{.Bucket}}")
  |> range(start: {{.Start}}, stop: {{.Stop}})
  |> filter(fn: (r) => r["_measurement"] == "price")
  |> filter(fn: (r) => r["_field"] == "price")
  |> first()
  |> yield(name: "open")

from(bucket: "{{.Bucket}}")
  |> range(start: {{.Start}}, stop: {{.Stop}})
  |> filter(fn: (r) => r["_measurement"] == "price")
  |> filter(fn: (r) => r["_field"] == "price")
  |> last()
  |> yield(name: "close")

from(bucket: "{{.Bucket}}")
  |> range(start: {{.Start}}, stop: {{.Stop}})
  |> filter(fn: (r) => r["_measurement"] == "price")
  |> filter(fn: (r) => r["_field"] == "price")
  |> max()
  |> yield(name: "high")

from(bucket: "{{.Bucket}}")
  |> range(start: {{.Start}}, stop: {{.Stop}})
  |> filter(fn: (r) => r["_measurement"] == "price")
  |> filter(fn: (r) => r["_field"] == "price")
  |> min()
  |> yield(name: "low")
`)
