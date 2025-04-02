package utils

import (
	"database/sql"
	"log"
)

func ConnectDB(dsn string) *sql.DB {
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Couldn't connect to db: %v", err)
	}
	if err = database.Ping(); err != nil {
		log.Fatalf("Db is unaccessible: %v", err)
	}
	return database
}
