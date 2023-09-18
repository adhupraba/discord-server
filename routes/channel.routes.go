package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterChannelRoutes() *chi.Mux {
	cc := controllers.ChannelController{}
	channelRoute := chi.NewRouter()

	channelRoute.Post("/", middlewares.Auth(cc.CreateChannel))

	return channelRoute
}
