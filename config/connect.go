package config

import (
	"log"
	"os"
	"sync"

	_ "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	DB        *sqlx.DB
	AWSSess   *session.Session
	S3_Bucket string
}

var (
	instance *Config
	once     sync.Once
)

func Init() *Config {
	once.Do(func() {
		dbURL := os.Getenv("DB_URL")
		conn, err := sqlx.Connect("postgres", dbURL)
		if err != nil {
			log.Fatalf("failed to connect to the database url: %v", err)
		}

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
		))
		instance = &Config{
			DB:        conn,
			AWSSess:   sess,
			S3_Bucket: os.Getenv("S3_BUCKET_NAME"),
		}
	})
	return instance
}

func Get() *Config {
	if instance == nil {
		log.Fatal("database instance is not initialized")
	}
	return instance
}
