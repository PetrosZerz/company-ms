package config

import (
	"github.com/spf13/viper"
)

// Config struct holds the application configuration
type Config struct {
	Database MongoDBConfig `mapstructure:"database"`
	Server   ServerConfig  `mapstructure:"server"`
	Kafka    KafkaConfig   `mapstructure:"kafka"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type MongoDBConfig struct {
	Host        string `mapstructure:"host"`
	Environment string `mapstructure:"environment"`
	Name        string `mapstructure:"name"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
}

type KafkaConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
