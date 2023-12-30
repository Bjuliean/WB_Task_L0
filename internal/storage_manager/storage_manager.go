package storagemanager

import (
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/storage"
)

type StorageManager struct {
	database *storage.Storage
	logsHandler *logs.Logger
}

func New(db *storage.Storage, logs *logs.Logger) StorageManager { //todo cache
	return StorageManager{
		database: db,
		logsHandler: logs,
	}
}

func (s *StorageManager) AddOrder(nworder models.Order) error {
	const ferr = "internal.storage_manager.AddOrder"
	
	if err := s.database.CreateOrder(nworder); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	return nil
}
