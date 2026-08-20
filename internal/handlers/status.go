package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received health check request", slog.String("method", r.Method), slog.String("url", r.URL.String()))
	render.JSON(w, r, HealthResponse{Status: "OK"})
}
