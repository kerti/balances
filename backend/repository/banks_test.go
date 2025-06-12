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
	bankTestNow        = time.Now()
	banksTestYesterday = time.Now().AddDate(0, 0, -1)
	banksTestUserID, _ = uuid.NewV7()
)

// bank accounts
var (
	bankAccountsStmtInsert = `INSERT INTO bank_accounts
	( entity_id, account_name, bank_name, account_holder_name, account_number, last_balance, last_balance_date, status, created, created_by, updated, updated_by, deleted, deleted_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`

	bankAccountsStmtUpdate = `
	UPDATE bank_accounts
	SET account_name = ?, bank_name = ?, account_holder_name = ?, account_number = ?, last_balance = ?, last_balance_date = ?, status = ?, created = ?, created_by = ?, updated = ?, updated_by = ?, deleted = ?, deleted_by = ?
	WHERE entity_id = ?`
)

// bank account balances
var (
	bankAccountBalancesStmtInsert = `INSERT INTO bank_account_balances
	( entity_id, bank_account_entity_id, date, balance, created, created_by, updated, updated_by, deleted, deleted_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )`

	bankAccountBalancesStmtUpdate = `
	UPDATE bank_account_balances
	SET bank_account_entity_id = ?, date = ?, balance = ?, created = ?, created_by = ?, updated = ?, updated_by = ?, deleted = ?, deleted_by = ?
	WHERE entity_id = ?`
)

var (
	banksTestAccountBalanceID1, _ = uuid.NewV7()
	banksTestAccountBalanceID2, _ = uuid.NewV7()

	banksTestBankAccountBalanceModel1 = model.BankAccountBalance{
		ID:            banksTestAccountBalanceID1,
		BankAccountID: banksTestAccountID1,
		Date:          bankTestNow,
		Balance:       float64(1000000),
		Created:       bankTestNow,
		CreatedBy:     banksTestUserID,
	}
	banksTestBankAccountBalanceModel2 = model.BankAccountBalance{
		ID:            banksTestAccountBalanceID2,
		BankAccountID: banksTestAccountID1,
		Date:          banksTestYesterday,
		Balance:       float64(1100000),
		Created:       banksTestYesterday,
		CreatedBy:     banksTestUserID,
	}
)

var (
	banksTestAccountID1, _    = uuid.NewV7()
	banksTestAccountID2, _    = uuid.NewV7()
	banksTestBankAccountModel = model.BankAccount{
		ID:                banksTestAccountID1,
		AccountName:       "Savings Account",
		BankName:          "First National Bank",
		AccountHolderName: "John Doe",
		AccountNumber:     "12345678790",
		LastBalance:       float64(1000000),
		LastBalanceDate:   bankTestNow,
		Status:            model.BankAccountStatusActive,
		Created:           bankTestNow,
		CreatedBy:         banksTestUserID,
		Balances:          []model.BankAccountBalance{banksTestBankAccountBalanceModel1},
	}
)

