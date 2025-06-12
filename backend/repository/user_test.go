package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/repository"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/stretchr/testify/assert"
)

var (
	userStmtInsert = `INSERT INTO users
	( entity_id, username, email, password, name, created, created_by, updated, updated_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )`

	userStmtUpdate = `
	UPDATE users
	SET username = ?, email = ?, password = ?, name = ?, created = ?, created_by = ?, updated = ?, updated_by = ?
	WHERE entity_id = ?`
)

var (
	userTestNow       = time.Now()
	userTestID1, _    = uuid.NewV7()
	userTestID2, _    = uuid.NewV7()
	userTestUserID, _ = uuid.NewV7()
	userTestModel     = model.User{
		ID:        userTestID1,
		Username:  "username",
		Email:     "email@example.com",
		Password:  "password",
		Name:      "John Doe",
		Created:   userTestNow,
		CreatedBy: userTestUserID,
	}
)

func TestUserRepository(t *testing.T) {

	t.Run("create", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(false))

			mock.
				ExpectPrepare(userStmtInsert).
				ExpectExec().
				WithArgs(
					userTestModel.ID,
					userTestModel.Username,
					userTestModel.Email,
					userTestModel.Password,
					userTestModel.Name,
					userTestModel.Created,
					userTestModel.CreatedBy,
					nil,
					nil).
				WillReturnResult(sqlmock.NewResult(1, 1))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(userTestModel)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorOnCheckExistence", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(userTestModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("alreadyExists", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(true))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(userTestModel)
			repo.Shutdown()

			assert.NotNil(t, err)
			assert.IsType(t, &failure.Failure{}, err)
			assert.Equal(t, failure.CodeOperationNotPermitted, err.(*failure.Failure).Code)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnPrepare", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(false))

			mock.
				ExpectPrepare(userStmtInsert).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(userTestModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnExec", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(false))

			mock.
				ExpectPrepare(userStmtInsert).
				ExpectExec().
				WithArgs(
					userTestModel.ID,
					userTestModel.Username,
					userTestModel.Email,
					userTestModel.Password,
					userTestModel.Name,
					userTestModel.Created,
					userTestModel.CreatedBy,
					nil,
					nil).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Create(userTestModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("existsByID", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(true))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ExistsByID(userTestID1)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("error", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ExistsByID(userTestID1)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("resolveByIDs", func(t *testing.T) {

		t.Run("normalNoID", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIDs([]uuid.UUID{})
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("normalSingleID", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{userTestID1}

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(userTestID1.String())

			mock.
				ExpectQuery(repository.QuerySelectUser + " WHERE users.entity_id IN (?)").
				WithArgs(userTestID1).
				WillReturnRows(result)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("normalMultipleID", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{userTestID1, userTestID2}

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(userTestID1.String()).
				AddRow(userTestID2.String())

			mock.
				ExpectQuery(repository.QuerySelectUser+" WHERE users.entity_id IN (?, ?)").
				WithArgs(userTestID1, userTestID2).
				WillReturnRows(result)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorExecutingSelect", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			ids := []uuid.UUID{userTestID1}

			mock.
				ExpectQuery(repository.QuerySelectUser + " WHERE users.entity_id IN (?)").
				WithArgs(userTestID1).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("resolveByIdentity", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			result := sqlmock.NewRows([]string{"email"}).AddRow(userTestModel.Email)
			mock.
				ExpectQuery(repository.QuerySelectUser+" WHERE users.username = ? OR users.email = ? LIMIT 1").
				WithArgs(userTestModel.Email, userTestModel.Email).
				WillReturnRows(result)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIdentity(userTestModel.Email)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("error", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery(repository.QuerySelectUser+" WHERE users.username = ? OR users.email = ? LIMIT 1").
				WithArgs(userTestModel.Email, userTestModel.Email).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			_, err := repo.ResolveByIdentity(userTestModel.Email)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("resolveByFilter", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			keyword := "example"
			likeKeyword := "%example%"

			dataResult := sqlmock.NewRows([]string{"entity_id"}).AddRow(userTestModel.ID)

			mock.
				ExpectQuery(repository.QuerySelectUser+" WHERE (((users.username LIKE ?) OR (users.email LIKE ?)) OR (users.name LIKE ?)) LIMIT ? OFFSET ?").
				WithArgs(likeKeyword, likeKeyword, likeKeyword, 10, 0).
				WillReturnRows(dataResult)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) FROM users WHERE (((users.username LIKE ?) OR (users.email LIKE ?)) OR (users.name LIKE ?))").
				WithArgs(likeKeyword, likeKeyword, likeKeyword).
				WillReturnRows(getCountResult(1))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			testFilter := model.UserFilterInput{}
			testFilter.Keyword = &keyword

			repo.Startup()
			_, _, err := repo.ResolveByFilter(testFilter.ToFilter())
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

	t.Run("update", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(true))

			mock.
				ExpectPrepare(userStmtUpdate).
				ExpectExec().
				WithArgs().
				WillReturnResult(sqlmock.NewResult(1, 1))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(userTestModel)
			repo.Shutdown()

			assert.Nil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("errorOnCheckExistence", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(userTestModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("doesNotExist", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(false))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(userTestModel)
			repo.Shutdown()

			assert.NotNil(t, err)
			assert.IsType(t, &failure.Failure{}, err)
			assert.Equal(t, failure.CodeEntityNotFound, err.(*failure.Failure).Code)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnPrepare", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(true))

			mock.
				ExpectPrepare(userStmtUpdate).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(userTestModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

		t.Run("failOnExec", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(getExistsResult(true))

			mock.
				ExpectPrepare(userStmtUpdate).
				ExpectExec().
				WithArgs().
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db

			repo.Startup()
			err := repo.Update(userTestModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			errMockExpectationsMet := mock.ExpectationsWereMet()

			assert.Nil(t, errMockExpectationsMet)
		})

	})

}
