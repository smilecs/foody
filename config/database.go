package config

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Database interface {
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	MustExec(query string, args ...interface{}) (sql.Result, error)
	Beginx() (*sqlx.Tx, error)
}

type SQLDatabase struct {
	DB *sqlx.DB
}

func (s *SQLDatabase) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return s.DB.QueryRowx(query, args...)
}

// Queryx executes a query that returns multiple rows.
func (s *SQLDatabase) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return s.DB.Queryx(query, args...)
}

// MustExec executes a query and returns an error if execution fails.
func (s *SQLDatabase) MustExec(query string, args ...interface{}) (sql.Result, error) {
	result, err := s.DB.Exec(query, args...)
	return result, err
}

func (s *SQLDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	return s.DB.Exec(query, args...)
}

func (s *SQLDatabase) Beginx() (*sqlx.Tx, error) {
	return s.DB.Beginx()
}
