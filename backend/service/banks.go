package service

import (
	"math"

	"github.com/google/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/logger"
)

// BankAccountImpl is the service provider implementation
type BankAccountImpl struct {
	Repository repository.BankAccount `inject:"bankAccountRepository"`
}

// Startup performs startup functions
func (s *BankAccountImpl) Startup() {
	logger.Trace("Bank Account Service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *BankAccountImpl) Shutdown() {
	logger.Trace("Bank Account Service shutting down...")
}

// Create creates a new Bank Account
func (s *BankAccountImpl) Create(input model.BankAccountInput, userID uuid.UUID) (*model.BankAccount, error) {
	bankAccount := model.NewBankAccountFromInput(input, userID)
	err := s.Repository.Create(bankAccount)
	if err != nil {
		return nil, err
	}
	return &bankAccount, err
}

// GetByID fetches a Bank Account by its ID
func (s *BankAccountImpl) GetByID(id uuid.UUID, withBalances bool, balanceStartDate, balanceEndDate cachetime.NCacheTime, pageSize *int) (*model.BankAccount, error) {
	bankAccounts, err := s.Repository.ResolveByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(bankAccounts) != 1 {
		return nil, failure.EntityNotFound("Bank Account")
	}

	bankAccount := bankAccounts[0]

	if withBalances {
		filter := model.BankAccountBalanceFilterInput{
			BankAccountIDs: &[]uuid.UUID{id},
		}

		if balanceStartDate.Valid {
			filter.StartDate = balanceStartDate
		}

		if balanceEndDate.Valid {
			filter.EndDate = balanceEndDate
		}

		if pageSize != nil {
			filter.PageSize = pageSize
		}

		balances, _, err := s.Repository.ResolveBalancesByFilter(filter.ToFilter())
		if err != nil {
			return nil, err
		}
		bankAccount.AttachBalances(balances, true)
	}

	return &bankAccount, nil
}

// GetByFilter fetches a set of Bank Accounts by its filter
func (s *BankAccountImpl) GetByFilter(input model.BankAccountFilterInput) ([]model.BankAccount, model.PageInfoOutput, error) {
	return s.Repository.ResolveByFilter(input.ToFilter())
}

// Update updates an existing Bank Account
func (s *BankAccountImpl) Update(input model.BankAccountInput, userID uuid.UUID) (*model.BankAccount, error) {
	bankAccounts, err := s.Repository.ResolveByIDs([]uuid.UUID{input.ID})
	if err != nil {
		return nil, err
	}

	if len(bankAccounts) != 1 {
		return nil, failure.EntityNotFound("Bank Account")
	}

	bankAccount := bankAccounts[0]

	err = bankAccount.Update(input, userID)
	if err != nil {
		return nil, err
	}

	err = s.Repository.Update(bankAccount)
	if err != nil {
		return nil, err
	}

	return &bankAccount, err
}

// Delete deletes an existing Bank Account. The method will find all the account's balances
// and delete all of them also.
func (s *BankAccountImpl) Delete(id uuid.UUID, userID uuid.UUID) (*model.BankAccount, error) {
	bankAccounts, err := s.Repository.ResolveByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(bankAccounts) != 1 {
		return nil, failure.EntityNotFound("Bank Account")
	}

	bankAccount := bankAccounts[0]

	// pre-validate to save one database call
	if !bankAccount.Deleted.Valid && !bankAccount.DeletedBy.Valid {
		filter := model.BankAccountBalanceFilterInput{}
		filter.BankAccountIDs = &[]uuid.UUID{bankAccount.ID}

		page := 1
		pageSize := math.MaxInt

		filter.Page = &page
		filter.PageSize = &pageSize

		balances, _, err := s.Repository.ResolveBalancesByFilter(filter.ToFilter())
		if err != nil {
			return nil, err
		}
		bankAccount.AttachBalances(balances, true)
	}

	err = bankAccount.Delete(userID)
	if err != nil {
		return nil, err
	}

	err = s.Repository.Update(bankAccount)
	if err != nil {
		return nil, err
	}

	return &bankAccount, err
}

// CreateBalance creates a new Bank Account Balance
func (s *BankAccountImpl) CreateBalance(input model.BankAccountBalanceInput, userID uuid.UUID) (*model.BankAccountBalance, error) {
	bankAccounts, err := s.Repository.ResolveByIDs([]uuid.UUID{input.BankAccountID})
	if err != nil {
		return nil, err
	}

	if len(bankAccounts) != 1 {
		return nil, failure.EntityNotFound("Bank Account")
	}

	bankAccount := bankAccounts[0]

	if bankAccount.Deleted.Valid {
		return nil, failure.OperationNotPermitted("add balance", "Bank Account", "the Bank Account is already deleted")
	}

	if bankAccount.Status == model.BankAccountStatusInactive {
		return nil, failure.OperationNotPermitted("add balance", "Bank Account", "the Bank Account is inactive")
	}

	lastBalances, err := s.Repository.ResolveLastBalancesByBankAccountID(bankAccount.ID, 1)
	if err != nil {
		return nil, err
	}

	if len(lastBalances) != 1 {
		return nil, failure.EntityNotFound("Bank Account Balance")
	}

	lastBalance := lastBalances[0]
	isNewerBalance := lastBalance.Date.Before(input.Date.Time())
	var bankAccountToUpdate *model.BankAccount

	if isNewerBalance {
		bankAccount.SetNewBalance(input, userID)
		bankAccountToUpdate = &bankAccount
	}

	bankAccountBalance := model.NewBankAccountBalanceFromInput(input, bankAccount.ID, userID)
	err = s.Repository.CreateBalance(bankAccountBalance, bankAccountToUpdate)
	if err != nil {
		return nil, err
	}

	return &bankAccountBalance, nil
}

// GetBalanceByID fetches a Bank Account Balance by its ID
func (s *BankAccountImpl) GetBalanceByID(id uuid.UUID) (*model.BankAccountBalance, error) {
	bankAccountBalances, err := s.Repository.ResolveBalancesByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(bankAccountBalances) != 1 {
		return nil, failure.EntityNotFound("Bank Account")
	}

	return &bankAccountBalances[0], nil
}

// GetBalancesByFilter fetches a set of Bank Account Balances by its filter
func (s *BankAccountImpl) GetBalancesByFilter(input model.BankAccountBalanceFilterInput) ([]model.BankAccountBalance, model.PageInfoOutput, error) {
	return s.Repository.ResolveBalancesByFilter(input.ToFilter())
}

// UpdateBalance updates an existing Bank Account Balance
func (s *BankAccountImpl) UpdateBalance(input model.BankAccountBalanceInput, userID uuid.UUID) (*model.BankAccountBalance, error) {
	bankAccounts, err := s.Repository.ResolveByIDs([]uuid.UUID{input.BankAccountID})
	if err != nil {
		return nil, err
	}

	if len(bankAccounts) != 1 {
		return nil, failure.EntityNotFound("Bank Account")
	}

	bankAccount := bankAccounts[0]

	if bankAccount.Deleted.Valid {
		return nil, failure.OperationNotPermitted("update", "Bank Account Balance", "the Bank Account is already deleted")
	}

	if bankAccount.Status == model.BankAccountStatusInactive {
		return nil, failure.OperationNotPermitted("update", "Bank Account Balance", "Bank Account is inactive")
	}

	bankAccountBalances, err := s.Repository.ResolveBalancesByIDs([]uuid.UUID{input.ID})
	if err != nil {
		return nil, err
	}

	if len(bankAccountBalances) != 1 {
		return nil, failure.EntityNotFound("Bank Account Balance")
	}

	bankAccountBalance := bankAccountBalances[0]

	if bankAccountBalance.Deleted.Valid {
		return nil, failure.OperationNotPermitted("update", "Bank Account Balance", "the Bank Account Balance is already deleted")
	}

	err = bankAccountBalance.Update(input, userID)
	if err != nil {
		return nil, err
	}

	lastBalances, err := s.Repository.ResolveLastBalancesByBankAccountID(bankAccount.ID, 1)
	if err != nil {
		return nil, err
	}

	if len(lastBalances) != 1 {
		return nil, failure.EntityNotFound("Bank Account Balance")
	}

	lastBalance := lastBalances[0]
	isNewerOrLastBalance := lastBalance.Date.Before(input.Date.Time()) || input.ID == lastBalance.ID
	var bankAccountToUpdate *model.BankAccount

	if isNewerOrLastBalance {
		bankAccount.SetNewBalance(input, userID)
		bankAccountToUpdate = &bankAccount
	}

	err = s.Repository.UpdateBalance(bankAccountBalance, bankAccountToUpdate)
	if err != nil {
		return nil, err
	}

	return &bankAccountBalance, nil
}

// DeleteBalance deletes an existing Bank Account Balance
func (s *BankAccountImpl) DeleteBalance(id uuid.UUID, userID uuid.UUID) (*model.BankAccountBalance, error) {
	bankAccountBalances, err := s.Repository.ResolveBalancesByIDs([]uuid.UUID{id})
	if err != nil {
		return nil, err
	}

	if len(bankAccountBalances) != 1 {
		return nil, failure.EntityNotFound("Bank Account Balance")
	}

	bankAccountBalance := bankAccountBalances[0]

	if bankAccountBalance.Deleted.Valid {
		return nil, failure.OperationNotPermitted("delete", "Bank Account Balance", "the Bank Account Balance is already deleted")
	}

	bankAccounts, err := s.Repository.ResolveByIDs([]uuid.UUID{bankAccountBalance.BankAccountID})
	if err != nil {
		return nil, err
	}

	if len(bankAccounts) != 1 {
		return nil, failure.EntityNotFound("Bank Account")
	}

	bankAccount := bankAccounts[0]

	if bankAccount.Deleted.Valid {
		return nil, failure.OperationNotPermitted("delete", "Bank Account Balance", "the Bank Account is already deleted")
	}

	if bankAccount.Status == model.BankAccountStatusInactive {
		return nil, failure.OperationNotPermitted("delete", "Bank Account Balance", "the Bank Account is inactive")
	}

	bankAccountBalance.Delete(userID)

	lastBalances, err := s.Repository.ResolveLastBalancesByBankAccountID(bankAccount.ID, 2)
	if err != nil {
		return nil, err
	}

	if len(lastBalances) < 1 {
		return nil, failure.EntityNotFound("Bank Account Last Balance")
	}

	if len(lastBalances) < 2 {
		return nil, failure.OperationNotPermitted("delete", "Bank Account Balance", "cannot delete the only Bank Account Balance belonging to a Bank Account")
	}

	lastBalanceDeleted := bankAccountBalance.ID.String() == lastBalances[0].ID.String()
	var bankAccountToUpdate *model.BankAccount

	if lastBalanceDeleted {
		newLastBalanceInput := model.BankAccountBalanceInput{
			ID:            lastBalances[1].ID,
			BankAccountID: lastBalances[1].BankAccountID,
			Balance:       lastBalances[1].Balance,
			Date:          cachetime.CacheTime(lastBalances[1].Date),
		}
		bankAccount.SetNewBalance(newLastBalanceInput, userID)
		bankAccountToUpdate = &bankAccount
	}

	err = s.Repository.UpdateBalance(bankAccountBalance, bankAccountToUpdate)
	if err != nil {
		return nil, err
	}

	return &bankAccountBalance, nil
}
