package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/data"
	"github.com/smilecs/foody/repository"
	"github.com/smilecs/foody/schema"
	"github.com/smilecs/foody/utils"
)

type UserHandler struct {
	Manager *repository.Manager
}

func NewUserHandler(manager *repository.Manager) *UserHandler {
	return &UserHandler{
		Manager: manager,
	}
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	dob := r.FormValue("date_of_birth")
	password := r.FormValue("password")

	// Validate required fields
	if username == "" || email == "" || dob == "" || password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !strings.Contains(email, "@") {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("media")
	if err != nil {
		http.Error(w, "Missing Profile Image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if file is empty
	if header.Size == 0 {
		http.Error(w, "Profile image cannot be empty", http.StatusBadRequest)
		return
	}

	cfg := config.Get()
	bucket := cfg.S3_Bucket
	key := fmt.Sprintf("users/%s/%s", username, header.Filename)

	url, err := data.UploadFileAndGetUrl(cfg.AWSSess, bucket, key, file, header.Size, header.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload error: %v", err), http.StatusInternalServerError)
		return
	}

	user_id := uuid.New()
	media_id := uuid.New()

	// Create media first
	media := schema.Media{
		Id:        media_id,
		URL:       url,
		MediaType: schema.Image,
		AuthorId:  user_id,
	}

	mediaID, err := u.Manager.MediaRepo.CreateMedia(media)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating media: %v", err), http.StatusInternalServerError)
		return
	}

	// Create user with media ID
	user := schema.User{
		Id:       user_id,
		Email:    email,
		Name:     username,
		Username: username,
		DOB:      dob,
	}

	_, err = u.Manager.UserRepo.CreateUser(user, password, mediaID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	authenticated, err := u.Manager.UserRepo.AuthUser(email, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication error: %v", err), http.StatusInternalServerError)
		return
	}

	if !authenticated {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Get user ID from repository
	user, err := u.Manager.UserRepo.GetUserByEmail(email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving user: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Id, user.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating token: %v", err), http.StatusInternalServerError)
		return
	}

	// Return token in response
	response := LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
