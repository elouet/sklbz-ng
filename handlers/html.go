package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"sklbz-ng/models"
)

// TemplateHandler manages HTML template rendering
type TemplateHandler struct {
	repo        *models.ArticleRepository
	templates   *template.Template
	templateDir string
}

// NewTemplateHandler creates a new TemplateHandler
func NewTemplateHandler(repo *models.ArticleRepository, templateDir string) *TemplateHandler {
	handler := &TemplateHandler{
		repo:      repo,
		templateDir: templateDir,
	}
	
	// Load templates
	handler.loadTemplates()
	
	return handler
}

// loadTemplates loads all HTML templates from the template directory
func (h *TemplateHandler) loadTemplates() {
	var err error
	
	// Pattern to match all .html files in the templates directory
	pattern := filepath.Join(h.templateDir, "*.html")
	
	h.templates, err = template.ParseGlob(pattern)
	if err != nil {
		log.Printf("Error loading templates: %v", err)
		// Create a minimal template for error display
		h.templates = template.Must(template.New("error").Parse("<html><body><h1>Error loading templates</h1></body></html>"))
	}
	
	log.Printf("Loaded templates from: %s", h.templateDir)
}

// RenderTemplate renders a template with the given data
func (h *TemplateHandler) RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", templateName, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// IndexHandler renders the home page
func (h *TemplateHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	h.RenderTemplate(w, "index.html", nil)
}

// ArticlesHandler renders the articles list page
func (h *TemplateHandler) ArticlesHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Error listing articles", http.StatusInternalServerError)
		return
	}

	// Convert to responses for template
	var responses []*models.ArticleResponse
	for _, article := range articles {
		responses = append(responses, article.ToResponse())
	}

	data := struct {
		Articles []*models.ArticleResponse
	}{
		Articles: responses,
	}

	h.RenderTemplate(w, "articles.html", data)
}

// ArticleHandler renders a single article page
func (h *TemplateHandler) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Try to parse as uint first
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		// Try string ID for backward compatibility
		article, err := h.repo.GetByStringID(idStr)
		if err != nil {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}
		h.RenderTemplate(w, "article.html", article.ToResponse())
		return
	}

	article, err := h.repo.Get(uint(id))
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	h.RenderTemplate(w, "article.html", article.ToResponse())
}

// RegisterHTMLRoutes registers HTML routes on the router
func RegisterHTMLRoutes(router *mux.Router, repo *models.ArticleRepository, templateDir string) {
	handler := NewTemplateHandler(repo, templateDir)

	// HTML routes
	router.HandleFunc("/", handler.IndexHandler).Methods("GET")
	router.HandleFunc("/articles", handler.ArticlesHandler).Methods("GET")
	router.HandleFunc("/articles/{id}", handler.ArticleHandler).Methods("GET")
}
