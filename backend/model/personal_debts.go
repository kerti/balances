package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/kerti/balances/backend/util/filter"
	"github.com/kerti/balances/backend/util/nuuid"
)

type PersonalDebtInterestType string
type PersonalDebtStatus string

const (
	// PersonalDebtInterestTypeNominal indicates a personal debt that uses nominal interest
	PersonalDebtInterestTypeNominal PersonalDebtInterestType = "nominal"
	// PersonalDebtInterestTypePercentage indicates a personal debt that uses percentage interest
	PersonalDebtInterestTypePercentage PersonalDebtInterestType = "percentage"
)

const (
	// PersonalDebtStatus indicates a personal debt that is active
	PersonalDebtStatusActive PersonalDebtStatus = "active"
	// PersonalDebtStatus indicates a personal debt that is paid
	PersonalDebtStatusPaid PersonalDebtStatus = "paid"
	// PersonalDebtStatus indicates a personal debt that is defaulted
	PersonalDebtStatusDefaulted PersonalDebtStatus = "defaulted"
	// PersonalDebtStatus indicates a personal debt that is written off
	PersonalDebtStatusWrittenOff PersonalDebtStatus = "written_off"
)

const (
	// PersonalDebtColumnID represents the corresponding column in the Personal Debts table
	PersonalDebtColumnID filter.Field = "personal_debts.entity_id"
	// PersonalDebtColumnName represents the corresponding column in the Personal Debts table
	PersonalDebtColumnName filter.Field = "personal_debts.name"
	// PersonalDebtColumnCreditor represents the corresponding column in the Personal Debts table
	PersonalDebtColumnCreditor filter.Field = "personal_debts.creditor"
	// PersonalDebtColumnContactInformation represents the corresponding column in the Personal Debts table
	PersonalDebtColumnContactInformation filter.Field = "personal_debts.contact_information"
	// PersonalDebtColumnPrincipal represents the corresponding column in the Personal Debts table
	PersonalDebtColumnPrincipal filter.Field = "personal_debts.principal"
	// PersonalDebtColumnInterest represents the corresponding column in the Personal Debts table
	PersonalDebtColumnInterest filter.Field = "personal_debts.interest"
	// PersonalDebtColumnInterestType represents the corresponding column in the Personal Debts table
	PersonalDebtColumnInterestType filter.Field = "personal_debts.interest_type"
	// PersonalDebtColumnDate represents the corresponding column in the Personal Debts table
	PersonalDebtColumnDate filter.Field = "personal_debts.date"
	// PersonalDebtColumnStatus represents the corresponding column in the Personal Debts table
	PersonalDebtColumnStatus filter.Field = "personal_debts.status"
	// PersonalDebtColumnCreated represents the corresponding column in the Personal Debts table
	PersonalDebtColumnCreated filter.Field = "personal_debts.created"
	// PersonalDebtColumnCreatedBy represents the corresponding column in the Personal Debts table
	PersonalDebtColumnCreatedBy filter.Field = "personal_debts.created_by"
	// PersonalDebtColumnUpdated represents the corresponding column in the Personal Debts table
	PersonalDebtColumnUpdated filter.Field = "personal_debts.updated"
	// PersonalDebtColumnUpdatedBy represents the corresponding column in the Personal Debts table
	PersonalDebtColumnUpdatedBy filter.Field = "personal_debts.updated_by"
	// PersonalDebtColumnDeleted represents the corresponding column in the Personal Debts table
	PersonalDebtColumnDeleted filter.Field = "personal_debts.deleted"
	// PersonalDebtColumnDeletedBy represents the corresponding column in the Personal Debts table
	PersonalDebtColumnDeletedBy filter.Field = "personal_debts."
)

const (
	// PersonalDebtBalanceColumnID represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnID filter.Field = "personal_debt_balances.entity_id"
	// PersonalDebtBalanceColumnPersonalDebtID represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnPersonalDebtID filter.Field = "personal_debt_balances.personal_debt_entity_id"
	// PersonalDebtBalanceColumnDate represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnDate filter.Field = "personal_debt_balances.date"
	// PersonalDebtBalanceColumnBalance represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnBalance filter.Field = "personal_debt_balances.balance"
	// PersonalDebtBalanceColumnCreated represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnCreated filter.Field = "personal_debt_balances.created"
	// PersonalDebtBalanceColumnCreatedBy represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnCreatedBy filter.Field = "personal_debt_balances.created_by"
	// PersonalDebtBalanceColumnUpdated represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnUpdated filter.Field = "personal_debt_balances.updated"
	// PersonalDebtBalanceColumnUpdatedBy represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnUpdatedBy filter.Field = "personal_debt_balances.updated_by"
	// PersonalDebtBalanceColumnDeleted represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnDeleted filter.Field = "personal_debt_balances.deleted"
	// PersonalDebtBalanceColumnDeletedBy represents the corresponding column in the Personal Debt Balances table
	PersonalDebtBalanceColumnDeletedBy filter.Field = "personal_debt_balances.deleted_by"
)

