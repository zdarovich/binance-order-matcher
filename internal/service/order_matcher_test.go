package service_test

import (
	"binance-order-matcher/internal/model"
	"binance-order-matcher/internal/service"
	"github.com/pkg/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	ErrInternal = errors.New("internal error")
)

type test struct {
	name   string
	mock   func()
	order  *model.Order
	book   *model.Book
	result *model.Order
	err    error
}

func orderMatcher(t *testing.T) (*service.OrderMatcher, *MockBookRepo, *MockSubscriptionAPI) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockBookRepo(mockCtl)
	webAPI := NewMockSubscriptionAPI(mockCtl)

	om := service.NewOrderMatcher(repo, webAPI)

	return om, repo, webAPI
}

func TestOrderMatcher_GetMatchedOrder(t *testing.T) {
	t.Parallel()

	om, bookRepo, _ := orderMatcher(t)
	bookBidSmall := &model.Book{UpdateID: 1, Symbol: "BNBUSDT", BidPrice: 3, BidQuantity: 1}
	bookBidBig := &model.Book{UpdateID: 2, Symbol: "BNBUSDT", BidPrice: 3, BidQuantity: 4}
	bookAskBig := &model.Book{UpdateID: 3, Symbol: "BNBUSDT", AskPrice: 3, AskQuantity: 4}

	tests := []test{
		{
			name: "GetMatchedOrder should return error if repo update fails",
			mock: func() {
				bookRepo.EXPECT().Save(bookBidSmall).Return(nil, ErrInternal)
			},
			order:  &model.Order{Id: "id", Type: model.SellOrder, Symbol: "BNBUSDT", Quantity: 3, Price: 3},
			book:   bookBidSmall,
			result: nil,
			err:    ErrInternal,
		},
		{
			name: "GetMatchedOrder should return nil if order doesn't match",
			mock: func() {
			},
			order:  &model.Order{Id: "id", Type: model.SellOrder, Symbol: "BNBUSDT", Quantity: 3, Price: 3},
			book:   bookAskBig,
			result: nil,
			err:    nil,
		},
		{
			name: "GetMatchedOrder should return updated order",
			mock: func() {
				bookRepo.EXPECT().Save(bookBidSmall).Return(bookBidSmall, nil)
			},
			order:  &model.Order{Id: "id", Type: model.SellOrder, Symbol: "BNBUSDT", Quantity: 3, Price: 3},
			book:   bookBidSmall,
			result: &model.Order{Id: "id", Type: model.SellOrder, Symbol: "BNBUSDT", Quantity: 2, Price: 3},
			err:    nil,
		},
		{
			name: "GetMatchedOrder should return updated order with quantity less than zero",
			mock: func() {
				bookRepo.EXPECT().Save(bookBidBig).Return(bookBidBig, nil)
			},
			order:  &model.Order{Id: "id", Type: model.SellOrder, Symbol: "BNBUSDT", Quantity: 3, Price: 3},
			book:   bookBidBig,
			result: &model.Order{Id: "id", Type: model.SellOrder, Symbol: "BNBUSDT", Quantity: -1, Price: 3},
			err:    nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			tc.mock()

			updatedOrder, err := om.GetMatchedOrder(tc.order, tc.book)

			require.Equal(t, tc.result, updatedOrder)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
