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

type Vehicle interface {
	Startup()
	Shutdown()
	HandleCreateVehicle(w http.ResponseWriter, r *http.Request)
	HandleGetVehicleByID(w http.ResponseWriter, r *http.Request)
	HandleGetVehicleByFilter(w http.ResponseWriter, r *http.Request)
	HandleUpdateVehicle(w http.ResponseWriter, r *http.Request)
	HandleDeleteVehicle(w http.ResponseWriter, r *http.Request)
	HandleCreateVehicleValue(w http.ResponseWriter, r *http.Request)
	HandleGetVehicleValueByID(w http.ResponseWriter, r *http.Request)
	HandleGetVehicleValueByFilter(w http.ResponseWriter, r *http.Request)
	HandleUpdateVehicleValue(w http.ResponseWriter, r *http.Request)
	HandleDeleteVehicleValue(w http.ResponseWriter, r *http.Request)
}

// VehicleImpl is the handler implementation for Vehicles
type VehicleImpl struct {
	Service service.Vehicle `inject:"vehicleService"`
}

// Startup performs startup functions
func (h *VehicleImpl) Startup() {
	logger.Trace("Vehicle Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *VehicleImpl) Shutdown() {
	logger.Trace("Vehicle Handler shutting down...")
}

// HandleCreateVehicle handles the request
func (h *VehicleImpl) HandleCreateVehicle(w http.ResponseWriter, r *http.Request) {
	input, err := h.getInputFromRequest(w, r)
	if err != nil {
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	vehicle, err := h.Service.Create(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, vehicle.ToOutput())
}

// HandleGetVehicleByID handles the request
func (h *VehicleImpl) HandleGetVehicleByID(w http.ResponseWriter, r *http.Request) {
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

	vehicle, err := h.Service.GetByID(id, withValues, valueStartDate, valueEndDate, pageSize)

	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, vehicle.ToOutput())
}

// HandleGetVehicleByFilter handles the request
func (h *VehicleImpl) HandleGetVehicleByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.VehicleFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	vehicles, pageInfo, err := h.Service.GetByFilter(input)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	outputs := make([]model.VehicleOutput, 0)
	for _, vehicle := range vehicles {
		output := vehicle.ToOutput()
		outputs = append(outputs, output)
	}

	pageOutput := model.PageOutput{
		Items:    outputs,
		PageInfo: pageInfo,
	}

	response.RespondWithJSON(w, http.StatusOK, pageOutput)
}

// HandleUpdateVehicle handles the request
func (h *VehicleImpl) HandleUpdateVehicle(w http.ResponseWriter, r *http.Request) {
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
	vehicle, err := h.Service.Update(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, vehicle.ToOutput())
}

// HandleDeleteVehicle handles the request
func (h *VehicleImpl) HandleDeleteVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	vehicle, err := h.Service.Delete(id, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, vehicle.ToOutput())
}

// HandleCreateVehicleValue handles the request
func (h *VehicleImpl) HandleCreateVehicleValue(w http.ResponseWriter, r *http.Request) {
	input, err := h.getValueInputFromRequest(w, r)
	if err != nil {
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	vehicleValue, err := h.Service.CreateValue(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, vehicleValue.ToOutput())
}

// HandleGetVehicleValueByID handles the request
func (h *VehicleImpl) HandleGetVehicleValueByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	vehicleValue, err := h.Service.GetValueByID(id)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, vehicleValue.ToOutput())
}

// HandleGetVehicleValueByFilter handles the request
func (h *VehicleImpl) HandleGetVehicleValueByFilter(w http.ResponseWriter, r *http.Request) {
	var input model.VehicleValueFilterInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	vehicleValues, pageInfo, err := h.Service.GetValuesByFilter(input)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	outputs := make([]model.VehicleValueOutput, 0)
	for _, vehicleValue := range vehicleValues {
		output := vehicleValue.ToOutput()
		outputs = append(outputs, output)
	}

	pageOutput := model.PageOutput{
		Items:    outputs,
		PageInfo: pageInfo,
	}

	response.RespondWithJSON(w, http.StatusOK, pageOutput)
}

// HandleUpdateVehicleValue handles the request
func (h *VehicleImpl) HandleUpdateVehicleValue(w http.ResponseWriter, r *http.Request) {
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
	vehicleValue, err := h.Service.UpdateValue(input, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, vehicleValue.ToOutput())
}

// HandleDeleteVehicleValue handles the request
func (h *VehicleImpl) HandleDeleteVehicleValue(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	userID := (r.Context().Value(ctxprops.PropUserID)).(*uuid.UUID)
	vehicleValue, err := h.Service.DeleteValue(id, *userID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, vehicleValue.ToOutput())
}

func (h *VehicleImpl) getInputFromRequest(w http.ResponseWriter, r *http.Request) (input model.VehicleInput, err error) {
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}

	return
}

func (h *VehicleImpl) getValueInputFromRequest(w http.ResponseWriter, r *http.Request) (input model.VehicleValueInput, err error) {
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}

	return
}
