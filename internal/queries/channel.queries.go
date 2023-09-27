package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
)

func (q *Queries) CreateChannel(ctx context.Context, data model.Channels) (model.Channels, error) {
	stmt := Channels.INSERT(
		Channels.AllColumns.Except(Channels.CreatedAt, Channels.UpdatedAt),
	).
		MODEL(data).
		RETURNING(Channels.AllColumns)

	var channel model.Channels
	err := stmt.QueryContext(ctx, q.db, &channel)

	return channel, err
}

type GetServerChannelParams struct {
	ChannelId uuid.UUID
	ServerId  *uuid.UUID
}

func (q *Queries) GetServerChannel(ctx context.Context, params GetServerChannelParams) (model.Channels, error) {
	var exp BoolExpression

	if params.ServerId == nil {
		exp = Bool(true)
	} else {
		exp = Channels.ServerID.EQ(UUID(params.ServerId))
	}

	stmt := SELECT(Channels.AllColumns).
		FROM(Channels).
		WHERE(
			Channels.ID.EQ(UUID(params.ChannelId)).
				AND(exp),
		)

	var channel model.Channels
	err := stmt.QueryContext(ctx, q.db, &channel)

	return channel, err
}

func (q *Queries) DeleteChannel(ctx context.Context, channelId uuid.UUID) error {
	stmt := Channels.DELETE().
		WHERE(Channels.ID.EQ(UUID(channelId)))

	_, err := stmt.ExecContext(ctx, q.db)

	return err
}

type UpdateChannelParams struct {
	ChannelId uuid.UUID
	Data      model.Channels
}

func (q *Queries) UpdateChannel(ctx context.Context, params UpdateChannelParams) (model.Channels, error) {
	stmt := Channels.UPDATE(Channels.Name, Channels.Type).
		MODEL(params.Data).
		WHERE(Channels.ID.EQ(UUID(params.ChannelId))).
		RETURNING(Channels.AllColumns)

	var channel model.Channels
	err := stmt.QueryContext(ctx, q.db, &channel)

	return channel, err
}

func (q *Queries) GetServerGeneralChannel(ctx context.Context, serverId uuid.UUID) (model.Channels, error) {
	stmt := SELECT(Channels.AllColumns).
		FROM(Channels).
		WHERE(
			Channels.ServerID.EQ(UUID(serverId)).
				AND(Channels.Name.EQ(String("general"))),
		)

	var channel model.Channels
	err := stmt.QueryContext(ctx, q.db, &channel)

	return channel, err
}
