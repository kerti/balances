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
	vehiclesStmtInsert = `INSERT INTO vehicles
	( entity_id, name, make, model, year, type, title_holder, license_plate_number, purchase_date, initial_value, initial_value_date, current_value, current_value_date, annual_depreciation_percent, status, created, created_by, updated, updated_by, deleted, deleted_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`

	vehicleValuesStmtInsert = `INSERT INTO vehicle_values
	( entity_id, vehicle_entity_id, date, value, created, created_by, updated, updated_by, deleted, deleted_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`

	vehiclesStmtUpdate = `UPDATE vehicle_values
	SET vehicle_entity_id = ?, date = ?, value = ?, created = ?, created_by = ?, updated = ?, updated_by = ?, deleted = ?, deleted_by = ?
	WHERE entity_id = ?`
)

type vehiclesRepositoryTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	repo               repository.Vehicle
	sqlmock            sqlmock.Sqlmock
	testUserID         uuid.UUID
	testVehicleID      uuid.UUID
	testVehicleValueID uuid.UUID
}

func TestVehiclesRepository(t *testing.T) {
	suite.Run(t, new(vehiclesRepositoryTestSuite))
}

func (t *vehiclesRepositoryTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	db, sqlmock := getMockedDriver(sqlmock.QueryMatcherEqual)
	repo := new(repository.VehicleMySQLRepo)
	repo.DB = &db
	t.repo = repo
	t.sqlmock = sqlmock
	t.testUserID, _ = uuid.NewV7()
	t.testVehicleID, _ = uuid.NewV7()
	t.testVehicleValueID, _ = uuid.NewV7()
	t.repo.Startup()
}

func (t *vehiclesRepositoryTestSuite) TearDownTest() {
	t.repo.Shutdown()
	t.ctrl.Finish()
}

func (t *vehiclesRepositoryTestSuite) getNewVehicleModel(id nuuid.NUUID, valueCount int) model.Vehicle {
	veh := model.Vehicle{}

	if id.Valid {
		veh.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		veh.ID = newID
	}

	veh.Name = "Test Name"
	veh.Make = "Test Make"
	veh.Model = "Test Model"
	veh.Year = 2020
	veh.Type = model.VehicleTypeCar
	veh.TitleHolder = "Test TitleHolder"
	veh.LicensePlateNumber = "Test LicensePlateNumber"
	veh.PurchaseDate = time.Now().AddDate(0, -1, -1)
	veh.InitialValue = float64(1000000)
	veh.InitialValueDate = time.Now().AddDate(0, 0, -1)
	veh.CurrentValue = float64(900000)
	veh.CurrentValueDate = time.Now().AddDate(0, 0, -1)
	veh.AnnualDepreciationPercent = 3.5
	veh.Status = model.VehicleStatusInUse
	veh.Created = time.Now().AddDate(0, -1, 0)
	veh.CreatedBy = t.testUserID
	veh.Updated = null.TimeFromPtr(nil)
	veh.UpdatedBy = nuuid.NUUID{Valid: false}
	veh.Deleted = null.TimeFromPtr(nil)
	veh.DeletedBy = nuuid.NUUID{Valid: false}

	for i := range valueCount {
		if i == valueCount-1 {
			veh.AttachValues(
				[]model.VehicleValue{
					t.getNewVehicleValueModel(
						nuuid.NUUID{Valid: false},
						nuuid.From(t.testVehicleID),
						null.TimeFrom(veh.CurrentValueDate),
						&veh.CurrentValue,
					),
				}, false)
		} else {
			veh.AttachValues(
				[]model.VehicleValue{
					t.getNewVehicleValueModel(
						nuuid.NUUID{Valid: false},
						nuuid.From(t.testVehicleID),
						null.TimeFromPtr(nil),
						nil,
					),
				}, false)
		}

	}

	return veh
}

