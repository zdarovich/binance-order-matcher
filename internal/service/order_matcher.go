package service

import (
	"binance-order-matcher/internal/model"
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	ErrTimeout = errors.New("timeout exception")
)

type OrderMatcher struct {
	client   SubscriptionAPI
	bookRepo BookRepo
	quit     chan bool
}

func NewOrderMatcher(br BookRepo, client SubscriptionAPI) *OrderMatcher {
	return &OrderMatcher{
		client:   client,
		bookRepo: br,
	}
}

func (t *OrderMatcher) Match(ctx context.Context, timeout time.Duration, order *model.Order) (bool, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, binanceStop, binanceApi, err := t.client.Connect(order.Symbol)
	if err != nil {
		return false, err
	}
	defer close(binanceStop)

	go func() {
		select {
		case <-time.After(timeout):
		case <-ctx.Done():
		}
		binanceStop <- struct{}{}
	}()

	var match bool
	for response := range binanceApi {
		if response.Error != nil {
			match = false
			err = response.Error
			break
		}
		book := response.Book
		log.Infof("Book received %+v", book)
		updatedOrder, matchErr := t.GetMatchedOrder(order, book)
		if err != nil {
			match = false
			err = matchErr
			cancel()
		}
		if updatedOrder != nil {
			order = updatedOrder
		}
		if order.Quantity <= 0 {
			match = true
			cancel()
		}
	}
	return match, err
}

func (t *OrderMatcher) GetMatchedOrder(order *model.Order, book *model.Book) (*model.Order, error) {
	if order.Quantity <= 0 {
		return nil, nil
	}
	isMatching, _, quantity := getMatchedPriceQuantity(order, book)
	if !isMatching {
		return nil, nil
	}

	err := t.createTransaction(order, book)
	if err != nil {
		return nil, err
	}
	updatedOrder := new(model.Order)
	updatedOrder.Id = order.Id
	updatedOrder.Symbol = order.Symbol
	updatedOrder.Price = order.Price
	updatedOrder.Type = order.Type
	updatedOrder.CreatedAt = order.CreatedAt
	updatedOrder.Quantity = order.Quantity - quantity

	return updatedOrder, nil
}

func (t *OrderMatcher) createTransaction(order *model.Order, book *model.Book) error {
	if book == nil {
		return nil
	}
	book.CreatedAt = time.Now().String()
	book.OrderId = order.Id
	savedBook, err := t.bookRepo.Save(book)
	fmt.Sprintf("save %+v", savedBook)
	return err
}

func getMatchedPriceQuantity(order *model.Order, book *model.Book) (bool, float64, float64) {
	switch order.Type {
	case model.BuyOrder:
		if order.Price >= book.AskPrice {
			return true, book.AskPrice, book.AskQuantity
		}
	case model.SellOrder:
		if order.Price <= book.BidPrice {
			return true, book.BidPrice, book.BidQuantity
		}
	}
	return false, 0, 0
}
