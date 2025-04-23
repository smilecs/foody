package repository

import "github.com/smilecs/foody/config"

type Manager struct {
	UserRepo  *UserRepository
	PostRepo  *PostRepository
	MediaRepo *MediaRepository
}

func NewManager(database config.Database) *Manager {
	return &Manager{
		UserRepo:  &UserRepository{Database: database},
		PostRepo:  &PostRepository{Database: database},
		MediaRepo: &MediaRepository{Database: database},
	}
}
