package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterDirectMessageRoutes() *chi.Mux {
	dc := controllers.DirectMessageController{}
	directMsgRoute := chi.NewRouter()

	directMsgRoute.Get("/", middlewares.Auth(dc.GetMessages))

	return directMsgRoute
}
