package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterConversationRoutes() *chi.Mux {
	cc := controllers.ConversationController{}
	conversationRoute := chi.NewRouter()

	conversationRoute.Get("/get-by-members", middlewares.Auth(cc.GetConversationByMembers))
	conversationRoute.Post("/", middlewares.Auth(cc.CreateNewConversation))

	return conversationRoute
}
