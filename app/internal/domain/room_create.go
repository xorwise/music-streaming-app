package domain

import (
	"context"
	"time"
)

type RoomCreateRequest struct {
	Name string `json:"name"`
}

type RoomCreateResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	OwnerID   int64      `json:"owner_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type RoomCreateUsecase interface {
	Create(ctx context.Context, room *Room) (int64, error)
	GenerateCode(ctx context.Context, roomID int64) string
	SetCode(ctx context.Context, roomID int64, code string) error
	AddRoomUser(ctx context.Context, roomID int64, userID int64) error
}
