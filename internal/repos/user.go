package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andrdru/go-template/internal/entities"
	"github.com/andrdru/go-template/internal/metrics"
	"github.com/andrdru/go-template/internal/tx"
)

type User struct {
	db transactor
}

func NewUser(db *sql.DB) *User {
	return &User{
		db: tx.NewTX(db),
	}
}

// CreateUser .
func (u *User) CreateUser(ctx context.Context, user entities.User) (id int64, err error) {
	start := time.Now()
	defer func() {
		u.handleMetric("user_create", time.Since(start), err)
	}()

	const query = `INSERT INTO users(email, passhash) VALUES($1, $2) RETURNING id`
	err = u.db.DB(ctx).QueryRowContext(ctx, query,
		user.Email,
		user.Passhash,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *User) Session(ctx context.Context, token string) (session entities.Session, err error) {
	start := time.Now()
	defer func() {
		u.handleMetric("session_get", time.Since(start), err)
	}()

	const query = `SELECT id,
       created_at,
       updated_at,
       deleted_at,
       user_id,
       token,
       extra
FROM sessions WHERE token = $1 and deleted_at IS NULL`

	err = u.db.DB(ctx).QueryRowContext(ctx, query, token).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.DeletedAt,
		&session.UserID,
		&session.Token,
		&session.Extra,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Session{}, entities.ErrNotFound
		}

		return entities.Session{}, err
	}

	return session, nil
}

func (u *User) CreateSession(ctx context.Context, session entities.Session) (err error) {
	start := time.Now()
	defer func() {
		u.handleMetric("session_create", time.Since(start), err)
	}()

	const query = `INSERT INTO sessions(user_id, token, extra) VALUES($1, $2, $3)`

	_, err = u.db.DB(ctx).ExecContext(ctx, query,
		session.UserID,
		session.Token,
		session.Extra,
	)

	return err
}

func (u *User) DeleteSession(ctx context.Context, accessToken string) (err error) {
	start := time.Now()
	defer func() {
		u.handleMetric("session_delete", time.Since(start), err)
	}()

	const query = `UPDATE sessions SET deleted_at = now() WHERE token = $1`

	res, err := u.db.DB(ctx).ExecContext(ctx, query, accessToken)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if count == 0 {
		return entities.ErrNotFound
	}

	return nil
}

func (u *User) User(ctx context.Context, email string) (user entities.User, err error) {
	start := time.Now()
	defer func() {
		u.handleMetric("user_get", time.Since(start), err)
	}()

	const query = `SELECT id,
       created_at,
       updated_at,
       deleted_at,
       email,
       passhash
FROM users WHERE email=$1 AND deleted_at IS NULL`

	err = u.db.DB(ctx).QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Email,
		&user.Passhash,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, entities.ErrNotFound
		}

		return entities.User{}, err
	}

	return user, nil
}

func (_ *User) handleMetric(name string, d time.Duration, err error) {
	metrics.HistogramObserverDB("postgres", name, entities.Err(err)).Observe(d.Seconds())
}
