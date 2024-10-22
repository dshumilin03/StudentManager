package handlers

import (
	"StudentManager/internal/domain"
	resp "StudentManager/internal/http/response"
	"StudentManager/internal/http/services"
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

type StudentIdRequst struct {
	Id int64 `json:"id" env-required:"true"`
}

type GetStudentRequest struct {
	FullName string `json:"full_name" env-required:"true"`
}

type StudentResponse struct {
	resp.Response
}

func CreateStudent(service services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateStudentRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Println("request body decoded", slog.Any("response", req))

		if req.Age != 0 && req.Email != "" && req.FullName != "" && req.GroupNumber != "" {
			student, err := service.Create(context.Background(), req.FullName, req.Age, req.GroupNumber, req.Email)
			if err != nil {
				if err.Error() == "student already exists" {
					w.WriteHeader(http.StatusBadRequest)
					render.JSON(w, r, resp.Error("student already exists"))
					return
				}

				log.Printf("failed to create student %v", err)

				render.JSON(w, r, resp.Error("failed to create student"))
				return
			}
			log.Println("student added", slog.Any("response", req))

			responseStudentCreated(w, r, student)
		} else {
			log.Println("invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
		}
	}
}

func GetAllStudents(service services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		students, err := service.GetAll(context.Background())
		if err != nil {
			log.Printf("failed to get students %v", err)

			render.JSON(w, r, resp.Error("failed to get students"))
			return
		}

		log.Println("received all students")

		responseFoundStudents(w, r, students)
	}
}

func GetStudentById(service services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req StudentIdRequst

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Println("request body decoded", slog.Any("request", req))

		student, err := service.GetById(context.Background(), req.Id)
		if err != nil {
			log.Printf("failed to get student %v", err)

			render.JSON(w, r, resp.Error("failed to get student"))
			return
		}

		log.Println("received student by id")

		responseFoundStudent(w, r, student)
	}
}

func UpdateStudent(service services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req UpdateStudentRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Println("request body decoded", slog.Any("response", req))

		student, err := service.Update(context.Background(), req.Id, req.FullName, req.Age, req.GroupNumber, req.Email)

		if err != nil {
			if err.Error() == "student doesn't exists" {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("student doesn't exists"))
				return
			}

			log.Printf("failed to update student %v", err)

			render.JSON(w, r, resp.Error("failed to update student"))
			return
		}

		log.Println("student updated", slog.Any("response", req))

		responseStudentUpdated(w, r, student)
	}
}

func DeleteStudentById(service services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req StudentIdRequst

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Println("request body decoded", slog.Any("request", req))

		err = service.DeleteById(context.Background(), req.Id)
		if err != nil {

			// TODO make own types of errors
			if err.Error() == "student does not exist" {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("student does not exist"))
				return
			}

			log.Printf("failed to delete student %v", err)

			render.JSON(w, r, resp.Error("failed to delete student"))
			return
		}

		log.Println("student deleted")

		w.WriteHeader(http.StatusNoContent)
	}
}

// TODO rewrite that responses into structs GroupResponse
func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, http.StatusOK)
}

func responseFoundStudents(w http.ResponseWriter, r *http.Request, students []domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, StudentResponse{
		Response: resp.FoundAllStudents(students),
	})
}

func responseFoundStudent(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, StudentResponse{
		Response: resp.FoundStudent(student),
	})
}

func responseStudentCreated(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, StudentResponse{
		Response: resp.StudentCreated(student),
	})
}

func responseStudentUpdated(w http.ResponseWriter, r *http.Request, student domain.Student) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, StudentResponse{
		Response: resp.StudentUpdated(student),
	})
}
