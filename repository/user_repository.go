package repository

import (
	"log"

	"github.com/google/uuid"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/schema"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Database config.Database
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*schema.User, error) {
	var user schema.User
	err := r.Database.QueryRowx("SELECT * FROM users WHERE user_id = $1", id).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user schema.User, password string, mediaID uuid.UUID) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("error hashing password: %v", err)
		return "", err
	}

	query := `
		INSERT INTO users (user_id, name, username, email, media_id, date_of_birth, password)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	var userID int
	err = r.Database.QueryRowx(query, user.Id, user.Name, user.Username, user.Email, mediaID, user.DOB, hashedPassword).Scan(&userID)
	if err != nil {
		log.Printf("error inserting user: %v", err)
		return "", err
	}

	log.Printf("User inserted with ID: %d", userID)
	return "", nil
}

func (r *UserRepository) AuthUser(email, password string) (bool, error) {
	var hashedPassword string

	query := `SELECT password FROM users WHERE email = $1`
	err := r.Database.QueryRowx(query, email).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, err
}

func (r *UserRepository) GetUserByEmail(email string) (*schema.User, error) {
	var user schema.User
	err := r.Database.QueryRowx("SELECT * FROM users WHERE email = $1", email).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