func (t *vehiclesRepositoryTestSuite) getNewVehicleValueModel(id nuuid.NUUID, vehicleID nuuid.NUUID, date null.Time, value *float64) model.VehicleValue {
	vv := model.VehicleValue{}

	if id.Valid {
		vv.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		vv.ID = newID
	}

	if vehicleID.Valid {
		vv.VehicleID = vehicleID.UUID
	} else {
		newID, _ := uuid.NewV7()
		vv.VehicleID = newID
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

func (t *vehiclesRepositoryTestSuite) getArgsFromVehicleModel(vehicle model.Vehicle) (args []driver.Value) {
	args = append(args, vehicle.ID)
	args = append(args, vehicle.Name)
	args = append(args, vehicle.Make)
	args = append(args, vehicle.Model)
	args = append(args, vehicle.Year)
	args = append(args, vehicle.Type)
	args = append(args, vehicle.TitleHolder)
	args = append(args, vehicle.LicensePlateNumber)
	args = append(args, vehicle.PurchaseDate)
	args = append(args, vehicle.InitialValue)
	args = append(args, vehicle.InitialValueDate)
	args = append(args, vehicle.CurrentValue)
	args = append(args, vehicle.CurrentValueDate)
	args = append(args, vehicle.AnnualDepreciationPercent)
	args = append(args, vehicle.Status)
	args = append(args, vehicle.Created)
	args = append(args, vehicle.CreatedBy)
	args = append(args, vehicle.Updated)
	args = append(args, vehicle.UpdatedBy)
	args = append(args, vehicle.Deleted)
	args = append(args, vehicle.DeletedBy)

	return
}

func (t *vehiclesRepositoryTestSuite) getArgsFromVehicleValueModel(vehicleValue model.VehicleValue) (args []driver.Value) {
	args = append(args, vehicleValue.ID)
	args = append(args, vehicleValue.VehicleID)
	args = append(args, vehicleValue.Date)
	args = append(args, vehicleValue.Value)
	args = append(args, vehicleValue.Created)
	args = append(args, vehicleValue.CreatedBy)
	args = append(args, vehicleValue.Updated)
	args = append(args, vehicleValue.UpdatedBy)
	args = append(args, vehicleValue.Deleted)
	args = append(args, vehicleValue.DeletedBy)

	return
}

func (t *vehiclesRepositoryTestSuite) TestCreate_Normal() {
	testModel := t.getNewVehicleModel(nuuid.From(t.testVehicleID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ? ").
		WithArgs(t.testVehicleID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(vehiclesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromVehicleModel(testModel)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for _, valueModel := range testModel.Values {
		t.sqlmock.
			ExpectPrepare(vehicleValuesStmtInsert).
			ExpectExec().
			WithArgs(t.getArgsFromVehicleValueModel(valueModel)...).
			WillReturnResult(sqlmock.NewResult(1, 1))

	}

	t.sqlmock.ExpectCommit()

	err := t.repo.Create(testModel)

	assert.NoError(t.T(), err)
}

func (t *vehiclesRepositoryTestSuite) TestCreate_ErrorOnCheckExistence() {
	errMsg := "failed checking existence of vehicle"
	testModel := t.getNewVehicleModel(nuuid.From(t.testVehicleID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ? ").
		WithArgs(t.testVehicleID).
		WillReturnError(errors.New(errMsg))

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "checking existence")
}

func (t *vehiclesRepositoryTestSuite) TestCreate_AlreadyExists() {
	testModel := t.getNewVehicleModel(nuuid.From(t.testVehicleID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ? ").
		WithArgs(t.testVehicleID).
		WillReturnRows(getExistsResult(true))

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeOperationNotPermitted, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "already exists")
}

func (t *vehiclesRepositoryTestSuite) TestCreate_FailOnPrepareVehicleStatement() {
	errMsg := "failed preparing statement to insert vehicle"
	testModel := t.getNewVehicleModel(nuuid.From(t.testVehicleID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ? ").
		WithArgs(t.testVehicleID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(vehiclesStmtInsert).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *vehiclesRepositoryTestSuite) TestCreate_FailOnExecVehicleStatement() {
	errMsg := "failed executing insert vehicle statement"
	testModel := t.getNewVehicleModel(nuuid.From(t.testVehicleID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ? ").
		WithArgs(t.testVehicleID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(vehiclesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromVehicleModel(testModel)...).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *vehiclesRepositoryTestSuite) TestCreate_FailOnPrepareVehicleValueStatement() {
	errMsg := "failed preparing insert vehicle value statement"
	testModel := t.getNewVehicleModel(nuuid.From(t.testVehicleID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ? ").
		WithArgs(t.testVehicleID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(vehiclesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromVehicleModel(testModel)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.
		ExpectPrepare(vehicleValuesStmtInsert).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *vehiclesRepositoryTestSuite) TestCreate_FailOnExecVehicleValueStatement() {
	errMsg := "failed executing insert vehicle value statement"
	testModel := t.getNewVehicleModel(nuuid.From(t.testVehicleID), 2)

	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ? ").
		WithArgs(t.testVehicleID).
		WillReturnRows(getExistsResult(false))

	t.sqlmock.ExpectBegin()

	t.sqlmock.
		ExpectPrepare(vehiclesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromVehicleModel(testModel)...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.sqlmock.
		ExpectPrepare(vehicleValuesStmtInsert).
		ExpectExec().
		WithArgs(t.getArgsFromVehicleValueModel(testModel.Values[0])...).
		WillReturnError(errors.New(errMsg))

	t.sqlmock.ExpectRollback()

	err := t.repo.Create(testModel)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "create", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

// TODO: test create value here

func (t *vehiclesRepositoryTestSuite) TestExistsByID_Normal() {
	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ?").
		WithArgs(t.testVehicleID).
		WillReturnRows(getExistsResult(true))

	_, err := t.repo.ExistsByID(t.testVehicleID)

	assert.NoError(t.T(), err)
}

func (t *vehiclesRepositoryTestSuite) TestExistsByID_Error() {
	errMsg := "failed checking existence of vehicle by ID"
	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ?").
		WithArgs(t.testVehicleID).
		WillReturnError(errors.New(errMsg))

	_, err := t.repo.ExistsByID(t.testVehicleID)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *vehiclesRepositoryTestSuite) TestExistsValueByID_Normal() {
	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicle_values WHERE vehicle_values.entity_id = ?").
		WithArgs(t.testVehicleValueID).
		WillReturnRows(getExistsResult(true))

	_, err := t.repo.ExistsValueByID(t.testVehicleValueID)

	assert.NoError(t.T(), err)
}

func (t *vehiclesRepositoryTestSuite) TestExistsValueByID_Error() {
	errMsg := "failed checking existence of vehicle value by ID"
	t.sqlmock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicle_values WHERE vehicle_values.entity_id = ?").
		WithArgs(t.testVehicleValueID).
		WillReturnError(errors.New(errMsg))

	_, err := t.repo.ExistsValueByID(t.testVehicleValueID)

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "exists by ID", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *vehiclesRepositoryTestSuite) TestResolveByIDs_Normal_NoID() {
	res, err := t.repo.ResolveByIDs([]uuid.UUID{})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 0)
}

func (t *vehiclesRepositoryTestSuite) TestResolveByIDs_Normal_SingleID() {
	t.sqlmock.ExpectQuery(repository.QuerySelectVehicle + " WHERE vehicles.entity_id IN (?)").
		WithArgs(t.testVehicleID).
		WillReturnRows(getSingleEntityIDResult(t.testVehicleID))

	res, err := t.repo.ResolveByIDs([]uuid.UUID{t.testVehicleID})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
}

func (t *vehiclesRepositoryTestSuite) TestResolveByIDs_Normal_MultipleIDs() {
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	t.sqlmock.ExpectQuery(repository.QuerySelectVehicle+" WHERE vehicles.entity_id IN (?, ?)").
		WithArgs(id1, id2).
		WillReturnRows(getMultiEntityIDResult([]uuid.UUID{id1, id2}))

	res, err := t.repo.ResolveByIDs([]uuid.UUID{id1, id2})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 2)
}

func (t *vehiclesRepositoryTestSuite) TestResolveByIDs_ErrorExecutingSelect() {
	errMsg := "failed resolving vehicles by IDs"
	t.sqlmock.ExpectQuery(repository.QuerySelectVehicle + " WHERE vehicles.entity_id IN (?)").
		WithArgs(t.testVehicleID).
		WillReturnError(errors.New(errMsg))

	res, err := t.repo.ResolveByIDs([]uuid.UUID{t.testVehicleID})

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by IDs", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)

	assert.Len(t.T(), res, 0)
}

func (t *vehiclesRepositoryTestSuite) TestResolveByFilter_Normal() {
	keyword := "example"
	likeKeyword := "%example%"

	mock.
		ExpectQuery(repository.QuerySelectVehicle+"WHERE (((((((vehicles.name LIKE ?) OR (vehicles.make LIKE ?)) OR (vehicles.model LIKE ?)) OR (vehciles.year LIKE ?)) OR (vehicles.type LIKE ?)) OR (vehicles.title_holder LIKE ?))) AND vehicles.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
		WillReturnRows(getSingleEntityIDResult(t.testVehicleID))

	mock.ExpectQuery("SELECT COUNT(entity_id) FROM vehicles WHERE (((((((vehicles.name LIKE ?) OR (vehicles.make LIKE ?)) OR (vehicles.model LIKE ?)) OR (vehciles.year LIKE ?)) OR (vehicles.type LIKE ?)) OR (vehicles.title_holder LIKE ?))) AND vehicles.deleted IS NULL").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword).
		WillReturnRows(getCountResult(1))

	testFilter := model.VehicleFilterInput{}
	testFilter.Keyword = &keyword

	res, pageInfo, err := t.repo.ResolveByFilter(testFilter.ToFilter())

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
	assert.Equal(t.T(), 1, pageInfo.Page)
	assert.Equal(t.T(), 1, pageInfo.PageCount)
	assert.Equal(t.T(), 1, pageInfo.TotalCount)
	assert.Equal(t.T(), 10, pageInfo.PageSize)
}

func (t *vehiclesRepositoryTestSuite) TestResolveByFilter_ErrorOnSelect() {
	errMsg := "failed resolving vehicles by filter"
	keyword := "example"
	likeKeyword := "%example%"

	mock.
		ExpectQuery(repository.QuerySelectVehicle+"WHERE (((((((vehicles.name LIKE ?) OR (vehicles.make LIKE ?)) OR (vehicles.model LIKE ?)) OR (vehciles.year LIKE ?)) OR (vehicles.type LIKE ?)) OR (vehicles.title_holder LIKE ?))) AND vehicles.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
		WillReturnError(errors.New(errMsg))

	testFilter := model.VehicleFilterInput{}
	testFilter.Keyword = &keyword

	res, pageInfo, err := t.repo.ResolveByFilter(testFilter.ToFilter())

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by filter", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
	assert.Equal(t.T(), 0, pageInfo.Page)
	assert.Equal(t.T(), 0, pageInfo.PageCount)
	assert.Equal(t.T(), 0, pageInfo.TotalCount)
	assert.Equal(t.T(), 0, pageInfo.PageSize)
}

func (t *vehiclesRepositoryTestSuite) TestResolveByFilter_ErrorOnCount() {
	errMsg := "failed resolving vehicles by filter"
	keyword := "example"
	likeKeyword := "%example%"

	mock.
		ExpectQuery(repository.QuerySelectVehicle+"WHERE (((((((vehicles.name LIKE ?) OR (vehicles.make LIKE ?)) OR (vehicles.model LIKE ?)) OR (vehciles.year LIKE ?)) OR (vehicles.type LIKE ?)) OR (vehicles.title_holder LIKE ?))) AND vehicles.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
		WillReturnRows(getSingleEntityIDResult(t.testVehicleID))

	mock.ExpectQuery("SELECT COUNT(entity_id) FROM vehicles WHERE (((((((vehicles.name LIKE ?) OR (vehicles.make LIKE ?)) OR (vehicles.model LIKE ?)) OR (vehciles.year LIKE ?)) OR (vehicles.type LIKE ?)) OR (vehicles.title_holder LIKE ?))) AND vehicles.deleted IS NULL").
		WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword, likeKeyword).
		WillReturnError(errors.New(errMsg))

	testFilter := model.VehicleFilterInput{}
	testFilter.Keyword = &keyword

	res, pageInfo, err := t.repo.ResolveByFilter(testFilter.ToFilter())

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by filter", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
	assert.Equal(t.T(), 0, pageInfo.Page)
	assert.Equal(t.T(), 0, pageInfo.PageCount)
	assert.Equal(t.T(), 0, pageInfo.TotalCount)
	assert.Equal(t.T(), 0, pageInfo.PageSize)
}

func (t *vehiclesRepositoryTestSuite) TestResolveValuesByIDs_Normal_NoID() {
	res, err := t.repo.ResolveValuesByIDs([]uuid.UUID{})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 0)
}

func (t *vehiclesRepositoryTestSuite) TestResolveValuesByIDs_Normal_SingleID() {
	t.sqlmock.ExpectQuery(repository.QuerySelectVehicleValues + " WHERE vehicle_values.entity_id IN (?)").
		WithArgs(t.testVehicleValueID).
		WillReturnRows(getSingleEntityIDResult(t.testVehicleValueID))

	res, err := t.repo.ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
}

func (t *vehiclesRepositoryTestSuite) TestResolveValuesByIDs_Normal_MultipleIDs() {
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	t.sqlmock.ExpectQuery(repository.QuerySelectVehicleValues+" WHERE vehicle_values.entity_id IN (?, ?)").
		WithArgs(id1, id2).
		WillReturnRows(getMultiEntityIDResult([]uuid.UUID{id1, id2}))

	res, err := t.repo.ResolveValuesByIDs([]uuid.UUID{id1, id2})

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 2)
}

func (t *vehiclesRepositoryTestSuite) TestResolveValuesByIDs_ErrorExecutingSelect() {
	errMsg := "failed resolving vehicle values by IDs"
	t.sqlmock.ExpectQuery(repository.QuerySelectVehicleValues + " WHERE vehicle_values.entity_id IN (?)").
		WithArgs(t.testVehicleValueID).
		WillReturnError(errors.New(errMsg))

	res, err := t.repo.ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID})

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by IDs", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)

	assert.Len(t.T(), res, 0)
}

func (t *vehiclesRepositoryTestSuite) TestResolveValuesByFilter_Normal() {
	mock.
		ExpectQuery(repository.QuerySelectVehicleValues+"WHERE ((vehicle_values.vehicle_entity_id IN (?))) AND vehicle_values.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(t.testVehicleID, 10, 0).
		WillReturnRows(getSingleEntityIDResult(t.testVehicleID))

	mock.ExpectQuery("SELECT COUNT(entity_id) FROM vehicle_values WHERE ((vehicle_values.vehicle_entity_id IN (?))) AND vehicle_values.deleted IS NULL").
		WithArgs(t.testVehicleID).
		WillReturnRows(getCountResult(1))

	testFilter := model.VehicleValueFilterInput{}
	testFilter.VehicleIDs = &[]uuid.UUID{t.testVehicleID}

	res, pageInfo, err := t.repo.ResolveValuesByFilter(testFilter.ToFilter())

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
	assert.Equal(t.T(), 1, pageInfo.Page)
	assert.Equal(t.T(), 1, pageInfo.PageCount)
	assert.Equal(t.T(), 1, pageInfo.TotalCount)
	assert.Equal(t.T(), 10, pageInfo.PageSize)
}

func (t *vehiclesRepositoryTestSuite) TestResolveValuesByFilter_ErrorOnSelect() {
	errMsg := "failed resolving vehicle values by filter"

	mock.
		ExpectQuery(repository.QuerySelectVehicleValues+"WHERE ((vehicle_values.vehicle_entity_id IN (?))) AND vehicle_values.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(t.testVehicleID, 10, 0).
		WillReturnError(errors.New(errMsg))

	testFilter := model.VehicleValueFilterInput{}
	testFilter.VehicleIDs = &[]uuid.UUID{t.testVehicleID}

	res, pageInfo, err := t.repo.ResolveValuesByFilter(testFilter.ToFilter())

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by filter", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
	assert.Equal(t.T(), 0, pageInfo.Page)
	assert.Equal(t.T(), 0, pageInfo.PageCount)
	assert.Equal(t.T(), 0, pageInfo.TotalCount)
	assert.Equal(t.T(), 0, pageInfo.PageSize)
}

func (t *vehiclesRepositoryTestSuite) TestResolveValuesByFilter_ErrorOnCount() {
	errMsg := "failed resolving vehicle values by filter"
	mock.
		ExpectQuery(repository.QuerySelectVehicleValues+"WHERE ((vehicle_values.vehicle_entity_id IN (?))) AND vehicle_values.deleted IS NULL LIMIT ? OFFSET ?").
		WithArgs(t.testVehicleID, 10, 0).
		WillReturnRows(getSingleEntityIDResult(t.testVehicleID))

	mock.ExpectQuery("SELECT COUNT(entity_id) FROM vehicle_values WHERE ((vehicle_values.vehicle_entity_id IN (?))) AND vehicle_values.deleted IS NULL").
		WithArgs(t.testVehicleID).
		WillReturnError(errors.New(errMsg))

	testFilter := model.VehicleValueFilterInput{}
	testFilter.VehicleIDs = &[]uuid.UUID{t.testVehicleID}

	res, pageInfo, err := t.repo.ResolveValuesByFilter(testFilter.ToFilter())

	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle Value", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by filter", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Len(t.T(), res, 0)
	assert.Equal(t.T(), 0, pageInfo.Page)
	assert.Equal(t.T(), 0, pageInfo.PageCount)
	assert.Equal(t.T(), 0, pageInfo.TotalCount)
	assert.Equal(t.T(), 0, pageInfo.PageSize)
}
