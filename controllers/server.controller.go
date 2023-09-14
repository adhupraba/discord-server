package controllers

import (
	"log"
	"net/http"

	"github.com/go-jet/jet/v2/qrm"

	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/utils"
)

type ServerController struct{}

func (hc *ServerController) Server(w http.ResponseWriter, r *http.Request) {
	clerkUser, errCode, err := utils.GetUserFromClerk(r.Context())

	if err != nil {
		utils.RespondWithError(w, errCode, err.Error())
		return
	}

	dbUser, err := queries.GetUserByClerkID(r.Context(), clerkUser.ID)

	if err != nil && err != qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if dbUser.UserID != "" {
		log.Println("sending existing user")
		utils.RespondWithJson(w, http.StatusOK, dbUser)
		return
	}

	utils.RespondWithJson(w, http.StatusOK, dbUser)
}
