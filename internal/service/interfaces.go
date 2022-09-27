//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=service_test
package service

import (
	"binance-order-matcher/internal/model"
)

type (
	SubscriptionAPI interface {
		Connect(symbol string) (<-chan struct{}, chan<- struct{}, <-chan ApiResponse, error)
	}

	ApiResponse struct {
		Book  *model.Book
		Error error
	}

	BookRepo interface {
		Save(book *model.Book) (*model.Book, error)
		GetByOrderId(orderId string) ([]*model.Book, error)
	}

	OrderRepo interface {
		Save(order *model.Order) (*model.Order, error)
		GetById(orderId string) ([]*model.Order, error)
	}
)
