package influxdb

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"
)

// Config defines InfluxDB config
type Config struct {
	URL          string `mapstructure:"url"`
	Token        string `mapstructure:"token"`
	Bucket       string `mapstructure:"bucket"`
	Organization string `mapstructure:"organization"`
}
type Service struct {
	config           *Config
	Client           influxdb2.Client
	queryAPI         api.QueryAPI
	WriteAPI         api.WriteAPI
	WriteAPIBlocking api.WriteAPIBlocking
}

func NewService(config *Config) (*Service, error) {
	if config.URL == "" {
		return nil, errors.New("URL is required")
	}
	if config.Token == "" {
		return nil, errors.New("Token is required")
	}
	if config.Bucket == "" {
		return nil, errors.New("Bucket is required")
	}
	if config.Organization == "" {
		return nil, errors.New("Organization is required")
	}

	client := influxdb2.NewClientWithOptions(config.URL, config.Token, influxdb2.DefaultOptions().
		SetFlushInterval(1000))
	queryAPI := client.QueryAPI(config.Organization)
	writeAPI := client.WriteAPI(config.Organization, config.Bucket)
	writeAPIBlocking := client.WriteAPIBlocking(config.Organization, config.Bucket)

	return &Service{
		config:           config,
		Client:           client,
		queryAPI:         queryAPI,
		WriteAPI:         writeAPI,
		WriteAPIBlocking: writeAPIBlocking,
	}, nil
}

func (s *Service) Query(ctx context.Context, template *QueryTemplate, args any) (*api.QueryTableResult, error) {
	query, err := template.GetQueryString(args)
	if err != nil {
		return nil, errors.Wrap(err, "can't get query string")
	}
	result, err := s.queryAPI.Query(ctx, query)
	if err != nil {
		err = errors.WithStack(err)
		return nil, errors.Wrap(err, "influxdb query failed")
	}
	return result, nil
}
