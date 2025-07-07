package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/handler"
	"github.com/kerti/balances/backend/handler/response"
	mock_service "github.com/kerti/balances/backend/mock/service"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/nuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type propertyHandlerTestSuite struct {
	suite.Suite
	ctrl                *gomock.Controller
	handler             handler.Property
	mockSvc             *mock_service.MockProperty
	testUserID          uuid.UUID
	testPropertyID      uuid.UUID
	testPropertyValueID uuid.UUID
}

func TestPropertyHandler(t *testing.T) {
	suite.Run(t, new(propertyHandlerTestSuite))
}

func (t *propertyHandlerTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockSvc = mock_service.NewMockProperty(t.ctrl)
	t.handler = &handler.PropertyImpl{
		Service: t.mockSvc,
	}
	t.testUserID, _ = uuid.NewV7()
	t.testPropertyID, _ = uuid.NewV7()
	t.testPropertyValueID, _ = uuid.NewV7()
	t.handler.Startup()
}

func (t *propertyHandlerTestSuite) TearDownTest() {
	t.handler.Shutdown()
	t.ctrl.Finish()
}

func (t *propertyHandlerTestSuite) getNewRequestWithContext(method, path string, input any, formParams *map[string]string, routeVarId nuuid.NUUID) (recorder *httptest.ResponseRecorder, request *http.Request) {
	var reqBody *bytes.Buffer
	var req *http.Request

	if method == http.MethodPost || method == http.MethodPatch {
		// write body for POST and PATCH
		jsonBody, err := json.Marshal(input)
		if err != nil {
			t.T().Fatal(err)
		}

		reqBody = bytes.NewBuffer(jsonBody)
		req = httptest.NewRequest(method, path, reqBody)
	} else {
		// inject params into URL for all else
		if formParams != nil {
			query := make(url.Values)
			for k, v := range *formParams {
				if k != "id" {
					query.Add(k, v)
				}
			}

			// Append query to URL
			fullPath := path
			if encoded := query.Encode(); encoded != "" {
				fullPath += "?" + encoded
			}

			req = httptest.NewRequest(method, fullPath, nil)
		} else {
			req = httptest.NewRequest(method, path, nil)
		}
	}

	// set ID route var
	if routeVarId.Valid {
		req = mux.SetURLVars(req, map[string]string{
			"id": routeVarId.UUID.String(),
		})
	}

	req.Header.Set("Content-Type", "application/json")

	// add context with user ID
	ctx := req.Context()
	ctx = context.WithValue(ctx, ctxprops.PropUserID, &t.testUserID)

	request = req.WithContext(ctx)
	recorder = httptest.NewRecorder()

	return
}

func (t *propertyHandlerTestSuite) getNewPropertyInput(id nuuid.NUUID) model.PropertyInput {
	veh := model.PropertyInput{}

	if id.Valid {
		veh.ID = id.UUID
	} else {
		veh.ID = t.testPropertyID
	}

	initialValueDate := time.Now().AddDate(-2, 0, 0) // defaults to 2 years ago

	veh.Name = "John's House"
	veh.Address = "1234 Main Street"
	veh.TotalArea = 1200
	veh.BuildingArea = 1000
	veh.AreaUnit = model.PropertyAreaUnitSQM
	veh.Type = model.PropertyTypeHouse
	veh.TitleHolder = "John Fitzgerald Doe"
	veh.TaxIdentifier = "TUNEMAN"
	veh.PurchaseDate = cachetime.CacheTime(initialValueDate)
	veh.InitialValue = 68000
	veh.InitialValueDate = cachetime.CacheTime(initialValueDate)
	veh.CurrentValue = 50000
	veh.CurrentValueDate = cachetime.CacheTime(time.Now())
	veh.AnnualAppreciationPercent = 3.5
	veh.Status = model.PropertyStatusInUse

	return veh
}

func (t *propertyHandlerTestSuite) getNewPropertyValueInput(id, propertyID nuuid.NUUID) model.PropertyValueInput {
	vvi := model.PropertyValueInput{}

	if id.Valid {
		vvi.ID = id.UUID
	} else {
		vvi.ID = t.testPropertyValueID
	}

	if propertyID.Valid {
		vvi.PropertyID = propertyID.UUID
	} else {
		vvi.PropertyID = t.testPropertyID
	}

	vvi.Date = cachetime.CacheTime(time.Now().AddDate(0, 0, -1))
	vvi.Value = float64(50000)

	return vvi
}

