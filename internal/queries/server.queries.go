package queries

import (
	"context"
	"database/sql"
	"errors"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/types"
)

type CreateServerData struct {
	model.Servers
	Channel model.Channels  `json:"channel"`
	Members []model.Members `json:"members"`
}

type GetServersOfUserParams struct {
	ProfileId uuid.UUID
	Opts      *types.PaginationOpts
}

func (q *Queries) GetServersOfUser(ctx context.Context, params GetServersOfUserParams) ([]model.Servers, error) {
	stmt := SELECT(Servers.AllColumns).
		FROM(
			Servers.
				LEFT_JOIN(Members, Members.ServerID.EQ(Servers.ID)),
		).
		WHERE(Members.ProfileID.EQ(UUID(params.ProfileId)))

	if params.Opts != nil && params.Opts.Limit > 0 {
		stmt = stmt.LIMIT(params.Opts.Limit)
	}

	if params.Opts != nil && params.Opts.Offset >= 0 {
		stmt = stmt.OFFSET(params.Opts.Offset)
	}

	servers := []model.Servers{}
	err := stmt.QueryContext(ctx, q.db, &servers)

	if err != nil && err == qrm.ErrNoRows {
		return []model.Servers{}, nil
	}

	return servers, err
}

type CreateServerWithTxParams struct {
	Db   *sql.DB
	Data model.Servers
}

func (q *Queries) CreateServerWithTx(ctx context.Context, params CreateServerWithTxParams) (*CreateServerData, error) {
	tx, err := params.Db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	qtx := q.WithTx(tx)

	stmt := Servers.INSERT(Servers.AllColumns.Except(Servers.CreatedAt, Servers.UpdatedAt)).
		MODEL(params.Data).
		RETURNING(Servers.AllColumns)

	var server model.Servers
	err = stmt.QueryContext(ctx, qtx.db, &server)

	if err != nil {
		return nil, errors.New("Error creating your server")
	}

	channelData := model.Channels{
		ID:        uuid.New(),
		Name:      "general",
		ProfileID: server.ProfileID,
		ServerID:  server.ID,
	}
	channel, err := qtx.CreateChannel(ctx, channelData)

	if err != nil {
		return nil, errors.New("Error creating default channel")
	}

	memberData := model.Members{
		ID:        uuid.New(),
		ProfileID: server.ProfileID,
		ServerID:  server.ID,
		Role:      model.MemberRoleADMIN,
	}
	member, err := qtx.CreateMember(ctx, memberData)

	if err != nil {
		return nil, errors.New("Error making you member in the server")
	}

	err = tx.Commit()

	if err != nil {
		return nil, errors.New("Error commiting the transaction")
	}

	res := &CreateServerData{
		server,
		channel,
		[]model.Members{member},
	}

	return res, nil
}

type GetServerParams struct {
	ServerId  uuid.UUID
	ProfileId uuid.UUID
}

func (q *Queries) GetServer(ctx context.Context, params GetServerParams) (model.Servers, error) {
	stmt := SELECT(Servers.AllColumns).
		FROM(
			Servers.
				LEFT_JOIN(Members, Members.ServerID.EQ(Servers.ID)),
		).
		WHERE(
			Servers.ID.EQ(UUID(params.ServerId)).
				AND(Members.ProfileID.EQ(UUID(params.ProfileId))),
		)

	var server model.Servers
	err := stmt.QueryContext(ctx, q.db, &server)

	return server, err
}

func (q *Queries) GetServerWithChannelsAndMembers(ctx context.Context, serverId uuid.UUID) (types.ServerWithChannelsAndMembers, error) {
	stmt := SELECT(
		Servers.AllColumns,
		Channels.AllColumns,
		Members.AllColumns,
		Profiles.AllColumns,
	).
		FROM(
			Servers.
				LEFT_JOIN(Channels, Channels.ServerID.EQ(Servers.ID)).
				LEFT_JOIN(Members, Members.ServerID.EQ(Servers.ID)).
				LEFT_JOIN(Profiles, Profiles.ID.EQ(Members.ProfileID)),
		).
		WHERE(Servers.ID.EQ(UUID(serverId))).
		ORDER_BY(Channels.Name, Members.Role)

	var server types.ServerWithChannelsAndMembers
	err := stmt.QueryContext(ctx, q.db, &server)

	return server, err
}

type UpdateServerInviteCodeParams struct {
	ServerId   uuid.UUID
	ProfileId  uuid.UUID
	InviteCode uuid.UUID
}

func (q *Queries) UpdateServerInviteCode(ctx context.Context, params UpdateServerInviteCodeParams) (model.Servers, error) {
	stmt := Servers.UPDATE(Servers.InviteCode).
		SET(Servers.InviteCode.SET(UUID(params.InviteCode))).
		WHERE(
			Servers.ID.EQ(UUID(params.ServerId)).
				AND(Servers.ProfileID.EQ(UUID(params.ProfileId))),
		).
		RETURNING(Servers.AllColumns)

	var server model.Servers
	err := stmt.QueryContext(ctx, q.db, &server)

	return server, err
}

type FindUserInServerWithInviteCodeParams struct {
	InviteCode uuid.UUID
	ProfileId  uuid.UUID
}

func (q *Queries) FindUserInServerWithInviteCode(ctx context.Context, params FindUserInServerWithInviteCodeParams) (model.Servers, error) {
	stmt := SELECT(Servers.AllColumns).
		FROM(
			Servers.
				LEFT_JOIN(Members, Members.ServerID.EQ(Servers.ID)),
		).
		WHERE(
			Servers.InviteCode.EQ(UUID(params.InviteCode)).
				AND(Members.ProfileID.EQ(UUID(params.ProfileId))),
		)

	var server model.Servers
	err := stmt.QueryContext(ctx, q.db, &server)

	return server, err
}

func (q *Queries) GetServerUsingInviteCode(ctx context.Context, inviteCode uuid.UUID) (model.Servers, error) {
	stmt := SELECT(Servers.ID).
		FROM(Servers).
		WHERE(Servers.InviteCode.EQ(UUID(inviteCode)))

	var server model.Servers
	err := stmt.QueryContext(ctx, q.db, &server)

	return server, err
}

type UpdateServerParams struct {
	ServerId  uuid.UUID
	ProfileId uuid.UUID
}

func (q *Queries) UpdateServer(ctx context.Context, params UpdateServerParams, data model.Servers) (model.Servers, error) {
	stmt := Servers.UPDATE(Servers.Name, Servers.ImageURL).
		MODEL(data).
		WHERE(
			Servers.ID.EQ(UUID(params.ServerId)).
				AND(Servers.ProfileID.EQ(UUID(params.ProfileId))),
		).
		RETURNING(Servers.AllColumns)

	var server model.Servers
	err := stmt.QueryContext(ctx, q.db, &server)

	return server, err
}

type DeleteServerParams struct {
	ServerId  uuid.UUID
	ProfileId uuid.UUID
}

func (q *Queries) DeleteServer(ctx context.Context, params DeleteServerParams) error {
	stmt := Servers.DELETE().
		WHERE(
			Servers.ID.EQ(UUID(params.ServerId)).
				AND(Servers.ProfileID.EQ(UUID(params.ProfileId))),
		)

	_, err := stmt.ExecContext(ctx, q.db)

	return err
}
