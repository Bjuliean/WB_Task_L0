// Получает данные с config/local.yml и загружает их
// в файл .env, который затем использует docker-compose.
// Позволяет задавать конфигурацию, которая используется
// и в главном приложении и docker-compose из одного файла - local.yml

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres      PostgresConfig      `yaml:"postgres"`
	NatsStreaming NatsStreamingConfig `yaml:"nats_streaming"`
}

type PostgresConfig struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	DBName        string `yaml:"dbname"`
	ContainerName string `yaml:"containername"`
}

type NatsStreamingConfig struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	ClientID      string `yaml:"client_id"`
	ClusterID     string `yaml:"cluster_id"`
	ContainerName string `yaml:"containername"`
}

func main() {

	const ferr = "cmd.config.main"

	cfg := createCfg()

	text := fmt.Sprintf("POSTGRES_PORTS=%s\n"+
		"POSTGRES_USER=%s\n"+
		"POSTGRES_PASSWORD=%s\n"+
		"POSTGRES_DB=%s\n"+
		"POSTGRES_CONTAINER_NAME=%s\n"+
		"NATS_CLUSTER=%s\n"+
		"NATS_PORTS=%s\n"+
		"NATS_CONTAINER_NAME=%s\n",
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.ContainerName,
		cfg.NatsStreaming.ClusterID,
		cfg.NatsStreaming.Port,
		cfg.NatsStreaming.ContainerName)

	file, err := os.Create(".env")
	if err != nil {
		log.Fatalf("%s: error while set environments: %s", ferr, err.Error())
	}
	defer file.Close()

	file.WriteString(text)
}

func createCfg() *Config {
	const ferr = "cmd.config.createCfg"

	cfgPath := os.Getenv("CONFIG_PATH")

	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		log.Fatalf("%s: config file is not exists: %s", ferr, cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("%s: error while reading config: %s", ferr, err.Error())
	}

	return &cfg
}
