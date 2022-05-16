package alert

import (
	"context"
	"time"

	"github.com/bombnp/cloud-final-services/api/config"
	"github.com/bombnp/cloud-final-services/api/repository"
	"github.com/bombnp/cloud-final-services/lib/influxdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type Service struct {
	repository *repository.Repository
}

func NewService(db *repository.Repository) *Service {
	return &Service{
		repository: db,
	}
}

func (s *Service) GetTokenAlertSummary(ctx context.Context) (map[common.Address]AlertSummary, error) {
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

	summaryMap := make(map[common.Address]AlertSummary)
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
			summary = AlertSummary{}
		}
		switch currentTable {
		case "high":
			summary.High = result.Record().Value().(float64)
			summary.HighTime = tm
		case "low":
			summary.Low = result.Record().Value().(float64)
			summary.LowTime = tm
		}
		summaryMap[address] = summary

	}
	for address, summary := range summaryMap {
		if summary.HighTime.Before(summary.LowTime) {
			summary.Change = (summary.High - summary.Low) / summary.High
		} else {
			summary.Change = (summary.High - summary.Low) / summary.Low
		}
		summaryMap[address] = summary
	}
	return summaryMap, nil
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
`)
