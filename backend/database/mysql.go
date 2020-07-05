package database

import (
	"fmt"

	// requred MySQL import
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/balances/backend/config"
	"github.com/kerti/balances/backend/util/logger"
)

// Result represent a query result object
type Result struct {
	Data  interface{}
	Error error
}

// Block contains a transaction block
type Block func(db *sqlx.Tx, c chan Result)

// MySQL is the MySQL database class
type MySQL struct {
	Config *config.Config
	db     *sqlx.DB
}

// Startup perform startup functions
func (m *MySQL) Startup() {
	logger.Trace("MySQL database driver starting up...")
	m.Config = config.Get()
	conf := m.Config
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", conf.DB.User, conf.DB.Pass, conf.DB.Host, conf.DB.Port, conf.DB.Name))
	info := fmt.Sprintf("%s@tcp(%s:%d)/%s", conf.DB.User, conf.DB.Host, conf.DB.Port, conf.DB.Name)
	if err != nil {
		logger.Warn("Failed to connect to %s mysql database at [%s]", conf.DB.Name, info)
	} else if err := db.Ping(); err != nil {
		logger.Err("Error while connecting to %s mysql database [%s]", conf.DB.Name, info)
	} else {
		logger.Info("Successfully connected to %s mysql [%s]", conf.DB.Name, info)
	}
	db.DB.SetMaxOpenConns(conf.DB.ConnLimit)
	m.db = db
}

// Shutdown cleans up everything and shuts down
func (m *MySQL) Shutdown() {
	logger.Trace("MySQL database driver shutting down...")
	if m.db != nil {
		m.db.Close()
	}
}

// WithTransaction performs queries with transaction
func (m *MySQL) WithTransaction(block Block) (result Result, err error) {
	c := make(chan Result)
	tx, err := m.db.Beginx()
	if err == nil {
		go block(tx, c)
	} else {
		logger.Err(err.Error())
		return
	}
	result = <-c
	if result.Error != nil {
		tx.Tx.Rollback()
	} else {
		tx.Tx.Commit()
	}
	return
}

// Get gets data
func (m *MySQL) Get(dest interface{}, query string, args ...interface{}) (err error) {
	return m.db.Get(dest, query, args...)
}

// Select selects records
func (m *MySQL) Select(dest interface{}, query string, args ...interface{}) (err error) {
	return m.db.Select(dest, query, args...)
}

// In performs queries with IN clause
func (m *MySQL) In(query string, params map[string]interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return query, args, err
	}
	return sqlx.In(query, args...)
}

// Prepare prepares an SQL statement
func (m *MySQL) Prepare(query string) (*sqlx.NamedStmt, error) {
	return m.db.PrepareNamed(query)
}

// PrepareBind prepares and binds an SQL statement
func (m *MySQL) PrepareBind(query string) (*sqlx.Stmt, error) {
	return m.db.Preparex(query)
}

// Rebind rebinds an SQL statement
func (m *MySQL) Rebind(query string) string {
	return m.db.Rebind(query)
}

// IsReady checks that the database is ready for operation
func (m *MySQL) IsReady() bool {
	if m.db == nil {
		return false
	}
	if err := m.db.Ping(); err != nil {
		return false
	}
	return true
}