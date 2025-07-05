package service

import (
	"math"

	"github.com/google/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

// VehicleImpl is the service provider implementation
type VehicleImpl struct {
	Repository repository.Vehicle `inject:"vehicleRepository"`
}

// Startup performs startup functions
func (s *VehicleImpl) Startup() {
	logger.Trace("Vehicle Service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *VehicleImpl) Shutdown() {
	logger.Trace("Vehicle Service shutting down...")
}

// Create creates a new Vehicle
func (s *VehicleImpl) Create(input model.VehicleInput, userID uuid.UUID) (*model.Vehicle, error) {
	vehicle := model.NewVehicleFromInput(input, userID)
	err := s.Repository.Create(vehicle)
	if err != nil {
		return nil, err
	}
	return &vehicle, err
}

// GetByID fetches a Vehicle by its ID
func (s *VehicleImpl) GetByID(id uuid.UUID, withValues bool, valueStartDate, valueEndDate cachetime.NCacheTime, pageSize *int) (*model.Vehicle, error) {
	vehicles, err := s.Repository.ResolveByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(vehicles) != 1 {
		return nil, failure.EntityNotFound("get by ID", "Vehicle")
	}

	vehicle := vehicles[0]

	if withValues {
		filter := model.VehicleValueFilterInput{
			VehicleIDs: &[]uuid.UUID{id},
		}

		if valueStartDate.Valid {
			filter.StartDate = valueStartDate
		}

		if valueEndDate.Valid {
			filter.EndDate = valueEndDate
		}

		if pageSize != nil {
			filter.PageSize = pageSize
		}

		values, _, err := s.Repository.ResolveValuesByFilter(filter.ToFilter())
		if err != nil {
			return nil, err
		}

		vehicle.AttachValues(values, true)
	}

	return &vehicle, nil
}

// GetByFilter fetches a set of Vehicles by its filter
func (s *VehicleImpl) GetByFilter(input model.VehicleFilterInput) ([]model.Vehicle, model.PageInfoOutput, error) {
	return s.Repository.ResolveByFilter(input.ToFilter())
}

// Update updates an existing Vehicle
func (s *VehicleImpl) Update(input model.VehicleInput, userID uuid.UUID) (*model.Vehicle, error) {
	vehicles, err := s.Repository.ResolveByIDs([]uuid.UUID{input.ID})
	if err != nil {
		return nil, err
	}

	if len(vehicles) != 1 {
		return nil, failure.EntityNotFound("update", "Vehicle")
	}

	vehicle := vehicles[0]

	err = vehicle.Update(input, userID)
	if err != nil {
		return nil, err
	}

	err = s.Repository.Update(vehicle)
	if err != nil {
		return nil, err
	}

	return &vehicle, err
}

// Delete deletes an existing Vehicle. The method will find all the vehicle's values
// and delete all of them also.
func (s *VehicleImpl) Delete(id uuid.UUID, userID uuid.UUID) (*model.Vehicle, error) {
	vehicles, err := s.Repository.ResolveByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(vehicles) != 1 {
		return nil, failure.EntityNotFound("delete", "Vehicle")
	}

	vehicle := vehicles[0]

	// pre-validate to save one database call
	if !vehicle.Deleted.Valid && !vehicle.DeletedBy.Valid {
		filter := model.VehicleValueFilterInput{}
		filter.VehicleIDs = &[]uuid.UUID{vehicle.ID}

		page := 1
		pageSize := math.MaxInt

		filter.Page = &page
		filter.PageSize = &pageSize

		values, _, err := s.Repository.ResolveValuesByFilter(filter.ToFilter())
		if err != nil {
			return nil, err
		}

		vehicle.AttachValues(values, true)
	}

	err = vehicle.Delete(userID)
	if err != nil {
		return nil, err
	}

	err = s.Repository.Update(vehicle)
	if err != nil {
		return nil, err
	}

	return &vehicle, err
}

// CreateValue creates a new Vehicle Value
func (s *VehicleImpl) CreateValue(input model.VehicleValueInput, userID uuid.UUID) (*model.VehicleValue, error) {
	vehicles, err := s.Repository.ResolveByIDs([]uuid.UUID{input.VehicleID})
	if err != nil {
		return nil, err
	}

	if len(vehicles) != 1 {
		return nil, failure.EntityNotFound("create balance", "Vehicle Value")
	}

	vehicle := vehicles[0]

	if vehicle.Deleted.Valid || vehicle.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("add balance", "Vehicle", "the Vehicle is already deleted")
	}

	if vehicle.Status == model.VehicleStatusSold {
		return nil, failure.OperationNotPermitted("add balance", "Vehicle", "the Vehicle has been sold")
	}

	lastValues, err := s.Repository.ResolveLastValuesByVehicleID(vehicle.ID, 1)
	if err != nil {
		return nil, err
	}

	if len(lastValues) != 1 {
		return nil, failure.EntityNotFound("create value", "Vehicle Last Value")
	}

	lastValue := lastValues[0]
	isNewerValue := lastValue.Date.Before(input.Date.Time())
	var vehicleToUpdate *model.Vehicle

	if isNewerValue {
		vehicle.SetCurrentValue(input, userID)
		vehicleToUpdate = &vehicle
	}

	vehicleValue := model.NewVehicleValueFromInput(input, vehicle.ID, userID)
	err = s.Repository.CreateValue(vehicleValue, vehicleToUpdate)
	if err != nil {
		return nil, err
	}

	return &vehicleValue, nil
}

// GetValueByID fetches a Vehicle Value by its ID
func (s *VehicleImpl) GetValueByID(id uuid.UUID) (*model.VehicleValue, error) {
	values, err := s.Repository.ResolveValuesByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(values) != 1 {
		return nil, failure.EntityNotFound("get by ID", "Vehicle Value")
	}

	return &values[0], nil
}

// GetValuesByFilter fetches a set of Vehicle Values by its filter
func (s *VehicleImpl) GetValuesByFilter(input model.VehicleValueFilterInput) ([]model.VehicleValue, model.PageInfoOutput, error) {
	return s.Repository.ResolveValuesByFilter(input.ToFilter())
}

// UpdateValue updates an existing Vehicle Value
func (s *VehicleImpl) UpdateValue(input model.VehicleValueInput, userID uuid.UUID) (*model.VehicleValue, error) {
	vehicles, err := s.Repository.ResolveByIDs([]uuid.UUID{input.VehicleID})
	if err != nil {
		return nil, err
	}

	if len(vehicles) != 1 {
		return nil, failure.EntityNotFound("update", "Vehicle Value")
	}

	vehicle := vehicles[0]

	if vehicle.Deleted.Valid || vehicle.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("update", "Vehicle Value", "the Vehicle is already deleted")
	}

	if vehicle.Status == model.VehicleStatusSold {
		return nil, failure.OperationNotPermitted("update", "Vehicle Value", "the Vehicle has been sold")
	}

	vehicleValues, err := s.Repository.ResolveValuesByIDs([]uuid.UUID{input.ID})
	if err != nil {
		return nil, err
	}

	if len(vehicleValues) != 1 {
		return nil, failure.EntityNotFound("update", "Vehicle Value")
	}

	vehicleValue := vehicleValues[0]

	if vehicleValue.Deleted.Valid || vehicleValue.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("update", "Vehicle Value", "the Vehicle Value is already deleted")
	}

	err = vehicleValue.Update(input, userID)
	if err != nil {
		return nil, err
	}

	currentValues, err := s.Repository.ResolveLastValuesByVehicleID(vehicle.ID, 1)
	if err != nil {
		return nil, err
	}

	if len(currentValues) != 1 {
		return nil, failure.EntityNotFound("update", "Vehicle Value")
	}

	currentValue := currentValues[0]
	isNewerOrCurrentValue := currentValue.Date.Before(input.Date.Time()) || input.ID == currentValue.ID
	var vehicleToUpdate *model.Vehicle

	if isNewerOrCurrentValue {
		vehicle.SetCurrentValue(input, userID)
		vehicleToUpdate = &vehicle
	}

	err = s.Repository.UpdateValue(vehicleValue, vehicleToUpdate)
	if err != nil {
		return nil, err
	}

	return &vehicleValue, nil
}

