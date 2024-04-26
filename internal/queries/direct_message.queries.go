package queries

import (
	"context"
	"fmt"
	"log"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/internal/helpers"
	"github.com/adhupraba/discord-server/types"
)

func transformDirectMsgToWsMsg(message model.DirectMessages) types.WsMessage {
	return types.WsMessage{
		ID:        message.ID,
		Content:   message.Content,
		FileUrl:   message.FileURL,
		MemberID:  message.MemberID,
		RoomId:    message.ConversationID,
		Deleted:   message.Deleted,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

type GetDirectMessagesParams struct {
	ConversationID  uuid.UUID
	LastMessageId   *uuid.UUID
	LastMessageDate *time.Time
}

func (q *Queries) GetDirectMessages(ctx context.Context, params GetDirectMessagesParams) (messages []types.WsMessageContent, nextCursor *string, err error) {
	var MESSAGES_BATCH = 10
	var exp BoolExpression = DirectMessages.ConversationID.EQ(UUID(params.ConversationID))
	var offset int64 = 0

	if params.LastMessageId != nil && params.LastMessageDate != nil {
		offset = 1

		exp = exp.AND(
			DirectMessages.CreatedAt.LT_EQ(TimestampzT(*params.LastMessageDate)).OR(
				DirectMessages.CreatedAt.EQ(TimestampzT(*params.LastMessageDate)).AND(DirectMessages.ID.EQ(UUID(*params.LastMessageId))),
			),
		)
	}

	stmt := SELECT(DirectMessages.AllColumns, Members.AllColumns, Profiles.AllColumns).
		FROM(
			DirectMessages.
				LEFT_JOIN(Members, Members.ID.EQ(DirectMessages.MemberID)).
				LEFT_JOIN(Profiles, Profiles.ID.EQ(Members.ProfileID)),
		).
		WHERE(exp).
		ORDER_BY(DirectMessages.CreatedAt.DESC()).
		LIMIT(int64(MESSAGES_BATCH)).
		OFFSET(offset)

	var dbMessages []types.DbDirectMessageWithMember

	err = stmt.QueryContext(ctx, q.db, &dbMessages)

	if err != nil {
		return []types.WsMessageContent{}, nil, err
	}

	wsMessages := []types.WsMessageContent{}

	for _, msg := range dbMessages {
		wsMessages = append(wsMessages, types.WsMessageContent{
			WsMessage: transformDirectMsgToWsMsg(msg.DirectMessages),
			Member:    msg.Member,
		})
	}

	if len(wsMessages) == MESSAGES_BATCH {
		last := wsMessages[len(wsMessages)-1]
		combined := fmt.Sprintf("%v&%v", last.ID.String(), last.CreatedAt.Format(time.RFC3339))
		next := helpers.Base64Encode(combined)
		nextCursor = &next
	}

	return wsMessages, nextCursor, nil
}

func (q *Queries) GetDirectMessageByID(ctx context.Context, messageId uuid.UUID) (message *model.DirectMessages, err error) {
	stmt := DirectMessages.SELECT(DirectMessages.AllColumns).WHERE(DirectMessages.ID.EQ(UUID(messageId)))

	var msg model.DirectMessages

	err = stmt.QueryContext(ctx, q.db, &msg)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (q *Queries) CreateDirectMessage(ctx context.Context, data model.DirectMessages) (types.WsMessageContent, error) {
	stmt := DirectMessages.INSERT(
		DirectMessages.AllColumns.Except(DirectMessages.CreatedAt, DirectMessages.UpdatedAt),
	).
		MODEL(data).
		RETURNING(DirectMessages.AllColumns)

	var message model.DirectMessages
	err := stmt.QueryContext(ctx, q.db, &message)

	if err != nil {
		return types.WsMessageContent{}, err
	}

	member, err := q.GetMemberWithProfileByMemberID(ctx, message.MemberID)

	if err != nil {
		return types.WsMessageContent{}, err
	}

	messageWithMember := types.WsMessageContent{
		WsMessage: transformDirectMsgToWsMsg(message),
		Member:    member,
	}

	return messageWithMember, err
}

func (q *Queries) UpdateDirectMessageByID(ctx context.Context, id uuid.UUID, content string) (message *types.WsOutgoingMessage, err error) {
	stmt := DirectMessages.UPDATE(DirectMessages.Content).
		MODEL(model.DirectMessages{Content: content}).
		WHERE(DirectMessages.ID.EQ(UUID(id))).
		RETURNING(DirectMessages.AllColumns)

	var updMessage model.DirectMessages
	err = stmt.QueryContext(ctx, q.db, &updMessage)

	if err != nil {
		log.Println("update direct message content error =>", err)
		return nil, err
	}

	member, err := q.GetMemberWithProfileByMemberID(ctx, updMessage.MemberID)

	if err != nil {
		log.Println("get member with profile error =>", err)
		return nil, err
	}

	message = &types.WsOutgoingMessage{
		Event: types.WsMessageEventMESSAGEMODIFIED,
		Message: &types.WsMessageContent{
			WsMessage: transformDirectMsgToWsMsg(updMessage),
			Member:    member,
		},
	}

	return message, err
}

func (q *Queries) DeleteDirectMessageByID(ctx context.Context, id uuid.UUID) (message *types.WsOutgoingMessage, err error) {
	stmt := DirectMessages.UPDATE(DirectMessages.FileURL, DirectMessages.Content, DirectMessages.Deleted).
		MODEL(model.DirectMessages{FileURL: nil, Content: "This message has been deleted.", Deleted: true}).
		WHERE(DirectMessages.ID.EQ(UUID(id))).
		RETURNING(DirectMessages.AllColumns)

	var updMessage model.DirectMessages
	err = stmt.QueryContext(ctx, q.db, &updMessage)

	if err != nil {
		return nil, err
	}

	member, err := q.GetMemberWithProfileByMemberID(ctx, updMessage.MemberID)

	if err != nil {
		return nil, err
	}

	message = &types.WsOutgoingMessage{
		Event: types.WsMessageEventMESSAGEMODIFIED,
		Message: &types.WsMessageContent{
			WsMessage: transformDirectMsgToWsMsg(updMessage),
			Member:    member,
		},
	}

	return message, err
}
