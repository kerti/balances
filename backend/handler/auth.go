package handler

import (
	"net/http"

	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/satori/uuid"

	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/logger"
)

// Auth handles all requests related to authentication and authorization
type Auth struct {
	Service     *service.Auth `inject:""`
	UserService *service.User `inject:""`
}

// Startup perform startup functions
func (h *Auth) Startup() {
	logger.Trace("Auth Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *Auth) Shutdown() {
	logger.Trace("Auth Handler shutting down...")
}

// HandlePreflight handles a preflight check for logins
func (h *Auth) HandlePreflight(w http.ResponseWriter, r *http.Request) {
	respondWithMessage(w, http.StatusOK, "OK")
}

// HandleAuthLogin performs a login action and returns the JWT token
func (h *Auth) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	token, err := h.Service.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithMessage(w, http.StatusOK, *token)
}

// HandleGetToken fetches a new token for the currently logged-in user
func (h *Auth) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxprops.PropUserID)
	if userID == nil {
		respondWithError(w, http.StatusUnauthorized, "user is not logged in")
		return
	}

	uuidSlice := []uuid.UUID{
		*(userID.(*uuid.UUID)),
	}
	users, err := h.UserService.GetByIDs(uuidSlice)
	if err != nil || len(users) != 1 {
		respondWithError(w, http.StatusUnauthorized, "user not found")
		return
	}

	token, err := h.Service.GetToken(users[0])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithMessage(w, http.StatusOK, *token)
}
