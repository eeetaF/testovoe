package main

import (
	"os"

	_ "github.com/lib/pq"

	"testovoe/src/models/migration"
	"testovoe/src/route"
	"testovoe/src/utils"
)

func main() {
	dsn := os.Getenv("DB_DSN")

	migration.RunMigrations(dsn)

	db := utils.ConnectDB(dsn)

	route.InitRoutes(db)
}
