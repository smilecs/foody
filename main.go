package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/middleware"
	"github.com/smilecs/foody/repository"
	"github.com/smilecs/foody/routes/handlers"
)

func main() {
	cfg := config.Init()

	// Create a SQLDatabase wrapper around the sqlx.DB
	dbWrapper := &config.SQLDatabase{
		DB: cfg.DB,
	}

	// Initialize the repository manager with the wrapped database
	manager := repository.NewManager(dbWrapper)

	// Create the user handler with the manager
	userHandler := handlers.NewUserHandler(manager)

	// Set up your routes
	router := chi.NewRouter()

	// Public routes
	router.Post("/signup", userHandler.CreateUser)
	router.Post("/login", userHandler.Login)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		// Add protected routes here
		// r.Post("/posts", postHandler.CreatePost)
		// r.Get("/posts", postHandler.GetPosts)
		// etc...
	})

	// Start the server
	http.ListenAndServe(":8080", router)
}
