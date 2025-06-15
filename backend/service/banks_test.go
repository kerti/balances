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

func (t *bankAccountsServiceTestSuite) getBankAccountSlice(count int) (res []model.BankAccount) {
	for range count {
		id, _ := uuid.NewV7()
		res = append(res, t.getNewBankAccount(nuuid.From(id), nil))
	}
	return
}

func (t *bankAccountsServiceTestSuite) getNewBankAccountBalance(id nuuid.NUUID, bankAccountID nuuid.NUUID, balance float64, date time.Time) model.BankAccountBalance {
	bal := model.BankAccountBalance{}

	if id.Valid {
		bal.ID = id.UUID
	} else {
		newID, _ := uuid.NewV7()
		bal.ID = newID
	}

	if bankAccountID.Valid {
		bal.BankAccountID = bankAccountID.UUID
	} else {
		newAccID, _ := uuid.NewV7()
		bal.BankAccountID = newAccID
	}

	bal.Balance = balance
	bal.Date = date

	return bal
}

func (t *bankAccountsServiceTestSuite) getDefaultPageInfo() model.PageInfoOutput {
	return model.PageInfoOutput{
		Page:       1,
		PageSize:   10,
		TotalCount: 1,
		PageCount:  1,
	}
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

func (t *bankAccountsServiceTestSuite) TestGetByID_Exists_NoBalance() {
	bankAccount := t.getNewBankAccount(nuuid.NUUID{}, nil)
	resolvedBankAccountSlice := []model.BankAccount{bankAccount}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return(resolvedBankAccountSlice, nil)

	_, err := t.svc.GetByID(t.testBankAccountID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_Exists_WithBalance_NoFilter() {
	balanceFilterInput := model.BankAccountBalanceFilterInput{
		BankAccountIDs: &[]uuid.UUID{t.testBankAccountID},
	}
	pageInfo := t.getDefaultPageInfo()

	bankAccount := t.getNewBankAccount(nuuid.NUUID{}, nil)
	balance1 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	balance2 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(2000), time.Now())
	balanceSlice := []model.BankAccountBalance{balance1, balance2}
	bankAccount.Balances = balanceSlice

	resolvedBankAccountSlice := []model.BankAccount{bankAccount}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return(resolvedBankAccountSlice, nil)
	t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
		Return(balanceSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testBankAccountID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_Exists_WithBalance_WithStartDate() {
	yesterday := cachetime.NCacheTime(null.TimeFrom(time.Now().AddDate(0, 0, -1)))
	balanceFilterInput := model.BankAccountBalanceFilterInput{
		BankAccountIDs: &[]uuid.UUID{t.testBankAccountID},
		StartDate:      yesterday,
	}
	pageInfo := t.getDefaultPageInfo()

	bankAccount := t.getNewBankAccount(nuuid.NUUID{}, nil)
	balance1 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	balance2 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(2000), time.Now())
	balanceSlice := []model.BankAccountBalance{balance1, balance2}
	bankAccount.Balances = balanceSlice

	resolvedBankAccountSlice := []model.BankAccount{bankAccount}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return(resolvedBankAccountSlice, nil)
	t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
		Return(balanceSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testBankAccountID, true, yesterday, cachetime.NCacheTime{}, nil)

	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_Exists_WithBalance_WithEndDate() {
	today := cachetime.NCacheTime(null.TimeFrom(time.Now()))
	balanceFilterInput := model.BankAccountBalanceFilterInput{
		BankAccountIDs: &[]uuid.UUID{t.testBankAccountID},
		EndDate:        today,
	}
	pageInfo := t.getDefaultPageInfo()

	bankAccount := t.getNewBankAccount(nuuid.NUUID{}, nil)
	balanceSlice := []model.BankAccountBalance{
		t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(2000), time.Now()),
	}
	bankAccount.Balances = balanceSlice

	resolvedBankAccountSlice := []model.BankAccount{bankAccount}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return(resolvedBankAccountSlice, nil)
	t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
		Return(balanceSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testBankAccountID, true, cachetime.NCacheTime{}, today, nil)

	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_Exists_WithBalance_WithBothDates() {
	yesterday := cachetime.NCacheTime(null.TimeFrom(time.Now().AddDate(0, 0, -1)))
	today := cachetime.NCacheTime(null.TimeFrom(time.Now()))
	balanceFilterInput := model.BankAccountBalanceFilterInput{
		BankAccountIDs: &[]uuid.UUID{t.testBankAccountID},
		StartDate:      yesterday,
		EndDate:        today,
	}
	pageInfo := t.getDefaultPageInfo()

	bankAccount := t.getNewBankAccount(nuuid.NUUID{}, nil)
	balance1 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	balance2 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(2000), time.Now())
	balanceSlice := []model.BankAccountBalance{balance1, balance2}
	bankAccount.Balances = balanceSlice

	resolvedBankAccountSlice := []model.BankAccount{bankAccount}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return(resolvedBankAccountSlice, nil)
	t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
		Return(balanceSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testBankAccountID, true, yesterday, today, nil)

	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_Exists_WithBalance_WithPageSize() {
	pageSize := 120
	balanceFilterInput := model.BankAccountBalanceFilterInput{
		BankAccountIDs: &[]uuid.UUID{t.testBankAccountID},
	}
	balanceFilterInput.PageSize = &pageSize
	pageInfo := model.PageInfoOutput{
		PageSize: pageSize,
	}

	bankAccount := t.getNewBankAccount(nuuid.NUUID{}, nil)
	balance1 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	balance2 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(2000), time.Now())
	balanceSlice := []model.BankAccountBalance{balance1, balance2}
	bankAccount.Balances = balanceSlice

	resolvedBankAccountSlice := []model.BankAccount{bankAccount}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return(resolvedBankAccountSlice, nil)
	t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
		Return(balanceSlice, pageInfo, nil)

	_, err := t.svc.GetByID(t.testBankAccountID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, &pageSize)

	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_Exists_WithBalance_RepoErrorResolvingAccount() {
	errMsg := "failed to resolve bank account by IDs"
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{}, errors.New(errMsg))

	res, err := t.svc.GetByID(t.testBankAccountID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_Exists_WithBalance_RepoErrorResolvingBalances() {
	errMsg := "error resolving balances"
	balanceFilterInput := model.BankAccountBalanceFilterInput{
		BankAccountIDs: &[]uuid.UUID{t.testBankAccountID},
	}

	bankAccount := t.getNewBankAccount(nuuid.NUUID{}, nil)
	balance1 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(1000), time.Now().AddDate(0, 0, -1))
	balance2 := t.getNewBankAccountBalance(nuuid.NUUID{}, nuuid.From(bankAccount.ID), float64(2000), time.Now())
	balanceSlice := []model.BankAccountBalance{balance1, balance2}
	bankAccount.Balances = balanceSlice
	resolvedBankAccountSlice := []model.BankAccount{bankAccount}

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return(resolvedBankAccountSlice, nil)
	t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).
		Return([]model.BankAccountBalance{}, t.getDefaultPageInfo(), errors.New(errMsg))

	res, err := t.svc.GetByID(t.testBankAccountID, true, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_ErrorResolvingByID() {
	errMsg := "repo failed resolving by IDs"
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return(nil, errors.New(errMsg))

	_, err := t.svc.GetByID(t.testBankAccountID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.Error(t.T(), err)
	assert.Equal(t.T(), err.Error(), errMsg)
}

func (t *bankAccountsServiceTestSuite) TestGetByID_NotExists() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{}, nil)

	_, err := t.svc.GetByID(t.testBankAccountID, false, cachetime.NCacheTime{}, cachetime.NCacheTime{}, nil)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
}

func (t *bankAccountsServiceTestSuite) TestGetByFilter_EmptyFilter() {
	filterInput := model.BankAccountFilterInput{}
	filter := filterInput.ToFilter()

	t.mockRepo.EXPECT().ResolveByFilter(filter).
		Return(t.getBankAccountSlice(2), t.getDefaultPageInfo(), nil)

	_, _, err := t.svc.GetByFilter(filterInput)

	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestGetByFilter_WithKeyword() {
	keyword := "example"
	filterInput := model.BankAccountFilterInput{}
	filterInput.Keyword = &keyword
	filter := filterInput.ToFilter()

	t.mockRepo.EXPECT().ResolveByFilter(filter).
		Return(t.getBankAccountSlice(2), t.getDefaultPageInfo(), nil)

	_, _, err := t.svc.GetByFilter(filterInput)

	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestUpdate_Normal() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{t.getNewBankAccount(nuuid.From(t.testBankAccountID), nil)}, nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).
		Return(nil)

	res, err := t.svc.Update(t.getNewBankAccountInput(nuuid.From(t.testBankAccountID), nil), t.testUserID)

	assert.NotNil(t.T(), res)
	assert.NoError(t.T(), err)
}

func (t *bankAccountsServiceTestSuite) TestUpdate_ErrorResolvingByIDs() {
	errMsg := "query failed"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{}, errors.New(errMsg))

	res, err := t.svc.Update(t.getNewBankAccountInput(nuuid.From(t.testBankAccountID), nil), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *bankAccountsServiceTestSuite) TestUpdate_AccountNotFound() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{}, nil)

	res, err := t.svc.Update(t.getNewBankAccountInput(nuuid.From(t.testBankAccountID), nil), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
}

func (t *bankAccountsServiceTestSuite) TestUpdate_AccountDeleted() {
	bankAccountInput := t.getNewBankAccountInput(nuuid.From(t.testBankAccountID), nil)
	deletedBankAccount := model.NewBankAccountFromInput(bankAccountInput, t.testUserID)
	deletedBankAccount.Delete(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{deletedBankAccount}, nil)

	res, err := t.svc.Update(bankAccountInput, t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "update")
	assert.Contains(t.T(), err.Error(), "Bank Account")
	assert.Contains(t.T(), err.Error(), "deleted")
}

func (t *bankAccountsServiceTestSuite) TestUpdate_RepoFailToUpdate() {
	errMsg := "failed to update"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{t.getNewBankAccount(nuuid.From(t.testBankAccountID), nil)}, nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).
		Return(errors.New(errMsg))

	res, err := t.svc.Update(t.getNewBankAccountInput(nuuid.From(t.testBankAccountID), nil), t.testUserID)

	assert.Nil(t.T(), res)
	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
}

func (t *bankAccountsServiceTestSuite) TestDelete_Normal() {
	balanceSlice := []model.BankAccountBalance{}

	balanceSlice = append(
		balanceSlice,
		t.getNewBankAccountBalance(
			nuuid.NUUID{},
			nuuid.From(t.testBankAccountID),
			float64(10000),
			time.Now().AddDate(0, 0, -1)))

	balanceSlice = append(
		balanceSlice,
		t.getNewBankAccountBalance(
			nuuid.NUUID{},
			nuuid.From(t.testBankAccountID),
			float64(12000),
			time.Now()))

	testBankAccount := t.getNewBankAccount(nuuid.From(t.testBankAccountID), &balanceSlice)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{testBankAccount}, nil)

	t.mockRepo.EXPECT().ResolveBalancesByFilter(gomock.Any()).
		Return(balanceSlice, t.getDefaultPageInfo(), nil)

	t.mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

	res, err := t.svc.Delete(t.testBankAccountID, t.testUserID)

	assert.NoError(t.T(), err)

	assert.NotNil(t.T(), res)
	assert.True(t.T(), res.Deleted.Valid)
	assert.True(t.T(), res.DeletedBy.Valid)

	assert.Equal(t.T(), 2, len(res.Balances))
	for _, resBalance := range res.Balances {
		assert.True(t.T(), resBalance.Deleted.Valid)
		assert.True(t.T(), resBalance.DeletedBy.Valid)
	}
}

func (t *bankAccountsServiceTestSuite) TestDelete_ErrorResolvingByIDs() {
	errMsg := "failed resolving by IDs"

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{}, errors.New(errMsg))

	res, err := t.svc.Delete(t.testBankAccountID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), errMsg)
	assert.Nil(t.T(), res)
}

