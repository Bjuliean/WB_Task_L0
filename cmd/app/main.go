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
	"wbl0/WB_Task_L0/internal/cache"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/server"
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

	cache := cache.New(logsHandler)

	storageManager := storagemanager.New(db, cache, logsHandler)

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
	time.Sleep(0 * time.Second)

	a, _ := uuid.Parse("d161bf21-dc63-41af-90d2-b025d1d49f4d")
	uhah, _ := db.GetOrder(a)
	fmt.Println(uhah)




	srv := server.New(cfg, server.HFuncList{
		OrderGetter: &storageManager,
	})
	srv.Start()


	// router := chi.NewRouter()

	// router.Use(middleware.Recoverer)
	// router.Use(middleware.RequestID)
	// router.Use(middleware.URLFormat)

	// router.Get("/{order_uid}", func(w http.ResponseWriter, r *http.Request) {
	// 	order_uid := chi.URLParam(r, "order_uid")
	// 	uid_val, _ := uuid.Parse(order_uid)
	// 	res, _ := storageManager.GetOrder(uid_val)
	// 	render.JSON(w, r, res)
	// })

	// srv := http.Server{
	// 	Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
	// 	Handler:      router,
	// }

	// srv.ListenAndServe()

}
