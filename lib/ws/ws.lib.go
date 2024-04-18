package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"

	"github.com/adhupraba/discord-server/types"
)

type WsClient struct {
	Conn     *websocket.Conn
	ID       string // Profile id
	MemberID string
	RoomID   string // will be empty when user initially establishes websocket connection. it will be updated when user opens a channel or a private conversation
	RoomType types.WsRoomType
	Message  chan *types.WsOutgoingMessage
}

type Hub struct {
	Clients    map[*WsClient]bool
	Register   chan *WsClient
	Unregister chan *WsClient
	Broadcast  chan *types.WsOutgoingMessage
}

var WsHub *Hub

func NewUserHub() {
	WsHub = &Hub{
		Clients:    make(map[*WsClient]bool),
		Register:   make(chan *WsClient),
		Unregister: make(chan *WsClient),
		Broadcast:  make(chan *types.WsOutgoingMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Clients[cl]; !ok {
				h.Clients[cl] = true
			}

		case cl := <-h.Unregister:
			delete(h.Clients, cl)

		case m := <-h.Broadcast:
			for cl := range h.Clients {
				if cl.RoomID == m.Message.RoomId.String() {
					cl.Message <- m
				}
			}
		}
	}
}

func (c *WsClient) WriteMessage() {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Message

		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

// read message from the websocket connection
func (c *WsClient) ReadMessage() {
	defer func() {
		WsHub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket err: %v\n", err)
			}

			break
		}

		var body types.WsIncomingMessage
		err = json.Unmarshal([]byte(m), &body)

		if err != nil {
			log.Println("invalid json message body in websocket =>", err)
			break
		}

		if body.Event == types.WsMessageEventJOINROOM {
			c.RoomID = body.RoomID
			c.RoomType = body.RoomType
		} else if body.Event == types.WsMessageEventNEWMESSAGE {
			_, err := BroadcastMessage(c.MemberID, c.RoomID, c.RoomType, *body.Message)

			if err != nil {
				log.Println("error broadcasting new message =>", err)
				break
			}
		}
	}
}
