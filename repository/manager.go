package repository

import "github.com/smilecs/foody/db"

type Manager struct{
UserRepo *UserRepository
PostRepo *PostRepository
}

func NewManager(database db.Database) *Manager{

}
