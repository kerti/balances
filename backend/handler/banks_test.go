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

type bankAccountHandlerTestSuite struct {
	suite.Suite
	ctrl                     *gomock.Controller
	handler                  handler.BankAccount
	mockSvc                  *mock_service.MockBankAccount
	testUserID               uuid.UUID
	testBankAccountID        uuid.UUID
	testBankAccountBalanceID uuid.UUID
}

func TestBankAccountHandler(t *testing.T) {
	suite.Run(t, new(bankAccountHandlerTestSuite))
}

func (t *bankAccountHandlerTestSuite) SetupTest() {
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

func (t *bankAccountHandlerTestSuite) TearDownTest() {
	t.handler.Shutdown()
	t.ctrl.Finish()
}

func (t *bankAccountHandlerTestSuite) getNewRequestWithContext(method, path string, input any, formParams *map[string]string, routeVarId nuuid.NUUID) (recorder *httptest.ResponseRecorder, request *http.Request) {
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

func (t *bankAccountHandlerTestSuite) getNewBankAccountInput(id nuuid.NUUID) model.BankAccountInput {
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

func (t *bankAccountHandlerTestSuite) getNewBankAccountBalanceInput(id, bankAccountID nuuid.NUUID) model.BankAccountBalanceInput {
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

func (t *bankAccountHandlerTestSuite) parseOutputToBankAccount(rr *httptest.ResponseRecorder) (actual *model.BankAccountOutput, fail *failure.Failure) {
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

func (t *bankAccountHandlerTestSuite) parseOutputToBankAccountBalance(rr *httptest.ResponseRecorder) (actual *model.BankAccountBalanceOutput, fail *failure.Failure) {
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

func (t *bankAccountHandlerTestSuite) parseOutputToBankAccountPage(rr *httptest.ResponseRecorder) (items []model.BankAccountOutput, pageInfo model.PageInfoOutput, fail *failure.Failure) {
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
		for _, bankAccountInterface := range actualSlice {
			bankAccountMap := (bankAccountInterface).(map[string]any)
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

func (t *bankAccountHandlerTestSuite) parseOutputToBankAccountBalancePage(rr *httptest.ResponseRecorder) (items []model.BankAccountBalanceOutput, pageInfo model.PageInfoOutput, fail *failure.Failure) {
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
		for _, bankAccountBalanceInterface := range actualSlice {
			bankAccountBalanceMap := (bankAccountBalanceInterface).(map[string]any)
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

func (t *bankAccountHandlerTestSuite) TestCreate_Normal() {
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

func (t *bankAccountHandlerTestSuite) TestCreate_FailedParsingRequestPayload() {
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

func (t *bankAccountHandlerTestSuite) TestCreate_ServiceFailedCreating() {
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

func (t *bankAccountHandlerTestSuite) TestGetByID_Normal_NoParams() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		nil,
		nuuid.From(t.testBankAccountID),
	)

	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	expectedResult := model.NewBankAccountFromInput(input, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().GetByID(t.testBankAccountID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).Return(&expectedResult, nil)

	t.handler.HandleGetBankAccountByID(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.AccountName, actual.AccountName)
	assert.Equal(t.T(), expected.BankName, actual.BankName)
	assert.Equal(t.T(), expected.AccountNumber, actual.AccountNumber)
	assert.Equal(t.T(), expected.AccountHolderName, actual.AccountHolderName)
	assert.Equal(t.T(), expected.LastBalance, actual.LastBalance)
	assert.Equal(t.T(), expected.LastBalanceDate.Time().Unix(), actual.LastBalanceDate.Time().Unix())
	assert.Equal(t.T(), expected.Status, actual.Status)
}

func (t *bankAccountHandlerTestSuite) TestGetByID_FailedParsingID() {
	formParams := make(map[string]string)
	formParams["id"] = t.testBankAccountID.String() + "123"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/"+t.testBankAccountID.String()+"123",
		&formParams,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetBankAccountByID(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestGetByID_Normal_WithBalances() {
	formParams := make(map[string]string)
	formParams["withBalances"] = "true"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		&formParams,
		nuuid.From(t.testBankAccountID),
	)

	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	expectedResult := model.NewBankAccountFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().
		GetByID(t.testBankAccountID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).
		Return(&expectedResult, nil)

	t.handler.HandleGetBankAccountByID(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *bankAccountHandlerTestSuite) TestGetByID_Normal_WithBalancesStartDate() {
	startDate := time.Unix(0, time.Now().AddDate(0, 0, -1).UnixMilli()*int64(time.Millisecond))
	formParams := make(map[string]string)
	formParams["balanceStartDate"] = strconv.FormatInt(startDate.UnixMilli(), 10)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		&formParams,
		nuuid.From(t.testBankAccountID),
	)

	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	expectedResult := model.NewBankAccountFromInput(input, t.testUserID)

	var nStartDate cachetime.NCacheTime
	nStartDate.Scan(startDate)
	t.mockSvc.EXPECT().
		GetByID(t.testBankAccountID, false, nStartDate, cachetime.NCacheTime{}, nil).
		Return(&expectedResult, nil)

	t.handler.HandleGetBankAccountByID(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *bankAccountHandlerTestSuite) TestGetByID_Normal_WithBalancesEndDate() {
	endDate := time.Unix(0, time.Now().UnixMilli()*int64(time.Millisecond))
	formParams := make(map[string]string)
	formParams["balanceEndDate"] = strconv.FormatInt(endDate.UnixMilli(), 10)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		&formParams,
		nuuid.From(t.testBankAccountID),
	)

	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	expectedResult := model.NewBankAccountFromInput(input, t.testUserID)

	var nEndDate cachetime.NCacheTime
	nEndDate.Scan(endDate)
	t.mockSvc.EXPECT().
		GetByID(t.testBankAccountID, false, cachetime.NCacheTime{}, nEndDate, nil).
		Return(&expectedResult, nil)

	t.handler.HandleGetBankAccountByID(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *bankAccountHandlerTestSuite) TestGetByID_Normal_WithPageSize() {
	pageSize := 10
	formParams := make(map[string]string)
	formParams["pageSize"] = strconv.Itoa(pageSize)
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		&formParams,
		nuuid.From(t.testBankAccountID),
	)

	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	expectedResult := model.NewBankAccountFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().
		GetByID(t.testBankAccountID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, &pageSize).
		Return(&expectedResult, nil)

	t.handler.HandleGetBankAccountByID(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *bankAccountHandlerTestSuite) TestGetByID_Normal_ServiceFailedResolving() {
	errMsg := "failed resolving bank account"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		nil,
		nuuid.From(t.testBankAccountID),
	)

	t.mockSvc.EXPECT().GetByID(t.testBankAccountID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil).
		Return(nil, failure.InternalError("get by ID", "Bank Account", errors.New(errMsg)))

	t.handler.HandleGetBankAccountByID(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	assert.NotNil(t.T(), err.Entity)
	assert.Equal(t.T(), "Bank Account", *err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	assert.NotNil(t.T(), err.Operation)
	assert.Equal(t.T(), "get by ID", *err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestGetByFilter_Normal() {
	keyword := "test keyword"
	input := model.BankAccountFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedBankAccounts := []model.BankAccount{}
	acc1 := model.NewBankAccountFromInput(t.getNewBankAccountInput(nuuid.NUUID{}), t.testUserID)
	acc2 := model.NewBankAccountFromInput(t.getNewBankAccountInput(nuuid.NUUID{}), t.testUserID)
	expectedBankAccounts = append(expectedBankAccounts, acc1)
	expectedBankAccounts = append(expectedBankAccounts, acc2)
	expectedPageInfo := model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 1,
		PageCount:  1,
	}

	t.mockSvc.EXPECT().GetByFilter(input).Return(expectedBankAccounts, expectedPageInfo, nil)

	t.handler.HandleGetBankAccountByFilter(rr, req)

	bankAccounts, pageInfo, err := t.parseOutputToBankAccountPage(rr)

	assert.Nil(t.T(), err)

	assert.Equal(t.T(), len(expectedBankAccounts), len(bankAccounts))
	assert.Equal(t.T(), expectedBankAccounts[0].ID, bankAccounts[0].ID)
	assert.Equal(t.T(), expectedBankAccounts[1].ID, bankAccounts[1].ID)

	assert.Equal(t.T(), 1, pageInfo.Page)
}

func (t *bankAccountHandlerTestSuite) TestGetByFilter_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetBankAccountByFilter(rr, req)

	bankAccounts, pageInfo, err := t.parseOutputToBankAccountPage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(bankAccounts))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *bankAccountHandlerTestSuite) TestGetByFilter_ServiceFailedResolving() {
	errMsg := "failed resolving bank accounts by filter"
	keyword := "test keyword"
	input := model.BankAccountFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().GetByFilter(input).
		Return(
			[]model.BankAccount{},
			model.PageInfoOutput{},
			failure.InternalError("get by filter", "Bank Account",
				errors.New(errMsg)))

	t.handler.HandleGetBankAccountByFilter(rr, req)

	bankAccounts, pageInfo, err := t.parseOutputToBankAccountPage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	assert.NotNil(t.T(), err.Entity)
	assert.Equal(t.T(), "Bank Account", *err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	assert.NotNil(t.T(), err.Operation)
	assert.Equal(t.T(), "get by filter", *err.Operation)

	assert.Equal(t.T(), 0, len(bankAccounts))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *bankAccountHandlerTestSuite) TestUpdate_Normal() {
	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/"+t.testBankAccountID.String(),
		input,
		nil,
		nuuid.From(t.testBankAccountID),
	)

	updatedBankAccount := model.NewBankAccountFromInput(input, t.testUserID)

	t.mockSvc.EXPECT().Update(gomock.Any(), t.testUserID).Return(&updatedBankAccount, nil)

	t.handler.HandleUpdateBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *bankAccountHandlerTestSuite) TestUpdate_FailedGettingIDFromRequest() {
	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/"+t.testBankAccountID.String(),
		input,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleUpdateBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestUpdate_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/"+t.testBankAccountID.String(),
		input,
		nil,
		nuuid.From(t.testBankAccountID),
	)

	t.handler.HandleUpdateBankAccount(rr, req)

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

func (t *bankAccountHandlerTestSuite) TestUpdate_MismatchedID() {
	input := t.getNewBankAccountInput(nuuid.NUUID{})
	newID, _ := uuid.NewV7()
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/"+t.testBankAccountID.String(),
		input,
		nil,
		nuuid.From(newID),
	)

	t.handler.HandleUpdateBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "id mismatch")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestUpdate_ServiceFailedUpdating() {
	errMsg := "failed updating bankAccount"
	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/"+t.testBankAccountID.String(),
		input,
		nil,
		nuuid.From(t.testBankAccountID),
	)

	t.mockSvc.EXPECT().Update(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleUpdateBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestDelete_Normal() {
	input := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID))
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		nil,
		nuuid.From(t.testBankAccountID),
	)

	deletedBankAccount := model.NewBankAccountFromInput(input, t.testUserID)
	deletedBankAccount.ID = t.testBankAccountID

	t.mockSvc.EXPECT().Delete(t.testBankAccountID, t.testUserID).Return(&deletedBankAccount, nil)

	t.handler.HandleDeleteBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), t.testBankAccountID, actual.ID)
	assert.Nil(t.T(), err)
}

func (t *bankAccountHandlerTestSuite) TestDelete_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleDeleteBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestDelete_ServiceFailedDeleting() {
	errMsg := "service failed deleting bankAccount"
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/bankAccounts/"+t.testBankAccountID.String(),
		nil,
		nil,
		nuuid.From(t.testBankAccountID),
	)

	t.mockSvc.EXPECT().Delete(t.testBankAccountID, t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleDeleteBankAccount(rr, req)

	actual, err := t.parseOutputToBankAccount(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestCreateBalance_Normal() {
	input := t.getNewBankAccountBalanceInput(nuuid.NUUID{Valid: false}, nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/balances",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedResult := model.NewBankAccountBalanceFromInput(input, input.BankAccountID, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().CreateBalance(gomock.Any(), t.testUserID).Return(&expectedResult, nil)

	t.handler.HandleCreateBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), http.StatusCreated, rr.Result().StatusCode)
	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), expected.ID, actual.ID)
	assert.Equal(t.T(), expected.BankAccountID, actual.BankAccountID)
	assert.Equal(t.T(), expected.Date.Time().Unix(), actual.Date.Time().Unix())
	assert.Equal(t.T(), expected.Balance, actual.Balance)
	assert.NotNil(t.T(), actual.Created)
	assert.NotNil(t.T(), actual.CreatedBy)
	assert.False(t.T(), actual.Updated.Valid)
	assert.False(t.T(), actual.UpdatedBy.Valid)
	assert.False(t.T(), actual.Deleted.Valid)
	assert.False(t.T(), actual.DeletedBy.Valid)
}

func (t *bankAccountHandlerTestSuite) TestCreateBalance_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/balances",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleCreateBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestCreateBalance_ServiceFailedCreatingBalance() {
	errMsg := "service failed creating bank account balances"
	input := t.getNewBankAccountBalanceInput(nuuid.NUUID{Valid: false}, nuuid.NUUID{Valid: false})
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/balances",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().CreateBalance(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleCreateBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestGetBalanceByID_Normal() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		nil,
		nil,
		nuuid.From(t.testBankAccountBalanceID),
	)

	input := t.getNewBankAccountBalanceInput(nuuid.From(t.testBankAccountBalanceID), nuuid.From(t.testBankAccountID))
	expectedResult := model.NewBankAccountBalanceFromInput(input, t.testBankAccountID, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().GetBalanceByID(t.testBankAccountBalanceID).Return(&expectedResult, nil)

	t.handler.HandleGetBankAccountBalanceByID(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.ID, actual.ID)
	assert.Equal(t.T(), expected.BankAccountID, actual.BankAccountID)
	assert.Equal(t.T(), expected.Date.Time().Unix(), actual.Date.Time().Unix())
	assert.Equal(t.T(), expected.Balance, actual.Balance)
	assert.Equal(t.T(), expected.Created.Time().Unix(), actual.Created.Time().Unix())
	assert.Equal(t.T(), expected.CreatedBy, actual.CreatedBy)
	assert.Equal(t.T(), expected.Updated, actual.Updated)
	assert.Equal(t.T(), expected.UpdatedBy, actual.UpdatedBy)
	assert.Equal(t.T(), expected.Deleted, actual.Deleted)
	assert.Equal(t.T(), expected.DeletedBy, actual.DeletedBy)
}

func (t *bankAccountHandlerTestSuite) TestGetBalanceByID_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		nil,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleGetBankAccountBalanceByID(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestGetBalanceByID_ServiceFailedResolving() {
	errMsg := "service failed resolving bank account balance"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		nil,
		nil,
		nuuid.From(t.testBankAccountBalanceID),
	)

	t.mockSvc.EXPECT().GetBalanceByID(t.testBankAccountBalanceID).Return(nil, errors.New(errMsg))

	t.handler.HandleGetBankAccountBalanceByID(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestGetBalanceByFilter_Normal() {
	keyword := "test keyword"
	input := model.BankAccountBalanceFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/balances/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	expectedBankAccountBalances := []model.BankAccountBalance{}
	vv1 := model.NewBankAccountBalanceFromInput(t.getNewBankAccountBalanceInput(nuuid.NUUID{}, nuuid.From(t.testBankAccountID)), t.testBankAccountID, t.testUserID)
	vv2 := model.NewBankAccountBalanceFromInput(t.getNewBankAccountBalanceInput(nuuid.NUUID{}, nuuid.From(t.testBankAccountID)), t.testBankAccountID, t.testUserID)
	expectedBankAccountBalances = append(expectedBankAccountBalances, vv1)
	expectedBankAccountBalances = append(expectedBankAccountBalances, vv2)
	expectedPageInfo := model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 1,
		PageCount:  1,
	}

	t.mockSvc.EXPECT().GetBalancesByFilter(input).Return(expectedBankAccountBalances, expectedPageInfo, nil)

	t.handler.HandleGetBankAccountBalanceByFilter(rr, req)

	bankAccountBalances, pageInfo, err := t.parseOutputToBankAccountBalancePage(rr)

	assert.Nil(t.T(), err)

	assert.Equal(t.T(), len(expectedBankAccountBalances), len(bankAccountBalances))
	assert.Equal(t.T(), expectedBankAccountBalances[0].ID, bankAccountBalances[0].ID)
	assert.Equal(t.T(), expectedBankAccountBalances[1].ID, bankAccountBalances[1].ID)

	assert.Equal(t.T(), 1, pageInfo.Page)
}

func (t *bankAccountHandlerTestSuite) TestGetBalanceByFilter_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/balances/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetBankAccountBalanceByFilter(rr, req)

	bankAccounts, pageInfo, err := t.parseOutputToBankAccountBalancePage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(bankAccounts))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *bankAccountHandlerTestSuite) TestGetBalanceByFilter_ServiceFailedResolving() {
	errMsg := "service failed resolving bank account balances"
	keyword := "test keyword"
	input := model.BankAccountBalanceFilterInput{}
	input.Keyword = &keyword
	rr, req := t.getNewRequestWithContext(
		http.MethodPost,
		"/bankAccounts/balances/search",
		input,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.mockSvc.EXPECT().GetBalancesByFilter(input).Return([]model.BankAccountBalance{}, model.PageInfoOutput{}, errors.New(errMsg))

	t.handler.HandleGetBankAccountBalanceByFilter(rr, req)

	vahicleBalances, pageInfo, err := t.parseOutputToBankAccountBalancePage(rr)

	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)

	assert.Equal(t.T(), 0, len(vahicleBalances))

	assert.Equal(t.T(), 0, pageInfo.Page)
}

func (t *bankAccountHandlerTestSuite) TestUpdateBalance_Normal() {
	input := t.getNewBankAccountBalanceInput(nuuid.From(t.testBankAccountBalanceID), nuuid.From(t.testBankAccountID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		input,
		nil,
		nuuid.From(t.testBankAccountBalanceID),
	)

	updatedBankAccountBalance := model.NewBankAccountBalanceFromInput(input, t.testBankAccountID, t.testUserID)

	t.mockSvc.EXPECT().UpdateBalance(gomock.Any(), t.testUserID).Return(&updatedBankAccountBalance, nil)

	t.handler.HandleUpdateBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.NotNil(t.T(), actual)
	assert.Nil(t.T(), err)
}

func (t *bankAccountHandlerTestSuite) TestUpdateBalance_FailedGettingIDFromRequest() {
	input := t.getNewBankAccountBalanceInput(nuuid.From(t.testBankAccountBalanceID), nuuid.From(t.testBankAccountID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		input,
		nil,
		nuuid.NUUID{},
	)

	t.handler.HandleUpdateBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestUpdateBalance_FailedParsingRequestPayload() {
	input := "test"
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		input,
		nil,
		nuuid.From(t.testBankAccountBalanceID),
	)

	t.handler.HandleUpdateBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "cannot unmarshal")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestUpdateBalance_MismatchedID() {
	input := t.getNewBankAccountBalanceInput(nuuid.NUUID{}, nuuid.NUUID{})
	newID, _ := uuid.NewV7()
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/"+t.testBankAccountBalanceID.String(),
		input,
		nil,
		nuuid.From(newID),
	)

	t.handler.HandleUpdateBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "id mismatch")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestUpdateBalance_ServiceFailedUpdating() {
	errMsg := "failed updating bank account balance"
	input := t.getNewBankAccountBalanceInput(nuuid.From(t.testBankAccountBalanceID), nuuid.From(t.testBankAccountID))
	rr, req := t.getNewRequestWithContext(
		http.MethodPatch,
		"/bankAccounts/"+t.testBankAccountBalanceID.String(),
		input,
		nil,
		nuuid.From(t.testBankAccountBalanceID),
	)

	t.mockSvc.EXPECT().UpdateBalance(gomock.Any(), t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleUpdateBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestDeleteBalance_Normal() {
	input := t.getNewBankAccountBalanceInput(nuuid.From(t.testBankAccountBalanceID), nuuid.From(t.testBankAccountID))
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		nil,
		nil,
		nuuid.From(t.testBankAccountBalanceID),
	)

	deletedBankAccountBalance := model.NewBankAccountBalanceFromInput(input, t.testBankAccountID, t.testUserID)
	deletedBankAccountBalance.ID = t.testBankAccountID

	t.mockSvc.EXPECT().DeleteBalance(t.testBankAccountBalanceID, t.testUserID).Return(&deletedBankAccountBalance, nil)

	t.handler.HandleDeleteBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.NotNil(t.T(), actual)
	assert.Equal(t.T(), t.testBankAccountID, actual.ID)
	assert.Nil(t.T(), err)
}

func (t *bankAccountHandlerTestSuite) TestDeleteBalance_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		nil,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleDeleteBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *bankAccountHandlerTestSuite) TestDeleteBalance_ServiceFailedDeleting() {
	errMsg := "service failed deleting bank account balance"
	rr, req := t.getNewRequestWithContext(
		http.MethodDelete,
		"/bankAccounts/balances/"+t.testBankAccountBalanceID.String(),
		nil,
		nil,
		nuuid.From(t.testBankAccountBalanceID),
	)

	t.mockSvc.EXPECT().DeleteBalance(t.testBankAccountBalanceID, t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleDeleteBankAccountBalance(rr, req)

	actual, err := t.parseOutputToBankAccountBalance(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeInternalError, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}
