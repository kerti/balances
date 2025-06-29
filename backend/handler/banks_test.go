package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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

type bankAccountsHandlerTestSuite struct {
	suite.Suite
	ctrl                     *gomock.Controller
	handler                  handler.BankAccount
	mockSvc                  *mock_service.MockBankAccount
	testUserID               uuid.UUID
	testBankAccountID        uuid.UUID
	testBankAccountBalanceID uuid.UUID
}

func TestBankAccountHandler(t *testing.T) {
	suite.Run(t, new(bankAccountsHandlerTestSuite))
}

func (t *bankAccountsHandlerTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockSvc = mock_service.NewMockBankAccount(t.ctrl)
	t.handler = &handler.BankAccountImpl{
		Service: t.mockSvc,
	}
	t.testUserID, _ = uuid.NewV7()
	t.testBankAccountID, _ = uuid.NewV7()
	t.testBankAccountBalanceID, _ = uuid.NewV7()
	t.handler.Startup()
}

func (t *bankAccountsHandlerTestSuite) TearDownTest() {
	t.handler.Shutdown()
	t.ctrl.Finish()
}

func (t *bankAccountsHandlerTestSuite) getNewRequestWithContext(method, path string, input any, formParams *map[string]string, routeVarId nuuid.NUUID) (recorder *httptest.ResponseRecorder, request *http.Request) {
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

func (t *bankAccountsHandlerTestSuite) getNewBankAccountInput(id nuuid.NUUID) model.BankAccountInput {
	acc := model.BankAccountInput{}

	if id.Valid {
		acc.ID = id.UUID
	} else {
		acc.ID = t.testBankAccountID
	}

	acc.AccountName = "John's Account"
	acc.BankName = "First National Bank"
	acc.AccountHolderName = "John Fitzgerald Doe"
	acc.AccountNumber = "123-456-7890"
	acc.LastBalance = float64(10000)
	acc.LastBalanceDate = cachetime.CacheTime(time.Now())
	acc.Status = model.BankAccountStatusActive

	return acc
}

func (t *bankAccountsHandlerTestSuite) getNewBankAccountBalanceInput(id, bankAccountID nuuid.NUUID) model.BankAccountBalanceInput {
	bbi := model.BankAccountBalanceInput{}

	if id.Valid {
		bbi.ID = id.UUID
	} else {
		bbi.ID = t.testBankAccountBalanceID
	}

	if bankAccountID.Valid {
		bbi.BankAccountID = bankAccountID.UUID
	} else {
		bbi.BankAccountID = t.testBankAccountBalanceID
	}

	bbi.Date = cachetime.CacheTime(time.Now().AddDate(0, 0, -1))
	bbi.Balance = float64(50000)

	return bbi
}

func (t *bankAccountsHandlerTestSuite) parseOutputToBankAccount(rr *httptest.ResponseRecorder) (actual *model.BankAccountOutput, fail *failure.Failure) {
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

func (t *bankAccountsHandlerTestSuite) parseOutputToBankAccountBalance(rr *httptest.ResponseRecorder) (actual *model.BankAccountBalanceOutput, fail *failure.Failure) {
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

func (t *bankAccountsHandlerTestSuite) parseOutputToBankAccountPage(rr *httptest.ResponseRecorder) (items []model.BankAccountOutput, pageInfo model.PageInfoOutput, fail *failure.Failure) {
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

		//convert interface{} to []model.BankAccountOutput
		actualSlice := (actual.Items).([]any)
		for _, vehicleInterface := range actualSlice {
			bankAccountMap := (vehicleInterface).(map[string]any)
			bankAccountJsonBytes, err := json.Marshal(bankAccountMap)
			if err != nil {
				t.T().Fatal(err)
			}
			var actualBankAccount model.BankAccountOutput
			err = json.Unmarshal(bankAccountJsonBytes, &actualBankAccount)
			if err != nil {
				t.T().Fatal(err)
			}
			items = append(items, actualBankAccount)
		}

		pageInfo = actual.PageInfo
	}

	if response.Error != nil {
		fail = response.Error
	}

	return
}

func (t *bankAccountsHandlerTestSuite) parseOutputToVehicleValuePage(rr *httptest.ResponseRecorder) (items []model.BankAccountBalanceOutput, pageInfo model.PageInfoOutput, fail *failure.Failure) {
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
			bankAccountBalanceMap := (vehicleValueInterface).(map[string]any)
			bankAccountBalanceJsonBytes, err := json.Marshal(bankAccountBalanceMap)
			if err != nil {
				t.T().Fatal(err)
			}
			var actualBankAccountBalance model.BankAccountBalanceOutput
			err = json.Unmarshal(bankAccountBalanceJsonBytes, &actualBankAccountBalance)
			if err != nil {
				t.T().Fatal(err)
			}
			items = append(items, actualBankAccountBalance)
		}

		pageInfo = actual.PageInfo
	}

	if response.Error != nil {
		fail = response.Error
	}

	return
}

func (t *bankAccountsHandlerTestSuite) TestCreate_Normal() {
	input := t.getNewBankAccountInput(nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedResult := model.NewBankAccountFromInput(input, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().Create(gomock.Any(), t.testUserID).Return(&expectedResult, nil)

	t.handler.HandleCreateBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusCreated, rr.Result().StatusCode)
	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), expected.AccountName, actual.AccountName)
	assert.Equal(t.T(), expected.BankName, actual.BankName)
	assert.Equal(t.T(), expected.AccountHolderName, actual.AccountHolderName)
	assert.Equal(t.T(), expected.AccountNumber, actual.AccountNumber)
	assert.Equal(t.T(), expected.LastBalance, actual.LastBalance)
	assert.Equal(t.T(), expected.LastBalanceDate.Time().Unix(), actual.LastBalanceDate.Time().Unix())
	assert.Equal(t.T(), expected.Status, actual.Status)
}

func (t *bankAccountsHandlerTestSuite) TestCreate_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleCreateBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountsHandlerTestSuite) TestCreate_ServiceFailedCreating() {
	errMsg := "service failed creating bank account"
	input := t.getNewBankAccountInput(nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().Create(gomock.Any(), t.testUserID).Return(nil, failure.InternalError("create", "Bank Account", errors.New(errMsg)))

	t.handler.HandleCreateBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	assert.NotNil(t.T(), err.Entity)
	assert.Equal(t.T(), "Bank Account", *err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	assert.NotNil(t.T(), err.Operation)
	assert.Equal(t.T(), "create", *err.Operation)
}
