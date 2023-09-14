package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
)

func RegisterProfileRoutes() *chi.Mux {
	pc := controllers.ProfileController{}
	profileRoute := chi.NewRouter()

	profileRoute.Get("/", pc.Profile)

	return profileRoute
}
