package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"wbl0/WB_Task_L0/internal/config"
	"wbl0/WB_Task_L0/internal/logs"
	"wbl0/WB_Task_L0/internal/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db          *sql.DB
	logsHandler *logs.Logger
}

func New(cfg *config.Config, logs *logs.Logger) *Storage {
	const ferr = "internal.storage.New"

	dataSrcName := fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		"dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName)

	db, err := sql.Open("postgres", dataSrcName)
	if err != nil {
		log.Fatalf("%s: error while opening db: %s", ferr, err.Error())
	}

	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Fatalf("%s: db connection failed: %s", ferr, err.Error())
	}

	return &Storage{
		db:          db,
		logsHandler: logs,
	}
}

func (s *Storage) CloseConnection() {
	s.db.Close()
}

func (s *Storage) CreateOrder(order models.Order) error {
	const op = "internal.storage.CreateOrder"
	ferr := op + fmt.Sprintf(" (%v)", order.OrderUID)

	if err := s.pushOrder(order); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	if err := s.pushItems(order); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	if err := s.pushPayment(order); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	if err := s.pushDelivery(order); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	s.logsHandler.WriteInfo(fmt.Sprintf("order created: %v", order.OrderUID))

	return nil
}

func (s *Storage) pushOrder(order models.Order) error {
	const ferr = "internal.storage.pushOrder"

	stOrder, err := s.db.Prepare("INSERT INTO orders(order_uid, track_number, entry, " +
		"locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, " +
		"oof_shard) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);")
	defer stOrder.Close()

	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	_, err = stOrder.Exec(order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OOFShard)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	return nil
}

func (s *Storage) pushItems(order models.Order) error {
	const ferr = "internal.storage.pushItems"

	stItems, err := s.db.Prepare("INSERT INTO items(order_uid, chrt_id, track_number, price, rid, " +
		"name, sale, size, total_price, nm_id, brand, status) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);")
	defer stItems.Close()

	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	for _, item := range order.Items {
		_, err := stItems.Exec(order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand,
			item.Status)
		if err != nil {
			s.logsHandler.WriteError(ferr, err.Error())
			return err
		}
	}

	return nil
}

func (s *Storage) pushPayment(order models.Order) error {
	const ferr = "internal.storage.pushPayment"

	stPayment, err := s.db.Prepare("INSERT INTO payments(order_uid, transaction, request_id, currency, " +
		"provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);")
	defer stPayment.Close()

	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	_, err = stPayment.Exec(order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	return nil
}

func (s *Storage) pushDelivery(order models.Order) error {
	const ferr = "internal.storage.pushDelivery"

	stDelivery, err := s.db.Prepare("INSERT INTO deliveries(order_uid, name, phone, zip, city, address, " +
		"region, email) VALUES($1, $2, $3, $4, $5, $6, $7, $8);")
	defer stDelivery.Close()

	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	_, err = stDelivery.Exec(order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	return nil
}
