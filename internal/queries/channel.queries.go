package queries

import (
	"context"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
)

func (q *Queries) CreateChannel(ctx context.Context, data model.Channels) (model.Channels, error) {
	stmt := Channels.INSERT(Channels.AllColumns.Except(Channels.CreatedAt, Channels.UpdatedAt)).MODEL(data).RETURNING(Channels.AllColumns)

	var channel model.Channels
	err := stmt.QueryContext(ctx, q.db, &channel)

	return channel, err
}
