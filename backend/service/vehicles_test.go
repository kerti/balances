package service_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	mock_repository "github.com/kerti/balances/backend/mock/repository"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/nuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type vehiclesServiceTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	svc                service.Vehicle
	mockRepo           *mock_repository.MockVehicle
	testUserID         uuid.UUID
	testVehicleID      uuid.UUID
	testVehicleValueID uuid.UUID
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
	t.testVehicleValueID, _ = uuid.NewV7()
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

func (t *vehiclesServiceTestSuite) getVehicleSlice(count int) (res []model.Vehicle) {
	for range count {
		id, _ := uuid.NewV7()
		res = append(res, t.getNewVehicle(nuuid.From(id), nil))
	}
	return
}

func (t *vehiclesServiceTestSuite) getNewVehicle(id nuuid.NUUID, values *[]model.VehicleValue) model.Vehicle {
	veh := model.Vehicle{}

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
	veh.PurchaseDate = initialValueDate
	veh.InitialValue = 68000
	veh.InitialValueDate = initialValueDate
	veh.CurrentValue = 50000
	veh.CurrentValueDate = time.Now()
	veh.AnnualDepreciationPercent = 3.5
	veh.Status = model.VehicleStatusInUse

	veh.Values = []model.VehicleValue{}
	if values != nil {
		for _, val := range *values {
			valCopy := val
			veh.Values = append(veh.Values, valCopy)
		}
	}

	return veh
}

func (t *vehiclesServiceTestSuite) getNewVehicleValue(id nuuid.NUUID, vehicleID nuuid.NUUID, value float64, date time.Time) model.VehicleValue {
	val := model.VehicleValue{}

	if id.Valid {
		val.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		val.ID = newID
	}

	if vehicleID.Valid {
		val.VehicleID = vehicleID.UUID
	} else {
		newVehID, _ := uuid.NewV7()
		val.VehicleID = newVehID
	}

	val.Value = value
	val.Date = date

	return val
}

func (t *vehiclesServiceTestSuite) getNewVehicleValueInput(id nuuid.NUUID, vehicleID nuuid.NUUID, value float64, date time.Time) model.VehicleValueInput {
	val := model.VehicleValueInput{}

	if id.Valid {
		val.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		val.ID = newID
	}

	if vehicleID.Valid {
		val.VehicleID = vehicleID.UUID
	} else {
		newVehicleID, _ := uuid.NewV7()
		val.VehicleID = newVehicleID
	}

	val.Value = value
	val.Date = cachetime.CacheTime(date)

	return val
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
	t.mockRepo.EXPECT().Create(gomock.Any()).Return(failure.InternalError("create", "Vehicle", errors.New(errMsg)))

	actual, err := t.svc.Create(testInput, t.testUserID)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), "Vehicle", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, errMsg)
	assert.Equal(t.T(), "create", *errAsFailure.Operation)
}