func (t *propertyHandlerTestSuite) parseOutputToProperty(rr *httptest.ResponseRecorder) (actual *model.PropertyOutput, fail *failure.Failure) {
	// read the response
	var response response.BaseResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.T().Fatal(err)
	}

	if response.Data != nil {
		// marshal the data to JSON
		actualMap := (*response.Data).(map[string]any)
		jsonBytes, err := json.Marshal(actualMap)
		if err != nil {
			t.T().Fatal(err)
		}
		// unmarshal back to the expected object
		err = json.Unmarshal(jsonBytes, &actual)
		if err != nil {
			t.T().Fatal(err)
		}
		return actual, nil
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return actual, nil
}

func (t *propertyHandlerTestSuite) parseOutputToPropertyValue(rr *httptest.ResponseRecorder) (actual *model.PropertyValueOutput, fail *failure.Failure) {
	// read the response
	var response response.BaseResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.T().Fatal(err)
	}

	if response.Data != nil {
		// marshal the data to JSON
		actualMap := (*response.Data).(map[string]any)
		jsonBytes, err := json.Marshal(actualMap)
		if err != nil {
			t.T().Fatal(err)
		}
		// unmarshal back to the expected object
		err = json.Unmarshal(jsonBytes, &actual)
		if err != nil {
			t.T().Fatal(err)
		}
		return actual, nil
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return actual, nil
}

func (t *propertyHandlerTestSuite) parseOutputToPropertyPage(rr *httptest.ResponseRecorder) (items []model.PropertyOutput, pageInfo model.PageInfoOutput, fail *failure.Failure) {
	// read the response
	var response response.BaseResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.T().Fatal(err)
	}

	if response.Data != nil {
		// marshal the data to JSON
		actualMap := (*response.Data).(map[string]any)
		jsonBytes, err := json.Marshal(actualMap)
		if err != nil {
			t.T().Fatal(err)
		}

		// unmarshal back to the expected object
		var actual model.PageOutput
		err = json.Unmarshal(jsonBytes, &actual)
		if err != nil {
			t.T().Fatal(err)
		}

		//convert interface{} to []model.PropertyOutput
		actualSlice := (actual.Items).([]any)
		for _, propertyInterface := range actualSlice {
			propertyMap := (propertyInterface).(map[string]any)
			propertyJsonBytes, err := json.Marshal(propertyMap)
			if err != nil {
				t.T().Fatal(err)
			}
			var actualProperty model.PropertyOutput
			err = json.Unmarshal(propertyJsonBytes, &actualProperty)
			if err != nil {
				t.T().Fatal(err)
			}
			items = append(items, actualProperty)
		}

		pageInfo = actual.PageInfo
	}

	if response.Error != nil {
		fail = response.Error
	}

	return
}

func (t *propertyHandlerTestSuite) parseOutputToPropertyValuePage(rr *httptest.ResponseRecorder) (items []model.PropertyValueOutput, pageInfo model.PageInfoOutput, fail *failure.Failure) {
	// read the response
	var response response.BaseResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.T().Fatal(err)
	}

	if response.Data != nil {
		// marshal the data to JSON
		actualMap := (*response.Data).(map[string]any)
		jsonBytes, err := json.Marshal(actualMap)
		if err != nil {
			t.T().Fatal(err)
		}

		// unmarshal back to the expected object
		var actual model.PageOutput
		err = json.Unmarshal(jsonBytes, &actual)
		if err != nil {
			t.T().Fatal(err)
		}

		//convert interface{} to []model.PropertyOutput
		actualSlice := (actual.Items).([]any)
		for _, propertyValueInterface := range actualSlice {
			propertyValueMap := (propertyValueInterface).(map[string]any)
			propertyValueJsonBytes, err := json.Marshal(propertyValueMap)
			if err != nil {
				t.T().Fatal(err)
			}
			var actualPropertyValue model.PropertyValueOutput
			err = json.Unmarshal(propertyValueJsonBytes, &actualPropertyValue)
			if err != nil {
				t.T().Fatal(err)
			}
			items = append(items, actualPropertyValue)
		}

		pageInfo = actual.PageInfo
	}

	if response.Error != nil {
		fail = response.Error
	}

	return
}

