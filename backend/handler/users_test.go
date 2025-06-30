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

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kerti/balances/backend/handler"
	"github.com/kerti/balances/backend/handler/response"
	mock_service "github.com/kerti/balances/backend/mock/service"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/ctxprops"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/nuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type userHandlerTestSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	handler    handler.User
	mockSvc    *mock_service.MockUser
	testUserID uuid.UUID
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(userHandlerTestSuite))
}

func (t *userHandlerTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockSvc = mock_service.NewMockUser(t.ctrl)
	t.handler = &handler.UserImpl{
		Service: t.mockSvc,
	}
	t.testUserID, _ = uuid.NewV7()
}

func (t *userHandlerTestSuite) TearDownTest() {
	t.handler.Shutdown()
	t.ctrl.Finish()
}

func (t *userHandlerTestSuite) getNewRequestWithContext(method, path string, input any, formParams *map[string]string, routeVarId nuuid.NUUID) (recorder *httptest.ResponseRecorder, request *http.Request) {
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

func (t *userHandlerTestSuite) getNewUserInput(id nuuid.NUUID) model.UserInput {
	usr := model.UserInput{}

	if id.Valid {
		usr.ID = id.UUID
	} else {
		usr.ID = t.testUserID
	}

	usr.Username = "johndoe"
	usr.Email = "johndoe@balances.com"
	usr.Password = "thisisjohnspassword"
	usr.Name = "John Fitzgerald Doe"

	return usr
}

func (t *userHandlerTestSuite) parseOutputToUser(rr *httptest.ResponseRecorder) (actual *model.UserOutput, fail *failure.Failure) {
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

func (t *userHandlerTestSuite) TestGetByID_Normal() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/users/"+t.testUserID.String(),
		nil,
		nil,
		nuuid.From(t.testUserID),
	)

	input := t.getNewUserInput(nuuid.From(t.testUserID))
	expectedResult := model.NewUserFromInput(input, t.testUserID)
	expected := expectedResult.ToOutput()

	t.mockSvc.EXPECT().GetByID(t.testUserID).Return(&expectedResult, nil)

	t.handler.HandleGetUserByID(rr, req)

	actual, err := t.parseOutputToUser(rr)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected.Username, actual.Username)
	assert.Equal(t.T(), expected.Email, actual.Email)
	assert.Equal(t.T(), expected.Password, actual.Password)
	assert.Equal(t.T(), expected.Name, actual.Name)

	assert.Equal(t.T(), expected.Created.Time().Unix(), actual.Created.Time().Unix())
	assert.Equal(t.T(), expected.CreatedBy, actual.CreatedBy)
	assert.Equal(t.T(), expected.Updated, actual.Updated)
	assert.Equal(t.T(), expected.UpdatedBy, actual.UpdatedBy)
}

func (t *userHandlerTestSuite) TestGetByID_FailedParsingID() {
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/users/"+t.testUserID.String()+"123",
		nil,
		nil,
		nuuid.NUUID{Valid: false},
	)

	t.handler.HandleGetUserByID(rr, req)

	actual, err := t.parseOutputToUser(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), failure.CodeBadRequest, err.Code)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, "invalid UUID length")
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}

func (t *userHandlerTestSuite) TestGetByID_ServiceFailedResolving() {
	errMsg := "service failed resolving"
	rr, req := t.getNewRequestWithContext(
		http.MethodGet,
		"/users/"+t.testUserID.String(),
		nil,
		nil,
		nuuid.From(t.testUserID),
	)

	t.mockSvc.EXPECT().GetByID(t.testUserID).Return(nil, errors.New(errMsg))

	t.handler.HandleGetUserByID(rr, req)

	actual, err := t.parseOutputToUser(rr)

	assert.Nil(t.T(), actual)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), http.StatusInternalServerError, rr.Result().StatusCode)
	// TODO: specify this
	assert.Nil(t.T(), err.Entity)
	assert.Contains(t.T(), err.Message, errMsg)
	// TODO: specify this
	assert.Nil(t.T(), err.Operation)
}