func (t *vehiclesServiceTestSuite) TestGetByID_Exists_NoValue() {
	vehicle := t.getNewVehicle(nuuid.NUUID{}, nil)
	resolvedVehicleSlice := []model.Vehicle{vehicle}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(resolvedVehicleSlice, nil)

	_, err := t.svc.GetByID(t.testVehicleID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestGetByID_Exists_WithValue_NoFilter() {
	valueFilterInput := model.VehicleValueFilterInput{
		VehicleIDs: &[]uuid.UUID{t.testVehicleID},
	}
	pageInfo := getDefaultPageInfo()

	vehicle := t.getNewVehicle(nuuid.NUUID{}, nil)
	value1 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(900), time.Now())
	valueSlice := []model.VehicleValue{value1, value2}
	vehicle.AttachValues(valueSlice, true)

	resolvedVehicleSlice := []model.Vehicle{vehicle}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(resolvedVehicleSlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testVehicleID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestGetByID_Exists_WithValue_WithStartDate() {
	yesterday := cachetime.NCacheTime(null.TimeFrom(time.Now().AddDate(0, 0, -1)))
	valueFilterInput := model.VehicleValueFilterInput{
		VehicleIDs: &[]uuid.UUID{t.testVehicleID},
		StartDate:  yesterday,
	}
	pageInfo := getDefaultPageInfo()

	vehicle := t.getNewVehicle(nuuid.NUUID{}, nil)
	value1 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(900), time.Now())
	valueSlice := []model.VehicleValue{value1, value2}
	vehicle.AttachValues(valueSlice, true)

	resolvedVehicleSlice := []model.Vehicle{vehicle}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(resolvedVehicleSlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testVehicleID, true, yesterday, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestGetByID_Exists_WithValue_WithEndDate() {
	today := cachetime.NCacheTime(null.TimeFrom(time.Now()))
	valueFilterInput := model.VehicleValueFilterInput{
		VehicleIDs: &[]uuid.UUID{t.testVehicleID},
		EndDate:    today,
	}
	pageInfo := getDefaultPageInfo()

	vehicle := t.getNewVehicle(nuuid.NUUID{}, nil)
	value1 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(900), time.Now())
	valueSlice := []model.VehicleValue{value1, value2}
	vehicle.AttachValues(valueSlice, true)

	resolvedVehicleSlice := []model.Vehicle{vehicle}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(resolvedVehicleSlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testVehicleID, true, cachetime.NCacheTime{}, today, nil)

	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestGetByID_Exists_WithValue_WithBothDates() {
	yesterday := cachetime.NCacheTime(null.TimeFrom(time.Now().AddDate(0, 0, -1)))
	today := cachetime.NCacheTime(null.TimeFrom(time.Now()))
	valueFilterInput := model.VehicleValueFilterInput{
		VehicleIDs: &[]uuid.UUID{t.testVehicleID},
		StartDate:  yesterday,
		EndDate:    today,
	}
	pageInfo := getDefaultPageInfo()

	vehicle := t.getNewVehicle(nuuid.NUUID{}, nil)
	value1 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(900), time.Now())
	valueSlice := []model.VehicleValue{value1, value2}
	vehicle.AttachValues(valueSlice, true)

	resolvedVehicleSlice := []model.Vehicle{vehicle}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(resolvedVehicleSlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testVehicleID, true, yesterday, today, nil)

	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestGetByID_Exists_WithValue_WithPageSize() {
	pageSize := 120
	valueFilterInput := model.VehicleValueFilterInput{
		VehicleIDs: &[]uuid.UUID{t.testVehicleID},
	}
	valueFilterInput.PageSize = &pageSize
	pageInfo := model.PageInfoOutput{
		PageSize: pageSize,
	}

	vehicle := t.getNewVehicle(nuuid.NUUID{}, nil)
	value1 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(900), time.Now())
	valueSlice := []model.VehicleValue{value1, value2}
	vehicle.AttachValues(valueSlice, true)

	resolvedVehicleSlice := []model.Vehicle{vehicle}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(resolvedVehicleSlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testVehicleID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, &pageSize)

	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestGetByID_Exists_With_Value_RepoErrorResolvingVehicle() {
	errMsg := "failed to resolve vehicle by IDs"
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, failure.InternalError("resolve by IDs", "Vehicle", errors.New(errMsg)))

	actual, err := t.svc.GetByID(t.testVehicleID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), "Vehicle", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, errMsg)
	assert.Equal(t.T(), "resolve by IDs", *errAsFailure.Operation)
}

func (t *vehiclesServiceTestSuite) TestGetByID_Exists_WithValue_RepoErrorResolvingValues() {
	errMsg := "error resolving values"
	valueFilterInput := model.VehicleValueFilterInput{
		VehicleIDs: &[]uuid.UUID{t.testVehicleID},
	}

	vehicle := t.getNewVehicle(nuuid.NUUID{}, nil)
	value1 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewVehicleValue(nuuid.NUUID{}, nuuid.From(vehicle.ID), float64(900), time.Now())
	valueSlice := []model.VehicleValue{value1, value2}
	vehicle.AttachValues(valueSlice, true)
	resolvedVehicleSlice := []model.Vehicle{vehicle}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(resolvedVehicleSlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return([]model.VehicleValue{}, getDefaultPageInfo(), failure.InternalError("resolve by filter", "Vehicle Value", errors.New(errMsg)))

	actual, err := t.svc.GetByID(t.testVehicleID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), "Vehicle Value", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, errMsg)
	assert.Equal(t.T(), "resolve by filter", *errAsFailure.Operation)
}

func (t *vehiclesServiceTestSuite) TestGetByID_NotExists() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, nil)

	actual, err := t.svc.GetByID(t.testVehicleID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), "Vehicle", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, "not found")
	assert.Equal(t.T(), "get by ID", *errAsFailure.Operation)
}

