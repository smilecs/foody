package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/smilecs/foody/repository"
)

func TestPostHandler_CreatePost(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewPostHandler(manager)
	userID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		formData       map[string]string
		fileContent    []byte
		userID         uuid.UUID
		expectedStatus int
	}{
		{
			name: "Valid post creation",
			formData: map[string]string{
				"title": "Test Post",
				"body":  "Test content",
				"tags":  "food,recipe",
			},
			fileContent:    []byte("fake image content"),
			userID:         userID,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			formData: map[string]string{
				"title": "Test Post",
				// Missing body
			},
			fileContent:    []byte("fake image content"),
			userID:         userID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing user ID",
			formData: map[string]string{
				"title": "Test Post",
				"body":  "Test content",
				"tags":  "food,recipe",
			},
			fileContent:    []byte("fake image content"),
			userID:         uuid.Nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create multipart request
			req := setupMultipartRequest(t, http.MethodPost, "/posts", tt.formData, "media", "post.jpg", tt.fileContent)
			if tt.userID != uuid.Nil {
				req = setupTestContext(req, tt.userID)
			}
			w := httptest.NewRecorder()

			// Execute request
			handler.CreatePost(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestPostHandler_GetPosts(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewPostHandler(manager)

	// Test cases
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
	}{
		{
			name:           "Get posts with default pagination",
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get posts with custom pagination",
			queryParams: map[string]string{
				"limit":  "5",
				"offset": "10",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get posts with invalid pagination",
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
			req := setupTestRequest(t, http.MethodGet, "/posts", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()

			// Execute request
			handler.GetPosts(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response structure
			var response struct {
				Posts      []repository.PostWithMedia `json:"posts"`
				Pagination struct {
					Total  int `json:"total"`
					Limit  int `json:"limit"`
					Offset int `json:"offset"`
				} `json:"pagination"`
			}
			readResponseBody(t, w, &response)
		})
	}
}

func TestPostHandler_GetPostByID(t *testing.T) {
	// Setup
	manager := mockRepositoryManager()
	handler := NewPostHandler(manager)
	postID := uuid.New()

	// Test cases
	tests := []struct {
		name           string
		postID         uuid.UUID
		expectedStatus int
	}{
		{
			name:           "Get existing post",
			postID:         postID,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get non-existent post",
			postID:         uuid.New(),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with post ID
			req := setupTestRequest(t, http.MethodGet, "/posts", nil)
			q := req.URL.Query()
			q.Add("id", tt.postID.String())
			req.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()

			// Execute request
			handler.GetPostByID(w, req)

			// Check response
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
