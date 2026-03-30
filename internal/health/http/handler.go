//revive:disable:package-comments
package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Handler serves HTTP health checks.
type Handler struct {
	log *slog.Logger
}

// New returns a health check handler.
func New(log *slog.Logger) *Handler {
	return &Handler{log: log}
}

type status struct {
	Status string `json:"status"`
}

// ServeHTTP reports service health.
func (h *Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(status{Status: "ok"}); err != nil {
		h.log.Error("failed to encode health response", "error", err)
	}
}
