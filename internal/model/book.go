package model

type Book struct {
	OrderId     string  `db:"order_id"`
	UpdateID    int     `db:"update_id"`
	Symbol      string  `db:"symbol"`
	BidPrice    float64 `db:"bid_price"`
	BidQuantity float64 `db:"bid_quantity"`
	AskPrice    float64 `db:"ask_price"`
	AskQuantity float64 `db:"ask_quantity"`
	CreatedAt   string  `db:"created_at"`
}
