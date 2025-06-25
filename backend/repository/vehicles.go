package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/balances/backend/database"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/filter"
	"github.com/kerti/balances/backend/util/logger"
)

const (
	QueryInsertVehicle = `
		INSERT INTO vehicles (
			entity_id,
			name,
			make,
			model,
			year,
			type,
			title_holder,
			license_plate_number,
			purchase_date,
			initial_value,
			initial_value_date,
			current_value,
			current_value_date,
			annual_depreciation_percent,
			status,
			created,
			created_by,
			updated,
			updated_by,
			deleted,
			deleted_by
		) VALUES (
			:entity_id,
			:name,
			:make,
			:model,
			:year,
			:type,
			:title_holder,
			:license_plate_number,
			:purchase_date,
			:initial_value,
			:initial_value_date,
			:current_value,
			:current_value_date,
			:annual_depreciation_percent,
			:status,
			:created,
			:created_by,
			:updated,
			:updated_by,
			:deleted,
			:deleted_by
		)`

	QueryInsertVehicleValue = `
		INSERT INTO vehicle_values(
			entity_id,
			vehicle_entity_id,
			date,
			value,
			created,
			created_by,
			updated,
			updated_by,
			deleted,
			deleted_by
		) VALUES (
			:entity_id,
			:vehicle_entity_id,
			:date,
			:value,
			:created,
			:created_by,
			:updated,
			:updated_by,
			:deleted,
			:deleted_by
		)`
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
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
	}
	return
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
	exists, err := r.ExistsByID(vehicle.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if exists {
		err = failure.OperationNotPermitted("create", "Vehicle", "already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreateVehicle(tx, vehicle); err != nil {
			e <- err
			return
		}

		for _, value := range vehicle.Values {
			if err := r.txCreateVehicleValue(tx, value); err != nil {
				e <- err
				return
			}
		}

		e <- nil
	})
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

func (r *VehicleMySQLRepo) txCreateVehicle(tx *sqlx.Tx, vehicle model.Vehicle) error {
	stmt, err := tx.PrepareNamed(QueryInsertVehicle)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(vehicle)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}

func (r *VehicleMySQLRepo) txCreateVehicleValue(tx *sqlx.Tx, vehicleValue model.VehicleValue) error {
	stmt, err := tx.PrepareNamed(QueryInsertVehicleValue)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(vehicleValue)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}
