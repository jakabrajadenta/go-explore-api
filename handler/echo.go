package handler

import (
	"encoding/json"
	"net/http"
)

type EchoResponse struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers,omitempty"`
	Query   map[string]string `json:"query,omitempty"`
	Body    any               `json:"body,omitempty"`
	Message string            `json:"message,omitempty"`
}

// HandleEchoGet echoes back the request's query parameters and headers.
func HandleEchoGet(w http.ResponseWriter, r *http.Request) {
	resp := EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: flattenHeaders(r),
		Query:   flattenQuery(r),
	}
	writeJSON(w, http.StatusOK, resp)
}

// HandleEchoPost decodes a JSON body and echoes it back.
func HandleEchoPost(w http.ResponseWriter, r *http.Request) {
	var body any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request body must be valid JSON"})
		return
	}

	resp := EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: flattenHeaders(r),
		Body:    body,
	}
	writeJSON(w, http.StatusOK, resp)
}

// HandleEchoPath echoes back a message from the URL path parameter.
func HandleEchoPath(w http.ResponseWriter, r *http.Request) {
	resp := EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Message: r.PathValue("message"),
	}
	writeJSON(w, http.StatusOK, resp)
}

func flattenHeaders(r *http.Request) map[string]string {
	out := make(map[string]string, len(r.Header))
	for key, vals := range r.Header {
		out[key] = vals[0]
	}
	return out
}

func flattenQuery(r *http.Request) map[string]string {
	params := r.URL.Query()
	out := make(map[string]string, len(params))
	for key, vals := range params {
		out[key] = vals[0]
	}
	return out
}
