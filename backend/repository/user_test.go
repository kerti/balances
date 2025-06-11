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
	testUserModel     = model.User{
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

	t.Run("existsByID", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			result := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(result)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			_, err := repo.ExistsByID(userTestID1)
			repo.Shutdown()

			assert.Nil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("error", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			_, err := repo.ExistsByID(userTestID1)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

	})

	t.Run("resolveByIDs", func(t *testing.T) {

		t.Run("normalNoID", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			_, err := repo.ResolveByIDs([]uuid.UUID{})
			repo.Shutdown()

			assert.Nil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
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
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.Nil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
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
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.Nil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
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
			_, err := repo.ResolveByIDs(ids)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

	})

	t.Run("resolveByIdentity", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			result := sqlmock.NewRows([]string{"email"}).AddRow(testUserModel.Email)
			mock.
				ExpectQuery(repository.QuerySelectUser+" WHERE users.username = ? OR users.email = ? LIMIT 1").
				WithArgs(testUserModel.Email, testUserModel.Email).
				WillReturnRows(result)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			_, err := repo.ResolveByIdentity(testUserModel.Email)
			repo.Shutdown()

			assert.Nil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("error", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery(repository.QuerySelectUser+" WHERE users.username = ? OR users.email = ? LIMIT 1").
				WithArgs(testUserModel.Email, testUserModel.Email).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			_, err := repo.ResolveByIdentity(testUserModel.Email)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

	})

	t.Run("create", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(userStmtInsert).
				ExpectExec().
				WithArgs(
					testUserModel.ID,
					testUserModel.Username,
					testUserModel.Email,
					testUserModel.Password,
					testUserModel.Name,
					testUserModel.Created,
					testUserModel.CreatedBy,
					nil,
					nil).
				WillReturnResult(sqlmock.NewResult(1, 1))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Create(testUserModel)
			repo.Shutdown()

			assert.Nil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("errorOnCheckExistence", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Create(testUserModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("alreadyExists", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(checkExistenceResult)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Create(testUserModel)
			repo.Shutdown()

			assert.NotNil(t, err)
			assert.IsType(t, &failure.Failure{}, err)
			assert.Equal(t, failure.CodeOperationNotPermitted, err.(*failure.Failure).Code)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("failOnPrepare", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(userStmtInsert).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Create(testUserModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("failOnExec", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(userStmtInsert).
				ExpectExec().
				WithArgs(
					testUserModel.ID,
					testUserModel.Username,
					testUserModel.Email,
					testUserModel.Password,
					testUserModel.Name,
					testUserModel.Created,
					testUserModel.CreatedBy,
					nil,
					nil).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Create(testUserModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

	})

	t.Run("update", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(userStmtUpdate).
				ExpectExec().
				WithArgs(
					testUserModel.Username,
					testUserModel.Email,
					testUserModel.Password,
					testUserModel.Name,
					testUserModel.Created,
					testUserModel.CreatedBy,
					nil,
					nil,
					testUserModel.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Update(testUserModel)
			repo.Shutdown()

			assert.Nil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("errorOnCheckExistence", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Update(testUserModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("doesNotExist", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(false)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(checkExistenceResult)

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Update(testUserModel)
			repo.Shutdown()

			assert.NotNil(t, err)
			assert.IsType(t, &failure.Failure{}, err)
			assert.Equal(t, failure.CodeEntityNotFound, err.(*failure.Failure).Code)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("failOnPrepare", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(userStmtUpdate).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Update(testUserModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

		t.Run("failOnExec", func(t *testing.T) {
			db, mock := getMockedDriver(sqlmock.QueryMatcherEqual)

			checkExistenceResult := sqlmock.
				NewRows([]string{"COUNT"}).
				AddRow(true)

			mock.
				ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").
				WithArgs(userTestID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(userStmtUpdate).
				ExpectExec().
				WithArgs(
					testUserModel.Username,
					testUserModel.Email,
					testUserModel.Password,
					testUserModel.Name,
					testUserModel.Created,
					testUserModel.CreatedBy,
					nil,
					nil,
					testUserModel.ID).
				WillReturnError(errors.New(""))

			repo := new(repository.UserMySQLRepo)
			repo.DB = &db
			err := repo.Update(testUserModel)
			repo.Shutdown()

			assert.NotNil(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all mock expectations met")
			}
		})

	})

}
