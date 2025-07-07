package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

type Property interface {
	Startup()
	Shutdown()
	HandleCreateProperty(w http.ResponseWriter, r *http.Request)
	HandleGetPropertyByID(w http.ResponseWriter, r *http.Request)
	HandleGetPropertyByFilter(w http.ResponseWriter, r *http.Request)
	HandleUpdateProperty(w http.ResponseWriter, r *http.Request)
	HandleDeleteProperty(w http.ResponseWriter, r *http.Request)
	HandleCreatePropertyValue(w http.ResponseWriter, r *http.Request)
	HandleGetPropertyValueByID(w http.ResponseWriter, r *http.Request)
	HandleGetPropertyValueByFilter(w http.ResponseWriter, r *http.Request)
	HandleUpdatePropertyValue(w http.ResponseWriter, r *http.Request)
	HandleDeletePropertyValue(w http.ResponseWriter, r *http.Request)
}

// PropertyImpl is the handler implementation for Properties
type PropertyImpl struct {
	Service service.Property `inject:"propertyService"`
}

// Startup performs startup functions
func (h *PropertyImpl) Startup() {
	logger.Trace("Property Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *PropertyImpl) Shutdown() {
	logger.Trace("Property Handler shutting down...")
}

// HandleCreateProperty handles the request
func (h *PropertyImpl) HandleCreateProperty(w http.ResponseWriter, r *http.Request) {
	input, err := h.getInputFromRequest(w, r)
	if err != nil {
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	property, err := h.Service.Create(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, property.ToOutput())
}

// HandleGetPropertyByID handles the request
func (h *PropertyImpl) HandleGetPropertyByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	r.ParseForm()
	_, withValues := r.Form["withValues"]
	valueStartDateStr, withValueStartDate := r.Form["valueStartDate"]
	valueEndDateStr, withValueEndDate := r.Form["valueEndDate"]
	pageSizeStr, withPageSize := r.Form["pageSize"]

	var valueStartDate cachetime.NCacheTime
	if withValueStartDate {
		valueStartDate.Scan(valueStartDateStr[0])
	}

	var valueEndDate cachetime.NCacheTime
	if withValueEndDate {
		valueEndDate.Scan(valueEndDateStr[0])
	}

	var pageSize *int = nil
	if withPageSize {
		parsedPageSize, err := strconv.Atoi(pageSizeStr[0])
		if err == nil {
			pageSize = &parsedPageSize
		}
	}

	property, err := h.Service.GetByID(id, withValues, valueStartDate, valueEndDate, pageSize)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, property.ToOutput())
}

// HandleGetPropertyByFilter handles the request
func (h *PropertyImpl) HandleGetPropertyByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.PropertyFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	properties, pageInfo, err := h.Service.GetByFilter(input)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	outputs := make([]model.PropertyOutput, 0)
	for _, property := range properties {
		output := property.ToOutput()
		outputs = append(outputs, output)
	}

	pageOutput := model.PageOutput{
		Items:    outputs,
		PageInfo: pageInfo,
	}

	response.RespondWithJSON(w, http.StatusOK, pageOutput)
}

// HandleUpdateProperty handles the request
func (h *PropertyImpl) HandleUpdateProperty(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	input, err := h.getInputFromRequest(w, r)
	if err != nil {
		return
	}

	if input.ID.String() != id.String() {
		response.RespondWithError(w, failure.BadRequestFromString("id mismatch"))
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	property, err := h.Service.Update(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, property.ToOutput())
}

// HandleDeleteProperty handles the request
func (h *PropertyImpl) HandleDeleteProperty(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	property, err := h.Service.Delete(id, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, property.ToOutput())
}

// HandleCreatePropertyValue handles the request
func (h *PropertyImpl) HandleCreatePropertyValue(w http.ResponseWriter, r *http.Request) {
	input, err := h.getValueInputFromRequest(w, r)
	if err != nil {
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	propertyValue, err := h.Service.CreateValue(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, propertyValue.ToOutput())
}

// HandleGetPropertyValueByID handles the request
func (h *PropertyImpl) HandleGetPropertyValueByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	propertyValue, err := h.Service.GetValueByID(id)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, propertyValue.ToOutput())
}

// HandleGetPropertyValueByFilter handles the request
func (h *PropertyImpl) HandleGetPropertyValueByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.PropertyValueFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	propertyValues, pageInfo, err := h.Service.GetValuesByFilter(input)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	outputs := make([]model.PropertyValueOutput, 0)
	for _, propertyValue := range propertyValues {
		output := propertyValue.ToOutput()
		outputs = append(outputs, output)
	}

	pageOutput := model.PageOutput{
		Items:    outputs,
		PageInfo: pageInfo,
	}

	response.RespondWithJSON(w, http.StatusOK, pageOutput)
}

// HandleUpdatePropertyValue handles the request
func (h *PropertyImpl) HandleUpdatePropertyValue(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	input, err := h.getValueInputFromRequest(w, r)
	if err != nil {
		return
	}

	if input.ID.String() != id.String() {
		response.RespondWithError(w, failure.BadRequestFromString("id mismatch"))
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	propertyValue, err := h.Service.UpdateValue(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, propertyValue.ToOutput())
}

// HandleDeletePropertyValue handles the request
func (h *PropertyImpl) HandleDeletePropertyValue(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	propertyValue, err := h.Service.DeleteValue(id, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, propertyValue.ToOutput())
}

func (h *PropertyImpl) getInputFromRequest(w http.ResponseWriter, r *http.Request) (input model.PropertyInput, err error) {
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}

	return
}

func (h *PropertyImpl) getValueInputFromRequest(w http.ResponseWriter, r *http.Request) (input model.PropertyValueInput, err error) {
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}

	return
}
