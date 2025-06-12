package handler

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/repository"
	"github.com/smilecs/foody/schema"
)

// TestConfigSingleton verifies that config.Get() works after TestMain sets the singleton
func TestConfigSingleton(t *testing.T) {
	t.Log("Starting TestConfigSingleton")

	// Reset config to ensure clean state
	config.Reset()

	// Create mock DB
	mockDB := &MockRepositoryManager{
		Users:     make(map[uuid.UUID]*schema.User),
		Posts:     make(map[uuid.UUID]*repository.PostWithMedia),
		Recipes:   make(map[uuid.UUID]*repository.RecipeWithMedia),
		MealPlans: make(map[uuid.UUID]*repository.MealPlanWithMedia),
		Media:     make(map[uuid.UUID]*schema.Media),
	}

	// Create a mock AWS session
	sess := session.Must(session.NewSession())

	// Set test instance
	config.SetTestInstance(mockDB)
	config.Get().AWSSess = sess
	config.Get().S3_Bucket = "test-bucket"

	// Try to get config
	cfg := config.Get()
	if cfg == nil {
		t.Fatal("config.Get() returned nil")
	}
	if cfg.DB == nil {
		t.Fatal("config.DB is nil")
	}
	t.Log("config.Get() returned a valid config instance")
}
