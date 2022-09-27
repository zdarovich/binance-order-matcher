package repo

import (
	"binance-order-matcher/internal/model"
	"binance-order-matcher/internal/service"
	"github.com/jmoiron/sqlx"
)

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) service.OrderRepo {
	return &orderRepo{db}
}

func (or *orderRepo) Save(order *model.Order) (*model.Order, error) {
	_, err := or.db.NamedExec(`
INSERT INTO orders (id,order_type,symbol,price,quantity,created_at) 
VALUES (:id,:order_type,:symbol,:price,:quantity,:created_at)
`,
		map[string]interface{}{
			"id":         order.Id,
			"order_type": order.Type,
			"symbol":     order.Symbol,
			"price":      order.Price,
			"quantity":   order.Quantity,
			"created_at": order.CreatedAt.String(),
		})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (or *orderRepo) GetById(orderId string) ([]*model.Order, error) {
	return nil, nil
}
