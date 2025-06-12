package repository_test

import (
	"github.com/DATA-DOG/go-sqlmock"
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

func getExistsResult(result bool) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"COUNT"}).AddRow(result)
}

func getCountResult(count int) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"COUNT"}).AddRow(count)
}
