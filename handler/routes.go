package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", HandleInfo)
	mux.HandleFunc("GET /health", HandleHealth)
	mux.HandleFunc("GET /echo", HandleEchoGet)
	mux.HandleFunc("POST /echo", HandleEchoPost)
	mux.HandleFunc("GET /echo/{message}", HandleEchoPath)
}
