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

type GroupHandlerImpl struct {
	service service.GroupService
}

func NewGroupHandlerImpl(service service.GroupService) *GroupHandlerImpl {
	return &GroupHandlerImpl{
		service: service,
	}
}

func (h *GroupHandlerImpl) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service
		var req CreateGroupRequest

		d := JsonDecoder[CreateGroupRequest, dto.GroupDto, domain.Group]{}
		err := d.Decode(w, r, req, h)
		if err != nil {
			return
		}

		if req.GroupNumber != "" {

			groupDto := dto.GroupDto{
				GroupNumber: req.GroupNumber,
			}

			group, err := groupService.Create(context.Background(), groupDto)
			if err != nil {
				if errors.Is(err, custom_errors.ErrGroupExists) {

					h.ResponseError(w, r, err.Error(), http.StatusBadRequest)
					return
				}

				h.ResponseError(w, r, "failed to create group", http.StatusInternalServerError)
				return
			}

			h.responseCreatedGroup(w, r, group)

		} else {
			log.Println("invalid request")

			h.ResponseError(w, r, "invalid request", http.StatusBadRequest)
			return
		}
	}
}

func (h *GroupHandlerImpl) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service

		groups, err := groupService.GetAll(context.Background())
		if err != nil {

			h.ResponseError(w, r, "failed to get groups", http.StatusInternalServerError)
			return
		}

		h.responseFoundGroups(w, r, groups)
	}
}

func (h *GroupHandlerImpl) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service

		var req GroupIdRequest

		d := JsonDecoder[GroupIdRequest, dto.GroupDto, domain.Group]{}
		err := d.Decode(w, r, req, h)
		if err != nil {
			return
		}

		group, err := groupService.GetById(context.Background(), req.Id)
		if err != nil {
			if errors.Is(err, custom_errors.ErrGroupNotFound) {
				h.ResponseError(w, r, err.Error(), http.StatusNotFound)
			}

			h.ResponseError(w, r, "failed to get group", http.StatusNotFound)
			return
		}

		h.responseFoundGroup(w, r, group)
	}
}

func (h *GroupHandlerImpl) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service
		// TODO write json decoder struct
		var req UpdateGroupRequest

		d := JsonDecoder[UpdateGroupRequest, dto.GroupDto, domain.Group]{}
		err := d.Decode(w, r, req, h)
		if err != nil {
			return
		}

		groupDto := dto.GroupDto{
			Id:          req.Id,
			GroupNumber: req.GroupNumber,
		}

		group, err := groupService.Update(context.Background(), groupDto)

		if err != nil {
			if errors.Is(err, custom_errors.ErrGroupNotFound) {
				h.ResponseError(w, r, err.Error(), http.StatusNotFound)
				return
			}

			h.ResponseError(w, r, "failed to update group", http.StatusInternalServerError)
			return
		}

		h.responseUpdatedGroup(w, r, group)
	}
}

func (h *GroupHandlerImpl) DeleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupService := h.service

		var req GroupIdRequest

		d := JsonDecoder[GroupIdRequest, dto.GroupDto, domain.Group]{}
		err := d.Decode(w, r, req, h)
		if err != nil {
			return
		}

		err = groupService.DeleteById(context.Background(), req.Id)
		if err != nil {

			if errors.Is(err, custom_errors.ErrGroupNotFound) {
				w.WriteHeader(http.StatusNotFound)
				h.ResponseError(w, r, err.Error(), http.StatusNotFound)
				return
			}

			h.ResponseError(w, r, "failed to delete group", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *GroupHandlerImpl) responseFoundGroups(w http.ResponseWriter, r *http.Request, groups []domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupsResponse(groups))
}

func (h *GroupHandlerImpl) responseFoundGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupResponse(group))
}

func (h *GroupHandlerImpl) responseCreatedGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, resp.GroupResponse(group))
}

func (h *GroupHandlerImpl) responseUpdatedGroup(w http.ResponseWriter, r *http.Request, group domain.Group) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.GroupResponse(group))
}

func (h *GroupHandlerImpl) ResponseError(w http.ResponseWriter, r *http.Request, msg string, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, resp.Error(msg))
}
