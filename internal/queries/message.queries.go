package queries

import (
	"context"
	"fmt"
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

	memberStmt := SELECT(Members.AllColumns, Profiles.AllColumns).
		FROM(
			Members.LEFT_JOIN(Profiles, Profiles.ID.EQ(Members.ProfileID)),
		).
		WHERE(Members.ID.EQ(UUID(data.MemberID)))

	var member types.MemberWithProfile
	err = memberStmt.QueryContext(ctx, q.db, &member)

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
			Messages.CreatedAt.LT_EQ(TimestampT(*params.LastMessageDate)).OR(
				Messages.CreatedAt.EQ(TimestampT(*params.LastMessageDate)).AND(Messages.ID.EQ(UUID(*params.LastMessageId))),
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
