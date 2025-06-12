package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/smilecs/foody/repository"
	"github.com/smilecs/foody/schema"
)

type MealPlanHandler struct {
	Manager *repository.Manager
}

func NewMealPlanHandler(manager *repository.Manager) *MealPlanHandler {
	return &MealPlanHandler{Manager: manager}
}

func (h *MealPlanHandler) CreateMealPlan(w http.ResponseWriter, r *http.Request) {
	var mealPlan schema.MealPlan
	if err := json.NewDecoder(r.Body).Decode(&mealPlan); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Set the author ID
	mealPlan.Id = uuid.New()
	mealPlan.AuthorId = userID

	if err := h.Manager.MealPlanRepo.CreateMealPlan(mealPlan); err != nil {
		http.Error(w, "Failed to create meal plan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mealPlan)
}

func (h *MealPlanHandler) GetMealPlansByAuthorID(w http.ResponseWriter, r *http.Request) {
	authorIDStr := chi.URLParam(r, "author_id")
	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		http.Error(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	mealPlans, err := h.Manager.MealPlanRepo.GetMealPlansByAuthorID(authorID)
	if err != nil {
		http.Error(w, "Failed to get meal plans", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(mealPlans)
}

func (h *MealPlanHandler) GetMealPlanByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid meal plan ID", http.StatusBadRequest)
		return
	}

	mealPlan, err := h.Manager.MealPlanRepo.GetMealPlanByID(id)
	if err != nil {
		http.Error(w, "Meal plan not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(mealPlan)
}

func (h *MealPlanHandler) UpdateMealPlan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid meal plan ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Get existing meal plan to verify ownership
	existingMealPlan, err := h.Manager.MealPlanRepo.GetMealPlanByID(id)
	if err != nil {
		http.Error(w, "Meal plan not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if existingMealPlan.AuthorId != userID {
		http.Error(w, "Unauthorized to update this meal plan", http.StatusForbidden)
		return
	}

	var mealPlan schema.MealPlan
	if err := json.NewDecoder(r.Body).Decode(&mealPlan); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mealPlan.Id = id
	mealPlan.AuthorId = userID

	if err := h.Manager.MealPlanRepo.UpdateMealPlan(mealPlan); err != nil {
		http.Error(w, "Failed to update meal plan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mealPlan)
}

func (h *MealPlanHandler) DeleteMealPlan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid meal plan ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Get existing meal plan to verify ownership
	existingMealPlan, err := h.Manager.MealPlanRepo.GetMealPlanByID(id)
	if err != nil {
		http.Error(w, "Meal plan not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if existingMealPlan.AuthorId != userID {
		http.Error(w, "Unauthorized to delete this meal plan", http.StatusForbidden)
		return
	}

	if err := h.Manager.MealPlanRepo.DeleteMealPlan(id); err != nil {
		http.Error(w, "Failed to delete meal plan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
