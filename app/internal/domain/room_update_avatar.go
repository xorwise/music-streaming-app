package domain

import "context"

type RoomUpdateAvatarUsecase interface {
	GetByID(ctx context.Context, id int64) (*Room, error)
	GetByUserIDandRoomID(ctx context.Context, userID int64, roomID int64) (*UserRoom, error)
	Update(ctx context.Context, room *Room) error
	SaveFile(ctx context.Context, fileData string, filename string) (string, error)
}
