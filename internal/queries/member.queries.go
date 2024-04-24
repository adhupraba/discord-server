package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/types"
)

func (q *Queries) CreateMember(ctx context.Context, data model.Members) (model.Members, error) {
	stmt := Members.INSERT(Members.ID, Members.ProfileID, Members.ServerID, Members.Role).MODEL(data).RETURNING(Members.AllColumns)

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
		)

	var member model.Members
	err := stmt.QueryContext(ctx, q.db, &member)

	return member, err
}

type UpdateMemberRoleParams struct {
	MemberId uuid.UUID
	ServerId uuid.UUID
	Role     model.MemberRole
}

func (q *Queries) UpdateMemberRole(ctx context.Context, params UpdateMemberRoleParams) error {
	stmt := Members.UPDATE(Members.Role).
		MODEL(model.Members{Role: params.Role}).
		WHERE(
			Members.ID.EQ(UUID(params.MemberId)).
				AND(Members.ServerID.EQ(UUID(params.ServerId))),
		)

	_, err := stmt.ExecContext(ctx, q.db)

	return err
}

type RemoveServerMemberParams struct {
	ServerId uuid.UUID
	MemberId uuid.UUID
}

func (q *Queries) RemoveServerMember(ctx context.Context, params RemoveServerMemberParams) error {
	stmt := Members.DELETE().
		WHERE(
			Members.ID.EQ(UUID(params.MemberId)).
				AND(Members.ServerID.EQ(UUID(params.ServerId))),
		)

	_, err := stmt.ExecContext(ctx, q.db)
	return err
}

type GetServerMemberWithProfileParams struct {
	ServerId  uuid.UUID
	ProfileId uuid.UUID
}

func (q *Queries) GetServerMemberWithProfile(ctx context.Context, params GetServerMemberWithProfileParams) (types.MemberWithProfile, error) {
	stmt := SELECT(Members.AllColumns, Profiles.AllColumns).
		FROM(
			Members.LEFT_JOIN(Profiles, Profiles.ID.EQ(Members.ProfileID)),
		).
		WHERE(
			Members.ServerID.EQ(UUID(params.ServerId)).
				AND(Members.ProfileID.EQ(UUID(params.ProfileId))),
		)

	var member types.MemberWithProfile
	err := stmt.QueryContext(ctx, q.db, &member)

	return member, err
}

func (q *Queries) GetMemberWithProfileByMemberID(ctx context.Context, memberId uuid.UUID) (types.MemberWithProfile, error) {
	stmt := SELECT(Members.AllColumns, Profiles.AllColumns).
		FROM(
			Members.LEFT_JOIN(Profiles, Profiles.ID.EQ(Members.ProfileID)),
		).
		WHERE(Members.ID.EQ(UUID(memberId)))

	var member types.MemberWithProfile
	err := stmt.QueryContext(ctx, q.db, &member)

	return member, err
}
