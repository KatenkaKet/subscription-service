package main

import (
	"database/sql"
	"log"
	"subscription-service/internal/env"
	"subscription-service/internal/models"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type application struct {
	port      int
	jwtSecret string
	allModels models.Models
}

func main() {
	db, err := sql.Open("pgx", "postgres://myuser:123@localhost:5432/subscription_service?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createdModels := models.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", ""),
		allModels: createdModels,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}

}