func (t *bankAccountsServiceTestSuite) TestDelete_AccountNotFound() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{}, nil)

	res, err := t.svc.Delete(t.testBankAccountID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "EntityNotFound")
	assert.Nil(t.T(), res)
}

func (t *bankAccountsServiceTestSuite) TestDelete_AccountAlreadyDeleted() {
	testDeletedBankAccount := t.getNewBankAccount(nuuid.From(t.testBankAccountID), nil)
	testDeletedBankAccount.Deleted = null.TimeFrom(time.Now())
	testDeletedBankAccount.DeletedBy = nuuid.From(t.testUserID)

	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{testDeletedBankAccount}, nil)

	res, err := t.svc.Delete(t.testBankAccountID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "delete")
	assert.Contains(t.T(), err.Error(), "Bank Account")
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}

func (t *bankAccountsServiceTestSuite) TestDelete_AccountBalanceAlreadyDeleted() {
	t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{t.testBankAccountID}).
		Return([]model.BankAccount{t.getNewBankAccount(nuuid.From(t.testBankAccountID), nil)}, nil)

	testDeletedNonLastBankAccountBalance := t.getNewBankAccountBalance(
		nuuid.NUUID{},
		nuuid.From(t.testBankAccountID),
		float64(10000),
		time.Now().AddDate(0, 0, -1))

	testDeletedNonLastBankAccountBalance.Deleted = null.TimeFrom(time.Now())
	testDeletedNonLastBankAccountBalance.DeletedBy = nuuid.From(t.testUserID)

	testLastBankAccountBalance := t.getNewBankAccountBalance(
		nuuid.NUUID{},
		nuuid.From(t.testBankAccountID),
		float64(12000),
		time.Now())

	t.mockRepo.EXPECT().ResolveBalancesByFilter(gomock.Any()).
		Return(
			[]model.BankAccountBalance{
				testDeletedNonLastBankAccountBalance,
				testLastBankAccountBalance,
			},
			t.getDefaultPageInfo(),
			nil)

	res, err := t.svc.Delete(t.testBankAccountID, t.testUserID)

	assert.Error(t.T(), err)
	assert.Contains(t.T(), err.Error(), "delete")
	assert.Contains(t.T(), err.Error(), "Bank Account Balance")
	assert.Contains(t.T(), err.Error(), "deleted")
	assert.Nil(t.T(), res)
}
