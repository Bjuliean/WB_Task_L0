package config

import "os"

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
	
	//var cfg Config

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {

	}

	return nil
}