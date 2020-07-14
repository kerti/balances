package service

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Startup() {}

func (m *mockUserRepository) Shutdown() {}

func (m *mockUserRepository) ExistsByID(id uuid.UUID) (exists bool, err error) {
	m.Called(id)
	return true, nil
}

func (m *mockUserRepository) ResolveByIDs(ids []uuid.UUID) (users []model.User, err error) {
	args := m.Called(ids)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *mockUserRepository) ResolveByIdentity(identity string) (user model.User, err error) {
	args := m.Called(identity)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *mockUserRepository) Create(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepository) Update(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

var (
	testNow       = time.Now()
	testID1, _    = uuid.NewV4()
	testUserID, _ = uuid.NewV4()
	testUserInput = model.UserInput{
		ID:       testID1,
		Username: "username",
		Email:    "email@example.com",
		Password: "password",
		Name:     "John Doe",
	}
	testUserModel = model.NewUserFromInput(testUserInput, testUserID)
)

func TestUserService(t *testing.T) {

	t.Run("getByID", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			mockRepo := mockUserRepository{}
			svc := UserImpl{Repository: &mockRepo}

			mockRepo.On("ResolveByIDs", []uuid.UUID{testID1}).Return([]model.User{testUserModel}, nil)

			result, err := svc.GetByID(testID1)

			assert.NotNil(t, result)
			assert.Nil(t, err)

			mockRepo.AssertCalled(t, "ResolveByIDs", []uuid.UUID{testID1})
			mockRepo.AssertNumberOfCalls(t, "ResolveByIDs", 1)
		})

		t.Run("errorOnRepo", func(t *testing.T) {
			mockRepo := mockUserRepository{}
			svc := UserImpl{Repository: &mockRepo}

			mockRepo.On("ResolveByIDs", []uuid.UUID{testID1}).Return([]model.User{}, errors.New("error"))

			result, err := svc.GetByID(testID1)

			assert.Nil(t, result)
			assert.NotNil(t, err)
			assert.IsType(t, errors.New(""), err)

			mockRepo.AssertCalled(t, "ResolveByIDs", []uuid.UUID{testID1})
			mockRepo.AssertNumberOfCalls(t, "ResolveByIDs", 1)
		})

		t.Run("notFound", func(t *testing.T) {
			mockRepo := mockUserRepository{}
			svc := UserImpl{Repository: &mockRepo}

			mockRepo.On("ResolveByIDs", []uuid.UUID{testID1}).Return([]model.User{}, nil)

			result, err := svc.GetByID(testID1)

			assert.Nil(t, result)
			assert.NotNil(t, err)
			assert.IsType(t, &failure.Failure{}, err)
			assert.Equal(t, failure.CodeEntityNotFound, err.(*failure.Failure).Code)

			mockRepo.AssertCalled(t, "ResolveByIDs", []uuid.UUID{testID1})
			mockRepo.AssertNumberOfCalls(t, "ResolveByIDs", 1)
		})

	})

	t.Run("create", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			mockRepo := mockUserRepository{}
			svc := UserImpl{Repository: &mockRepo}

			mockRepo.On("Create", mock.AnythingOfType("model.User")).Return(nil)

			result, err := svc.Create(testUserInput, testUserID)

			assert.NotNil(t, result)
			assert.Nil(t, err)

			mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("model.User"))
			mockRepo.AssertNumberOfCalls(t, "Create", 1)
		})

		t.Run("error", func(t *testing.T) {
			mockRepo := mockUserRepository{}
			svc := UserImpl{Repository: &mockRepo}

			mockRepo.On("Create", mock.AnythingOfType("model.User")).Return(errors.New(""))

			result, err := svc.Create(testUserInput, testUserID)

			assert.NotNil(t, result)
			assert.NotNil(t, err)

			mockRepo.AssertCalled(t, "Create", mock.AnythingOfType("model.User"))
			mockRepo.AssertNumberOfCalls(t, "Create", 1)
		})

	})

	t.Run("update", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			mockRepo := mockUserRepository{}
			svc := UserImpl{Repository: &mockRepo}

			mockRepo.On("ResolveByIDs", []uuid.UUID{testID1}).Return([]model.User{testUserModel}, nil)
			mockRepo.On("Update", mock.AnythingOfType("model.User")).Return(nil)

			result, err := svc.Update(testUserInput, testUserID)

			assert.NotNil(t, result)
			assert.Nil(t, err)

			mockRepo.AssertCalled(t, "ResolveByIDs", []uuid.UUID{testID1})
			mockRepo.AssertNumberOfCalls(t, "ResolveByIDs", 1)
			mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("model.User"))
			mockRepo.AssertNumberOfCalls(t, "Update", 1)
		})

		t.Run("errorOnResolve", func(t *testing.T) {
			mockRepo := mockUserRepository{}
			svc := UserImpl{Repository: &mockRepo}

			mockRepo.On("ResolveByIDs", []uuid.UUID{testID1}).Return([]model.User{testUserModel}, errors.New(""))

			result, err := svc.Update(testUserInput, testUserID)

			assert.NotNil(t, result)
			assert.NotNil(t, err)

			mockRepo.AssertCalled(t, "ResolveByIDs", []uuid.UUID{testID1})
			mockRepo.AssertNumberOfCalls(t, "ResolveByIDs", 1)
			mockRepo.AssertNotCalled(t, "Update", mock.AnythingOfType("model.User"))
		})

		t.Run("notFound", func(t *testing.T) {
			mockRepo := mockUserRepository{}
			svc := UserImpl{Repository: &mockRepo}

			mockRepo.On("ResolveByIDs", []uuid.UUID{testID1}).Return([]model.User{}, nil)

			result, err := svc.Update(testUserInput, testUserID)

			assert.NotNil(t, result)
			assert.NotNil(t, err)
			assert.IsType(t, &failure.Failure{}, err)
			assert.Equal(t, failure.CodeEntityNotFound, err.(*failure.Failure).Code)

			mockRepo.AssertCalled(t, "ResolveByIDs", []uuid.UUID{testID1})
			mockRepo.AssertNumberOfCalls(t, "ResolveByIDs", 1)
			mockRepo.AssertNotCalled(t, "Update", mock.AnythingOfType("model.User"))
		})

	})
}
