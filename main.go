package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	_ "github.com/lib/pq"

	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/routes"
)

func init() {
	lib.LoadEnv()
	lib.ConnectDb()
	lib.InitClerkClient()
	lib.NewHub()

	go lib.WsHub.Run()
}

func main() {
	if lib.SqlConn != nil {
		defer lib.SqlConn.Close()
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   lib.EnvConfig.CorsAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	router.Use(httprate.LimitByIP(100, 1*time.Minute))

	addr := "0.0.0.0:" + lib.EnvConfig.Port

	serve := http.Server{
		Handler: router,
		Addr:    addr,
		// Addr:    ":" + lib.EnvConfig.Port,
	}

	log.Printf("Http server running on http://%s", addr)
	log.Printf("Websocket server running on ws://%s", addr)

	apiRouter := chi.NewRouter()
	apiRouter.Mount("/profile", routes.RegisterProfileRoutes())
	apiRouter.Mount("/server", routes.RegisterServerRoutes())
	apiRouter.Mount("/member", routes.RegisterMemberRoutes())
	apiRouter.Mount("/channel", routes.RegisterChannelRoutes())
	apiRouter.Mount("/conversation", routes.RegisterConversationRoutes())
	apiRouter.Mount("/message", routes.RegisterMessageRoutes())
	apiRouter.Mount("/direct-message", routes.RegisterDirectMessageRoutes())

	wsRouter := chi.NewRouter()
	wsRouter.Mount("/", routes.RegisterWsRoutes())

	router.Mount("/api/health", routes.RegisterHealthRoutes())
	router.Mount("/api", lib.InjectActiveSession(apiRouter))
	router.Mount("/ws", wsRouter)

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