func (t *propertyHandlerTestSuite) TestCreate_Normal() {
	input := t.getNewPropertyInput(nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedResult := model.NewPropertyFromInput(input, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().Create(gomock.Any(), t.testUserID).Return(&expectedResult, nil)

	t.handler.HandleCreateProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusCreated, rr.Result().StatusCode)
	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), expected.Name, actual.Name)
	assert.Equal(t.T(), expected.Address, actual.Address)
	assert.Equal(t.T(), expected.TotalArea, actual.TotalArea)
	assert.Equal(t.T(), expected.BuildingArea, actual.BuildingArea)
	assert.Equal(t.T(), expected.AreaUnit, actual.AreaUnit)
	assert.Equal(t.T(), expected.TitleHolder, actual.TitleHolder)
	assert.Equal(t.T(), expected.TaxIdentifier, actual.TaxIdentifier)
	assert.Equal(t.T(), expected.PurchaseDate.Time().Unix(), actual.PurchaseDate.Time().Unix())
	assert.Equal(t.T(), expected.InitialValue, actual.InitialValue)
	assert.Equal(t.T(), expected.InitialValueDate.Time().Unix(), actual.InitialValueDate.Time().Unix())
	assert.Equal(t.T(), expected.CurrentValue, actual.CurrentValue)
	assert.Equal(t.T(), expected.CurrentValueDate.Time().Unix(), actual.CurrentValueDate.Time().Unix())
	assert.Equal(t.T(), expected.AnnualAppreciationPercent, actual.AnnualAppreciationPercent)
	assert.Equal(t.T(), expected.Status, actual.Status)
	assert.NotNil(t.T(), actual.Created)
	assert.NotNil(t.T(), actual.CreatedBy)
	assert.False(t.T(), actual.Updated.Valid)
	assert.False(t.T(), actual.UpdatedBy.Valid)
	assert.False(t.T(), actual.Deleted.Valid)
	assert.False(t.T(), actual.DeletedBy.Valid)
}

