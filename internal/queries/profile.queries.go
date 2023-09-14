package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/lib"
)

func GetUserByClerkID(ctx context.Context, clerkID string) (model.Profiles, error) {
	stmt := SELECT(Profiles.AllColumns).
		FROM(Profiles).
		WHERE(
			Profiles.UserID.EQ(String(clerkID)),
		)

	var profile model.Profiles
	err := stmt.QueryContext(ctx, lib.DB, &profile)

	return profile, err
}

func InsertUserProfile(ctx context.Context, data model.Profiles) (model.Profiles, error) {
	stmt := Profiles.INSERT(Profiles.AllColumns.Except(Profiles.CreatedAt, Profiles.UpdatedAt)).
		MODEL(data).
		RETURNING(Profiles.AllColumns)

	var profile model.Profiles
	err := stmt.QueryContext(ctx, lib.DB, &profile)

	return profile, err
}
