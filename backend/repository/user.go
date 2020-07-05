package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kerti/balances/backend/database"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/logger"
	"github.com/satori/uuid"
)

const (
	querySelect = `
		SELECT
			users.entity_id,
			users.username,
			users.email,
			users.password,
			users.name,
			users.created,
			users.created_by,
			users.updated,
			users.updated_by
		FROM
			users `

	queryInsert = `
		INSERT INTO users (
			entity_id,
			username,
			email,
			password,
			name,
			created,
			created_by,
			updated,
			updated_by
		) VALUES (
			:entity_id,
			:username,
			:email,
			:password,
			:name,
			:created,
			:created_by,
			:updated,
			:updated_by
		)`

	queryUpdate = `
		UPDATE users
		SET
			username = :username,
			email = :email,
			password = :password,
			name = :name,
			created = :created,
			created_by = :created_by,
			updated = :updated,
			updated_by = :updated_by
		WHERE entity_id = :entity_id`
)

// User is the repository for users
type User struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup perform startup functions
func (r *User) Startup() {
	logger.Trace("User Repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *User) Shutdown() {
	logger.Trace("User Repository shutting down...")
}

// ExistsByID checks the existence of a user by its ID
func (r *User) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
	}
	return
}

// ResolveByIDs resolves users by their IDs
func (r *User) ResolveByIDs(ids []uuid.UUID) (users []model.User, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := sqlx.In(querySelect+" WHERE users.entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	err = r.DB.Select(&users, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	return
}

// ResolveByIdentity resolves a user by its username or email
func (r *User) ResolveByIdentity(identity string) (user model.User, err error) {
	err = r.DB.Get(
		&user,
		querySelect+" WHERE users.username = ? OR users.email = ? LIMIT 1",
		identity,
		identity,
	)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}
	return
}

// Create creates a user
func (r *User) Create(user model.User) error {
	exists, err := r.ExistsByID(user.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	var stmt *sqlx.NamedStmt

	if exists {
		err = errors.New("user already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	stmt, err = r.DB.Prepare(queryInsert)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}

// Update updates a user
func (r *User) Update(user model.User) error {
	exists, err := r.ExistsByID(user.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if !exists {
		err = errors.New("user does not exist")
		logger.ErrNoStack("%v", err)
		return err
	}

	stmt, err := r.DB.Prepare(queryUpdate)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}
