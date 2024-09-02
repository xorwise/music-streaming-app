package domain

import "context"

type RoomUsersResponse struct {
	Users []UserMeResponse `json:"users,omitempty"`
}

type RoomUsersUsecase interface {
	ListRoomUsers(ctx context.Context, roomID int64, limit int, offset int) ([]*User, error)
	GetByID(ctx context.Context, id int64) (*Room, error)
	GetByUserIDandRoomID(ctx context.Context, id int64, userID int64) (*UserRoom, error)
}
