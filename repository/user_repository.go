package repository

import (
	"github.com/smilecs/foody/db"
	"github.com/smilecs/foody/schema"
	"log"
)

type UserRepository struct {
	Database db.Database
}

func (r *UserRepository) GetUserByID(id int) (*schema.User, error) {
	var user schema.User
	err := r.Database.QueryRowx("SELECT * FROM users WHERE id = $1", id).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user schema.User) (string, error) {
	query := `
		INSERT INTO users (name, email, media, date_of_birth)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var userID int
	err := r.Database.QueryRowx(query, user.Name, user.Email, user.Media, user.DOB).Scan(&userID)
	if err != nil {
		log.Printf("error inserting user: %v", err)
		return "", err
	}

	log.Printf("User inserted with ID: %d", userID)
	return "", nil
}

func (r *UserRepository) AuthUser(user schema.User) (string, error){

}
