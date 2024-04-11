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
	lib.NewChannelHub()

	go lib.HubChannel.Run()
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

	addr := "127.0.0.1:" + lib.EnvConfig.Port

	serve := http.Server{
		Handler: router,
		Addr:    addr,
		// Addr:    ":" + lib.EnvConfig.Port,
	}

	log.Printf("Http server running on http://%s", addr)
	log.Printf("Websocket server running on ws://%s", addr)

	apiRouter := chi.NewRouter()
	apiRouter.Mount("/health", routes.RegisterHealthRoutes())
	apiRouter.Mount("/profile", routes.RegisterProfileRoutes())
	apiRouter.Mount("/server", routes.RegisterServerRoutes())
	apiRouter.Mount("/member", routes.RegisterMemberRoutes())
	apiRouter.Mount("/channel", routes.RegisterChannelRoutes())
	apiRouter.Mount("/conversation", routes.RegisterConversationRoutes())

	wsRouter := chi.NewRouter()
	wsRouter.Mount("/", routes.RegisterWsRoutes())

	router.Mount("/api", lib.InjectActiveSession(apiRouter))
	router.Mount("/ws", wsRouter)

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
