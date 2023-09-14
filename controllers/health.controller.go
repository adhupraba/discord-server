package controllers

import (
	"net/http"

	"github.com/adhupraba/discord-server/types"
	"github.com/adhupraba/discord-server/utils"
)

type HealthController struct{}

func (hc *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "success"})
}
