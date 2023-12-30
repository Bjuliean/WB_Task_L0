package cache

import (
	"fmt"
	"sync"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"

	"github.com/google/uuid"
)

type Cache struct {
	mu sync.RWMutex
	storage     map[string]models.Order
	logsHandler *logs.Logger
}

func New(logs *logs.Logger) *Cache {
	return &Cache{
		storage:     make(map[string]models.Order),
		logsHandler: logs,
	}
}

func (c *Cache) CreateOrder(order models.Order) error {
	const ferr = "internal.cache.New"

	c.mu.RLock()
	if _, exists := c.storage[order.OrderUID.String()]; exists == true {
		msg := fmt.Sprintf("%s (%v): already exists", ferr, order.OrderUID)
		c.logsHandler.WriteInfo(msg)
		return fmt.Errorf(msg)
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.storage[order.OrderUID.String()] = order

	return nil
}

func (c *Cache) GetOrders() ([]models.Order, error) {
	const ferr = "internal.cache.GetOrders"
	var res []models.Order

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, item := range c.storage {
		res = append(res, item)
	}

	return res, nil
}

func (c *Cache) GetOrder(uid uuid.UUID) (*models.Order, error) {
	const ferr = "internal.cache.GetOrder"
	var res models.Order

	c.mu.RLock()
	defer c.mu.RUnlock()

	res, exists := c.storage[uid.String()]
	if exists == false {
		msg := fmt.Sprintf("%s (%v): could not be found", ferr, uid)
		c.logsHandler.WriteInfo(msg)
		return &models.Order{}, fmt.Errorf(msg)
	}

	return &res, nil
}

func (c *Cache) ReloadCache(orders []models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for k := range c.storage {
		delete(c.storage, k)
	}
	for _, order := range orders {
		c.storage[order.OrderUID.String()] = order
	}
}
