package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func New(config *Config) (db *gorm.DB, err error) {
	if config.Host == "" {
		return nil, errors.New("Host is required")
	}
	if config.Port == "" {
		return nil, errors.New("Port is required")
	}
	if config.User == "" {
		return nil, errors.New("User is required")
	}
	if config.Password == "" {
		return nil, errors.New("Password is required")
	}
	if config.DBName == "" {
		return nil, errors.New("DBName is required")
	}
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize db session")
	}
	return db, nil
}
