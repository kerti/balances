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
	"github.com/kerti/balances/backend/util/nuuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type bankAccountsServiceTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	svc               service.BankAccount
	mockRepo          *mock_repository.MockBankAccount
	testUserID        uuid.UUID
	testBankAccountID uuid.UUID
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
	t.testUserID, _ = uuid.NewV7()
	t.testBankAccountID, _ = uuid.NewV7()
	t.svc.Startup()
}

func (t *bankAccountsServiceTestSuite) TearDownTest() {
	t.svc.Shutdown()
	t.ctrl.Finish()
}

func (t *bankAccountsServiceTestSuite) getNewBankAccountInput(id nuuid.NUUID, balances *[]model.BankAccountBalanceInput) model.BankAccountInput {
	acc := model.BankAccountInput{}

	if id.Valid {
		acc.ID = id.UUID
	} else {
		acc.ID = t.testBankAccountID
	}

	lastBalanceDate := time.Now().AddDate(0, -1, 0) //defaults to last month

	acc.AccountName = "Savings Account"
	acc.BankName = "First National Bank"
	acc.AccountHolderName = "John Doe"
	acc.AccountNumber = "123-456-7890"
	acc.LastBalance = float64(1000000)
	acc.LastBalanceDate = cachetime.CacheTime(lastBalanceDate)
	acc.Status = model.BankAccountStatusActive

	acc.Balances = []model.BankAccountBalanceInput{}
	if balances != nil {
		for _, bal := range *balances {
			balCopy := bal
			acc.Balances = append(acc.Balances, balCopy)
		}
	}

	return acc
}

func (t *bankAccountsServiceTestSuite) getNewBankAccount(id nuuid.NUUID, balances *[]model.BankAccountBalance) model.BankAccount {
	acc := model.BankAccount{}

	if id.Valid {
		acc.ID = id.UUID
	} else {
		acc.ID = t.testBankAccountID
	}

	acc.AccountName = "Savings Account"
	acc.BankName = "First National Bank"
	acc.AccountHolderName = "John Doe"
	acc.AccountNumber = "123-456-7890"
	acc.LastBalance = float64(1000000)
	acc.LastBalanceDate = time.Now().AddDate(0, 1, 0) // defaults to last month
	acc.Status = model.BankAccountStatusActive
	acc.Created = time.Now()
	acc.CreatedBy = t.testUserID

	acc.Balances = []model.BankAccountBalance{}
	if balances != nil {
		for _, bal := range *balances {
			balCopy := bal
			acc.Balances = append(acc.Balances, balCopy)
		}
	}

	return acc
}

func (t *bankAccountsServiceTestSuite) TestCreate_Normal() {
	testInput := t.getNewBankAccountInput(nuuid.NUUID{Valid: false}, nil)
	t.mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	res, err := t.svc.Create(testInput, t.testUserID)

	assert.NoError(t.T(), err)
	assert.Equal(t.T(), testInput.AccountName, res.AccountName)
	assert.Equal(t.T(), testInput.BankName, res.BankName)
	assert.Equal(t.T(), testInput.AccountHolderName, res.AccountHolderName)
	assert.Equal(t.T(), testInput.AccountNumber, res.AccountNumber)
	assert.Equal(t.T(), testInput.LastBalance, res.LastBalance)
	assert.Equal(t.T(), testInput.LastBalanceDate.Time(), res.LastBalanceDate)
	assert.Equal(t.T(), testInput.Status, res.Status)
	assert.Equal(t.T(), len(res.Balances), 1)
	assert.Equal(t.T(), testInput.LastBalanceDate.Time(), res.Balances[0].Date)
	assert.Equal(t.T(), testInput.LastBalance, res.Balances[0].Balance)
}

