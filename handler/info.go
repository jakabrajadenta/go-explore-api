package handler

import "net/http"

type InfoResponse struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Language    string     `json:"language"`
	Endpoints   []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	resp := InfoResponse{
		Name:        "go-explore-api",
		Version:     "1.0.0",
		Description: "A learning project for building REST APIs with native Go — no framework, just stdlib.",
		Language:    "Go 1.24",
		Endpoints: []Endpoint{
			{Method: "GET", Path: "/", Description: "API info and available endpoints"},
			{Method: "GET", Path: "/health", Description: "Health check with uptime"},
			{Method: "GET", Path: "/echo", Description: "Echo query params and request headers"},
			{Method: "POST", Path: "/echo", Description: "Echo JSON request body back to caller"},
			{Method: "GET", Path: "/echo/{message}", Description: "Echo a path parameter as message"},
		},
	}
	writeJSON(w, http.StatusOK, resp)
}
