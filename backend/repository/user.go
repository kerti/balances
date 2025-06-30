package repository

import (
	"github.com/google/uuid"
	"github.com/kerti/balances/backend/database"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/kerti/balances/backend/util/filter"
	"github.com/kerti/balances/backend/util/logger"
)

const (
	QuerySelectUser = `
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

	QueryInsertUser = `
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

	QueryUpdateUser = `
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

// UserMySQLRepo is the repository for Users implemented with MySQL backend
type UserMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup perform startup functions
func (r *UserMySQLRepo) Startup() {
	logger.Trace("User Repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *UserMySQLRepo) Shutdown() {
	logger.Trace("User Repository shutting down...")
}

// ExistsByID checks the existence of a User by its ID
func (r *UserMySQLRepo) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
	}
	return
}

// ResolveByIDs resolves Users by their IDs
func (r *UserMySQLRepo) ResolveByIDs(ids []uuid.UUID) (users []model.User, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(QuerySelectUser+" WHERE users.entity_id IN (?)", ids)
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

// ResolveByIdentity resolves a User by its username or email
func (r *UserMySQLRepo) ResolveByIdentity(identity string) (user model.User, err error) {
	err = r.DB.Get(
		&user,
		QuerySelectUser+" WHERE users.username = ? OR users.email = ? LIMIT 1",
		identity,
		identity,
	)

	if err != nil {
		logger.Warn("[userRepo] unsuccessful user resolution using identity: %s", identity)
	}

	return
}

// ResolveByFilter resolves Users by a specified filter
func (r *UserMySQLRepo) ResolveByFilter(filter filter.Filter) (users []model.User, pageInfo model.PageInfoOutput, err error) {
	filterQueryString, err := filter.ToQueryString()
	if err != nil {
		return users, pageInfo, err
	}

	filterArgs := filter.GetArgs(true)
	query, args, err := r.DB.In(
		QuerySelectUser+filterQueryString+filter.Pagination.ToQueryString(),
		filterArgs...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("get by filter", "user", err)
		return
	}

	err = r.DB.Select(&users, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("get by filter", "user", err)
		return
	}

	var count int
	filterArgsNoPagination := filter.GetArgs(false)
	err = r.DB.Get(
		&count,
		"SELECT COUNT(entity_id) FROM users "+filterQueryString,
		filterArgsNoPagination...)
	if err != nil {
		logger.ErrNoStack("%v", err)
		err = failure.InternalError("get by filter", "user", err)
		return
	}

	pageInfo = model.PageInfoOutput{
		Page:       filter.Pagination.Page,
		PageSize:   filter.Pagination.PageSize,
		TotalCount: count,
		PageCount:  filter.Pagination.GetPageCount(count),
	}

	return
}

// Create creates a User
func (r *UserMySQLRepo) Create(user model.User) error {
	exists, err := r.ExistsByID(user.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if exists {
		err = failure.OperationNotPermitted("create", "User", "already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	stmt, err := r.DB.Prepare(QueryInsertUser)
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

// Update updates a User
func (r *UserMySQLRepo) Update(user model.User) error {
	exists, err := r.ExistsByID(user.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if !exists {
		err = failure.EntityNotFound("update", "User")
		logger.ErrNoStack("%v", err)
		return err
	}

	stmt, err := r.DB.Prepare(QueryUpdateUser)
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
