package handler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/smilecs/foody/repository"
	"github.com/smilecs/foody/schema"
)

// MockRepositoryManager implements repository.Manager for testing
type MockRepositoryManager struct {
	Users     map[uuid.UUID]*schema.User
	Posts     map[uuid.UUID]*repository.PostWithMedia
	Recipes   map[uuid.UUID]*repository.RecipeWithMedia
	MealPlans map[uuid.UUID]*repository.MealPlanWithMedia
	Media     map[uuid.UUID]*schema.Media
}

// NewMockRepositoryManager creates a new mock repository manager
func NewMockRepositoryManager() *repository.Manager {
	mock := &MockRepositoryManager{
		Users:     make(map[uuid.UUID]*schema.User),
		Posts:     make(map[uuid.UUID]*repository.PostWithMedia),
		Recipes:   make(map[uuid.UUID]*repository.RecipeWithMedia),
		MealPlans: make(map[uuid.UUID]*repository.MealPlanWithMedia),
		Media:     make(map[uuid.UUID]*schema.Media),
	}

	// Create mock repositories
	userRepo := &MockUserRepository{manager: mock}
	postRepo := &MockPostRepository{manager: mock}
	recipeRepo := &MockRecipeRepository{manager: mock}
	mealPlanRepo := &MockMealPlanRepository{manager: mock}
	mediaRepo := &MockMediaRepository{manager: mock}

	return &repository.Manager{
		UserRepo:     userRepo,
		PostRepo:     postRepo,
		RecipeRepo:   recipeRepo,
		MealPlanRepo: mealPlanRepo,
		MediaRepo:    mediaRepo,
	}
}

// MockUserRepository implements repository.UserRepository for testing
type MockUserRepository struct {
	manager   *MockRepositoryManager
	passwords map[uuid.UUID]string // Store passwords for authentication
}

func (r *MockUserRepository) GetUserByID(id uuid.UUID) (*schema.User, error) {
	if user, ok := r.manager.Users[id]; ok {
		return user, nil
	}
	return nil, sql.ErrNoRows
}

func (r *MockUserRepository) CreateUser(user schema.User, password string, mediaID uuid.UUID) (string, error) {
	// Store user
	r.manager.Users[user.Id] = &user

	// Store password for authentication
	if r.passwords == nil {
		r.passwords = make(map[uuid.UUID]string)
	}
	r.passwords[user.Id] = password

	// Store media
	media := &schema.Media{
		Id:        mediaID,
		URL:       "https://test-bucket.s3.amazonaws.com/users/" + user.Username + "/profile.jpg",
		MediaType: schema.Image,
		AuthorId:  user.Id,
	}
	r.manager.Media[mediaID] = media

	fmt.Printf("[DEBUG] Users map after CreateUser: %+v\n", r.manager.Users)
	return user.Id.String(), nil
}

func (r *MockUserRepository) AuthUser(email, password string) (bool, error) {
	// Find user by email
	var userID uuid.UUID
	for _, user := range r.manager.Users {
		if user.Email == email {
			userID = user.Id
			break
		}
	}

	if userID == uuid.Nil {
		return false, nil
	}

	// Check password
	if storedPassword, ok := r.passwords[userID]; ok {
		return storedPassword == password, nil
	}

	return false, nil
}

func (r *MockUserRepository) GetUserByEmail(email string) (*schema.User, error) {
	fmt.Printf("[DEBUG] MockUserRepository.GetUserByEmail called with email: %s\n", email)
	for _, user := range r.manager.Users {
		fmt.Printf("[DEBUG] Checking user: %+v\n", user)
		if user.Email == email {
			fmt.Printf("[DEBUG] Found user by email: %+v\n", user)
			return user, nil
		}
	}
	fmt.Printf("[DEBUG] No user found for email: %s\n", email)
	return nil, sql.ErrNoRows
}

// MockPostRepository implements repository.PostRepository for testing
type MockPostRepository struct {
	manager *MockRepositoryManager
}

