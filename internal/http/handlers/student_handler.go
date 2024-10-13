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

type GetStudentRequest struct {
	FullName string `json:"full_name" env-required:"true"`
}

type Response struct {
	resp.Response
}

func CreateStudent(service services.StudentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateStudentRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("response body is empty")

			render.JSON(w, r, resp.Error("empty response"))

			return
		}
		if err != nil {
			log.Printf("failed to decode response body %v", err)

			render.JSON(w, r, resp.Error("failed to decode response"))

			return
		}

		log.Println("response body decoded", slog.Any("response", req))

		student, err := service.Create(context.Background(), req.FullName, req.Age, req.GroupNumber, req.Email)

		if err != nil {
			log.Printf("failed to create student %v", err)

			render.JSON(w, r, resp.Error("failed to create student"))
			return
		}

		log.Println("student added", slog.Any("response", req))

		responseCreated(w, r, student)
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

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}

func responseFoundStudents(w http.ResponseWriter, r *http.Request, students []domain.Student) {
	render.JSON(w, r, Response{
		Response: resp.FoundAll(students),
	})
}

func responseCreated(w http.ResponseWriter, r *http.Request, student domain.Student) {
	render.JSON(w, r, Response{
		Response: resp.Created(student),
	})
}
