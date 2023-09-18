package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
	"github.com/adhupraba/discord-server/middlewares"
)

func RegisterMemberRoutes() *chi.Mux {
	mc := controllers.MemberController{}
	memberRoute := chi.NewRouter()

	memberRoute.Patch("/{memberId}", middlewares.Auth(mc.UpdateMemberRole))
	memberRoute.Delete("/{memberId}", middlewares.Auth(mc.RemoveServerMember))

	return memberRoute
}
