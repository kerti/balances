package handler

import (
	"net/http"

	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/util/logger"
)

// Health handles all requests related to server health check
type Health struct {
	isPreparingShutdown bool
	isHealthy           bool
}

// Startup perform startup functions
func (h *Health) Startup() {
	logger.Trace("Health Handler starting up...")
	h.isPreparingShutdown = false
	h.isHealthy = true
}

// PrepareShutdown prepares the service for shutdown
func (h *Health) PrepareShutdown() {
	logger.Trace("Health Handler preparing for shutdown...")
	h.isPreparingShutdown = true
	h.isHealthy = false
}

// Shutdown cleans up everything and shuts down
func (h *Health) Shutdown() {
	h.isPreparingShutdown = false
	logger.Trace("Health Handler shutting down...")
}

// HandleHealthCheck handles the request
func (h *Health) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if h.isHealthy {
		response.RespondWithMessage(w, http.StatusOK, "OK")
	} else {
		if h.isPreparingShutdown {
			response.RespondWithPreparingShutdown(w)
		} else {
			response.RespondWithUnhealthy(w)
		}
	}
}
