package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterProfileRoutes() *chi.Mux {
	pc := controllers.ProfileController{}
	profileRoute := chi.NewRouter()

	profileRoute.Get("/upsert", pc.UpsertProfile)
	profileRoute.Get("/", middlewares.Auth(pc.GetProfile))

	return profileRoute
}
