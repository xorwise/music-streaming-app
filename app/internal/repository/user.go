package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) (int64, error) {
	stmt, err := ur.db.PrepareContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, user.Username, user.PassHash)

	var id int64
	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, domain.ErrUserAlreadyExists
		} else if pgErr.Code == "23502" {
			return 0, domain.ErrFieldRequired
		}
		return 0, err
	}
	return id, nil
}

func (ur *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user *domain.User = &domain.User{}
	stmt, err := ur.db.PrepareContext(ctx, "SELECT id, username, avatar_path, password FROM users WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	if err := row.Scan(&user.ID, &user.Username, &user.Avatar, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user *domain.User = &domain.User{}
	stmt, err := ur.db.PrepareContext(ctx, "SELECT id, username, avatar_path, password FROM users WHERE username = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, username)
	if err := row.Scan(&user.ID, &user.Username, &user.Avatar, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) Update(ctx context.Context, user *domain.User) error {
	stmt, err := ur.db.PrepareContext(ctx, "UPDATE users SET avatar_path = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, user.Avatar, user.ID)
	return err
}
