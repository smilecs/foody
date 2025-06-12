package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/handler"
	"github.com/smilecs/foody/middleware"
	"github.com/smilecs/foody/repository"
)

func main() {
	cfg := config.Init()

	// Initialize manager with all repositories
	manager := repository.NewManager(cfg.DB)

	// Initialize handlers with manager
	userHandler := handler.NewUserHandler(manager)
	postHandler := handler.NewPostHandler(manager)
	recipeHandler := handler.NewRecipeHandler(manager)
	mealPlanHandler := handler.NewMealPlanHandler(manager)

	router := chi.NewRouter()

	// Public routes
	router.Post("/signup", userHandler.CreateUser)
	router.Post("/login", userHandler.Login)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Post routes
		r.Post("/posts", postHandler.CreatePost)
		r.Get("/posts", postHandler.GetPosts)
		r.Get("/posts/{id}", postHandler.GetPostByID)
		r.Put("/posts/{id}", postHandler.UpdatePost)
		r.Delete("/posts/{id}", postHandler.DeletePost)

		// Recipe routes
		r.Route("/api/recipes", func(r chi.Router) {
			r.Post("/", recipeHandler.CreateRecipe)
			r.Get("/", recipeHandler.GetRecipes)
			r.Get("/{id}", recipeHandler.GetRecipeByID)
			r.Get("/author/{author_id}", recipeHandler.GetRecipesByAuthorID)
			r.Put("/{id}", recipeHandler.UpdateRecipe)
			r.Delete("/{id}", recipeHandler.DeleteRecipe)
		})

		// Meal Plan routes
		r.Route("/api/meal-plans", func(r chi.Router) {
			r.Post("/", mealPlanHandler.CreateMealPlan)
			r.Get("/author/{author_id}", mealPlanHandler.GetMealPlansByAuthorID)
			r.Get("/{id}", mealPlanHandler.GetMealPlanByID)
			r.Put("/{id}", mealPlanHandler.UpdateMealPlan)
			r.Delete("/{id}", mealPlanHandler.DeleteMealPlan)
		})
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if PORT env var is not set
	}

	log.Println("Server is starting on port 8080...")
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
