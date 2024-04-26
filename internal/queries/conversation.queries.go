package queries

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	. "github.com/adhupraba/discord-server/internal/discord/public/table"
	"github.com/adhupraba/discord-server/types"
)

func (q *Queries) CreateConversation(ctx context.Context, data model.Conversations) (model.Conversations, error) {
	stmt := Conversations.INSERT(
		Conversations.AllColumns.Except(Conversations.CreatedAt, Conversations.UpdatedAt),
	).
		MODEL(data).
		RETURNING(Conversations.AllColumns)

	var conversation model.Conversations
	err := stmt.QueryContext(ctx, q.db, &conversation)

	return conversation, err
}

type GetConversationByMembersParams struct {
	MemberOneId uuid.UUID
	MemberTwoId uuid.UUID
}

func (q *Queries) GetConversationByMembers(ctx context.Context, params GetConversationByMembersParams) (types.ConversationWithMemberAndProfile, error) {
	memberOne := Members.AS("member_one")
	profileOne := Profiles.AS("profile_one")
	memberTwo := Members.AS("member_two")
	profileTwo := Profiles.AS("profile_two")

	stmt := SELECT(
		Conversations.AllColumns,
		memberOne.AllColumns,
		profileOne.AllColumns,
		memberTwo.AllColumns,
		profileTwo.AllColumns,
	).
		FROM(
			Conversations.
				LEFT_JOIN(memberOne, memberOne.ID.EQ(Conversations.MemberOneID)).
				LEFT_JOIN(profileOne, profileOne.ID.EQ(memberOne.ProfileID)).
				LEFT_JOIN(memberTwo, memberTwo.ID.EQ(Conversations.MemberTwoID)).
				LEFT_JOIN(profileTwo, profileTwo.ID.EQ(memberTwo.ProfileID)),
		).
		WHERE(
			Conversations.MemberOneID.EQ(UUID(params.MemberOneId)).
				AND(Conversations.MemberTwoID.EQ(UUID(params.MemberTwoId))),
		)

	var conversation types.ConversationWithMemberAndProfile
	err := stmt.QueryContext(ctx, q.db, &conversation)

	return conversation, err
}

type GetConversationWithMembersByIDParams struct {
	ConversationID uuid.UUID
	ProfileID      *uuid.UUID
}

func (q *Queries) GetConversationWithMembersByID(ctx context.Context, params GetConversationWithMembersByIDParams) (types.ConversationWithMemberAndProfile, error) {
	memberOne := Members.AS("member_one")
	profileOne := Profiles.AS("profile_one")
	memberTwo := Members.AS("member_two")
	profileTwo := Profiles.AS("profile_two")

	exp := Conversations.ID.EQ(UUID(params.ConversationID))

	if params.ProfileID != nil {
		exp.AND(
			OR(
				memberOne.ProfileID.EQ(UUID(params.ProfileID)),
				memberTwo.ProfileID.EQ(UUID(params.ProfileID)),
			),
		)
	}

	stmt := SELECT(
		Conversations.AllColumns,
		memberOne.AllColumns,
		profileOne.AllColumns,
		memberTwo.AllColumns,
		profileTwo.AllColumns,
	).
		FROM(
			Conversations.
				LEFT_JOIN(memberOne, memberOne.ID.EQ(Conversations.MemberOneID)).
				LEFT_JOIN(profileOne, profileOne.ID.EQ(memberOne.ProfileID)).
				LEFT_JOIN(memberTwo, memberTwo.ID.EQ(Conversations.MemberTwoID)).
				LEFT_JOIN(profileTwo, profileTwo.ID.EQ(memberTwo.ProfileID)),
		).
		WHERE(exp)

	var conversation types.ConversationWithMemberAndProfile
	err := stmt.QueryContext(ctx, q.db, &conversation)

	return conversation, err
}
