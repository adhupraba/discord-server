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
	lib.InitClerkClient()
}

func main() {
	if lib.SqlConn != nil {
		defer lib.SqlConn.Close()
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
		Addr:    "127.0.0.1:" + lib.EnvConfig.Port,
		// Addr:    ":" + lib.EnvConfig.Port,
	}

	log.Printf("Server listening on port %s", lib.EnvConfig.Port)

	apiRouter := chi.NewRouter()

	apiRouter.Mount("/health", routes.RegisterHealthRoutes())
	apiRouter.Mount("/profile", routes.RegisterProfileRoutes())
	apiRouter.Mount("/server", routes.RegisterServerRoutes())
	apiRouter.Mount("/member", routes.RegisterMemberRoutes())
	apiRouter.Mount("/channel", routes.RegisterChannelRoutes())
	apiRouter.Mount("/conversation", routes.RegisterConversationRoutes())

	router.Mount("/api", lib.InjectActiveSession(apiRouter))

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
