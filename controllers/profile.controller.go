package controllers

import (
	"net/http"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/constants"
	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
	"github.com/adhupraba/discord-server/utils"
)

type ProfileController struct{}

func (pc *ProfileController) UpsertProfile(w http.ResponseWriter, r *http.Request) {
	clerkUser, errCode, err := utils.GetUserFromClerk(r.Context())

	if err != nil {
		utils.RespondWithError(w, errCode, err.Error())
		return
	}

	dbUser, err := lib.DB.GetUserByClerkID(r.Context(), clerkUser.ID)

	if err != nil && err != qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if dbUser.ID.String() != constants.EmptyUUID {
		servers, err := lib.DB.GetServersOfUser(r.Context(), queries.GetServersOfUserParams{
			ProfileId: dbUser.ID,
			Opts:      &types.PaginationOpts{Limit: 1},
		})

		if err != nil && err != qrm.ErrNoRows {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching server")
			return
		}

		var server *model.Servers

		if len(servers) > 0 {
			server = &servers[0]
		}

		utils.RespondWithJson(w, http.StatusOK, types.ProfileAndServer{
			Profiles: dbUser,
			Server:   server,
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

	dbUser, err = lib.DB.InsertUserProfile(r.Context(), data)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, types.ProfileAndServer{
		Profiles: dbUser,
		Server:   nil,
	})
}

func (pc *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	utils.RespondWithJson(w, http.StatusOK, profile)
}
