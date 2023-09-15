package queries

import (
	"context"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
)

func (q *Queries) CreateMember(ctx context.Context, data model.Members) (model.Members, error) {
	stmt := Members.INSERT(Members.AllColumns.Except(Members.CreatedAt, Members.UpdatedAt)).MODEL(data).RETURNING(Members.AllColumns)

	var member model.Members
	err := stmt.QueryContext(ctx, q.db, &member)

	return member, err
}
