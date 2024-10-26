package handler

import (
	"StudentManager/internal/custom_errors"
	"StudentManager/internal/domain"
	"StudentManager/internal/dto"
	resp "StudentManager/internal/http/response"
	"StudentManager/internal/http/service"
	"context"
	"errors"
	"github.com/go-chi/render"
	"log"
	"log/slog"
	"net/http"
)

type CreateStudentRequest struct {
	FullName    string `json:"full_name" env-required:"true"`
	Age         int    `json:"age" env-required:"true"`
	GroupNumber string `json:"group_number"`
	Email       string `json:"email" env-required:"true"`
}

type UpdateStudentRequest struct {
	Id          int64  `json:"id" env-required:"true"`
	FullName    string `json:"full_name" env-required:"true"`
	Age         int    `json:"age" env-required:"true"`
	GroupNumber string `json:"group_number"`
	Email       string `json:"email" env-required:"true"`
}

type StudentIdRequest struct {
	Id int64 `json:"id" env-required:"true"`
}

type GetStudentRequest struct {
	FullName string `json:"full_name" env-required:"true"`
}

type StudentHandlerImpl struct {
	service service.StudentService
}

func NewStudentHandlerImpl(service service.StudentService) *StudentHandlerImpl {
	return &StudentHandlerImpl{service}
}

func (h *StudentHandlerImpl) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		var req CreateStudentRequest

		d := JsonDecoder[CreateStudentRequest, dto.GroupDto, domain.Group]{}
		err := d.Decode(w, r, req, h)
		if err != nil {
			return
		}

		if req.Age == 0 || req.Email == "" || req.FullName == "" || req.GroupNumber == "" {
			log.Println("invalid request")

			h.ResponseError(w, r, "invalid request", http.StatusBadRequest)
			return
		}

		studentDto := dto.StudentDto{
			FullName:    req.FullName,
			Age:         req.Age,
			GroupNumber: req.GroupNumber,
			Email:       req.Email,
		}

		student, err := studentService.Create(context.Background(), studentDto)
		if err != nil {
			if errors.Is(err, custom_errors.ErrStudentExists) || errors.Is(err, custom_errors.ErrGroupExists) {
				h.ResponseError(w, r, err.Error(), http.StatusBadRequest)
				return
			}

			h.ResponseError(w, r, "failed to create student", http.StatusInternalServerError)
			return
		}

		h.responseStudentCreated(w, r, student)
	}
}

func (h *StudentHandlerImpl) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		students, err := studentService.GetAll(context.Background())
		if err != nil {

			h.ResponseError(w, r, "failed to get student", http.StatusNotFound)
			return
		}

		h.responseFoundStudents(w, r, students)
	}
}

func (h *StudentHandlerImpl) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		var req StudentIdRequest

		d := JsonDecoder[StudentIdRequest, dto.GroupDto, domain.Group]{}
		err := d.Decode(w, r, req, h)
		if err != nil {
			return
		}

		student, err := studentService.GetById(context.Background(), req.Id)
		if err != nil {
			if errors.Is(err, custom_errors.ErrStudentNotFound) {
				h.ResponseError(w, r, err.Error(), http.StatusNotFound)
			}

			h.ResponseError(w, r, "failed to get student", http.StatusNotFound)
			return
		}

		h.responseFoundStudent(w, r, student)
	}
}

func (h *StudentHandlerImpl) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		var req UpdateStudentRequest

		d := JsonDecoder[UpdateStudentRequest, dto.GroupDto, domain.Group]{}
		err := d.Decode(w, r, req, h)
		if err != nil {
			return
		}

		log.Println("request body decoded", slog.Any("response", req))

		studentDto := dto.StudentDto{
			Id:          req.Id,
			FullName:    req.FullName,
			Age:         req.Age,
			GroupNumber: req.GroupNumber,
			Email:       req.Email,
		}

		student, err := studentService.Update(context.Background(), studentDto)

		if err != nil {
			if errors.Is(err, custom_errors.ErrStudentNotFound) {
				h.ResponseError(w, r, err.Error(), http.StatusNotFound)
				return
			}

			h.ResponseError(w, r, "failed to update student", http.StatusInternalServerError)
			return
		}

		h.responseStudentUpdated(w, r, student)
	}
}

func (h *StudentHandlerImpl) DeleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service
		var req StudentIdRequest

		d := JsonDecoder[StudentIdRequest, dto.GroupDto, domain.Group]{}
		err := d.Decode(w, r, req, h)
		if err != nil {
			return
		}

		err = studentService.DeleteById(context.Background(), req.Id)
		if err != nil {

			if errors.Is(err, custom_errors.ErrStudentNotFound) {

				h.ResponseError(w, r, err.Error(), http.StatusNotFound)
				return
			}

			h.ResponseError(w, r, "failed to delete student", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *StudentHandlerImpl) responseFoundStudents(w http.ResponseWriter, r *http.Request, students []domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.StudentsResponse(students))
}

func (h *StudentHandlerImpl) responseFoundStudent(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.StudentResponse(student))
}

func (h *StudentHandlerImpl) responseStudentCreated(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resp.StudentResponse(student))
}

func (h *StudentHandlerImpl) responseStudentUpdated(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.StudentResponse(student))
}

func (h *StudentHandlerImpl) ResponseError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, resp.Error(msg))
}
