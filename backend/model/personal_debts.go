package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/failure"
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
	// PersonalDebtColumnCurrentBalance represents the corresponding column in the Personal Debts table
	PersonalDebtColumnCurrentBalance filter.Field = "personal_debts.current_balance"
	// PersonalDebtColumnCurrentBalanceDate represents the corresponding column in the Personal Debts table
	PersonalDebtColumnCurrentBalanceDate filter.Field = "personal_debts.current_balance_date"
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
	// PersonalDebtPaymentColumnID represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnID filter.Field = "personal_debt_payments.entity_id"
	// PersonalDebtPaymentColumnPersonalDebtID represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnPersonalDebtID filter.Field = "personal_debt_payments.personal_debt_entity_id"
	// PersonalDebtPaymentColumnPaymentDate represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnPaymentDate filter.Field = "personal_debt_payments.payment_date"
	// PersonalDebtPaymentColumnPaymentAmount represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnPaymentAmount filter.Field = "personal_debt_payments.payment_amount"
	// PersonalDebtPaymentColumnCreated represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnCreated filter.Field = "personal_debt_payments.created"
	// PersonalDebtPaymentColumnCreatedBy represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnCreatedBy filter.Field = "personal_debt_payments.created_by"
	// PersonalDebtPaymentColumnUpdated represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnUpdated filter.Field = "personal_debt_payments.updated"
	// PersonalDebtPaymentColumnUpdatedBy represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnUpdatedBy filter.Field = "personal_debt_payments.updated_by"
	// PersonalDebtPaymentColumnDeleted represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnDeleted filter.Field = "personal_debt_payments.deleted"
	// PersonalDebtPaymentColumnDeletedBy represents the corresponding column in the Personal Debt Payments table
	PersonalDebtPaymentColumnDeletedBy filter.Field = "personal_debt_payments.deleted_by"
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
	CurrentBalance     float64                  `db:"current_balance"`
	CurrentBalanceDate time.Time                `db:"current_balance_date"`
	Status             PersonalDebtStatus       `db:"status"`
	Created            time.Time                `db:"created"`
	CreatedBy          uuid.UUID                `db:"created_by" validate:"min=36,max=36"`
	Updated            null.Time                `db:"updated"`
	UpdatedBy          nuuid.NUUID              `db:"updated_by" validate:"min=36,max=36"`
	Deleted            null.Time                `db:"deleted"`
	DeletedBy          nuuid.NUUID              `db:"deleted_by" validate:"min=36,max=36"`
	Balances           []PersonalDebtBalance
	Payments           []PersonalDebtPayment
}

func NewPersonalDebtFromInput(input PersonalDebtInput, userID uuid.UUID) (pd PersonalDebt) {
	now := time.Now()
	newUUID, _ := uuid.NewV7()

	pd = PersonalDebt{
		ID:                 newUUID,
		Name:               input.Name,
		Creditor:           input.Creditor,
		ContactInformation: input.ContactInformation,
		Principal:          input.Principal,
		Interest:           input.Interest,
		InterestType:       input.InterestType,
		Date:               input.Date.Time(),
		CurrentBalance:     input.CurrentBalance,
		CurrentBalanceDate: input.CurrentBalanceDate.Time(),
		Status:             input.Status,
		Created:            now,
		CreatedBy:          userID,
	}

	balances := make([]PersonalDebtBalance, 0)

	// add principal + interest as firsst balance
	balances = append(balances, NewPersonalDebtBalanceFromInput(PersonalDebtBalanceInput{
		Date:    input.Date,
		Balance: input.Principal + input.Interest,
	}, pd.ID, userID))

	pd.Balances = balances

	// TODO: assume no payments initially?

	return
}

// AttachBalances attaches Personal Debt Balances to a Personal Debt
func (pd *PersonalDebt) AttachBalances(balances []PersonalDebtBalance, clearBeforeAttach bool) {
	if clearBeforeAttach {
		pd.Balances = []PersonalDebtBalance{}
	}

	for _, balance := range balances {
		if balance.PersonalDebtID == pd.ID {
			pd.Balances = append(pd.Balances, balance)
		}
	}
}

// AttachPayments attaches Personal Debt Payments to a Personal Debt
func (pd *PersonalDebt) AttachPayments(payments []PersonalDebtPayment, clearBeforeAttach bool) {
	if clearBeforeAttach {
		pd.Payments = []PersonalDebtPayment{}
	}

	for _, payment := range payments {
		if payment.PersonalDebtID == pd.ID {
			pd.Payments = append(pd.Payments, payment)
		}
	}
}