func (t *bankAccountsServiceTestSuite) TestCreate_RepoFailToCreate() {
	errMsg := "failed to create bank account"
	testInput := t.getNewBankAccountInput(nuuid.NUUID{Valid: false}, nil)
	t.mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New(errMsg))

	res, err := t.svc.Create(testInput, t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *bankAccountsServiceTestSuite) TestGetByID() {

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

			assert.NoError(t.T(), err)
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

				assert.NoError(t.T(), err)
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

				assert.NoError(t.T(), err)
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

				assert.NoError(t.T(), err)
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

				assert.NoError(t.T(), err)
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

				assert.NoError(t.T(), err)
			})

			t.Run("RepoErrorResolvingAccount", func() {
				t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
					Return([]model.BankAccount{}, nil)

				res, err := t.svc.GetByID(testID, true, nilCacheTime, nilCacheTime, nil)

				assert.Nil(t.T(), res)
				assert.Error(t.T(), err)
				assert.Contains(t.T(), err.Error(), "EntityNotFound")
			})

			t.Run("RepoErrorResolvingBalances", func() {
				errMsg := "error resolving balances"
				balanceFilterInput := model.BankAccountBalanceFilterInput{
					BankAccountIDs: &[]uuid.UUID{testID},
				}

				t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
					Return(resolvedBankAccountSlice, nil)
				t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
					Return([]model.BankAccountBalance{}, model.PageInfoOutput{}, errors.New(errMsg))

				res, err := t.svc.GetByID(testID, true, nilCacheTime, nilCacheTime, nil)

				assert.Nil(t.T(), res)
				assert.Error(t.T(), err)
				assert.Contains(t.T(), err.Error(), errMsg)
			})

		})

	})

	t.Run("ErrorResolvingByID", func() {
		errMsg := "repo failed resolving by IDs"
		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
			Return(nil, errors.New(errMsg))

		_, err := t.svc.GetByID(testID, false, nilCacheTime, nilCacheTime, nil)

		assert.Error(t.T(), err)
		assert.Equal(t.T(), err.Error(), errMsg)
	})

	t.Run("NotExists", func() {
		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).
			Return([]model.BankAccount{}, nil)

		_, err := t.svc.GetByID(testID, false, nilCacheTime, nilCacheTime, nil)

		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), "EntityNotFound")
	})

}

func (t *bankAccountsServiceTestSuite) TestGetByFilter() {

	bankAccount1 := model.BankAccount{}
	bankAccount2 := model.BankAccount{}
	bankAccountSlice := []model.BankAccount{bankAccount1, bankAccount2}
	defaultPageInfo := model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 2,
		PageCount:  1,
	}

	t.Run("EmptyFilter", func() {
		filterInput := model.BankAccountFilterInput{}
		filter := filterInput.ToFilter()

		t.mockRepo.EXPECT().ResolveByFilter(filter).
			Return(bankAccountSlice, defaultPageInfo, nil)

		_, _, err := t.svc.GetByFilter(filterInput)

		assert.NoError(t.T(), err)
	})

	t.Run("WithKeyword", func() {
		keyword := "example"
		filterInput := model.BankAccountFilterInput{}
		filterInput.Keyword = &keyword
		filter := filterInput.ToFilter()

		t.mockRepo.EXPECT().ResolveByFilter(filter).
			Return(bankAccountSlice, defaultPageInfo, nil)

		_, _, err := t.svc.GetByFilter(filterInput)

		assert.NoError(t.T(), err)
	})

}

