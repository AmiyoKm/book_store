package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QueryTimeDuration = time.Second * 30
	ErrorNotFound     = errors.New("resource not found")
)

type Storage struct {
	Books interface {
		Create(context.Context, *Book) error
		GetByID(context.Context, int) (*Book, error)
		Update(context.Context , *Book) error
		Delete (context.Context , int) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Books: &BookStore{db},
	}
}