// Update performs an update on a Personal Debt
func (pd *PersonalDebt) Update(input PersonalDebtInput, userID uuid.UUID) error {
	if pd.Deleted.Valid || pd.DeletedBy.Valid {
		return failure.OperationNotPermitted("update", "Personal Debt", "already deleted")
	}

	now := time.Now()

	pd.Name = input.Name
	pd.Creditor = input.Creditor
	pd.ContactInformation = input.ContactInformation
	pd.Principal = input.Principal
	pd.Interest = input.Interest
	pd.InterestType = input.InterestType
	pd.Date = input.Date.Time()
	pd.Status = input.Status
	pd.Updated = null.TimeFrom(now)
	pd.UpdatedBy = nuuid.From(userID)

	return nil
}

// Delete performs a delete on a Personal Debt
func (pd *PersonalDebt) Delete(userID uuid.UUID) error {
	if pd.Deleted.Valid || pd.DeletedBy.Valid {
		return failure.OperationNotPermitted("delete", "Personal Debt", "already deleted")
	}

	now := time.Now()

	pd.Deleted = null.TimeFrom(now)
	pd.DeletedBy = nuuid.From(userID)

	deletedBalances := make([]PersonalDebtBalance, 0)
	for _, balance := range pd.Balances {
		err := balance.Delete(userID)
		if err != nil {
			return err
		}

		deletedBalances = append(deletedBalances, balance)
	}

	pd.Balances = deletedBalances

	deletedPayments := make([]PersonalDebtPayment, 0)
	for _, payment := range pd.Payments {
		err := payment.Delete(userID)
		if err != nil {
			return err
		}

		deletedPayments = append(deletedPayments, payment)
	}

	pd.Payments = deletedPayments

	return nil
}

// SetBalance sets a new current balance and current balance date on a Personal Debt
func (pd *PersonalDebt) SetBalance(input PersonalDebtBalanceInput, userID uuid.UUID) error {
	now := time.Now()

	pd.CurrentBalance = input.Balance
	pd.CurrentBalanceDate = input.Date.Time()
	pd.Updated = null.TimeFrom(now)
	pd.UpdatedBy = nuuid.From(userID)

	// TODO: Validate?

	return nil
}

