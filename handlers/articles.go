package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"sklbz-ng/models"
)

// ArticleHandler manages article-related HTTP handlers
type ArticleHandler struct {
	storage *models.ArticleStorage
}

// NewArticleHandler creates a new ArticleHandler
func NewArticleHandler(storage *models.ArticleStorage) *ArticleHandler {
	return &ArticleHandler{storage: storage}
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

	// Generate ID (using timestamp for simplicity)
	id := fmt.Sprintf("article-%d", time.Now().UnixNano())

	// Create article
	article := &models.Article{
		ID:         id,
		Title:      req.Title,
		Content:    req.Content,
		Author:     req.Author,
		Date:       time.Now(),
		Visibility: req.Visibility,
		Categories: req.Categories,
		Tags:       req.Tags,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to storage
	if err := h.storage.Save(article); err != nil {
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
	id := vars["id"]

	article, err := h.storage.Get(id)
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
	id := vars["id"]

	// Get existing article
	article, err := h.storage.Get(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Article not found: %v", err), http.StatusNotFound)
		return
	}

	// Decode update request
	var req models.ArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Update fields
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Author != "" {
		article.Author = req.Author
	}
	if req.Visibility != "" {
		article.Visibility = req.Visibility
	}
	if req.Categories != nil {
		article.Categories = req.Categories
	}
	if req.Tags != nil {
		article.Tags = req.Tags
	}

	article.UpdatedAt = time.Now()

	// Save updated article
	if err := h.storage.Save(article); err != nil {
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
	id := vars["id"]

	if err := h.storage.Delete(id); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting article: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListArticles returns all articles
func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	query := r.URL.Query()
	author := query.Get("author")
	category := query.Get("category")
	tag := query.Get("tag")
	visibility := query.Get("visibility")

	articles, err := h.storage.ListAll()
	if err != nil {
		log.Printf("Error listing articles: %v", err)
		http.Error(w, fmt.Sprintf("Error listing articles: %v", err), http.StatusInternalServerError)
		return
	}

	// Apply filters
	var filtered []*models.Article
	for _, article := range articles {
		// Filter by author
		if author != "" && article.Author != author {
			continue
		}

		// Filter by category
		if category != "" {
			found := false
			for _, cat := range article.Categories {
				if cat == category {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by tag
		if tag != "" {
			found := false
			for _, t := range article.Tags {
				if t == tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filter by visibility
		if visibility != "" && article.Visibility != models.Visibility(visibility) {
			continue
		}

		filtered = append(filtered, article)
	}

	// Convert to responses
	var responses []*models.ArticleResponse
	for _, article := range filtered {
		responses = append(responses, article.ToResponse())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// RegisterArticleRoutes registers article routes on the router
func RegisterArticleRoutes(router *mux.Router, storage *models.ArticleStorage) {
	handler := NewArticleHandler(storage)

	router.HandleFunc("/api/articles", handler.CreateArticle).Methods("POST")
	router.HandleFunc("/api/articles", handler.ListArticles).Methods("GET")
	router.HandleFunc("/api/articles/{id}", handler.GetArticle).Methods("GET")
	router.HandleFunc("/api/articles/{id}", handler.UpdateArticle).Methods("PUT")
	router.HandleFunc("/api/articles/{id}", handler.DeleteArticle).Methods("DELETE")
}