func (t *vehiclesServiceTestSuite) TestGetByFilter_EmptyFilter() {
	filterInput := model.VehicleFilterInput{}
	filter := filterInput.ToFilter()

	t.mockRepo.EXPECT().ResolveByFilter(filter).
		Return(t.getVehicleSlice(2), getDefaultPageInfo(), nil)

	_, _, err := t.svc.GetByFilter(filterInput)

	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestGetByFilter_WithKeyword() {
	keyword := "example"
	filterInput := model.VehicleFilterInput{}
	filterInput.Keyword = &keyword
	filter := filterInput.ToFilter()

	t.mockRepo.EXPECT().ResolveByFilter(filter).
		Return(t.getVehicleSlice(2), getDefaultPageInfo(), nil)

	_, _, err := t.svc.GetByFilter(filterInput)

	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestUpdate_Normal() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{t.getNewVehicle(nuuid.From(t.testVehicleID), nil)}, nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).
		Return(nil)

	res, err := t.svc.Update(t.getNewVehicleInput(nuuid.From(t.testVehicleID)), t.testUserID)

	assert.NotNil(t.T(), res)
	assert.NoError(t.T(), err)
}

func (t *vehiclesServiceTestSuite) TestUpdate_RepoErrorResolvingByIDs() {
	errMsg := "query failed"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, failure.InternalError("resolve by IDs", "Vehicle", errors.New(errMsg)))

	res, err := t.svc.Update(t.getNewVehicleInput(nuuid.From(t.testVehicleID)), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by IDs", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *vehiclesServiceTestSuite) TestUpdate_VehicleNotFound() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, nil)

	res, err := t.svc.Update(t.getNewVehicleInput(nuuid.From(t.testVehicleID)), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeEntityNotFound, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "not found")
}

func (t *vehiclesServiceTestSuite) TestUpdate_VehicleDeleted() {
	vehicleInput := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	deletedVehicle := model.NewVehicleFromInput(vehicleInput, t.testUserID)
	deletedVehicle.Delete(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{deletedVehicle}, nil)

	res, err := t.svc.Update(t.getNewVehicleInput(nuuid.From(t.testVehicleID)), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeOperationNotPermitted, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "already deleted")
}

func (t *vehiclesServiceTestSuite) TestUpdate_RepoErrorUpdating() {
	errMsg := "failed to update"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{t.getNewVehicle(nuuid.From(t.testVehicleID), nil)}, nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).
		Return(failure.InternalError("update", "Vehicle", errors.New(errMsg)))

	res, err := t.svc.Update(t.getNewVehicleInput(nuuid.From(t.testVehicleID)), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Vehicle", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *vehiclesServiceTestSuite) TestDelete_Normal() {
	valuesSlice := []model.VehicleValue{}

	valuesSlice = append(
		valuesSlice,
		t.getNewVehicleValue(
			nuuid.NUUID{},
			nuuid.From(t.testVehicleID),
			float64(10000),
			time.Now().AddDate(0, 0, -1)))

	valuesSlice = append(
		valuesSlice,
		t.getNewVehicleValue(
			nuuid.NUUID{},
			nuuid.From(t.testVehicleID),
			float64(9000),
			time.Now()))

	testVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), &valuesSlice)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testVehicle}, nil)

	t.mockRepo.EXPECT().ResolveValuesByFilter(gomock.Any()).
		Return(valuesSlice, getDefaultPageInfo(), nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

	res, err := t.svc.Delete(t.testVehicleID, t.testUserID)

	assert.NoError(t.T(), err)

	assert.NotNil(t.T(), res)
	assert.True(t.T(), res.Deleted.Valid)
	assert.True(t.T(), res.DeletedBy.Valid)

	assert.Len(t.T(), res.Values, 2)
	for _, resValue := range res.Values {
		assert.True(t.T(), resValue.Deleted.Valid)
		assert.True(t.T(), resValue.DeletedBy.Valid)
	}
}

