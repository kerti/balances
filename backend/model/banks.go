package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/filter"
	"github.com/kerti/balances/backend/util/nuuid"
)

// BankAccountStatus indicates the status of a Bank Account
type BankAccountStatus string

const (
	// BankAccountStatusActive indicates an active Bank Account
	BankAccountStatusActive BankAccountStatus = "active"
	// BankAccountStatusInactive indicates an inactive Bank Account
	BankAccountStatusInactive BankAccountStatus = "inactive"
)

const (
	// BankAccountColumnID represents the corresponding column in Bank Account table
	BankAccountColumnID filter.Field = "bank_accounts.entity_id"
	// BankAccountColumnAccountName represents the corresponding column in Bank Account table
	BankAccountColumnAccountName filter.Field = "bank_accounts.account_name"
	// BankAccountColumnBankName represents the corresponding column in Bank Account table
	BankAccountColumnBankName filter.Field = "bank_accounts.bank_name"
	// BankAccountColumnAccountHolderName represents the corresponding column in Bank Account table
	BankAccountColumnAccountHolderName filter.Field = "bank_accounts.account_holder_name"
	// BankAccountColumnAccountNumber represents the corresponding column in Bank Account table
	BankAccountColumnAccountNumber filter.Field = "bank_accounts.account_number"
	// BankAccountColumnLastBalance represents the corresponding column in Bank Account table
	BankAccountColumnLastBalance filter.Field = "bank_accounts.last_balance"
	// BankAccountColumnLastBalanceDate represents the corresponding column in Bank Account table
	BankAccountColumnLastBalanceDate filter.Field = "bank_accounts.last_balance_date"
	// BankAccountColumnStatus represents the corresponding column in Bank Account table
	BankAccountColumnStatus filter.Field = "bank_accounts.status"
	// BankAccountColumnCreated represents the corresponding column in Bank Account table
	BankAccountColumnCreated filter.Field = "bank_accounts.created"
	// BankAccountColumnCreatedBy represents the corresponding column in Bank Account table
	BankAccountColumnCreatedBy filter.Field = "bank_accounts.created_by"
	// BankAccountColumnUpdated represents the corresponding column in Bank Account table
	BankAccountColumnUpdated filter.Field = "bank_accounts.updated"
	// BankAccountColumnUpdatedBy represents the corresponding column in Bank Account table
	BankAccountColumnUpdatedBy filter.Field = "bank_accounts.updated_by"
	// BankAccountColumnDeleted represents the corresponding column in Bank Account table
	BankAccountColumnDeleted filter.Field = "bank_accounts.deleted"
	// BankAccountColumnDeletedBy represents the corresponding column in Bank Account table
	BankAccountColumnDeletedBy filter.Field = "bank_accounts.deleted_by"
)

const (
	// BankAccountBalanceColumnID represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnID filter.Field = "bank_account_balances.entity_id"
	// BankAccountBalanceColumnBankAccountID represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnBankAccountID filter.Field = "bank_account_balances.bank_account_entity_id"
	// BankAccountBalanceColumnDate represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnDate filter.Field = "bank_account_balances.date"
	// BankAccountBalanceColumnBalance represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnBalance filter.Field = "bank_account_balances.balance"
	// BankAccountBalanceColumnCreated represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnCreated filter.Field = "bank_account_balances.created"
	// BankAccountBalanceColumnCreatedBy represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnCreatedBy filter.Field = "bank_account_balances.created_by"
	// BankAccountBalanceColumnUpdated represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnUpdated filter.Field = "bank_account_balances.updated"
	// BankAccountBalanceColumnUpdatedBy represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnUpdatedBy filter.Field = "bank_account_balances.updated_by"
	// BankAccountBalanceColumnDeleted represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnDeleted filter.Field = "bank_account_balances.deleted"
	// BankAccountBalanceColumnDeletedBy represents the corresponding column in Bank Account Balances table
	BankAccountBalanceColumnDeletedBy filter.Field = "bank_account_balances.deleted_by"
)

