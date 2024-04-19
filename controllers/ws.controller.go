package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
	"github.com/adhupraba/discord-server/utils"
)

type WsController struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wc *WsController) Connect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connect endpoint reached")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error establishing websocket channel connection")
		return
	}

	fmt.Println("websocket conn established")

	mem := &lib.WsClient{
		Conn:    conn,
		Message: make(chan *types.WsOutgoingMessage),
	}

	lib.WsHub.Register <- mem

	go mem.WriteMessage()

	mem.ReadMessage()
}
