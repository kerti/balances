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
	QuerySelectVehicle = `
		SELECT
			vehicles.entity_id,
			vehicles.name,
			vehicles.make,
			vehicles.model,
			vehicles.year,
			vehicles.type,
			vehicles.title_holder,
			vehicles.license_plate_number,
			vehicles.purchase_date,
			vehicles.initial_value,
			vehicles.initial_value_date,
			vehicles.current_value,
			vehicles.current_value_date,
			vehicles.annual_depreciation_percent,
			vehicles.status,
			vehicles.created,
			vehicles.created_by,
			vehicles.updated,
			vehicles.updated_by,
			vehicles.deleted,
			vehicles.deleted_by
		FROM
			vehicles `

	QuerySelectVehicleValues = `
		SELECT
			vehicle_values.entity_id,
			vehicle_values.vehicle_entity_id,
			vehicle_values.date,
			vehicle_values.value,
			vehicle_values.created,
			vehicle_values.created_by,
			vehicle_values.updated,
			vehicle_values.updated_by,
			vehicle_values.deleted,
			vehicle_values.deleted_by
		FROM
			vehicle_values `

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
		INSERT INTO vehicle_values (
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

	QueryUpdateVehicle = `
		UPDATE vehicles
		SET
			name = :name,
			make = :make,
			model = :model,
			year = :year,
			type = :type,
			title_holder = :title_holder,
			license_plate_number = :license_plate_number,
			purchase_date = :purchase_date,
			initial_value = :initial_value,
			initial_value_date = :initial_value_date,
			current_value = :current_value,
			current_value_date = :current_value_date,
			annual_depreciation_percent = :annual_depreciation_percent,
			status = :status,
			created = :created,
			created_by = :created_by,
			updated = :updated,
			updated_by = :updated_by,
			deleted = :deleted,
			deleted_by = :deleted_by
		WHERE entity_id = :entity_id`
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
		err = failure.InternalError("exists by ID", "Vehicle", err)
	}
	return
}

// ExistsValueByID checks the existence of a Vehicle Value by its ID
func (r *VehicleMySQLRepo) ExistsValueByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM vehicle_values WHERE vehicle_values.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("exists by ID", "Vehicle Value", err)
	}
	return
}

// ResolveByIDs resolves Vehicles by their IDs
func (r *VehicleMySQLRepo) ResolveByIDs(ids []uuid.UUID) (vehicles []model.Vehicle, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(QuerySelectVehicle+" WHERE vehicles.entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by IDs", "Vehicle", err)
		return
	}

	err = r.DB.Select(&vehicles, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by IDs", "Vehicle", err)
	}

	return
}

// ResolveValuesByIDs resolves Vehicle Values by their IDs
func (r *VehicleMySQLRepo) ResolveValuesByIDs(ids []uuid.UUID) (vehicleValues []model.VehicleValue, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(QuerySelectVehicleValues+" WHERE vehicle_values.entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by IDs", "Vehicle Value", err)
		return
	}

	err = r.DB.Select(&vehicleValues, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by IDs", "Vehicle Value", err)
	}

	return
}

// ResolveByFilter resolves Vehicles by a specified filter
func (r *VehicleMySQLRepo) ResolveByFilter(filter filter.Filter) (vehicles []model.Vehicle, pageInfo model.PageInfoOutput, err error) {
	filterQueryString, err := filter.ToQueryString()
	if err != nil {
		err = failure.InternalError("resolve by filter", "Vehicle", err)
		return vehicles, pageInfo, err
	}

	filterArgs := filter.GetArgs(true)
	query, args, err := r.DB.In(
		QuerySelectVehicle+filterQueryString+filter.Pagination.ToQueryString(),
		filterArgs...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Vehicle", err)
		return
	}

	err = r.DB.Select(&vehicles, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Vehicle", err)
		return
	}

	var count int
	filterArgsNoPagination := filter.GetArgs(false)
	err = r.DB.Get(
		&count,
		"SELECT COUNT(entity_id) FROM vehicles "+filterQueryString,
		filterArgsNoPagination...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Vehicle", err)
		vehicles = []model.Vehicle{}
		return
	}

	pageInfo = model.PageInfoOutput{
		Page:       filter.Pagination.Page,
		PageSize:   filter.Pagination.PageSize,
		TotalCount: count,
		PageCount:  filter.Pagination.GetPageCount(count),
	}

	return
}

