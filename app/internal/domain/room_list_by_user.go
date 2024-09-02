package domain

import "context"

type RoomListByUserUsecase interface {
	ListByUser(ctx context.Context, userID int64, limit int, offset int) ([]*Room, error)
}

type RoomListByUserResponse struct {
	Rooms []*RoomCreateResponse `json:"rooms,omitempty"`
}
