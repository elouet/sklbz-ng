package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"sklbz-ng/models"
)

// ArticleHandler manages article-related HTTP handlers
type ArticleHandler struct {
	repo *models.ArticleRepository
}

// NewArticleHandler creates a new ArticleHandler
func NewArticleHandler(repo *models.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{repo: repo}
}

// CreateArticle creates a new article
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req models.ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" || req.Content == "" || req.Author == "" {
		http.Error(w, "Title, content, and author are required", http.StatusBadRequest)
		return
	}

	// Set default visibility if not provided
	if req.Visibility == "" {
		req.Visibility = models.VisibilityPublic
	}

	// Create article
	article := &models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Author:     req.Author,
		Date:       time.Now(),
		Visibility: req.Visibility,
		Categories: req.Categories,
		Tags:       req.Tags,
	}

	// Save to database
	if err := h.repo.Create(article); err != nil {
		log.Printf("Error saving article: %v", err)
		http.Error(w, fmt.Sprintf("Error saving article: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article.ToResponse())
}

// GetArticle retrieves an article by ID
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Try to parse as uint first
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// Try string ID for backward compatibility
		article, err := h.repo.GetByStringID(idStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Article not found: %v", err), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(article.ToResponse())
		return
	}

	article, err := h.repo.Get(uint(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Article not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article.ToResponse())
}

// UpdateArticle updates an existing article
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Try to parse as uint first
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// Try string ID for backward compatibility
		article, err := h.repo.GetByStringID(idStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Article not found: %v", err), http.StatusNotFound)
			return
		}
		id = uint(article.ID)
	}

	// Decode update request
	var req models.ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Author != "" {
		updates["author"] = req.Author
	}
	if req.Visibility != "" {
		updates["visibility"] = req.Visibility
	}
	if req.Categories != nil {
		updates["categories"] = req.Categories
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	updates["updated_at"] = time.Now()

	// Update article
	article, err := h.repo.Update(id, updates)
	if err != nil {
		log.Printf("Error updating article: %v", err)
		http.Error(w, fmt.Sprintf("Error updating article: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article.ToResponse())
}

// DeleteArticle deletes an article by ID
func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Try to parse as uint first
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// Try string ID for backward compatibility
		err := h.repo.DeleteByStringID(idStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting article: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting article: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListArticles returns all articles
func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	query := r.URL.Query()
	filters := make(map[string]interface{})
	
	if author := query.Get("author"); author != "" {
		filters["author"] = author
	}
	if category := query.Get("category"); category != "" {
		filters["category"] = category
	}
	if tag := query.Get("tag"); tag != "" {
		filters["tag"] = tag
	}
	if visibility := query.Get("visibility"); visibility != "" {
		filters["visibility"] = visibility
	}

	var articles []*models.Article
	var err error
	
	if len(filters) > 0 {
		articles, err = h.repo.List(filters)
	} else {
		articles, err = h.repo.ListAll()
	}
	
	if err != nil {
		log.Printf("Error listing articles: %v", err)
		http.Error(w, fmt.Sprintf("Error listing articles: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert to responses
	var responses []*models.ArticleResponse
	for _, article := range articles {
		responses = append(responses, article.ToResponse())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// RegisterArticleRoutes registers article routes on the router
func RegisterArticleRoutes(router *mux.Router, repo *models.ArticleRepository) {
	handler := NewArticleHandler(repo)

	router.HandleFunc("/api/articles", handler.CreateArticle).Methods("POST")
	router.HandleFunc("/api/articles", handler.ListArticles).Methods("GET")
	router.HandleFunc("/api/articles/{id}", handler.GetArticle).Methods("GET")
	router.HandleFunc("/api/articles/{id}", handler.UpdateArticle).Methods("PUT")
	router.HandleFunc("/api/articles/{id}", handler.DeleteArticle).Methods("DELETE")
}
