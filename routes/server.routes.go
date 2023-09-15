package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterServerRoutes() *chi.Mux {
	sc := controllers.ServerController{}
	serverRoute := chi.NewRouter()

	serverRoute.Post("/", middlewares.Auth(sc.CreateServer))

	return serverRoute
}
