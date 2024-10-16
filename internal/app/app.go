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
	services := services.NewServices(repos.Students)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/students", func(r chi.Router) {
		studentService := services.Students
		r.Post("/", handlers.CreateStudent(studentService))
		r.Get("/", handlers.GetAllStudents(studentService))

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", handlers.GetStudentById(studentService)) //TODO add path variable
			r.Delete("/", handlers.DeleteStudentById(studentService))
			r.Put("/", handlers.UpdateStudent(studentService))
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
