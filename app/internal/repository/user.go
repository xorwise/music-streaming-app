package repository

import (
	"context"
	"database/sql"

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
		// TODO: custom errors
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, user.Username, user.PassHash)

	var id int64
	if err := row.Scan(&id); err != nil {
		// TODO: custom errors
		return 0, err
	}
	return id, nil
}

func (ur *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user *domain.User = &domain.User{}
	stmt, err := ur.db.PrepareContext(ctx, "SELECT id, username, password FROM users WHERE id = $1")
	if err != nil {
		// TODO: custom errors
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	if err := row.Scan(&user.ID, &user.Username, &user.PassHash); err != nil {
		// TODO: custom errors
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user *domain.User = &domain.User{}
	stmt, err := ur.db.PrepareContext(ctx, "SELECT id, username, password FROM users WHERE username = $1")
	if err != nil {
		// TODO: custom errors
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, username)
	if err := row.Scan(&user.ID, &user.Username, &user.PassHash); err != nil {
		// TODO: custom errors
		return nil, err
	}
	return user, nil
}
