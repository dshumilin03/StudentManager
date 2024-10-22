package app

import (
	"StudentManager/internal/config"
	"StudentManager/internal/http/handlers"
	"StudentManager/internal/http/services"
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
	// TODO provide just repos
	appServices := services.NewServices(repos.Students, repos.Groups)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/students", func(r chi.Router) {
		studentService := appServices.Students
		r.Post("/", handlers.CreateStudent(studentService))
		r.Get("/", handlers.GetAllStudents(studentService))

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", handlers.GetStudentById(studentService)) //TODO add path variable
			r.Delete("/", handlers.DeleteStudentById(studentService))
			r.Put("/", handlers.UpdateStudent(studentService))
		})
	})

	// Я закончил на добавлении групп надо потестить запросы к ним
	/* 1) Доделать группы (проверить при создании студента есть ли группа в бд, также добавить проверку при апдейте студента
	   2) Доделать все TODO
	   3) По-хорошему написать тесты бы и функциональные и юнит
	*/
	r.Route("/groups", func(r chi.Router) {
		groupService := appServices.Groups
		r.Post("/", handlers.CreateGroup(groupService))
		r.Get("/", handlers.GetAllGroups(groupService))

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", handlers.GetGroupById(groupService)) //TODO add path variable
			r.Delete("/", handlers.DeleteGroupById(groupService))
			r.Put("/", handlers.UpdateGroup(groupService))
		})
	})

	log.Println("Listening on server address " + cfg.Address)
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
