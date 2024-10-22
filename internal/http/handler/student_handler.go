package handler

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/dto"
	resp "StudentManager/internal/http/response"
	"StudentManager/internal/http/service"
	"context"
	"errors"
	"github.com/go-chi/render"
	"io"
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

type StudentHandler struct {
	service service.StudentService
}

func NewStudentHandler(service service.StudentService) *StudentHandler {
	return &StudentHandler{service}
}

func (h *StudentHandler) CreateStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		var req CreateStudentRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			h.responseError(w, r, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to decode request body: %v", err)

			h.responseError(w, r, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Println("request body decoded", slog.Any("response", req))

		if req.Age == 0 || req.Email == "" || req.FullName == "" || req.GroupNumber == "" {
			log.Println("invalid request")

			h.responseError(w, r, "invalid request", http.StatusBadRequest)
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
			if err.Error() == "student already exists" {
				h.responseError(w, r, "student already exists", http.StatusBadRequest)
				return
			} else if err.Error() == "group doesn't exist" {
				h.responseError(w, r, "group doesn't exist", http.StatusBadRequest)
				return
			}

			h.responseError(w, r, "failed to create student", http.StatusInternalServerError)
			return
		}

		h.responseStudentCreated(w, r, student)
	}
}

func (h *StudentHandler) GetAllStudents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		students, err := studentService.GetAll(context.Background())
		if err != nil {

			h.responseError(w, r, "failed to get student", http.StatusNotFound)
			return
		}

		h.responseFoundStudents(w, r, students)
	}
}

func (h *StudentHandler) GetStudentById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		var req StudentIdRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			h.responseError(w, r, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to decode request body: %v", err)

			h.responseError(w, r, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Println("request body decoded", slog.Any("request", req))

		student, err := studentService.GetById(context.Background(), req.Id)
		if err != nil {
			if err.Error() == "student doesn't exist" {
				h.responseError(w, r, "student doesn't exist", http.StatusNotFound)
			}

			h.responseError(w, r, "failed to get student", http.StatusNotFound)
			return
		}

		h.responseFoundStudent(w, r, student)
	}
}

func (h *StudentHandler) UpdateStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		var req UpdateStudentRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			h.responseError(w, r, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to decode request body: %v", err)

			h.responseError(w, r, "failed to decode request", http.StatusBadRequest)
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
			if err.Error() == "student doesn't exist" {

				h.responseError(w, r, "student doesn't exist", http.StatusNotFound)
				return
			}

			h.responseError(w, r, "failed to update student", http.StatusInternalServerError)
			return
		}

		h.responseStudentUpdated(w, r, student)
	}
}

func (h *StudentHandler) DeleteStudentById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentService := h.service

		var req StudentIdRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			h.responseError(w, r, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to decode request body: %v", err)

			h.responseError(w, r, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Println("request body decoded", slog.Any("request", req))

		err = studentService.DeleteById(context.Background(), req.Id)
		if err != nil {

			if err.Error() == "student does not exist" {

				h.responseError(w, r, "student doesn't exist", http.StatusNotFound)
				return
			}

			h.responseError(w, r, "failed to delete student", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *StudentHandler) responseFoundStudents(w http.ResponseWriter, r *http.Request, students []domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.StudentsResponse(students))
}

func (h *StudentHandler) responseFoundStudent(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.StudentResponse(student))
}

func (h *StudentHandler) responseStudentCreated(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resp.StudentResponse(student))
}

func (h *StudentHandler) responseStudentUpdated(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.StudentResponse(student))
}

func (h *StudentHandler) responseError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, resp.Error(msg))
}
