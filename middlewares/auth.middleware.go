package middlewares

import (
	"log"
	"net/http"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/utils"
)

type NextFunc func(http.ResponseWriter, *http.Request, model.Profiles)

func Auth(handler NextFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clerkUser, errCode, err := utils.GetUserFromClerk(r.Context())

		if err != nil {
			log.Println("get user from clerk error =>", r.URL, err)
			utils.RespondWithError(w, errCode, err.Error())
			return
		}

		profile, err := lib.DB.GetUserByClerkID(r.Context(), clerkUser.ID)

		if err != nil {
			log.Println("get user by clerk id error =>", err)
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		handler(w, r, profile)
	}
}
