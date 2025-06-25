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

type vehicleHandlerTestSuite struct {
	suite.Suite
	ctrl               *gomock.Controller
	handler            handler.Vehicle
	mockSvc            *mock_service.MockVehicle
	testUserID         uuid.UUID
	testVehicleID      uuid.UUID
	testVehicleValueID uuid.UUID
}

func TestVehicleHandler(t *testing.T) {
	suite.Run(t, new(vehicleHandlerTestSuite))
}

func (t *vehicleHandlerTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockSvc = mock_service.NewMockVehicle(t.ctrl)
	t.handler = &handler.VehicleImpl{
		Service: t.mockSvc,
	}
	t.testUserID, _ = uuid.NewV7()
	t.testVehicleID, _ = uuid.NewV7()
	t.testVehicleValueID, _ = uuid.NewV7()
	t.handler.Startup()
}

func (t *vehicleHandlerTestSuite) TearDownTest() {
	t.handler.Shutdown()
	t.ctrl.Finish()
}

func (t *vehicleHandlerTestSuite) getNewRequestWithContext(method, path string, input any, formParams *map[string]string, routeVarId nuuid.NUUID) (recorder *httptest.ResponseRecorder, request *http.Request) {
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

func (t *vehicleHandlerTestSuite) getNewVehicleInput(id nuuid.NUUID) model.VehicleInput {
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

func (t *vehicleHandlerTestSuite) getNewVehicleValueInput(id, vehicleID nuuid.NUUID) model.VehicleValueInput {
	vvi := model.VehicleValueInput{}

	if id.Valid {
		vvi.ID = id.UUID
	} else {
		vvi.ID = t.testVehicleValueID
	}

	if vehicleID.Valid {
		vvi.VehicleID = vehicleID.UUID
	} else {
		vvi.VehicleID = t.testVehicleID
	}

	vvi.Date = cachetime.CacheTime(time.Now().AddDate(0, 0, -1))
	vvi.Value = float64(50000)

	return vvi
}

func (t *vehicleHandlerTestSuite) parseOutputToVehicle(rr *httptest.ResponseRecorder) (actual *model.VehicleOutput, fail *failure.Failure) {
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

func (t *vehicleHandlerTestSuite) parseOutputToVehicleValue(rr *httptest.ResponseRecorder) (actual *model.VehicleValueOutput, fail *failure.Failure) {
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

func (t *vehicleHandlerTestSuite) parseOutputToVehiclePage(rr *httptest.ResponseRecorder) (items []model.VehicleOutput, pageInfo model.PageInfoOutput, fail *failure.Failure) {
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

		//convert interface{} to []model.VehicleOutput
		actualSlice := (actual.Items).([]any)
		for _, vehicleInterface := range actualSlice {
			vehicleMap := (vehicleInterface).(map[string]any)
			vehicleJsonBytes, err := json.Marshal(vehicleMap)
			if err != nil {
				t.T().Fatal(err)
			}
			var actualVehicle model.VehicleOutput
			err = json.Unmarshal(vehicleJsonBytes, &actualVehicle)
			if err != nil {
				t.T().Fatal(err)
			}
			items = append(items, actualVehicle)
		}

		pageInfo = actual.PageInfo
	}

	if response.Error != nil {
		fail = response.Error
	}

	return
}

func (t *vehicleHandlerTestSuite) parseOutputToVehicleValuePage(rr *httptest.ResponseRecorder) (items []model.VehicleValueOutput, pageInfo model.PageInfoOutput, fail *failure.Failure) {
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

		//convert interface{} to []model.VehicleOutput
		actualSlice := (actual.Items).([]any)
		for _, vehicleValueInterface := range actualSlice {
			vehicleValueMap := (vehicleValueInterface).(map[string]any)
			vehicleValueJsonBytes, err := json.Marshal(vehicleValueMap)
			if err != nil {
				t.T().Fatal(err)
			}
			var actualVehicleValue model.VehicleValueOutput
			err = json.Unmarshal(vehicleValueJsonBytes, &actualVehicleValue)
			if err != nil {
				t.T().Fatal(err)
			}
			items = append(items, actualVehicleValue)
		}

		pageInfo = actual.PageInfo
	}

	if response.Error != nil {
		fail = response.Error
	}

	return
}

func (t *vehicleHandlerTestSuite) TestCreate_Normal() {
	input := t.getNewVehicleInput(nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedResult := model.NewVehicleFromInput(input, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().Create(gomock.Any(), t.testUserID).Return(&expectedResult, nil)

	t.handler.HandleCreateVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusCreated, rr.Result().StatusCode)
	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), expected.Name, actual.Name)
	assert.Equal(t.T(), expected.Make, actual.Make)
	assert.Equal(t.T(), expected.Model, actual.Model)
	assert.Equal(t.T(), expected.Year, actual.Year)
	assert.Equal(t.T(), expected.TitleHolder, actual.TitleHolder)
	assert.Equal(t.T(), expected.LicensePlateNumber, actual.LicensePlateNumber)
	assert.Equal(t.T(), expected.PurchaseDate.Time().Unix(), actual.PurchaseDate.Time().Unix())
	assert.Equal(t.T(), expected.InitialValue, actual.InitialValue)
	assert.Equal(t.T(), expected.InitialValueDate.Time().Unix(), actual.InitialValueDate.Time().Unix())
	assert.Equal(t.T(), expected.CurrentValue, actual.CurrentValue)
	assert.Equal(t.T(), expected.CurrentValueDate.Time().Unix(), actual.CurrentValueDate.Time().Unix())
	assert.Equal(t.T(), expected.AnnualDepreciationPercent, actual.AnnualDepreciationPercent)
	assert.Equal(t.T(), expected.Status, actual.Status)
	assert.NotNil(t.T(), actual.Created)
	assert.NotNil(t.T(), actual.CreatedBy)
	assert.False(t.T(), actual.Updated.Valid)
	assert.False(t.T(), actual.UpdatedBy.Valid)
	assert.False(t.T(), actual.Deleted.Valid)
	assert.False(t.T(), actual.DeletedBy.Valid)
}

func (t *vehicleHandlerTestSuite) TestCreate_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleCreateVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestCreate_ServiceFailedCreating() {
	errMsg := "service failed creating vehicle"
	input := t.getNewVehicleInput(nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().Create(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleCreateVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestGetByID_Normal_NoParams() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		nil,
		nuuid.From(t.testVehicleID),
	)

	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	expectedResult := model.NewVehicleFromInput(input, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().GetByID(t.testVehicleID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).Return(&expectedResult, nil)

	t.handler.HandleGetVehicleByID(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.Name, actual.Name)
	assert.Equal(t.T(), expected.Make, actual.Make)
	assert.Equal(t.T(), expected.Model, actual.Model)
	assert.Equal(t.T(), expected.Year, actual.Year)
	assert.Equal(t.T(), expected.TitleHolder, actual.TitleHolder)
	assert.Equal(t.T(), expected.LicensePlateNumber, actual.LicensePlateNumber)
	assert.Equal(t.T(), expected.PurchaseDate.Time().Unix(), actual.PurchaseDate.Time().Unix())
	assert.Equal(t.T(), expected.InitialValue, actual.InitialValue)
	assert.Equal(t.T(), expected.InitialValueDate.Time().Unix(), actual.InitialValueDate.Time().Unix())
	assert.Equal(t.T(), expected.CurrentValue, actual.CurrentValue)
	assert.Equal(t.T(), expected.CurrentValueDate.Time().Unix(), actual.CurrentValueDate.Time().Unix())
	assert.Equal(t.T(), expected.AnnualDepreciationPercent, actual.AnnualDepreciationPercent)
	assert.Equal(t.T(), expected.Status, actual.Status)
	assert.Equal(t.T(), expected.Created.Time().Unix(), actual.Created.Time().Unix())
	assert.Equal(t.T(), expected.CreatedBy, actual.CreatedBy)
	assert.Equal(t.T(), expected.Updated, actual.Updated)
	assert.Equal(t.T(), expected.UpdatedBy, actual.UpdatedBy)
	assert.Equal(t.T(), expected.Deleted, actual.Deleted)
	assert.Equal(t.T(), expected.DeletedBy, actual.DeletedBy)
}

func (t *vehicleHandlerTestSuite) TestGetByID_FailedParsingID() {
	formParams := make(map[string]string)
	formParams["id"] = t.testVehicleID.String() + "123"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/"+t.testVehicleID.String()+"123",
		formParams,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetVehicleByID(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestGetByID_Normal_WithValues() {
	formParams := make(map[string]string)
	formParams["withValues"] = "true"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		&formParams,
		nuuid.From(t.testVehicleID),
	)

	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	expectedResult := model.NewVehicleFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().GetByID(t.testVehicleID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).Return(&expectedResult, nil)

	t.handler.HandleGetVehicleByID(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *vehicleHandlerTestSuite) TestGetByID_Normal_WithValueStartDate() {
	startDate := time.Unix(0, time.Now().AddDate(0, 0, -1).UnixMilli()*int64(time.Millisecond))
	formParams := make(map[string]string)
	formParams["valueStartDate"] = strconv.FormatInt(startDate.UnixMilli(), 10)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		&formParams,
		nuuid.From(t.testVehicleID),
	)

	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	expectedResult := model.NewVehicleFromInput(input, t.testUserID)

	var nStartDate cachetime.NCacheTime
	nStartDate.Scan(startDate)
	t.mockSvc.EXPECT().GetByID(t.testVehicleID, false, nStartDate, cachetime.NCacheTime{}, nil).Return(&expectedResult, nil)

	t.handler.HandleGetVehicleByID(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *vehicleHandlerTestSuite) TestGetByID_Normal_WithValueEndDate() {
	endDate := time.Unix(0, time.Now().AddDate(0, 0, -1).UnixMilli()*int64(time.Millisecond))
	formParams := make(map[string]string)
	formParams["valueEndDate"] = strconv.FormatInt(endDate.UnixMilli(), 10)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		&formParams,
		nuuid.From(t.testVehicleID),
	)

	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	expectedResult := model.NewVehicleFromInput(input, t.testUserID)

	var nEndDate cachetime.NCacheTime
	nEndDate.Scan(endDate)
	t.mockSvc.EXPECT().GetByID(t.testVehicleID, false, cachetime.NCacheTime{}, nEndDate, nil).Return(&expectedResult, nil)

	t.handler.HandleGetVehicleByID(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *vehicleHandlerTestSuite) TestGetByID_Normal_WithPageSize() {
	pageSize := 10
	formParams := make(map[string]string)
	formParams["pageSize"] = strconv.Itoa(pageSize)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		&formParams,
		nuuid.From(t.testVehicleID),
	)

	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	expectedResult := model.NewVehicleFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().GetByID(t.testVehicleID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, &pageSize).Return(&expectedResult, nil)

	t.handler.HandleGetVehicleByID(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *vehicleHandlerTestSuite) TestGetByID_ServiceFailedResolving() {
	errMsg := "service failed resolving"
	formParams := make(map[string]string)
	formParams["id"] = t.testVehicleID.String()
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		&formParams,
		nuuid.From(t.testVehicleID),
	)

	t.mockSvc.EXPECT().GetByID(t.testVehicleID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).Return(nil, errors.New(errMsg))

	t.handler.HandleGetVehicleByID(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestGetByFilter_Normal() {
	keyword := "test keyword"
	input := model.VehicleFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedVehicles := []model.Vehicle{}
	v1 := model.NewVehicleFromInput(t.getNewVehicleInput(nuuid.NUUID{}), t.testUserID)
	v2 := model.NewVehicleFromInput(t.getNewVehicleInput(nuuid.NUUID{}), t.testUserID)
	expectedVehicles = append(expectedVehicles, v1)
	expectedVehicles = append(expectedVehicles, v2)
	expectedPageInfo := model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 1,
		PageCount:  1,
	}

	t.mockSvc.EXPECT().GetByFilter(input).Return(expectedVehicles, expectedPageInfo, nil)

	t.handler.HandleGetVehicleByFilter(rr, req)

	vehicles, pageInfo, err := t.parseOutputToVehiclePage(rr)

	assert.Nil(t.T(), err)

	assert.Equal(t.T(), len(expectedVehicles), len(vehicles))
	assert.Equal(t.T(), expectedVehicles[0].ID, vehicles[0].ID)
	assert.Equal(t.T(), expectedVehicles[1].ID, vehicles[1].ID)

	assert.Equal(t.T(), 1, pageInfo.Page)
}

func (t *vehicleHandlerTestSuite) TestGetByFilter_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetVehicleByFilter(rr, req)

	vehicles, pageInfo, err := t.parseOutputToVehiclePage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(vehicles))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *vehicleHandlerTestSuite) TestGetByFilter_ServiceFailedResolving() {
	errMsg := "failed resolving vehicles by filter"
	keyword := "test keyword"
	input := model.VehicleFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().GetByFilter(input).Return([]model.Vehicle{}, model.PageInfoOutput{}, errors.New(errMsg))

	t.handler.HandleGetVehicleByFilter(rr, req)

	vehicles, pageInfo, err := t.parseOutputToVehiclePage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(vehicles))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *vehicleHandlerTestSuite) TestUpdate_Normal() {
	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/"+t.testVehicleID.String(),
		input,
		nil,
		nuuid.From(t.testVehicleID),
	)

	updatedVehicle := model.NewVehicleFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().Update(gomock.Any(), t.testUserID).Return(&updatedVehicle, nil)

	t.handler.HandleUpdateVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *vehicleHandlerTestSuite) TestUpdate_FailedGettingIDFromRequest() {
	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/"+t.testVehicleID.String(),
		input,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleUpdateVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestUpdate_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/"+t.testVehicleID.String(),
		input,
		nil,
		nuuid.From(t.testVehicleID),
	)

	t.handler.HandleUpdateVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestUpdate_MismatchedID() {
	input := t.getNewVehicleInput(nuuid.NUUID{})
	newID, _ := uuid.NewV7()
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/"+t.testVehicleID.String(),
		input,
		nil,
		nuuid.From(newID),
	)

	t.handler.HandleUpdateVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "id mismatch")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestUpdate_ServiceFailedUpdating() {
	errMsg := "failed updating vehicle"
	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/"+t.testVehicleID.String(),
		input,
		nil,
		nuuid.From(t.testVehicleID),
	)

	t.mockSvc.EXPECT().Update(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleUpdateVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestDelete_Normal() {
	input := t.getNewVehicleInput(nuuid.From(t.testVehicleID))
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		nil,
		nuuid.From(t.testVehicleID),
	)

	deletedVehicle := model.NewVehicleFromInput(input, t.testUserID)
	deletedVehicle.ID = t.testVehicleID

	t.mockSvc.EXPECT().Delete(t.testVehicleID, t.testUserID).Return(&deletedVehicle, nil)

	t.handler.HandleDeleteVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), t.testVehicleID, actual.ID)
	assert.Nil(t.T(), err)
}

