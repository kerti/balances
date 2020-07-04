package pgsql

import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/satori/uuid"

	gorsk "github.com/kerti/balances/backend"
)

// User represents the client for user table
type User struct{}

// View returns single user by ID
func (u User) View(db orm.DB, id uuid.UUID) (gorsk.User, error) {
	user := gorsk.User{ID: id}
	err := db.Select(&user)
	return user, err
}

// Update updates user's info
func (u User) Update(db orm.DB, user gorsk.User) error {
	return db.Update(&user)
}
