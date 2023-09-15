package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
)

func RegisterHealthRoutes() *chi.Mux {
	hc := controllers.HealthController{}
	healthRoute := chi.NewRouter()

	healthRoute.Get("/", hc.Health)

	return healthRoute
}