func (t *bankAccountsServiceTestSuite) TestUpdate() {

	testUserID, _ := uuid.NewV7()
	testOldUserID, _ := uuid.NewV7()
	testBankAccountID, _ := uuid.NewV7()

	testLastBalanceDate := time.Now().AddDate(0, 0, -2)
	testLastBalance := float64(1000000)

	testNewLastBalanceDate := time.Now()
	testNewLastBalance := float64(1100000)

	bankAccountInput := model.BankAccountInput{
		ID:                testBankAccountID,
		AccountName:       "Savings Account Updated",
		BankName:          "First National Bank Updated",
		AccountHolderName: "John Doe Updated",
		AccountNumber:     "123-456-7890-updated",
		LastBalance:       testNewLastBalance,
		LastBalanceDate:   cachetime.CacheTime(testNewLastBalanceDate),
		Status:            model.BankAccountStatusActive,
		Balances:          []model.BankAccountBalanceInput{},
	}

	existingBankAccount := model.BankAccount{
		ID:                testBankAccountID,
		AccountName:       "Savings Account",
		BankName:          "First National Bank",
		AccountHolderName: "John Doe",
		AccountNumber:     "123-456-7890",
		LastBalance:       testLastBalance,
		LastBalanceDate:   testLastBalanceDate,
		Status:            model.BankAccountStatusActive,
		Created:           time.Now().AddDate(0, 0, -1),
		CreatedBy:         testOldUserID,
	}

	deletedBankAccount := model.BankAccount{
		ID:                testBankAccountID,
		AccountName:       "Savings Account",
		BankName:          "First National Bank",
		AccountHolderName: "John Doe",
		AccountNumber:     "123-456-7890",
		LastBalance:       testLastBalance,
		LastBalanceDate:   testLastBalanceDate,
		Status:            model.BankAccountStatusActive,
		Created:           time.Now().AddDate(0, 0, -2),
		CreatedBy:         testOldUserID,
		Deleted:           null.TimeFrom(time.Now().AddDate(0, 0, -1)),
		DeletedBy:         nuuid.From(testOldUserID),
	}

	t.Run("Normal", func() {

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{existingBankAccount}, nil)

		t.mockRepo.EXPECT().Update(gomock.Any()).
			Return(nil)

		res, err := t.svc.Update(bankAccountInput, testUserID)

		assert.NotNil(t.T(), res)
		assert.NoError(t.T(), err)

	})

	t.Run("ErrorResolvingByIDs", func() {

		errMsg := "query failed"

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{}, errors.New(errMsg))

		res, err := t.svc.Update(bankAccountInput, testUserID)

		assert.Nil(t.T(), res)
		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), errMsg)

	})

	t.Run("AccountNotFound", func() {

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{}, nil)

		res, err := t.svc.Update(bankAccountInput, testUserID)

		assert.Nil(t.T(), res)
		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), "EntityNotFound")

	})

	t.Run("AccountDeleted", func() {

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{deletedBankAccount}, nil)

		res, err := t.svc.Update(bankAccountInput, testUserID)

		assert.Nil(t.T(), res)
		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), "update")
		assert.Contains(t.T(), err.Error(), "Bank Account")
		assert.Contains(t.T(), err.Error(), "deleted")

	})

	t.Run("RepoFailToUpdate", func() {
		errMsg := "failed to update"

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{existingBankAccount}, nil)

		t.mockRepo.EXPECT().Update(gomock.Any()).
			Return(errors.New(errMsg))

		res, err := t.svc.Update(bankAccountInput, testUserID)

		assert.Nil(t.T(), res)
		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), errMsg)

	})

}

