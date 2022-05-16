package alert

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bombnp/cloud-final-services/api/config"
	"github.com/bombnp/cloud-final-services/api/repository"
	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"github.com/bombnp/cloud-final-services/lib/pubsub"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v8"
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

func (s *Service) SendAlerts(ctx context.Context, alerts []PriceAlert) error {
	if len(alerts) == 0 {
		return nil
	}
	pairSubMap, err := s.repository.QueryPairSubscriptionsMap()
	if err != nil {
		return errors.Wrap(err, "can't get pair subscriptions map")
	}
	pairNames, err := s.repository.QueryPairNames()
	if err != nil {
		return errors.Wrap(err, "can't get pair subscriptions map")
	}

	var alertMessages []pubsub.PriceAlertMsg

	tm := time.Now().Unix()

	for _, alert := range alerts {

		lastTime, err := s.repository.Redis.Get(ctx, "pair:"+alert.Address.String()).Int64()

		if err != nil && err != redis.Nil {
			return errors.Wrap(err, "redis error")
		}
		if tm-lastTime < 3600 {
			continue
		}

		err = s.repository.Redis.Set(ctx, "pair:"+alert.Address.String(), tm, 0).Err()
		if err != nil {
			return errors.Wrap(err, "redis error")
		}

		pairSubs, ok := pairSubMap[alert.Address]
		if !ok {
			continue
		}
		for _, pairSub := range pairSubs {
			alertMsg := pubsub.PriceAlertMsg{
				ServerId:    pairSub.ServerId,
				PoolAddress: pairSub.PoolAddress,
				ChannelId:   pairSub.ChannelId,
				PairName:    pairNames[alert.Address],
				Change:      alert.Change,
				Since:       alert.Since.Unix(),
			}
			alertMessages = append(alertMessages, alertMsg)
		}
	}
	err = s.publishAlertMessages(ctx, alertMessages)
	if err != nil {
		return errors.Wrap(err, "can't publish alerts")
	}
	return nil
}

func (s *Service) publishAlertMessages(ctx context.Context, alertMessages []pubsub.PriceAlertMsg) error {
	if len(alertMessages) == 0 {
		return nil
	}
	conf := config.InitConfig()
	pub := s.pub
	out, err := json.Marshal(alertMessages)
	if err != nil {
		return errors.Wrap(err, "can't marshal alerts to json")
	}

	msg := message.NewMessage(watermill.NewUUID(), out)
	var orderingKey string
	if conf.Publisher.EnableMessageOrdering {
		orderingKey = pubsub.PriceAlertsTopic
	}
	if err = pub.Publish(ctx, pubsub.PriceAlertsTopic, orderingKey, msg); err != nil {
		return errors.Wrap(err, "can't publish message to pubsub")
	}
	log.Printf("published %d alerts\n", len(alertMessages))
	return nil
}

type PriceSummary struct {
	Close    float64
	High     float64
	HighTime time.Time
	Low      float64
	LowTime  time.Time
}

type PriceAlert struct {
	Address common.Address
	Change  float64
	Since   time.Time
}

const changeThreshold = 0.05

func (s *Service) GetTokenAlerts(ctx context.Context) ([]PriceAlert, error) {
	conf := config.InitConfig()
	currentTime := time.Now()
	args := &alertSummaryTemplateArgs{
		Bucket: conf.Database.InfluxDB.Bucket,
		Start:  currentTime.Add(-6 * time.Hour).Unix(),
		Stop:   currentTime.Unix(),
	}
	result, err := s.repository.InfluxDB.Query(ctx, alertTemplate, args)
	if err != nil {
		return nil, errors.Wrap(err, "error during influx query")
	}

	summaryMap := make(map[common.Address]PriceSummary)
	var currentTable string
	for result.Next() {
		if result.TableChanged() {
			columns := result.TableMetadata().Columns()
			currentTable = columns[0].DefaultValue()
		}
		address := common.HexToAddress(result.Record().ValueByKey("address").(string))
		tm := result.Record().Time()
		summary, ok := summaryMap[address]
		if !ok {
			summary = PriceSummary{}
		}
		switch currentTable {
		case "high":
			summary.High = result.Record().Value().(float64)
			summary.HighTime = tm
		case "low":
			summary.Low = result.Record().Value().(float64)
			summary.LowTime = tm
		case "close":
			summary.Close = result.Record().Value().(float64)
		}
		summaryMap[address] = summary

	}
	var alerts []PriceAlert
	for address, summary := range summaryMap {
		changeLow := (summary.Close - summary.Low) / summary.Low
		changeHigh := (summary.Close - summary.High) / summary.High
		if address == common.HexToAddress("0x05faf555522fa3f93959f86b41a380866609") {
			log.Println(changeLow, changeHigh)
		}
		if summary.HighTime.After(summary.LowTime) {
			if math.Abs(changeHigh) >= changeThreshold {
				alerts = append(alerts, PriceAlert{
					Address: address,
					Change:  changeHigh,
					Since:   summary.HighTime,
				})
			} else if math.Abs(changeLow) >= changeThreshold {
				alerts = append(alerts, PriceAlert{
					Address: address,
					Change:  changeLow,
					Since:   summary.LowTime,
				})
			}
		} else {
			if math.Abs(changeLow) >= changeThreshold {
				alerts = append(alerts, PriceAlert{
					Address: address,
					Change:  changeLow,
					Since:   summary.LowTime,
				})
			} else if math.Abs(changeHigh) >= changeThreshold {
				alerts = append(alerts, PriceAlert{
					Address: address,
					Change:  changeHigh,
					Since:   summary.HighTime,
				})
			}
		}
	}
	return alerts, nil
}

type alertSummaryTemplateArgs struct {
	Bucket string
	Start  int64
	Stop   int64
}

var alertTemplate = influxdb.NewQueryTemplate("prices", `
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
from(bucket: "{{.Bucket}}")
  |> range(start: {{.Start}}, stop: {{.Stop}})
  |> filter(fn: (r) => r["_measurement"] == "price")
  |> filter(fn: (r) => r["_field"] == "price")
  |> last()
  |> yield(name: "close")
`)
