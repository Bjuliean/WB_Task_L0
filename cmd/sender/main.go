package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
	"wbl0/WB_Task_L0/internal/models"

	"github.com/google/uuid"
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

const (
	lettersKit  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 "
	ferr       = "cmd.sender.main"
	clientID   = "2"
	totalTestsPrepared = 3
	totalTestsRandom = 300
	tFilesPath = "./misc/test"
	maxItemsQuantity = 30
)

func main() {
	cfg := createCfg()

	sc, err := stan.Connect(cfg.NatsStreaming.ClusterID, clientID)
	if err != nil {
		log.Printf("%s: failed to connect nats: %s", ferr, err.Error())
		return
	}

	pTests := preparedTests()

	for _, test := range pTests {
		sc.Publish(cfg.NatsStreaming.SubscribeSubject, test)
	}

	for i := 0; i < totalTestsRandom; i++ {
		rOrder := randomOrder()

		dt, err := json.Marshal(rOrder)
		if err != nil {
			log.Printf("%s: failed to code order: %s", ferr, err.Error())
		}

		sc.Publish(cfg.NatsStreaming.SubscribeSubject, dt)
	}

}

func preparedTests() [][]byte {
	var res [][]byte
	
	for i := 0; i < totalTestsPrepared; i++ {
		filePath := tFilesPath + strconv.Itoa(i) + ".json"
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("%s: failed to open files: %s", ferr, err.Error())
			return nil
		}

		dt, err := io.ReadAll(file)
		if err != nil {
			log.Printf("%s: error while reading file: %s", ferr, err.Error())
			return nil
		}

		res = append(res, dt)
	}

	return res
}

func randomOrder() models.Order {
	u, _ := uuid.NewRandom()

	return models.Order{
		OrderUID:    u,
		TrackNumber: u.String(),
		Entry: randomStringVal(6, true),
		Locale: randomStringVal(4, true),
		InternalSignature: randomStringVal(10, false),
		CustomerID: randomStringVal(10, false),
		DeliveryService: randomStringVal(20, true),
		Shardkey: randomStringVal(10, false),
		SmID: rand.Intn(1000),
		DateCreated: randomDate(),
		OOFShard: randomStringVal(10, true),
		Payment: randomPayment(u),
		Delivery: randomDelivery(u),
		Items: randomItems(u, u.String(), maxItemsQuantity),
	}
}

func randomItems(ru uuid.UUID, tr string, quantity int) []models.Item {
	var res []models.Item
	quantity = rand.Intn(quantity)

	for i := 0; i < quantity; i++ {
		res = append(res, randomItem(ru, tr))
	}

	return res
}

func randomItem(ru uuid.UUID, tr string) models.Item {
	return models.Item{
		OrderUID: ru,
		ChrtID: rand.Intn(10000),
		TrackNumber: tr,
		Price: rand.Float64(),
		Rid: randomStringVal(20, false),
		Name: randomStringVal(15, true),
		Sale: rand.Float64(),
		Size: randomStringVal(4, true),
		TotalPrice: rand.Float64(),
		NmID: rand.Intn(10000),
		Brand: randomStringVal(20, true),
		Status: rand.Intn(999),
	}
}

func randomDelivery(ru uuid.UUID) models.Delivery {
	return models.Delivery{
		OrderUID: ru,
		Name: randomStringVal(20, true),
		Phone: randomStringVal(10, false),
		Zip: randomStringVal(20, true),
		City: randomStringVal(20, true),
		Address: randomStringVal(20, true),
		Region: randomStringVal(20, true),
		Email: randomStringVal(20, true),
	}
}

func randomPayment(ru uuid.UUID) models.Payment {
	u, _ := uuid.NewRandom()

	return models.Payment{
		OrderUID: ru,
		Transaction: u,
		RequestID: randomStringVal(5, false),
		Currency: randomCurrency(),
		Provider: randomStringVal(10, true),
		Amount: rand.Float64(),
		PaymentDT: rand.Intn(100000),
		Bank: randomStringVal(20, true),
		DeliveryCost: rand.Float64(),
		CustomFee: rand.Intn(10000),
	}
}

func randomCurrency() string {
	var curKit []string = []string{
		"USD",
		"EUR",
		"RUB",
		"XAF",
		"AUD",
	}

	return curKit[rand.Intn(len(curKit))]
}

func randomDate() time.Time {
	startDate := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	dur := endDate.Sub(startDate)

	randdur := time.Duration(rand.Int63n(int64(dur)))

	return startDate.Add(randdur)
}

func randomStringVal(size int, roll_size bool) string {
	if roll_size {
		size = rand.Intn(size)
	}

	res := make([]byte, size)

	for i := 0; i < size; i++ {
		res[i] = lettersKit[rand.Intn(len(lettersKit))]
	}

	return string(res)
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
