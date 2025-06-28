package service

import (
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
	return nil, failure.Unimplemented("service unimplemented for this method")
}

// Delete deletes an existing Vehicle. The method will find all the vehicle's values
// and delete all of them also.
func (s *VehicleImpl) Delete(id uuid.UUID, userID uuid.UUID) (*model.Vehicle, error) {
	return nil, failure.Unimplemented("service unimplemented for this method")
}

// CreateValue creates a new Vehicle Value
func (s *VehicleImpl) CreateValue(input model.VehicleValueInput, userID uuid.UUID) (*model.VehicleValue, error) {
	return nil, failure.Unimplemented("service unimplemented for this method")
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
	return nil, failure.Unimplemented("service unimplemented for this method")
}

// DeleteValue deletes an existing Vehicle Value
func (s *VehicleImpl) DeleteValue(id uuid.UUID, userID uuid.UUID) (*model.VehicleValue, error) {
	return nil, failure.Unimplemented("service unimplemented for this method")
}
