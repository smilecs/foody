package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/data"
	"github.com/smilecs/foody/repository"
	"github.com/smilecs/foody/schema"
)

func init() {
	data.UploadFileAndGetUrl = func(sess *session.Session, bucket, key string, file multipart.File, size int64, contentType string) (string, error) {
		return "https://mocked-url.com/profile.jpg", nil
	}
}

// TestMain initializes config before running tests
func TestMain(m *testing.M) {
	// Reset the config singleton for a clean test environment
	config.Reset()

	// Create a mock AWS session
	sess := session.Must(session.NewSession())

	// Create a mock database
	mockDB := &MockRepositoryManager{
		Users:     make(map[uuid.UUID]*schema.User),
		Posts:     make(map[uuid.UUID]*repository.PostWithMedia),
		Recipes:   make(map[uuid.UUID]*repository.RecipeWithMedia),
		MealPlans: make(map[uuid.UUID]*repository.MealPlanWithMedia),
		Media:     make(map[uuid.UUID]*schema.Media),
	}

	// Set the mock config with a dummy session and bucket
	config.SetTestInstance(mockDB)
	config.Get().AWSSess = sess
	config.Get().S3_Bucket = "test-bucket"

	os.Exit(m.Run())
}

// setupTestRequest creates a test request with the given method, path, and body
func setupTestRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

// setupMultipartRequest creates a test multipart form request
func setupMultipartRequest(t *testing.T, method, path string, formFields map[string]string, fileField, filename string, fileContent []byte) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range formFields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatalf("Failed to write form field: %v", err)
		}
	}

	// Add file only if fileContent is not nil
	if fileContent != nil {
		part, err := writer.CreateFormFile(fileField, filename)
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		if _, err := part.Write(fileContent); err != nil {
			t.Fatalf("Failed to write file content: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// setupTestContext adds a user ID to the request context
func setupTestContext(req *http.Request, userID uuid.UUID) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, "user_id", userID)
	return req.WithContext(ctx)
}

// readResponseBody reads and unmarshals the response body
func readResponseBody(t *testing.T, w *httptest.ResponseRecorder, v interface{}) {
	if err := json.NewDecoder(w.Body).Decode(v); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
}

// mockRepositoryManager returns a new mock repository manager
func mockRepositoryManager() *repository.Manager {
	// Reset config to ensure clean state
	config.Reset()

	// Create a mock database
	mockDB := &MockRepositoryManager{
		Users:     make(map[uuid.UUID]*schema.User),
		Posts:     make(map[uuid.UUID]*repository.PostWithMedia),
		Recipes:   make(map[uuid.UUID]*repository.RecipeWithMedia),
		MealPlans: make(map[uuid.UUID]*repository.MealPlanWithMedia),
		Media:     make(map[uuid.UUID]*schema.Media),
	}

	// Create a mock AWS session
	sess := session.Must(session.NewSession())

	// Set the mock config
	config.SetTestInstance(mockDB)
	config.Get().AWSSess = sess
	config.Get().S3_Bucket = "test-bucket"

	// Return the same manager instance that is set in the config
	return &repository.Manager{
		UserRepo:     &MockUserRepository{manager: mockDB},
		PostRepo:     &MockPostRepository{manager: mockDB},
		RecipeRepo:   &MockRecipeRepository{manager: mockDB},
		MealPlanRepo: &MockMealPlanRepository{manager: mockDB},
		MediaRepo:    &MockMediaRepository{manager: mockDB},
	}
}
