package repository_test

import (
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/nuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	propertiesStmtInsert = `INSERT INTO properties
	( entity_id, name, address, total_area, building_area, area_unit, type, title_holder, tax_identifier, purchase_date, initial_value, initial_value_date, current_value, current_value_date, annual_appreciation_percent, status, created, created_by, updated, updated_by, deleted, deleted_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`

	propertyValuesStmtInsert = `INSERT INTO property_values
	( entity_id, property_entity_id, date, value, created, created_by, updated, updated_by, deleted, deleted_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`

	propertiesStmtUpdate = `UPDATE properties
	SET name = ?, address = ?, total_area = ?, building_area = ?, area_unit = ?, type = ?, title_holder = ?, tax_identifier = ?, purchase_date = ?, initial_value = ?, initial_value_date = ?, current_value = ?, current_value_date = ?, annual_appreciation_percent = ?, status = ?, created = ?, created_by = ?, updated = ?, updated_by = ?, deleted = ?, deleted_by = ?
	WHERE entity_id = ?`

	propertyValuesStmtUpdate = `UPDATE property_values
	SET property_entity_id = ?, date = ?, value = ?, created = ?, created_by = ?, updated = ?, updated_by = ?, deleted = ?, deleted_by = ?
	WHERE entity_id = ?`
)

type propertiesRepositoryTestSuite struct {
	suite.Suite
	ctrl                *gomock.Controller
	repo                repository.Property
	sqlmock             sqlmock.Sqlmock
	testUserID          uuid.UUID
	testPropertyID      uuid.UUID
	testPropertyValueID uuid.UUID
}

func TestPropertiesRepository(t *testing.T) {
	suite.Run(t, new(propertiesRepositoryTestSuite))
}

func (t *propertiesRepositoryTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	db, sqlmock := getMockedDriver(sqlmock.QueryMatcherEqual)
	repo := new(repository.PropertyMySQLRepo)
	repo.DB = &db
	t.repo = repo
	t.sqlmock = sqlmock
	t.testUserID, _ = uuid.NewV7()
	t.testPropertyID, _ = uuid.NewV7()
	t.testPropertyValueID, _ = uuid.NewV7()
	t.repo.Startup()
}

func (t *propertiesRepositoryTestSuite) TearDownTest() {
	t.repo.Shutdown()
	t.ctrl.Finish()
}

func (t *propertiesRepositoryTestSuite) getNewPropertyModel(id nuuid.NUUID, valueCount int) model.Property {
	prop := model.Property{}

	if id.Valid {
		prop.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		prop.ID = newID
	}

	prop.Name = "Test Name"
	prop.Address = "Test Address"
	prop.TotalArea = 1500
	prop.BuildingArea = 1200
	prop.AreaUnit = model.PropertyAreaUnitSQM
	prop.Type = model.PropertyTypeHouse
	prop.TitleHolder = "Test TitleHolder"
	prop.TaxIdentifier = "Test TaxIdentifier"
	prop.PurchaseDate = time.Now().AddDate(0, -1, -1)
	prop.InitialValue = float64(1000000)
	prop.InitialValueDate = time.Now().AddDate(0, 0, -1)
	prop.CurrentValue = float64(900000)
	prop.CurrentValueDate = time.Now().AddDate(0, 0, -1)
	prop.AnnualAppreciationPercent = 3.5
	prop.Status = model.PropertyStatusInUse
	prop.Created = time.Now().AddDate(0, -1, 0)
	prop.CreatedBy = t.testUserID
	prop.Updated = null.TimeFromPtr(nil)
	prop.UpdatedBy = nuuid.NUUID{Valid: false}
	prop.Deleted = null.TimeFromPtr(nil)
	prop.DeletedBy = nuuid.NUUID{Valid: false}

	for i := range valueCount {
		if i == valueCount-1 {
			prop.AttachValues(
				[]model.PropertyValue{
					t.getNewPropertyValueModel(
						nuuid.NUUID{Valid: false},
						nuuid.From(t.testPropertyID),
						null.TimeFrom(prop.CurrentValueDate),
						&prop.CurrentValue,
					),
				}, false)
		} else {
			prop.AttachValues(
				[]model.PropertyValue{
					t.getNewPropertyValueModel(
						nuuid.NUUID{Valid: false},
						nuuid.From(t.testPropertyID),
						null.TimeFromPtr(nil),
						nil,
					),
				}, false)
		}

	}

	return prop
}

