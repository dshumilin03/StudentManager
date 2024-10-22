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

type GroupResponse struct {
	resp.Response
}

func CreateGroup(service services.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req CreateGroupRequest

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

		if req.GroupNumber != "" {
			group, err := service.Create(context.Background(), req.GroupNumber)
			if err != nil {
				if err.Error() == "group already exists" {
					w.WriteHeader(http.StatusBadRequest)
					render.JSON(w, r, resp.Error("group already exists"))
					return
				}

				// TODO disable double notification from service and handler
				log.Printf("failed to create group %v", err)

				render.JSON(w, r, resp.Error("failed to create group"))
				return
			}
			log.Println("group added", slog.Any("response", req))

			responseCreatedGroup(w, r, group)

		} else {
			log.Println("invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
		}
	}
}

func GetAllGroups(service services.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		groups, err := service.GetAll(context.Background())
		if err != nil {
			log.Printf("failed to get groups %v", err)

			render.JSON(w, r, resp.Error("failed to get groups"))
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

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Println("request body decoded", slog.Any("request", req))

		group, err := service.GetById(context.Background(), req.Id)
		if err != nil {
			log.Printf("failed to get group %v", err)

			render.JSON(w, r, resp.Error("failed to get group"))
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

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Printf("failed to decode request body %v", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Println("request body decoded", slog.Any("response", req))

		group, err := service.Update(context.Background(), req.Id, req.GroupNumber)

		if err != nil {
			if err.Error() == "group doesn't exists" {
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("group doesn't exists"))
				return
			}

			log.Printf("failed to update group %v", err)

			render.JSON(w, r, resp.Error("failed to update group"))
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
			if err.Error() == "group does not exist" {
				w.WriteHeader(http.StatusNotFound)
				// TODO write instead of word group or student, domain
				render.JSON(w, r, resp.Error("group does not exist"))
				return
			}

			log.Printf("failed to delete group %v", err)

			render.JSON(w, r, resp.Error("failed to delete group"))
			return
		}

		log.Println("group deleted")

		w.WriteHeader(http.StatusNoContent)
	}
}

func responseFoundGroups(w http.ResponseWriter, r *http.Request, groups []domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, GroupResponse{
		Response: resp.FoundAllGroups(groups),
	})
}

func responseFoundGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, GroupResponse{
		Response: resp.FoundGroup(group),
	})
}

func responseCreatedGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, GroupResponse{
		Response: resp.CreatedGroup(group),
	})
}

func responseUpdatedGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, GroupResponse{
		Response: resp.UpdatedGroup(group),
	})
}
