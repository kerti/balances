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

type propertiesServiceTestSuite struct {
	suite.Suite
	ctrl                *gomock.Controller
	svc                 service.Property
	mockRepo            *mock_repository.MockProperty
	testUserID          uuid.UUID
	testPropertyID      uuid.UUID
	testPropertyValueID uuid.UUID
}

func TestPropertiesService(t *testing.T) {
	suite.Run(t, new(propertiesServiceTestSuite))
}

func (t *propertiesServiceTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockRepo = mock_repository.NewMockProperty(t.ctrl)
	t.svc = &service.PropertyImpl{
		Repository: t.mockRepo,
	}
	t.testUserID, _ = uuid.NewV7()
	t.testPropertyID, _ = uuid.NewV7()
	t.testPropertyValueID, _ = uuid.NewV7()
	t.svc.Startup()
}

func (t *propertiesServiceTestSuite) TearDownTest() {
	t.svc.Shutdown()
	t.ctrl.Finish()
}

func (t *propertiesServiceTestSuite) getNewPropertyInput(id nuuid.NUUID) model.PropertyInput {
	prop := model.PropertyInput{}

	if id.Valid {
		prop.ID = id.UUID
	} else {
		prop.ID = t.testPropertyID
	}

	initialValueDate := time.Now().AddDate(-2, 0, 0) // defaults to 2 years ago

	prop.Name = "John's House"
	prop.Address = "1234 Main Street"
	prop.TotalArea = 1200
	prop.BuildingArea = 1000
	prop.AreaUnit = model.PropertyAreaUnitSQM
	prop.Type = model.PropertyTypeHouse
	prop.TitleHolder = "John Fitzgerald Doe"
	prop.TaxIdentifier = "TAX-12345"
	prop.PurchaseDate = cachetime.CacheTime(initialValueDate)
	prop.InitialValue = 68000
	prop.InitialValueDate = cachetime.CacheTime(initialValueDate)
	prop.CurrentValue = 50000
	prop.CurrentValueDate = cachetime.CacheTime(time.Now())
	prop.AnnualAppreciationPercent = 3.5
	prop.Status = model.PropertyStatusInUse

	return prop
}

func (t *propertiesServiceTestSuite) getPropertySlice(count int) (res []model.Property) {
	for range count {
		id, _ := uuid.NewV7()
		res = append(res, t.getNewProperty(nuuid.From(id), nil))
	}
	return
}

func (t *propertiesServiceTestSuite) getNewProperty(id nuuid.NUUID, values *[]model.PropertyValue) model.Property {
	prop := model.Property{}

	if id.Valid {
		prop.ID = id.UUID
	} else {
		prop.ID = t.testPropertyID
	}

	initialValueDate := time.Now().AddDate(-2, 0, 0) // defaults to 2 years ago

	prop.Name = "John's House"
	prop.Address = "1234 Main Street"
	prop.TotalArea = 1200
	prop.BuildingArea = 1000
	prop.AreaUnit = model.PropertyAreaUnitSQM
	prop.Type = model.PropertyTypeHouse
	prop.TitleHolder = "John Fitzgerald Doe"
	prop.TaxIdentifier = "TAX-12345"
	prop.PurchaseDate = initialValueDate
	prop.InitialValue = 68000
	prop.InitialValueDate = initialValueDate
	prop.CurrentValue = 50000
	prop.CurrentValueDate = time.Now()
	prop.AnnualAppreciationPercent = 3.5
	prop.Status = model.PropertyStatusInUse

	prop.Values = []model.PropertyValue{}
	if values != nil {
		for _, val := range *values {
			valCopy := val
			prop.Values = append(prop.Values, valCopy)
		}
	}

	return prop
}

func (t *propertiesServiceTestSuite) getNewPropertyValue(id nuuid.NUUID, propertyID nuuid.NUUID, value float64, date time.Time) model.PropertyValue {
	val := model.PropertyValue{}

	if id.Valid {
		val.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		val.ID = newID
	}

	if propertyID.Valid {
		val.PropertyID = propertyID.UUID
	} else {
		newPropID, _ := uuid.NewV7()
		val.PropertyID = newPropID
	}

	val.Value = value
	val.Date = date

	return val
}

func (t *propertiesServiceTestSuite) getNewPropertyValueInput(id nuuid.NUUID, propertyID nuuid.NUUID, value float64, date time.Time) model.PropertyValueInput {
	val := model.PropertyValueInput{}

	if id.Valid {
		val.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		val.ID = newID
	}

	if propertyID.Valid {
		val.PropertyID = propertyID.UUID
	} else {
		newPropertyID, _ := uuid.NewV7()
		val.PropertyID = newPropertyID
	}

	val.Value = value
	val.Date = cachetime.CacheTime(date)

	return val
}

