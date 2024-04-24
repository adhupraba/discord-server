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

	wcApis := chi.NewRouter()
	wcApis.Post("/send", middlewares.Auth(wc.SendMessage))
	wcApis.Patch("/{messageId}", middlewares.Auth(wc.EditMessage))
	wcApis.Delete("/{messageId}", middlewares.Auth(wc.DeleteMessage))

	wcRoute.Mount("/message", lib.InjectActiveSession(wcApis))

	return wcRoute
}
