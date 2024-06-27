package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"github.com/knbr13/company-service-go/config"
	"github.com/knbr13/company-service-go/internal/handlers"
	"github.com/knbr13/company-service-go/internal/kafka"
)

type app struct {
	hndlrs   *handlers.Handlers
	cfg      *config.Config
	producer sarama.SyncProducer
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %s\n", err.Error())
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

	producer, err := kafka.ConnectProducer([]string{cfg.KafkaBroker})
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %s\n", err.Error())
	}

	errCh := make(chan error, 32)

	app := &app{
		hndlrs:   handlers.NewHandlers(db, cfg, producer, errCh),
		cfg:      cfg,
		producer: producer,
	}

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: app.SetupRoutes(),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		log.Println("caught signal", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- httpServer.Shutdown(ctx)
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server error: %s\n", err)
		}
	}()

	for {
		select {
		case err := <-errCh:
			log.Printf("Error received: %s\n", err)
		case err := <-shutdownError:
			if err != nil {
				log.Fatalf("HTTP server shutdown error: %s\n", err)
			}
		}
	}
}
