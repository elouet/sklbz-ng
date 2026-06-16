package server

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"sklbz-ng/handlers"
	"sklbz-ng/models"
)

// Server represents the HTTP server
type Server struct {
	router   *mux.Router
	storage  *models.ArticleStorage
	port     string
	staticDir string
}

// NewServer creates a new Server instance
func NewServer(port, dataDir, staticDir string) *Server {
	router := mux.NewRouter()
	storage := models.NewArticleStorage(dataDir)
	
	return &Server{
		router:   router,
		storage:  storage,
		port:     port,
		staticDir: staticDir,
	}
}

// SetupRoutes configures all server routes
func (s *Server) SetupRoutes() {
	// API routes for articles
	handlers.RegisterArticleRoutes(s.router, s.storage)

	// Static file routes
	// Serve JavaScript files
	jsPath := filepath.Join(s.staticDir, "js")
	s.router.PathPrefix("/static/js/").Handler(http.StripPrefix("/static/js/", http.FileServer(http.Dir(jsPath))))
	
	// Serve CSS files
	cssPath := filepath.Join(s.staticDir, "css")
	s.router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir(cssPath))))
	
	// Serve other static files (images, etc.)
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticDir))))

	// Health check endpoint
	s.router.HandleFunc("/health", s.healthCheck).Methods("GET")

	// Root endpoint
	s.router.HandleFunc("/", s.rootHandler).Methods("GET")
}

// healthCheck returns server health status
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

// rootHandler serves the root endpoint
func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Welcome to sklbz-ng API", "endpoints": {"articles": "/api/articles", "static": "/static/"}}`))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.SetupRoutes()
	
	log.Printf("Starting server on port %s", s.port)
	log.Printf("Data directory: %s", s.storage.basePath)
	log.Printf("Static files directory: %s", s.staticDir)
	
	return http.ListenAndServe(":"+s.port, s.router)
}

// GetRouter returns the router for testing purposes
func (s *Server) GetRouter() *mux.Router {
	return s.router
}
