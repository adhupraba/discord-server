package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/types"
)

func (q *Queries) GetUserByClerkID(ctx context.Context, clerkID string) (model.Profiles, error) {
	stmt := SELECT(Profiles.AllColumns).
		FROM(Profiles).
		WHERE(Profiles.UserID.EQ(String(clerkID)))

	var profile model.Profiles
	err := stmt.QueryContext(ctx, q.db, &profile)

	return profile, err
}

func (q *Queries) InsertUserProfile(ctx context.Context, data model.Profiles) (model.Profiles, error) {
	stmt := Profiles.INSERT(Profiles.AllColumns.Except(Profiles.CreatedAt, Profiles.UpdatedAt)).MODEL(data).RETURNING(Profiles.AllColumns)

	var profile model.Profiles
	err := stmt.QueryContext(ctx, q.db, &profile)

	return profile, err
}

func (q *Queries) GetUserAndServers(ctx context.Context, profileID uuid.UUID) (types.ProfileAndServer, error) {
	stmt := SELECT(Profiles.AllColumns).
		FROM(
			Profiles.LEFT_JOIN(
				Servers, Servers.ProfileID.EQ(Profiles.ID),
			),
		).
		WHERE(Profiles.ID.EQ(UUID(profileID)))

	var profile types.ProfileAndServer
	err := stmt.QueryContext(ctx, q.db, &profile)

	return profile, err
}
