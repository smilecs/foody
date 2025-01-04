package repository

import "github.com/smilecs/foody/db"

type PostRepository struct {
	Database db.SQLDatabase
}