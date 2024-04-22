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

func transformMsgToWsMsg(message model.Messages) types.WsMessage {
	return types.WsMessage{
		ID:        message.ID,
		Content:   message.Content,
		FileUrl:   message.FileURL,
		MemberID:  message.MemberID,
		RoomId:    message.ChannelID,
		Deleted:   message.Deleted,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

func (q *Queries) CreateChannelMessage(ctx context.Context, data model.Messages) (types.WsMessageContent, error) {
	stmt := Messages.INSERT(
		Messages.AllColumns.Except(Messages.CreatedAt, Messages.UpdatedAt),
	).
		MODEL(data).
		RETURNING(Messages.AllColumns)

	var message model.Messages
	err := stmt.QueryContext(ctx, q.db, &message)

	if err != nil {
		return types.WsMessageContent{}, err
	}

	member, err := q.GetMemberWithProfileByMemberID(ctx, message.MemberID)

	if err != nil {
		return types.WsMessageContent{}, err
	}

	messageWithMember := types.WsMessageContent{
		WsMessage: transformMsgToWsMsg(message),
		Member:    member,
	}

	return messageWithMember, err
}

type GetMessagesParams struct {
	ChannelId       uuid.UUID
	LastMessageId   *uuid.UUID
	LastMessageDate *time.Time
}

func (q *Queries) GetMessages(ctx context.Context, params GetMessagesParams) (messages []types.WsMessageContent, nextCursor *string, err error) {
	var MESSAGES_BATCH = 10
	var exp BoolExpression = Messages.ChannelID.EQ(UUID(params.ChannelId))
	var offset int64 = 0

	if params.LastMessageId != nil && params.LastMessageDate != nil {
		offset = 1

		exp = exp.AND(
			Messages.CreatedAt.LT_EQ(TimestampzT(*params.LastMessageDate)).OR(
				Messages.CreatedAt.EQ(TimestampzT(*params.LastMessageDate)).AND(Messages.ID.EQ(UUID(*params.LastMessageId))),
			),
		)
	}

	stmt := SELECT(Messages.AllColumns, Members.AllColumns, Profiles.AllColumns).
		FROM(
			Messages.
				LEFT_JOIN(Members, Members.ID.EQ(Messages.MemberID)).
				LEFT_JOIN(Profiles, Profiles.ID.EQ(Members.ProfileID)),
		).
		WHERE(exp).
		ORDER_BY(Messages.CreatedAt.DESC()).
		LIMIT(int64(MESSAGES_BATCH)).
		OFFSET(offset)

	var dbMessages []types.DbMessageWithMember

	err = stmt.QueryContext(ctx, q.db, &dbMessages)

	if err != nil {
		return []types.WsMessageContent{}, nil, err
	}

	wsMessages := []types.WsMessageContent{}

	for _, msg := range dbMessages {
		wsMessages = append(wsMessages, types.WsMessageContent{
			WsMessage: transformMsgToWsMsg(msg.Messages),
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

type GetMessageByIDParams struct {
	ID        uuid.UUID
	ChannelID *uuid.UUID
}

func (q *Queries) GetMessageByID(ctx context.Context, params GetMessageByIDParams) (message *model.Messages, err error) {
	var exp BoolExpression = Messages.ID.EQ(UUID(params.ID))

	if params.ChannelID != nil {
		exp = exp.AND(Messages.ChannelID.EQ(UUID(params.ChannelID)))
	}

	stmt := Messages.SELECT(Messages.AllColumns).WHERE(exp)

	var msg model.Messages

	err = stmt.QueryContext(ctx, q.db, &msg)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (q *Queries) UpdateMessageByID(ctx context.Context, id uuid.UUID, content string) (message *types.WsOutgoingMessage, err error) {
	stmt := Messages.UPDATE(Messages.Content).
		MODEL(model.Messages{Content: content}).
		WHERE(Messages.ID.EQ(UUID(id))).
		RETURNING(Messages.AllColumns)

	var updMessage model.Messages
	err = stmt.QueryContext(ctx, q.db, &updMessage)

	if err != nil {
		log.Println("update message content error =>", err)
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
			WsMessage: transformMsgToWsMsg(updMessage),
			Member:    member,
		},
	}

	log.Println("message constructed")

	return message, err
}

func (q *Queries) DeleteMessageByID(ctx context.Context, id uuid.UUID) (message *types.WsOutgoingMessage, err error) {
	stmt := Messages.UPDATE(Messages.FileURL, Messages.Content, Messages.Deleted).
		MODEL(model.Messages{FileURL: nil, Content: "This message has been deleted.", Deleted: true}).
		WHERE(Messages.ID.EQ(UUID(id))).
		RETURNING(Messages.AllColumns)

	var updMessage model.Messages
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
			WsMessage: transformMsgToWsMsg(updMessage),
			Member:    member,
		},
	}

	return message, err
}
