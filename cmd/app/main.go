package main

import (
	"wbl0/WB_Task_L0/internal/broker"
	"wbl0/WB_Task_L0/internal/cache"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/server"
	"wbl0/WB_Task_L0/internal/storage"
	storagemanager "wbl0/WB_Task_L0/internal/storage_manager"
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
	logsHandler.WriteInfo("connected to database")

	cache := cache.New(logsHandler)
	logsHandler.WriteInfo("cache created")

	storageManager := storagemanager.New(db, cache, logsHandler)

	nats := broker.New(cfg, &storageManager, logsHandler)
	nats.SubscribeAndHandle()
	defer nats.CloseConnection()
	logsHandler.WriteInfo("connected to nats-streaming")

	srv := server.New(cfg, server.HFuncList{
		OrderGetter: &storageManager,
		OrdersGetter: &storageManager,
		OrderSaver: &storageManager,
	})

	logsHandler.WriteInfo("SERVER STARTED")

	srv.Start()
	
	logsHandler.WriteInfo("SERVER CLOSED")
}
