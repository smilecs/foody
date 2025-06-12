package repository

import (
	"github.com/google/uuid"
	"github.com/smilecs/foody/schema"
)

// UserRepositoryInterface defines the methods for user repository
// (add more methods as needed)
type UserRepositoryInterface interface {
	GetUserByID(id uuid.UUID) (*schema.User, error)
	CreateUser(user schema.User, password string, mediaID uuid.UUID) (string, error)
	AuthUser(email, password string) (bool, error)
	GetUserByEmail(email string) (*schema.User, error)
}

type PostRepositoryInterface interface {
	CreatePost(post schema.Post, mediaID uuid.UUID, mediaURL string) error
	GetPostByID(id uuid.UUID) (*PostWithMedia, error)
	GetPosts(limit, offset int) ([]PostWithMedia, error)
	UpdatePost(post schema.Post) error
	DeletePost(id uuid.UUID) error
	GetTotalPostsCount() (int, error)
}

type RecipeRepositoryInterface interface {
	CreateRecipe(recipe schema.Recipe, mediaID uuid.UUID, mediaURL string) error
	GetRecipeByID(id uuid.UUID) (*RecipeWithMedia, error)
	GetRecipes(limit, offset int) ([]RecipeWithMedia, error)
	GetRecipesByAuthorID(authorID uuid.UUID) ([]RecipeWithMedia, error)
	UpdateRecipe(recipe schema.Recipe) error
	DeleteRecipe(id uuid.UUID) error
}

type MealPlanRepositoryInterface interface {
	CreateMealPlan(mealPlan schema.MealPlan) error
	GetMealPlanByID(id uuid.UUID) (*MealPlanWithMedia, error)
	GetMealPlansByAuthorID(authorID uuid.UUID) ([]MealPlanWithMedia, error)
	UpdateMealPlan(mealPlan schema.MealPlan) error
	DeleteMealPlan(id uuid.UUID) error
}

type MediaRepositoryInterface interface {
	CreateMedia(media schema.Media) (uuid.UUID, error)
	GetMediaByID(id uuid.UUID) (*schema.Media, error)
	GetMediaByAuthorID(authorID uuid.UUID) ([]schema.Media, error)
}
