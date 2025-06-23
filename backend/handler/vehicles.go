package handler

import (
	"errors"
	"net/http"

	"github.com/kerti/balances/backend/handler/response"
	"github.com/kerti/balances/backend/service"
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

type VehicleImpl struct {
	Service service.Vehicle `inject:"vehicleService"`
}

func (h *VehicleImpl) Startup() {
	logger.Trace("Vehicle Handler starting up...")
}

func (h *VehicleImpl) Shutdown() {
	logger.Trace("Vehicle Handler shutting down...")
}

// HandleCreateVehicle handles the request
func (h *VehicleImpl) HandleCreateVehicle(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleGetVehicleByID handles the request
func (h *VehicleImpl) HandleGetVehicleByID(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleGetVehicleByFilter handles the request
func (h *VehicleImpl) HandleGetVehicleByFilter(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleUpdateVehicle handles the request
func (h *VehicleImpl) HandleUpdateVehicle(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleDeleteVehicle handles the request
func (h *VehicleImpl) HandleDeleteVehicle(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleCreateVehicleValue handles the request
func (h *VehicleImpl) HandleCreateVehicleValue(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleGetVehicleValueByID handles the request
func (h *VehicleImpl) HandleGetVehicleValueByID(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleGetVehicleValueByFilter handles the request
func (h *VehicleImpl) HandleGetVehicleValueByFilter(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleUpdateVehicleValue handles the request
func (h *VehicleImpl) HandleUpdateVehicleValue(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}

// HandleDeleteVehicleValue handles the request
func (h *VehicleImpl) HandleDeleteVehicleValue(w http.ResponseWriter, r *http.Request) {
	response.RespondWithError(w, errors.New("handler unimplemented for this method"))
}
