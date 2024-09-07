package controller

import (
	"github.com/xorwise/music-streaming-service/internal/bootstrap"
	"golang.org/x/net/websocket"
)

type WSBaseController struct {
	Cfg *bootstrap.Config
}

func (wsc *WSBaseController) Handle(ws *websocket.Conn) {
	var message string
	for {
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			break
		}

		err = websocket.Message.Send(ws, message)
		if err != nil {
			break
		}
	}

}
