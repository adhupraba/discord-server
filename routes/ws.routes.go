package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
)

func RegisterWsRoutes() *chi.Mux {
	wc := controllers.WsController{}
	wcRoute := chi.NewRouter()

	wcRoute.Get("/connect", wc.Connect)

	return wcRoute
}
