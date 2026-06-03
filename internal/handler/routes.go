package handler

import (
	"net/http"

	"github.com/jakabrajadenta/go-explore-api/internal/service"
)

func RegisterRoutes(mux *http.ServeMux, userSvc service.UserService) {
	// Root info
	mux.HandleFunc("GET /{$}", HandleInfo)

	// Health check
	mux.HandleFunc("GET /health", HandleHealth)

	// Echo (learning / debug endpoints)
	mux.HandleFunc("GET /echo", HandleEchoGet)
	mux.HandleFunc("POST /echo", HandleEchoPost)
	mux.HandleFunc("GET /echo/{message}", HandleEchoPath)

	// User management
	uh := NewUserHandler(userSvc)
	mux.HandleFunc("GET /api/v1/users", uh.GetAll)
	mux.HandleFunc("GET /api/v1/users/{id}", uh.GetByID)
	mux.HandleFunc("POST /api/v1/users", uh.Create)
	mux.HandleFunc("PUT /api/v1/users/{id}", uh.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", uh.Delete)
}
