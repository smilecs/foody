package repository

import "github.com/smilecs/foody/config"

type Manager struct {
	UserRepo   *UserRepository
	PostRepo   *PostRepository
	MediaRepo  *MediaRepository
	RecipeRepo *RecipeRepository
}

func NewManager(database config.Database) *Manager {
	return &Manager{
		UserRepo:   &UserRepository{Database: database},
		PostRepo:   &PostRepository{Database: database},
		MediaRepo:  &MediaRepository{Database: database},
		RecipeRepo: &RecipeRepository{Database: database},
	}
}
