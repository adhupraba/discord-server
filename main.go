package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/routes"
)

func init() {
	lib.LoadEnv()
	lib.ConnectDb()
}

func main() {
	if lib.DB != nil {
		defer lib.DB.Close()
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   lib.EnvConfig.CorsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	serve := http.Server{
		Handler: router,
		Addr:    ":" + lib.EnvConfig.Port,
	}

	log.Printf("Server listening on port %s", lib.EnvConfig.Port)

	apiRouter := chi.NewRouter()

	apiRouter.Mount("/health", routes.RegisterHealthRoutes())

	router.Mount("/api", apiRouter)

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