func TestBanksRepository(t *testing.T) {

	t.Run("createBankAccount", func(t *testing.T) {

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
					banksTestBankAccountBalanceModel1.ID,
					banksTestBankAccountBalanceModel1.BankAccountID,
					banksTestBankAccountBalanceModel1.Date,
					banksTestBankAccountBalanceModel1.Balance,
					banksTestBankAccountBalanceModel1.Created,
					banksTestBankAccountBalanceModel1.CreatedBy,
					banksTestBankAccountBalanceModel1.Updated,
					banksTestBankAccountBalanceModel1.UpdatedBy,
					banksTestBankAccountBalanceModel1.Deleted,
					banksTestBankAccountBalanceModel1.DeletedBy,
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
					banksTestBankAccountBalanceModel1.ID,
					banksTestBankAccountBalanceModel1.BankAccountID,
					banksTestBankAccountBalanceModel1.Date,
					banksTestBankAccountBalanceModel1.Balance,
					banksTestBankAccountBalanceModel1.Created,
					banksTestBankAccountBalanceModel1.CreatedBy,
					banksTestBankAccountBalanceModel1.Updated,
					banksTestBankAccountBalanceModel1.UpdatedBy,
					banksTestBankAccountBalanceModel1.Deleted,
					banksTestBankAccountBalanceModel1.DeletedBy,
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

	t.Run("createBankAccountBalance", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestBankAccountBalanceModel2.ID.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountBalancesStmtInsert).
				ExpectExec().
				WithArgs(
					banksTestBankAccountBalanceModel2.ID,
					banksTestBankAccountBalanceModel2.BankAccountID,
					banksTestBankAccountBalanceModel2.Date,
					banksTestBankAccountBalanceModel2.Balance,
					banksTestBankAccountBalanceModel2.Created,
					banksTestBankAccountBalanceModel2.CreatedBy,
					nil,
					nil,
					nil,
					nil,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.
				ExpectPrepare(bankAccountsStmtUpdate).
				ExpectExec().
				WithArgs(
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
					banksTestBankAccountModel.ID,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.CreateBalance(banksTestBankAccountBalanceModel2, &banksTestBankAccountModel)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("normalNoAccountUpdate", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestBankAccountBalanceModel2.ID.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountBalancesStmtInsert).
				ExpectExec().
				WithArgs(
					banksTestBankAccountBalanceModel2.ID,
					banksTestBankAccountBalanceModel2.BankAccountID,
					banksTestBankAccountBalanceModel2.Date,
					banksTestBankAccountBalanceModel2.Balance,
					banksTestBankAccountBalanceModel2.Created,
					banksTestBankAccountBalanceModel2.CreatedBy,
					nil,
					nil,
					nil,
					nil,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.CreateBalance(banksTestBankAccountBalanceModel2, nil)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)

		})

	})

	t.Run("existsBankAccountByID", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			result := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(result)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ExistsByID(banksTestAccountID1)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("error", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ExistsByID(banksTestAccountID1)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("existsBankAccountBalanceByID", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			result := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1.String()).
				WillReturnRows(result)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ExistsBalanceByID(banksTestAccountBalanceID1)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("error", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ExistsBalanceByID(banksTestAccountBalanceID1)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("resolveBankAccountByIDs", func(t *testing.T) {

		t.Run("normalNoID", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIDs([]uuid.UUID{})
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("normalSingleID", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{banksTestAccountID1}

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(banksTestAccountID1.String())

			mock.ExpectQuery(repository.QuerySelectBankAccount + " WHERE bank_accounts.entity_id IN (?)").
				WithArgs(banksTestAccountID1).
				WillReturnRows(result)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("normalMultipleIDs", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{banksTestAccountID1, banksTestAccountID2}

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(banksTestAccountID1.String()).
				AddRow(banksTestAccountID2.String())

			mock.ExpectQuery(repository.QuerySelectBankAccount+" WHERE bank_accounts.entity_id IN (?, ?)").
				WithArgs(banksTestAccountID1, banksTestAccountID2).
				WillReturnRows(result)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorExecutingSelect", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{banksTestAccountID1, banksTestAccountID2}

			mock.ExpectQuery(repository.QuerySelectBankAccount+" WHERE bank_accounts.entity_id IN (?, ?)").
				WithArgs(banksTestAccountID1, banksTestAccountID2).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("resolveBankAccountByFilter", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			keyword := "example"
			likeKeyword := "%example%"

			countResult := sqlmock.NewRows([]string{"COUNT"}).AddRow(1)
			dataResult := sqlmock.NewRows([]string{"entity_id"}).AddRow(banksTestAccountID1)

			mock.
				ExpectQuery(repository.QuerySelectBankAccount+"WHERE (((((bank_accounts.account_name LIKE ?) OR (bank_accounts.bank_name LIKE ?)) OR (bank_accounts.account_number LIKE ?)) OR (bank_accounts.account_holder_name LIKE ?))) AND bank_accounts.deleted IS NULL LIMIT ? OFFSET ?").
				WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
				WillReturnRows(dataResult)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) FROM bank_accounts WHERE (((((bank_accounts.account_name LIKE ?) OR (bank_accounts.bank_name LIKE ?)) OR (bank_accounts.account_number LIKE ?)) OR (bank_accounts.account_holder_name LIKE ?))) AND bank_accounts.deleted IS NULL").
				WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword).
				WillReturnRows(countResult)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			testFilter := model.BankAccountFilterInput{}
			testFilter.Keyword = &keyword

			repo.Startup()
			_, _, err := repo.ResolveByFilter(testFilter.ToFilter())
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorOnSelect", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			keyword := "example"
			likeKeyword := "%example%"

			mock.
				ExpectQuery(repository.QuerySelectBankAccount+"WHERE (((((bank_accounts.account_name LIKE ?) OR (bank_accounts.bank_name LIKE ?)) OR (bank_accounts.account_number LIKE ?)) OR (bank_accounts.account_holder_name LIKE ?))) AND bank_accounts.deleted IS NULL LIMIT ? OFFSET ?").
				WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			testFilter := model.BankAccountFilterInput{}
			testFilter.Keyword = &keyword

			repo.Startup()
			_, _, err := repo.ResolveByFilter(testFilter.ToFilter())
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorOnCount", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			keyword := "example"
			likeKeyword := "%example%"

			dataResult := sqlmock.NewRows([]string{"entity_id"}).AddRow(banksTestAccountID1)

			mock.
				ExpectQuery(repository.QuerySelectBankAccount+"WHERE (((((bank_accounts.account_name LIKE ?) OR (bank_accounts.bank_name LIKE ?)) OR (bank_accounts.account_number LIKE ?)) OR (bank_accounts.account_holder_name LIKE ?))) AND bank_accounts.deleted IS NULL LIMIT ? OFFSET ?").
				WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword, 10, 0).
				WillReturnRows(dataResult)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) FROM bank_accounts WHERE (((((bank_accounts.account_name LIKE ?) OR (bank_accounts.bank_name LIKE ?)) OR (bank_accounts.account_number LIKE ?)) OR (bank_accounts.account_holder_name LIKE ?))) AND bank_accounts.deleted IS NULL").
				WithArgs(likeKeyword, likeKeyword, likeKeyword, likeKeyword).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			testFilter := model.BankAccountFilterInput{}
			testFilter.Keyword = &keyword

			repo.Startup()
			_, _, err := repo.ResolveByFilter(testFilter.ToFilter())
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("resolveBankAccountBalancesByIDs", func(t *testing.T) {

		t.Run("normalNoID", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveBalancesByIDs([]uuid.UUID{})
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("normalSingleID", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{banksTestAccountBalanceID1}

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(banksTestAccountBalanceID1.String())

			mock.ExpectQuery(repository.QuerySelectBankAccountBalance + " WHERE bank_account_balances.entity_id IN (?)").
				WithArgs(banksTestAccountBalanceID1).
				WillReturnRows(result)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveBalancesByIDs(ids)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("normalMultipleIDs", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{banksTestAccountBalanceID1, banksTestAccountBalanceID2}

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(banksTestAccountBalanceID1.String()).
				AddRow(banksTestAccountBalanceID2.String())

			mock.ExpectQuery(repository.QuerySelectBankAccountBalance+" WHERE bank_account_balances.entity_id IN (?, ?)").
				WithArgs(banksTestAccountBalanceID1, banksTestAccountBalanceID2).
				WillReturnRows(result)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveBalancesByIDs(ids)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorExecutingSelect", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{banksTestAccountBalanceID1, banksTestAccountBalanceID2}

			mock.ExpectQuery(repository.QuerySelectBankAccountBalance+" WHERE bank_account_balances.entity_id IN (?, ?)").
				WithArgs(banksTestAccountBalanceID1, banksTestAccountBalanceID2).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveBalancesByIDs(ids)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("resolveBankAccountBalancesByFilter", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			countResult := sqlmock.NewRows([]string{"COUNT"}).AddRow(1)
			dataResult := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(banksTestAccountBalanceID1)

			mock.
				ExpectQuery(repository.QuerySelectBankAccountBalance+"WHERE ((bank_account_balances.bank_account_entity_id IN (?))) AND bank_account_balances.deleted IS NULL LIMIT ? OFFSET ?").
				WithArgs(banksTestAccountID1, 10, 0).
				WillReturnRows(dataResult)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) FROM bank_account_balances WHERE ((bank_account_balances.bank_account_entity_id IN (?))) AND bank_account_balances.deleted IS NULL").
				WithArgs(banksTestAccountID1).
				WillReturnRows(countResult)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			testFilter := model.BankAccountBalanceFilterInput{}
			testFilter.BankAccountIDs = &[]uuid.UUID{banksTestAccountID1}

			repo.Startup()
			_, _, err := repo.ResolveBalancesByFilter(testFilter.ToFilter())
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorOnSelect", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery(repository.QuerySelectBankAccountBalance+"WHERE ((bank_account_balances.bank_account_entity_id IN (?))) AND bank_account_balances.deleted IS NULL LIMIT ? OFFSET ?").
				WithArgs(banksTestAccountID1, 10, 0).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			testFilter := model.BankAccountBalanceFilterInput{}
			testFilter.BankAccountIDs = &[]uuid.UUID{banksTestAccountID1}

			repo.Startup()
			_, _, err := repo.ResolveBalancesByFilter(testFilter.ToFilter())
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorOnCount", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			dataResult := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(banksTestAccountBalanceID1)

			mock.
				ExpectQuery(repository.QuerySelectBankAccountBalance+"WHERE ((bank_account_balances.bank_account_entity_id IN (?))) AND bank_account_balances.deleted IS NULL LIMIT ? OFFSET ?").
				WithArgs(banksTestAccountID1, 10, 0).
				WillReturnRows(dataResult)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) FROM bank_account_balances WHERE ((bank_account_balances.bank_account_entity_id IN (?))) AND bank_account_balances.deleted IS NULL").
				WithArgs(banksTestAccountID1).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			testFilter := model.BankAccountBalanceFilterInput{}
			testFilter.BankAccountIDs = &[]uuid.UUID{banksTestAccountID1}

			repo.Startup()
			_, _, err := repo.ResolveBalancesByFilter(testFilter.ToFilter())
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("resolveBankAccountLastBalancesByBankAccountID", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			count := 2
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(banksTestAccountBalanceID2).
				AddRow(banksTestAccountBalanceID1)

			mock.
				ExpectQuery(repository.QuerySelectBankAccountBalance+"WHERE bank_account_balances.bank_account_entity_id = ? ORDER BY bank_account_balances.date DESC LIMIT ?").
				WithArgs(banksTestAccountID1, count).
				WillReturnRows(result)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveLastBalancesByBankAccountID(banksTestAccountID1, count)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("countZero", func(t *testing.T) {
			count := 0
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveLastBalancesByBankAccountID(banksTestAccountID1, count)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnSelect", func(t *testing.T) {
			count := 2
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery(repository.QuerySelectBankAccountBalance+"WHERE bank_account_balances.bank_account_entity_id = ? ORDER BY bank_account_balances.date DESC LIMIT ?").
				WithArgs(banksTestAccountID1, count).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveLastBalancesByBankAccountID(banksTestAccountID1, count)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("updateBankAccount", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountsStmtUpdate).
				ExpectExec().
				WithArgs(
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
					banksTestBankAccountModel.ID,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(banksTestBankAccountModel)
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
			err := repo.Update(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("doesNotExist", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnPrepare", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountsStmtUpdate).
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnExec", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?").
				WithArgs(banksTestAccountID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountsStmtUpdate).
				ExpectExec().
				WithArgs(
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
					banksTestBankAccountModel.ID,
				).
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("updateBankAccountBalance", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.NewRows([]string{"COUNT"}).AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountBalancesStmtUpdate).
				ExpectExec().
				WithArgs().
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.
				ExpectPrepare(bankAccountsStmtUpdate).
				ExpectExec().
				WithArgs().
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.UpdateBalance(banksTestBankAccountBalanceModel1, &banksTestBankAccountModel)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("normalNoAccountUpdate", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.NewRows([]string{"COUNT"}).AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountBalancesStmtUpdate).
				ExpectExec().
				WithArgs().
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.UpdateBalance(banksTestBankAccountBalanceModel1, nil)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorOnCheckExistence", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1).
				WillReturnError(errors.New(""))

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.UpdateBalance(banksTestBankAccountBalanceModel1, &banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("doesNotExist", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.NewRows([]string{"COUNT"}).AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1).
				WillReturnRows(checkExistenceResult)

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.UpdateBalance(banksTestBankAccountBalanceModel1, &banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnPrepare", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.NewRows([]string{"COUNT"}).AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountBalancesStmtUpdate).
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.UpdateBalance(banksTestBankAccountBalanceModel1, &banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnExec", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.NewRows([]string{"COUNT"}).AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountBalancesStmtUpdate).
				ExpectExec().
				WithArgs().
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.UpdateBalance(banksTestBankAccountBalanceModel1, &banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnExecAccountUpdate", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.NewRows([]string{"COUNT"}).AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?").
				WithArgs(banksTestAccountBalanceID1).
				WillReturnRows(checkExistenceResult)

			mock.ExpectBegin()

			mock.
				ExpectPrepare(bankAccountBalancesStmtUpdate).
				ExpectExec().
				WithArgs().
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.
				ExpectPrepare(bankAccountsStmtUpdate).
				ExpectExec().
				WithArgs().
				WillReturnError(errors.New(""))

			mock.ExpectRollback()

			repo := new(repository.BankAccountMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.UpdateBalance(banksTestBankAccountBalanceModel1, &banksTestBankAccountModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

}
