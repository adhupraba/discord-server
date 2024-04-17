package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/types"
)

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

	wsMessage := types.WsMessage{
		ID:        message.ID,
		Content:   message.Content,
		FileUrl:   message.FileUrl,
		MemberID:  message.MemberID,
		RoomId:    message.ChannelID,
		Deleted:   message.Deleted,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}

	messageWithMember := types.WsMessageContent{
		WsMessage: wsMessage,
		Member:    member,
	}

	return messageWithMember, err
}
