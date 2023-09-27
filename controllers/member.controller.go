package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
	"github.com/adhupraba/discord-server/utils"
)

type MemberController struct{}

type ServerWithMembersResponse struct {
	model.Servers
	Members []types.MemberWithProfile `json:"members"`
}

func (mc *MemberController) UpdateMemberRole(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	mIdQ := chi.URLParam(r, "memberId")
	memberId, err := uuid.Parse(mIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid member id")
		return
	}

	sIdQ := r.URL.Query().Get("serverId")

	if sIdQ == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Server id is not present")
		return
	}

	serverId, err := uuid.Parse(sIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	type updateMemberBody struct {
		Role string `json:"role" validate:"required,oneof=MODERATOR GUEST"`
	}

	var body updateMemberBody
	err = utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid data")
		return
	}

	profileMember, err := lib.DB.GetServerMember(r.Context(), queries.GetServerMemberParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	// the user is not in the server, then invalid request
	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusUnauthorized, "Only an admin of a server can update a role of a member")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when validating user")
		return
	}

	// the user is the admin of the server, then do not allow to change the role
	if profileMember.ProfileID == memberId && profileMember.Role == model.MemberRoleADMIN {
		utils.RespondWithError(w, http.StatusForbidden, "Admins cannot change their own role")
		return
	}

	// the user is not an admin of the server, do not allow to update role
	if profileMember.Role != model.MemberRoleADMIN {
		utils.RespondWithError(w, http.StatusForbidden, "Only an admin can change a member's role")
		return
	}

	err = lib.DB.UpdateMemberRole(r.Context(), queries.UpdateMemberRoleParams{
		MemberId: memberId,
		ServerId: serverId,
		Role:     model.MemberRole(body.Role),
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when updating server data")
		return
	}

	server, err := lib.DB.GetServerWithChannelsAndMembers(r.Context(), serverId)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching server data")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, ServerWithMembersResponse{
		Servers: server.Servers,
		Members: server.Members,
	})
}

func (mc *MemberController) RemoveServerMember(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	mIdQ := chi.URLParam(r, "memberId")
	memberId, err := uuid.Parse(mIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid member id")
		return
	}

	sIdQ := r.URL.Query().Get("serverId")

	if sIdQ == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Server id is not present")
		return
	}

	serverId, err := uuid.Parse(sIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	profileMember, err := lib.DB.GetServerMember(r.Context(), queries.GetServerMemberParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	// the user is not in the server, then invalid request
	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusUnauthorized, "Only an admin of a server can remove a member")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when validating user")
		return
	}

	// the user is the admin of the server, then do not allow to remove themselves from server
	if profileMember.ProfileID == memberId && profileMember.Role == model.MemberRoleADMIN {
		utils.RespondWithError(w, http.StatusForbidden, "Admins cannot remove themselves from the server")
		return
	}

	// the user is not an admin of the server, do not allow to remove a member
	if profileMember.Role != model.MemberRoleADMIN {
		utils.RespondWithError(w, http.StatusForbidden, "Only an admin can remove a member of a server")
		return
	}

	err = lib.DB.RemoveServerMember(r.Context(), queries.RemoveServerMemberParams{
		ServerId: serverId,
		MemberId: memberId,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when removing member")
		return
	}

	server, err := lib.DB.GetServerWithChannelsAndMembers(r.Context(), serverId)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching server data")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, ServerWithMembersResponse{
		Servers: server.Servers,
		Members: server.Members,
	})
}

func (mc *MemberController) GetServerMember(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	sIdQ := chi.URLParam(r, "serverId")
	serverId, err := uuid.Parse(sIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	member, err := lib.DB.GetServerMemberWithProfile(r.Context(), queries.GetServerMemberWithProfileParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusNotFound, "Member not found")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching member")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, member)
}