func (t *propertiesRepositoryTestSuite) getNewPropertyValueModel(id nuuid.NUUID, propertyID nuuid.NUUID, date null.Time, value *float64) model.PropertyValue {
	vv := model.PropertyValue{}

	if id.Valid {
		vv.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		vv.ID = newID
	}

	if propertyID.Valid {
		vv.PropertyID = propertyID.UUID
	} else {
		newID, _ := uuid.NewV7()
		vv.PropertyID = newID
	}

	if date.Valid {
		vv.Date = date.Time
	} else {
		vv.Date = time.Now()
	}

	if value != nil {
		vv.Value = *value
	} else {
		vv.Value = 123123123
	}

	vv.Created = time.Now().AddDate(0, -1, 0)
	vv.CreatedBy = t.testUserID
	vv.Updated = null.TimeFromPtr(nil)
	vv.UpdatedBy = nuuid.NUUID{Valid: false}
	vv.Deleted = null.TimeFromPtr(nil)
	vv.DeletedBy = nuuid.NUUID{Valid: false}

	return vv
}

func (t *propertiesRepositoryTestSuite) getArgsFromPropertyModel(property model.Property, setIdLast bool) (args []driver.Value) {
	if !setIdLast {
		args = append(args, property.ID)
	}

	args = append(args, property.Name)
	args = append(args, property.Address)
	args = append(args, property.TotalArea)
	args = append(args, property.BuildingArea)
	args = append(args, property.AreaUnit)
	args = append(args, property.Type)
	args = append(args, property.TitleHolder)
	args = append(args, property.TaxIdentifier)
	args = append(args, property.PurchaseDate)
	args = append(args, property.InitialValue)
	args = append(args, property.InitialValueDate)
	args = append(args, property.CurrentValue)
	args = append(args, property.CurrentValueDate)
	args = append(args, property.AnnualAppreciationPercent)
	args = append(args, property.Status)
	args = append(args, property.Created)
	args = append(args, property.CreatedBy)
	args = append(args, property.Updated)
	args = append(args, property.UpdatedBy)
	args = append(args, property.Deleted)
	args = append(args, property.DeletedBy)

	if setIdLast {
		args = append(args, property.ID)
	}

	return
}

func (t *propertiesRepositoryTestSuite) getArgsFromPropertyValueModel(propertyValue model.PropertyValue, setIdLast bool) (args []driver.Value) {
	if !setIdLast {
		args = append(args, propertyValue.ID)
	}

	args = append(args, propertyValue.PropertyID)
	args = append(args, propertyValue.Date)
	args = append(args, propertyValue.Value)
	args = append(args, propertyValue.Created)
	args = append(args, propertyValue.CreatedBy)
	args = append(args, propertyValue.Updated)
	args = append(args, propertyValue.UpdatedBy)
	args = append(args, propertyValue.Deleted)
	args = append(args, propertyValue.DeletedBy)

	if setIdLast {
		args = append(args, propertyValue.ID)
	}

	return
}

func (t *propertiesRepositoryTestSuite) TestCreate_Normal() {
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ? ").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertiesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyModel(testModel, false)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for _, valueModel := range testModel.Values {
		t.sqlmock.
			ExpectPrepare(propertyValuesStmtInsert).
			ExpectExec().
			WithArgs(t.getArgsFromPropertyValueModel(valueModel, false)...).
			WillReturnResult(sqlmock.NewResult(1, 1))

	}

	t.sqlmock.ExpectCommit()

	err := t.repo.Create(testModel)

	assert.NoError(t.T(), err)
}

func (t *propertiesRepositoryTestSuite) TestCreate_ErrorOnCheckExistence() {
	errMsg := "failed checking existence of property"
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ? ").
		WithArgs(t.testPropertyID).
		WillReturnError(errors.New(errMsg))

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "checking existence")
}

func (t *propertiesRepositoryTestSuite) TestCreate_AlreadyExists() {
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ? ").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(true))

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeOperationNotPermitted, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "already exists")
}

