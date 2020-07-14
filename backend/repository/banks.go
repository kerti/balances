package repository

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/balances/backend/database"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/filter"
	"github.com/kerti/balances/backend/util/logger"
)

const (
	querySelectBankAccount = `
		SELECT
			bank_accounts.entity_id,
			bank_accounts.account_name,
			bank_accounts.bank_name,
			bank_accounts.account_holder_name,
			bank_accounts.account_number,
			bank_accounts.last_balance,
			bank_accounts.last_balance_date,
			bank_accounts.status,
			bank_accounts.created,
			bank_accounts.created_by,
			bank_accounts.updated,
			bank_accounts.updated_by,
			bank_accounts.deleted,
			bank_accounts.deleted_by
		FROM
			bank_accounts `

	querySelectBankAccountBalance = `
		SELECT
			bank_account_balances.entity_id,
			bank_account_balances.bank_account_entity_id,
			bank_account_balances.date,
			bank_account_balances.balance,
			bank_account_balances.created,
			bank_account_balances.created_by,
			bank_account_balances.updated,
			bank_account_balances.updated_by,
			bank_account_balances.deleted,
			bank_account_balances.deleted_by
		FROM
			bank_account_balances `

	queryInsertBankAccount = `
		INSERT INTO bank_accounts (
			entity_id,
			account_name,
			bank_name,
			account_holder_name,
			account_number,
			last_balance,
			last_balance_date,
			status,
			created,
			created_by,
			updated,
			updated_by,
			deleted,
			deleted_by
		) VALUES (
			:entity_id,
			:account_name,
			:bank_name,
			:account_holder_name,
			:account_number,
			:last_balance,
			:last_balance_date,
			:status,
			:created,
			:created_by,
			:updated,
			:updated_by,
			:deleted,
			:deleted_by
		)`

	queryInsertBankAccountBalance = `
		INSERT INTO bank_account_balances (
			entity_id,
			bank_account_entity_id,
			date,
			balance,
			created,
			created_by,
			updated,
			updated_by,
			deleted,
			deleted_by
		) VALUES (
			:entity_id,
			:bank_account_entity_id,
			:date,
			:balance,
			:created,
			:created_by,
			:updated,
			:updated_by,
			:deleted,
			:deleted_by
		)`

	queryUpdateBankAccount = `
		UPDATE bank_accounts
		SET
			account_name = :account_name,
			bank_name = :bank_name,
			account_holder_name = :account_holder_name,
			account_number = :account_number,
			last_balance = :last_balance,
			last_balance_date = :last_balance_date,
			status = :status,
			created = :created,
			created_by = :created_by,
			updated = :updated,
			updated_by = :updated_by,
			deleted = :deleted,
			deleted_by = :deleted_by
		WHERE entity_id = :entity_id`

	queryUpdateBankAccountBalance = `
		UPDATE bank_account_balances
		SET
			entity_id = :entity_id,
			bank_account_entity_id = :bank_account_entity_id,
			date = :date,
			balance = :balance,
			created = :created,
			created_by = :created_by,
			updated = :updated,
			updated_by = :updated_by,
			deleted = :deleted,
			deleted_by = :deleted_by
		WHERE entity_id = :entity_id`
)

// BankAccount is the Bank Account repository interface
type BankAccount interface {
	Startup()
	Shutdown()
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ExistsBalanceByID(id uuid.UUID) (exists bool, err error)
	ResolveByIDs(ids []uuid.UUID) (bankAccounts []model.BankAccount, err error)
	ResolveBalancesByIDs(ids []uuid.UUID) (bankAccountBalances []model.BankAccountBalance, err error)
	ResolveByFilter(filter filter.Filter) (bankAccounts []model.BankAccount, pageInfo model.PageInfoOutput, err error)
	ResolveBalancesByFilter(filter filter.Filter) (bankAccountBalances []model.BankAccountBalance, pageInfo model.PageInfoOutput, err error)
	ResolveLastBalancesByBankAccountID(id uuid.UUID, count int) (bankAccountBalances []model.BankAccountBalance, err error)
	Create(bankAccount model.BankAccount) error
	Update(bankAccount model.BankAccount) error
	CreateBalance(bankAccountBalance model.BankAccountBalance, bankAccount *model.BankAccount) error
	UpdateBalance(bankAccountBalance model.BankAccountBalance, bankAccount *model.BankAccount) error
}

