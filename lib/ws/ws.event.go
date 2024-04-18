package ws

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
)

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
		newMessage, err = lib.DB.CreateChannelMessage(context.Background(), model.Messages{
			Content:   body.Content,
			FileUrl:   &body.FileUrl,
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

	WsHub.Broadcast <- msg

	return msg, nil
}
