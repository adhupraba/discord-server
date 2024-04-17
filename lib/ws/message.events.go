package ws

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
)

func BroadcastMessage(client *WsClient, body types.WsIncomingMessageBody) error {
	roomId, err := uuid.Parse(client.RoomID)

	if err != nil {
		log.Print("invalid channel uuid")
		return err
	}

	memberId, err := uuid.Parse(client.ID)

	if err != nil {
		log.Print("invalid member uuid")
		return err
	}

	newMessage, err := lib.DB.CreateChannelMessage(context.Background(), model.Messages{
		Content:   body.Content,
		FileUrl:   &body.FileUrl,
		MemberID:  memberId,
		ChannelID: roomId,
		Deleted:   false,
	})

	if err != nil {
		log.Print("failed to save message to db =>", err)
		return err
	}

	WsHub.Broadcast <- &types.WsOutgoingMessage{
		Event:   types.WsMessageEventBROADCAST,
		UserID:  client.ID,
		Message: &newMessage,
	}

	return nil
}