// ResolveValuesByFilter resolves Vehicle Values by a specified filter
func (r *VehicleMySQLRepo) ResolveValuesByFilter(filter filter.Filter) (vehicleValues []model.VehicleValue, pageInfo model.PageInfoOutput, err error) {
	filterQueryString, err := filter.ToQueryString()
	if err != nil {
		err = failure.InternalError("resolve by filter", "Vehicle Value", err)
		return
	}

	filterArgs := filter.GetArgs(true)
	query, args, err := r.DB.In(
		QuerySelectVehicleValues+filterQueryString+filter.Pagination.ToQueryString(),
		filterArgs...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Vehicle Value", err)
		return
	}

	err = r.DB.Select(&vehicleValues, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Vehicle Value", err)
		return
	}

	var count int
	filterArgsNoPagination := filter.GetArgs(false)
	query, args, err = r.DB.In(
		"SELECT COUNT(entity_id) FROM vehicle_values "+filterQueryString,
		filterArgsNoPagination...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Vehicle Value", err)
		vehicleValues = []model.VehicleValue{}
		return
	}

	err = r.DB.Get(&count, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Vehicle Value", err)
		vehicleValues = []model.VehicleValue{}
		return
	}

	pageInfo = model.PageInfoOutput{
		Page:       filter.Pagination.Page,
		PageSize:   filter.Pagination.PageSize,
		TotalCount: count,
		PageCount:  filter.Pagination.GetPageCount(count),
	}

	return
}

// ResolveLastValuesByVehicleID resolves last X Vehicle Values by their Vehicle ID and count param
func (r *VehicleMySQLRepo) ResolveLastValuesByVehicleID(id uuid.UUID, count int) (vehicleValues []model.VehicleValue, err error) {
	return []model.VehicleValue{}, failure.Unimplemented("repository unimplemented for this method: resolveLastValuesByVehicleID")
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
			wrappedErr := failure.InternalError("create", "Vehicle", err)
			e <- wrappedErr
			return
		}

		for _, value := range vehicle.Values {
			if err := r.txCreateVehicleValue(tx, value); err != nil {
				wrappedErr := failure.InternalError("create", "Vehicle", err)
				e <- wrappedErr
				return
			}
		}

		e <- nil
	})
}

// Update updates a Vehicle
func (r *VehicleMySQLRepo) Update(vehicle model.Vehicle) error {
	exists, err := r.ExistsByID(vehicle.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if !exists {
		err = failure.EntityNotFound("update", "Vehicle")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdateVehicle(tx, vehicle); err != nil {
			err = failure.InternalError("update", "Vehicle", err)
			e <- err
			return
		}

		e <- nil
	})
}

// CreateValue creates a new Vehicle Value and optionally updates the Vehicle transactionally
func (r *VehicleMySQLRepo) CreateValue(vehicleValue model.VehicleValue, vehicle *model.Vehicle) error {
	return failure.Unimplemented("repository unimplemented for this method: createValue")
}

// UpdateValue updates an existing Vehicle Value and optionally updates the Vehicle transactionally
func (r *VehicleMySQLRepo) UpdateValue(vehicleValue model.VehicleValue, vehicle *model.Vehicle) error {
	return failure.Unimplemented("repository unimplemented for this method: updateValue")
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

func (r *VehicleMySQLRepo) txUpdateVehicle(tx *sqlx.Tx, vehicle model.Vehicle) error {
	stmt, err := tx.PrepareNamed(QueryUpdateVehicle)
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
