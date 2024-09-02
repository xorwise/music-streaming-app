package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type roomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) domain.RoomRepository {
	return &roomRepository{
		db: db,
	}
}

func (rr *roomRepository) Create(ctx context.Context, room *domain.Room) (int64, error) {
	stmt, err := rr.db.PrepareContext(ctx, "INSERT INTO rooms (name, code, owner_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRowContext(ctx, room.Name, room.Code, room.OwnerID, room.CreatedAt).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23502" {
			return 0, domain.ErrFieldRequired
		}
		return 0, err
	}
	return id, nil
}

func (rr *roomRepository) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	var room *domain.Room = &domain.Room{}
	stmt, err := rr.db.PrepareContext(ctx, "SELECT id, name, code, owner_id, created_at, updated_at FROM rooms WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	if err := row.Scan(&room.ID, &room.Name, &room.Code, &room.OwnerID, &room.CreatedAt, &room.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRoomNotFound
		}
		return nil, err
	}
	return room, nil
}

func (rr *roomRepository) GetByCode(ctx context.Context, code string) (*domain.Room, error) {
	var room *domain.Room = &domain.Room{}
	stmt, err := rr.db.PrepareContext(ctx, "SELECT id, name, code, owner_id, created_at, updated_at FROM rooms WHERE code = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, code)
	if err := row.Scan(&room.ID, &room.Name, &room.Code, &room.OwnerID, &room.CreatedAt, &room.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRoomNotFound
		}
		return nil, err
	}
	return room, nil
}

func (rr *roomRepository) ListByOwnerID(ctx context.Context, ownerID int64) ([]*domain.Room, error) {
	var rooms []*domain.Room
	stmt, err := rr.db.PrepareContext(ctx, "SELECT id, name, code, owner_id, created_at, updated_at FROM rooms WHERE owner_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var room *domain.Room = &domain.Room{}
		if err := rows.Scan(&room.ID, &room.Name, &room.Code, &room.OwnerID, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (rr *roomRepository) ListRoomUsers(ctx context.Context, roomID int64) ([]*domain.User, error) {
	var users []*domain.User
	stmt, err := rr.db.PrepareContext(ctx, "SELECT id, username FROM users WHERE id IN (SELECT user_id FROM users_rooms WHERE room_id = $1)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, roomID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user *domain.User = &domain.User{}
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (rr *roomRepository) SetCode(ctx context.Context, roomID int64, code string) error {
	stmt, err := rr.db.PrepareContext(ctx, "UPDATE rooms SET code = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, code, roomID)
	return err
}

func (rr *roomRepository) AddRoomUser(ctx context.Context, roomID int64, userID int64) error {
	stmt, err := rr.db.PrepareContext(ctx, "INSERT INTO users_rooms (user_id, room_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, userID, roomID)
	return err
}
