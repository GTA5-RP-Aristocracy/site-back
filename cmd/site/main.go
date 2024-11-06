//go:generate swagger generate spec -o swagger.json
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/GTA5-RP-Aristocracy/site-back/db"
	"github.com/GTA5-RP-Aristocracy/site-back/user"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/rs/zerolog"
)

func main() {
	// This is the entry point of the site command.
	// It starts the web server and listens for incoming requests

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Define the configuration options.
	dbConfig := db.Config{
		User:     "postgres",
		Password: "postgres",
		Host:     "localhost",
		Database: "gta_site",
	}
	err := env.Parse(&dbConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse the database configuration")
	}

	// Connect to the database.
	db, err := db.ConnectDB(dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to the database")
	}
	defer db.Close()

	logger.Info().Msg("connected to the database")

	// Create a new user repository.
	userRepo := user.NewRepository(db)

	// Create a new user service.
	userService := user.NewService(userRepo)

	// Create a new user http handler.
	userHandler := user.NewHandler(userService)

	loggerRouter := httplog.NewLogger("gta-site-api", httplog.Options{
		JSON:     true,
		LogLevel: slog.LevelDebug,
		Concise:  true,
		// RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "v0.0.1",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})

	// Start the web server.
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(loggerRouter))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	userHandler.RegisterUserRouter(r)

	// TODO add signal handling for graceful shutdown
	logger.Info().Msg("starting the web server")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatal().Err(err).Msg("failed to start the web server")
	}
}
