package repository

import (
	"errors"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/balances/backend/database"
	"github.com/stretchr/testify/assert"
)

var (
	once sync.Once
	db   database.MySQL
	mock sqlmock.Sqlmock
)

func getMockedDriver() (database.MySQL, sqlmock.Sqlmock) {
	once.Do(func() {
		mockDB, sqlmock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mockSqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		mock = sqlmock

		db.Startup()
		db.Shutdown()
		db.DB = mockSqlxDB
	})
	return db, mock
}

func TestUserRepository(t *testing.T) {

	t.Run("existsByID", func(t *testing.T) {

		t.Run("normal", func(t *testing.T) {
			db, mock := getMockedDriver()

			id, _ := uuid.NewV4()

			result := sqlmock.NewRows([]string{"COUNT"}).AddRow(true)
			mock.ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").WithArgs(id.String()).WillReturnRows(result)

			repo := new(User)
			repo.DB = &db
			_, err := repo.ExistsByID(id)

			assert.Nil(t, err)
		})

		t.Run("error", func(t *testing.T) {
			db, mock := getMockedDriver()

			id, _ := uuid.NewV4()

			// result := sqlmock.NewRows([]string{"COUNT"}).AddRow(true)
			mock.ExpectQuery("SELECT COUNT(entity_id) > 0 FROM users WHERE users.entity_id = ?").WithArgs(id.String()).WillReturnError(errors.New(""))

			repo := new(User)
			repo.DB = &db
			_, err := repo.ExistsByID(id)
			assert.NotNil(t, err)
		})
	})

}
