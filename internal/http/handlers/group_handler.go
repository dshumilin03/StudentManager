package handlers

import (
	"StudentManager/internal/domain"
	"StudentManager/internal/dto"
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

type CreateGroupRequest struct {
	GroupNumber string `json:"group_number" env-required:"true"`
}

type UpdateGroupRequest struct {
	Id          int64  `json:"id" env-required:"true"`
	GroupNumber string `json:"group_number" env-required:"true"`
}

type GroupIdRequest struct {
	Id int64 `json:"id" env-required:"true"`
}

// TODO make StudentResponse also
type GetGroupRequest struct {
	GroupNumber string `json:"group_number" env-required:"true"`
}

func CreateGroup(service services.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateGroupRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Println("request body is empty")

			responseError(w, r, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			responseError(w, r, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Println("request body decoded", slog.Any("response", req))

		if req.GroupNumber != "" {

			groupDto := dto.GroupDto{
				GroupNumber: req.GroupNumber,
			}

			group, err := service.Create(context.Background(), groupDto)
			if err != nil {
				if err.Error() == "group already exists" {
					responseError(w, r, "group already exists", http.StatusBadRequest)
					return
				}

				// TODO disable double notification from service and handler
				log.Printf("failed to create group %v", err)

				responseError(w, r, "failed to create group", http.StatusInternalServerError)
				return
			}
			log.Println("group added", slog.Any("response", req))

			responseCreatedGroup(w, r, group)

		} else {
			log.Println("invalid request")
			responseError(w, r, "invalid request", http.StatusBadRequest)
			return
		}
	}
}

func GetAllGroups(service services.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		groups, err := service.GetAll(context.Background())
		if err != nil {
			log.Printf("failed to get groups %v", err)

			responseError(w, r, "failed to get groups", http.StatusInternalServerError)
			return
		}

		log.Println("received all groups")

		responseFoundGroups(w, r, groups)
	}
}

func GetGroupById(service services.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req GroupIdRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			responseError(w, r, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			responseError(w, r, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Println("request body decoded", slog.Any("request", req))

		group, err := service.GetById(context.Background(), req.Id)
		if err != nil {
			log.Printf("failed to get group %v", err)

			responseError(w, r, "failed to get group", http.StatusNotFound)
			return
		}

		log.Println("received group by id")

		responseFoundGroup(w, r, group)
	}
}

func UpdateGroup(service services.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req UpdateGroupRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			responseError(w, r, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			responseError(w, r, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Println("request body decoded", slog.Any("response", req))

		groupDto := dto.GroupDto{
			Id:          req.Id,
			GroupNumber: req.GroupNumber,
		}

		group, err := service.Update(context.Background(), groupDto)

		if err != nil {
			if err.Error() == "group doesn't exists" {
				responseError(w, r, "group doesn't exists", http.StatusNotFound)
				return
			}

			log.Printf("failed to update group %v", err)

			responseError(w, r, "failed to update group", http.StatusInternalServerError)
			return
		}

		log.Println("group updated", slog.Any("response", req))

		responseUpdatedGroup(w, r, group)
	}
}

func DeleteGroupById(service services.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req GroupIdRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			log.Println("request body is empty")

			responseError(w, r, "empty request", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			responseError(w, r, "failed to decode request", http.StatusBadRequest)

			return
		}

		log.Println("request body decoded", slog.Any("request", req))

		err = service.DeleteById(context.Background(), req.Id)
		if err != nil {

			// TODO make own types of errors
			if err.Error() == "group does not exist" {
				w.WriteHeader(http.StatusNotFound)
				// TODO write instead of word group or student, domain
				responseError(w, r, "group doesn't exists", http.StatusNotFound)
				return
			}

			log.Printf("failed to delete group %v", err)

			responseError(w, r, "failed to delete group", http.StatusInternalServerError)
			return
		}

		log.Println("group deleted")

		w.WriteHeader(http.StatusNoContent)
	}
}

func responseFoundGroups(w http.ResponseWriter, r *http.Request, groups []domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupsResponse(groups))
}

func responseFoundGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupResponse(group))
}

func responseCreatedGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resp.GroupResponse(group))
}

func responseUpdatedGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupResponse(group))
}

func responseError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, resp.Error(msg))
}
