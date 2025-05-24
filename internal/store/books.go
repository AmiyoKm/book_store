package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          string    `json:"isbn"`
	Price         float32   `json:"price"`
	Tags          []string  `json:"tags"`
	Description   string    `json:"description"`
	CoverImageUrl string    `json:"cover_image_url"`
	Pages         int       `json:"pages"`
	Stock         int       `json:"stock"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Version       int       `json:"version"`
}

type BookStore struct {
	db *sql.DB
}

func (s *BookStore) Create(ctx context.Context, book *Book) error {
	query := `insert into books ( title, author, isbn, description, price, stock, tags, pages, cover_image_url) values ($1 , $2 , $3 , $4 , $5 ,$6 ,$7 , $8 , $9) RETURNING id , created_at , updated_at ;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, book.Title, book.Author, book.ISBN, book.Description, book.Price, book.Stock, pq.Array(book.Tags), book.Pages, book.CoverImageUrl).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *BookStore) GetByID(ctx context.Context, bookID int) (*Book, error) {
	query := `select id ,  title , author , isbn , description , price , stock , tags , pages , cover_image_url , created_at , updated_at , version from books where id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()
	book := &Book{}
	err := s.db.QueryRowContext(ctx, query, bookID).Scan(
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
		&book.CreatedAt,
		&book.UpdatedAt,
		&book.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return book, nil
}

func (s *BookStore) Update(ctx context.Context, book *Book) error {
	query := `update books set title=$1 , author=$2 , isbn=$3 , description=$4 , price=$5 , stock=$6 , tags=$7 , pages=$8 , cover_image_url=$9 , version = version+1 where id = 1 and version=$10 RETURNING version;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		book.Title,
		book.Author,
		book.ISBN,
		book.Description,
		book.Price,
		book.Stock,
		pq.Array(book.Tags),
		book.Pages,
		book.CoverImageUrl,
		book.Version,
	).Scan(&book.Version)

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

func (s *BookStore) Delete(ctx context.Context, bookID int) error {
	query := `delete from books where id = $1;`

	_, err := s.db.ExecContext(ctx, query, bookID)

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

type BooksBySearchPayload struct {
	Query    string
	Title    string
	Author   string
	Tags     []string
	MinPrice float32
	MaxPrice float32
	InStock  *bool
}

func (s *BookStore) SearchByBooks(ctx context.Context, filters BooksBySearchPayload) ([]*Book, error) {
	query := `
	SELECT id, title, author, isbn, price, tags, description,
		   cover_image_url, pages, stock, created_at, updated_at, version
	FROM books
	WHERE 1=1
`
	args := []any{}
	argID := 1
	if filters.Query != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR author ILIKE $%d OR description ILIKE $%d)", argID, argID+1, argID+2)
		search := "%" + filters.Query + "%"
		args = append(args, search, search, search)
		argID += 3
	}
	if filters.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argID)
		args = append(args, "%"+filters.Title+"%")
		argID++
	}
	if filters.Author != "" {
		query += fmt.Sprintf(" AND author ILIKE $%d", argID)
		args = append(args, "%"+filters.Author+"%")
		argID++
	}
	if len(filters.Tags) > 0 {
		query += fmt.Sprintf(" AND tags && $%d", argID)
		args = append(args, pq.Array(filters.Tags))
		argID++
	}
	if filters.MinPrice > 0 {
		query += fmt.Sprintf(" AND price >=$%d", argID)
		args = append(args, filters.MinPrice)
		argID++
	}
	if filters.MaxPrice > 0 {
		query += fmt.Sprintf(" AND price <=$%d", argID)
		args = append(args, filters.MaxPrice)
		argID++
	}
	if filters.InStock != nil {
		if *filters.InStock {
			query += " AND stock > 0"
		}
	}
	query += " ORDER BY title ASC"
	rows, err := s.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Book
	for rows.Next() {
		b := &Book{}

		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Author,
			&b.ISBN,
			&b.Price,
			pq.Array(&b.Tags),
			&b.Description,
			&b.CoverImageUrl,
			&b.Pages,
			&b.Stock,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.Version,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil

}
