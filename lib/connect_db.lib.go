package lib

import (
	"database/sql"
	"log"

	"github.com/adhupraba/discord-server/internal/queries"
)

var DB *queries.Queries
var SqlConn *sql.DB

func ConnectDb() {
	log.Println("db url", EnvConfig.DbUrl)
	conn, err := sql.Open("postgres", EnvConfig.DbUrl)

	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	DB = queries.New(conn)
	SqlConn = conn
}
