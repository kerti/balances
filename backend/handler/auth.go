package handler

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

// Auth handles all requests related to authentication and authorization
type Auth struct {
	Service     *service.Auth `inject:"authService"`
	UserService *service.User `inject:"userService"`
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
	response.RespondWithMessage(w, http.StatusOK, "OK")
}

// HandleAuthLogin performs a login action and returns the JWT token
func (h *Auth) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	authInfo, err := h.Service.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		response.RespondWithError(w, failure.Unauthorized(err.Error()))
		return
	}

	responseObject := response.LoginResponse{
		Expiration: cachetime.CacheTime(*authInfo.Expiration),
		Token:      *authInfo.Token,
		User:       authInfo.User.ToOutput(),
	}
	response.RespondWithJSON(w, http.StatusOK, responseObject)
}

// HandleGetToken fetches a new token for the currently logged-in user
func (h *Auth) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)

	user, err := h.UserService.GetByID(*userID)
	if err != nil {
		response.RespondWithError(w, failure.Unauthorized("user not found"))
		return
	}

	token, expiration, err := h.Service.GetToken(*user)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	responseObject := response.TokenResponse{
		Expiration: cachetime.CacheTime(*expiration),
		Token:      *token,
	}
	response.RespondWithJSON(w, http.StatusOK, responseObject)
}
