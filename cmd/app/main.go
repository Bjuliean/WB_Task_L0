package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"wbl0/WB_Task_L0/internal/broker"
	"wbl0/WB_Task_L0/internal/cache"
	"wbl0/WB_Task_L0/internal/closer"
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
	//logsHandler.SilenceOperatingMode(true) //-- Логи только в logs/logs.txt

	logsHandler.WriteInfo("loading config...")
	cfg := config.New()

	logsHandler.WriteInfo("connecting to database...")
	db := storage.New(cfg, logsHandler)

	logsHandler.WriteInfo("creating cache...")
	cache := cache.New(logsHandler)
	storageManager := storagemanager.New(db, cache, logsHandler)

	logsHandler.WriteInfo("connecting to nats-streaming...")
	nats := broker.New(cfg, &storageManager, logsHandler)
	nats.SubscribeAndHandle()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv := server.New(cfg, server.HFuncList{
		OrderGetter:  &storageManager,
		OrdersGetter: &storageManager,
		OrderSaver:   &storageManager,
	}, &ctx)

	closer := closer.New()
	closer.AddToList(logsHandler.Close, db.CloseConnection,
		nats.CloseConnection, srv.Stop)

	logsHandler.WriteInfo(fmt.Sprintf("SERVER STARTED [%s:%s]",
		cfg.Server.Host, cfg.Server.Port))

	go func() {
		srv.Start()
	}()

	<-ctx.Done()
	logsHandler.WriteInfo("CLOSING SERVER...")
	closer.Shutdown()
}
