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
)

type CreateServerData struct {
	model.Servers
	Channel model.Channels  `json:"channel"`
	Members []model.Members `json:"members"`
}

func (q *Queries) GetServersOfUser(ctx context.Context, profileId uuid.UUID) ([]model.Servers, error) {
	stmt := SELECT(Servers.AllColumns).
		FROM(
			Servers.
				LEFT_JOIN(Members, Members.ServerID.EQ(Servers.ID)),
		).
		WHERE(
			Members.ProfileID.EQ(UUID(profileId)),
		)

	servers := []model.Servers{}
	err := stmt.QueryContext(ctx, q.db, &servers)

	if err != nil && err == qrm.ErrNoRows {
		return []model.Servers{}, nil
	}

	return servers, err
}

func (q *Queries) CreateServerWithTx(ctx context.Context, db *sql.DB, data model.Servers) (*CreateServerData, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	qtx := q.WithTx(tx)

	stmt := Servers.INSERT(Servers.AllColumns.Except(Servers.CreatedAt, Servers.UpdatedAt)).MODEL(data).RETURNING(Servers.AllColumns)

	var server model.Servers
	err = stmt.QueryContext(ctx, qtx.db, &server)

	if err != nil {
		return nil, errors.New("Error creating your server")
	}

	channelData := model.Channels{
		ID:        uuid.New(),
		Name:      "general",
		Type:      model.ChannelTypeTEXT,
		ProfileID: server.ProfileID,
		ServerID:  server.ID,
	}
	channel, err := qtx.CreateChannel(ctx, channelData)

	if err != nil {
		return nil, errors.New("Error creating default channel")
	}

	memberData := model.Members{
		ID:        uuid.New(),
		Role:      model.MemberRoleADMIN,
		ProfileID: server.ProfileID,
		ServerID:  server.ID,
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
