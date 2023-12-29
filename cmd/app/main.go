package main

// b563feb7b2b84b6test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"wbl0/WB_Task_L0/internal/broker"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/storage"

	"github.com/google/uuid"
)

const (
	logsPath = "./logs/logs.txt"
)

func main() {
	logsHandler := logs.New(logsPath)
	defer logsHandler.Close()

	cfg := config.New()

	db := storage.New(cfg, logsHandler)
	defer db.CloseConnection()

	nats := broker.New(cfg, logsHandler)
	defer nats.CloseConnection()

	file, err := os.Open("./misc/test1.json") // model, test2

	byts, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("read err: %s", err.Error())
	}

	var o models.Order
	err = json.Unmarshal(byts, &o)
	if err != nil {
		log.Fatalf("marshal err: %s", err.Error())
	}

	err = db.CreateOrder(o)

	a, _ := uuid.Parse("d161bf21-dc63-41af-90d2-b025d1d49f4d")
	uhah, err := db.GetOrder(a)
	fmt.Println(uhah)

}
