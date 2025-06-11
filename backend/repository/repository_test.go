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
