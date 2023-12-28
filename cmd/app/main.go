package main

import (
	"wbl0/WB_Task_L0/internal/logs"
)

const(
	logsPath = "./logs/logs.txt"
)

func main() {
	logsHandler := logs.New(logsPath)
	defer logsHandler.Close()
	logsHandler.WriteInfo("AAAAA")
	logsHandler.WriteError("DDDDDD")
}