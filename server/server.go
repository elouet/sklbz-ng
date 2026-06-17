package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"sklbz-ng/config"
	"sklbz-ng/handlers"
	"sklbz-ng/middleware"
	"sklbz-ng/models"
)

// Server represents the HTTP server
type Server struct {
	router      *mux.Router
	repo        *models.ArticleRepository
	port        string
	staticDir   string
	templateDir string
	mtlsConfig  *config.MTLSConfig
	useHTTPS    bool
	certFile    string
	keyFile     string
}

// NewServer creates a new Server instance
func NewServer(port, dbPath, staticDir, templateDir string) *Server {
	router := mux.NewRouter()
	
	// Initialize database
	db := config.InitDB(dbPath)
	
	// Auto-migrate database
	repo := models.NewArticleRepository(db)
	if err := repo.AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	
	return &Server{
		router:      router,
		repo:        repo,
		port:        port,
		staticDir:   staticDir,
		templateDir: templateDir,
		mtlsConfig:  config.DefaultMTLSConfig(),
		useHTTPS:    false,
	}
}

// NewMTLSServer creates a new Server instance with mTLS support
func NewMTLSServer(port, dbPath, staticDir, templateDir, certFile, keyFile, caCertFile string) *Server {
	server := NewServer(port, dbPath, staticDir, templateDir)
	
	// Load mTLS configuration
	mtlsConfig, err := config.LoadMTLSConfig(certFile, keyFile, caCertFile)
	if err != nil {
		log.Printf("Warning: Failed to load mTLS config: %v", err)
		mtlsConfig = config.DefaultMTLSConfig()
	}
	
	server.mtlsConfig = mtlsConfig
	server.useHTTPS = mtlsConfig.Enabled
	server.certFile = certFile
	server.keyFile = keyFile
	
	return server
}

// SetupRoutes configures all server routes
func (s *Server) SetupRoutes() {
	// Apply mTLS middleware to all routes
	// This middleware will check for client certificates on write operations
	writeAuthMiddleware := middleware.WriteOnlyAuth(s.mtlsConfig)
	
	// API routes for articles - with write authentication
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	apiRouter.Use(writeAuthMiddleware)
	handlers.RegisterArticleRoutes(apiRouter, s.repo)

	// HTML routes for articles
	handlers.RegisterHTMLRoutes(s.router, s.repo, s.templateDir)

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
	w.Write([]byte(`{"message": "Welcome to sklbz-ng API with GORM", "endpoints": {"articles": "/api/articles", "static": "/static/"}}`))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.SetupRoutes()
	
	if s.useHTTPS {
		log.Printf("Starting HTTPS server on port %s with mTLS", s.port)
		log.Printf("Database: SQLite (GORM)")
		log.Printf("Static files directory: %s", s.staticDir)
		log.Printf("Templates directory: %s", s.templateDir)
		log.Printf("Server certificate: %s", s.certFile)
		log.Printf("Server key: %s", s.keyFile)
		
		// Create HTTPS server with mTLS
		server, err := config.CreateMTLSListener(":"+s.port, s.certFile, s.keyFile, s.mtlsConfig)
		if err != nil {
			return fmt.Errorf("failed to create HTTPS server: %v", err)
		}
		
		server.Handler = s.router
		
		return server.ListenAndServeTLS("", "")
	}
	
	log.Printf("Starting HTTP server on port %s", s.port)
	log.Printf("Database: SQLite (GORM)")
	log.Printf("Static files directory: %s", s.staticDir)
	log.Printf("Templates directory: %s", s.templateDir)
	if s.mtlsConfig.Enabled {
		log.Printf("Warning: mTLS is configured but server is running in HTTP mode. Use HTTPS for mTLS to work properly.")
	}
	
	return http.ListenAndServe(":"+s.port, s.router)
}

// GetRouter returns the router for testing purposes
func (s *Server) GetRouter() *mux.Router {
	return s.router
}

// Close closes the database connection
func (s *Server) Close() {
	config.CloseDB()
}