func (t *propertiesServiceTestSuite) TestCreate_Normal() {
	testInput := t.getNewPropertyInput(nuuid.NUUID{Valid: false})
	t.mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	res, err := t.svc.Create(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.Equal(t.T(), testInput.Name, res.Name)
	assert.Equal(t.T(), testInput.Address, res.Address)
	assert.Equal(t.T(), testInput.TotalArea, res.TotalArea)
	assert.Equal(t.T(), testInput.BuildingArea, res.BuildingArea)
	assert.Equal(t.T(), testInput.AreaUnit, res.AreaUnit)
	assert.Equal(t.T(), testInput.Type, res.Type)
	assert.Equal(t.T(), testInput.TitleHolder, res.TitleHolder)
	assert.Equal(t.T(), testInput.TaxIdentifier, res.TaxIdentifier)
	assert.Equal(t.T(), testInput.PurchaseDate.Time().Unix(), res.PurchaseDate.Unix())
	assert.Equal(t.T(), testInput.InitialValue, res.InitialValue)
	assert.Equal(t.T(), testInput.InitialValueDate.Time().Unix(), res.InitialValueDate.Unix())
	assert.Equal(t.T(), testInput.CurrentValue, res.CurrentValue)
	assert.Equal(t.T(), testInput.CurrentValueDate.Time().Unix(), res.CurrentValueDate.Unix())
	assert.Equal(t.T(), testInput.AnnualAppreciationPercent, res.AnnualAppreciationPercent)
	assert.Equal(t.T(), testInput.Status, res.Status)
}

func (t *propertiesServiceTestSuite) TestCreate_RepoFailToCreate() {
	errMsg := "repo failed to create property"
	testInput := t.getNewPropertyInput(nuuid.NUUID{Valid: false})
	t.mockRepo.EXPECT().Create(gomock.Any()).Return(failure.InternalError("create", "Property", errors.New(errMsg)))

	actual, err := t.svc.Create(testInput, t.testUserID)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), "Property", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, errMsg)
	assert.Equal(t.T(), "create", *errAsFailure.Operation)
}

func (t *propertiesServiceTestSuite) TestGetByID_Exists_NoValue() {
	property := t.getNewProperty(nuuid.NUUID{}, nil)
	resolvedPropertySlice := []model.Property{property}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(resolvedPropertySlice, nil)

	_, err := t.svc.GetByID(t.testPropertyID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestGetByID_Exists_WithValue_NoFilter() {
	valueFilterInput := model.PropertyValueFilterInput{
		PropertyIDs: &[]uuid.UUID{t.testPropertyID},
	}
	pageInfo := getDefaultPageInfo()

	property := t.getNewProperty(nuuid.NUUID{}, nil)
	value1 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(900), time.Now())
	valueSlice := []model.PropertyValue{value1, value2}
	property.AttachValues(valueSlice, true)

	resolvedPropertySlice := []model.Property{property}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(resolvedPropertySlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testPropertyID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestGetByID_Exists_WithValue_WithStartDate() {
	yesterday := cachetime.NCacheTime(null.TimeFrom(time.Now().AddDate(0, 0, -1)))
	valueFilterInput := model.PropertyValueFilterInput{
		PropertyIDs: &[]uuid.UUID{t.testPropertyID},
		StartDate:   yesterday,
	}
	pageInfo := getDefaultPageInfo()

	property := t.getNewProperty(nuuid.NUUID{}, nil)
	value1 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(900), time.Now())
	valueSlice := []model.PropertyValue{value1, value2}
	property.AttachValues(valueSlice, true)

	resolvedPropertySlice := []model.Property{property}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(resolvedPropertySlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testPropertyID, true, yesterday, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestGetByID_Exists_WithValue_WithEndDate() {
	today := cachetime.NCacheTime(null.TimeFrom(time.Now()))
	valueFilterInput := model.PropertyValueFilterInput{
		PropertyIDs: &[]uuid.UUID{t.testPropertyID},
		EndDate:     today,
	}
	pageInfo := getDefaultPageInfo()

	property := t.getNewProperty(nuuid.NUUID{}, nil)
	value1 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(900), time.Now())
	valueSlice := []model.PropertyValue{value1, value2}
	property.AttachValues(valueSlice, true)

	resolvedPropertySlice := []model.Property{property}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(resolvedPropertySlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testPropertyID, true, cachetime.NCacheTime{}, today, nil)

	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestGetByID_Exists_WithValue_WithBothDates() {
	yesterday := cachetime.NCacheTime(null.TimeFrom(time.Now().AddDate(0, 0, -1)))
	today := cachetime.NCacheTime(null.TimeFrom(time.Now()))
	valueFilterInput := model.PropertyValueFilterInput{
		PropertyIDs: &[]uuid.UUID{t.testPropertyID},
		StartDate:   yesterday,
		EndDate:     today,
	}
	pageInfo := getDefaultPageInfo()

	property := t.getNewProperty(nuuid.NUUID{}, nil)
	value1 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(900), time.Now())
	valueSlice := []model.PropertyValue{value1, value2}
	property.AttachValues(valueSlice, true)

	resolvedPropertySlice := []model.Property{property}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(resolvedPropertySlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testPropertyID, true, yesterday, today, nil)

	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestGetByID_Exists_WithValue_WithPageSize() {
	pageSize := 120
	valueFilterInput := model.PropertyValueFilterInput{
		PropertyIDs: &[]uuid.UUID{t.testPropertyID},
	}
	valueFilterInput.PageSize = &pageSize
	pageInfo := model.PageInfoOutput{
		PageSize: pageSize,
	}

	property := t.getNewProperty(nuuid.NUUID{}, nil)
	value1 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(900), time.Now())
	valueSlice := []model.PropertyValue{value1, value2}
	property.AttachValues(valueSlice, true)

	resolvedPropertySlice := []model.Property{property}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(resolvedPropertySlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return(valueSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testPropertyID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, &pageSize)

	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestGetByID_Exists_With_Value_RepoErrorResolvingProperty() {
	errMsg := "failed to resolve property by IDs"
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, failure.InternalError("resolve by IDs", "Property", errors.New(errMsg)))

	actual, err := t.svc.GetByID(t.testPropertyID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), "Property", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, errMsg)
	assert.Equal(t.T(), "resolve by IDs", *errAsFailure.Operation)
}

func (t *propertiesServiceTestSuite) TestGetByID_Exists_WithValue_RepoErrorResolvingValues() {
	errMsg := "error resolving values"
	valueFilterInput := model.PropertyValueFilterInput{
		PropertyIDs: &[]uuid.UUID{t.testPropertyID},
	}

	property := t.getNewProperty(nuuid.NUUID{}, nil)
	value1 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	value2 := t.getNewPropertyValue(nuuid.NUUID{}, nuuid.From(property.ID), float64(900), time.Now())
	valueSlice := []model.PropertyValue{value1, value2}
	property.AttachValues(valueSlice, true)
	resolvedPropertySlice := []model.Property{property}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(resolvedPropertySlice, nil)
	t.mockRepo.EXPECT().ResolveValuesByFilter(valueFilterInput.ToFilter()).
		Return([]model.PropertyValue{}, getDefaultPageInfo(), failure.InternalError("resolve by filter", "Property Value", errors.New(errMsg)))

	actual, err := t.svc.GetByID(t.testPropertyID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), "Property Value", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, errMsg)
	assert.Equal(t.T(), "resolve by filter", *errAsFailure.Operation)
}

func (t *propertiesServiceTestSuite) TestGetByID_NotExists() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, nil)

	actual, err := t.svc.GetByID(t.testPropertyID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), "Property", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, "not found")
	assert.Equal(t.T(), "get by ID", *errAsFailure.Operation)
}

