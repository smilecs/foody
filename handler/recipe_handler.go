package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/smilecs/foody/repository"
	"github.com/smilecs/foody/schema"
)

type RecipeHandler struct {
	Manager *repository.Manager
}

func NewRecipeHandler(manager *repository.Manager) *RecipeHandler {
	return &RecipeHandler{Manager: manager}
}

func (h *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe schema.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate UUID for recipe
	recipe.Id = uuid.New()

	// Handle media upload if present
	var mediaID uuid.UUID
	var mediaURL string
	if recipe.MediaId != uuid.Nil {
		media, err := h.Manager.MediaRepo.GetMediaByID(recipe.MediaId)
		if err != nil {
			http.Error(w, "Media not found", http.StatusNotFound)
			return
		}
		mediaID = media.Id
		mediaURL = media.URL
	}

	if err := h.Manager.RecipeRepo.CreateRecipe(recipe, mediaID, mediaURL); err != nil {
		http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	recipes, err := h.Manager.RecipeRepo.GetRecipes(limit, offset)
	if err != nil {
		http.Error(w, "Failed to get recipes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) GetRecipeByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	recipe, err := h.Manager.RecipeRepo.GetRecipeByID(id)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) GetRecipesByAuthorID(w http.ResponseWriter, r *http.Request) {
	authorIDStr := chi.URLParam(r, "author_id")
	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		http.Error(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	recipes, err := h.Manager.RecipeRepo.GetRecipesByAuthorID(authorID)
	if err != nil {
		http.Error(w, "Failed to get recipes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	var recipe schema.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	recipe.Id = id
	if err := h.Manager.RecipeRepo.UpdateRecipe(recipe); err != nil {
		http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipe)
}

func (h *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	if err := h.Manager.RecipeRepo.DeleteRecipe(id); err != nil {
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
