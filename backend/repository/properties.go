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
	QuerySelectProperty = `
		SELECT
			properties.entity_id,
			properties.name,
			properties.address,
			properties.total_area,
			properties.building_area,
			properties.area_unit,
			properties.type,
			properties.title_holder,
			properties.tax_identifier,
			properties.purchase_date,
			properties.initial_value,
			properties.initial_value_date,
			properties.current_value,
			properties.current_value_date,
			properties.annual_appreciation_percent,
			properties.status,
			properties.created,
			properties.created_by,
			properties.updated,
			properties.updated_by,
			properties.deleted,
			properties.deleted_by
		FROM
			properties `

	QuerySelectPropertyValues = `
		SELECT
			property_values.entity_id,
			property_values.property_entity_id,
			property_values.date,
			property_values.value,
			property_values.created,
			property_values.created_by,
			property_values.updated,
			property_values.updated_by,
			property_values.deleted,
			property_values.deleted_by
		FROM
			property_values `

	QueryInsertProperty = `
		INSERT INTO properties (
			entity_id,
			name,
			address,
			total_area,
			building_area,
			area_unit,
			type,
			title_holder,
			tax_identifier,
			purchase_date,
			initial_value,
			initial_value_date,
			current_value,
			current_value_date,
			annual_appreciation_percent,
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
			:address,
			:total_area,
			:building_area,
			:area_unit,
			:type,
			:title_holder,
			:tax_identifier,
			:purchase_date,
			:initial_value,
			:initial_value_date,
			:current_value,
			:current_value_date,
			:annual_appreciation_percent,
			:status,
			:created,
			:created_by,
			:updated,
			:updated_by,
			:deleted,
			:deleted_by
		)`

	QueryInsertPropertyValue = `
		INSERT INTO property_values (
			entity_id,
			property_entity_id,
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
			:property_entity_id,
			:date,
			:value,
			:created,
			:created_by,
			:updated,
			:updated_by,
			:deleted,
			:deleted_by
		)`

	QueryUpdateProperty = `
		UPDATE properties
		SET
			name = :name,
			address = :address,
			total_area = :total_area,
			building_area = :building_area,
			area_unit = :area_unit,
			type = :type,
			title_holder = :title_holder,
			tax_identifier = :tax_identifier,
			purchase_date = :purchase_date,
			initial_value = :initial_value,
			initial_value_date = :initial_value_date,
			current_value = :current_value,
			current_value_date = :current_value_date,
			annual_appreciation_percent = :annual_appreciation_percent,
			status = :status,
			created = :created,
			created_by = :created_by,
			updated = :updated,
			updated_by = :updated_by,
			deleted = :deleted,
			deleted_by = :deleted_by
		WHERE entity_id = :entity_id`

	QueryUpdatePropertyValue = `
		UPDATE property_values
		SET
			property_entity_id = :property_entity_id,
			date = :date,
			value = :value,
			created = :created,
			created_by = :created_by,
			updated = :updated,
			updated_by = :updated_by,
			deleted = :deleted,
			deleted_by = :deleted_by
		WHERE entity_id = :entity_id`
)

// PropertyMySQLRepo is the repository for Propertie implemented with MySQL backend
type PropertyMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup perform startup functions
func (r *PropertyMySQLRepo) Startup() {
	logger.Trace("Property repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *PropertyMySQLRepo) Shutdown() {
	logger.Trace("Property repository shutting down...")
}

// ExistsByID checks the existence of a Property by its ID
func (r *PropertyMySQLRepo) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("exists by ID", "Property", err)
	}
	return
}

// ExistsValueByID checks the existence of a Property Value by its ID
func (r *PropertyMySQLRepo) ExistsValueByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("exists by ID", "Property Value", err)
	}
	return
}

// ResolveByIDs resolves Properties by their IDs
func (r *PropertyMySQLRepo) ResolveByIDs(ids []uuid.UUID) (properties []model.Property, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(QuerySelectProperty+" WHERE properties.entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by IDs", "Property", err)
		return
	}

	err = r.DB.Select(&properties, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by IDs", "Property", err)
	}

	return
}

// ResolveValuesByIDs resolves Property Values by their IDs
func (r *PropertyMySQLRepo) ResolveValuesByIDs(ids []uuid.UUID) (vehicleValues []model.PropertyValue, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(QuerySelectPropertyValues+" WHERE property_values.entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by IDs", "Property Value", err)
		return
	}

	err = r.DB.Select(&vehicleValues, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by IDs", "Property Value", err)
	}

	return
}

