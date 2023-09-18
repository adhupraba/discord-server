package controllers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/constants"
	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
	"github.com/adhupraba/discord-server/utils"
)

type ServerController struct{}

func (sc *ServerController) CreateServer(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	type createServerBody struct {
		Name     string `json:"name" validate:"required,min=1,max=128"`
		ImageURL string `json:"imageUrl" validate:"required,url"`
	}

	var body createServerBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid data sent")
		return
	}

	serverData := model.Servers{
		ID:         uuid.New(),
		Name:       body.Name,
		ImageURL:   body.ImageURL,
		InviteCode: uuid.New(),
		ProfileID:  profile.ID,
	}

	server, err := lib.DB.CreateServerWithTx(r.Context(), queries.CreateServerWithTxParams{
		Db:   lib.SqlConn,
		Data: serverData,
	})

	if err != nil {
		log.Println("server creation error:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating your server")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, server)
}

func (sc *ServerController) GetUserMemberServers(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	servers, err := lib.DB.GetServersOfUser(r.Context(), queries.GetServersOfUserParams{
		ProfileId: profile.ID,
		Opts:      nil,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding servers")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, servers)
}

func (sc *ServerController) GetServer(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	idQ := chi.URLParam(r, "serverId")
	serverId, err := uuid.Parse(idQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	server, err := lib.DB.GetServer(r.Context(), queries.GetServerParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusNotFound, "Server does not exist")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding the server")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, server)
}

func (sc *ServerController) GetFullServerDetails(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	idQ := chi.URLParam(r, "serverId")

	serverId, err := uuid.Parse(idQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	server, err := lib.DB.GetServerWithChannelsAndMembers(r.Context(), serverId)

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusNotFound, "Server does not exist")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding the server")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, server)
}

func (sc *ServerController) UpdateInviteCode(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	idQ := chi.URLParam(r, "serverId")
	serverId, err := uuid.Parse(idQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	server, err := lib.DB.UpdateServerInviteCode(r.Context(), queries.UpdateServerInviteCodeParams{
		ServerId:   serverId,
		ProfileId:  profile.ID,
		InviteCode: uuid.New(),
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error while updating new invite code")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, server)
}

func (sc *ServerController) VerifyAndAcceptInviteCode(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	idQ := chi.URLParam(r, "inviteCode")
	inviteCode, err := uuid.Parse(idQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid invite code")
		return
	}

	existing, err := lib.DB.FindUserInServerWithInviteCode(r.Context(), queries.FindUserInServerWithInviteCodeParams{
		InviteCode: inviteCode,
		ProfileId:  profile.ID,
	})

	if err != nil && err != qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error validating invite code")
		return
	}

	if existing.ID.String() != constants.EmptyUUID {
		utils.RespondWithJson(w, http.StatusFound, types.Json{
			"existing": true,
			"server":   existing,
		})
		return
	}

	server, err := lib.DB.GetServerUsingInviteCode(r.Context(), inviteCode)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when adding as member into server")
		return
	}

	member, err := lib.DB.CreateMember(r.Context(), model.Members{
		ID:        uuid.New(),
		ProfileID: profile.ID,
		ServerID:  server.ID,
		Role:      model.MemberRoleGUEST,
	})

	if err != nil {
		log.Println("error adding as member =>", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error adding as member")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, types.Json{
		"existing": false,
		"member":   member,
	})
}

func (sc *ServerController) UpdateServer(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	idQ := chi.URLParam(r, "serverId")
	serverId, err := uuid.Parse(idQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	type updateServerBody struct {
		Name     string `json:"name" validate:"required,min=1,max=128"`
		ImageURL string `json:"imageUrl" validate:"required,url"`
	}

	var body updateServerBody
	err = utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid data")
		return
	}

	params := queries.UpdateServerParams{ServerId: serverId, ProfileId: profile.ID}
	data := model.Servers{Name: body.Name, ImageURL: body.ImageURL}

	server, err := lib.DB.UpdateServer(r.Context(), params, data)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when updating server data")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, server)
}

func (sc *ServerController) MemberLeaveServer(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	idQ := chi.URLParam(r, "serverId")
	serverId, err := uuid.Parse(idQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	member, err := lib.DB.GetServerMember(r.Context(), queries.GetServerMemberParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error valdiating user")
		return
	}

	if member.Role == model.MemberRoleADMIN {
		utils.RespondWithError(w, http.StatusForbidden, "An admin cannot leave the server")
		return
	}

	err = lib.DB.RemoveServerMember(r.Context(), queries.RemoveServerMemberParams{
		ServerId: serverId,
		MemberId: member.ID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when removing member")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Member has left the server."})
}

func (sc *ServerController) DeleteServer(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	idQ := chi.URLParam(r, "serverId")
	serverId, err := uuid.Parse(idQ)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid server id")
		return
	}

	err = lib.DB.DeleteServer(r.Context(), queries.DeleteServerParams{
		ServerId:  serverId,
		ProfileId: profile.ID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when deleting server")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "Deleted the server."})
}
