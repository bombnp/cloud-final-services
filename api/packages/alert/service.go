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

type alertSummary struct {
	High   float64
	Low    float64
	Change float64
}

func (s *Service) GetTokenAlertSummary(ctx context.Context) (map[common.Address]alertSummary, error) {
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

	summaryMap := make(map[common.Address]alertSummary)
	var currentTable string
	for result.Next() {
		if result.TableChanged() {
			columns := result.TableMetadata().Columns()
			currentTable = columns[0].DefaultValue()
		}
		address := common.HexToAddress(result.Record().ValueByKey("address").(string))
		summary, ok := summaryMap[address]
		if !ok {
			summary = alertSummary{}
		}
		switch currentTable {
		case "high":
			summary.High = result.Record().Value().(float64)
		case "low":
			summary.Low = result.Record().Value().(float64)
		}
		summaryMap[address] = summary
	}
	for address, summary := range summaryMap {
		summary.Change = (summary.High - summary.Low) / summary.Low
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
