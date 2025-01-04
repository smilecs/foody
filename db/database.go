package db

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
    return s.DB.QueryRowx(query, args)
}

func (s *SQLDatabase) MustExec(query string, args ...interface{}) sql.Result {
    return s.DB.MustExec(query, args)
}
