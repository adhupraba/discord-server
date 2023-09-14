package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/adhupraba/discord-server/controllers"
)

func RegisterHealthRoutes() *chi.Mux {
	hr := controllers.HealthController{}
	healthRoute := chi.NewRouter()

	healthRoute.Get("/", hr.Health)

	return healthRoute
}
