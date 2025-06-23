package repository

import (
	"github.com/google/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/filter"
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

// User is the User repository interface
type User interface {
	Startup()
	Shutdown()
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByIDs(ids []uuid.UUID) (users []model.User, err error)
	ResolveByIdentity(identity string) (user model.User, err error)
	ResolveByFilter(filter filter.Filter) (users []model.User, pageInfo model.PageInfoOutput, err error)
	Create(user model.User) error
	Update(user model.User) error
}

// Vehicle is the Vehicle repository interface
type Vehicle interface {
	Startup()
	Shutdown()
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ExistsValueByID(id uuid.UUID) (exists bool, err error)
	ResolveByIDs(ids []uuid.UUID) (vehicles []model.Vehicle, err error)
	ResolveValuesByIDs(ids []uuid.UUID) (vehicleValues []model.VehicleValue, err error)
	ResolveByFilter(filter filter.Filter) (vehicles []model.Vehicle, pageInfo model.PageInfoOutput, err error)
	ResolveValuesByFilter(filter filter.Filter) (vehicleValues []model.VehicleValue, pageInfo model.PageInfoOutput, err error)
	ResolveLastValuesByVehicleID(id uuid.UUID, count int) (vehicleValues []model.VehicleValue, err error)
	Create(vehicle model.Vehicle) error
	Update(vehicle model.Vehicle) error
	CreateValue(vehicleValue model.VehicleValue, vehicle *model.Vehicle) error
	UpdateValue(vehicleValue model.VehicleValue, vehicle *model.Vehicle) error
}
