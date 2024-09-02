package domain

import (
	"context"
	"time"
)

type Room struct {
	ID        int64
	Name      string
	Code      string
	OwnerID   int64
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type RoomRepository interface {
	Create(ctx context.Context, room *Room) (int64, error)
	GetByID(ctx context.Context, id int64) (*Room, error)
	GetByCode(ctx context.Context, code string) (*Room, error)
	ListByOwnerID(ctx context.Context, ownerID int64) ([]*Room, error)
	ListRoomUsers(ctx context.Context, roomID int64) ([]*User, error)
	SetCode(ctx context.Context, roomID int64, code string) error
	AddRoomUser(ctx context.Context, roomID int64, userID int64) error
}
