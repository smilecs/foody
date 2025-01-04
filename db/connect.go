package db

import (
    "github.com/jmoiron/sqlx"
    "log"
    "os"
)

type Config struct {
    Conn *sqlx.DB
}

var instance Database

func Init() Database {
    dbURL := os.Getenv("DB_URL")
    conn, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }
    instance = &SQLDatabase{DB: conn}
    return instance
}

func GetInstance() Database {
    if instance == nil {
        log.Fatal("database instance is not initialized")
    }
    return instance
}
