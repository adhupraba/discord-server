package lib

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type envConfig struct {
	Port               string
	DbUrl              string
	Env                string
	CorsAllowedOrigins []string
	ClerkSecretKey     string
}

var EnvConfig envConfig

func LoadEnv() {
	var env string = os.Getenv("ENV")
	var port string = os.Getenv("PORT")
	log.Println("env, port =>", env, port)

	if env != "production" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Unable to load .env:", err)
		}
	}

	if env == "" {
		env = "development"
	}

	if port == "" {
		port = os.Getenv("ENV_PORT")
	}

	if port == "" {
		port = "8000"
	}

	EnvConfig = envConfig{
		Port:               os.Getenv("ENV_PORT"),
		DbUrl:              os.Getenv("DB_URL"),
		Env:                env,
		CorsAllowedOrigins: []string{},
		ClerkSecretKey:     os.Getenv("CLERK_SECRET_KEY"),
	}

	if EnvConfig.Env == "" {
		EnvConfig.Env = "development"
	}

	if os.Getenv("CORS_ALLOWED_ORIGINGS") == "" {
		EnvConfig.CorsAllowedOrigins = []string{"*"}
	} else {
		EnvConfig.CorsAllowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINGS"), " ")
	}

	if EnvConfig.Port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	if EnvConfig.DbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	if EnvConfig.ClerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY is not found in the environment")
	}
}
