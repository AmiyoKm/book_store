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
		SearchByBooks(ctx context.Context, filters BooksBySearchPayload) ([]*Book, error)
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
		UpdatePassword(context.Context, int, *Password) error
		MarkPasswordRequestAsUsed(ctx context.Context, hashToken string) error
		Activate(context.Context, string) error
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
	Reviews interface {
		Create(context.Context, *Review) error
		GetByBookID(context.Context, int) ([]*Review, error)
		GetByID(ctx context.Context, reviewID int) (*Review, error)
		Delete(context.Context, int, int) error
		Update(context.Context, *Review) error
	}
	Carts interface {
		GetOrCreateCart(ctx context.Context, userID int) (*Cart, error)
		GetCartItem(ctx context.Context, cartID int) (*CartItem, error)
		InsertOrUpdateCartItem(ctx context.Context, cartID, bookID, quantity int) error
		GetCartItemsWithBooks(ctx context.Context, cartID int) ([]CartItemWithBook, error)
		DeleteCartItem(context.Context, int, int) error
		DeleteCart(context.Context, int) error
		UpdateQuantity(ctx context.Context, quantity int, itemID, userID int) error
	}
	WishLists interface {
		Create(ctx context.Context, wishlist *Wishlist) error
		GetWishlistBooks(context.Context, int) ([]*Book, error)
		Delete(ctx context.Context, userID int, BookID int) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Books:     &BookStore{db},
		Users:     &UserStore{db},
		Roles:     &RoleStore{db},
		Orders:    &OrderStore{db},
		Reviews:   &ReviewStore{db},
		Carts:     &CartStore{db},
		WishLists: &WishlistStore{db},
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
