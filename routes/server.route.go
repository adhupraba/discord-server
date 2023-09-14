package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
)

func RegisterServerRoutes() *chi.Mux {
	sc := controllers.ServerController{}
	serverRoute := chi.NewRouter()

	serverRoute.Get("/", sc.Server)

	return serverRoute
}
