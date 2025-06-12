package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smilecs/foody/schema"
)

// TestFullUserJourney tests a complete user journey including:
// 1. User registration
// 2. User login
// 3. Creating a recipe
// 4. Creating a post with the recipe
// 5. Creating a meal plan
// 6. Retrieving and updating resources
func TestFullUserJourney(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	userHandler := NewUserHandler(manager)
	recipeHandler := NewRecipeHandler(manager)
	postHandler := NewPostHandler(manager)
	mealPlanHandler := NewMealPlanHandler(manager)

	// Step 1: Register a new user
	userID := uuid.New()
	formData := map[string]string{
		"username":      "testuser",
		"email":         "test@example.com",
		"date_of_birth": "1990-01-01",
		"password":      "password123",
	}
	req := setupMultipartRequest(t, http.MethodPost, "/users", formData, "media", "profile.jpg", []byte("fake image content"))
	w := httptest.NewRecorder()
	userHandler.CreateUser(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create user: expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Step 2: Login to get token
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	req = setupTestRequest(t, http.MethodPost, "/login", loginData)
	w = httptest.NewRecorder()
	userHandler.Login(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to login: expected status %d, got %d", http.StatusOK, w.Code)
	}

	var loginResponse LoginResponse
	readResponseBody(t, w, &loginResponse)
	if loginResponse.Token == "" {
		t.Fatal("Expected token in login response, got empty string")
	}

	// Step 3: Create a recipe
	recipe := schema.Recipe{
		Id:          uuid.New(),
		Title:       "Test Recipe",
		Description: "Test description",
		Ingredients: []schema.Ingredient{
			{
				Name:     "Ingredient 1",
				Quantity: 1.0,
				Unit:     "cup",
			},
		},
		Steps: []schema.Step{
			{
				Order:       1,
				Description: "Step 1",
			},
		},
		PrepTime:  func() *time.Duration { d := 30 * time.Minute; return &d }(),
		CookTime:  func() *time.Duration { d := 1 * time.Hour; return &d }(),
		TotalTime: func() *time.Duration { d := 90 * time.Minute; return &d }(),
		Servings:  &[]int{4}[0],
		AuthorId:  userID,
	}

	req = setupTestRequest(t, http.MethodPost, "/recipes", recipe)
	req = setupTestContext(req, userID)
	w = httptest.NewRecorder()
	recipeHandler.CreateRecipe(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create recipe: expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Step 4: Create a post with the recipe
	postFormData := map[string]string{
		"title": "Test Post",
		"body":  "Test content",
		"tags":  "food,recipe",
	}
	req = setupMultipartRequest(t, http.MethodPost, "/posts", postFormData, "media", "post.jpg", []byte("fake image content"))
	req = setupTestContext(req, userID)
	w = httptest.NewRecorder()
	postHandler.CreatePost(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create post: expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Step 5: Create a meal plan
	mealPlan := schema.MealPlan{
		Id:       uuid.New(),
		RecipeId: recipe.Id,
		AuthorId: userID,
		MealType: schema.Breakfast,
		Date:     time.Now(),
		Verified: false,
	}

	req = setupTestRequest(t, http.MethodPost, "/meal-plans", mealPlan)
	req = setupTestContext(req, userID)
	w = httptest.NewRecorder()
	mealPlanHandler.CreateMealPlan(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create meal plan: expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Step 6: Get all created resources
	// Get recipe
	req = setupTestRequest(t, http.MethodGet, "/recipes/"+recipe.Id.String(), nil)
	w = httptest.NewRecorder()
	recipeHandler.GetRecipeByID(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to get recipe: expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Get posts
	req = setupTestRequest(t, http.MethodGet, "/posts", nil)
	w = httptest.NewRecorder()
	postHandler.GetPosts(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to get posts: expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Get meal plans
	req = setupTestRequest(t, http.MethodGet, "/meal-plans/author/"+userID.String(), nil)
	w = httptest.NewRecorder()
	mealPlanHandler.GetMealPlansByAuthorID(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to get meal plans: expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Step 7: Update resources
	// Update recipe
	recipe.Title = "Updated Recipe Title"
	req = setupTestRequest(t, http.MethodPut, "/recipes/"+recipe.Id.String(), recipe)
	req = setupTestContext(req, userID)
	w = httptest.NewRecorder()
	recipeHandler.UpdateRecipe(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to update recipe: expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Update meal plan
	mealPlan.MealType = schema.Lunch
	req = setupTestRequest(t, http.MethodPut, "/meal-plans/"+mealPlan.Id.String(), mealPlan)
	req = setupTestContext(req, userID)
	w = httptest.NewRecorder()
	mealPlanHandler.UpdateMealPlan(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Failed to update meal plan: expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Step 8: Delete resources
	// Delete meal plan
	req = setupTestRequest(t, http.MethodDelete, "/meal-plans/"+mealPlan.Id.String(), nil)
	req = setupTestContext(req, userID)
	w = httptest.NewRecorder()
	mealPlanHandler.DeleteMealPlan(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("Failed to delete meal plan: expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Delete recipe
	req = setupTestRequest(t, http.MethodDelete, "/recipes/"+recipe.Id.String(), nil)
	req = setupTestContext(req, userID)
	w = httptest.NewRecorder()
	recipeHandler.DeleteRecipe(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("Failed to delete recipe: expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

// TestUnauthorizedAccess tests that unauthorized users cannot access protected resources
func TestUnauthorizedAccess(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	recipeHandler := NewRecipeHandler(manager)
	postHandler := NewPostHandler(manager)
	mealPlanHandler := NewMealPlanHandler(manager)

	// Test unauthorized recipe creation
	recipe := schema.Recipe{
		Id:    uuid.New(),
		Title: "Test Recipe",
	}
	req := setupTestRequest(t, http.MethodPost, "/recipes", recipe)
	w := httptest.NewRecorder()
	recipeHandler.CreateRecipe(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected unauthorized access for recipe creation, got status %d", w.Code)
	}

	// Test unauthorized post creation
	postFormData := map[string]string{
		"title": "Test Post",
		"body":  "Test content",
	}
	req = setupMultipartRequest(t, http.MethodPost, "/posts", postFormData, "media", "post.jpg", []byte("fake image content"))
	w = httptest.NewRecorder()
	postHandler.CreatePost(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected unauthorized access for post creation, got status %d", w.Code)
	}

	// Test unauthorized meal plan creation
	mealPlan := schema.MealPlan{
		Id:       uuid.New(),
		RecipeId: uuid.New(),
		MealType: schema.Breakfast,
		Date:     time.Now(),
	}
	req = setupTestRequest(t, http.MethodPost, "/meal-plans", mealPlan)
	w = httptest.NewRecorder()
	mealPlanHandler.CreateMealPlan(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected unauthorized access for meal plan creation, got status %d", w.Code)
	}
}

// TestResourceOwnership tests that users can only modify their own resources
func TestResourceOwnership(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	recipeHandler := NewRecipeHandler(manager)
	mealPlanHandler := NewMealPlanHandler(manager)

	// Create two users
	user1ID := uuid.New()
	user2ID := uuid.New()

	// Create a recipe as user1
	recipe := schema.Recipe{
		Id:          uuid.New(),
		Title:       "Test Recipe",
		Description: "Test description",
		AuthorId:    user1ID,
	}

	req := setupTestRequest(t, http.MethodPost, "/recipes", recipe)
	req = setupTestContext(req, user1ID)
	w := httptest.NewRecorder()
	recipeHandler.CreateRecipe(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create recipe: expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Try to update recipe as user2
	recipe.Title = "Updated by user2"
	req = setupTestRequest(t, http.MethodPut, "/recipes/"+recipe.Id.String(), recipe)
	req = setupTestContext(req, user2ID)
	w = httptest.NewRecorder()
	recipeHandler.UpdateRecipe(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected forbidden access for updating another user's recipe, got status %d", w.Code)
	}

	// Create a meal plan as user1
	mealPlan := schema.MealPlan{
		Id:       uuid.New(),
		RecipeId: recipe.Id,
		AuthorId: user1ID,
		MealType: schema.Breakfast,
		Date:     time.Now(),
	}

	req = setupTestRequest(t, http.MethodPost, "/meal-plans", mealPlan)
	req = setupTestContext(req, user1ID)
	w = httptest.NewRecorder()
	mealPlanHandler.CreateMealPlan(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create meal plan: expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Try to delete meal plan as user2
	req = setupTestRequest(t, http.MethodDelete, "/meal-plans/"+mealPlan.Id.String(), nil)
	req = setupTestContext(req, user2ID)
	w = httptest.NewRecorder()
	mealPlanHandler.DeleteMealPlan(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected forbidden access for deleting another user's meal plan, got status %d", w.Code)
	}
}
