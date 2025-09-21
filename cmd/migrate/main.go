package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"subscription-service/internal/env"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction^ 'up' or 'down'")
	}

	direction := os.Args[1]

	db, err := sql.Open("postgres", env.GetPostgresDSN())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
	if err != nil {
		log.Fatal(err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a version number to force")
		}
		var version uint
		_, err := fmt.Sscan(os.Args[2], &version)
		if err != nil {
			log.Fatal("Invalid version number")
		}
		if err := m.Force(int(version)); err != nil {
			log.Fatal(err)
		}
		log.Printf("Database dirty state cleared. Forced version: %d\n", version)
	default:
		log.Fatal("Please provide a migration direction: 'up', 'down' or 'force'")
	}
}
