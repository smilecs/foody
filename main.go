package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/handler"
	"github.com/smilecs/foody/middleware"
	"github.com/smilecs/foody/repository"
)

func main() {
	cfg := config.Init()

	// Create a SQLDatabase wrapper around the sqlx.DB
	dbWrapper := &config.SQLDatabase{
		DB: cfg.DB,
	}

	// Initialize manager with all repositories
	manager := repository.NewManager(dbWrapper)

	// Initialize handlers with manager
	userHandler := handler.NewUserHandler(manager)
	postHandler := handler.NewPostHandler(manager)
	recipeHandler := handler.NewRecipeHandler(manager)

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
	})

	// Start the server
	http.ListenAndServe(":8080", router)
}
