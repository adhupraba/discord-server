package lib

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func ConnectDb() {
	log.Println("db url", EnvConfig.DbUrl)
	conn, err := sql.Open("postgres", EnvConfig.DbUrl)

	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	DB = conn
}