// BankAccount represents a Bank Account object
type BankAccount struct {
	ID                uuid.UUID            `db:"entity_id" validate:"min=36,max=36"`
	AccountName       string               `db:"account_name" validate:"max=255"`
	BankName          string               `db:"bank_name" validate:"max=255"`
	AccountHolderName string               `db:"account_holder_name" validate:"max=255"`
	AccountNumber     string               `db:"account_number" validate:"max=255"`
	LastBalance       float64              `db:"last_balance" validate:"min=0"`
	LastBalanceDate   time.Time            `db:"last_balance_date"`
	Status            BankAccountStatus    `db:"status"`
	Created           time.Time            `db:"created"`
	CreatedBy         uuid.UUID            `db:"created_by" validate:"min=36,max=36"`
	Updated           null.Time            `db:"updated"`
	UpdatedBy         nuuid.NUUID          `db:"updated_by" validate:"min=36,max=36"`
	Deleted           null.Time            `db:"deleted"`
	DeletedBy         nuuid.NUUID          `db:"deleted_by" validate:"min=36,max=36"`
	Balances          []BankAccountBalance `db:"-"`
}

// NewBankAccountFromInput creates a new Bank Account from its input object
func NewBankAccountFromInput(input BankAccountInput, userID uuid.UUID) (b BankAccount) {
	now := time.Now()
	newUUID, _ := uuid.NewV7()

	b = BankAccount{
		ID:                newUUID,
		AccountName:       input.AccountName,
		BankName:          input.BankName,
		AccountHolderName: input.AccountHolderName,
		AccountNumber:     input.AccountNumber,
		LastBalance:       input.LastBalance,
		LastBalanceDate:   input.LastBalanceDate.Time(),
		Status:            input.Status,
		Created:           now,
		CreatedBy:         userID,
	}

	balances := make([]BankAccountBalance, 0)
	for _, bbInput := range input.Balances {
		balances = append(balances, NewBankAccountBalanceFromInput(bbInput, b.ID, userID))
	}

	b.Balances = balances

	// TODO: Validate ?

	return
}

// AttachBalances attaches Bank Account Balances to a Bank Account
func (b *BankAccount) AttachBalances(balances []BankAccountBalance) BankAccount {
	for _, balance := range balances {
		if balance.BankAccountID == b.ID {
			b.Balances = append(b.Balances, balance)
		}
	}
	return *b
}

// Update performs an update on a Bank Account
func (b *BankAccount) Update(input BankAccountInput, userID uuid.UUID) error {
	now := time.Now()

	b.AccountName = input.AccountName
	b.BankName = input.BankName
	b.AccountHolderName = input.AccountHolderName
	b.AccountNumber = input.AccountNumber
	b.Status = input.Status
	b.Updated = null.TimeFrom(now)
	b.UpdatedBy = nuuid.From(userID)

	// TODO: Validate ?

	return nil
}

// Delete performs a delete on a Bank Account
func (b *BankAccount) Delete(userID uuid.UUID) error {
	now := time.Now()

	b.Deleted = null.TimeFrom(now)
	b.DeletedBy = nuuid.From(userID)

	deletedBalances := make([]BankAccountBalance, 0)
	for _, balance := range b.Balances {
		err := balance.Delete(userID)
		if err != nil {
			return err
		}

		deletedBalances = append(deletedBalances, balance)
	}

	b.Balances = deletedBalances

	// TODO: Validate ?

	return nil
}

// SetNewBalance sets a new balance and balance date on a Bank Account
func (b *BankAccount) SetNewBalance(input BankAccountBalanceInput, userID uuid.UUID) error {
	now := time.Now()

	b.LastBalance = input.Balance
	b.LastBalanceDate = input.Date.Time()
	b.Updated = null.TimeFrom(now)
	b.UpdatedBy = nuuid.From(userID)

	// TODO: Validate ?

	return nil
}

// ToOutput converts a Bank Account to its JSON-compatible object representation
func (b *BankAccount) ToOutput() BankAccountOutput {
	o := BankAccountOutput{
		ID:                b.ID,
		AccountName:       b.AccountName,
		BankName:          b.BankName,
		AccountHolderName: b.AccountHolderName,
		AccountNumber:     b.AccountNumber,
		LastBalance:       b.LastBalance,
		LastBalanceDate:   cachetime.CacheTime(b.LastBalanceDate),
		Status:            b.Status,
		Created:           cachetime.CacheTime(b.Created),
		CreatedBy:         b.CreatedBy,
		Updated:           cachetime.NCacheTime(b.Updated),
		UpdatedBy:         b.UpdatedBy,
		Deleted:           cachetime.NCacheTime(b.Deleted),
		DeletedBy:         b.DeletedBy,
	}

	bbOutput := make([]BankAccountBalanceOutput, 0)
	for _, bb := range b.Balances {
		bbOutput = append(bbOutput, bb.ToOutput())
	}

	o.Balances = bbOutput

	return o
}

