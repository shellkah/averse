package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Cache  CacheConfig
	Log    LogConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type CacheConfig struct {
	Capacity int
}

type LogConfig struct {
	Level string
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.AutomaticEnv()

	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("log.level", "LOG_LEVEL")
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("cache.capacity", "CACHE_CAPACITY")

	viper.SetDefault("server.port", 50051)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("cache.capacity", 1000)
	viper.SetDefault("log.level", "info")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No configuration file found: %v", err)
	}

	viper.WatchConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
