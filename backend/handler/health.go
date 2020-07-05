package handler

import (
	"net/http"

	"github.com/kerti/balances/backend/util/logger"
)

// Health handles all requests related to server health check
type Health struct {
	isHealthy bool
}

// Startup perform startup functions
func (h *Health) Startup() {
	logger.Trace("Health Handler starting up...")
	h.isHealthy = true
}

// PrepareShutdown prepares the service for shutdown
func (h *Health) PrepareShutdown() {
	logger.Trace("Health Handler preparing for shutdown...")
	h.isHealthy = false
}

// Shutdown cleans up everything and shuts down
func (h *Health) Shutdown() {
	logger.Trace("Health Handler shutting down...")
}

// HandleHealthCheck handles the request
func (h *Health) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if h.isHealthy {
		respondWithMessage(w, http.StatusOK, "OK")
	} else {
		respondWithMessage(w, http.StatusServiceUnavailable, "SERVER UNHEALTHY")
	}
}