func (t *propertiesRepositoryTestSuite) TestCreate_FailOnPreparePropertyStatement() {
	errMsg := "failed preparing statement to insert property"
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ? ").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertiesStmtInsert).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestCreate_FailOnExecPropertyStatement() {
	errMsg := "failed executing insert property statement"
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ? ").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertiesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyModel(testModel, false)...).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestCreate_FailOnPreparePropertyValueStatement() {
	errMsg := "failed preparing insert property value statement"
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ? ").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertiesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyModel(testModel, false)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtInsert).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestCreate_FailOnExecPropertyValueStatement() {
	errMsg := "failed executing insert property value statement"
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ? ").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertiesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyModel(testModel, false)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyValueModel(testModel.Values[0], false)...).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestCreateValue_Normal() {
	newValue := t.getNewPropertyValueModel(nuuid.NUUID{Valid: false}, nuuid.From(t.testPropertyID), null.TimeFromPtr(nil), nil)
	property := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)
	property.CurrentValue = newValue.Value
	property.CurrentValueDate = newValue.Date
	property.Updated = null.TimeFrom(time.Now())

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(newValue.ID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyValueModel(newValue, false)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.
		ExpectPrepare(propertiesStmtUpdate).
		ExpectExec().
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.ExpectCommit()

	err := t.repo.CreateValue(newValue, &property)

	assert.NoError(t.T(), err)
}

