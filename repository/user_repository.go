package repository

import (
	"github.com/google/uuid"
	"github.com/smilecs/foody/db"
	"github.com/smilecs/foody/schema"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepository struct {
	Database db.Database
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*schema.User, error) {
	var user schema.User
	err := r.Database.QueryRowx("SELECT * FROM users WHERE user_id = $1", id).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user schema.User) (string, error) {
	query := `
		INSERT INTO users (user_id, name, email, media, date_of_birth)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var userID int
	err := r.Database.QueryRowx(query, user.Id, user.Name, user.Email, user.Media, user.DOB).Scan(&userID)
	if err != nil {
		log.Printf("error inserting user: %v", err)
		return "", err
	}

	log.Printf("User inserted with ID: %d", userID)
	return "", nil
}

func (r *UserRepository) AuthUser(email, password string) (bool, error){
var hashedPassword string

query := `SELECT password FROM users WHERE email = $1`
err := r.Database.QueryRowx(query, email).Scan(&hashedPassword)
if err != nil{
	return false, err
}

err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
if err != nil{
	return false, err
}
return true, err
}