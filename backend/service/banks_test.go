package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	mock_repository "github.com/kerti/balances/backend/mock"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/service"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type bankAccountsServiceTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	svc      service.BankAccount
	mockRepo *mock_repository.MockBankAccount
}

func TestBankAccountsService(t *testing.T) {
	suite.Run(t, new(bankAccountsServiceTestSuite))
}

func (t *bankAccountsServiceTestSuite) SetupTest() {
	t.ctrl = gomock.NewController(t.T())
	t.mockRepo = mock_repository.NewMockBankAccount(t.ctrl)
	t.svc = &service.BankAccountImpl{
		Repository: t.mockRepo,
	}
	t.svc.Startup()
}

func (t *bankAccountsServiceTestSuite) TearDownTest() {
	t.svc.Shutdown()
	t.ctrl.Finish()
}

func (t *bankAccountsServiceTestSuite) TestResolveByID() {

	testID, _ := uuid.NewV7()
	testBalanceID1, _ := uuid.NewV7()
	testBalanceID2, _ := uuid.NewV7()
	nilCacheTime := cachetime.NCacheTime{}

	resolvedBankAccount := model.BankAccount{
		ID: testID,
	}

	resolvedBankAccountSlice := []model.BankAccount{resolvedBankAccount}

	resolvedBankAccountBalance1 := model.BankAccountBalance{
		ID: testBalanceID1,
	}

	resolvedBankAccountBalance2 := model.BankAccountBalance{
		ID: testBalanceID2,
	}

	resolvedBankAccountBalanceSlice := []model.BankAccountBalance{
		resolvedBankAccountBalance1,
		resolvedBankAccountBalance2,
	}

	t.Run("Exists", func() {

		t.Run("NoBalance", func() {
			t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
				Return(resolvedBankAccountSlice, nil)

			_, err := t.svc.GetByID(testID, false, nilCacheTime, nilCacheTime, nil)

			assert.NoError(t.T(), err, "should not error")
		})

		t.Run("WithBalance", func() {

			t.Run("NoFilter", func() {
				balanceFilterInput := model.BankAccountBalanceFilterInput{
					BankAccountIDs: &[]uuid.UUID{testID},
				}
				pageInfo := model.PageInfoOutput{}

				t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
					Return(resolvedBankAccountSlice, nil)
				t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
					Return(resolvedBankAccountBalanceSlice, pageInfo, nil)

				_, err := t.svc.GetByID(testID, true, nilCacheTime, nilCacheTime, nil)

				assert.NoError(t.T(), err, "should not error")
			})

			t.Run("WithStartDate", func() {
				yesterday := cachetime.NCacheTime(null.TimeFrom(time.Now().AddDate(0, 0, -1)))
				balanceFilterInput := model.BankAccountBalanceFilterInput{
					BankAccountIDs: &[]uuid.UUID{testID},
					StartDate:      yesterday,
				}
				pageInfo := model.PageInfoOutput{}

				t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
					Return(resolvedBankAccountSlice, nil)
				t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
					Return(resolvedBankAccountBalanceSlice, pageInfo, nil)

				_, err := t.svc.GetByID(testID, true, yesterday, nilCacheTime, nil)

				assert.NoError(t.T(), err, "should not error")
			})

			t.Run("WithEndDate", func() {
				now := cachetime.NCacheTime(null.TimeFrom(time.Now()))
				balanceFilterInput := model.BankAccountBalanceFilterInput{
					BankAccountIDs: &[]uuid.UUID{testID},
					EndDate:        now,
				}
				pageInfo := model.PageInfoOutput{}

				t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
					Return(resolvedBankAccountSlice, nil)
				t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
					Return(resolvedBankAccountBalanceSlice, pageInfo, nil)

				_, err := t.svc.GetByID(testID, true, nilCacheTime, now, nil)

				assert.NoError(t.T(), err, "should not error")
			})

			t.Run("WithBothDates", func() {
				yesterday := cachetime.NCacheTime(null.TimeFrom(time.Now().AddDate(0, 0, -1)))
				now := cachetime.NCacheTime(null.TimeFrom(time.Now()))
				balanceFilterInput := model.BankAccountBalanceFilterInput{
					BankAccountIDs: &[]uuid.UUID{testID},
					StartDate:      yesterday,
					EndDate:        now,
				}
				pageInfo := model.PageInfoOutput{}

				t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
					Return(resolvedBankAccountSlice, nil)
				t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
					Return(resolvedBankAccountBalanceSlice, pageInfo, nil)

				_, err := t.svc.GetByID(testID, true, yesterday, now, nil)

				assert.NoError(t.T(), err, "should not error")
			})

			t.Run("WithPageSize", func() {
				pageSize := 120
				balanceFilterInput := model.BankAccountBalanceFilterInput{
					BankAccountIDs: &[]uuid.UUID{testID},
				}
				balanceFilterInput.PageSize = &pageSize
				pageInfo := model.PageInfoOutput{}

				t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
					Return(resolvedBankAccountSlice, nil)
				t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
					Return(resolvedBankAccountBalanceSlice, pageInfo, nil)

				_, err := t.svc.GetByID(testID, true, nilCacheTime, nilCacheTime, &pageSize)

				assert.NoError(t.T(), err, "should not error")
			})

			t.Run("ErrorResolvingBalances", func() {
				errMsg := "error resolving balances"
				balanceFilterInput := model.BankAccountBalanceFilterInput{
					BankAccountIDs: &[]uuid.UUID{testID},
				}

				t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
					Return(resolvedBankAccountSlice, nil)
				t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
					Return([]model.BankAccountBalance{}, model.PageInfoOutput{}, errors.New(errMsg))

				_, err := t.svc.GetByID(testID, true, nilCacheTime, nilCacheTime, nil)

				assert.Error(t.T(), err, "should return error")
				assert.Contains(t.T(), err.Error(), errMsg)
			})

		})

	})

	t.Run("ErrorResolvingByID", func() {
		errMsg := "repo failed resolving by IDs"
		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
			Return(nil, errors.New(errMsg))

		_, err := t.svc.GetByID(testID, false, nilCacheTime, nilCacheTime, nil)

		assert.Error(t.T(), err, "should return error")
		assert.Equal(t.T(), err.Error(), errMsg)
	})

	t.Run("NotExists", func() {
		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
			Return([]model.BankAccount{}, nil)

		_, err := t.svc.GetByID(testID, false, nilCacheTime, nilCacheTime, nil)

		assert.Error(t.T(), err, "should return error")
		assert.Contains(t.T(), err.Error(), "EntityNotFound")
	})

}