const (
	// PersonalDebtRepaymentColumnID represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnID filter.Field = "personal_debt_repayments.entity_id"
	// PersonalDebtRepaymentColumnPersonalDebtID represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnPersonalDebtID filter.Field = "personal_debt_repayments.personal_debt_entity_id"
	// PersonalDebtRepaymentColumnPaymentDate represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnPaymentDate filter.Field = "personal_debt_repayments.payment_date"
	// PersonalDebtRepaymentColumnPaymentAmount represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnPaymentAmount filter.Field = "personal_debt_repayments.payment_amount"
	// PersonalDebtRepaymentColumnCreated represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnCreated filter.Field = "personal_debt_repayments.created"
	// PersonalDebtRepaymentColumnCreatedBy represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnCreatedBy filter.Field = "personal_debt_repayments.created_by"
	// PersonalDebtRepaymentColumnUpdated represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnUpdated filter.Field = "personal_debt_repayments.updated"
	// PersonalDebtRepaymentColumnUpdatedBy represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnUpdatedBy filter.Field = "personal_debt_repayments.updated_by"
	// PersonalDebtRepaymentColumnDeleted represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnDeleted filter.Field = "personal_debt_repayments.deleted"
	// PersonalDebtRepaymentColumnDeletedBy represents the corresponding column in the Personal Debt Repayments table
	PersonalDebtRepaymentColumnDeletedBy filter.Field = "personal_debt_repayments.deleted_by"
)

// PersonalDebt represents a Personal Debt object
type PersonalDebt struct {
	ID                 uuid.UUID                `db:"entity_id" validate:"min=36,max=36"`
	Name               string                   `db:"name" validate:"max=255"`
	Creditor           string                   `db:"creditor" validate:"max=255"`
	ContactInformation string                   `db:"contact_information" validate:"max=255"`
	Principal          float64                  `db:"principal" validate:"min=0"`
	Interest           float64                  `db:"interest" validate:"min=0"`
	InterestType       PersonalDebtInterestType `db:"interest_type"`
	Date               time.Time                `db:"date"`
	Status             VehicleStatus            `db:"status"`
	Created            time.Time                `db:"created"`
	CreatedBy          uuid.UUID                `db:"created_by" validate:"min=36,max=36"`
	Updated            null.Time                `db:"updated"`
	UpdatedBy          nuuid.NUUID              `db:"updated_by" validate:"min=36,max=36"`
	Deleted            null.Time                `db:"deleted"`
	DeletedBy          nuuid.NUUID              `db:"deleted_by" validate:"min=36,max=36"`
	Balances           []PersonalDebtBalance
	Repayments         []PersonalDebtRepayment
}

// PersonalDebtBalance represents a Personal Debt Balance object
type PersonalDebtBalance struct {
	ID             uuid.UUID   `db:"entity_id" validate:"min=36,max=36"`
	PersonalDebtID uuid.UUID   `db:"personal_debt_entity_id" validate:"min=36,max=36"`
	Date           time.Time   `db:"date"`
	Balance        float64     `db:"balance" validate:"min=0"`
	Created        time.Time   `db:"created"`
	CreatedBy      uuid.UUID   `db:"created_by" validate:"min=36,max=36"`
	Updated        null.Time   `db:"updated"`
	UpdatedBy      nuuid.NUUID `db:"updated_by" validate:"min=36,max=36"`
	Deleted        null.Time   `db:"deleted"`
	DeletedBy      nuuid.NUUID `db:"deleted_by" validate:"min=36,max=36"`
}

// PersonalDebtRepayment represents a Personal Debt Repayment object
type PersonalDebtRepayment struct {
	ID             uuid.UUID   `db:"entity_id" validate:"min=36,max=36"`
	PersonalDebtID uuid.UUID   `db:"personal_debt_entity_id" validate:"min=36,max=36"`
	PaymentDate    time.Time   `db:"payment_date"`
	PaymentAmount  float64     `db:"payment_amount" validate:"min=0"`
	Created        time.Time   `db:"created"`
	CreatedBy      uuid.UUID   `db:"created_by" validate:"min=36,max=36"`
	Updated        null.Time   `db:"updated"`
	UpdatedBy      nuuid.NUUID `db:"updated_by" validate:"min=36,max=36"`
	Deleted        null.Time   `db:"deleted"`
	DeletedBy      nuuid.NUUID `db:"deleted_by" validate:"min=36,max=36"`
}
