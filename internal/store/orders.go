package store

import (
	"context"
	"database/sql"
	"fmt"
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
func (s *OrderStore) Update(ctx context.Context, order *Order) error {
	query := `UPDATE orders
SET
    shipping_address = $1,
    payment_method = $2,
    status = $3
WHERE
    id = $4 AND user_id = $5;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, order.ShippingAddress, order.PaymentMethod, order.Status, order.ID, order.UserID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrorNotFound
		default:
			return err
		}
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}
func (s *OrderStore) Get(ctx context.Context, userID int) ([]Order, error) {
	query := `
		SELECT o.id, o.user_id, o.total_amount, o.status, o.payment_method,
		o.shipping_address, o.placed_at, o.updated_at,
		oi.id, oi.order_id, oi.book_id, oi.quantity, oi.price
		FROM orders o
		LEFT JOIN order_items oi ON o.id = oi.order_id
		WHERE o.user_id = $1
		ORDER BY o.id, oi.id;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[int]*Order)

	for rows.Next() {
		var order Order

		// Nullable fields for order items
		var itemID sql.NullInt64
		var itemOrderID sql.NullInt64
		var itemBookID sql.NullInt64
		var itemQuantity sql.NullInt64
		var itemPrice sql.NullFloat64

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.Status,
			&order.PaymentMethod,
			&order.ShippingAddress,
			&order.PlacedAt,
			&order.UpdatedAt,
			&itemID,
			&itemOrderID,
			&itemBookID,
			&itemQuantity,
			&itemPrice,
		)
		if err != nil {
			return nil, err
		}

		if existingOrder, exists := ordersMap[order.ID]; exists {
			if itemID.Valid {
				existingOrder.Items = append(existingOrder.Items, OrderItem{
					ID:       int(itemID.Int64),
					OrderID:  int(itemOrderID.Int64),
					BookID:   int(itemBookID.Int64),
					Quantity: int(itemQuantity.Int64),
					Price:    itemPrice.Float64,
				})
			}
		} else {
			if itemID.Valid {
				order.Items = append(order.Items, OrderItem{
					ID:       int(itemID.Int64),
					OrderID:  int(itemOrderID.Int64),
					BookID:   int(itemBookID.Int64),
					Quantity: int(itemQuantity.Int64),
					Price:    itemPrice.Float64,
				})
			}
			ordersMap[order.ID] = &order
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Directly return the list of orders from the map
	finalOrders := make([]Order, 0, len(ordersMap))
	for _, orderPtr := range ordersMap {
		finalOrders = append(finalOrders, *orderPtr)
	}

	return finalOrders, nil
}
