package handlers

import (
	"StudentManager/internal/domain"
	resp "StudentManager/internal/http/response"
	"StudentManager/internal/repository"
	"context"
	"errors"
	"github.com/go-chi/render"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type Request struct {
	FullName    string `json:"full_name" env-required:"true"`
	Age         int    `json:"age" env-required:"true"`
	GroupNumber string `json:"group_number"`
	Email       string `json:"email" env-required:"true"`
}

type Response struct {
	resp.Response
}

func CreateStudent(repo repository.StudentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request

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

		student := domain.Student{
			FullName:    req.FullName,
			Age:         req.Age,
			GroupNumber: req.GroupNumber,
			Email:       req.Email,
		}
		_ = repo.Create(context.Background(), student)
		// TODO check if exists

		//TODO check other exceptions

		log.Println("student added", slog.Any("response", req))

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
