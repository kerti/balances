package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
}

func (t *bankAccountsServiceTestSuite) TearDownTest() {
	t.ctrl.Finish()
}

func (t *bankAccountsServiceTestSuite) TestResolveByID() {

	t.Run("TestExistsNoBalance", func() {
		testID, _ := uuid.NewV7()
		nilCacheTime := cachetime.NCacheTime{}
		bankAccounts := []model.BankAccount{
			{
				ID: testID,
			},
		}

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).Return(bankAccounts, nil)

		_, err := t.svc.GetByID(testID, false, nilCacheTime, nilCacheTime, nil)

		assert.NoError(t.T(), err, "should not error")
	})

	t.Run("TestExistsWithBalance", func() {
		testID, _ := uuid.NewV7()
		testBalanceID1, _ := uuid.NewV7()
		testBalanceID2, _ := uuid.NewV7()
		nilCacheTime := cachetime.NCacheTime{}

		bankAccounts := []model.BankAccount{
			{
				ID: testID,
			},
		}
		bankAccountBalances := []model.BankAccountBalance{
			{
				ID: testBalanceID1,
			},
			{
				ID: testBalanceID2,
			},
		}
		balanceFilterInput := model.BankAccountBalanceFilterInput{
			BankAccountIDs: &[]uuid.UUID{testID},
		}
		pageInfo := model.PageInfoOutput{}

		t.mockRepo.EXPECT().ResolveByIDs([]uuid.UUID{testID}).Return(bankAccounts, nil)
		t.mockRepo.EXPECT().ResolveBalancesByFilter(balanceFilterInput.ToFilter()).Return(bankAccountBalances, pageInfo, nil)

		_, err := t.svc.GetByID(testID, true, nilCacheTime, nilCacheTime, nil)

		assert.NoError(t.T(), err, "should not error")
	})

}
