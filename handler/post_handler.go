package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/smilecs/foody/config"
	"github.com/smilecs/foody/data"
	"github.com/smilecs/foody/repository"
	"github.com/smilecs/foody/schema"
)

type PostHandler struct {
	Manager *repository.Manager
}

func NewPostHandler(manager *repository.Manager) *PostHandler {
	return &PostHandler{Manager: manager}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	body := r.FormValue("body")
	tags := r.FormValue("tags")

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("media")
	if err != nil {
		http.Error(w, "Missing media file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine media type based on content type
	contentType := header.Header.Get("Content-Type")
	var mediaType schema.MediaType

	if strings.HasPrefix(contentType, "image/") {
		mediaType = schema.Image
	} else if strings.HasPrefix(contentType, "video/") {
		mediaType = schema.Video
	} else {
		http.Error(w, "Unsupported media type. Only images and videos are allowed.", http.StatusBadRequest)
		return
	}

	cfg := config.Get()
	bucket := cfg.S3_Bucket
	key := fmt.Sprintf("posts/%s/%s", userID.String(), header.Filename)

	url, err := data.UploadFileAndGetUrl(cfg.AWSSess, bucket, key, file, header.Size, contentType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload error: %v", err), http.StatusInternalServerError)
		return
	}

	postID := uuid.New()
	mediaID := uuid.New()

	media := schema.Media{
		Id:        mediaID,
		URL:       url,
		MediaType: mediaType,
		AuthorId:  userID,
	}

	mediaID, err = h.Manager.MediaRepo.CreateMedia(media)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating media: %v", err), http.StatusInternalServerError)
		return
	}

	post := schema.Post{
		Id:       postID,
		Title:    title,
		Body:     body,
		Tags:     tags,
		MediaId:  mediaID,
		AuthorId: userID,
	}

	err = h.Manager.PostRepo.CreatePost(post, mediaID, url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating post: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit := 10
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get posts from repository
	posts, err := h.Manager.PostRepo.GetPosts(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	// Get total count for pagination metadata
	totalCount, err := h.Manager.PostRepo.GetTotalPostsCount()
	if err != nil {
		http.Error(w, "Failed to get total posts count", http.StatusInternalServerError)
		return
	}

	// Create response with pagination metadata
	response := struct {
		Posts      []repository.PostWithMedia `json:"posts"`
		Pagination struct {
			Total  int `json:"total"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		} `json:"pagination"`
	}{
		Posts: posts,
		Pagination: struct {
			Total  int `json:"total"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Total:  totalCount,
			Limit:  limit,
			Offset: offset,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.Manager.PostRepo.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var post schema.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify the post belongs to the user
	existingPost, err := h.Manager.PostRepo.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if existingPost.AuthorId != userID {
		http.Error(w, "Unauthorized to update this post", http.StatusForbidden)
		return
	}

	post.Id = postID
	post.AuthorId = userID

	err = h.Manager.PostRepo.UpdatePost(post)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Verify the post belongs to the user
	existingPost, err := h.Manager.PostRepo.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if existingPost.AuthorId != userID {
		http.Error(w, "Unauthorized to delete this post", http.StatusForbidden)
		return
	}

	err = h.Manager.PostRepo.DeletePost(postID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
