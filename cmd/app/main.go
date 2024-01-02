package main

import (
	"fmt"
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
	//logsHandler.SilenceOperatingMode(true) //-- Выключает логи

	logsHandler.WriteInfo("loading config...")
	cfg := config.New()

	logsHandler.WriteInfo("connecting to database...")
	db := storage.New(cfg, logsHandler)
	defer db.CloseConnection()

	logsHandler.WriteInfo("creating cache...")
	cache := cache.New(logsHandler)
	storageManager := storagemanager.New(db, cache, logsHandler)

	logsHandler.WriteInfo("connecting to nats-streaming...")
	nats := broker.New(cfg, &storageManager, logsHandler)
	nats.SubscribeAndHandle()
	defer nats.CloseConnection()

	srv := server.New(cfg, server.HFuncList{
		OrderGetter:  &storageManager,
		OrdersGetter: &storageManager,
		OrderSaver:   &storageManager,
	})

	logsHandler.WriteInfo(fmt.Sprintf("SERVER STARTED [%s:%s]",
		cfg.Server.Host, cfg.Server.Port))

	srv.Start()

	logsHandler.WriteInfo("SERVER CLOSED")
}
