package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/xorwise/music-streaming-service/internal/config"
	"github.com/xorwise/music-streaming-service/internal/database"
	"github.com/xorwise/music-streaming-service/internal/domain/models"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Config) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_DATABASE))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, username string, passHash []byte) (int64, error) {
	const op = "storage.postgresql.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(username, pass_hash) VALUES($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, username, string(passHash[:]))
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Class() == "23000" {
			return 0, fmt.Errorf("%s: %w", op, database.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.LastInsertId()
}

func (s *Storage) User(ctx context.Context, username string) (models.User, error) {
	const op = "storage.postgresql.User"

	stmt, err := s.db.Prepare("SELECT id, username, pass_hash FROM users where username = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, username)

	var user models.User
	err = row.Scan(&user.ID, &user.Username, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, database.ErrNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil

}
