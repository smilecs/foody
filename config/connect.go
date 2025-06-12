package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	DB        Database
	AWSSess   *session.Session
	S3_Bucket string
	Port      string
}

var (
	instance *Config
	once     sync.Once
)

// SetTestInstance sets a mock instance for testing
func SetTestInstance(mockDB Database) {
	fmt.Println("[DEBUG] SetTestInstance called, setting singleton instance")
	instance = &Config{
		DB:        mockDB,
		AWSSess:   nil,
		S3_Bucket: "test-bucket",
		Port:      "8080",
	}
}

func Init() *Config {
	once.Do(func() {
		fmt.Println("[DEBUG] Init called, setting singleton instance")
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
			DB:        &SQLDatabase{DB: conn},
			AWSSess:   sess,
			S3_Bucket: os.Getenv("S3_BUCKET_NAME"),
			Port:      os.Getenv("PORT"),
		}
	})
	return instance
}

func Get() *Config {
	if instance == nil {
		fmt.Println("[DEBUG] Get called, but instance is nil!")
		log.Fatal("database instance is not initialized")
	}
	fmt.Println("[DEBUG] Get called, returning singleton instance")
	return instance
}

// Reset clears the singleton instance for testing
func Reset() {
	instance = nil
	once = sync.Once{}
}
