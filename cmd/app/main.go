package main

import (
	"fmt"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
)

const(
	logsPath = "./logs/logs.txt"
)

func main() {
	logsHandler := logs.New(logsPath)
	defer logsHandler.Close()
	
	cfg := config.New()

	fmt.Println(cfg.Postgres)
}