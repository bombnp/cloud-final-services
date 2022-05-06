package config

import (
	"log"
	"strings"
	"sync"

	"github.com/bombnp/cloud-final-services/lib/postgres"
	"github.com/bombnp/cloud-final-services/lib/redis"
	"github.com/spf13/viper"
)

var configOnce sync.Once
var config *Config

type Config struct {
	Postgres *postgres.Config `mapstructure:"postgres"`
	Redis    *redis.Config    `mapstructure:"redis"`
	Chain    Chain            `mapstructure:"chain"`
}

type Chain struct {
	NodeURL           string `mapstructure:"node_url"`
	MaxBlocksPerQuery uint64 `mapstructure:"max_blocks_per_query"`
}

func InitConfig() *Config {
	configOnce.Do(func() {

		viper.SetConfigName("config")   // name of config file without extension
		viper.AddConfigPath("./config") // path to look for config file

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Config file not found", err.Error())
		}
		viper.AutomaticEnv()

		viper.WatchConfig() // Watch for changes to the configuration file and recompile
		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
		log.Println("Config initialized!")
	})
	return config
}
