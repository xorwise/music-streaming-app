package domain

import "context"

type RoomEnterRequest struct {
	Code string `json:"code"`
}

type RoomEnterUsecase interface {
	GetByCode(ctx context.Context, code string) (*Room, error)
	GetByUserIDandRoomID(ctx context.Context, id int64, userID int64) (*UserRoom, error)
	AddRoomUser(ctx context.Context, roomID int64, userID int64) error
}
