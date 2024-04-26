package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterWsRoutes() *chi.Mux {
	wc := controllers.WsController{}
	wcRoute := chi.NewRouter()
	wcRoute.Get("/connect", wc.Connect)

	chanMsgApis := chi.NewRouter()
	chanMsgApis.Post("/send", middlewares.Auth(wc.SendChanMessage))
	chanMsgApis.Patch("/{messageId}", middlewares.Auth(wc.EditChanMessage))
	chanMsgApis.Delete("/{messageId}", middlewares.Auth(wc.DeleteChanMessage))

	directMsgApis := chi.NewRouter()
	directMsgApis.Post("/send", middlewares.Auth(wc.SendDirectMessage))
	directMsgApis.Patch("/{messageId}", middlewares.Auth(wc.EditDirectMessage))
	directMsgApis.Delete("/{messageId}", middlewares.Auth(wc.DeleteDirectMessage))

	wcRoute.Mount("/message", lib.InjectActiveSession(chanMsgApis))
	wcRoute.Mount("/direct-message", lib.InjectActiveSession(directMsgApis))

	return wcRoute
}
