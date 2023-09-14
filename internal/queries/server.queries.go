package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/lib"
)

func GetFirstServerOfUser(ctx context.Context, profileId uuid.UUID) (*model.Servers, error) {
	stmt := SELECT(Servers.AllColumns).
		FROM(Servers).
		WHERE(
			Servers.ProfileID.EQ(UUID(profileId)),
		).LIMIT(1)

	var server model.Servers
	err := stmt.QueryContext(ctx, lib.DB, &server)

	if err != nil && err == qrm.ErrNoRows {
		return nil, nil
	}

	return &server, err
}
