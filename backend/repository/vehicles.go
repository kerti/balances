package repository

import (
	"github.com/google/uuid"
	"github.com/kerti/balances/backend/database"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/filter"
	"github.com/kerti/balances/backend/util/logger"
)

// VehicleMySQLRepo is the repository for Vehicles implemented with MySQL backend
type VehicleMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup perform startup functions
func (r *VehicleMySQLRepo) Startup() {
	logger.Trace("Vehicle repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *VehicleMySQLRepo) Shutdown() {
	logger.Trace("Vehicle repository shutting down...")
}

// ExistsByID checks the existence of a Vehicle by its ID
func (r *VehicleMySQLRepo) ExistsByID(id uuid.UUID) (exists bool, err error) {
	return false, failure.Unimplemented("repository unimplemented for this method")
}

// ExistsValueByID checks the existence of a Vehicle Value by its ID
func (r *VehicleMySQLRepo) ExistsValueByID(id uuid.UUID) (exists bool, err error) {
	return false, failure.Unimplemented("repository unimplemented for this method")
}

// ResolveByIDs resolves Vehicles by their IDs
func (r *VehicleMySQLRepo) ResolveByIDs(ids []uuid.UUID) (vehicles []model.Vehicle, err error) {
	return []model.Vehicle{}, failure.Unimplemented("repository unimplemented for this method")
}

// ResolveValuesByIDs resolves Vehicle Values by their IDs
func (r *VehicleMySQLRepo) ResolveValuesByIDs(ids []uuid.UUID) (vehicleValues []model.VehicleValue, err error) {
	return []model.VehicleValue{}, failure.Unimplemented("repository unimplemented for this method")
}

// ResolveByFilter resolves Vehicles by a specified filter
func (r *VehicleMySQLRepo) ResolveByFilter(filter filter.Filter) (vehicles []model.Vehicle, pageInfo model.PageInfoOutput, err error) {
	return []model.Vehicle{}, model.PageInfoOutput{}, failure.Unimplemented("repository unimplemented for this method")
}

// ResolveValuesByFilter resolves Vehicle Values by a specified filter
func (r *VehicleMySQLRepo) ResolveValuesByFilter(filter filter.Filter) (vehicleValues []model.VehicleValue, pageInfo model.PageInfoOutput, err error) {
	return []model.VehicleValue{}, model.PageInfoOutput{}, failure.Unimplemented("repository unimplemented for this method")
}

// ResolveLastValuesByVehicleID resolves last X Vehicle Values by their Vehicle ID and count param
func (r *VehicleMySQLRepo) ResolveLastValuesByVehicleID(id uuid.UUID, count int) (vehicleValues []model.VehicleValue, err error) {
	return []model.VehicleValue{}, failure.Unimplemented("repository unimplemented for this method")
}

// Create creates a Vehicle
func (r *VehicleMySQLRepo) Create(vehicle model.Vehicle) error {
	return failure.Unimplemented("repository unimplemented for this method")
}

// Update updates a Vehicle
func (r *VehicleMySQLRepo) Update(vehicle model.Vehicle) error {
	return failure.Unimplemented("repository unimplemented for this method")
}

// CreateValue creates a new Vehicle Value and optionally updates the Vehicle transactionally
func (r *VehicleMySQLRepo) CreateValue(vehicleValue model.VehicleValue, vehicle *model.Vehicle) error {
	return failure.Unimplemented("repository unimplemented for this method")
}

// UpdateValue updates an existing Vehicle Value and optionally updates the Vehicle transactionally
func (r *VehicleMySQLRepo) UpdateValue(vehicleValue model.VehicleValue, vehicle *model.Vehicle) error {
	return failure.Unimplemented("repository unimplemented for this method")
}
