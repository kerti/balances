package model

import (
	"log"
	"time"

	"github.com/guregu/null"
	"github.com/kerti/balances/backend/util/cachetime"
	"github.com/kerti/balances/backend/util/nuuid"
	"github.com/satori/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user entity object
type User struct {
	ID        uuid.UUID   `db:"entity_id" validate:"min=36,max=36"`
	Username  string      `db:"username"`
	Email     string      `db:"email"`
	Password  string      `db:"password"`
	Name      string      `db:"name"`
	Created   time.Time   `db:"created"`
	CreatedBy uuid.UUID   `db:"created_by"`
	Updated   null.Time   `db:"updated"`
	UpdatedBy nuuid.NUUID `db:"updated_by"`
}

// NewUserFromInput creates a new user from its input struct
func NewUserFromInput(input UserInput, userID uuid.UUID) (u User) {
	now := time.Now()

	u = User{
		ID:        uuid.NewV4(),
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

// Update performs an update on a user
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

// SetPassword sets a user's password
func (u *User) SetPassword(password string, userID uuid.UUID) error {
	now := time.Now()

	u.Password = password
	u.Updated = null.TimeFrom(now)
	u.UpdatedBy = nuuid.From(userID)

	u.hashPassword()

	return nil
}

// ComparePassword compares a user's password in storage against an input
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

// ToOutput converts a user to its JSON-compatible object representation
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

// UserInput represents an input struct for user entity
type UserInput struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
}

// UserOutput is the JSON-compatible object representation of user
type UserOutput struct {
	ID        uuid.UUID            `json:"id"`
	Username  string               `json:"username"`
	Email     string               `json:"email"`
	Password  string               `json:"password"`
	Name      string               `json:"name"`
	Created   cachetime.CacheTime  `json:"created"`
	CreatedBy uuid.UUID            `json:"createdBy"`
	Updated   cachetime.NCacheTime `json:"updated"`
	UpdatedBy nuuid.NUUID          `json:"updatedBy"`
}
