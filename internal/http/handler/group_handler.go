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

type GetGroupRequest struct {
	GroupNumber string `json:"group_number" env-required:"true"`
}

type GroupHandler struct {
	service service.GroupService
}

func NewGroupHandler(service service.GroupService) *GroupHandler {
	return &GroupHandler{
		service: service,
	}
}

func (h *GroupHandler) CreateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service
		var req CreateGroupRequest

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

		if req.GroupNumber != "" {

			groupDto := dto.GroupDto{
				GroupNumber: req.GroupNumber,
			}

			group, err := groupService.Create(context.Background(), groupDto)
			if err != nil {
				if err.Error() == "group already exists" {

					h.responseError(w, r, err.Error(), http.StatusBadRequest)
					return
				}

				h.responseError(w, r, "failed to create group", http.StatusInternalServerError)
				return
			}

			h.responseCreatedGroup(w, r, group)

		} else {
			log.Println("invalid request")

			h.responseError(w, r, "invalid request", http.StatusBadRequest)
			return
		}
	}
}

func (h *GroupHandler) GetAllGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service

		groups, err := groupService.GetAll(context.Background())
		if err != nil {

			h.responseError(w, r, "failed to get groups", http.StatusInternalServerError)
			return
		}

		h.responseFoundGroups(w, r, groups)
	}
}

func (h *GroupHandler) GetGroupById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service

		var req GroupIdRequest

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

		group, err := groupService.GetById(context.Background(), req.Id)
		if err != nil {

			h.responseError(w, r, "failed to get group", http.StatusNotFound)
			return
		}

		h.responseFoundGroup(w, r, group)
	}
}

func (h *GroupHandler) UpdateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service
		// TODO write json decoder struct
		var req UpdateGroupRequest

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

		groupDto := dto.GroupDto{
			Id:          req.Id,
			GroupNumber: req.GroupNumber,
		}

		group, err := groupService.Update(context.Background(), groupDto)

		if err != nil {
			if err.Error() == "group doesn't exist" {
				h.responseError(w, r, "group doesn't exist", http.StatusNotFound)
				return
			}

			h.responseError(w, r, "failed to update group", http.StatusInternalServerError)
			return
		}

		h.responseUpdatedGroup(w, r, group)
	}
}

func (h *GroupHandler) DeleteGroupById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service

		var req GroupIdRequest

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

		err = groupService.DeleteById(context.Background(), req.Id)
		if err != nil {

			// TODO make own types of errors
			if err.Error() == "group does not exist" {
				w.WriteHeader(http.StatusNotFound)
				// TODO write instead of word group or student, domain
				h.responseError(w, r, "group doesn't exist", http.StatusNotFound)
				return
			}

			h.responseError(w, r, "failed to delete group", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *GroupHandler) responseFoundGroups(w http.ResponseWriter, r *http.Request, groups []domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupsResponse(groups))
}

func (h *GroupHandler) responseFoundGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupResponse(group))
}

func (h *GroupHandler) responseCreatedGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resp.GroupResponse(group))
}

func (h *GroupHandler) responseUpdatedGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupResponse(group))
}

func (h *GroupHandler) responseError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, resp.Error(msg))
}
