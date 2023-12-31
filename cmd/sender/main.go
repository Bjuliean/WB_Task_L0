package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nats-io/stan.go"
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
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	ClientID         string `yaml:"client_id"`
	ClusterID        string `yaml:"cluster_id"`
	SubscribeSubject string `yaml:"subscribe_subject"`
	ContainerName    string `yaml:"containername"`
}

func main() {
	const (
		ferr       = "cmd.sender.main"
		clientID   = "2"
		totalTests = 3
		tFilesPath = "./misc/test"
	)

	cfg := createCfg()
	var testFiles [][]byte

	sc, err := stan.Connect(cfg.NatsStreaming.ClusterID, clientID, stan.NatsURL(fmt.Sprintf("%s:%s", cfg.NatsStreaming.Host, cfg.NatsStreaming.Port)))
	if err != nil {
		log.Printf("%s: failed to connect nats: %s", ferr, err.Error())
		return
	}

	for i := 0; i < totalTests; i++ {
		filePath := tFilesPath + strconv.Itoa(i) + ".json"
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("%s: failed to open files: %s", ferr, err.Error())
			return
		}

		dt, err := io.ReadAll(file)
		if err != nil {
			log.Printf("%s: error while reading file: %s", ferr, err.Error())
			return
		}

		testFiles = append(testFiles, dt)
	}
	
	for _, test := range testFiles {
		sc.Publish(cfg.NatsStreaming.SubscribeSubject, test)
	}

	// file, err := os.Open("./misc/model.json")
	// if err != nil {
	// 	log.Printf("%s: failed to open files: %s", ferr, err.Error())
	// 	return
	// }

	// dt, err := io.ReadAll(file)
	// if err != nil {
	// 	log.Printf("%s: error while reading file: %s", ferr, err.Error())
	// 	return
	// }

	// sc, err := stan.Connect(cfg.NatsStreaming.ClusterID, clientID, stan.NatsURL(fmt.Sprintf("%s:%s", cfg.NatsStreaming.Host, cfg.NatsStreaming.Port)))
	// if err != nil {
	// 	log.Printf("%s: failed to connect nats: %s", ferr, err.Error())
	// 	return
	// }

	// err = sc.Publish(cfg.NatsStreaming.SubscribeSubject, dt)
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