// DeleteValue deletes an existing Vehicle Value
func (s *VehicleImpl) DeleteValue(id uuid.UUID, userID uuid.UUID) (*model.VehicleValue, error) {
	vehicleValues, err := s.Repository.ResolveValuesByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(vehicleValues) != 1 {
		return nil, failure.EntityNotFound("delete", "Vehicle Value")
	}

	vehicleValue := vehicleValues[0]

	if vehicleValue.Deleted.Valid || vehicleValue.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("delete", "Vehicle Value", "the Vehicle Value is already deleted")
	}

	vehicles, err := s.Repository.ResolveByIDs([]uuid.UUID{vehicleValue.VehicleID})
	if err != nil {
		return nil, err
	}

	if len(vehicles) != 1 {
		return nil, failure.EntityNotFound("delete", "Vehicle")
	}

	vehicle := vehicles[0]

	if vehicle.Deleted.Valid || vehicle.DeletedBy.Valid {
		return nil, failure.OperationNotPermitted("delete", "Vehicle Value", "the Vehicle is already deleted")
	}

	if vehicle.Status == model.VehicleStatusSold {
		return nil, failure.OperationNotPermitted("delete", "Vehicle Value", "the Vehicle has been sold")
	}

	vehicleValue.Delete(userID)

	currentValues, err := s.Repository.ResolveLastValuesByVehicleID(vehicle.ID, 2)
	if err != nil {
		return nil, err
	}

	if len(currentValues) < 1 {
		return nil, failure.EntityNotFound("delete", "Vehicle Current Value")
	}

	if len(currentValues) < 2 {
		return nil, failure.OperationNotPermitted("delete", "Vehicle Value", "cannot delete the only Vehicle Value belonging to a Vehicle")
	}

	currentValueDeleted := vehicleValue.ID.String() == currentValues[0].ID.String()
	var vehicleToUpdate *model.Vehicle

	if currentValueDeleted {
		newCurrentValueInput := model.VehicleValueInput{
			ID:        currentValues[1].ID,
			VehicleID: currentValues[1].VehicleID,
			Value:     currentValues[1].Value,
			Date:      cachetime.CacheTime(currentValues[1].Date),
		}
		vehicle.SetCurrentValue(newCurrentValueInput, userID)
		vehicleToUpdate = &vehicle
	}

	err = s.Repository.UpdateValue(vehicleValue, vehicleToUpdate)
	if err != nil {
		return nil, err
	}

	return &vehicleValue, nil
}
