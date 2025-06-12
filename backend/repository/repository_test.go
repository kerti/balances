package repository_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/balances/backend/database"
)

var (
	db   database.MySQL
	mock sqlmock.Sqlmock
)

func getMockedDriver(matcher sqlmock.QueryMatcher) (database.MySQL, sqlmock.Sqlmock) {
	mockDB, sqlmock, _ := sqlmock.New(sqlmock.QueryMatcherOption(matcher))
	mockSqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mock = sqlmock
	db.DB = mockSqlxDB
	return db, mock
}

func getExistsResult(exists bool) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"COUNT"}).AddRow(exists)
}

func getCountResult(count int) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"COUNT"}).AddRow(count)
}

func getSingleEntityIDResult(id uuid.UUID) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"entity_id"}).AddRow(id)
}

func getMultiEntityIDResult(ids []uuid.UUID) *sqlmock.Rows {
	result := sqlmock.NewRows([]string{"entity_id"})
	for _, id := range ids {
		result.AddRow(id)
	}
	return result
}