// BankAccountMySQLRepo is the repository for Bank Accounts implemented with MySQL backend
type BankAccountMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup perform startup functions
func (r *BankAccountMySQLRepo) Startup() {
	logger.Trace("BankAccount repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *BankAccountMySQLRepo) Shutdown() {
	logger.Trace("BankAccount repository shutting down...")
}

// ExistsByID checks the existence of a Bank Account by its ID
func (r *BankAccountMySQLRepo) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM bank_accounts WHERE bank_accounts.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
	}
	return
}

// ExistsBalanceByID checks the existence of a Bank Account Balance by its ID
func (r *BankAccountMySQLRepo) ExistsBalanceByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM bank_account_balances WHERE bank_account_balances.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
	}
	return
}

// ResolveByIDs resolves Bank Accounts by their IDs
func (r *BankAccountMySQLRepo) ResolveByIDs(ids []uuid.UUID) (bankAccounts []model.BankAccount, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(querySelectBankAccount+" WHERE bank_accounts.entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	err = r.DB.Select(&bankAccounts, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	return
}

// ResolveBalancesByIDs resoloves Bank Account Balances by their IDs
func (r *BankAccountMySQLRepo) ResolveBalancesByIDs(ids []uuid.UUID) (bankAccountBalances []model.BankAccountBalance, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(querySelectBankAccountBalance+" WHERE bank_account_balances.entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	err = r.DB.Select(&bankAccountBalances, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	return
}

// ResolveByFilter resolves Banks Accounts by a specified filter
func (r *BankAccountMySQLRepo) ResolveByFilter(filter filter.Filter) (bankAccounts []model.BankAccount, pageInfo model.PageInfoOutput, err error) {
	filterQueryString, err := filter.ToQueryString()
	if err != nil {
		return bankAccounts, pageInfo, err
	}

	filterArgs := filter.GetArgs(true)
	query, args, err := r.DB.In(
		querySelectBankAccount+filterQueryString+filter.Pagination.ToQueryString(),
		filterArgs...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	err = r.DB.Select(&bankAccounts, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	var count int
	filterArgsNoPagination := filter.GetArgs(false)
	err = r.DB.Get(
		&count,
		"SELECT COUNT(entity_id) FROM bank_accounts "+filterQueryString,
		filterArgsNoPagination...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	pageInfo = model.PageInfoOutput{
		Page:       filter.Pagination.Page,
		PageSize:   filter.Pagination.PageSize,
		TotalCount: count,
		TotalPages: filter.Pagination.GetPageCount(count),
	}

	return
}

// ResolveBalancesByFilter resolves Banks Account Balances by a specified filter
func (r *BankAccountMySQLRepo) ResolveBalancesByFilter(filter filter.Filter) (bankAccountBalances []model.BankAccountBalance, pageInfo model.PageInfoOutput, err error) {
	filterQueryString, err := filter.ToQueryString()
	if err != nil {
		return bankAccountBalances, pageInfo, err
	}

	filterArgs := filter.GetArgs(true)
	logger.Warn(querySelectBankAccountBalance + filterQueryString + filter.Pagination.ToQueryString())
	logger.Warn("%#v", filterArgs)
	query, args, err := r.DB.In(
		querySelectBankAccountBalance+filterQueryString+filter.Pagination.ToQueryString(),
		filterArgs...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	err = r.DB.Select(&bankAccountBalances, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	var count int
	filterArgsNoPagination := filter.GetArgs(false)
	err = r.DB.Get(
		&count,
		"SELECT COUNT(entity_id) FROM bank_account_balances "+filterQueryString,
		filterArgsNoPagination...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	pageInfo = model.PageInfoOutput{
		Page:       filter.Pagination.Page,
		PageSize:   filter.Pagination.PageSize,
		TotalCount: count,
		TotalPages: filter.Pagination.GetPageCount(count),
	}

	return
}

// ResolveLastBalancesByBankAccountID resolves last X Bank Account Balances by their Bank Account ID and count param
func (r *BankAccountMySQLRepo) ResolveLastBalancesByBankAccountID(id uuid.UUID, count int) (bankAccountBalances []model.BankAccountBalance, err error) {
	if count == 0 {
		return
	}

	whereClause := " WHERE bank_account_balances.bank_account_entity_id = ? ORDER BY bank_account_balances.date DESC LIMIT ?"
	query, args, err := r.DB.In(
		querySelectBankAccountBalance+whereClause, id, count)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	err = r.DB.Select(&bankAccountBalances, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	return
}

// Create creates a Bank Account
func (r *BankAccountMySQLRepo) Create(bankAccount model.BankAccount) error {
	exists, err := r.ExistsByID(bankAccount.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if exists {
		err = failure.OperationNotPermitted("create", "Bank Account", "already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreateBankAccount(tx, bankAccount); err != nil {
			e <- err
			return
		}

		if err := r.txCreateBankAccountBalance(tx, bankAccount.Balances[0]); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

// Update updates a bank account
func (r *BankAccountMySQLRepo) Update(bankAccount model.BankAccount) error {
	exists, err := r.ExistsByID(bankAccount.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if !exists {
		err = failure.EntityNotFound("Bank Account")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdateBankAccount(tx, bankAccount); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

// CreateBalance creates a new Bank Account Balance and optionally updates the Bank Account transactionally
func (r *BankAccountMySQLRepo) CreateBalance(bankAccountBalance model.BankAccountBalance, bankAccount *model.BankAccount) error {
	exists, err := r.ExistsBalanceByID(bankAccountBalance.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if exists {
		err = failure.OperationNotPermitted("create", "Bank Account Balance", "already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreateBankAccountBalance(tx, bankAccountBalance); err != nil {
			e <- err
			return
		}

		if bankAccount != nil {
			if err := r.txUpdateBankAccount(tx, *bankAccount); err != nil {
				e <- err
				return
			}
		}

		e <- nil
	})
}

// UpdateBalance updates an existing Bank Account Balance and optionally updates the Bank Account transactionally
func (r *BankAccountMySQLRepo) UpdateBalance(bankAccountBalance model.BankAccountBalance, bankAccount *model.BankAccount) error {
	exists, err := r.ExistsBalanceByID(bankAccountBalance.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if !exists {
		err = failure.EntityNotFound("Bank Account Balance")
		logger.ErrNoStack("%v", err)
		return err
	}

	return r.DB.WithTransaction(r.DB, func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdateBankAccountBalance(tx, bankAccountBalance); err != nil {
			e <- err
			return
		}

		if bankAccount != nil {
			if err := r.txUpdateBankAccount(tx, *bankAccount); err != nil {
				e <- err
				return
			}
		}

		e <- nil
	})
}

func (r *BankAccountMySQLRepo) txCreateBankAccount(tx *sqlx.Tx, bankAccount model.BankAccount) error {
	stmt, err := tx.PrepareNamed(queryInsertBankAccount)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(bankAccount)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}

func (r *BankAccountMySQLRepo) txCreateBankAccountBalance(tx *sqlx.Tx, bankAccountBalance model.BankAccountBalance) error {
	stmt, err := tx.PrepareNamed(queryInsertBankAccountBalance)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(bankAccountBalance)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}

func (r *BankAccountMySQLRepo) txUpdateBankAccount(tx *sqlx.Tx, bankAccount model.BankAccount) error {
	stmt, err := tx.PrepareNamed(queryUpdateBankAccount)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(bankAccount)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}

func (r *BankAccountMySQLRepo) txUpdateBankAccountBalance(tx *sqlx.Tx, bankAccountBalance model.BankAccountBalance) error {
	stmt, err := tx.PrepareNamed(queryUpdateBankAccountBalance)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(bankAccountBalance)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}
