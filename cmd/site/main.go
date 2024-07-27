package main

import (
	"net/http"
	"os"
	"time"

	"github.com/GTA5-RP-Aristocracy/site-back/db"
	"github.com/GTA5-RP-Aristocracy/site-back/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func main() {
	// This is the entry point of the site command.
	// It starts the web server and listens for incoming requests.

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Connect to the database.
	db, err := db.ConnectDB("postgres", "postgres", "localhost", "gta_site")
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

	// Start the web server.
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	userHandler.RegisterUserRouter(r)

	// TODO add signal handling for graceful shutdown
	logger.Info().Msg("starting the web server")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatal().Err(err).Msg("failed to start the web server")
	}
}
