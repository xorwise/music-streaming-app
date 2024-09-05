package domain

import "context"

type RoomGetByIDUsecase interface {
	GetByID(ctx context.Context, id int64) (*Room, error)
	GetUserIDandRoomID(ctx context.Context, id int64, userID int64) (*UserRoom, error)
}

type RoomGetByIDResponse RoomCreateResponse
