package controllers

import (
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/utils"
)

type ServerController struct{}

func (hc *ServerController) CreateServer(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
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
		InviteCode: uuid.New().String(),
		ProfileID:  profile.ID,
	}

	server, err := lib.DB.CreateServerWithTx(r.Context(), lib.SqlConn, serverData)

	if err != nil {
		log.Println("server creation error:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating your server")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, server)
}
