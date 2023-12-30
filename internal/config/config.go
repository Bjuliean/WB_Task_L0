package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres      PostgresConfig      `yaml:"postgres"`
	NatsStreaming NatsStreamingConfig `yaml:"nats_streaming"`
	Server        ServerConfig        `yaml:"server"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type NatsStreamingConfig struct {
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	ClientID         string `yaml:"client_id"`
	ClusterID        string `yaml:"cluster_id"`
	SubscribeSubject string `yaml:"subscribe_subject"`
}

type ServerConfig struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func New() *Config {
	const ferr = "internal.config.New"

	var cfg Config

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatalf("%s: CONFIG_PATH is not exists", ferr)
	}

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("%s: error while reading config: %s", ferr, err.Error())
	}

	return &cfg
}
