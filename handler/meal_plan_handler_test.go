package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smilecs/foody/schema"
)

func TestMealPlanHandler_CreateMealPlan(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewMealPlanHandler(manager)
	userID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		mealPlan       schema.MealPlan
		userID         uuid.UUID
		expectedStatus int
	}{
		{
			name: "Valid meal plan creation",
			mealPlan: schema.MealPlan{
				Id:       uuid.New(),
				RecipeId: uuid.New(),
				AuthorId: userID,
				MealType: schema.Breakfast,
				Date:     time.Now(),
				Verified: false,
				PhotoId:  nil,
			},
			userID:         userID,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			mealPlan: schema.MealPlan{
				Id:   uuid.New(),
				Date: time.Now(),
				// Missing RecipeId and MealType
			},
			userID:         userID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing user ID",
			mealPlan: schema.MealPlan{
				Id:       uuid.New(),
				RecipeId: uuid.New(),
				MealType: schema.Breakfast,
				Date:     time.Now(),
				Verified: false,
			},
			userID:         uuid.Nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := setupTestRequest(t, http.MethodPost, "/meal-plans", tt.mealPlan)
			if tt.userID != uuid.Nil {
				req = setupTestContext(req, tt.userID)
			}
			w := httptest.NewRecorder()

			// Execute request
			handler.CreateMealPlan(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestMealPlanHandler_GetMealPlansByAuthorID(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewMealPlanHandler(manager)
	authorID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		authorID       uuid.UUID
		expectedStatus int
	}{
		{
			name:           "Get meal plans by existing author",
			authorID:       authorID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get meal plans by non-existent author",
			authorID:       uuid.New(),
			expectedStatus: http.StatusOK, // Should return empty list, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with author ID
			req := setupTestRequest(t, http.MethodGet, "/meal-plans/author/"+tt.authorID.String(), nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.GetMealPlansByAuthorID(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestMealPlanHandler_GetMealPlanByID(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewMealPlanHandler(manager)
	mealPlanID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		mealPlanID     uuid.UUID
		expectedStatus int
	}{
		{
			name:           "Get existing meal plan",
			mealPlanID:     mealPlanID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get non-existent meal plan",
			mealPlanID:     uuid.New(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with meal plan ID
			req := setupTestRequest(t, http.MethodGet, "/meal-plans/"+tt.mealPlanID.String(), nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.GetMealPlanByID(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestMealPlanHandler_UpdateMealPlan(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewMealPlanHandler(manager)
	userID := uuid.New()
	mealPlanID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		mealPlanID     uuid.UUID
		mealPlan       schema.MealPlan
		userID         uuid.UUID
		expectedStatus int
	}{
		{
			name:       "Update existing meal plan",
			mealPlanID: mealPlanID,
			mealPlan: schema.MealPlan{
				Id:       mealPlanID,
				RecipeId: uuid.New(),
				AuthorId: userID,
				MealType: schema.Lunch,
				Date:     time.Now(),
				Verified: true,
			},
			userID:         userID,
			expectedStatus: http.StatusOK,
		},
		{
			name:       "Update non-existent meal plan",
			mealPlanID: uuid.New(),
			mealPlan: schema.MealPlan{
				Id:       uuid.New(),
				RecipeId: uuid.New(),
				AuthorId: userID,
				MealType: schema.Dinner,
				Date:     time.Now(),
				Verified: false,
			},
			userID:         userID,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:       "Update meal plan without authorization",
			mealPlanID: mealPlanID,
			mealPlan: schema.MealPlan{
				Id:       mealPlanID,
				RecipeId: uuid.New(),
				AuthorId: uuid.New(), // Different author ID
				MealType: schema.Breakfast,
				Date:     time.Now(),
				Verified: false,
			},
			userID:         userID,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := setupTestRequest(t, http.MethodPut, "/meal-plans/"+tt.mealPlanID.String(), tt.mealPlan)
			if tt.userID != uuid.Nil {
				req = setupTestContext(req, tt.userID)
			}
			w := httptest.NewRecorder()

			// Execute request
			handler.UpdateMealPlan(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestMealPlanHandler_DeleteMealPlan(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewMealPlanHandler(manager)
	userID := uuid.New()
	mealPlanID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		mealPlanID     uuid.UUID
		userID         uuid.UUID
		expectedStatus int
	}{
		{
			name:           "Delete existing meal plan",
			mealPlanID:     mealPlanID,
			userID:         userID,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Delete non-existent meal plan",
			mealPlanID:     uuid.New(),
			userID:         userID,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Delete meal plan without authorization",
			mealPlanID:     mealPlanID,
			userID:         uuid.New(), // Different user ID
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := setupTestRequest(t, http.MethodDelete, "/meal-plans/"+tt.mealPlanID.String(), nil)
			if tt.userID != uuid.Nil {
				req = setupTestContext(req, tt.userID)
			}
			w := httptest.NewRecorder()

			// Execute request
			handler.DeleteMealPlan(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
