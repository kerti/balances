package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
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

func (t *vehiclesRepositoryTestSuite) TestCreate_Normal() {
	testModel := t.getNewVehicleModel(nuuid.From(t.testVehicleID), 2)

	mock.
		ExpectQuery("SELECT COUNT(entity_id) > 0 FROM vehicles WHERE vehicles.entity_id = ? ").
		WithArgs(t.testVehicleID).
		WillReturnRows(getExistsResult(false))

	mock.ExpectBegin()

	mock.
		ExpectPrepare(vehiclesStmtInsert).
		ExpectExec().
		WithArgs(
			testModel.ID,
			testModel.Name,
			testModel.Make,
			testModel.Model,
			testModel.Year,
			testModel.Type,
			testModel.TitleHolder,
			testModel.LicensePlateNumber,
			testModel.PurchaseDate,
			testModel.InitialValue,
			testModel.InitialValueDate,
			testModel.CurrentValue,
			testModel.CurrentValueDate,
			testModel.AnnualDepreciationPercent,
			testModel.Status,
			testModel.Created,
			testModel.CreatedBy,
			testModel.Updated,
			testModel.UpdatedBy,
			testModel.Deleted,
			testModel.DeletedBy,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	for _, valueModel := range testModel.Values {
		mock.
			ExpectPrepare(vehicleValuesStmtInsert).
			ExpectExec().
			WithArgs(
				valueModel.ID,
				valueModel.VehicleID,
				valueModel.Date,
				valueModel.Value,
				valueModel.Created,
				valueModel.CreatedBy,
				valueModel.Updated,
				valueModel.UpdatedBy,
				valueModel.Deleted,
				valueModel.DeletedBy,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

	}

	mock.ExpectCommit()

	err := t.repo.Create(testModel)

	assert.NoError(t.T(), err)
}