// BankAccountInput represents an input struct for Bank Account entity
type BankAccountInput struct {
	ID                uuid.UUID                 `json:"id"`
	AccountName       string                    `json:"accountName"`
	BankName          string                    `json:"bankName"`
	AccountHolderName string                    `json:"accountHolderName"`
	AccountNumber     string                    `json:"accountNumber"`
	LastBalance       float64                   `json:"lastBalance"`
	LastBalanceDate   cachetime.CacheTime       `json:"lastBalanceDate"`
	Status            BankAccountStatus         `json:"status"`
	Balances          []BankAccountBalanceInput `json:"balances"`
}

// BankAccountOutput is the JSON-compatible object representation of Bank Account
type BankAccountOutput struct {
	ID                uuid.UUID                  `json:"id"`
	AccountName       string                     `json:"accountName"`
	BankName          string                     `json:"bankName"`
	AccountHolderName string                     `json:"accountHolderName"`
	AccountNumber     string                     `json:"accountNumber"`
	LastBalance       float64                    `json:"lastBalance"`
	LastBalanceDate   cachetime.CacheTime        `json:"lastBalanceDate"`
	Status            BankAccountStatus          `json:"status"`
	Created           cachetime.CacheTime        `json:"created"`
	CreatedBy         uuid.UUID                  `json:"createdBy"`
	Updated           cachetime.NCacheTime       `json:"updated,omitempty"`
	UpdatedBy         nuuid.NUUID                `json:"updatedBy,omitempty"`
	Deleted           cachetime.NCacheTime       `json:"deleted,omitempty"`
	DeletedBy         nuuid.NUUID                `json:"deletedBy,omitempty"`
	Balances          []BankAccountBalanceOutput `json:"balances"`
}

// BankAccountBalance represents a snapshot of a Bank Account's balance at a given time
type BankAccountBalance struct {
	ID            uuid.UUID   `db:"entity_id" validate:"min=36,max=36"`
	BankAccountID uuid.UUID   `db:"bank_account_entity_id" validate:"min=36,max=36"`
	Date          time.Time   `db:"date"`
	Balance       float64     `db:"balance"`
	Created       time.Time   `db:"created"`
	CreatedBy     uuid.UUID   `db:"created_by" validate:"min=36,max=36"`
	Updated       null.Time   `db:"updated"`
	UpdatedBy     nuuid.NUUID `db:"updated_by" validate:"min=36,max=36"`
	Deleted       null.Time   `db:"deleted"`
	DeletedBy     nuuid.NUUID `db:"deleted_by" validate:"min=36,max=36"`
}

// NewBankAccountBalanceFromInput creates a new Bank Account Balance from its input object
func NewBankAccountBalanceFromInput(input BankAccountBalanceInput, bankAccountID uuid.UUID, userID uuid.UUID) (bb BankAccountBalance) {
	now := time.Now()
	newUUID, _ := uuid.NewV7()

	bb = BankAccountBalance{
		ID:            newUUID,
		BankAccountID: bankAccountID,
		Date:          input.Date.Time(),
		Balance:       input.Balance,
		Created:       now,
		CreatedBy:     userID,
	}

	// TODO: Validate ?

	return
}

// Update performs an update on a Bank Account Balance
func (bb BankAccountBalance) Update(input BankAccountBalanceInput, userID uuid.UUID) error {
	now := time.Now()

	bb.Date = input.Date.Time()
	bb.Balance = input.Balance
	bb.Updated = null.TimeFrom(now)
	bb.UpdatedBy = nuuid.From(userID)

	// TODO: Validate ?

	return nil
}

// Delete performs a delete on a Bank Account Balance
func (bb BankAccountBalance) Delete(userID uuid.UUID) error {
	now := time.Now()

	bb.Deleted = null.TimeFrom(now)
	bb.DeletedBy = nuuid.From(userID)

	// TODO: Validate ?

	return nil
}

// ToOutput converts a Bank Account Balance to its JSON-compatible object representation
func (bb *BankAccountBalance) ToOutput() BankAccountBalanceOutput {
	return BankAccountBalanceOutput{
		ID:            bb.ID,
		BankAccountID: bb.BankAccountID,
		Date:          cachetime.CacheTime(bb.Date),
		Balance:       bb.Balance,
		Created:       cachetime.CacheTime(bb.Created),
		CreatedBy:     bb.CreatedBy,
		Updated:       cachetime.NCacheTime(bb.Updated),
		UpdatedBy:     bb.UpdatedBy,
		Deleted:       cachetime.NCacheTime(bb.Deleted),
		DeletedBy:     bb.DeletedBy,
	}
}

