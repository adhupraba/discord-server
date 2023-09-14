package lib

import (
	"log"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

var ClerkClient clerk.Client
var InjectActiveSession func(handler http.Handler) http.Handler

func InitClerkClient() {
	client, err := clerk.NewClient(EnvConfig.ClerkSecretKey)

	if err != nil {
		log.Fatal("error initialising clerk client =>", err)
	}

	ClerkClient = client
	InjectActiveSession = clerk.WithSession(client)
}
