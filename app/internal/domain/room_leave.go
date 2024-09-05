package domain

import "context"

type RoomLeaveUsecase interface {
	GetByID(ctx context.Context, id int64) (*Room, error)
	GetUserIDandRoomID(ctx context.Context, id int64, userID int64) (*UserRoom, error)
	RemoveRoomUser(ctx context.Context, roomID int64, userID int64) error
}