// BankAccountBalanceInput represents an input struct for Bank Account Balance entity
type BankAccountBalanceInput struct {
	ID            uuid.UUID           `json:"id"`
	BankAccountID uuid.UUID           `json:"bankAccountId"`
	Date          cachetime.CacheTime `json:"date"`
	Balance       float64             `json:"balance"`
}

// BankAccountBalanceOutput is the JSON-compatible object representation of Bank Account Balance
type BankAccountBalanceOutput struct {
	ID            uuid.UUID            `json:"id"`
	BankAccountID uuid.UUID            `json:"bankAccountId"`
	Date          cachetime.CacheTime  `json:"date"`
	Balance       float64              `json:"balance"`
	Created       cachetime.CacheTime  `json:"created"`
	CreatedBy     uuid.UUID            `json:"createdBy"`
	Updated       cachetime.NCacheTime `json:"updated,omitempty"`
	UpdatedBy     nuuid.NUUID          `json:"updatedBy,omitempty"`
	Deleted       cachetime.NCacheTime `json:"deleted,omitempty"`
	DeletedBy     nuuid.NUUID          `json:"deletedBy,omitempty"`
}

// BankAccountFilterInput is the filter input object for Bank Accounts
type BankAccountFilterInput struct {
	filter.BaseFilterInput
}

// ToFilter converts this entity-specific filter into a generic filter.Filter object
func (f *BankAccountFilterInput) ToFilter() filter.Filter {
	keywordFields := []filter.Field{
		BankAccountColumnAccountName,
		BankAccountColumnBankName,
		BankAccountColumnAccountNumber,
		BankAccountColumnAccountHolderName,
	}

	return filter.Filter{
		TableName:      "bank_accounts",
		Clause:         f.BaseFilterInput.GetKeywordFilter(keywordFields, false),
		IncludeDeleted: f.GetIncludeDeleted(),
		Pagination:     f.BaseFilterInput.GetPagination(),
	}
}

// BankAccountBalanceFilterInput is the filter input object for Bank Account Balances
type BankAccountBalanceFilterInput struct {
	filter.BaseFilterInput
	BankAccountID nuuid.NUUID          `json:"bankAccountId,omitempty"`
	StartDate     cachetime.NCacheTime `json:"startDate,omitempty"`
	EndDate       cachetime.NCacheTime `json:"endDate,omitempty"`
	BalanceMin    *float64             `json:"balanceMin,omitempty"`
	BalanceMax    *float64             `json:"balanceMax,omitempty"`
}

// ToFilter converts this entity-specific filter into a generic filter.Filter object
func (f *BankAccountBalanceFilterInput) ToFilter() filter.Filter {
	theFilter := filter.Filter{
		TableName:      "bank_account_balances",
		IncludeDeleted: f.GetIncludeDeleted(),
		Pagination:     f.BaseFilterInput.GetPagination(),
	}

	if f.BankAccountID.Valid {
		theFilter.AddClause(filter.Clause{
			Operand1: BankAccountBalanceColumnBankAccountID,
			Operand2: f.BankAccountID.UUID,
			Operator: filter.OperatorEqual,
		}, filter.OperatorAnd)
	}

	if f.StartDate.Valid {
		theFilter.AddClause(filter.Clause{
			Operand1: BankAccountBalanceColumnDate,
			Operand2: f.StartDate.Time,
			Operator: filter.OperatorGreaterThanEqual,
		}, filter.OperatorAnd)
	}

	if f.EndDate.Valid {
		theFilter.AddClause(filter.Clause{
			Operand1: BankAccountBalanceColumnDate,
			Operand2: f.EndDate.Time,
			Operator: filter.OperatorLessThanEqual,
		}, filter.OperatorAnd)
	}

	if f.BalanceMin != nil {
		theFilter.AddClause(filter.Clause{
			Operand1: BankAccountBalanceColumnBalance,
			Operand2: *f.BalanceMin,
			Operator: filter.OperatorGreaterThanEqual,
		}, filter.OperatorAnd)
	}

	if f.BalanceMax != nil {
		theFilter.AddClause(filter.Clause{
			Operand1: BankAccountBalanceColumnBalance,
			Operand2: *f.BalanceMax,
			Operator: filter.OperatorLessThanEqual,
		}, filter.OperatorAnd)
	}

	return theFilter
}
