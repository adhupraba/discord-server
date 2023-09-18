package controllers

import (
	"net/http"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/utils"
)

type ChannelController struct{}

func (cc *ChannelController) CreateChannel(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	idQ := r.URL.Query().Get("serverId")
	serverId, err := uuid.Parse(idQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	type createChannelBody struct {
		Name string            `json:"name" validate:"required,min=1,max=128"`
		Type model.ChannelType `json:"type" validate:"required,oneof=AUDIO VIDEO TEXT"`
	}

	var body createChannelBody
	err = utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid data")
		return
	}

	if body.Name == "general" {
		utils.RespondWithError(w, http.StatusBadRequest, "Channel name cannot be 'general'")
		return
	}

	profileMember, err := lib.DB.GetServerMember(r.Context(), queries.GetServerMemberParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusUnauthorized, "You are not part of the server to create a channel")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when validating user")
		return
	}

	if profileMember.Role != model.MemberRoleADMIN && profileMember.Role != model.MemberRoleMODERATOR {
		utils.RespondWithError(w, http.StatusForbidden, "Only admins and moderators can create a channel")
		return
	}

	channel, err := lib.DB.CreateChannel(r.Context(), model.Channels{
		ID:        uuid.New(),
		Name:      body.Name,
		Type:      body.Type,
		ProfileID: profile.ID,
		ServerID:  serverId,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when creating channel")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, channel)
}
