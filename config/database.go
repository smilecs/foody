package config

import (
    "database/sql"
    "github.com/jmoiron/sqlx"
)

type Database interface {
    QueryRowx(query string, args ...interface{}) *sqlx.Row
    Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
    MustExec(query string, args ...interface{}) (sql.Result, error)
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
