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
func (s *UserStore) GetByID(ctx context.Context, ID int) (*User, error) {
	query := `select users.id , username , email , password , created_at , roles.* from users
    join roles on (users.role_id = roles.id)
    where users.id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, ID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
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
func (s *UserStore) Update(ctx context.Context, user *User) error {
	query := `UPDATE users
	SET username = $1 WHERE id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	row, err := s.db.ExecContext(ctx, query, user.Username, user.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrorNotFound
		default:
			return err
		}
	}
	affected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrorNotFound
	}
	return nil
}
func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, user); err != nil {
			return err
		}

		if err := s.createAndInvitation(ctx, tx, token, exp, user.ID); err != nil {
			return err
		}
		return nil
	})
}
func (s *UserStore) Delete(ctx context.Context, userID int) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}
		if err := s.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int) error {
	query := `delete from user_invitations where user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, userID int) error {
	query := `delete from users where id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) createAndInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int) error {

	query := `INSERT INTO user_invitations (token , user_id , expiry)
	VALUES ($1 , $2 , $3)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}
	return nil
}
