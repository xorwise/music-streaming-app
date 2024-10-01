package domain

import (
	"context"
	"time"
)

type Room struct {
	ID        int64
	Name      string
	Avatar    string
	Code      string
	OwnerID   int64
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UserRoom struct {
	RoomID int64
	UserID int64
}

type RoomRepository interface {
	Create(ctx context.Context, room *Room) (int64, error)
	GetByID(ctx context.Context, id int64) (*Room, error)
	GetByCode(ctx context.Context, code string) (*Room, error)
	ListByOwnerID(ctx context.Context, ownerID int64) ([]*Room, error)
	ListRoomUsers(ctx context.Context, roomID int64, limit int, offset int) ([]*User, error)
	SetCode(ctx context.Context, roomID int64, code string) error
	AddRoomUser(ctx context.Context, roomID int64, userID int64) error
	GetByUserIDandRoomID(ctx context.Context, id int64, userID int64) (*UserRoom, error)
	ListByUserID(ctx context.Context, userID int64, limit int, offset int) ([]*Room, error)
	RemoveRoomUser(ctx context.Context, roomID int64, userID int64) error
	Update(ctx context.Context, room *Room) error
}
