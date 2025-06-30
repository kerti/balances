package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

// User is the handler interface for Users
type User interface {
	Startup()
	Shutdown()
	HandleGetUserByID(w http.ResponseWriter, r *http.Request)
	HandleGetUserByFilter(w http.ResponseWriter, r *http.Request)
	HandleCreateUser(w http.ResponseWriter, r *http.Request)
	HandleUpdateUser(w http.ResponseWriter, r *http.Request)
}

// UserImpl is the handler implementation for Users
type UserImpl struct {
	Service service.User `inject:"userService"`
}

// Startup perform startup functions
func (h *UserImpl) Startup() {
	logger.Trace("User Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *UserImpl) Shutdown() {
	logger.Trace("User Handler shutting down...")
}

// HandleGetUserByID handles the request
func (h *UserImpl) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
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

// HandleGetUserByFilter handles the request
func (h *UserImpl) HandleGetUserByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.UserFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	users, pageInfo, err := h.Service.GetByFilter(input)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	outputs := make([]model.UserOutput, 0)
	for _, user := range users {
		outputs = append(outputs, user.ToOutput())
	}

	pageOutput := model.PageOutput{
		Items:    outputs,
		PageInfo: pageInfo,
	}

	response.RespondWithJSON(w, http.StatusOK, pageOutput)
}

// HandleCreateUser handles the request
func (h *UserImpl) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
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
}

// HandleUpdateUser handles the request
func (h *UserImpl) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
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
}