func (t *bankAccountsServiceTestSuite) TestDelete() {

	testUserID, _ := uuid.NewV7()
	testBankAccountID, _ := uuid.NewV7()
	testBankAccountBalanceID1, _ := uuid.NewV7()
	testBankAccountBalanceID2, _ := uuid.NewV7()

	lastBalance := float64(1000000)
	lastBalanceDate := time.Now()

	nonLastBalance := float64(900000)
	nonLastBalanceDate := time.Now().AddDate(0, 0, -1)

	testLastBankAccountBalance := model.BankAccountBalance{
		ID:            testBankAccountBalanceID1,
		BankAccountID: testBankAccountID,
		Date:          lastBalanceDate,
		Balance:       lastBalance,
		Created:       time.Now().AddDate(0, 0, -1),
		CreatedBy:     testUserID,
	}

	testNonLastBankAccountBalance := model.BankAccountBalance{
		ID:            testBankAccountBalanceID2,
		BankAccountID: testBankAccountID,
		Date:          nonLastBalanceDate,
		Balance:       nonLastBalance,
		Created:       time.Now().AddDate(0, 0, -1),
		CreatedBy:     testUserID,
	}

	testDeletedNonLastBankAccountBalance := model.BankAccountBalance{
		ID:            testBankAccountBalanceID2,
		BankAccountID: testBankAccountID,
		Date:          nonLastBalanceDate,
		Balance:       nonLastBalance,
		Created:       time.Now().AddDate(0, 0, -1),
		CreatedBy:     testUserID,
		Deleted:       null.TimeFrom(time.Now()),
		DeletedBy:     nuuid.From(testUserID),
	}

	testBankAccount := model.BankAccount{
		ID:                testBankAccountID,
		AccountName:       "Savings Account",
		BankName:          "First National Bank",
		AccountHolderName: "John Doe",
		AccountNumber:     "123-456-7890",
		LastBalance:       lastBalance,
		LastBalanceDate:   lastBalanceDate,
		Status:            model.BankAccountStatusActive,
		Created:           time.Now().AddDate(0, 0, -1),
		CreatedBy:         testUserID,
	}

	testDeletedBankAccount := model.BankAccount{
		ID:                testBankAccountID,
		AccountName:       "Savings Account",
		BankName:          "First National Bank",
		AccountHolderName: "John Doe",
		AccountNumber:     "123-456-7890",
		LastBalance:       lastBalance,
		LastBalanceDate:   lastBalanceDate,
		Status:            model.BankAccountStatusActive,
		Created:           time.Now().AddDate(0, 0, -1),
		CreatedBy:         testUserID,
		Deleted:           null.TimeFrom(time.Now()),
		DeletedBy:         nuuid.From(testUserID),
	}

	t.Run("Normal", func() {

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{testBankAccount}, nil)

		t.mockRepo.EXPECT().ResolveBalancesByFilter(gomock.Any()).
			Return(
				[]model.BankAccountBalance{
					testNonLastBankAccountBalance,
					testLastBankAccountBalance,
				},
				model.PageInfoOutput{},
				nil).Times(2)

		t.mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

		res, err := t.svc.Delete(testBankAccountID, testUserID)

		assert.NoError(t.T(), err)

		assert.NotNil(t.T(), res)
		assert.True(t.T(), res.Deleted.Valid)
		assert.True(t.T(), res.DeletedBy.Valid)

		assert.Equal(t.T(), 2, len(res.Balances))
		for _, resBalance := range res.Balances {
			assert.True(t.T(), resBalance.Deleted.Valid)
			assert.True(t.T(), resBalance.DeletedBy.Valid)
		}
	})

	t.Run("ErrorResolvingByIDs", func() {
		errMsg := "failed resolving by IDs"

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{}, errors.New(errMsg))

		res, err := t.svc.Delete(testBankAccountID, testUserID)

		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), errMsg)
		assert.Nil(t.T(), res)
	})

	t.Run("AccountNotFound", func() {
		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{}, nil)

		res, err := t.svc.Delete(testBankAccountID, testUserID)

		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), "EntityNotFound")
		assert.Nil(t.T(), res)
	})

	t.Run("AccountAlreadyDeleted", func() {

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{testDeletedBankAccount}, nil)

		res, err := t.svc.Delete(testBankAccountID, testUserID)

		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), "delete")
		assert.Contains(t.T(), err.Error(), "Bank Account")
		assert.Contains(t.T(), err.Error(), "deleted")
		assert.Nil(t.T(), res)
	})

	t.Run("AccountBalanceAlreadyDeleted", func() {

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testBankAccountID}).
			Return([]model.BankAccount{testBankAccount}, nil)

		t.mockRepo.EXPECT().ResolveBalancesByFilter(gomock.Any()).
			Return(
				[]model.BankAccountBalance{
					testDeletedNonLastBankAccountBalance,
					testLastBankAccountBalance,
				},
				model.PageInfoOutput{},
				nil)

		res, err := t.svc.Delete(testBankAccountID, testUserID)

		assert.Error(t.T(), err)
		assert.Contains(t.T(), err.Error(), "delete")
		assert.Contains(t.T(), err.Error(), "Bank Account Balance")
		assert.Contains(t.T(), err.Error(), "deleted")
		assert.Nil(t.T(), res)
	})
}