func (t *propertiesServiceTestSuite) TestGetByFilter_EmptyFilter() {
	filterInput := model.PropertyFilterInput{}
	filter := filterInput.ToFilter()

	t.mockRepo.EXPECT().ResolveByFilter(filter).
		Return(t.getPropertySlice(2), getDefaultPageInfo(), nil)

	_, _, err := t.svc.GetByFilter(filterInput)

	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestGetByFilter_WithKeyword() {
	keyword := "example"
	filterInput := model.PropertyFilterInput{}
	filterInput.Keyword = &keyword
	filter := filterInput.ToFilter()

	t.mockRepo.EXPECT().ResolveByFilter(filter).
		Return(t.getPropertySlice(2), getDefaultPageInfo(), nil)

	_, _, err := t.svc.GetByFilter(filterInput)

	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestUpdate_Normal() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{t.getNewProperty(nuuid.From(t.testPropertyID), nil)}, nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).
		Return(nil)

	res, err := t.svc.Update(t.getNewPropertyInput(nuuid.From(t.testPropertyID)), t.testUserID)

	assert.NotNil(t.T(), res)
	assert.NoError(t.T(), err)
}

func (t *propertiesServiceTestSuite) TestUpdate_RepoErrorResolvingByIDs() {
	errMsg := "query failed"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, failure.InternalError("resolve by IDs", "Property", errors.New(errMsg)))

	res, err := t.svc.Update(t.getNewPropertyInput(nuuid.From(t.testPropertyID)), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "resolve by IDs", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesServiceTestSuite) TestUpdate_PropertyNotFound() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, nil)

	res, err := t.svc.Update(t.getNewPropertyInput(nuuid.From(t.testPropertyID)), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeEntityNotFound, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "not found")
}

func (t *propertiesServiceTestSuite) TestUpdate_PropertyDeleted() {
	propertyInput := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	deletedProperty := model.NewPropertyFromInput(propertyInput, t.testUserID)
	deletedProperty.Delete(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{deletedProperty}, nil)

	res, err := t.svc.Update(t.getNewPropertyInput(nuuid.From(t.testPropertyID)), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeOperationNotPermitted, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), "already deleted")
}

