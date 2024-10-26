package handler

import (
	"errors"
	"github.com/go-chi/render"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type JsonDecoder[T any, Dto any, Domain any] struct{}

func (j JsonDecoder[T, Dto, Domain]) Decode(w http.ResponseWriter, r *http.Request, data T, h Handler[Dto, Domain]) error {

	err := render.DecodeJSON(r.Body, &data)
	if errors.Is(err, io.EOF) {

		log.Println("request body is empty")

		h.ResponseError(w, r, "empty request", http.StatusBadRequest)
		return err
	}
	if err != nil {
		log.Printf("failed to decode request body: %v", err)

		h.ResponseError(w, r, "failed to decode request", http.StatusBadRequest)

		return err
	}

	log.Println("request body decoded", slog.Any("request", data))
	return nil
}
