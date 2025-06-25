package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mock_repository "github.com/kerti/balances/backend/mock/repository"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/nuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type vehiclesServiceTestSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	svc           service.Vehicle
	mockRepo      *mock_repository.MockVehicle
	testUserID    uuid.UUID
	testVehicleID uuid.UUID
}

func TestVehiclesService(t *testing.T) {
	suite.Run(t, new(vehiclesServiceTestSuite))
}

func (t *vehiclesServiceTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockRepo = mock_repository.NewMockVehicle(t.ctrl)
	t.svc = &service.VehicleImpl{
		Repository: t.mockRepo,
	}
	t.testUserID, _ = uuid.NewV7()
	t.testVehicleID, _ = uuid.NewV7()
	t.svc.Startup()
}

func (t *vehiclesServiceTestSuite) TearDownTest() {
	t.svc.Shutdown()
	t.ctrl.Finish()
}

func (t *vehiclesServiceTestSuite) getNewVehicleInput(id nuuid.NUUID) model.VehicleInput {
	veh := model.VehicleInput{}

	if id.Valid {
		veh.ID = id.UUID
	} else {
		veh.ID = t.testVehicleID
	}

	initialValueDate := time.Now().AddDate(-2, 0, 0) // defaults to 2 years ago

	veh.Name = "John's Car"
	veh.Make = "Hyundai"
	veh.Model = "Palisade"
	veh.Year = 2020
	veh.Type = model.VehicleTypeCar
	veh.TitleHolder = "John Fitzgerald Doe"
	veh.LicensePlateNumber = "TUNEMAN"
	veh.PurchaseDate = cachetime.CacheTime(initialValueDate)
	veh.InitialValue = 68000
	veh.InitialValueDate = cachetime.CacheTime(initialValueDate)
	veh.CurrentValue = 50000
	veh.CurrentValueDate = cachetime.CacheTime(time.Now())
	veh.AnnualDepreciationPercent = 3.5
	veh.Status = model.VehicleStatusInUse

	return veh
}

func (t *vehiclesServiceTestSuite) TestCreate_Normal() {
	testInput := t.getNewVehicleInput(nuuid.NUUID{Valid: false})
	t.mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	res, err := t.svc.Create(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.Equal(t.T(), testInput.Name, res.Name)
	assert.Equal(t.T(), testInput.Make, res.Make)
	assert.Equal(t.T(), testInput.Model, res.Model)
	assert.Equal(t.T(), testInput.Year, res.Year)
	assert.Equal(t.T(), testInput.Type, res.Type)
	assert.Equal(t.T(), testInput.TitleHolder, res.TitleHolder)
	assert.Equal(t.T(), testInput.LicensePlateNumber, res.LicensePlateNumber)
	assert.Equal(t.T(), testInput.PurchaseDate.Time().Unix(), res.PurchaseDate.Unix())
	assert.Equal(t.T(), testInput.InitialValue, res.InitialValue)
	assert.Equal(t.T(), testInput.InitialValueDate.Time().Unix(), res.InitialValueDate.Unix())
	assert.Equal(t.T(), testInput.CurrentValue, res.CurrentValue)
	assert.Equal(t.T(), testInput.CurrentValueDate.Time().Unix(), res.CurrentValueDate.Unix())
	assert.Equal(t.T(), testInput.AnnualDepreciationPercent, res.AnnualDepreciationPercent)
	assert.Equal(t.T(), testInput.Status, res.Status)
}

func (t *vehiclesServiceTestSuite) TestCreate_RepoFailToCreate() {
	errMsg := "repo failed to create vehicle"
	testInput := t.getNewVehicleInput(nuuid.NUUID{Valid: false})
	t.mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New(errMsg))

	res, err := t.svc.Create(testInput, t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
}
