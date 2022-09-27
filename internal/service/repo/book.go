package repo

import (
	"binance-order-matcher/internal/model"
	"binance-order-matcher/internal/service"
	"github.com/jmoiron/sqlx"
)

type bookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) service.BookRepo {
	return &bookRepo{db}
}

func (b *bookRepo) Save(book *model.Book) (*model.Book, error) {
	_, err := b.db.NamedExec(`
INSERT INTO books (order_id,update_id,symbol,bid_price,bid_quantity,ask_price,ask_quantity,created_at) 
VALUES (:order_id,:update_id,:symbol,:bid_price,:bid_quantity,:ask_price,:ask_quantity,:created_at)
`,
		map[string]interface{}{
			"order_id":     book.OrderId,
			"update_id":    book.UpdateID,
			"symbol":       book.Symbol,
			"bid_price":    book.BidPrice,
			"bid_quantity": book.BidQuantity,
			"ask_price":    book.AskPrice,
			"ask_quantity": book.AskQuantity,
			"created_at":   book.CreatedAt,
		})
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (b *bookRepo) GetByOrderId(orderId string) ([]*model.Book, error) {
	rows, err := b.db.NamedQuery(`SELECT * FROM books WHERE order_id=:order_id`, map[string]interface{}{"order_id": orderId})
	if err != nil {
		return nil, err
	}
	result := []*model.Book{}
	for rows.Next() {
		book := &model.Book{}
		err := rows.StructScan(book)
		if err != nil {
			return nil, err
		}
		result = append(result, book)
	}
	return result, nil
}
