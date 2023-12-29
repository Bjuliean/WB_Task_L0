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

	tx, err := s.db.Begin()
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		tx.Rollback()
		return err
	}

	if err := s.pushOrder(tx, order); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		tx.Rollback()
		return err
	}

	if err := s.pushItems(tx, order); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		tx.Rollback()
		return err
	}

	if err := s.pushPayment(tx, order); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		tx.Rollback()
		return err
	}

	if err := s.pushDelivery(tx, order); err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		tx.Rollback()
		return err
	}

	tx.Commit()
	s.logsHandler.WriteInfo(fmt.Sprintf("order created: %v", order.OrderUID))

	return nil
}

func (s *Storage) GetOrders() ([]models.Order, error) {
	const ferr = "internal.storage.GetOrders"
	var res []models.Order

	rowsOrder, err := s.db.Query("SELECT * FROM orders;")
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return nil, err
	}

	for rowsOrder.Next() {
		var singleOrder models.Order

		err := rowsOrder.Scan(&singleOrder.OrderUID, &singleOrder.TrackNumber,
			&singleOrder.Entry, &singleOrder.Locale, &singleOrder.InternalSignature,
			&singleOrder.CustomerID, &singleOrder.DeliveryService, &singleOrder.Shardkey,
			&singleOrder.SmID, &singleOrder.DateCreated, &singleOrder.OOFShard)
		if err != nil {
			s.logsHandler.WriteError(ferr, err.Error())
			return nil, err
		}
		err = s.assembleOrder(&singleOrder)
		if err != nil {
			s.logsHandler.WriteError(ferr, err.Error())
			return nil, err
		}
		res = append(res, singleOrder)
	}

	return res, nil
}

func (s *Storage) assembleOrder(order *models.Order) error {
	const ferr = "internal.storage.assembleOrder"
	var err error

	order.Items, err = s.searchItems(order)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	order.Payment, err = s.searchPayment(order)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	order.Delivery, err = s.searchDelivery(order)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	return nil
}

func (s *Storage) searchItems(order *models.Order) ([]models.Item, error) {
	const ferr = "internal.storage.searchItems"
	var res []models.Item

	st, err := s.db.Prepare("SELECT * FROM items WHERE order_uid = $1;")
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return []models.Item{}, err
	}

	rows, err := st.Query(order.OrderUID)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return []models.Item{}, err
	}

	for rows.Next() {
		var singleItem models.Item

		err := rows.Scan(&singleItem.OrderUID, &singleItem.ChrtID, &singleItem.TrackNumber,
			&singleItem.Price, &singleItem.Rid, &singleItem.Name, &singleItem.Sale, &singleItem.Size,
			&singleItem.TotalPrice, &singleItem.NmID, &singleItem.Brand, &singleItem.Status)
		if err != nil {
			s.logsHandler.WriteError(ferr, err.Error())
			return []models.Item{}, err
		}

		res = append(res, singleItem)
	}

	return res, nil
}

func (s *Storage) searchPayment(order *models.Order) (models.Payment, error) {
	const ferr = "internal.storage.searchPayment"
	var res models.Payment

	row, err := s.db.Prepare("SELECT * FROM payments WHERE order_uid = $1;")
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return models.Payment{}, err
	}
	err = row.QueryRow(order.OrderUID).Scan(&res.OrderUID, &res.Transaction, &res.RequestID,
		&res.Currency, &res.Provider, &res.Amount, &res.PaymentDT, &res.Bank, &res.DeliveryCost,
		&res.GoodsTotal, &res.CustomFee)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return models.Payment{}, err
	}

	return res, nil
}

func (s *Storage) searchDelivery(order *models.Order) (models.Delivery, error) {
	const ferr = "internal.storage.searchDelivery"
	var res models.Delivery

	row, err := s.db.Prepare("SELECT * FROM deliveries WHERE order_uid = $1;")
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return models.Delivery{}, err
	}

	err = row.QueryRow(order.OrderUID).Scan(&res.OrderUID, &res.Name, &res.Phone, &res.Zip, &res.City,
		&res.Address, &res.Region, &res.Email)
	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return models.Delivery{}, err
	}

	return res, nil
}

func (s *Storage) pushOrder(tx *sql.Tx, order models.Order) error {
	const ferr = "internal.storage.pushOrder"

	stOrder, err := tx.Prepare("INSERT INTO orders(order_uid, track_number, entry, " +
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

func (s *Storage) pushItems(tx *sql.Tx, order models.Order) error {
	const ferr = "internal.storage.pushItems"

	stItems, err := tx.Prepare("INSERT INTO items(order_uid, chrt_id, track_number, price, rid, " +
		"name, sale, size, total_price, nm_id, brand, status) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);")
	defer stItems.Close()

	if err != nil {
		s.logsHandler.WriteError(ferr, err.Error())
		return err
	}

	for _, item := range order.Items {
		_, err := stItems.Exec(order.OrderUID, item.ChrtID, order.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand,
			item.Status)
		if err != nil {
			s.logsHandler.WriteError(ferr, err.Error())
			return err
		}
	}

	return nil
}

func (s *Storage) pushPayment(tx *sql.Tx, order models.Order) error {
	const ferr = "internal.storage.pushPayment"

	stPayment, err := tx.Prepare("INSERT INTO payments(order_uid, transaction, request_id, currency, " +
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

func (s *Storage) pushDelivery(tx *sql.Tx, order models.Order) error {
	const ferr = "internal.storage.pushDelivery"

	stDelivery, err := tx.Prepare("INSERT INTO deliveries(order_uid, name, phone, zip, city, address, " +
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
