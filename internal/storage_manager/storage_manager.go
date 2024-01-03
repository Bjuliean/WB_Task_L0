package storagemanager

import (
	"fmt"
	"log"
	"wbl0/WB_Task_L0/internal/cache"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/storage"

	"github.com/google/uuid"
)

type StorageManager struct {
	database     *storage.Storage
	cacheStorage *cache.Cache
	logsHandler  *logs.Logger
}

func New(db *storage.Storage, cch *cache.Cache, logs *logs.Logger) StorageManager {
	tmp, err := db.GetOrders()
	if err != nil {
		log.Fatalf("cache initialize failed")
	}
	cch.ReloadCache(tmp)

	return StorageManager{
		database:     db,
		cacheStorage: cch,
		logsHandler:  logs,
	}
}

func (s *StorageManager) AddOrder(nworder models.Order) error {
	const ferr = "internal.storage_manager.AddOrder"

	err := s.database.CreateOrder(nworder)
	if err != nil {
		msg := fmt.Sprintf("%s (%v)", ferr, nworder.OrderUID)
		s.logsHandler.WriteError(msg, err.Error())
		return err
	}

	err = s.cacheStorage.CreateOrder(nworder)
	if err != nil {
		msg := fmt.Sprintf("%s (%v)", ferr, nworder.OrderUID)
		s.logsHandler.WriteError(msg, err.Error())
		return err
	}

	return nil
}

func (s *StorageManager) GetOrder(uid uuid.UUID) (*models.Order, error) {
	const ferr = "internal.storage_manager.GetOrder"
	var res *models.Order

	res, err := s.cacheStorage.GetOrder(uid)
	if err != nil {
		res, err = s.database.GetOrder(uid)
		if err != nil {
			msg := fmt.Sprintf("%s (%v)", ferr, uid)
			s.logsHandler.WriteError(msg, err.Error())
			return &models.Order{}, err
		}

		s.logsHandler.WriteInfo(fmt.Sprintf("%v: not found in cache, recorded in db", uid))
		s.cacheStorage.CreateOrder(*res)
	}

	return res, nil
}

func (s *StorageManager) GetOrders() ([]models.Order, error) {
	const ferr = "internal.storage_manager.GetOrders"

	var res []models.Order

	res, err := s.cacheStorage.GetOrders()
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return nil, err
	}

	return res, nil
}
