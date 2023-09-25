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

func (cc *ChannelController) DeleteChannel(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	cIdQ := chi.URLParam(r, "channelId")
	channelId, err := uuid.Parse(cIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid channel id")
		return
	}

	sIdQ := r.URL.Query().Get("serverId")
	serverId, err := uuid.Parse(sIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	channel, err := lib.DB.GetServerChannel(r.Context(), queries.GetServerChannelParams{
		ChannelId: channelId,
		ServerId:  serverId,
	})

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusNotFound, "Channel not found")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when validating channel data")
		return
	}

	if channel.Name == "general" {
		utils.RespondWithError(w, http.StatusBadRequest, "general channel cannot be deleted")
		return
	}

	profileMember, err := lib.DB.GetServerMember(r.Context(), queries.GetServerMemberParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusUnauthorized, "You are not part of the server to delete a channel")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when validating user")
		return
	}

	if profileMember.Role != model.MemberRoleADMIN && profileMember.Role != model.MemberRoleMODERATOR {
		utils.RespondWithError(w, http.StatusForbidden, "You do not have permission to delete the channel")
		return
	}

	err = lib.DB.DeleteChannel(r.Context(), channel.ID)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when deleting channel")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, types.Json{"message": "Deleted channel successfully"})
}

func (cc *ChannelController) UpdateChannel(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	cIdQ := chi.URLParam(r, "channelId")
	channelId, err := uuid.Parse(cIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid channel id")
		return
	}

	sIdQ := r.URL.Query().Get("serverId")
	serverId, err := uuid.Parse(sIdQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	type updateChannelBody struct {
		Name string            `json:"name" validate:"required,min=1,max=128"`
		Type model.ChannelType `json:"type" validate:"required,oneof=AUDIO VIDEO TEXT"`
	}

	var body updateChannelBody
	err = utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid data")
		return
	}

	if body.Name == "general" {
		utils.RespondWithError(w, http.StatusBadRequest, "Channel name cannot be 'general'")
		return
	}

	channel, err := lib.DB.GetServerChannel(r.Context(), queries.GetServerChannelParams{
		ChannelId: channelId,
		ServerId:  serverId,
	})

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusNotFound, "Channel not found")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when validating channel data")
		return
	}

	profileMember, err := lib.DB.GetServerMember(r.Context(), queries.GetServerMemberParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusUnauthorized, "You are not part of the server to update a channel")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when validating user")
		return
	}

	if profileMember.Role != model.MemberRoleADMIN && profileMember.Role != model.MemberRoleMODERATOR {
		utils.RespondWithError(w, http.StatusForbidden, "You do not have permission to update the channel")
		return
	}

	channel, err = lib.DB.UpdateChannel(r.Context(), queries.UpdateChannelParams{
		ChannelId: channelId,
		Data: model.Channels{
			Name: body.Name,
			Type: body.Type,
		},
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when creating channel")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, channel)
}
