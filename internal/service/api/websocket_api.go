package api

import (
	"binance-order-matcher/internal/model"
	"binance-order-matcher/internal/service"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type WebsocketApi struct {
	host string
}

type WsBookTickerHandler func(event *Book)
type ErrHandler func(err error)

type Book struct {
	UpdateID    uint    `json:"u,omitempty"`
	Symbol      string  `json:"s,omitempty"`
	BidPrice    float64 `json:"b,string,omitempty"`
	BidQuantity float64 `json:"B,string,omitempty"`
	AskPrice    float64 `json:"a,string,omitempty"`
	AskQuantity float64 `json:"A,string,omitempty"`
}

func New(host string) *WebsocketApi {
	return &WebsocketApi{
		host: host,
	}
}

func (br *WebsocketApi) Connect(symbol string) (<-chan struct{}, chan<- struct{}, <-chan service.ApiResponse, error) {

	c, err := br.connectToStream(symbol)
	if err != nil {
		return nil, nil, nil, err
	}
	c.SetReadLimit(655350)
	doneC := make(chan struct{})
	stopC := make(chan struct{})
	apiChan := make(chan service.ApiResponse)
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer close(doneC)
		defer close(apiChan)
		defer log.Info("binance api exited")
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			c.Close()
		}()
		for {
			book := new(Book)
			err := c.ReadJSON(book)
			if err != nil {
				if !silent {
					apiChan <- service.ApiResponse{Error: err}
				}
				return
			}
			apiChan <- service.ApiResponse{Book: ApiBookToBook(book)}
		}
	}()
	return doneC, stopC, apiChan, nil
}

func (br *WebsocketApi) connectToStream(symbol string) (*websocket.Conn, error) {
	url := fmt.Sprintf("%s/%s@bookTicker", br.host, strings.ToLower(symbol))
	conn, _, err := websocket.DefaultDialer.DialContext(context.Background(), url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to stream")
	}

	return conn, nil
}

func ApiBookToBook(book *Book) *model.Book {
	return &model.Book{
		UpdateID:    int(book.UpdateID),
		Symbol:      book.Symbol,
		BidPrice:    book.BidPrice,
		BidQuantity: book.BidQuantity,
		AskPrice:    book.AskPrice,
		AskQuantity: book.AskQuantity,
	}
}
