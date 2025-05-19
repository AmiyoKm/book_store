package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QueryTimeDuration    = time.Second * 30
	ErrorNotFound        = errors.New("resource not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

type Storage struct {
	Books interface {
		Create(context.Context, *Book) error
		GetByID(context.Context, int) (*Book, error)
		Update(context.Context, *Book) error
		Delete(context.Context, int) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByEmail(context.Context, string) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Delete(context.Context, int) error
		GetByID(ctx context.Context, ID int) (*User, error)
		Update(context.Context, *User) error
		CreatePasswordRequest(context.Context, *User, string, time.Duration) (*int, error)
		DeletePasswordRequest(context.Context, int) error
		GetPasswordRequest(context.Context, string) (*PasswordChangeRequest, error)
	}
	Roles interface {
		GetByName(ctx context.Context, name string) (*Role, error)
	}
	Orders interface {
		GetByID(context.Context, int) (*Order, error)
		Create(ctx context.Context, order *Order) error
		Get(ctx context.Context, userID int) ([]Order, error)
		Update(context.Context, *Order) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Books:  &BookStore{db},
		Users:  &UserStore{db},
		Roles:  &RoleStore{db},
		Orders: &OrderStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
