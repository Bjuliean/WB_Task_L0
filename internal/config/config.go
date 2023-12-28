package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres	PostgresConfig	`yaml:"postgres"`
}

type PostgresConfig struct {
	Host	  string		`yaml:"host"`
	Port	  string		`yaml:"port"`
	User	  string		`yaml:"user"`
	Password  string		`yaml:"password"`
}

func New() *Config {
	const ferr = "internal.config.New"
	
	var cfg Config

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatalf("%s: CONFIG_PATH is not exists")
	}

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("%s: error while reading config: %s", ferr, err.Error())
	}

	return &cfg
}