package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/types"
)

func (q *Queries) CreateChannelMessage(ctx context.Context, data model.Messages) (types.MessageWithMember, error) {
	stmt := Messages.INSERT(
		Messages.AllColumns.Except(Messages.CreatedAt, Messages.UpdatedAt),
	).
		MODEL(data).
		RETURNING(Messages.AllColumns)

	var message model.Messages
	err := stmt.QueryContext(ctx, q.db, &message)

	if err != nil {
		return types.MessageWithMember{}, err
	}

	memberStmt := SELECT(Members.AllColumns, Profiles.AllColumns).
		FROM(
			Members.LEFT_JOIN(Profiles, Profiles.ID.EQ(Members.ProfileID)),
		).
		WHERE(Members.ID.EQ(UUID(data.MemberID)))

	var member types.MemberWithProfile
	err = memberStmt.QueryContext(ctx, q.db, &member)

	messageWithMember := types.MessageWithMember{
		Messages: message,
		Member:   member,
	}

	return messageWithMember, err
}
