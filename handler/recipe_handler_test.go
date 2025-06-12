package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smilecs/foody/schema"
)

func TestRecipeHandler_CreateRecipe(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewRecipeHandler(manager)

	// Test cases
	tests := []struct {
		name           string
		recipe         schema.Recipe
		expectedStatus int
	}{
		{
			name: "Valid recipe creation",
			recipe: schema.Recipe{
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
				AuthorId:  uuid.New(),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			recipe: schema.Recipe{
				Id:    uuid.New(),
				Title: "Test Recipe",
				// Missing description and other required fields
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := setupTestRequest(t, http.MethodPost, "/recipes", tt.recipe)
			w := httptest.NewRecorder()

			// Execute request
			handler.CreateRecipe(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRecipeHandler_GetRecipes(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewRecipeHandler(manager)

	// Test cases
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
	}{
		{
			name:           "Get recipes with default pagination",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get recipes with custom pagination",
			queryParams: map[string]string{
				"limit":  "5",
				"offset": "10",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get recipes with invalid pagination",
			queryParams: map[string]string{
				"limit":  "invalid",
				"offset": "invalid",
			},
			expectedStatus: http.StatusOK, // Should still work with default values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with query parameters
			req := setupTestRequest(t, http.MethodGet, "/recipes", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()

			// Execute request
			handler.GetRecipes(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRecipeHandler_GetRecipeByID(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewRecipeHandler(manager)
	recipeID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		recipeID       uuid.UUID
		expectedStatus int
	}{
		{
			name:           "Get existing recipe",
			recipeID:       recipeID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get non-existent recipe",
			recipeID:       uuid.New(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with recipe ID
			req := setupTestRequest(t, http.MethodGet, "/recipes/"+tt.recipeID.String(), nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.GetRecipeByID(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRecipeHandler_GetRecipesByAuthorID(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewRecipeHandler(manager)
	authorID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		authorID       uuid.UUID
		expectedStatus int
	}{
		{
			name:           "Get recipes by existing author",
			authorID:       authorID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get recipes by non-existent author",
			authorID:       uuid.New(),
			expectedStatus: http.StatusOK, // Should return empty list, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with author ID
			req := setupTestRequest(t, http.MethodGet, "/recipes/author/"+tt.authorID.String(), nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.GetRecipesByAuthorID(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
