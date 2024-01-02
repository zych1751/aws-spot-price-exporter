package internal

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Regions           []string
	IncludeAllRegions bool
	IncludeAllZones   bool
	MapZoneId         bool
	MetricName        string
	Port              int
}

var defaultConfig = Config{
	Regions:           []string{},
	IncludeAllRegions: false,
	IncludeAllZones:   false,
	MapZoneId:         true,
	MetricName:        "aws_spot_price",
	Port:              8080,
}

func ReadConfigOrDie() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/")
	viper.AddConfigPath("/etc/spot-price-exporter/")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("SPOT")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ignore
		} else {
			log.Fatalf("unable to read config, %v", err)
		}
	}

	cfg := defaultConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to unmarshal config, %v", err)
	}
	log.Printf("config: %+v", cfg)
	return &cfg
}
