package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
	"github.com/satori/uuid"
)

// User handles all requests related to Users
type User struct {
	Service *service.User `inject:"userService"`
}

// Startup perform startup functions
func (h *User) Startup() {
	logger.Trace("User Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *User) Shutdown() {
	logger.Trace("User Handler shutting down...")
}

// HandleGetUserByID handles the request
func (h *User) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	user, err := h.Service.GetByID(id)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, user.ToOutput())
}

// HandleCreateUser handles the request
func (h *User) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input model.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	user, err := h.Service.Create(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, user.ToOutput())
	return
}

// HandleUpdateUser handles the request
func (h *User) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	var input model.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}
	input.ID = id

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	user, err := h.Service.Update(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, user.ToOutput())
	return
}
