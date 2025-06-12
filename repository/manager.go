package repository

import "github.com/smilecs/foody/config"

// Use interfaces from interfaces.go

type Manager struct {
	UserRepo     UserRepositoryInterface
	PostRepo     PostRepositoryInterface
	MediaRepo    MediaRepositoryInterface
	RecipeRepo   RecipeRepositoryInterface
	MealPlanRepo MealPlanRepositoryInterface
}

func NewManager(database config.Database) *Manager {
	return &Manager{
		UserRepo:     &UserRepository{Database: database},
		PostRepo:     &PostRepository{Database: database},
		MediaRepo:    &MediaRepository{Database: database},
		RecipeRepo:   &RecipeRepository{Database: database},
		MealPlanRepo: &MealPlanRepository{Database: database},
	}
}
