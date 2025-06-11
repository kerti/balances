package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/stretchr/testify/assert"
)

// common
var (
	banksTestNow       = time.Now()
	banksTestUserID, _ = uuid.NewV7()
)

// bank accounts
var (
	bankAccountsStmtInsert = `INSERT INTO bank_accounts
	( entity_id, account_name, bank_name, account_holder_name, account_number, last_balance, last_balance_date, status, created, created_by, updated, updated_by, deleted, deleted_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`
)

// bank account balances
var (
	bankAccountBalancesStmtInsert = `INSERT INTO bank_account_balances
	( entity_id, bank_account_entity_id, date, balance, created, created_by, updated, updated_by, deleted, deleted_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`
)

var (
	banksTestAccountBalanceID1, _    = uuid.NewV7()
	banksTestBankAccountBalanceModel = model.BankAccountBalance{
		ID:            banksTestAccountBalanceID1,
		BankAccountID: banksTestAccountID1,
		Date:          banksTestNow,
		Balance:       float64(1000000),
		Created:       banksTestNow,
		CreatedBy:     banksTestUserID,
	}
)

var (
	banksTestAccountID1, _ = uuid.NewV7()
	// banksTestAccountID2, _    = uuid.NewV7()
	banksTestBankAccountModel = model.BankAccount{
		ID:                banksTestAccountID1,
		AccountName:       "Savings Account",
		BankName:          "First National Bank",
		AccountHolderName: "John Doe",
		AccountNumber:     "12345678790",
		LastBalance:       float64(1000000),
		LastBalanceDate:   banksTestNow,
		Status:            model.BankAccountStatusActive,
		Created:           banksTestNow,
		CreatedBy:         banksTestUserID,
		Balances:          []model.BankAccountBalance{banksTestBankAccountBalanceModel},
	}
)

func TestBanksRepository(t *testing.T) {

	t.Run("create", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountsStmtInsert).
				ExpectExec().
				WithArgs(
					banksTestBankAccountModel.ID,
					banksTestBankAccountModel.AccountName,
					banksTestBankAccountModel.BankName,
					banksTestBankAccountModel.AccountHolderName,
					banksTestBankAccountModel.AccountNumber,
					banksTestBankAccountModel.LastBalance,
					banksTestBankAccountModel.LastBalanceDate,
					banksTestBankAccountModel.Status,
					banksTestBankAccountModel.Created,
					banksTestBankAccountModel.CreatedBy,
					banksTestBankAccountModel.Updated,
					banksTestBankAccountModel.UpdatedBy,
					banksTestBankAccountModel.Deleted,
					banksTestBankAccountModel.DeletedBy,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.
				ExpectPrepare(bankAccountBalancesStmtInsert).
				ExpectExec().
				WithArgs(
					banksTestBankAccountBalanceModel.ID,
					banksTestBankAccountBalanceModel.BankAccountID,
					banksTestBankAccountBalanceModel.Date,
					banksTestBankAccountBalanceModel.Balance,
					banksTestBankAccountBalanceModel.Created,
					banksTestBankAccountBalanceModel.CreatedBy,
					banksTestBankAccountBalanceModel.Updated,
					banksTestBankAccountBalanceModel.UpdatedBy,
					banksTestBankAccountBalanceModel.Deleted,
					banksTestBankAccountBalanceModel.DeletedBy,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(banksTestBankAccountModel)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorOnCheckExistence", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("alreadyExists", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)
			assert.IsType(t, &failure.Failure{}, err)
			assert.Equal(t, failure.CodeOperationNotPermitted, err.(*failure.Failure).Code)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnPrepareBankAccountStatement", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountsStmtInsert).
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnExecBankAccountStatement", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountsStmtInsert).
				ExpectExec().
				WithArgs(
					banksTestBankAccountModel.ID,
					banksTestBankAccountModel.AccountName,
					banksTestBankAccountModel.BankName,
					banksTestBankAccountModel.AccountHolderName,
					banksTestBankAccountModel.AccountNumber,
					banksTestBankAccountModel.LastBalance,
					banksTestBankAccountModel.LastBalanceDate,
					banksTestBankAccountModel.Status,
					banksTestBankAccountModel.Created,
					banksTestBankAccountModel.CreatedBy,
					banksTestBankAccountModel.Updated,
					banksTestBankAccountModel.UpdatedBy,
					banksTestBankAccountModel.Deleted,
					banksTestBankAccountModel.DeletedBy,
				).
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnPrepareBankAccountBalanceStatement", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountsStmtInsert).
				ExpectExec().
				WithArgs(
					banksTestBankAccountModel.ID,
					banksTestBankAccountModel.AccountName,
					banksTestBankAccountModel.BankName,
					banksTestBankAccountModel.AccountHolderName,
					banksTestBankAccountModel.AccountNumber,
					banksTestBankAccountModel.LastBalance,
					banksTestBankAccountModel.LastBalanceDate,
					banksTestBankAccountModel.Status,
					banksTestBankAccountModel.Created,
					banksTestBankAccountModel.CreatedBy,
					banksTestBankAccountModel.Updated,
					banksTestBankAccountModel.UpdatedBy,
					banksTestBankAccountModel.Deleted,
					banksTestBankAccountModel.DeletedBy,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.
				ExpectPrepare(bankAccountBalancesStmtInsert).
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnExecBankAccountBalanceStatement", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountsStmtInsert).
				ExpectExec().
				WithArgs(
					banksTestBankAccountModel.ID,
					banksTestBankAccountModel.AccountName,
					banksTestBankAccountModel.BankName,
					banksTestBankAccountModel.AccountHolderName,
					banksTestBankAccountModel.AccountNumber,
					banksTestBankAccountModel.LastBalance,
					banksTestBankAccountModel.LastBalanceDate,
					banksTestBankAccountModel.Status,
					banksTestBankAccountModel.Created,
					banksTestBankAccountModel.CreatedBy,
					banksTestBankAccountModel.Updated,
					banksTestBankAccountModel.UpdatedBy,
					banksTestBankAccountModel.Deleted,
					banksTestBankAccountModel.DeletedBy,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.
				ExpectPrepare(bankAccountBalancesStmtInsert).
				ExpectExec().
				WithArgs(
					banksTestBankAccountBalanceModel.ID,
					banksTestBankAccountBalanceModel.BankAccountID,
					banksTestBankAccountBalanceModel.Date,
					banksTestBankAccountBalanceModel.Balance,
					banksTestBankAccountBalanceModel.Created,
					banksTestBankAccountBalanceModel.CreatedBy,
					banksTestBankAccountBalanceModel.Updated,
					banksTestBankAccountBalanceModel.UpdatedBy,
					banksTestBankAccountBalanceModel.Deleted,
					banksTestBankAccountBalanceModel.DeletedBy,
				).
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})
	})

}
