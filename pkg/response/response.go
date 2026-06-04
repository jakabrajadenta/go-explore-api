package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jakabrajadenta/go-explore-api/pkg/logger"
)

// Response is the standard envelope for all API responses.
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Meta    Meta   `json:"meta"`
}

// ListResponse wraps paginated list results.
type ListResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    any      `json:"data"`
	Errors  any      `json:"errors,omitempty"`
	Meta    ListMeta `json:"meta"`
}

type Meta struct {
	TraceID   string `json:"trace_id,omitempty"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
}

type ListMeta struct {
	TraceID    string     `json:"trace_id,omitempty"`
	Timestamp  string     `json:"timestamp"`
	Path       string     `json:"path"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ── 2xx ──────────────────────────────────────────────────────

func OK(w http.ResponseWriter, r *http.Request, message string, data any) {
	write(w, http.StatusOK, Response{
		Success: true, Message: message, Data: data, Meta: newMeta(r),
	})
}

func Created(w http.ResponseWriter, r *http.Request, message string, data any) {
	write(w, http.StatusCreated, Response{
		Success: true, Message: message, Data: data, Meta: newMeta(r),
	})
}

func OKList(w http.ResponseWriter, r *http.Request, message string, data any, page, perPage, total int) {
	totalPages := 0
	if perPage > 0 {
		totalPages = total / perPage
		if total%perPage != 0 {
			totalPages++
		}
	}
	write(w, http.StatusOK, ListResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: ListMeta{
			TraceID:   logger.TraceIDFrom(r.Context()),
			Timestamp: now(),
			Path:      r.URL.Path,
			Pagination: Pagination{
				Page: page, PerPage: perPage,
				Total: total, TotalPages: totalPages,
			},
		},
	})
}

// ── 4xx ──────────────────────────────────────────────────────

func BadRequest(w http.ResponseWriter, r *http.Request, errors any) {
	write(w, http.StatusBadRequest, Response{
		Success: false, Message: "Validation failed", Errors: errors, Meta: newMeta(r),
	})
}

func NotFound(w http.ResponseWriter, r *http.Request, message string) {
	write(w, http.StatusNotFound, Response{
		Success: false, Message: message, Meta: newMeta(r),
	})
}

func Conflict(w http.ResponseWriter, r *http.Request, message string) {
	write(w, http.StatusConflict, Response{
		Success: false, Message: message, Meta: newMeta(r),
	})
}

// ── 5xx ──────────────────────────────────────────────────────

func InternalError(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusInternalServerError, Response{
		Success: false, Message: "Internal server error", Meta: newMeta(r),
	})
}

// ── helpers ──────────────────────────────────────────────────

func newMeta(r *http.Request) Meta {
	return Meta{
		TraceID:   logger.TraceIDFrom(r.Context()),
		Timestamp: now(),
		Path:      r.URL.Path,
	}
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
