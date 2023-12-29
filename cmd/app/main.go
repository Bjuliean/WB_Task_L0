package main

// b563feb7b2b84b6test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/storage"
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

	file, err := os.Open("./misc/test1.json")

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

	uhah, err := db.GetOrders()
	if err != nil {
		log.Fatalf("marshal err: %s", err.Error())
	}
	for _, item := range uhah {
		fmt.Println(item)
	}

	if err != nil {
		log.Fatalf("failed to create order: %s", err.Error())
	}
}
