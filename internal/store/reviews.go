package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Review struct {
	ID        int
	UserID    int
	BookID    int
	Content   string
	Rating    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReviewStore struct {
	db *sql.DB
}

func (s *ReviewStore) Create(ctx context.Context, review *Review) error {
	query := `INSERT INTO reviews(user_id , book_id , rating, content) VALUES(
		$1,$2,$3,$4
	) RETURNING id , created_at , updated_at;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, review.UserID, review.BookID, review.Rating, review.Content).Scan(
		&review.ID,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *ReviewStore) GetByBookID(ctx context.Context, bookID int) ([]*Review, error) {
	query := `SELECT id , user_id , book_id , rating , content , created_at , updated_at FROM reviews WHERE book_id = $1 ORDER BY created_at ASC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, bookID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()
	var reviews []*Review
	for rows.Next() {
		review := &Review{}

		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.BookID,
			&review.Rating,
			&review.Content,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}
func (s *ReviewStore) GetByID(ctx context.Context, reviewID int) (*Review, error) {
	query := `SELECT id , user_id , book_id , rating , content , created_at , updated_at FROM reviews WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()
	review := &Review{}
	err := s.db.QueryRowContext(ctx, query, reviewID).Scan(
		&review.ID,
		&review.UserID,
		&review.BookID,
		&review.Rating,
		&review.Content,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return review, nil

}

func (s *ReviewStore) Delete(ctx context.Context, reviewID int, userID int) error {
	query := `delete from reviews where id = $1 and user_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, reviewID, userID)
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
		return ErrorNotFound
	}
	return nil
}

func (s *ReviewStore) Update(ctx context.Context, review *Review) error {
	query := `
		UPDATE reviews
		SET content = $1, rating = $2
		WHERE id = $3 AND user_id = $4
		RETURNING created_at, updated_at;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		review.Content,
		review.Rating,
		review.ID,
		review.UserID,
	).Scan(&review.CreatedAt, &review.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNotFound
		}
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}