func (t *propertyHandlerTestSuite) TestCreate_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleCreateProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestCreate_ServiceFailedCreating() {
	errMsg := "service failed creating property"
	input := t.getNewPropertyInput(nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().Create(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleCreateProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestGetByID_Normal_NoParams() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/"+t.testPropertyID.String(),
		nil,
		nil,
		nuuid.From(t.testPropertyID),
	)

	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	expectedResult := model.NewPropertyFromInput(input, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().GetByID(t.testPropertyID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).Return(&expectedResult, nil)

	t.handler.HandleGetPropertyByID(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.Name, actual.Name)
	assert.Equal(t.T(), expected.Address, actual.Address)
	assert.Equal(t.T(), expected.TotalArea, actual.TotalArea)
	assert.Equal(t.T(), expected.BuildingArea, actual.BuildingArea)
	assert.Equal(t.T(), expected.AreaUnit, actual.AreaUnit)
	assert.Equal(t.T(), expected.TitleHolder, actual.TitleHolder)
	assert.Equal(t.T(), expected.TaxIdentifier, actual.TaxIdentifier)
	assert.Equal(t.T(), expected.PurchaseDate.Time().Unix(), actual.PurchaseDate.Time().Unix())
	assert.Equal(t.T(), expected.InitialValue, actual.InitialValue)
	assert.Equal(t.T(), expected.InitialValueDate.Time().Unix(), actual.InitialValueDate.Time().Unix())
	assert.Equal(t.T(), expected.CurrentValue, actual.CurrentValue)
	assert.Equal(t.T(), expected.CurrentValueDate.Time().Unix(), actual.CurrentValueDate.Time().Unix())
	assert.Equal(t.T(), expected.AnnualAppreciationPercent, actual.AnnualAppreciationPercent)
	assert.Equal(t.T(), expected.Status, actual.Status)
	assert.Equal(t.T(), expected.Created.Time().Unix(), actual.Created.Time().Unix())
	assert.Equal(t.T(), expected.CreatedBy, actual.CreatedBy)
	assert.Equal(t.T(), expected.Updated, actual.Updated)
	assert.Equal(t.T(), expected.UpdatedBy, actual.UpdatedBy)
	assert.Equal(t.T(), expected.Deleted, actual.Deleted)
	assert.Equal(t.T(), expected.DeletedBy, actual.DeletedBy)
}

func (t *propertyHandlerTestSuite) TestGetByID_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/"+t.testPropertyID.String()+"123",
		nil,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetPropertyByID(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestGetByID_Normal_WithValues() {
	formParams := make(map[string]string)
	formParams["withValues"] = "true"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/"+t.testPropertyID.String(),
		nil,
		&formParams,
		nuuid.From(t.testPropertyID),
	)

	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	expectedResult := model.NewPropertyFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().GetByID(t.testPropertyID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).Return(&expectedResult, nil)

	t.handler.HandleGetPropertyByID(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *propertyHandlerTestSuite) TestGetByID_Normal_WithValueStartDate() {
	startDate := time.Unix(0, time.Now().AddDate(0, 0, -1).UnixMilli()*int64(time.Millisecond))
	formParams := make(map[string]string)
	formParams["valueStartDate"] = strconv.FormatInt(startDate.UnixMilli(), 10)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/"+t.testPropertyID.String(),
		nil,
		&formParams,
		nuuid.From(t.testPropertyID),
	)

	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	expectedResult := model.NewPropertyFromInput(input, t.testUserID)

	var nStartDate cachetime.NCacheTime
	nStartDate.Scan(startDate)
	t.mockSvc.EXPECT().GetByID(t.testPropertyID, false, nStartDate, cachetime.NCacheTime{}, nil).Return(&expectedResult, nil)

	t.handler.HandleGetPropertyByID(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *propertyHandlerTestSuite) TestGetByID_Normal_WithValueEndDate() {
	endDate := time.Unix(0, time.Now().AddDate(0, 0, -1).UnixMilli()*int64(time.Millisecond))
	formParams := make(map[string]string)
	formParams["valueEndDate"] = strconv.FormatInt(endDate.UnixMilli(), 10)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/"+t.testPropertyID.String(),
		nil,
		&formParams,
		nuuid.From(t.testPropertyID),
	)

	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	expectedResult := model.NewPropertyFromInput(input, t.testUserID)

	var nEndDate cachetime.NCacheTime
	nEndDate.Scan(endDate)
	t.mockSvc.EXPECT().GetByID(t.testPropertyID, false, cachetime.NCacheTime{}, nEndDate, nil).Return(&expectedResult, nil)

	t.handler.HandleGetPropertyByID(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *propertyHandlerTestSuite) TestGetByID_Normal_WithPageSize() {
	pageSize := 10
	formParams := make(map[string]string)
	formParams["pageSize"] = strconv.Itoa(pageSize)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/"+t.testPropertyID.String(),
		nil,
		&formParams,
		nuuid.From(t.testPropertyID),
	)

	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	expectedResult := model.NewPropertyFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().GetByID(t.testPropertyID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, &pageSize).Return(&expectedResult, nil)

	t.handler.HandleGetPropertyByID(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *propertyHandlerTestSuite) TestGetByID_ServiceFailedResolving() {
	errMsg := "service failed resolving"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/"+t.testPropertyID.String(),
		nil,
		nil,
		nuuid.From(t.testPropertyID),
	)

	t.mockSvc.EXPECT().GetByID(t.testPropertyID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).Return(nil, errors.New(errMsg))

	t.handler.HandleGetPropertyByID(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestGetByFilter_Normal() {
	keyword := "test keyword"
	input := model.PropertyFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedProperties := []model.Property{}
	v1 := model.NewPropertyFromInput(t.getNewPropertyInput(nuuid.NUUID{}), t.testUserID)
	v2 := model.NewPropertyFromInput(t.getNewPropertyInput(nuuid.NUUID{}), t.testUserID)
	expectedProperties = append(expectedProperties, v1)
	expectedProperties = append(expectedProperties, v2)
	expectedPageInfo := model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 1,
		PageCount:  1,
	}

	t.mockSvc.EXPECT().GetByFilter(input).Return(expectedProperties, expectedPageInfo, nil)

	t.handler.HandleGetPropertyByFilter(rr, req)

	properties, pageInfo, err := t.parseOutputToPropertyPage(rr)

	assert.Nil(t.T(), err)

	assert.Equal(t.T(), len(expectedProperties), len(properties))
	assert.Equal(t.T(), expectedProperties[0].ID, properties[0].ID)
	assert.Equal(t.T(), expectedProperties[1].ID, properties[1].ID)

	assert.Equal(t.T(), 1, pageInfo.Page)
}

func (t *propertyHandlerTestSuite) TestGetByFilter_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetPropertyByFilter(rr, req)

	properties, pageInfo, err := t.parseOutputToPropertyPage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(properties))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *propertyHandlerTestSuite) TestGetByFilter_ServiceFailedResolving() {
	errMsg := "failed resolving properties by filter"
	keyword := "test keyword"
	input := model.PropertyFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().GetByFilter(input).Return([]model.Property{}, model.PageInfoOutput{}, errors.New(errMsg))

	t.handler.HandleGetPropertyByFilter(rr, req)

	properties, pageInfo, err := t.parseOutputToPropertyPage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(properties))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *propertyHandlerTestSuite) TestUpdate_Normal() {
	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/"+t.testPropertyID.String(),
		input,
		nil,
		nuuid.From(t.testPropertyID),
	)

	updatedProperty := model.NewPropertyFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().Update(gomock.Any(), t.testUserID).Return(&updatedProperty, nil)

	t.handler.HandleUpdateProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *propertyHandlerTestSuite) TestUpdate_FailedGettingIDFromRequest() {
	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/"+t.testPropertyID.String(),
		input,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleUpdateProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestUpdate_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/"+t.testPropertyID.String(),
		input,
		nil,
		nuuid.From(t.testPropertyID),
	)

	t.handler.HandleUpdateProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestUpdate_MismatchedID() {
	input := t.getNewPropertyInput(nuuid.NUUID{})
	newID, _ := uuid.NewV7()
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/"+t.testPropertyID.String(),
		input,
		nil,
		nuuid.From(newID),
	)

	t.handler.HandleUpdateProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "id mismatch")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestUpdate_ServiceFailedUpdating() {
	errMsg := "failed updating property"
	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/"+t.testPropertyID.String(),
		input,
		nil,
		nuuid.From(t.testPropertyID),
	)

	t.mockSvc.EXPECT().Update(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleUpdateProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestDelete_Normal() {
	input := t.getNewPropertyInput(nuuid.From(t.testPropertyID))
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/properties/"+t.testPropertyID.String(),
		nil,
		nil,
		nuuid.From(t.testPropertyID),
	)

	deletedProperty := model.NewPropertyFromInput(input, t.testUserID)
	deletedProperty.ID = t.testPropertyID

	t.mockSvc.EXPECT().Delete(t.testPropertyID, t.testUserID).Return(&deletedProperty, nil)

	t.handler.HandleDeleteProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), t.testPropertyID, actual.ID)
	assert.Nil(t.T(), err)
}

