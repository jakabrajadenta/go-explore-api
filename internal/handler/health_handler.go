package handler

import (
	"net/http"
	"time"

	"github.com/jakabrajadenta/go-explore-api/pkg/response"
)

var startTime = time.Now()

type healthData struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	data := healthData{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(startTime).Round(time.Second).String(),
	}
	response.OK(w, r, "Service is healthy", data)
}
