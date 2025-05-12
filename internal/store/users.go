package store

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	IsActive  bool      `json:"is_active"`
	RoleID    int       `json:"role_id"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore struct {
	db *sql.DB
}
type Password struct {
	Text *string
	Hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Text = &text
	p.Hash = hash
	return nil
}
func (p *Password) ComparePassword(pass string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(pass))
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (username , password , email , role_id)
	VALUES($1  , $2 , $3 , (SELECT id FROM roles WHERE name = $4)) RETURNING id , created_at , updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, user.Username, user.Password.Hash, user.Email, user.Role.Name).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `select id , username , email , password , created_at ,role_id from users where email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
		&user.RoleID,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
