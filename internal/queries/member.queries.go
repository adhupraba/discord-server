package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
)

func (q *Queries) CreateMember(ctx context.Context, data model.Members) (model.Members, error) {
	stmt := Members.INSERT(Members.ID, Members.ProfileID, Members.ServerID).MODEL(data).RETURNING(Members.AllColumns)

	var member model.Members
	err := stmt.QueryContext(ctx, q.db, &member)

	return member, err
}

type GetServerMemberParams struct {
	ServerId  uuid.UUID
	ProfileId uuid.UUID
}

func (q *Queries) GetServerMember(ctx context.Context, params GetServerMemberParams) (model.Members, error) {
	stmt := SELECT(Members.AllColumns).
		FROM(Members).
		WHERE(
			Members.ServerID.EQ(UUID(params.ServerId)).
				AND(Members.ProfileID.EQ(UUID(params.ProfileId))),
		).LIMIT(1)

	var member model.Members
	err := stmt.QueryContext(ctx, q.db, &member)

	return member, err
}
