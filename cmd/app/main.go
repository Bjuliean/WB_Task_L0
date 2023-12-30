package main

// b563feb7b2b84b6test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"wbl0/WB_Task_L0/internal/broker"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/storage"
	storagemanager "wbl0/WB_Task_L0/internal/storage_manager"

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

	storageManager := storagemanager.New(db, logsHandler)

	nats := broker.New(cfg, &storageManager, logsHandler)
	nats.SubscribeAndHandle()
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
	time.Sleep(5 * time.Second)

	a, _ := uuid.Parse("d101af40-dc63-51af-90d2-a125d1d49f4d")
	uhah, _ := db.GetOrder(a)
	fmt.Println(uhah)

}