func (t *propertyHandlerTestSuite) TestDelete_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/properties/"+t.testPropertyID.String(),
		nil,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleDeleteProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestDelete_ServiceFailedDeleting() {
	errMsg := "service failed deleting property"
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/properties/"+t.testPropertyID.String(),
		nil,
		nil,
		nuuid.From(t.testPropertyID),
	)

	t.mockSvc.EXPECT().Delete(t.testPropertyID, t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleDeleteProperty(rr, req)

	actual, err := t.parseOutputToProperty(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestCreateValue_Normal() {
	input := t.getNewPropertyValueInput(nuuid.NUUID{Valid: false}, nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/values",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedResult := model.NewPropertyValueFromInput(input, input.PropertyID, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().CreateValue(gomock.Any(), t.testUserID).Return(&expectedResult, nil)

	t.handler.HandleCreatePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusCreated, rr.Result().StatusCode)
	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), expected.ID, actual.ID)
	assert.Equal(t.T(), expected.PropertyID, actual.PropertyID)
	assert.Equal(t.T(), expected.Date.Time().Unix(), actual.Date.Time().Unix())
	assert.Equal(t.T(), expected.Value, actual.Value)
	assert.NotNil(t.T(), actual.Created)
	assert.NotNil(t.T(), actual.CreatedBy)
	assert.False(t.T(), actual.Updated.Valid)
	assert.False(t.T(), actual.UpdatedBy.Valid)
	assert.False(t.T(), actual.Deleted.Valid)
	assert.False(t.T(), actual.DeletedBy.Valid)
}

