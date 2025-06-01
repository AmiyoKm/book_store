package store

import (
	"context"
	"database/sql"
	"time"
)

type Cart struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Items     []CartItem `json:"items"`
}

type CartItem struct {
	ID        int       `json:"id"`
	CartID    int       `json:"cart_id"`
	BookID    int       `json:"book_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CartItemWithBook struct {
	ID            int       `json:"id"`
	CartID        int       `json:"cart_id"`
	BookID        int       `json:"book_id"`
	Quantity      int       `json:"quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Price         float64   `json:"price"`
	CoverImageUrl string    `json:"cover_image_url"`
	Stock         int       `json:"stock"`
}
type CartStore struct {
	db *sql.DB
}

func (s *CartStore) GetOrCreateCart(ctx context.Context, userID int) (*Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	cart := &Cart{}
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			insertQuery := `INSERT INTO carts (user_id) VALUES ($1) RETURNING id, user_id, created_at, updated_at`
			err = s.db.QueryRowContext(ctx, insertQuery, userID).Scan(
				&cart.ID,
				&cart.UserID,
				&cart.CreatedAt,
				&cart.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return cart, nil
}

func (s *CartStore) GetCartItem(ctx context.Context, cartID int) (*CartItem, error) {
	query := `SELECT id, cart_id, book_id, quantity, created_at, updated_at
	          FROM cart_items WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	var item CartItem
	err := s.db.QueryRowContext(ctx, query, cartID).Scan(
		&item.ID,
		&item.CartID,
		&item.BookID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
func (s *CartStore) InsertOrUpdateCartItem(ctx context.Context, cartID, bookID, quantity int) error {
	query := `INSERT INTO cart_items (cart_id , book_id , quantity) VALUES ($1 , $2 , $3 ) ON CONFLICT (cart_id , book_id)
	DO UPDATE
	SET quantity = cart_items.quantity + EXCLUDED.quantity , updated_at = CURRENT_TIMESTAMP
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, cartID, bookID, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *CartStore) GetCartItemsWithBooks(ctx context.Context, cartID int) ([]CartItemWithBook, error) {
	query := `
	SELECT
	ci.id, ci.cart_id, ci.book_id, ci.quantity, ci.created_at, ci.updated_at,
	b.title, b.author, b.price , b.cover_image_url , b.stock
	FROM cart_items ci
	JOIN books b ON ci.book_id = b.id
	WHERE ci.cart_id = $1 ORDER BY ci.created_at ASC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItemWithBook
	for rows.Next() {
		var item CartItemWithBook
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.BookID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Title,
			&item.Author,
			&item.Price,
			&item.CoverImageUrl,
			&item.Stock,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *CartStore) DeleteCartItem(ctx context.Context, userID int, itemID int) error {
	query := `DELETE FROM cart_items ci USING carts c
			WHERE ci.cart_id = c.id
			AND c.user_id = $1 AND ci.id= $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, itemID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *CartStore) DeleteCart(ctx context.Context, userID int) error {
	query := `DELETE FROM carts WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *CartStore) UpdateQuantity(ctx context.Context, quantity int, itemID, userID int) error {
	query := `
	UPDATE cart_items ci
	SET quantity = $1
	FROM carts c
	WHERE c.id = ci.cart_id
  	AND ci.id = $2
  	AND c.user_id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, quantity, itemID, userID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}
