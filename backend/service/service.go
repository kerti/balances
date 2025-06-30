package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/cachetime"
)

// Auth is the service provider interface
type Auth interface {
	Startup()
	Shutdown()
	Authenticate(basic string) (authInfo *model.AuthenticationInfo, err error)
	Authorize(bearer string) (userID *uuid.UUID, err error)
	GetToken(user model.User) (token *string, expiration *time.Time, err error)
}

// BankAccount is the service provider interface
type BankAccount interface {
	Startup()
	Shutdown()
	Create(input model.BankAccountInput, userID uuid.UUID) (*model.BankAccount, error)
	GetByID(id uuid.UUID, withBalances bool, balanceStartDate, balanceEndDate cachetime.NCacheTime, pageSize *int) (*model.BankAccount, error)
	GetByFilter(input model.BankAccountFilterInput) ([]model.BankAccount, model.PageInfoOutput, error)
	Update(input model.BankAccountInput, userID uuid.UUID) (*model.BankAccount, error)
	Delete(id uuid.UUID, userID uuid.UUID) (*model.BankAccount, error)
	CreateBalance(input model.BankAccountBalanceInput, userID uuid.UUID) (*model.BankAccountBalance, error)
	GetBalanceByID(id uuid.UUID) (*model.BankAccountBalance, error)
	GetBalancesByFilter(input model.BankAccountBalanceFilterInput) ([]model.BankAccountBalance, model.PageInfoOutput, error)
	UpdateBalance(input model.BankAccountBalanceInput, userID uuid.UUID) (*model.BankAccountBalance, error)
	DeleteBalance(id uuid.UUID, userID uuid.UUID) (*model.BankAccountBalance, error)
}

// User is the service provider interface
type User interface {
	Startup()
	Shutdown()
	GetByID(id uuid.UUID) (*model.User, error)
	GetByFilter(input model.UserFilterInput) ([]model.User, model.PageInfoOutput, error)
	Create(input model.UserInput, userID uuid.UUID) (*model.User, error)
	Update(input model.UserInput, userID uuid.UUID) (*model.User, error)
}

// Vehicle is the service provider interface
type Vehicle interface {
	Startup()
	Shutdown()
	Create(input model.VehicleInput, userID uuid.UUID) (*model.Vehicle, error)
	GetByID(id uuid.UUID, withValues bool, valueStartDate, valueEndDate cachetime.NCacheTime, pageSize *int) (*model.Vehicle, error)
	GetByFilter(input model.VehicleFilterInput) ([]model.Vehicle, model.PageInfoOutput, error)
	Update(input model.VehicleInput, userID uuid.UUID) (*model.Vehicle, error)
	Delete(id uuid.UUID, userID uuid.UUID) (*model.Vehicle, error)
	CreateValue(input model.VehicleValueInput, userID uuid.UUID) (*model.VehicleValue, error)
	GetValueByID(id uuid.UUID) (*model.VehicleValue, error)
	GetValuesByFilter(input model.VehicleValueFilterInput) ([]model.VehicleValue, model.PageInfoOutput, error)
	UpdateValue(input model.VehicleValueInput, userID uuid.UUID) (*model.VehicleValue, error)
	DeleteValue(id uuid.UUID, userID uuid.UUID) (*model.VehicleValue, error)
}
