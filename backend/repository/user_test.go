package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/kerti/balances/backend/model"
	"github.com/kerti/balances/backend/util/failure"
	"github.com/stretchr/testify/assert"
)

var (
	testInsertUserStatement = `INSERT INTO users
	( entity_id, username, email, password, name, created, created_by, updated, updated_by )
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )`

	testUpdateUserStatement = `
	UPDATE users
	SET username = ?, email = ?, password = ?, name = ?, created = ?, created_by = ?, updated = ?, updated_by = ?
	WHERE entity_id = ?`
)

var (
	testNow       = time.Now()
	testID1, _    = uuid.NewV4()
	testID2, _    = uuid.NewV4()
	testUserID, _ = uuid.NewV4()
	testUserModel = model.User{
		ID:        testID1,
		Username:  "username",
		Email:     "email@example.com",
		Password:  "password",
		Name:      "John Doe",
		Created:   testNow,
		CreatedBy: testUserID,
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
				WithArgs(testID1.String()).
				WillReturnRows(result)

			repo := new(UserMySQLRepo)
			repo.DB = &db
			_, err := repo.ExistsByID(testID1)
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
				WithArgs(testID1.String()).
				WillReturnError(errors.New(""))

			repo := new(UserMySQLRepo)
			repo.DB = &db
			_, err := repo.ExistsByID(testID1)
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

			repo := new(UserMySQLRepo)
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

			ids := []uuid.UUID{testID1}

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(testID1.String())

			mock.
				ExpectQuery(querySelectUser + " WHERE users.entity_id IN (?)").
				WithArgs(testID1).
				WillReturnRows(result)

			repo := new(UserMySQLRepo)
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

			ids := []uuid.UUID{testID1, testID2}

			result := sqlmock.
				NewRows([]string{"entity_id"}).
				AddRow(testID1.String()).
				AddRow(testID2.String())

			mock.
				ExpectQuery(querySelectUser+" WHERE users.entity_id IN (?, ?)").
				WithArgs(testID1, testID2).
				WillReturnRows(result)

			repo := new(UserMySQLRepo)
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

			ids := []uuid.UUID{testID1}

			mock.
				ExpectQuery(querySelectUser + " WHERE users.entity_id IN (?)").
				WithArgs(testID1).
				WillReturnError(errors.New(""))

			repo := new(UserMySQLRepo)
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
				ExpectQuery(querySelectUser+" WHERE users.username = ? OR users.email = ? LIMIT 1").
				WithArgs(testUserModel.Email, testUserModel.Email).
				WillReturnRows(result)

			repo := new(UserMySQLRepo)
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
				ExpectQuery(querySelectUser+" WHERE users.username = ? OR users.email = ? LIMIT 1").
				WithArgs(testUserModel.Email, testUserModel.Email).
				WillReturnError(errors.New(""))

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(testInsertUserStatement).
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

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnError(errors.New(""))

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnRows(checkExistenceResult)

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(testInsertUserStatement).
				WillReturnError(errors.New(""))

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(testInsertUserStatement).
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

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(testUpdateUserStatement).
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

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnError(errors.New(""))

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnRows(checkExistenceResult)

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(testUpdateUserStatement).
				WillReturnError(errors.New(""))

			repo := new(UserMySQLRepo)
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
				WithArgs(testID1.String()).
				WillReturnRows(checkExistenceResult)

			mock.
				ExpectPrepare(testUpdateUserStatement).
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

			repo := new(UserMySQLRepo)
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