func (r *MockPostRepository) CreatePost(post schema.Post, mediaID uuid.UUID, mediaURL string) error {
	r.manager.Posts[post.Id] = &repository.PostWithMedia{
		Post:      post,
		MediaURL:  mediaURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return nil
}

func (r *MockPostRepository) GetPostByUserID(id uuid.UUID) (*repository.PostWithMedia, error) {
	for _, post := range r.manager.Posts {
		if post.AuthorId == id {
			return post, nil
		}
	}
	return nil, nil
}

func (r *MockPostRepository) GetPosts(limit, offset int) ([]repository.PostWithMedia, error) {
	var posts []repository.PostWithMedia
	for _, post := range r.manager.Posts {
		posts = append(posts, *post)
	}
	return posts, nil
}

func (r *MockPostRepository) GetPostByID(id uuid.UUID) (*repository.PostWithMedia, error) {
	if post, ok := r.manager.Posts[id]; ok {
		return post, nil
	}
	return nil, nil
}

func (r *MockPostRepository) UpdatePost(post schema.Post) error {
	if existingPost, ok := r.manager.Posts[post.Id]; ok {
		existingPost.Post = post
		existingPost.UpdatedAt = time.Now()
	}
	return nil
}

func (r *MockPostRepository) DeletePost(id uuid.UUID) error {
	delete(r.manager.Posts, id)
	return nil
}

func (r *MockPostRepository) GetTotalPostsCount() (int, error) {
	return len(r.manager.Posts), nil
}

// MockRecipeRepository implements repository.RecipeRepository for testing
type MockRecipeRepository struct {
	manager *MockRepositoryManager
}

func (r *MockRecipeRepository) CreateRecipe(recipe schema.Recipe, mediaID uuid.UUID, mediaURL string) error {
	r.manager.Recipes[recipe.Id] = &repository.RecipeWithMedia{
		Recipe:   recipe,
		MediaURL: mediaURL,
	}
	return nil
}

func (r *MockRecipeRepository) GetRecipes(limit, offset int) ([]repository.RecipeWithMedia, error) {
	var recipes []repository.RecipeWithMedia
	for _, recipe := range r.manager.Recipes {
		recipes = append(recipes, *recipe)
	}
	return recipes, nil
}

func (r *MockRecipeRepository) GetRecipeByID(id uuid.UUID) (*repository.RecipeWithMedia, error) {
	if recipe, ok := r.manager.Recipes[id]; ok {
		return recipe, nil
	}
	return nil, nil
}

func (r *MockRecipeRepository) GetRecipesByAuthorID(authorID uuid.UUID) ([]repository.RecipeWithMedia, error) {
	var recipes []repository.RecipeWithMedia
	for _, recipe := range r.manager.Recipes {
		if recipe.AuthorId == authorID {
			recipes = append(recipes, *recipe)
		}
	}
	return recipes, nil
}

func (r *MockRecipeRepository) UpdateRecipe(recipe schema.Recipe) error {
	if existingRecipe, ok := r.manager.Recipes[recipe.Id]; ok {
		existingRecipe.Recipe = recipe
		existingRecipe.UpdatedAt = time.Now()
	}
	return nil
}

func (r *MockRecipeRepository) DeleteRecipe(id uuid.UUID) error {
	delete(r.manager.Recipes, id)
	return nil
}

func (r *MockRecipeRepository) GetTotalRecipesCount() (int, error) {
	return len(r.manager.Recipes), nil
}

// MockMealPlanRepository implements repository.MealPlanRepository for testing
type MockMealPlanRepository struct {
	manager *MockRepositoryManager
}

func (r *MockMealPlanRepository) CreateMealPlan(mealPlan schema.MealPlan) error {
	r.manager.MealPlans[mealPlan.Id] = &repository.MealPlanWithMedia{
		MealPlan: mealPlan,
	}
	return nil
}

func (r *MockMealPlanRepository) GetMealPlansByAuthorID(authorID uuid.UUID) ([]repository.MealPlanWithMedia, error) {
	var mealPlans []repository.MealPlanWithMedia
	for _, mealPlan := range r.manager.MealPlans {
		if mealPlan.AuthorId == authorID {
			mealPlans = append(mealPlans, *mealPlan)
		}
	}
	return mealPlans, nil
}

func (r *MockMealPlanRepository) GetMealPlanByID(id uuid.UUID) (*repository.MealPlanWithMedia, error) {
	if mealPlan, ok := r.manager.MealPlans[id]; ok {
		return mealPlan, nil
	}
	return nil, nil
}

func (r *MockMealPlanRepository) UpdateMealPlan(mealPlan schema.MealPlan) error {
	if existingMealPlan, ok := r.manager.MealPlans[mealPlan.Id]; ok {
		existingMealPlan.MealPlan = mealPlan
		existingMealPlan.UpdatedAt = time.Now()
	}
	return nil
}

func (r *MockMealPlanRepository) DeleteMealPlan(id uuid.UUID) error {
	delete(r.manager.MealPlans, id)
	return nil
}

// MockMediaRepository implements repository.MediaRepository for testing
type MockMediaRepository struct {
	manager *MockRepositoryManager
}

func (r *MockMediaRepository) CreateMedia(media schema.Media) (uuid.UUID, error) {
	r.manager.Media[media.Id] = &media
	return media.Id, nil
}

func (r *MockMediaRepository) GetMediaByID(id uuid.UUID) (*schema.Media, error) {
	if media, ok := r.manager.Media[id]; ok {
		return media, nil
	}
	return nil, nil
}

func (r *MockMediaRepository) GetMediaByAuthorID(authorID uuid.UUID) ([]schema.Media, error) {
	var mediaList []schema.Media
	for _, media := range r.manager.Media {
		if media.AuthorId == authorID {
			mediaList = append(mediaList, *media)
		}
	}
	return mediaList, nil
}

// Implement config.Database interface
func (m *MockRepositoryManager) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return &sqlx.Row{}
}

func (m *MockRepositoryManager) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return &sqlx.Rows{}, nil
}

func (m *MockRepositoryManager) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (m *MockRepositoryManager) Beginx() (*sqlx.Tx, error) {
	return &sqlx.Tx{}, nil
}

func (m *MockRepositoryManager) MustExec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
