package model

import (
	"log"
	"time"

	"github.com/kerti/balances/backend/util/filter"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/nuuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// UserColumnID represents the corresponding column in User table
	UserColumnID filter.Field = "users.entity_id"
	// UserColumnUsername represents the corresponding column in User table
	UserColumnUsername filter.Field = "users.username"
	// UserColumnEmail represents the corresponding column in User table
	UserColumnEmail filter.Field = "users.email"
	// UserColumnPassword represents the corresponding column in User table
	UserColumnPassword filter.Field = "users.password"
	// UserColumnName represents the corresponding column in User table
	UserColumnName filter.Field = "users.name"
	// UserColumnCreated represents the corresponding column in User table
	UserColumnCreated filter.Field = "users.created"
	// UserColumnCreatedBy represents the corresponding column in User table
	UserColumnCreatedBy filter.Field = "users.created_by"
	// UserColumnUpdated represents the corresponding column in User table
	UserColumnUpdated filter.Field = "users.updated"
	// UserColumnUpdatedBy represents the corresponding column in User table
	UserColumnUpdatedBy filter.Field = "users.updated_by"
)

// User represents a User entity object
type User struct {
	ID        uuid.UUID   `db:"entity_id" validate:"min=36,max=36"`
	Username  string      `db:"username"`
	Email     string      `db:"email"`
	Password  string      `db:"password"`
	Name      string      `db:"name"`
	Created   time.Time   `db:"created"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"min=36,max=36"`
	Updated   null.Time   `db:"updated"`
	UpdatedBy nuuid.NUUID `db:"updated_by" validate:"min=36,max=36"`
}

// NewUserFromInput creates a new User from its input struct
func NewUserFromInput(input UserInput, userID uuid.UUID) (u User) {
	now := time.Now()
	newUUID, _ := uuid.NewV4()

	u = User{
		ID:        newUUID,
		Username:  input.Username,
		Email:     input.Email,
		Password:  input.Password,
		Name:      input.Name,
		Created:   now,
		CreatedBy: userID,
	}
	u.hashPassword()

	// TODO: Validate?

	return
}

// Update performs an update on a User
func (u *User) Update(input UserInput, userID uuid.UUID) error {
	now := time.Now()

	u.Username = input.Username
	u.Email = input.Email
	u.Name = input.Name
	u.Updated = null.TimeFrom(now)
	u.UpdatedBy = nuuid.From(userID)

	// TODO: Validate?

	return nil
}

// SetPassword sets a User's password
func (u *User) SetPassword(password string, userID uuid.UUID) error {
	now := time.Now()

	u.Password = password
	u.Updated = null.TimeFrom(now)
	u.UpdatedBy = nuuid.From(userID)

	u.hashPassword()

	// TODO: Validate?

	return nil
}

// ComparePassword compares a User's password in storage against an input
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) hashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}
	u.Password = string(hash)
	return nil
}

// ToOutput converts a User to its JSON-compatible object representation
func (u *User) ToOutput() UserOutput {
	return UserOutput{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  "********", // never expose password to public
		Name:      u.Name,
		Created:   cachetime.CacheTime(u.Created),
		CreatedBy: u.CreatedBy,
		Updated:   cachetime.NCacheTime(u.Updated),
		UpdatedBy: u.UpdatedBy,
	}
}

// UserInput represents an input struct for User entity
type UserInput struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
}

// UserOutput is the JSON-compatible object representation of User
type UserOutput struct {
	ID        uuid.UUID            `json:"id"`
	Username  string               `json:"username"`
	Email     string               `json:"email"`
	Password  string               `json:"password"`
	Name      string               `json:"name"`
	Created   cachetime.CacheTime  `json:"created"`
	CreatedBy uuid.UUID            `json:"createdBy"`
	Updated   cachetime.NCacheTime `json:"updated,omitempty"`
	UpdatedBy nuuid.NUUID          `json:"updatedBy,omitempty"`
}

// UserFilterInput is the filter input object for Users
type UserFilterInput struct {
	filter.BaseFilterInput
	IDs *[]uuid.UUID `json:"ids,omitempty"`
}

// ToFilter converts this entity-specific filter into a generic filter.Filter object
func (f *UserFilterInput) ToFilter() filter.Filter {
	theFilter := filter.Filter{
		TableName:      "users",
		IncludeDeleted: true,
		Pagination:     f.BaseFilterInput.GetPagination(),
	}

	if f.IDs != nil {
		for _, id := range *f.IDs {
			theFilter.AddClause(filter.Clause{
				Operand1: UserColumnID,
				Operand2: id,
				Operator: filter.OperatorEqual,
			}, filter.OperatorOr)
		}
	}

	keywordFields := []filter.Field{
		UserColumnUsername,
		UserColumnEmail,
		UserColumnName,
	}
	keywordClause := f.BaseFilterInput.GetKeywordFilter(keywordFields, false)
	theFilter.AddClause(*keywordClause, filter.OperatorAnd)

	return theFilter
}