func (t *propertiesServiceTestSuite) TestUpdate_RepoErrorUpdating() {
	errMsg := "failed to update"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{t.getNewProperty(nuuid.From(t.testPropertyID), nil)}, nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).
		Return(failure.InternalError("update", "Property", errors.New(errMsg)))

	res, err := t.svc.Update(t.getNewPropertyInput(nuuid.From(t.testPropertyID)), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.IsType(t.T(), &failure.Failure{}, err)
	assert.Equal(t.T(), failure.CodeInternalError, err.(*failure.Failure).Code)
	assert.Equal(t.T(), "Property", *err.(*failure.Failure).Entity)
	assert.Equal(t.T(), "update", *err.(*failure.Failure).Operation)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *propertiesServiceTestSuite) TestDelete_Normal() {
	valuesSlice := []model.PropertyValue{}

	valuesSlice = append(
		valuesSlice,
		t.getNewPropertyValue(
			nuuid.NUUID{},
			nuuid.From(t.testPropertyID),
			float64(10000),
			time.Now().AddDate(0, 0, -1)))

	valuesSlice = append(
		valuesSlice,
		t.getNewPropertyValue(
			nuuid.NUUID{},
			nuuid.From(t.testPropertyID),
			float64(9000),
			time.Now()))

	testProperty := t.getNewProperty(nuuid.From(t.testPropertyID), &valuesSlice)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testProperty}, nil)

	t.mockRepo.EXPECT().ResolveValuesByFilter(gomock.Any()).
		Return(valuesSlice, getDefaultPageInfo(), nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

	res, err := t.svc.Delete(t.testPropertyID, t.testUserID)

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

func (t *propertiesServiceTestSuite) TestDelete_RepoErrorResolvingByIDs() {
	errMsg := "failed resolving by IDs"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, errors.New(errMsg))

	res, err := t.svc.Delete(t.testPropertyID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDelete_RepoErrorResolvingValuesByFilter() {
	errMsg := "failed resolving property values"

	testProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testProperty}, nil)

	t.mockRepo.EXPECT().ResolveValuesByFilter(gomock.Any()).
		Return([]model.PropertyValue{}, model.PageInfoOutput{}, errors.New(errMsg))

	res, err := t.svc.Delete(t.testPropertyID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)

	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDelete_PropertyNotFound() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, nil)

	res, err := t.svc.Delete(t.testPropertyID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDelete_PropertyAlreadyDeleted() {
	testDeletedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	testDeletedProperty.Deleted = null.TimeFrom(time.Now())
	testDeletedProperty.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testDeletedProperty}, nil)

	res, err := t.svc.Delete(t.testPropertyID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "delete")
	assert.Contains(t.T(), err.Error(), "Property")
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDelete_PropertyValueAlreadyDeleted() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{t.getNewProperty(nuuid.From(t.testPropertyID), nil)}, nil)

	testDeletedNonLastPropertyValue := t.getNewPropertyValue(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(10000),
		time.Now().AddDate(0, 0, -1))

	testDeletedNonLastPropertyValue.Deleted = null.TimeFrom(time.Now())
	testDeletedNonLastPropertyValue.DeletedBy = nuuid.From(t.testUserID)

	testLastPropertyValue := t.getNewPropertyValue(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(12000),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByFilter(gomock.Any()).
		Return(
			[]model.PropertyValue{
				testDeletedNonLastPropertyValue,
				testLastPropertyValue,
			},
			getDefaultPageInfo(),
			nil)

	res, err := t.svc.Delete(t.testPropertyID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "delete")
	assert.Contains(t.T(), err.Error(), "Property Value")
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDelete_RepoErrorUpdating() {
	errMsg := "failed updating property"
	valueSlice := []model.PropertyValue{}

	valueSlice = append(
		valueSlice,
		t.getNewPropertyValue(
			nuuid.NUUID{},
			nuuid.From(t.testPropertyID),
			float64(10000),
			time.Now().AddDate(0, 0, -1)))

	valueSlice = append(
		valueSlice,
		t.getNewPropertyValue(
			nuuid.NUUID{},
			nuuid.From(t.testPropertyID),
			float64(12000),
			time.Now()))

	testProperty := t.getNewProperty(nuuid.From(t.testPropertyID), &valueSlice)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testProperty}, nil)

	t.mockRepo.EXPECT().ResolveValuesByFilter(gomock.Any()).
		Return(valueSlice, getDefaultPageInfo(), nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).Return(errors.New(errMsg))

	res, err := t.svc.Delete(t.testPropertyID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)

	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_Normal_CurrentValue() {
	testValueDate := time.Now()
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)
	testValue := t.getNewPropertyValue(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)

	testPropertyToUpdate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	testPropertyToUpdate.CurrentValue = float64(900)
	testPropertyToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	testPropertyAfterUpate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	testPropertyAfterUpate.PurchaseDate = testPropertyToUpdate.PurchaseDate
	testPropertyAfterUpate.InitialValueDate = testPropertyToUpdate.InitialValueDate
	testPropertyAfterUpate.CurrentValue = testValue.Value
	testPropertyAfterUpate.CurrentValueDate = testValue.Date

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testPropertyToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return(
			[]model.PropertyValue{
				t.getNewPropertyValue(
					nuuid.NUUID{},
					nuuid.From(t.testPropertyID),
					float64(900),
					time.Now().AddDate(0, 0, -1)),
			},
			nil)

	t.mockRepo.EXPECT().CreateValue(
		propertyValueMatcher{testValue},
		propertyPointerMatcher{testPropertyAfterUpate}).
		Return(nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_Normal_NotCurrentValue() {
	testValueDate := time.Now().AddDate(0, 0, -12)
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)
	testValue := t.getNewPropertyValue(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)

	testPropertyToUpdate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	testPropertyToUpdate.CurrentValue = float64(900)
	testPropertyToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testPropertyToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return(
			[]model.PropertyValue{
				t.getNewPropertyValue(
					nuuid.NUUID{},
					nuuid.From(t.testPropertyID),
					float64(900),
					time.Now().AddDate(0, 0, -1)),
			},
			nil)

	t.mockRepo.EXPECT().
		CreateValue(propertyValueMatcher{testValue}, nil).
		Return(nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_RepoFailedResolvingByIDs() {
	errMsg := "failed to resolve property by IDs"
	testValueDate := time.Now()
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, errors.New(errMsg))

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_PropertyNotFound() {
	testValueDate := time.Now()
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_PropertyDeleted() {
	testValueDate := time.Now()
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)

	deletedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	deletedProperty.Deleted = null.TimeFrom(time.Now())
	deletedProperty.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{deletedProperty}, nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_PropertySold() {
	testValueDate := time.Now()
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)

	deletedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	deletedProperty.Status = model.PropertyStatusSold

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{deletedProperty}, nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "sold")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_RepoFailedResolvingCurrentValue() {
	errMsg := "failed resolving property values by property ID"
	testValueDate := time.Now().AddDate(0, 0, -12)
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)

	testPropertyToUpdate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	testPropertyToUpdate.CurrentValue = float64(900)
	testPropertyToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testPropertyToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return([]model.PropertyValue{}, errors.New(errMsg))

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_CurrentValueNotFound() {
	testValueDate := time.Now().AddDate(0, 0, -12)
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)

	testPropertyToUpdate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	testPropertyToUpdate.CurrentValue = float64(900)
	testPropertyToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testPropertyToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return([]model.PropertyValue{}, nil)

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestCreateValue_RepoFailedCreatingPropertyValue() {
	errMsg := "failed creating property value"
	testValueDate := time.Now().AddDate(0, 0, -12)
	testInput := t.getNewPropertyValueInput(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate)
	testValue := t.getNewPropertyValue(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(1234.56),
		testValueDate,
	)

	testPropertyToUpdate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	testPropertyToUpdate.CurrentValue = float64(900)
	testPropertyToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{testPropertyToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return(
			[]model.PropertyValue{
				t.getNewPropertyValue(
					nuuid.NUUID{},
					nuuid.From(t.testPropertyID),
					float64(900),
					time.Now().AddDate(0, 0, -1))},
			nil)

	t.mockRepo.EXPECT().CreateValue(propertyValueMatcher{testValue}, nil).
		Return(errors.New(errMsg))

	res, err := t.svc.CreateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestGetValueByID_Normal() {
	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return(
			[]model.PropertyValue{
				t.getNewPropertyValue(
					nuuid.From(t.testPropertyValueID),
					nuuid.From(t.testPropertyID),
					float64(1000000), time.Now(),
				),
			},
			nil)

	res, err := t.svc.GetValueByID(t.testPropertyValueID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestGetValueByID_RepoFailedResolvingValue() {
	errMsg := "failed to resolve property values"
	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return(
			[]model.PropertyValue{},
			failure.InternalError("resolve by filter", "Property Value", errors.New(errMsg)))

	actual, err := t.svc.GetValueByID(t.testPropertyValueID)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, errAsFailure.Code)
	assert.Equal(t.T(), "Property Value", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, errMsg)
	assert.Equal(t.T(), "resolve by filter", *errAsFailure.Operation)
}

func (t *propertiesServiceTestSuite) TestGetValueByID_ValueNotFound() {
	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return(
			[]model.PropertyValue{},
			nil)

	actual, err := t.svc.GetValueByID(t.testPropertyValueID)
	errAsFailure, ok := err.(*failure.Failure)
	if !ok {
		t.T().Fatal("failed converting error to failure object")
	}

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), failure.CodeEntityNotFound, errAsFailure.Code)
	assert.Equal(t.T(), "Property Value", *errAsFailure.Entity)
	assert.Contains(t.T(), errAsFailure.Message, "not found")
	assert.Equal(t.T(), "get by ID", *errAsFailure.Operation)
}

func (t *propertiesServiceTestSuite) TestGetValueBVyFilter_Normal() {
	filter := model.PropertyValueFilterInput{}
	t.mockRepo.EXPECT().ResolveValuesByFilter(filter.ToFilter()).
		Return(
			[]model.PropertyValue{
				t.getNewPropertyValue(
					nuuid.From(t.testPropertyValueID),
					nuuid.From(t.testPropertyValueID),
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

func (t *propertiesServiceTestSuite) TestUpdateValue_Normal_CurrentValue() {
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	propertyToUpdate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	propertyToUpdate.CurrentValue = float64(900)
	propertyToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	valueToUpdate := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		propertyToUpdate.CurrentValue,
		propertyToUpdate.CurrentValueDate)

	updatedValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		testInput.Value,
		testInput.Date.Time())

	updatedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	updatedProperty.PurchaseDate = propertyToUpdate.PurchaseDate
	updatedProperty.InitialValueDate = propertyToUpdate.InitialValueDate
	updatedProperty.CurrentValue = updatedValue.Value
	updatedProperty.CurrentValueDate = updatedValue.Date

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{propertyToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return([]model.PropertyValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		propertyValueMatcher{updatedValue},
		propertyPointerMatcher{updatedProperty}).
		Return(nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_Normal_NonCurrentValue() {
	newValueID, _ := uuid.NewV7()
	testInput := t.getNewPropertyValueInput(
		nuuid.From(newValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now().AddDate(0, 0, -2),
	)

	propertyToUpdate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	propertyToUpdate.CurrentValue = float64(900)
	propertyToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	valueToUpdate := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		propertyToUpdate.CurrentValue,
		propertyToUpdate.CurrentValueDate)

	updatedValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		testInput.Value,
		testInput.Date.Time())

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{propertyToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{newValueID}).
		Return([]model.PropertyValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return([]model.PropertyValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		propertyValueMatcher{updatedValue},
		nil).
		Return(nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_RepoFailedResolvingByIDs() {
	errMsg := "failed resolving properties by IDs"
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, errors.New(errMsg))

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_PropertyNotFound() {
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Property")
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_PropertyDeleted() {
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.CurrentValue = float64(900)
	resolvedProperty.CurrentValueDate = time.Now().AddDate(0, 0, -1)
	resolvedProperty.Deleted = null.TimeFrom(time.Now())
	resolvedProperty.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{resolvedProperty}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_PropertySold() {
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.CurrentValue = float64(900)
	resolvedProperty.CurrentValueDate = time.Now().AddDate(0, 0, -1)
	resolvedProperty.Status = model.PropertyStatusSold

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{resolvedProperty}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "sold")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_RepoFailedResolvingValues() {
	errMsg := "failed resolving property values"
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.CurrentValue = float64(900)
	resolvedProperty.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{resolvedProperty}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{}, errors.New(errMsg))

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_ValueNotFound() {
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.CurrentValue = float64(900)
	resolvedProperty.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{resolvedProperty}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Property Value")
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_ValueDeleted() {
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.CurrentValue = float64(900)
	resolvedProperty.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	resolvedValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		resolvedProperty.CurrentValue,
		resolvedProperty.CurrentValueDate)
	resolvedValue.Deleted = null.TimeFrom(time.Now())
	resolvedValue.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{resolvedProperty}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{resolvedValue}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Property Value")
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_RepoFailedResolvingCurrentValues() {
	errMsg := "failed resolving property last value"
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.CurrentValue = float64(900)
	resolvedProperty.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	resolvedPropertyValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		resolvedProperty.CurrentValue,
		resolvedProperty.CurrentValueDate)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{resolvedProperty}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{resolvedPropertyValue}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return([]model.PropertyValue{}, errors.New(errMsg))

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_CurrentValueNotFound() {
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.CurrentValue = float64(900)
	resolvedProperty.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	resolvedPropertyValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		resolvedProperty.CurrentValue,
		resolvedProperty.CurrentValueDate)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{resolvedProperty}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{resolvedPropertyValue}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return([]model.PropertyValue{}, nil)

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Property Value")
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestUpdateValue_RepoFailedUpdatingValue() {
	errMsg := "failed updating property value"
	testInput := t.getNewPropertyValueInput(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(1000),
		time.Now(),
	)

	propertyToUpdate := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	propertyToUpdate.CurrentValue = float64(900)
	propertyToUpdate.CurrentValueDate = time.Now().AddDate(0, 0, -1)

	valueToUpdate := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		propertyToUpdate.CurrentValue,
		propertyToUpdate.CurrentValueDate)

	updatedValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		testInput.Value,
		testInput.Date.Time())

	updatedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	updatedProperty.PurchaseDate = propertyToUpdate.PurchaseDate
	updatedProperty.InitialValueDate = propertyToUpdate.InitialValueDate
	updatedProperty.CurrentValue = updatedValue.Value
	updatedProperty.CurrentValueDate = updatedValue.Date

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{propertyToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 1).
		Return([]model.PropertyValue{valueToUpdate}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		propertyValueMatcher{updatedValue},
		propertyPointerMatcher{updatedProperty}).
		Return(errors.New(errMsg))

	res, err := t.svc.UpdateValue(testInput, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_Normal_CurrentValue() {

	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	secondToCurrentValue := t.getNewPropertyValue(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(456),
		time.Now().AddDate(0, 0, -12))

	deletedCurrentValue := t.getNewPropertyValue(
		nuuid.From(lastValue.ID),
		nuuid.From(lastValue.PropertyID),
		lastValue.Value,
		lastValue.Date)
	deletedCurrentValue.Deleted = null.TimeFrom(time.Now())
	deletedCurrentValue.DeletedBy = nuuid.From(t.testUserID)

	property := t.getNewProperty(nuuid.From(t.testPropertyID), nil)

	updatedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	updatedProperty.PurchaseDate = property.PurchaseDate
	updatedProperty.InitialValueDate = property.InitialValueDate
	updatedProperty.CurrentValue = secondToCurrentValue.Value
	updatedProperty.CurrentValueDate = secondToCurrentValue.Date

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{property}, nil)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 2).
		Return([]model.PropertyValue{lastValue, secondToCurrentValue}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		propertyValueMatcher{deletedCurrentValue},
		propertyPointerMatcher{updatedProperty}).
		Return(nil)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_Normal_NonCurrentValue() {

	lastValue := t.getNewPropertyValue(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	secondToCurrentValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(456),
		time.Now().AddDate(0, 0, -12))

	deletedNonCurrentValue := t.getNewPropertyValue(
		nuuid.From(secondToCurrentValue.ID),
		nuuid.From(secondToCurrentValue.PropertyID),
		secondToCurrentValue.Value,
		secondToCurrentValue.Date)
	deletedNonCurrentValue.Deleted = null.TimeFrom(time.Now())
	deletedNonCurrentValue.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{secondToCurrentValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(
			[]model.Property{t.getNewProperty(nuuid.From(t.testPropertyID), nil)},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 2).
		Return([]model.PropertyValue{lastValue, secondToCurrentValue}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		propertyValueMatcher{deletedNonCurrentValue},
		nil).
		Return(nil)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.NoError(t.T(), err)
	assert.NotNil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_RepoFailedResolvingValueByIDs() {
	errMsg := "failed resolving property values by IDs"

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{}, errors.New(errMsg))

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_ValueNotFound() {
	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{}, nil)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_ValueAlreadyDeleted() {

	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())
	lastValue.Deleted = null.TimeFrom(time.Now())
	lastValue.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Property Value")
	assert.Contains(t.T(), err.Error(), "already deleted")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_RepoFailedResolvingByIDs() {
	errMsg := "failed resolving properties by IDs"
	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, errors.New(errMsg))

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_PropertyNotFound() {
	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return([]model.Property{}, nil)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_PropertyDeleted() {

	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.Deleted = null.TimeFrom(time.Now())
	resolvedProperty.DeletedBy = nuuid.From(t.testPropertyID)

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(
			[]model.Property{resolvedProperty},
			nil,
		)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Property")
	assert.Contains(t.T(), err.Error(), "already deleted")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_PropertySold() {

	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	resolvedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	resolvedProperty.Status = model.PropertyStatusSold

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(
			[]model.Property{resolvedProperty},
			nil,
		)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Property")
	assert.Contains(t.T(), err.Error(), "sold")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_RepoFailedResolvingLastValues() {
	errMsg := "failed resolving property last values"
	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(
			[]model.Property{t.getNewProperty(nuuid.From(t.testPropertyID), nil)},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 2).
		Return([]model.PropertyValue{}, errors.New(errMsg))

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_CurrentValueNotFound() {

	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(
			[]model.Property{t.getNewProperty(nuuid.From(t.testPropertyID), nil)},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 2).
		Return([]model.PropertyValue{}, nil)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "Property Current Value")
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_CannotDeleteTheOnlyValue() {

	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(
			[]model.Property{t.getNewProperty(nuuid.From(t.testPropertyID), nil)},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 2).
		Return([]model.PropertyValue{lastValue}, nil)

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "OperationNotPermitted")
	assert.Contains(t.T(), err.Error(), "delete")
	assert.Contains(t.T(), err.Error(), "Property Value")
	assert.Contains(t.T(), err.Error(), "the only")
	assert.Nil(t.T(), res)
}

func (t *propertiesServiceTestSuite) TestDeleteValue_RepoFailedUpdatingValue() {
	errMsg := "failed updating property value"
	lastValue := t.getNewPropertyValue(
		nuuid.From(t.testPropertyValueID),
		nuuid.From(t.testPropertyID),
		float64(123),
		time.Now())

	secondToCurrentValue := t.getNewPropertyValue(
		nuuid.NUUID{},
		nuuid.From(t.testPropertyID),
		float64(456),
		time.Now().AddDate(0, 0, -12))

	deletedCurrentValue := t.getNewPropertyValue(
		nuuid.From(lastValue.ID),
		nuuid.From(lastValue.PropertyID),
		lastValue.Value,
		lastValue.Date)
	deletedCurrentValue.Deleted = null.TimeFrom(time.Now())
	deletedCurrentValue.DeletedBy = nuuid.From(t.testUserID)

	property := t.getNewProperty(nuuid.From(t.testPropertyID), nil)

	updatedProperty := t.getNewProperty(nuuid.From(t.testPropertyID), nil)
	updatedProperty.PurchaseDate = property.PurchaseDate
	updatedProperty.InitialValueDate = property.InitialValueDate
	updatedProperty.CurrentValue = secondToCurrentValue.Value
	updatedProperty.CurrentValueDate = secondToCurrentValue.Date

	t.mockRepo.EXPECT().ResolveValuesByIDs([]uuid.UUID{t.testPropertyValueID}).
		Return([]model.PropertyValue{lastValue}, nil)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testPropertyID}).
		Return(
			[]model.Property{property},
			nil,
		)

	t.mockRepo.EXPECT().ResolveLastValuesByPropertyID(t.testPropertyID, 2).
		Return([]model.PropertyValue{lastValue, secondToCurrentValue}, nil)

	t.mockRepo.EXPECT().UpdateValue(
		propertyValueMatcher{deletedCurrentValue},
		propertyPointerMatcher{updatedProperty}).
		Return(errors.New(errMsg))

	res, err := t.svc.DeleteValue(t.testPropertyValueID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

// matchers

type propertyPointerMatcher struct {
	expected model.Property
}

func (m propertyPointerMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*model.Property)
	if !ok {
		return false
	}

	return actual.Name == m.expected.Name &&
		actual.Address == m.expected.Address &&
		actual.TotalArea == m.expected.TotalArea &&
		actual.BuildingArea == m.expected.BuildingArea &&
		actual.AreaUnit == m.expected.AreaUnit &&
		actual.Type == m.expected.Type &&
		actual.TitleHolder == m.expected.TitleHolder &&
		actual.TaxIdentifier == m.expected.TaxIdentifier &&
		actual.PurchaseDate.Equal(m.expected.PurchaseDate) &&
		actual.InitialValue == m.expected.InitialValue &&
		actual.InitialValueDate.Equal(m.expected.InitialValueDate) &&
		actual.CurrentValue == m.expected.CurrentValue &&
		actual.CurrentValueDate.Equal(m.expected.CurrentValueDate) &&
		actual.AnnualAppreciationPercent == m.expected.AnnualAppreciationPercent &&
		actual.Status == m.expected.Status
}

func (m propertyPointerMatcher) String() string {
	return fmt.Sprintf(
		"is Property with Name=%s, Address=%s, TotalArea=%f, BuildingArea=%f, AreaUnit=%s, Type=%s, TitleHolder=%s, TaxIdentifier=%s, PurchaseDate=%s, InitialValue=%f, InitialValueDate=%s, CurrentValue=%f, CurrentValueDate=%s, AnnualAppreciationPercent=%f, Status=%s",
		m.expected.Name,
		m.expected.Address,
		m.expected.TotalArea,
		m.expected.BuildingArea,
		m.expected.AreaUnit,
		m.expected.Type,
		m.expected.TitleHolder,
		m.expected.TaxIdentifier,
		m.expected.PurchaseDate,
		m.expected.InitialValue,
		m.expected.InitialValueDate,
		m.expected.CurrentValue,
		m.expected.CurrentValueDate,
		m.expected.AnnualAppreciationPercent,
		m.expected.Status)
}

type propertyValueMatcher struct {
	expected model.PropertyValue
}

func (m propertyValueMatcher) Matches(x interface{}) bool {
	actual, ok := x.(model.PropertyValue)
	if !ok {
		return false
	}

	return actual.Date.Equal(m.expected.Date) &&
		actual.Value == m.expected.Value &&
		actual.PropertyID == m.expected.PropertyID
}

func (m propertyValueMatcher) String() string {
	return fmt.Sprintf(
		"is PropertyValue with Value=%.2f, Date=%v, PropertyID=%v",
		m.expected.Value,
		m.expected.Date,
		m.expected.PropertyID)
}
