package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderUID          uuid.UUID `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OOFShard          string    `json:"oof_shard"`
	Payment           Payment   `json:"payment"`
	Delivery          Delivery  `json:"delivery"`
	Items             []Item    `json:"items"`
}

type Item struct {
	OrderUID    uuid.UUID `json:"order_uid"`
	ChrtID      int       `json:"chrt_id"`
	TrackNumber string    `json:"track_number"`
	Price       float64   `json:"price"`
	Rid         string    `json:"rid"`
	Name        string    `json:"name"`
	Sale        float64   `json:"sale"`
	Size        uint      `json:"size"`
	TotalPrice  float64   `json:"total_price"`
	NmID        int       `json:"nm_id"`
	Brand       string    `json:"brand"`
	Status      int       `json:"status"`
}

type Payment struct {
	OrderUID     uuid.UUID `json:"order_uid"`
	Transaction  uuid.UUID `json:"transaction"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDT    int       `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost float64   `json:"delivery_cost"`
	GoodsTotal   uint      `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}

type Delivery struct {
	OrderUID uuid.UUID `json:"order_uid"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Zip      string    `json:"zip"`
	City     string    `json:"city"`
	Address  string    `json:"address"`
	Region   string    `json:"region"`
	Email    string    `json:"email"`
}