func (t *propertyHandlerTestSuite) TestCreateValue_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/values",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleCreatePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestCreateValue_ServiceFailedCreatingValue() {
	errMsg := "service failed creating property value"
	input := t.getNewPropertyValueInput(nuuid.NUUID{Valid: false}, nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/values",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().CreateValue(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleCreatePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestGetValueByID_Normal() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/values/"+t.testPropertyValueID.String(),
		nil,
		nil,
		nuuid.From(t.testPropertyValueID),
	)

	input := t.getNewPropertyValueInput(nuuid.From(t.testPropertyValueID), nuuid.From(t.testPropertyID))
	expectedResult := model.NewPropertyValueFromInput(input, t.testPropertyID, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().GetValueByID(t.testPropertyValueID).Return(&expectedResult, nil)

	t.handler.HandleGetPropertyValueByID(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.ID, actual.ID)
	assert.Equal(t.T(), expected.PropertyID, actual.PropertyID)
	assert.Equal(t.T(), expected.Date.Time().Unix(), actual.Date.Time().Unix())
	assert.Equal(t.T(), expected.Value, actual.Value)
	assert.Equal(t.T(), expected.Created.Time().Unix(), actual.Created.Time().Unix())
	assert.Equal(t.T(), expected.CreatedBy, actual.CreatedBy)
	assert.Equal(t.T(), expected.Updated, actual.Updated)
	assert.Equal(t.T(), expected.UpdatedBy, actual.UpdatedBy)
	assert.Equal(t.T(), expected.Deleted, actual.Deleted)
	assert.Equal(t.T(), expected.DeletedBy, actual.DeletedBy)
}

