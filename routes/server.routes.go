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
	serverRoute.Get("/user-servers", middlewares.Auth(sc.GetUserMemberServers))
	serverRoute.Get("/{serverId}", middlewares.Auth(sc.GetServer))
	serverRoute.Patch("/{serverId}", middlewares.Auth(sc.UpdateServer))
	serverRoute.Get("/{serverId}/channels-and-members", middlewares.Auth(sc.GetFullServerDetails))
	serverRoute.Patch("/{serverId}/invite-code", middlewares.Auth(sc.UpdateInviteCode))
	serverRoute.Patch("/{inviteCode}/verify", middlewares.Auth(sc.VerifyAndAcceptInviteCode))
	serverRoute.Patch("/{serverId}/leave", middlewares.Auth(sc.MemberLeaveServer))
	serverRoute.Delete("/{serverId}", middlewares.Auth(sc.DeleteServer))
	serverRoute.Get("/{serverId}/general-channel", middlewares.Auth(sc.GetServerGeneralChannel))

	return serverRoute
}