func (t *vehiclesServiceTestSuite) TestDelete_RepoErrorResolvingByIDs() {
	errMsg := "failed resolving by IDs"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, errors.New(errMsg))

	res, err := t.svc.Delete(t.testVehicleID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDelete_RepoErrorResolvingValuesByFilter() {
	errMsg := "failed resolving vehicle values"

	testVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testVehicle}, nil)

	t.mockRepo.EXPECT().ResolveValuesByFilter(gomock.Any()).
		Return([]model.VehicleValue{}, model.PageInfoOutput{}, errors.New(errMsg))

	res, err := t.svc.Delete(t.testVehicleID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)

	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDelete_VehicleNotFound() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, nil)

	res, err := t.svc.Delete(t.testVehicleID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDelete_VehicleAlreadyDeleted() {
	testDeletedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	testDeletedVehicle.Deleted = null.TimeFrom(time.Now())
	testDeletedVehicle.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testDeletedVehicle}, nil)

	res, err := t.svc.Delete(t.testVehicleID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "delete")
	assert.Contains(t.T(), err.Error(), "Vehicle")
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDelete_VehicleValueAlreadyDeleted() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{t.getNewVehicle(nuuid.From(t.testVehicleID), nil)}, nil)

	testDeletedNonLastVehicleValue := t.getNewVehicleValue(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(10000),
		time.Now().AddDate(0, 0, -1))

	testDeletedNonLastVehicleValue.Deleted = null.TimeFrom(time.Now())
	testDeletedNonLastVehicleValue.DeletedBy = nuuid.From(t.testUserID)

	testLastVehicleValue := t.getNewVehicleValue(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(12000),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByFilter(gomock.Any()).
		Return(
			[]model.VehicleValue{
				testDeletedNonLastVehicleValue,
				testLastVehicleValue,
			},
			getDefaultPageInfo(),
			nil)

	res, err := t.svc.Delete(t.testVehicleID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "delete")
	assert.Contains(t.T(), err.Error(), "Vehicle Value")
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDelete_RepoErrorUpdating() {
	errMsg := "failed updating vehicle"
	valueSlice := []model.VehicleValue{}

	valueSlice = append(
		valueSlice,
		t.getNewVehicleValue(
			nuuid.NUUID{},
			nuuid.From(t.testVehicleID),
			float64(10000),
			time.Now().AddDate(0, 0, -1)))

	valueSlice = append(
		valueSlice,
		t.getNewVehicleValue(
			nuuid.NUUID{},
			nuuid.From(t.testVehicleID),
			float64(12000),
			time.Now()))

	testVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), &valueSlice)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testVehicle}, nil)

	t.mockRepo.EXPECT().ResolveValuesByFilter(gomock.Any()).
		Return(valueSlice, getDefaultPageInfo(), nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).Return(errors.New(errMsg))

	res, err := t.svc.Delete(t.testVehicleID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)

	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_Normal_CurrentValue() {
	testValueDate := time.Now()
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)
	testValue := t.getNewVehicleValue(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)

	testVehicleToUpdate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	testVehicleToUpdate.CurrentValue = float64(900)
	testVehicleToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	testVehicleAfterUpate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	testVehicleAfterUpate.PurchaseDate = testVehicleToUpdate.PurchaseDate
	testVehicleAfterUpate.InitialValueDate = testVehicleToUpdate.InitialValueDate
	testVehicleAfterUpate.CurrentValue = testValue.Value
	testVehicleAfterUpate.CurrentValueDate = testValue.Date

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testVehicleToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return(
			[]model.VehicleValue{
				t.getNewVehicleValue(
					nuuid.NUUID{},
					nuuid.From(t.testVehicleID),
					float64(900),
					time.Now().AddDate(0, 0, -1)),
			},
			nil)

	t.mockRepo.EXPECT().CreateValue(
		vehicleValueMatcher{testValue},
		vehiclePointerMatcher{testVehicleAfterUpate}).
		Return(nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_Normal_NotCurrentValue() {
	testValueDate := time.Now().AddDate(0, 0, -12)
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)
	testValue := t.getNewVehicleValue(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)

	testVehicleToUpdate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	testVehicleToUpdate.CurrentValue = float64(900)
	testVehicleToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testVehicleToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return(
			[]model.VehicleValue{
				t.getNewVehicleValue(
					nuuid.NUUID{},
					nuuid.From(t.testVehicleID),
					float64(900),
					time.Now().AddDate(0, 0, -1)),
			},
			nil)

	t.mockRepo.EXPECT().
		CreateValue(vehicleValueMatcher{testValue}, nil).
		Return(nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_RepoFailedResolvingByIDs() {
	errMsg := "failed to resolve vehicle by IDs"
	testValueDate := time.Now()
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, errors.New(errMsg))

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_VehicleNotFound() {
	testValueDate := time.Now()
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_VehicleDeleted() {
	testValueDate := time.Now()
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)

	deletedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	deletedVehicle.Deleted = null.TimeFrom(time.Now())
	deletedVehicle.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{deletedVehicle}, nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_VehicleSold() {
	testValueDate := time.Now()
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)

	deletedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	deletedVehicle.Status = model.VehicleStatusSold

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{deletedVehicle}, nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "sold")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_RepoFailedResolvingCurrentValue() {
	errMsg := "failed resolving vehicle values by vehicle ID"
	testValueDate := time.Now().AddDate(0, 0, -12)
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)

	testVehicleToUpdate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	testVehicleToUpdate.CurrentValue = float64(900)
	testVehicleToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testVehicleToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return([]model.VehicleValue{}, errors.New(errMsg))

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_CurrentValueNotFound() {
	testValueDate := time.Now().AddDate(0, 0, -12)
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)

	testVehicleToUpdate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	testVehicleToUpdate.CurrentValue = float64(900)
	testVehicleToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testVehicleToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return([]model.VehicleValue{}, nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestCreateValue_RepoFailedCreatingVehicleValue() {
	errMsg := "failed creating vehicle value"
	testValueDate := time.Now().AddDate(0, 0, -12)
	testInput := t.getNewVehicleValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate)
	testValue := t.getNewVehicleValue(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(1234.56),
		testValueDate,
	)

	testVehicleToUpdate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	testVehicleToUpdate.CurrentValue = float64(900)
	testVehicleToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{testVehicleToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return(
			[]model.VehicleValue{
				t.getNewVehicleValue(
					nuuid.NUUID{},
					nuuid.From(t.testVehicleID),
					float64(900),
					time.Now().AddDate(0, 0, -1))},
			nil)

	t.mockRepo.EXPECT().CreateValue(vehicleValueMatcher{testValue}, nil).
		Return(errors.New(errMsg))

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestGetValueByID_Normal() {
	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return(
			[]model.VehicleValue{
				t.getNewVehicleValue(
					nuuid.From(t.testVehicleValueID),
					nuuid.From(t.testVehicleID),
					float64(1000000), time.Now(),
				),
			},
			nil)

	res, err := t.svc.GetValueByID(t.testVehicleValueID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestGetValueByID_RepoFailedResolvingValue() {
	errMsg := "failed to resolve vehicle values"
	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return(
			[]model.VehicleValue{},
			failure.InternalError("resolve by filter", "Vehicle Value", errors.New(errMsg)))

	actual, err := t.svc.GetValueByID(t.testVehicleValueID)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, errAsFailure.Code)
	assert.Equal(t.T(), "Vehicle Value", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, errMsg)
	assert.Equal(t.T(), "resolve by filter", *errAsFailure.Operation)
}

func (t *vehiclesServiceTestSuite) TestGetValueByID_ValueNotFound() {
	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return(
			[]model.VehicleValue{},
			nil)

	actual, err := t.svc.GetValueByID(t.testVehicleValueID)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), failure.CodeEntityNotFound, errAsFailure.Code)
	assert.Equal(t.T(), "Vehicle Value", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, "not found")
	assert.Equal(t.T(), "get by ID", *errAsFailure.Operation)
}

func (t *vehiclesServiceTestSuite) TestGetValueBVyFilter_Normal() {
	filter := model.VehicleValueFilterInput{}
	t.mockRepo.EXPECT().ResolveValuesByFilter(filter.ToFilter()).
		Return(
			[]model.VehicleValue{
				t.getNewVehicleValue(
					nuuid.From(t.testVehicleValueID),
					nuuid.From(t.testVehicleValueID),
					float64(1000000),
					time.Now())},
			getDefaultPageInfo(),
			nil,
		)

	res, pageInfo, err := t.svc.GetValuesByFilter(filter)

	assert.NoError(t.T(), err)
	assert.Len(t.T(), res, 1)
	assert.Equal(t.T(), pageInfo.TotalCount, 1)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_Normal_CurrentValue() {
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	vehicleToUpdate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	vehicleToUpdate.CurrentValue = float64(900)
	vehicleToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	valueToUpdate := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		vehicleToUpdate.CurrentValue,
		vehicleToUpdate.CurrentValueDate)

	updatedValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		testInput.Value,
		testInput.Date.Time())

	updatedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	updatedVehicle.CurrentValue = updatedValue.Value
	updatedVehicle.CurrentValueDate = updatedValue.Date

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{vehicleToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return([]model.VehicleValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		vehicleValueMatcher{updatedValue},
		vehiclePointerMatcher{updatedVehicle}).
		Return(nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_Normal_NonCurrentValue() {
	newValueID, _ := uuid.NewV7()
	testInput := t.getNewVehicleValueInput(
		nuuid.From(newValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now().AddDate(0, 0, -2),
	)

	vehicleToUpdate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	vehicleToUpdate.CurrentValue = float64(900)
	vehicleToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	valueToUpdate := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		vehicleToUpdate.CurrentValue,
		vehicleToUpdate.CurrentValueDate)

	updatedValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		testInput.Value,
		testInput.Date.Time())

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{vehicleToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{newValueID}).
		Return([]model.VehicleValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return([]model.VehicleValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		vehicleValueMatcher{updatedValue},
		nil).
		Return(nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_RepoFailedResolvingByIDs() {
	errMsg := "failed resolving vehicles by IDs"
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, errors.New(errMsg))

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_VehicleNotFound() {
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Vehicle")
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_VehicleDeleted() {
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.CurrentValue = float64(900)
	resolvedVehicle.CurrentValueDate = time.Now().AddDate(0, 0, -1)
	resolvedVehicle.Deleted = null.TimeFrom(time.Now())
	resolvedVehicle.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{resolvedVehicle}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_VehicleSold() {
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.CurrentValue = float64(900)
	resolvedVehicle.CurrentValueDate = time.Now().AddDate(0, 0, -1)
	resolvedVehicle.Status = model.VehicleStatusSold

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{resolvedVehicle}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "sold")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_RepoFailedResolvingValues() {
	errMsg := "failed resolving vehicle values"
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.CurrentValue = float64(900)
	resolvedVehicle.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{resolvedVehicle}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{}, errors.New(errMsg))

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_ValueNotFound() {
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.CurrentValue = float64(900)
	resolvedVehicle.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{resolvedVehicle}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Vehicle Value")
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_ValueDeleted() {
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.CurrentValue = float64(900)
	resolvedVehicle.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	resolvedValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		resolvedVehicle.CurrentValue,
		resolvedVehicle.CurrentValueDate)
	resolvedValue.Deleted = null.TimeFrom(time.Now())
	resolvedValue.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{resolvedVehicle}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{resolvedValue}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Vehicle Value")
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_RepoFailedResolvingCurrentValues() {
	errMsg := "failed resolving vehicle last value"
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.CurrentValue = float64(900)
	resolvedVehicle.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	resolvedVehicleValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		resolvedVehicle.CurrentValue,
		resolvedVehicle.CurrentValueDate)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{resolvedVehicle}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{resolvedVehicleValue}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return([]model.VehicleValue{}, errors.New(errMsg))

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_CurrentValueNotFound() {
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.CurrentValue = float64(900)
	resolvedVehicle.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	resolvedVehicleValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		resolvedVehicle.CurrentValue,
		resolvedVehicle.CurrentValueDate)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{resolvedVehicle}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{resolvedVehicleValue}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return([]model.VehicleValue{}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Vehicle Value")
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestUpdateValue_RepoFailedUpdatingValue() {
	errMsg := "failed updating vehicle value"
	testInput := t.getNewVehicleValueInput(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(1000),
		time.Now(),
	)

	vehicleToUpdate := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	vehicleToUpdate.CurrentValue = float64(900)
	vehicleToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	valueToUpdate := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		vehicleToUpdate.CurrentValue,
		vehicleToUpdate.CurrentValueDate)

	updatedValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		testInput.Value,
		testInput.Date.Time())

	updatedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	updatedVehicle.PurchaseDate = vehicleToUpdate.PurchaseDate
	updatedVehicle.InitialValueDate = vehicleToUpdate.InitialValueDate
	updatedVehicle.CurrentValue = updatedValue.Value
	updatedVehicle.CurrentValueDate = updatedValue.Date

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{vehicleToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 1).
		Return([]model.VehicleValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		vehicleValueMatcher{updatedValue},
		vehiclePointerMatcher{updatedVehicle}).
		Return(errors.New(errMsg))

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_Normal_CurrentValue() {

	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	secondToCurrentValue := t.getNewVehicleValue(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(456),
		time.Now().AddDate(0, 0, -12))

	deletedCurrentValue := t.getNewVehicleValue(
		nuuid.From(lastValue.ID),
		nuuid.From(lastValue.VehicleID),
		lastValue.Value,
		lastValue.Date)
	deletedCurrentValue.Deleted = null.TimeFrom(time.Now())
	deletedCurrentValue.DeletedBy = nuuid.From(t.testUserID)

	vehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)

	updatedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	updatedVehicle.PurchaseDate = vehicle.PurchaseDate
	updatedVehicle.InitialValueDate = vehicle.InitialValueDate
	updatedVehicle.CurrentValue = secondToCurrentValue.Value
	updatedVehicle.CurrentValueDate = secondToCurrentValue.Date

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{vehicle}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 2).
		Return([]model.VehicleValue{lastValue, secondToCurrentValue}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		vehicleValueMatcher{deletedCurrentValue},
		vehiclePointerMatcher{updatedVehicle}).
		Return(nil)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_Normal_NonCurrentValue() {

	lastValue := t.getNewVehicleValue(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	secondToCurrentValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(456),
		time.Now().AddDate(0, 0, -12))

	deletedNonCurrentValue := t.getNewVehicleValue(
		nuuid.From(secondToCurrentValue.ID),
		nuuid.From(secondToCurrentValue.VehicleID),
		secondToCurrentValue.Value,
		secondToCurrentValue.Date)
	deletedNonCurrentValue.Deleted = null.TimeFrom(time.Now())
	deletedNonCurrentValue.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{secondToCurrentValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(
			[]model.Vehicle{t.getNewVehicle(nuuid.From(t.testVehicleID), nil)},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 2).
		Return([]model.VehicleValue{lastValue, secondToCurrentValue}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		vehicleValueMatcher{deletedNonCurrentValue},
		nil).
		Return(nil)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_RepoFailedResolvingValueByIDs() {
	errMsg := "failed resolving vehicle values by IDs"

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{}, errors.New(errMsg))

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_ValueNotFound() {
	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{}, nil)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_ValueAlreadyDeleted() {

	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())
	lastValue.Deleted = null.TimeFrom(time.Now())
	lastValue.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Vehicle Value")
	assert.Contains(t.T(), err.Error(), "already deleted")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_RepoFailedResolvingByIDs() {
	errMsg := "failed resolving vehicles by IDs"
	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, errors.New(errMsg))

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_VehicleNotFound() {
	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return([]model.Vehicle{}, nil)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_VehicleDeleted() {

	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.Deleted = null.TimeFrom(time.Now())
	resolvedVehicle.DeletedBy = nuuid.From(t.testVehicleID)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(
			[]model.Vehicle{resolvedVehicle},
			nil,
		)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Vehicle")
	assert.Contains(t.T(), err.Error(), "already deleted")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_VehicleSold() {

	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	resolvedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	resolvedVehicle.Status = model.VehicleStatusSold

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(
			[]model.Vehicle{resolvedVehicle},
			nil,
		)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Vehicle")
	assert.Contains(t.T(), err.Error(), "sold")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_RepoFailedResolvingLastValues() {
	errMsg := "failed resolving vehicle last values"
	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(
			[]model.Vehicle{t.getNewVehicle(nuuid.From(t.testVehicleID), nil)},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 2).
		Return([]model.VehicleValue{}, errors.New(errMsg))

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_CurrentValueNotFound() {

	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(
			[]model.Vehicle{t.getNewVehicle(nuuid.From(t.testVehicleID), nil)},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 2).
		Return([]model.VehicleValue{}, nil)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Vehicle Current Value")
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_CannotDeleteTheOnlyValue() {

	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(
			[]model.Vehicle{t.getNewVehicle(nuuid.From(t.testVehicleID), nil)},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 2).
		Return([]model.VehicleValue{lastValue}, nil)

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "OperationNotPermitted")
	assert.Contains(t.T(), err.Error(), "delete")
	assert.Contains(t.T(), err.Error(), "Vehicle Value")
	assert.Contains(t.T(), err.Error(), "the only")
	assert.Nil(t.T(), res)
}

func (t *vehiclesServiceTestSuite) TestDeleteValue_RepoFailedUpdatingValue() {
	errMsg := "failed updating vehicle value"
	lastValue := t.getNewVehicleValue(
		nuuid.From(t.testVehicleValueID),
		nuuid.From(t.testVehicleID),
		float64(123),
		time.Now())

	secondToCurrentValue := t.getNewVehicleValue(
		nuuid.NUUID{},
		nuuid.From(t.testVehicleID),
		float64(456),
		time.Now().AddDate(0, 0, -12))

	deletedCurrentValue := t.getNewVehicleValue(
		nuuid.From(lastValue.ID),
		nuuid.From(lastValue.VehicleID),
		lastValue.Value,
		lastValue.Date)
	deletedCurrentValue.Deleted = null.TimeFrom(time.Now())
	deletedCurrentValue.DeletedBy = nuuid.From(t.testUserID)

	vehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)

	updatedVehicle := t.getNewVehicle(nuuid.From(t.testVehicleID), nil)
	updatedVehicle.PurchaseDate = vehicle.PurchaseDate
	updatedVehicle.InitialValueDate = vehicle.InitialValueDate
	updatedVehicle.CurrentValue = secondToCurrentValue.Value
	updatedVehicle.CurrentValueDate = secondToCurrentValue.Date

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testVehicleValueID}).
		Return([]model.VehicleValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testVehicleID}).
		Return(
			[]model.Vehicle{vehicle},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByVehicleID(t.testVehicleID, 2).
		Return([]model.VehicleValue{lastValue, secondToCurrentValue}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		vehicleValueMatcher{deletedCurrentValue},
		vehiclePointerMatcher{updatedVehicle}).
		Return(errors.New(errMsg))

	res, err := t.svc.DeleteValue(t.testVehicleValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

// matchers

type vehiclePointerMatcher struct {
	expected model.Vehicle
}

func (m vehiclePointerMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*model.Vehicle)
	if !ok {
		return false
	}

	return actual.Name == m.expected.Name &&
		actual.Make == m.expected.Make &&
		actual.Model == m.expected.Model &&
		actual.Year == m.expected.Year &&
		actual.Type == m.expected.Type &&
		actual.TitleHolder == m.expected.TitleHolder &&
		actual.LicensePlateNumber == m.expected.LicensePlateNumber &&
		actual.PurchaseDate.Equal(m.expected.PurchaseDate) &&
		actual.InitialValue == m.expected.InitialValue &&
		actual.InitialValueDate.Equal(m.expected.InitialValueDate) &&
		actual.CurrentValue == m.expected.CurrentValue &&
		actual.CurrentValueDate.Equal(m.expected.CurrentValueDate) &&
		actual.AnnualDepreciationPercent == m.expected.AnnualDepreciationPercent &&
		actual.Status == m.expected.Status
}

func (m vehiclePointerMatcher) String() string {
	return fmt.Sprintf(
		"is Vehicle with Name=%s, Make=%s, Model=%s, Year=%d, Type=%s, TitleHolder=%s, LicensePlateNumber=%s, PurchaseDate=%s, InitialValue=%f, InitialValueDate=%s, CurrentValue=%f, CurrentValueDate=%s, AnnualDepreciationPercent=%f, Status=%s",
		m.expected.Name,
		m.expected.Make,
		m.expected.Model,
		m.expected.Year,
		m.expected.Type,
		m.expected.TitleHolder,
		m.expected.LicensePlateNumber,
		m.expected.PurchaseDate,
		m.expected.InitialValue,
		m.expected.InitialValueDate,
		m.expected.CurrentValue,
		m.expected.CurrentValueDate,
		m.expected.AnnualDepreciationPercent,
		m.expected.Status)
}

type vehicleValueMatcher struct {
	expected model.VehicleValue
}

func (m vehicleValueMatcher) Matches(x interface{}) bool {
	actual, ok := x.(model.VehicleValue)
	if !ok {
		return false
	}

	return actual.Date.Equal(m.expected.Date) &&
		actual.Value == m.expected.Value &&
		actual.VehicleID == m.expected.VehicleID
}

func (m vehicleValueMatcher) String() string {
	return fmt.Sprintf(
		"is VehicleValue with Value=%.2f, Date=%v, VehicleID=%v",
		m.expected.Value,
		m.expected.Date,
		m.expected.VehicleID)
}
