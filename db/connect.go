package db

import (
    "github.com/jmoiron/sqlx"
    "log"
    "os"
)

type Config struct {
    Conn *sqlx.DB
}

var config = Config{}
var Instance *SQLDatabase

func Init() {
    dbUrl := os.Getenv("DB_URL")
    conn, err := sqlx.Connect("postgres", dbUrl)
    if err != nil {
        log.Fatal(err)
    }
    config.Conn = conn
    Instance = &SQLDatabase{DB: conn}
}

func Get() *Config {
    return &config
}
