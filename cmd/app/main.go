package main

import (
	"encoding/json"
	"fmt"
	//"fmt"
	"log"
	"os"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/storage"
)

const(
	logsPath = "./logs/logs.txt"
)

func main() {
	logsHandler := logs.New(logsPath)
	defer logsHandler.Close()
	
	cfg := config.New()

	db := storage.New(cfg, logsHandler)
	defer db.CloseConnection()

	file, _ := os.Open("./misc/model.json")
	var o models.Order
	json.NewDecoder(file).Decode(&o)

	err := db.CreateOrder(o)

	fmt.Println(o)

	if err != nil {
		log.Fatalf("failed to create order: %s", err.Error())
	}
}