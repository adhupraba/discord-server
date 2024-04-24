package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/types"
)

type WsClient struct {
	Conn     *websocket.Conn
	MemberID string
	RoomID   string // will be empty when user initially establishes websocket connection. it will be updated when user opens a channel or a private conversation
	RoomType types.WsRoomType
	Message  chan *types.WsOutgoingMessage
}

type Hub struct {
	Clients    map[*websocket.Conn]*WsClient
	Register   chan *WsClient
	Unregister chan *WsClient
	Broadcast  chan *types.WsOutgoingMessage
}

var WsHub *Hub

func NewHub() {
	WsHub = &Hub{
		Clients:    make(map[*websocket.Conn]*WsClient),
		Register:   make(chan *WsClient),
		Unregister: make(chan *WsClient),
		Broadcast:  make(chan *types.WsOutgoingMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Clients[cl.Conn]; !ok {
				h.Clients[cl.Conn] = cl
			}

		case cl := <-h.Unregister:
			delete(h.Clients, cl.Conn)

		case m := <-h.Broadcast:
			for conn := range h.Clients {
				cl := h.Clients[conn]

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

		if body.Event == types.WsMessageEventAUTHENTICATE && body.AuthToken != "" {
			sessClaims, err := ClerkClient.VerifyToken(body.AuthToken)

			if err != nil {
				log.Println("invalid auth token =>", err)
				break
			}

			_, err = ClerkClient.Users().Read(sessClaims.Subject)

			if err != nil {
				fmt.Println("clerk user read error =>", err)
				break
			}

			go func() {
				c.Message <- &types.WsOutgoingMessage{
					Event:   types.WsMessageEventACKNOWLEDGED,
					Message: nil,
				}
			}()
		} else if body.Event == types.WsMessageEventJOINROOM {
			c.MemberID = body.MemberID
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

func BroadcastMessage(member_id string, room_id string, room_type types.WsRoomType, body types.WsIncomingMessageBody) (*types.WsOutgoingMessage, error) {
	roomId, err := uuid.Parse(room_id)

	if err != nil {
		log.Print("invalid channel uuid")
		return nil, err
	}

	memberId, err := uuid.Parse(member_id)

	if err != nil {
		log.Print("invalid member uuid")
		return nil, err
	}

	var newMessage types.WsMessageContent

	if room_type == types.WsRoomTypeCHANNEL {
		newMessage, err = DB.CreateChannelMessage(context.Background(), model.Messages{
			ID:        uuid.New(),
			Content:   body.Content,
			FileURL:   body.FileUrl,
			MemberID:  memberId,
			ChannelID: roomId,
			Deleted:   false,
		})
	} else {
		// ! save message to conversation table
	}

	if err != nil {
		log.Print("failed to save message to db =>", err)
		return nil, err
	}

	msg := &types.WsOutgoingMessage{
		Event:   types.WsMessageEventBROADCAST,
		Message: &newMessage,
	}

	go func() {
		WsHub.Broadcast <- msg
	}()

	return msg, nil
}
