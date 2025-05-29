package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Wishlist struct {
	ID        int
	BookID    int
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
type WishlistStore struct {
	db *sql.DB
}

func (s *WishlistStore) Create(ctx context.Context, wishlist *Wishlist) error {
	query := `INSERT INTO wishlists(book_id , user_id) VALUES($1 , $2) RETURNING id , created_at , updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, wishlist.BookID, wishlist.UserID).Scan(
		&wishlist.ID,
		&wishlist.CreatedAt,
		&wishlist.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *WishlistStore) GetWishlistBooks(ctx context.Context, userID int) ([]*Book, error) {
	query := `SELECT b.id, b.title, b.author, b.isbn, b.description, b.price, b.stock, b.tags, b.pages, b.cover_image_url FROM books b JOIN wishlists w ON w.book_id = b.id
	WHERE w.user_id = $1 ORDER BY w.created_at DESC`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	var books []*Book
	defer rows.Close()

	for rows.Next() {
		book := &Book{}

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Description,
			&book.Price,
			&book.Stock,
			pq.Array(&book.Tags),
			&book.Pages,
			&book.CoverImageUrl,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return books, nil
}
func (s *WishlistStore) Delete(ctx context.Context, userID int, BookID int) error {
	query := `DELETE FROM wishlists WHERE book_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, BookID, userID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrorNotFound
		default:
			return err
		}
	}
	return err
}
