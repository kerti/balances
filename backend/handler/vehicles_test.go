package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
	ctrl          *gomock.Controller
	handler       handler.Vehicle
	mockSvc       *mock_service.MockVehicle
	testUserID    uuid.UUID
	testVehicleID uuid.UUID
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
	t.handler.Startup()
}

func (t *vehicleHandlerTestSuite) TearDownTest() {
	t.handler.Shutdown()
	t.ctrl.Finish()
}

func (t *vehicleHandlerTestSuite) getNewRequestWithContext(input interface{}, method, path string) (recorder *httptest.ResponseRecorder, request *http.Request) {
	jsonBody, err := json.Marshal(input)
	if err != nil {
		t.T().Fatal(err)
	}

	reqBody := bytes.NewBuffer(jsonBody)
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

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

func (t *vehicleHandlerTestSuite) parseOutputToVehicle(rr *httptest.ResponseRecorder, actual *model.VehicleOutput) {
	// read the response
	var response response.BaseResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.T().Fatal(err)
	}

	// marshal the data to JSON
	actualMap := (*response.Data).(map[string]interface{})
	jsonBytes, err := json.Marshal(actualMap)
	if err != nil {
		t.T().Fatal(err)
	}

	// unmarshal back to the expected object
	err = json.Unmarshal(jsonBytes, actual)
	if err != nil {
		t.T().Fatal(err)
	}
}

func (t *vehicleHandlerTestSuite) parseOutputToError(rr *httptest.ResponseRecorder) (actual *failure.Failure) {
	// read the response
	var response response.BaseResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.T().Fatal(err)
	}

	actual = response.Error
	return
}

func (t *vehicleHandlerTestSuite) TestCreate_Normal() {
	input := t.getNewVehicleInput(nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		input,
		http.MethodPost,
		"/vehicles",
	)

	expectedResult := model.NewVehicleFromInput(input, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().Create(gomock.Any(), t.testUserID).Return(&expectedResult, nil)

	t.handler.HandleCreateVehicle(rr, req)

	var actual model.VehicleOutput
	t.parseOutputToVehicle(rr, &actual)

	assert.Equal(t.T(), http.StatusCreated, rr.Result().StatusCode)
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
	input := time.Now()
	rr, req := t.getNewRequestWithContext(
		input,
		http.MethodPost,
		"/vehicles",
	)

	t.handler.HandleCreateVehicle(rr, req)

	actual := t.parseOutputToError(rr)

	assert.Equal(t.T(), failure.CodeBadRequest, actual.Code)
	// TODO: specify this
	assert.Nil(t.T(), actual.Entity)
	assert.Contains(t.T(), actual.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), actual.Operation)
}

func (t *vehicleHandlerTestSuite) TestCreate_ServiceFailedCreating() {
	errMsg := "service failed creating vehicle"
	input := t.getNewVehicleInput(nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		input,
		http.MethodPost,
		"/vehicles",
	)

	t.mockSvc.EXPECT().Create(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleCreateVehicle(rr, req)

	actual := t.parseOutputToError(rr)

	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), actual.Entity)
	assert.Contains(t.T(), actual.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), actual.Operation)
}
