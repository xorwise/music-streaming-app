package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type natsUtils struct {
	conn      *nats.Conn
	broadcast chan *domain.RoomBroadcastResponse
	wsh       domain.WebSocketHandler
}

func NewNatsUtils(conn *nats.Conn, broadcast chan *domain.RoomBroadcastResponse, wsh domain.WebSocketHandler) domain.MessageBrokerUtils {
	return &natsUtils{
		conn:      conn,
		broadcast: broadcast,
		wsh:       wsh,
	}
}

func (nu *natsUtils) SubscribeToNats() error {
	_, err := nu.conn.Subscribe("broadcast.*", func(msg *nats.Msg) {
		roomID, err := strconv.Atoi(msg.Subject[len("broadcast."):])
		if err != nil {
			return
		}
		var response domain.WSRoomResponse
		err = json.Unmarshal(msg.Data, &response)
		if err != nil {
			return
		}

		nu.broadcast <- &domain.RoomBroadcastResponse{
			RoomID:   int64(roomID),
			Response: &response,
		}
	})
	return err
}

func (nu *natsUtils) BroadcastMessage(roomID int64, msg *domain.WSRoomResponse) error {
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	err = nu.conn.Publish("broadcast."+fmt.Sprintf("%d", roomID), byteMsg)
	if err != nil {
		return err
	}
	return nil
}

func (nu *natsUtils) GetClientsInRoom(ctx context.Context, roomID int64) ([]int64, error) {
	var clientsInRoom []int64
	responseCh := make(chan []int64)

	subject := "room_clients." + fmt.Sprintf("%d", roomID)

	inbox := nu.conn.NewRespInbox()

	sub, err := nu.conn.SubscribeSync(inbox)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	err = nu.conn.PublishRequest(subject, inbox, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg, err := sub.NextMsgWithContext(ctx)
			if err != nil {
				break
			}

			var clients []int64
			err = json.Unmarshal(msg.Data, &clients)
			if err == nil {
				responseCh <- clients
			}
		}
		close(responseCh)
	}()

	for clients := range responseCh {
		clientsInRoom = append(clientsInRoom, clients...)
	}

	return clientsInRoom, nil
}

func (nu *natsUtils) HandleRoomClientRequests() {
	nu.conn.Subscribe("room_clients.*", func(msg *nats.Msg) {
		roomID, err := strconv.Atoi(msg.Subject[len("room_clients."):])
		if err != nil {
			return
		}

		localClients := nu.wsh.BroadcastClients(int64(roomID))
		if localClients != nil {
			data, _ := json.Marshal(localClients)
			msg.Respond(data)
		}
	})
}
