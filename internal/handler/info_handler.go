package handler

import (
	"net/http"

	"github.com/jakabrajadenta/go-explore-api/pkg/response"
)

type endpointInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type apiInfo struct {
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Description string         `json:"description"`
	Language    string         `json:"language"`
	Endpoints   []endpointInfo `json:"endpoints"`
}

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	info := apiInfo{
		Name:        "go-explore-api",
		Version:     "1.0.0",
		Description: "Learning project — REST API with native Go, no framework.",
		Language:    "Go 1.24",
		Endpoints: []endpointInfo{
			{Method: "GET", Path: "/", Description: "API info"},
			{Method: "GET", Path: "/health", Description: "Health check & uptime"},
			{Method: "GET", Path: "/echo", Description: "Echo query params & headers"},
			{Method: "POST", Path: "/echo", Description: "Echo JSON body"},
			{Method: "GET", Path: "/echo/{message}", Description: "Echo path parameter"},
			{Method: "GET", Path: "/api/v1/users", Description: "List users (page, per_page)"},
			{Method: "GET", Path: "/api/v1/users/{id}", Description: "Get user by ID"},
			{Method: "POST", Path: "/api/v1/users", Description: "Create user"},
			{Method: "PUT", Path: "/api/v1/users/{id}", Description: "Update user"},
			{Method: "DELETE", Path: "/api/v1/users/{id}", Description: "Delete user"},
		},
	}
	response.OK(w, r, "go-explore-api is running", info)
}
