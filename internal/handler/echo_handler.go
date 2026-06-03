package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jakabrajadenta/go-explore-api/pkg/response"
)

type echoData struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers,omitempty"`
	Query   map[string]string `json:"query,omitempty"`
	Body    any               `json:"body,omitempty"`
	Message string            `json:"message,omitempty"`
}

func HandleEchoGet(w http.ResponseWriter, r *http.Request) {
	data := echoData{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: flattenHeaders(r),
		Query:   flattenQuery(r),
	}
	response.OK(w, r, "Echo GET", data)
}

func HandleEchoPost(w http.ResponseWriter, r *http.Request) {
	var body any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, r, map[string]string{"body": "request body must be valid JSON"})
		return
	}
	data := echoData{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: flattenHeaders(r),
		Body:    body,
	}
	response.OK(w, r, "Echo POST", data)
}

func HandleEchoPath(w http.ResponseWriter, r *http.Request) {
	data := echoData{
		Method:  r.Method,
		Path:    r.URL.Path,
		Message: r.PathValue("message"),
	}
	response.OK(w, r, "Echo path", data)
}

func flattenHeaders(r *http.Request) map[string]string {
	out := make(map[string]string, len(r.Header))
	for k, v := range r.Header {
		out[k] = v[0]
	}
	return out
}

func flattenQuery(r *http.Request) map[string]string {
	params := r.URL.Query()
	out := make(map[string]string, len(params))
	for k, v := range params {
		out[k] = v[0]
	}
	return out
}
