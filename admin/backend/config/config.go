package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig ServerConfig `mapstructure:"server"`
	JWTConfig    JWTConfig    `mapstructure:"jwt"`
	MongoConfig  MongoConfig  `mapstructure:"mongodb"`
	RedisConfig  RedisConfig  `mapstructure:"redis"`
	LogConfig    LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type MongoConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

var Conf = new(Config)

func Init() error {
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(Conf); err != nil {
		return err
	}
	return nil
}
