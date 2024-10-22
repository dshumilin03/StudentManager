package app

import (
	"StudentManager/internal/config"
	"StudentManager/internal/http/handler"
	"StudentManager/internal/http/service"
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
	appServices := service.NewServices(repos)
	handlers := handler.NewHandlers(appServices)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers.InitRoutes(r)

	// Я закончил на добавлении групп надо потестить запросы к ним
	/* 1) Доделать группы (проверить при создании студента есть ли группа в бд, также добавить проверку при апдейте студента
	   2) Доделать все TODO
	   3) По-хорошему написать тесты бы и функциональные и юнит
	*/

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	log.Println("Listening on server address " + cfg.Address)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server")
	}

	log.Fatalf("server stopped")
}
