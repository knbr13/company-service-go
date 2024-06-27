package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/knbr13/company-service-go/config"
	"github.com/knbr13/company-service-go/internal/handlers"
)

type app struct {
	hndlrs *handlers.Handlers
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load.env file: %s\n", err.Error())
	}

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err.Error())
	}

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to open database: %s\n", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %s\n", err.Error())
	}

	app := &app{
		hndlrs: handlers.NewHandlers(db, cfg),
	}

	log.Println("Starting app server on port :8080")
	http.ListenAndServe(":8080", app.SetupRoutes())
}