func (t *vehicleHandlerTestSuite) TestDelete_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleDeleteVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestDelete_ServiceFailedDeleting() {
	errMsg := "service failed deleting vehicle"
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/vehicles/"+t.testVehicleID.String(),
		nil,
		nil,
		nuuid.From(t.testVehicleID),
	)

	t.mockSvc.EXPECT().Delete(t.testVehicleID, t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleDeleteVehicle(rr, req)

	actual, err := t.parseOutputToVehicle(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestCreateValue_Normal() {
	input := t.getNewVehicleValueInput(nuuid.NUUID{Valid: false}, nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/values",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedResult := model.NewVehicleValueFromInput(input, input.VehicleID, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().CreateValue(gomock.Any(), t.testUserID).Return(&expectedResult, nil)

	t.handler.HandleCreateVehicleValue(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusCreated, rr.Result().StatusCode)
	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), expected.ID, actual.ID)
	assert.Equal(t.T(), expected.VehicleID, actual.VehicleID)
	assert.Equal(t.T(), expected.Date.Time().Unix(), actual.Date.Time().Unix())
	assert.Equal(t.T(), expected.Value, actual.Value)
	assert.NotNil(t.T(), actual.Created)
	assert.NotNil(t.T(), actual.CreatedBy)
	assert.False(t.T(), actual.Updated.Valid)
	assert.False(t.T(), actual.UpdatedBy.Valid)
	assert.False(t.T(), actual.Deleted.Valid)
	assert.False(t.T(), actual.DeletedBy.Valid)
}

func (t *vehicleHandlerTestSuite) TestCreateValue_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/values",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleCreateVehicleValue(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestCreateValue_ServiceFailedCreatingValue() {
	errMsg := "service failed creating vehicle value"
	input := t.getNewVehicleValueInput(nuuid.NUUID{Valid: false}, nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/values",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().CreateValue(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleCreateVehicleValue(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestGetValueByID_Normal() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/values/"+t.testVehicleValueID.String(),
		nil,
		nil,
		nuuid.From(t.testVehicleValueID),
	)

	input := t.getNewVehicleValueInput(nuuid.From(t.testVehicleValueID), nuuid.From(t.testVehicleID))
	expectedResult := model.NewVehicleValueFromInput(input, t.testVehicleID, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().GetValueByID(t.testVehicleValueID).Return(&expectedResult, nil)

	t.handler.HandleGetVehicleValueByID(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.ID, actual.ID)
	assert.Equal(t.T(), expected.VehicleID, actual.VehicleID)
	assert.Equal(t.T(), expected.Date.Time().Unix(), actual.Date.Time().Unix())
	assert.Equal(t.T(), expected.Value, actual.Value)
	assert.Equal(t.T(), expected.Created.Time().Unix(), actual.Created.Time().Unix())
	assert.Equal(t.T(), expected.CreatedBy, actual.CreatedBy)
	assert.Equal(t.T(), expected.Updated, actual.Updated)
	assert.Equal(t.T(), expected.UpdatedBy, actual.UpdatedBy)
	assert.Equal(t.T(), expected.Deleted, actual.Deleted)
	assert.Equal(t.T(), expected.DeletedBy, actual.DeletedBy)
}

func (t *vehicleHandlerTestSuite) TestGetValueByID_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/values/"+t.testVehicleValueID.String(),
		nil,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleGetVehicleValueByID(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestGetValueByID_ServiceFailedResolving() {
	errMsg := "service failed resolving vehicle value"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/vehicles/values/"+t.testVehicleValueID.String(),
		nil,
		nil,
		nuuid.From(t.testVehicleValueID),
	)

	t.mockSvc.EXPECT().GetValueByID(t.testVehicleValueID).Return(nil, errors.New(errMsg))

	t.handler.HandleGetVehicleValueByID(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestGetValueByFilter_Normal() {
	keyword := "test keyword"
	input := model.VehicleValueFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/values/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedVehicleValues := []model.VehicleValue{}
	vv1 := model.NewVehicleValueFromInput(t.getNewVehicleValueInput(nuuid.NUUID{}, nuuid.From(t.testVehicleID)), t.testVehicleID, t.testUserID)
	vv2 := model.NewVehicleValueFromInput(t.getNewVehicleValueInput(nuuid.NUUID{}, nuuid.From(t.testVehicleID)), t.testVehicleID, t.testUserID)
	expectedVehicleValues = append(expectedVehicleValues, vv1)
	expectedVehicleValues = append(expectedVehicleValues, vv2)
	expectedPageInfo := model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 1,
		PageCount:  1,
	}

	t.mockSvc.EXPECT().GetValuesByFilter(input).Return(expectedVehicleValues, expectedPageInfo, nil)

	t.handler.HandleGetVehicleValueByFilter(rr, req)

	vehicleValues, pageInfo, err := t.parseOutputToVehicleValuePage(rr)

	assert.Nil(t.T(), err)

	assert.Equal(t.T(), len(expectedVehicleValues), len(vehicleValues))
	assert.Equal(t.T(), expectedVehicleValues[0].ID, vehicleValues[0].ID)
	assert.Equal(t.T(), expectedVehicleValues[1].ID, vehicleValues[1].ID)

	assert.Equal(t.T(), 1, pageInfo.Page)
}

func (t *vehicleHandlerTestSuite) TestGetValueByFilter_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/values/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetVehicleValueByFilter(rr, req)

	vehicles, pageInfo, err := t.parseOutputToVehicleValuePage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(vehicles))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *vehicleHandlerTestSuite) TestGetValueByFilter_ServiceFailedResolving() {
	errMsg := "service failed resolving vehicle values"
	keyword := "test keyword"
	input := model.VehicleValueFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/vehicles/values/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().GetValuesByFilter(input).Return([]model.VehicleValue{}, model.PageInfoOutput{}, errors.New(errMsg))

	t.handler.HandleGetVehicleValueByFilter(rr, req)

	vahicleValues, pageInfo, err := t.parseOutputToVehicleValuePage(rr)

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

func (t *vehicleHandlerTestSuite) TestUpdateValue_Normal() {
	input := t.getNewVehicleValueInput(nuuid.From(t.testVehicleValueID), nuuid.From(t.testVehicleID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/values/"+t.testVehicleValueID.String(),
		input,
		nil,
		nuuid.From(t.testVehicleValueID),
	)

	updatedVehicleValue := model.NewVehicleValueFromInput(input, t.testVehicleID, t.testUserID)

	t.mockSvc.EXPECT().UpdateValue(gomock.Any(), t.testUserID).Return(&updatedVehicleValue, nil)

	t.handler.HandleUpdateVehicleValue(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *vehicleHandlerTestSuite) TestUpdateValue_FailedGettingIDFromRequest() {
	input := t.getNewVehicleValueInput(nuuid.From(t.testVehicleValueID), nuuid.From(t.testVehicleID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/values/"+t.testVehicleValueID.String(),
		input,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleUpdateVehicleValue(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestUpdateValue_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/values/"+t.testVehicleValueID.String(),
		input,
		nil,
		nuuid.From(t.testVehicleValueID),
	)

	t.handler.HandleUpdateVehicleValue(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestUpdateValue_MismatchedID() {
	input := t.getNewVehicleValueInput(nuuid.NUUID{}, nuuid.NUUID{})
	newID, _ := uuid.NewV7()
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/"+t.testVehicleValueID.String(),
		input,
		nil,
		nuuid.From(newID),
	)

	t.handler.HandleUpdateVehicleValue(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "id mismatch")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *vehicleHandlerTestSuite) TestUpdateValue_ServiceFailedUpdating() {
	errMsg := "failed updating vehicle value"
	input := t.getNewVehicleValueInput(nuuid.From(t.testVehicleValueID), nuuid.From(t.testVehicleID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/vehicles/"+t.testVehicleValueID.String(),
		input,
		nil,
		nuuid.From(t.testVehicleValueID),
	)

	t.mockSvc.EXPECT().UpdateValue(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleUpdateVehicleValue(rr, req)

	actual, err := t.parseOutputToVehicleValue(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}
