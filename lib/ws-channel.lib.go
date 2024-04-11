package lib

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/types"
)

type ChannelMember struct {
	Conn      *websocket.Conn
	Message   chan *types.MessageWithMember
	ID        string `json:"id"` // channel member id
	ChannelID string `json:"channelId"`
	ProfileID string `json:"profileId"`
}

type Channel struct {
	ID      string `json:"id"`
	Members map[string]*ChannelMember
}

type ChannelHub struct {
	Channels   map[string]*Channel
	Register   chan *ChannelMember
	Unregister chan *ChannelMember
	Broadcast  chan *types.MessageWithMember
}

var HubChannel *ChannelHub

func NewChannelHub() {
	HubChannel = &ChannelHub{
		Channels:   make(map[string]*Channel),
		Register:   make(chan *ChannelMember),
		Unregister: make(chan *ChannelMember),
		Broadcast:  make(chan *types.MessageWithMember),
	}
}

func (h *ChannelHub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if ch, ok := h.Channels[cl.ChannelID]; ok {
				if _, ok := ch.Members[cl.ID]; !ok {
					ch.Members[cl.ID] = cl
				}
			}

		case cl := <-h.Unregister:
			if ch, ok := h.Channels[cl.ChannelID]; ok {
				if _, ok := ch.Members[cl.ID]; ok {
					delete(ch.Members, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if ch, ok := h.Channels[m.ChannelID.String()]; ok {
				for _, cl := range ch.Members {
					cl.Message <- m
				}
			}
		}
	}
}

// write message to the websocket connection
func (c *ChannelMember) WriteMessage() {
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
func (c *ChannelMember) ReadMessage() {
	defer func() {
		HubChannel.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket err: %v", err)
			}

			break
		}

		var body types.SendMessageBody
		err = json.Unmarshal([]byte(m), &body)

		if err != nil {
			log.Print("invalid json message body in websocket =>", err)
			break
		}

		memberId, err := uuid.Parse(c.ID)

		if err != nil {
			log.Print("invalid uuid member id")
			break
		}

		channelId, err := uuid.Parse(c.ChannelID)

		if err != nil {
			log.Print("invalid uuid channel id")
			break
		}

		newMessage, err := DB.CreateChannelMessage(context.Background(), model.Messages{
			Content:   body.Content,
			FileUrl:   &body.FileUrl,
			MemberID:  memberId,
			ChannelID: channelId,
			Deleted:   false,
		})

		if err != nil {
			log.Print("failed to save message to db =>", err)
			break
		}

		HubChannel.Broadcast <- &newMessage
	}
}