// AddPayment adds a new payment and optionally sets the current balance date on a Personal Debt
func (pd *PersonalDebt) AddPayment(input PersonalDebtPaymentInput, userID uuid.UUID, alsoSetBalance bool) error {
	payment := NewPersonalDebtPaymentFromInput(input, pd.ID, userID)

	pd.Payments = append(pd.Payments, payment)

	if alsoSetBalance {
		balanceID, _ := uuid.NewV7()
		balanceInput := PersonalDebtBalanceInput{
			ID:             balanceID,
			PersonalDebtID: pd.ID,
			Date:           input.PaymentDate,
			Balance:        pd.CurrentBalance - input.PaymentAmount,
		}

		err := pd.SetBalance(balanceInput, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

// ToOutput converts a Personal Debt to its JSON-compatible object representation
func (pd *PersonalDebt) ToOutput() PersonalDebtOutput {
	o := PersonalDebtOutput{
		ID:                 pd.ID,
		Name:               pd.Name,
		Creditor:           pd.Creditor,
		ContactInformation: pd.ContactInformation,
		Principal:          pd.Principal,
		Interest:           pd.Interest,
		InterestType:       pd.InterestType,
		Date:               cachetime.CacheTime(pd.Date),
		CurrentBalance:     pd.CurrentBalance,
		CurrentBalanceDate: cachetime.CacheTime(pd.CurrentBalanceDate),
		Status:             pd.Status,
		Created:            cachetime.CacheTime(pd.Created),
		CreatedBy:          pd.CreatedBy,
		Updated:            cachetime.NCacheTime(pd.Updated),
		UpdatedBy:          pd.UpdatedBy,
		Deleted:            cachetime.NCacheTime(pd.Deleted),
		DeletedBy:          pd.DeletedBy,
	}

	pdbOutput := make([]PersonalDebtBalanceOutput, 0)
	for _, pdb := range pd.Balances {
		pdbOutput = append(pdbOutput, pdb.ToOutput())
	}

	o.Balances = pdbOutput

	pdpOutput := make([]PersonalDebtPaymentOutput, 0)
	for _, pdp := range pd.Payments {
		pdpOutput = append(pdpOutput, pdp.ToOutput())
	}

	o.Payments = pdpOutput

	return o
}

// PersonalDebtInput represents an input struct for Personal Debt entity
type PersonalDebtInput struct {
	ID                 uuid.UUID                `json:"ID"`
	Name               string                   `json:"name"`
	Creditor           string                   `json:"creditor"`
	ContactInformation string                   `json:"contactInformation"`
	Principal          float64                  `json:"principal"`
	Interest           float64                  `json:"interest"`
	InterestType       PersonalDebtInterestType `json:"interestType"`
	Date               cachetime.CacheTime      `json:"date"`
	CurrentBalance     float64                  `json:"currentBalance"`
	CurrentBalanceDate cachetime.CacheTime      `json:"currentBalanceDate"`
	Status             PersonalDebtStatus       `json:"status"`
}

// PersonalDebtOutput is the JSON-compatible representation of Personal Debt
type PersonalDebtOutput struct {
	ID                 uuid.UUID                   `json:"id"`
	Name               string                      `json:"name"`
	Creditor           string                      `json:"creditor"`
	ContactInformation string                      `json:"contactInformation"`
	Principal          float64                     `json:"principal"`
	Interest           float64                     `json:"interest"`
	InterestType       PersonalDebtInterestType    `json:"interestType"`
	Date               cachetime.CacheTime         `json:"date"`
	CurrentBalance     float64                     `json:"currentBalance"`
	CurrentBalanceDate cachetime.CacheTime         `json:"currentBalanceDate"`
	Status             PersonalDebtStatus          `json:"status"`
	Created            cachetime.CacheTime         `json:"created"`
	CreatedBy          uuid.UUID                   `json:"createdBy"`
	Updated            cachetime.NCacheTime        `json:"updated,omitempty"`
	UpdatedBy          nuuid.NUUID                 `json:"updatedBy,omitempty"`
	Deleted            cachetime.NCacheTime        `json:"deleted,omitempty"`
	DeletedBy          nuuid.NUUID                 `json:"deletedBy,omitempty"`
	Balances           []PersonalDebtBalanceOutput `json:"balances"`
	Payments           []PersonalDebtPaymentOutput `json:"payments"`
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

// NewPersonalDebtBalanceFromInput creates a new Personal Debt Balance from its input object
func NewPersonalDebtBalanceFromInput(input PersonalDebtBalanceInput, personalDebtID uuid.UUID, userID uuid.UUID) (pdb PersonalDebtBalance) {
	now := time.Now()
	newUUID, _ := uuid.NewV7()

	pdb = PersonalDebtBalance{
		ID:             newUUID,
		PersonalDebtID: personalDebtID,
		Date:           input.Date.Time(),
		Balance:        input.Balance,
		Created:        now,
		CreatedBy:      userID,
	}

	// TODO: Validate?

	return
}

// Update performs an update on a Personal Debt Balance
func (pdb *PersonalDebtBalance) Update(input PersonalDebtBalanceInput, userID uuid.UUID) error {
	if pdb.Deleted.Valid || pdb.DeletedBy.Valid {
		return failure.OperationNotPermitted("update", "Personal Debt Balance", "already deleted")
	}

	now := time.Now()

	pdb.Date = input.Date.Time()
	pdb.Balance = input.Balance
	pdb.Updated = null.TimeFrom(now)
	pdb.UpdatedBy = nuuid.From(userID)

	return nil
}

// Delete performs a delete on a Personal Debt Balance
func (pdb *PersonalDebtBalance) Delete(userID uuid.UUID) error {
	if pdb.Deleted.Valid || pdb.DeletedBy.Valid {
		return failure.OperationNotPermitted("delete", "Personal Debt Balance", "already deleted")
	}

	now := time.Now()

	pdb.Deleted = null.TimeFrom(now)
	pdb.DeletedBy = nuuid.From(userID)

	return nil
}

// ToOutput converts a Personal DebtBalance to its JSON-compatible object representation
func (pdb *PersonalDebtBalance) ToOutput() PersonalDebtBalanceOutput {
	return PersonalDebtBalanceOutput{
		ID:             pdb.ID,
		PersonalDebtID: pdb.PersonalDebtID,
		Date:           cachetime.CacheTime(pdb.Date),
		Balance:        pdb.Balance,
		Created:        cachetime.CacheTime(pdb.Created),
		CreatedBy:      pdb.CreatedBy,
		Updated:        cachetime.NCacheTime(pdb.Updated),
		UpdatedBy:      pdb.UpdatedBy,
		Deleted:        cachetime.NCacheTime(pdb.Deleted),
		DeletedBy:      pdb.DeletedBy,
	}
}

// PersonalDebtBalanceInput represents an input struct for a Personal Debt Balance entity
type PersonalDebtBalanceInput struct {
	ID             uuid.UUID           `json:"id"`
	PersonalDebtID uuid.UUID           `json:"personalDebtId"`
	Date           cachetime.CacheTime `json:"date"`
	Balance        float64             `json:"balance"`
}

// PersonalDebtBalanceOutput is the JSON-compatible object representation of Personal Debt Balance
type PersonalDebtBalanceOutput struct {
	ID             uuid.UUID            `json:"id"`
	PersonalDebtID uuid.UUID            `json:"personalDebtId"`
	Date           cachetime.CacheTime  `json:"date"`
	Balance        float64              `json:"balance"`
	Created        cachetime.CacheTime  `json:"created"`
	CreatedBy      uuid.UUID            `json:"createdBy"`
	Updated        cachetime.NCacheTime `json:"updated,omitempty"`
	UpdatedBy      nuuid.NUUID          `json:"updatedBy,omitempty"`
	Deleted        cachetime.NCacheTime `json:"deleted,omitempty"`
	DeletedBy      nuuid.NUUID          `json:"deletedBy,omitempty"`
}

// PersonalDebtPayment represents a Personal Debt Payment object
type PersonalDebtPayment struct {
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

// NewPersonalDebtPaymentFromInput creates a new Personal Debt Payment from its input object
func NewPersonalDebtPaymentFromInput(input PersonalDebtPaymentInput, personalDebtID uuid.UUID, userID uuid.UUID) (pdp PersonalDebtPayment) {
	now := time.Now()
	newUUID, _ := uuid.NewV7()

	pdp = PersonalDebtPayment{
		ID:             newUUID,
		PersonalDebtID: personalDebtID,
		PaymentDate:    input.PaymentDate.Time(),
		PaymentAmount:  input.PaymentAmount,
		Created:        now,
		CreatedBy:      userID,
	}

	// TODO: validate?

	return
}

// Update performs an update on a Personal Debt Payment
func (pdp *PersonalDebtPayment) Update(input PersonalDebtPaymentInput, userID uuid.UUID) error {
	if pdp.Deleted.Valid || pdp.DeletedBy.Valid {
		return failure.OperationNotPermitted("update", "Personal Debt Payment", "already deleted")
	}

	now := time.Now()

	pdp.PaymentDate = input.PaymentDate.Time()
	pdp.PaymentAmount = input.PaymentAmount
	pdp.Updated = null.TimeFrom(now)
	pdp.UpdatedBy = nuuid.From(userID)

	return nil
}

// Delete performs a delete on a Personal Debt Payment
func (pdp *PersonalDebtPayment) Delete(userID uuid.UUID) error {
	if pdp.Deleted.Valid || pdp.DeletedBy.Valid {
		return failure.OperationNotPermitted("delete", "Personal Debt Payment", "already deleted")
	}

	now := time.Now()

	pdp.Deleted = null.TimeFrom(now)
	pdp.DeletedBy = nuuid.From(userID)

	return nil
}

// ToOutput converts a Personal Debt Payment to its JSON-compatible object representation
func (pdp *PersonalDebtPayment) ToOutput() PersonalDebtPaymentOutput {
	return PersonalDebtPaymentOutput{
		ID:             pdp.ID,
		PersonalDebtID: pdp.PersonalDebtID,
		PaymentDate:    cachetime.CacheTime(pdp.PaymentDate),
		PaymentAmount:  pdp.PaymentAmount,
		Created:        cachetime.CacheTime(pdp.Created),
		CreatedBy:      pdp.CreatedBy,
		Updated:        cachetime.NCacheTime(pdp.Updated),
		UpdatedBy:      pdp.UpdatedBy,
		Deleted:        cachetime.NCacheTime(pdp.Deleted),
		DeletedBy:      pdp.DeletedBy,
	}
}

// PersonalDebtPaymentInput represents an input struct for Personal Debt Payment entity
type PersonalDebtPaymentInput struct {
	ID             uuid.UUID           `json:"id"`
	PersonalDebtID uuid.UUID           `json:"personalDebtId"`
	PaymentDate    cachetime.CacheTime `json:"paymentDate"`
	PaymentAmount  float64             `json:"paymentAmount"`
}

// PersonalDebtPaymentOutput is the JSON-compatible object representation of Personal Debt Payment
type PersonalDebtPaymentOutput struct {
	ID             uuid.UUID            `json:"id"`
	PersonalDebtID uuid.UUID            `json:"personalDebtId"`
	PaymentDate    cachetime.CacheTime  `json:"paymentDate"`
	PaymentAmount  float64              `json:"paymentAmount"`
	Created        cachetime.CacheTime  `json:"created"`
	CreatedBy      uuid.UUID            `json:"createdBy"`
	Updated        cachetime.NCacheTime `json:"updated,omitempty"`
	UpdatedBy      nuuid.NUUID          `json:"updatedBy,omitempty"`
	Deleted        cachetime.NCacheTime `json:"deleted,omitempty"`
	DeletedBy      nuuid.NUUID          `json:"deletedBy,omitempty"`
}
