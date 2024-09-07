package controller

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"github.com/xorwise/music-streaming-service/internal/domain"
	"golang.org/x/net/websocket"
)

type WSRoomController struct {
	Usecase domain.WSRoomUsecase
	Cfg     *bootstrap.Config
	Clients *domain.WSClients
}

func (wsc *WSRoomController) Handle(ws *websocket.Conn) {
	ctx := ws.Request().Context()
	user := ctx.Value("user").(*domain.User)

	roomID, err := strconv.Atoi(ws.Request().URL.Query().Get("id"))
	if err != nil {
		json.NewEncoder(ws).Encode(domain.ErrorResponse{Error: err.Error()})
		ws.Close()
		return
	}
	id := int64(roomID)

	room, err := wsc.Usecase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrRoomNotFound) {
			json.NewEncoder(ws).Encode(domain.ErrorResponse{Error: err.Error()})
			ws.Close()
		} else {
			ws.Close()
		}
		return
	}

	_, err = wsc.Usecase.GetUserIDandRoomID(ctx, room.ID, user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotUserInRoom) {
			json.NewEncoder(ws).Encode(domain.ErrorResponse{Error: err.Error()})
			ws.Close()
		} else {
			ws.Close()
		}
		return
	}
	defer func() {
		ws.Close()
		wsc.Usecase.LoggedOut(ctx, room.ID, user.ID)
		wsc.Clients.Remove(room.ID, user.ID)
	}()

	wsc.Clients.Add(id, user.ID, ws)
	wsc.Usecase.LoggedIn(ctx, room.ID, user.ID)

	wsc.Usecase.Handle(ws, room, user)
}
