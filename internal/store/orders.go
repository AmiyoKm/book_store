package store

import (
	"context"
	"database/sql"
	"time"
)

type Order struct {
	ID              int         `json:"id"`
	UserID          int         `json:"user_id"`
	TotalAmount     float64     `json:"total_amount"`
	Status          string      `json:"status"`
	PaymentMethod   string      `json:"payment_method"`
	ShippingAddress string      `json:"shipping_address"`
	PlacedAt        time.Time   `json:"placed_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Items           []OrderItem `json:"order_items"`
}
type OrderItem struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	BookID   int     `json:"book_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type OrderStore struct {
	db *sql.DB
}

func (s *OrderStore) Create(ctx context.Context, order *Order) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		query := `insert into orders ( user_id, total_amount, payment_method ,shipping_address)
		values ($1 , $2 , $3  , $4 ) returning id , placed_at , updated_at;`

		err := tx.QueryRowContext(ctx, query, order.UserID, order.TotalAmount, order.PaymentMethod, order.ShippingAddress).Scan(
			&order.ID,
			&order.PlacedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return err
		}
		for i := range order.Items {
			order.Items[i].OrderID = order.ID
			if err := s.createOrderItem(ctx, tx, &order.Items[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *OrderStore) createOrderItem(ctx context.Context, tx *sql.Tx, orderItem *OrderItem) error {
	query := `insert into order_items ( order_id, book_id, quantity, price)
	values ($1 , $2 , $3 , $4) returning id;`
	err := tx.QueryRowContext(ctx, query, orderItem.OrderID, orderItem.BookID, orderItem.Quantity, orderItem.Price).Scan(
		&orderItem.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *OrderStore) GetByID(ctx context.Context, ID int) (*Order, error) {

	query := `
	SELECT id, user_id, total_amount, status, payment_method, shipping_address, placed_at, updated_at
	FROM orders
	WHERE id = $1;
	`
	order := &Order{}

	err := s.db.QueryRowContext(ctx, query, ID).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Status,
		&order.PaymentMethod,
		&order.ShippingAddress,
		&order.PlacedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	order.Items = []OrderItem{}
	itemsQuery := `
		SELECT id, order_id, book_id, quantity, price
		FROM order_items
		WHERE order_id = $1;
	`
	rows, err := s.db.QueryContext(ctx, itemsQuery, ID)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()
	for rows.Next() {
		var item OrderItem

		if err := rows.Scan(&item.ID, &item.OrderID, &item.BookID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return order, nil
}