func (t *propertyHandlerTestSuite) TestGetValueByID_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/values/"+t.testPropertyValueID.String(),
		nil,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleGetPropertyValueByID(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestGetValueByID_ServiceFailedResolving() {
	errMsg := "service failed resolving property value"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/properties/values/"+t.testPropertyValueID.String(),
		nil,
		nil,
		nuuid.From(t.testPropertyValueID),
	)

	t.mockSvc.EXPECT().GetValueByID(t.testPropertyValueID).Return(nil, errors.New(errMsg))

	t.handler.HandleGetPropertyValueByID(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestGetValueByFilter_Normal() {
	keyword := "test keyword"
	input := model.PropertyValueFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/values/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedPropertyValues := []model.PropertyValue{}
	vv1 := model.NewPropertyValueFromInput(t.getNewPropertyValueInput(nuuid.NUUID{}, nuuid.From(t.testPropertyID)), t.testPropertyID, t.testUserID)
	vv2 := model.NewPropertyValueFromInput(t.getNewPropertyValueInput(nuuid.NUUID{}, nuuid.From(t.testPropertyID)), t.testPropertyID, t.testUserID)
	expectedPropertyValues = append(expectedPropertyValues, vv1)
	expectedPropertyValues = append(expectedPropertyValues, vv2)
	expectedPageInfo := model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 1,
		PageCount:  1,
	}

	t.mockSvc.EXPECT().GetValuesByFilter(input).Return(expectedPropertyValues, expectedPageInfo, nil)

	t.handler.HandleGetPropertyValueByFilter(rr, req)

	propertyValues, pageInfo, err := t.parseOutputToPropertyValuePage(rr)

	assert.Nil(t.T(), err)

	assert.Equal(t.T(), len(expectedPropertyValues), len(propertyValues))
	assert.Equal(t.T(), expectedPropertyValues[0].ID, propertyValues[0].ID)
	assert.Equal(t.T(), expectedPropertyValues[1].ID, propertyValues[1].ID)

	assert.Equal(t.T(), 1, pageInfo.Page)
}

func (t *propertyHandlerTestSuite) TestGetValueByFilter_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/values/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetPropertyValueByFilter(rr, req)

	properties, pageInfo, err := t.parseOutputToPropertyValuePage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(properties))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *propertyHandlerTestSuite) TestGetValueByFilter_ServiceFailedResolving() {
	errMsg := "service failed resolving property values"
	keyword := "test keyword"
	input := model.PropertyValueFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/properties/values/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().GetValuesByFilter(input).Return([]model.PropertyValue{}, model.PageInfoOutput{}, errors.New(errMsg))

	t.handler.HandleGetPropertyValueByFilter(rr, req)

	vahicleValues, pageInfo, err := t.parseOutputToPropertyValuePage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(vahicleValues))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *propertyHandlerTestSuite) TestUpdateValue_Normal() {
	input := t.getNewPropertyValueInput(nuuid.From(t.testPropertyValueID), nuuid.From(t.testPropertyID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/values/"+t.testPropertyValueID.String(),
		input,
		nil,
		nuuid.From(t.testPropertyValueID),
	)

	updatedPropertyValue := model.NewPropertyValueFromInput(input, t.testPropertyID, t.testUserID)

	t.mockSvc.EXPECT().UpdateValue(gomock.Any(), t.testUserID).Return(&updatedPropertyValue, nil)

	t.handler.HandleUpdatePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *propertyHandlerTestSuite) TestUpdateValue_FailedGettingIDFromRequest() {
	input := t.getNewPropertyValueInput(nuuid.From(t.testPropertyValueID), nuuid.From(t.testPropertyID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/values/"+t.testPropertyValueID.String(),
		input,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleUpdatePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestUpdateValue_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/values/"+t.testPropertyValueID.String(),
		input,
		nil,
		nuuid.From(t.testPropertyValueID),
	)

	t.handler.HandleUpdatePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestUpdateValue_MismatchedID() {
	input := t.getNewPropertyValueInput(nuuid.NUUID{}, nuuid.NUUID{})
	newID, _ := uuid.NewV7()
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/"+t.testPropertyValueID.String(),
		input,
		nil,
		nuuid.From(newID),
	)

	t.handler.HandleUpdatePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "id mismatch")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestUpdateValue_ServiceFailedUpdating() {
	errMsg := "failed updating property value"
	input := t.getNewPropertyValueInput(nuuid.From(t.testPropertyValueID), nuuid.From(t.testPropertyID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/properties/"+t.testPropertyValueID.String(),
		input,
		nil,
		nuuid.From(t.testPropertyValueID),
	)

	t.mockSvc.EXPECT().UpdateValue(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleUpdatePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestDeleteValue_Normal() {
	input := t.getNewPropertyValueInput(nuuid.From(t.testPropertyValueID), nuuid.From(t.testPropertyID))
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/properties/values/"+t.testPropertyValueID.String(),
		nil,
		nil,
		nuuid.From(t.testPropertyValueID),
	)

	deletedPropertyValue := model.NewPropertyValueFromInput(input, t.testPropertyID, t.testUserID)
	deletedPropertyValue.ID = t.testPropertyID

	t.mockSvc.EXPECT().DeleteValue(t.testPropertyValueID, t.testUserID).Return(&deletedPropertyValue, nil)

	t.handler.HandleDeletePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), t.testPropertyID, actual.ID)
	assert.Nil(t.T(), err)
}

func (t *propertyHandlerTestSuite) TestDeleteValue_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/properties/values/"+t.testPropertyValueID.String(),
		nil,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleDeletePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *propertyHandlerTestSuite) TestDeleteValue_ServiceFailedDeleting() {
	errMsg := "service failed deleting property value"
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/properties/values/"+t.testPropertyValueID.String(),
		nil,
		nil,
		nuuid.From(t.testPropertyValueID),
	)

	t.mockSvc.EXPECT().DeleteValue(t.testPropertyValueID, t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleDeletePropertyValue(rr, req)

	actual, err := t.parseOutputToPropertyValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}
