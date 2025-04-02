package migration

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	migrationsPath := "file://migrations/"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Couldn't run migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatalf("Can't run migration: %v", err)
	}

	log.Println("Migration completed")
}
