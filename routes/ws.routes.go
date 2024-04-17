package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterWsRoutes() *chi.Mux {
	wc := controllers.WsController{}
	wcRoute := chi.NewRouter()

	wcRoute.Get("/connect", middlewares.Auth(wc.Connect))
	wcRoute.Post("/messages", middlewares.Auth(wc.SendMessage))

	return wcRoute
}
