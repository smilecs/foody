package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "github.com/smilecs/foody/db"
    "github.com/smilecs/foody/repository"
    "github.com/smilecs/foody/utils"
    "net/http"
)

func main() {
	database := db.Init()
	repository.NewManager(database)
}
