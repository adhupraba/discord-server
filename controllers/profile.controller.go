package controllers

import (
	"net/http"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/utils"
)

type ProfileController struct{}

type profileResponse struct {
	model.Profiles
	Server *model.Servers `json:"server"`
}

func (hc *ProfileController) Profile(w http.ResponseWriter, r *http.Request) {
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
		server, err := queries.GetFirstServerOfUser(r.Context(), dbUser.ID)

		if err != nil && err != qrm.ErrNoRows {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching server")
			return
		}

		utils.RespondWithJson(w, http.StatusOK, profileResponse{
			dbUser,
			server,
		})
		return
	}

	data := model.Profiles{
		ID:       uuid.New(),
		UserID:   clerkUser.ID,
		Name:     *clerkUser.FirstName + " " + *clerkUser.LastName,
		ImageURL: *clerkUser.ImageURL,
		Email:    *&clerkUser.EmailAddresses[0].EmailAddress,
	}

	dbUser, err = queries.InsertUserProfile(r.Context(), data)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, profileResponse{
		dbUser,
		nil,
	})
}
