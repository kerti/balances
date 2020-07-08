package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kerti/balances/backend/model"

	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/logger"
	"github.com/satori/uuid"
)

// User handles all requests related to users
type User struct {
	Service *service.User `inject:""`
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
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	users, err := h.Service.GetByIDs([]uuid.UUID{id})
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(users) != 1 {
		response.RespondWithError(w, http.StatusInternalServerError, "user not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, users[0].ToOutput())
}

// HandleCreateUser handles the request
func (h *User) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input model.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	// TODO: FIX THIS!
	userID := uuid.NewV4()
	user, err := h.Service.Create(input, userID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
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
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var input model.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err.Error())
	}
	input.ID = id

	// TODO: FIX THIS!
	userID := uuid.NewV4()
	user, err := h.Service.Update(input, userID)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, user.ToOutput())
	return
}
