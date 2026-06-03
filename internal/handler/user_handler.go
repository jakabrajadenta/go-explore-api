package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jakabrajadenta/go-explore-api/internal/dto"
	"github.com/jakabrajadenta/go-explore-api/internal/service"
	"github.com/jakabrajadenta/go-explore-api/pkg/response"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	users, total, err := h.svc.GetAll(r.Context(), page, perPage)
	if err != nil {
		response.InternalError(w, r)
		return
	}
	response.OKList(w, r, "Users retrieved successfully", users, page, perPage, total)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.BadRequest(w, r, map[string]string{"id": "id must be a valid integer"})
		return
	}

	user, err := h.svc.GetByID(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, r, "User not found")
		return
	}
	if err != nil {
		response.InternalError(w, r)
		return
	}
	response.OK(w, r, "User retrieved successfully", user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, map[string]string{"body": "request body must be valid JSON"})
		return
	}
	if errs := req.Validate(); len(errs) > 0 {
		response.BadRequest(w, r, errs)
		return
	}

	user, err := h.svc.Create(r.Context(), req)
	if errors.Is(err, service.ErrUsernameConflict) {
		response.Conflict(w, r, "Username is already taken")
		return
	}
	if errors.Is(err, service.ErrEmailConflict) {
		response.Conflict(w, r, "Email is already registered")
		return
	}
	if err != nil {
		response.InternalError(w, r)
		return
	}
	response.Created(w, r, "User created successfully", user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.BadRequest(w, r, map[string]string{"id": "id must be a valid integer"})
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, map[string]string{"body": "request body must be valid JSON"})
		return
	}
	if errs := req.Validate(); len(errs) > 0 {
		response.BadRequest(w, r, errs)
		return
	}

	user, err := h.svc.Update(r.Context(), id, req)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, r, "User not found")
		return
	}
	if errors.Is(err, service.ErrUsernameConflict) {
		response.Conflict(w, r, "Username is already taken")
		return
	}
	if errors.Is(err, service.ErrEmailConflict) {
		response.Conflict(w, r, "Email is already registered")
		return
	}
	if err != nil {
		response.InternalError(w, r)
		return
	}
	response.OK(w, r, "User updated successfully", user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.BadRequest(w, r, map[string]string{"id": "id must be a valid integer"})
		return
	}

	err = h.svc.Delete(r.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(w, r, "User not found")
		return
	}
	if err != nil {
		response.InternalError(w, r)
		return
	}
	response.OK(w, r, "User deleted successfully", nil)
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
