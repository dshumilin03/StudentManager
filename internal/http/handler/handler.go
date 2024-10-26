package handler

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/dto"
	"StudentManager/internal/http/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Handler[Dto any, Domain any] interface {
	ResponseError(w http.ResponseWriter, r *http.Request, msg string, status int)
	Create() http.HandlerFunc
	GetAll() http.HandlerFunc
	GetById() http.HandlerFunc
	Update() http.HandlerFunc
	DeleteById() http.HandlerFunc
}

type StudentHandler interface {
	Handler[dto.StudentDto, domain.Student]
}

type GroupHandler interface {
	Handler[dto.GroupDto, domain.Group]
}

type Handlers struct {
	Students StudentHandler
	Groups   GroupHandler
}

func NewHandlers(services *service.Services) *Handlers {
	log.Printf("Handlers are created")
	return &Handlers{
		Students: NewStudentHandlerImpl(services.Students),
		Groups:   NewGroupHandlerImpl(services.Groups),
	}
}

func (h *Handlers) InitRoutes(r chi.Router) {

	r.Route("/students", func(r chi.Router) {
		studentHandler := h.Students
		r.Post("/", studentHandler.Create())
		r.Get("/", studentHandler.GetAll())

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", studentHandler.GetById()) //TODO add path variable
			r.Delete("/", studentHandler.DeleteById())
			r.Put("/", studentHandler.Update())
		})
	})

	r.Route("/groups", func(r chi.Router) {
		groupHandler := h.Groups
		r.Post("/", groupHandler.Create())
		r.Get("/", groupHandler.GetAll())

		r.Route("/{Id}", func(r chi.Router) {
			r.Get("/", groupHandler.GetById()) //TODO add path variable
			r.Delete("/", groupHandler.DeleteById())
			r.Put("/", groupHandler.Update())
		})
	})
}
