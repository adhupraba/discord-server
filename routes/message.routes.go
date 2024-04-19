package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterMessageRoutes() *chi.Mux {
	mc := controllers.MessageController{}
	messageRoute := chi.NewRouter()

	messageRoute.Get("/", middlewares.Auth(mc.GetMessages))
	messageRoute.Post("/send", middlewares.Auth(mc.SendMessage))

	return messageRoute
}
