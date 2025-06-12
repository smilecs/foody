package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/schema"
)

// Add a package-level variable for userRepo access in checkResponse
var userRepoForTest *MockUserRepository

func TestUserHandler_CreateUser(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewUserHandler(manager)

	// Ensure the handler's Manager is the same as the one set in the config singleton
	if handler.Manager.UserRepo.(*MockUserRepository).manager != config.Get().DB.(*MockRepositoryManager) {
		t.Fatal("Handler is not using the correct mock repository manager from config singleton")
	}

	// Test cases
	tests := []struct {
		name           string
		formData       map[string]string
		fileContent    []byte
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Valid user creation",
			formData: map[string]string{
				"username":      "testuser",
				"email":         "test@example.com",
				"date_of_birth": "1990-01-01",
				"password":      "password123",
			},
			fileContent:    []byte("fake image content"),
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				// Verify user was created
				user, err := userRepoForTest.GetUserByEmail("test@example.com")
				if err != nil {
					t.Errorf("Failed to get created user: %v", err)
				}
				if user == nil {
					t.Error("User was not created")
				}
				if user.Username != "testuser" {
					t.Errorf("Expected username 'testuser', got '%s'", user.Username)
				}
			},
		},
		{
			name: "Missing required fields",
			formData: map[string]string{
				"username": "testuser",
				// Omit email, dob, and password
			},
			fileContent:    []byte("fake image content"),
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				// Verify no user was created
				user, _ := userRepoForTest.GetUserByEmail("test@example.com")
				if user != nil {
					t.Error("User was created despite missing required fields")
				}
			},
		},
		{
			name: "Invalid email format",
			formData: map[string]string{
				"username":      "testuser",
				"email":         "invalid-email",
				"date_of_birth": "1990-01-01",
				"password":      "password123",
			},
			fileContent:    []byte("fake image content"),
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				// Verify no user was created
				user, _ := userRepoForTest.GetUserByEmail("invalid-email")
				if user != nil {
					t.Error("User was created with invalid email")
				}
			},
		},
		{
			name: "Missing profile image",
			formData: map[string]string{
				"username":      "testuser",
				"email":         "test@example.com",
				"date_of_birth": "1990-01-01",
				"password":      "password123",
			},
			fileContent:    nil, // No file content
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				// Verify no user was created
				user, _ := userRepoForTest.GetUserByEmail("test@example.com")
				if user != nil {
					t.Error("User was created without profile image")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the user map before each test for isolation
			mgr := config.Get().DB.(*MockRepositoryManager)
			for k := range mgr.Users {
				delete(mgr.Users, k)
			}

			// Create multipart request
			req := setupMultipartRequest(t, http.MethodPost, "/users", tt.formData, "media", "profile.jpg", tt.fileContent)
			w := httptest.NewRecorder()

			// Execute request
			handler.CreateUser(w, req)

			// Get the actual repository manager from config singleton
			userRepo := &MockUserRepository{manager: mgr}

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Run additional checks if provided
			if tt.checkResponse != nil {
				userRepoForTest = userRepo
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewUserHandler(manager)

	// Create a test user first
	testUser := schema.User{
		Id:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		DOB:      "1990-01-01",
	}
	_, err := manager.UserRepo.CreateUser(testUser, "password123", uuid.New())
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test cases
	tests := []struct {
		name           string
		formData       map[string]string
		expectedStatus int
		expectToken    bool
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "Valid login",
			formData: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response LoginResponse
				readResponseBody(t, w, &response)
				if response.Token == "" {
					t.Error("Expected token in response, got empty string")
				}
			},
		},
		{
			name: "Invalid credentials",
			formData: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response LoginResponse
				readResponseBody(t, w, &response)
				if response.Token != "" {
					t.Error("Expected no token for invalid credentials")
				}
			},
		},
		{
			name: "Missing credentials",
			formData: map[string]string{
				"email": "test@example.com",
				// Missing password
			},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response LoginResponse
				readResponseBody(t, w, &response)
				if response.Token != "" {
					t.Error("Expected no token for missing credentials")
				}
			},
		},
		{
			name: "Non-existent user",
			formData: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response LoginResponse
				readResponseBody(t, w, &response)
				if response.Token != "" {
					t.Error("Expected no token for non-existent user")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create form request
			req := setupTestRequest(t, http.MethodPost, "/login", tt.formData)
			w := httptest.NewRecorder()

			// Execute request
			handler.Login(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Run additional checks if provided
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
