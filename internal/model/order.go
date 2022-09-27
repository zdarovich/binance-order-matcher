package model

import "time"

const (
	BuyOrder = iota
	SellOrder
)

type Order struct {
	Id        string
	Type      int
	Symbol    string
	Quantity  float64
	Price     float64
	CreatedAt time.Time
}

func CreateOrder(id string, orderType int, symbol string, size, priceLimit float64) *Order {
	order := Order{
		Id:        id,
		Type:      orderType,
		Symbol:    symbol,
		Quantity:  size,
		Price:     priceLimit,
		CreatedAt: time.Now(),
	}
	return &order
}
