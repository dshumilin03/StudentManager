package handler

import (
	"StudentManager/internal/http/service"
	"github.com/go-chi/chi/v5"
	"log"
)

type Handlers struct {
	Students StudentHandler
	Groups   GroupHandler
}

func NewHandlers(services *service.Services) *Handlers {
	log.Printf("Handlers are created")
	return &Handlers{
		Students: *NewStudentHandler(services.Students),
		Groups:   *NewGroupHandler(services.Groups),
	}
}

func (h *Handlers) InitRoutes(r chi.Router) {

	r.Route("/students", func(r chi.Router) {
		studentHandler := h.Students
		r.Post("/", studentHandler.CreateStudent())
		r.Get("/", studentHandler.GetAllStudents())

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", studentHandler.GetStudentById()) //TODO add path variable
			r.Delete("/", studentHandler.DeleteStudentById())
			r.Put("/", studentHandler.UpdateStudent())
		})
	})

	r.Route("/groups", func(r chi.Router) {
		groupHandler := h.Groups
		r.Post("/", groupHandler.CreateGroup())
		r.Get("/", groupHandler.GetAllGroups())

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", groupHandler.GetGroupById()) //TODO add path variable
			r.Delete("/", groupHandler.DeleteGroupById())
			r.Put("/", groupHandler.UpdateGroup())
		})
	})
}