func (t *propertiesRepositoryTestSuite) TestCreateValue_NoPropertyUpdate() {
	newValue := t.getNewPropertyValueModel(nuuid.NUUID{Valid: false}, nuuid.From(t.testPropertyID), null.TimeFromPtr(nil), nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(newValue.ID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyValueModel(newValue, false)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.ExpectCommit()

	err := t.repo.CreateValue(newValue, nil)

	assert.NoError(t.T(), err)
}

func (t *propertiesRepositoryTestSuite) TestCreateValue_ErrorOnCheckExistence() {
	errMsg := "failed checking existence of property value"
	newValue := t.getNewPropertyValueModel(nuuid.NUUID{Valid: false}, nuuid.From(t.testPropertyID), null.TimeFromPtr(nil), nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(newValue.ID).
		WillReturnError(failure.InternalError("exists by ID", "Property Value", errors.New(errMsg)))

	err := t.repo.CreateValue(newValue, nil)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestCreateValue_AlreadyExists() {
	newValue := t.getNewPropertyValueModel(nuuid.NUUID{Valid: false}, nuuid.From(t.testPropertyID), null.TimeFromPtr(nil), nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(newValue.ID).
		WillReturnRows(getExistsResult(true))

	err := t.repo.CreateValue(newValue, nil)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeOperationNotPermitted, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "already exists")
}

func (t *propertiesRepositoryTestSuite) TestCreateValue_FailOnPrepare() {
	errMsg := "failed preparing statement for creating property value"
	newValue := t.getNewPropertyValueModel(nuuid.NUUID{Valid: false}, nuuid.From(t.testPropertyID), null.TimeFromPtr(nil), nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(newValue.ID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtInsert).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.CreateValue(newValue, nil)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestCreateValue_FailOnExec() {
	errMsg := "failed executing statement to create property value"
	newValue := t.getNewPropertyValueModel(nuuid.NUUID{Valid: false}, nuuid.From(t.testPropertyID), null.TimeFromPtr(nil), nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(newValue.ID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyValueModel(newValue, false)...).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.CreateValue(newValue, nil)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestCreateValue_FailOnPropertyUpdate() {
	errMsg := "failed executing statement to update property"
	newValue := t.getNewPropertyValueModel(nuuid.NUUID{Valid: false}, nuuid.From(t.testPropertyID), null.TimeFromPtr(nil), nil)
	property := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)
	property.CurrentValue = newValue.Value
	property.CurrentValueDate = newValue.Date
	property.Updated = null.TimeFrom(time.Now())

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(newValue.ID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyValueModel(newValue, false)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.
		ExpectPrepare(propertiesStmtUpdate).
		ExpectExec().
		WithArgs().
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.CreateValue(newValue, &property)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestExistsByID_Normal() {
	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ?").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(true))

	_, err := t.repo.ExistsByID(t.testPropertyID)

	assert.NoError(t.T(), err)
}

func (t *propertiesRepositoryTestSuite) TestExistsByID_Error() {
	errMsg := "failed checking existence of property by ID"
	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ?").
		WithArgs(t.testPropertyID).
		WillReturnError(errors.New(errMsg))

	_, err := t.repo.ExistsByID(t.testPropertyID)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestExistsValueByID_Normal() {
	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnRows(getExistsResult(true))

	_, err := t.repo.ExistsValueByID(t.testPropertyValueID)

	assert.NoError(t.T(), err)
}

func (t *propertiesRepositoryTestSuite) TestExistsValueByID_Error() {
	errMsg := "failed checking existence of property value by ID"
	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnError(errors.New(errMsg))

	_, err := t.repo.ExistsValueByID(t.testPropertyValueID)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestResolveByIDs_Normal_NoID() {
	res, err := t.repo.ResolveByIDs([]uuid.UUID{})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 0)
}

func (t *propertiesRepositoryTestSuite) TestResolveByIDs_Normal_SingleID() {
	t.sqlmock.ExpectQuery(repository.QuerySelectProperty + " WHERE properties.entity_id IN (?)").
		WithArgs(t.testPropertyID).
		WillReturnRows(getSingleEntityIDResult(t.testPropertyID))

	res, err := t.repo.ResolveByIDs([]uuid.UUID{t.testPropertyID})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
}

func (t *propertiesRepositoryTestSuite) TestResolveByIDs_Normal_MultipleIDs() {
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	t.sqlmock.ExpectQuery(repository.QuerySelectProperty+" WHERE properties.entity_id IN (?, ?)").
		WithArgs(id1, id2).
		WillReturnRows(getMultiEntityIDResult([]uuid.UUID{id1, id2}))

	res, err := t.repo.ResolveByIDs([]uuid.UUID{id1, id2})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 2)
}

func (t *propertiesRepositoryTestSuite) TestResolveByIDs_ErrorExecutingSelect() {
	errMsg := "failed resolving properties by IDs"
	t.sqlmock.ExpectQuery(repository.QuerySelectProperty + " WHERE properties.entity_id IN (?)").
		WithArgs(t.testPropertyID).
		WillReturnError(errors.New(errMsg))

	res, err := t.repo.ResolveByIDs([]uuid.UUID{t.testPropertyID})

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by IDs", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)

	assert.Len(t.T(), res, 0)
}

func (t *propertiesRepositoryTestSuite) TestResolveByFilter_Normal() {
	keyword := "example"
	likeKeyword := "%example%"

	t.sqlmock.
		ExpectQuery(repository.QuerySelectProperty+"WHERE ((((((properties.name LIKE ?) OR (properties.address LIKE ?)) OR (properties.type LIKE ?)) OR (properties.title_holder LIKE ?)) OR (properties.tax_identifier LIKE ?))) AND properties.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
		WillReturnRows(getSingleEntityIDResult(t.testPropertyID))

	t.sqlmock.ExpectQuery("SELECT COUNT(entity_id) FROM properties WHERE ((((((properties.name LIKE ?) OR (properties.address LIKE ?)) OR (properties.type LIKE ?)) OR (properties.title_holder LIKE ?)) OR (properties.tax_identifier LIKE ?))) AND properties.deleted IS NULL").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword).
		WillReturnRows(getCountResult(1))

	testFilter := model.PropertyFilterInput{}
	testFilter.Keyword = &keyword

	res, pageInfo, err := t.repo.ResolveByFilter(testFilter.ToFilter())

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
	assert.Equal(t.T(), 1, pageInfo.Page)
	assert.Equal(t.T(), 1, pageInfo.PageCount)
	assert.Equal(t.T(), 1, pageInfo.TotalCount)
	assert.Equal(t.T(), 10, pageInfo.PageSize)
}

func (t *propertiesRepositoryTestSuite) TestResolveByFilter_ErrorOnSelect() {
	errMsg := "failed resolving properties by filter"
	keyword := "example"
	likeKeyword := "%example%"

	t.sqlmock.
		ExpectQuery(repository.QuerySelectProperty+"WHERE ((((((properties.name LIKE ?) OR (properties.address LIKE ?)) OR (properties.type LIKE ?)) OR (properties.title_holder LIKE ?)) OR (properties.tax_identifier LIKE ?))) AND properties.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
		WillReturnError(errors.New(errMsg))

	testFilter := model.PropertyFilterInput{}
	testFilter.Keyword = &keyword

	res, pageInfo, err := t.repo.ResolveByFilter(testFilter.ToFilter())

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by filter", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
	assert.Equal(t.T(), 0, pageInfo.Page)
	assert.Equal(t.T(), 0, pageInfo.PageCount)
	assert.Equal(t.T(), 0, pageInfo.TotalCount)
	assert.Equal(t.T(), 0, pageInfo.PageSize)
}

func (t *propertiesRepositoryTestSuite) TestResolveByFilter_ErrorOnCount() {
	errMsg := "failed resolving properties by filter"
	keyword := "example"
	likeKeyword := "%example%"

	t.sqlmock.
		ExpectQuery(repository.QuerySelectProperty+"WHERE ((((((properties.name LIKE ?) OR (properties.address LIKE ?)) OR (properties.type LIKE ?)) OR (properties.title_holder LIKE ?)) OR (properties.tax_identifier LIKE ?))) AND properties.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
		WillReturnRows(getSingleEntityIDResult(t.testPropertyID))

	t.sqlmock.ExpectQuery("SELECT COUNT(entity_id) FROM properties WHERE ((((((properties.name LIKE ?) OR (properties.address LIKE ?)) OR (properties.type LIKE ?)) OR (properties.title_holder LIKE ?)) OR (properties.tax_identifier LIKE ?))) AND properties.deleted IS NULL").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword).
		WillReturnError(errors.New(errMsg))

	testFilter := model.PropertyFilterInput{}
	testFilter.Keyword = &keyword

	res, pageInfo, err := t.repo.ResolveByFilter(testFilter.ToFilter())

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by filter", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
	assert.Equal(t.T(), 0, pageInfo.Page)
	assert.Equal(t.T(), 0, pageInfo.PageCount)
	assert.Equal(t.T(), 0, pageInfo.TotalCount)
	assert.Equal(t.T(), 0, pageInfo.PageSize)
}

func (t *propertiesRepositoryTestSuite) TestResolveValuesByIDs_Normal_NoID() {
	res, err := t.repo.ResolveValuesByIDs([]uuid.UUID{})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 0)
}

func (t *propertiesRepositoryTestSuite) TestResolveValuesByIDs_Normal_SingleID() {
	t.sqlmock.ExpectQuery(repository.QuerySelectPropertyValues + " WHERE property_values.entity_id IN (?)").
		WithArgs(t.testPropertyValueID).
		WillReturnRows(getSingleEntityIDResult(t.testPropertyValueID))

	res, err := t.repo.ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
}

func (t *propertiesRepositoryTestSuite) TestResolveValuesByIDs_Normal_MultipleIDs() {
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	t.sqlmock.ExpectQuery(repository.QuerySelectPropertyValues+" WHERE property_values.entity_id IN (?, ?)").
		WithArgs(id1, id2).
		WillReturnRows(getMultiEntityIDResult([]uuid.UUID{id1, id2}))

	res, err := t.repo.ResolveValuesByIDs([]uuid.UUID{id1, id2})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 2)
}

func (t *propertiesRepositoryTestSuite) TestResolveValuesByIDs_ErrorExecutingSelect() {
	errMsg := "failed resolving property values by IDs"
	t.sqlmock.ExpectQuery(repository.QuerySelectPropertyValues + " WHERE property_values.entity_id IN (?)").
		WithArgs(t.testPropertyValueID).
		WillReturnError(errors.New(errMsg))

	res, err := t.repo.ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID})

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by IDs", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)

	assert.Len(t.T(), res, 0)
}

func (t *propertiesRepositoryTestSuite) TestResolveValuesByFilter_Normal() {
	t.sqlmock.
		ExpectQuery(repository.QuerySelectPropertyValues+"WHERE ((property_values.property_entity_id IN (?))) AND property_values.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(t.testPropertyID, 10, 0).
		WillReturnRows(getSingleEntityIDResult(t.testPropertyID))

	t.sqlmock.ExpectQuery("SELECT COUNT(entity_id) FROM property_values WHERE ((property_values.property_entity_id IN (?))) AND property_values.deleted IS NULL").
		WithArgs(t.testPropertyID).
		WillReturnRows(getCountResult(1))

	testFilter := model.PropertyValueFilterInput{}
	testFilter.PropertyIDs = &[]uuid.UUID{t.testPropertyID}

	res, pageInfo, err := t.repo.ResolveValuesByFilter(testFilter.ToFilter())

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
	assert.Equal(t.T(), 1, pageInfo.Page)
	assert.Equal(t.T(), 1, pageInfo.PageCount)
	assert.Equal(t.T(), 1, pageInfo.TotalCount)
	assert.Equal(t.T(), 10, pageInfo.PageSize)
}

func (t *propertiesRepositoryTestSuite) TestResolveValuesByFilter_ErrorOnSelect() {
	errMsg := "failed resolving property values by filter"

	t.sqlmock.
		ExpectQuery(repository.QuerySelectPropertyValues+"WHERE ((property_values.property_entity_id IN (?))) AND property_values.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(t.testPropertyID, 10, 0).
		WillReturnError(errors.New(errMsg))

	testFilter := model.PropertyValueFilterInput{}
	testFilter.PropertyIDs = &[]uuid.UUID{t.testPropertyID}

	res, pageInfo, err := t.repo.ResolveValuesByFilter(testFilter.ToFilter())

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by filter", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
	assert.Equal(t.T(), 0, pageInfo.Page)
	assert.Equal(t.T(), 0, pageInfo.PageCount)
	assert.Equal(t.T(), 0, pageInfo.TotalCount)
	assert.Equal(t.T(), 0, pageInfo.PageSize)
}

func (t *propertiesRepositoryTestSuite) TestResolveValuesByFilter_ErrorOnCount() {
	errMsg := "failed resolving property values by filter"
	t.sqlmock.
		ExpectQuery(repository.QuerySelectPropertyValues+"WHERE ((property_values.property_entity_id IN (?))) AND property_values.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(t.testPropertyID, 10, 0).
		WillReturnRows(getSingleEntityIDResult(t.testPropertyID))

	t.sqlmock.ExpectQuery("SELECT COUNT(entity_id) FROM property_values WHERE ((property_values.property_entity_id IN (?))) AND property_values.deleted IS NULL").
		WithArgs(t.testPropertyID).
		WillReturnError(errors.New(errMsg))

	testFilter := model.PropertyValueFilterInput{}
	testFilter.PropertyIDs = &[]uuid.UUID{t.testPropertyID}

	res, pageInfo, err := t.repo.ResolveValuesByFilter(testFilter.ToFilter())

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by filter", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
	assert.Equal(t.T(), 0, pageInfo.Page)
	assert.Equal(t.T(), 0, pageInfo.PageCount)
	assert.Equal(t.T(), 0, pageInfo.TotalCount)
	assert.Equal(t.T(), 0, pageInfo.PageSize)
}

func (t *propertiesRepositoryTestSuite) TestResolveLastValuesByPropertyID_Normal() {
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	t.sqlmock.ExpectQuery(repository.QuerySelectPropertyValues+"WHERE property_values.property_entity_id = ? and property_values.deleted IS NULL AND property_values.deleted_by IS NULL ORDER BY property_values.date DESC LIMIT ?").
		WithArgs(t.testPropertyID, 2).
		WillReturnRows(getMultiEntityIDResult([]uuid.UUID{id1, id2}))

	res, err := t.repo.ResolveLastValuesByPropertyID(t.testPropertyID, 2)

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 2)
}

func (t *propertiesRepositoryTestSuite) TestResolveLastValuesByPropertyID_CountZero() {
	res, err := t.repo.ResolveLastValuesByPropertyID(t.testPropertyID, 0)

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 0)
}

func (t *propertiesRepositoryTestSuite) TestResolveLastValuesByPropertyID_FailOnSelect() {
	errMsg := "failed resolving last values"

	t.sqlmock.ExpectQuery(repository.QuerySelectPropertyValues+"WHERE property_values.property_entity_id = ? and property_values.deleted IS NULL AND property_values.deleted_by IS NULL ORDER BY property_values.date DESC LIMIT ?").
		WithArgs(t.testPropertyID, 2).
		WillReturnError(errors.New(errMsg))

	res, err := t.repo.ResolveLastValuesByPropertyID(t.testPropertyID, 2)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve last values", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
}

func (t *propertiesRepositoryTestSuite) TestUpdate_Normal() {
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ?").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(true))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertiesStmtUpdate).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyModel(testModel, true)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.ExpectCommit()

	err := t.repo.Update(testModel)

	assert.NoError(t.T(), err)
}

func (t *propertiesRepositoryTestSuite) TestUpdate_ErrorOnCheckExistence() {
	errMsg := "failed checking the existence of property"
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ?").
		WithArgs(t.testPropertyID).
		WillReturnError(errors.New(errMsg))

	err := t.repo.Update(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestUpdate_DoesNotExist() {
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ?").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(false))

	err := t.repo.Update(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeEntityNotFound, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "Record not found")
}

func (t *propertiesRepositoryTestSuite) TestUpdate_FailOnPrepare() {
	errMsg := "failed preparing update statemtnt for property"
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ?").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(true))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertiesStmtUpdate).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Update(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestUpdate_FailOnExec() {
	errMsg := "failed executing update statement for property"
	testModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM properties WHERE properties.entity_id = ?").
		WithArgs(t.testPropertyID).
		WillReturnRows(getExistsResult(true))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertiesStmtUpdate).
		ExpectExec().
		WithArgs(t.getArgsFromPropertyModel(testModel, true)...).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Update(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestUpdateValue_Normal() {
	propertyValueModel := t.getNewPropertyValueModel(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		null.TimeFromPtr(nil),
		nil)

	propertyModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnRows(getExistsResult(true))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtUpdate).
		ExpectExec().
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.
		ExpectPrepare(propertiesStmtUpdate).
		ExpectExec().
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.ExpectCommit()

	err := t.repo.UpdateValue(propertyValueModel, &propertyModel)

	assert.NoError(t.T(), err)
}

func (t *propertiesRepositoryTestSuite) TestUpdateValue_Normal_NoAccountUpdate() {
	propertyValueModel := t.getNewPropertyValueModel(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		null.TimeFromPtr(nil),
		nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnRows(getExistsResult(true))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtUpdate).
		ExpectExec().
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.ExpectCommit()

	err := t.repo.UpdateValue(propertyValueModel, nil)

	assert.NoError(t.T(), err)
}

func (t *propertiesRepositoryTestSuite) TestUpdateValue_ErrorOnCheckExistence() {
	errMsg := "failed checking the existence of property value"
	propertyValueModel := t.getNewPropertyValueModel(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		null.TimeFromPtr(nil),
		nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnError(failure.InternalError("exists by ID", "Property Value", errors.New(errMsg)))

	err := t.repo.UpdateValue(propertyValueModel, nil)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestUpdateValue_DoesNotExist() {
	propertyValueModel := t.getNewPropertyValueModel(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		null.TimeFromPtr(nil),
		nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnRows(getExistsResult(false))

	err := t.repo.UpdateValue(propertyValueModel, nil)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeEntityNotFound, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "not found")
}

func (t *propertiesRepositoryTestSuite) TestUpdateValue_FailOnPrepare() {
	errMsg := "failed preparig statement to update property value"
	propertyValueModel := t.getNewPropertyValueModel(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		null.TimeFromPtr(nil),
		nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnRows(getExistsResult(true))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtUpdate).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.UpdateValue(propertyValueModel, nil)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestUpdateValue_FailOnExec() {
	errMsg := "failed preparig statement to update property value"
	propertyValueModel := t.getNewPropertyValueModel(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		null.TimeFromPtr(nil),
		nil)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnRows(getExistsResult(true))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtUpdate).
		ExpectExec().
		WithArgs().
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.UpdateValue(propertyValueModel, nil)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesRepositoryTestSuite) TestUpdateValue_FailOnPropertyUpdate() {
	errMsg := "failed executing statement to update property"
	propertyValueModel := t.getNewPropertyValueModel(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		null.TimeFromPtr(nil),
		nil)

	propertyModel := t.getNewPropertyModel(nuuid.From(t.testPropertyID), 0)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM property_values WHERE property_values.entity_id = ?").
		WithArgs(t.testPropertyValueID).
		WillReturnRows(getExistsResult(true))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(propertyValuesStmtUpdate).
		ExpectExec().
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.
		ExpectPrepare(propertiesStmtUpdate).
		ExpectExec().
		WithArgs().
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.UpdateValue(propertyValueModel, &propertyModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}
