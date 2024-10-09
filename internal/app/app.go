package app

import (
	"StudentManager/internal/config"
	"StudentManager/internal/http/handlers"
	"StudentManager/internal/repository"
	"StudentManager/pkg/database/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func Run() {
	cfg := config.Init()
	db, err := postgres.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
		return
	}
	repos := repository.NewRepositories(db)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/home", handlers.CreateStudent(repos.Students))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server")
	}

	log.Fatalf("server stopped")
}