// ResolveByFilter resolves Properties by a specified filter
func (r *PropertyMySQLRepo) ResolveByFilter(filter filter.Filter) (properties []model.Property, pageInfo model.PageInfoOutput, err error) {
	filterQueryString, err := filter.ToQueryString()
	if err != nil {
		err = failure.InternalError("resolve by filter", "Property", err)
		return properties, pageInfo, err
	}

	filterArgs := filter.GetArgs(true)
	query, args, err := r.DB.In(
		QuerySelectProperty+filterQueryString+filter.Pagination.ToQueryString(),
		filterArgs...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Property", err)
		return
	}

	err = r.DB.Select(&properties, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Property", err)
		return
	}

	var count int
	filterArgsNoPagination := filter.GetArgs(false)
	err = r.DB.Get(
		&count,
		"SELECT COUNT(entity_id) FROM properties "+filterQueryString,
		filterArgsNoPagination...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Property", err)
		properties = []model.Property{}
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

// ResolveValuesByFilter resolves Property Values by a specified filter
func (r *PropertyMySQLRepo) ResolveValuesByFilter(filter filter.Filter) (vehicleValues []model.PropertyValue, pageInfo model.PageInfoOutput, err error) {
	filterQueryString, err := filter.ToQueryString()
	if err != nil {
		err = failure.InternalError("resolve by filter", "Property Value", err)
		return
	}

	filterArgs := filter.GetArgs(true)
	query, args, err := r.DB.In(
		QuerySelectPropertyValues+filterQueryString+filter.Pagination.ToQueryString(),
		filterArgs...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Property Value", err)
		return
	}

	err = r.DB.Select(&vehicleValues, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Property Value", err)
		return
	}

	var count int
	filterArgsNoPagination := filter.GetArgs(false)
	query, args, err = r.DB.In(
		"SELECT COUNT(entity_id) FROM property_values "+filterQueryString,
		filterArgsNoPagination...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Property Value", err)
		vehicleValues = []model.PropertyValue{}
		return
	}

	err = r.DB.Get(&count, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve by filter", "Property Value", err)
		vehicleValues = []model.PropertyValue{}
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

// ResolveLastValuesByPropertyID resolves last X Property Values by their Property ID and count param
func (r *PropertyMySQLRepo) ResolveLastValuesByPropertyID(id uuid.UUID, count int) (vehicleValues []model.PropertyValue, err error) {
	if count == 0 {
		return
	}

	whereClause := " WHERE property_values.property_entity_id = ? and property_values.deleted IS NULL AND property_values.deleted_by IS NULL ORDER BY property_values.date DESC LIMIT ?"
	query, args, err := r.DB.In(QuerySelectPropertyValues+whereClause, id, count)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve last values", "Property Value", err)
		return
	}

	err = r.DB.Select(&vehicleValues, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("resolve last values", "Property Value", err)
		return
	}

	return
}

// Create creates a Property
func (r *PropertyMySQLRepo) Create(vehicle model.Property) error {
	exists, err := r.ExistsByID(vehicle.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if exists {
		err = failure.OperationNotPermitted("create", "Property", "already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreateProperty(tx, vehicle); err != nil {
			wrappedErr := failure.InternalError("create", "Property", err)
			e <- wrappedErr
			return
		}

		for _, value := range vehicle.Values {
			if err := r.txCreatePropertyValue(tx, value); err != nil {
				wrappedErr := failure.InternalError("create", "Property", err)
				e <- wrappedErr
				return
			}
		}

		e <- nil
	})
}

// Update updates a Property
func (r *PropertyMySQLRepo) Update(vehicle model.Property) error {
	exists, err := r.ExistsByID(vehicle.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if !exists {
		err = failure.EntityNotFound("update", "Property")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdateProperty(tx, vehicle); err != nil {
			err = failure.InternalError("update", "Property", err)
			e <- err
			return
		}

		e <- nil
	})
}

// CreateValue creates a new Property Value and optionally updates the Property transactionally
func (r *PropertyMySQLRepo) CreateValue(vehicleValue model.PropertyValue, vehicle *model.Property) error {
	exists, err := r.ExistsValueByID(vehicleValue.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if exists {
		err = failure.OperationNotPermitted("create", "Property Value", "already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreatePropertyValue(tx, vehicleValue); err != nil {
			err = failure.InternalError("create", "Property Value", err)
			e <- err
			return
		}

		if vehicle != nil {
			if err := r.txUpdateProperty(tx, *vehicle); err != nil {
				err = failure.InternalError("create", "Property Value", err)
				e <- err
				return
			}
		}

		e <- nil
	})
}

// UpdateValue updates an existing Property Value and optionally updates the Property transactionally
func (r *PropertyMySQLRepo) UpdateValue(vehicleValue model.PropertyValue, vehicle *model.Property) error {
	exists, err := r.ExistsValueByID(vehicleValue.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if !exists {
		err = failure.EntityNotFound("update", "Property Value")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdatePropertyValue(tx, vehicleValue); err != nil {
			err = failure.InternalError("update", "Property Value", err)
			e <- err
			return
		}

		if vehicle != nil {
			if err := r.txUpdateProperty(tx, *vehicle); err != nil {
				err = failure.InternalError("update", "Property Value", err)
				e <- err
				return
			}
		}

		e <- nil
	})
}

func (r *PropertyMySQLRepo) txCreateProperty(tx *sqlx.Tx, vehicle model.Property) error {
	stmt, err := tx.PrepareNamed(QueryInsertProperty)
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

func (r *PropertyMySQLRepo) txCreatePropertyValue(tx *sqlx.Tx, vehicleValue model.PropertyValue) error {
	stmt, err := tx.PrepareNamed(QueryInsertPropertyValue)
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

func (r *PropertyMySQLRepo) txUpdateProperty(tx *sqlx.Tx, vehicle model.Property) error {
	stmt, err := tx.PrepareNamed(QueryUpdateProperty)
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

func (r *PropertyMySQLRepo) txUpdatePropertyValue(tx *sqlx.Tx, vehicleValue model.PropertyValue) error {
	stmt, err := tx.PrepareNamed(QueryUpdatePropertyValue)
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
