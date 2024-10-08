package domain

import "context"

type MessageBrokerUtils interface {
	SubscribeToNats() error
	GetClientsInRoom(ctx context.Context, roomID int64) ([]int64, error)
	HandleRoomClientRequests()
	BroadcastMessage(roomID int64, msg *WSRoomResponse) error
}
