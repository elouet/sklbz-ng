package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sklbz-ng/server"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Port to listen on")
	dbPath := flag.String("db", "./sklbz.db", "Path to SQLite database file")
	staticDir := flag.String("static", "./static", "Directory to serve static files from")
	templateDir := flag.String("templates", "./templates", "Directory containing HTML templates")
	
	// mTLS flags
	mtlsEnabled := flag.Bool("mtls", false, "Enable mutual TLS authentication")
	certFile := flag.String("cert", "", "Path to server certificate file (PEM)")
	keyFile := flag.String("key", "", "Path to server private key file (PEM)")
	caCertFile := flag.String("ca-cert", "", "Path to CA certificate file for client verification (PEM)")

	flag.Parse()

	// Ensure directories exist
	if err := os.MkdirAll(*staticDir, 0755); err != nil {
		log.Fatalf("Failed to create static directory: %v", err)
	}
	if err := os.MkdirAll(*templateDir, 0755); err != nil {
		log.Fatalf("Failed to create templates directory: %v", err)
	}

	var s *server.Server
	
	// Check if mTLS is enabled via flag or environment variable
	if *mtlsEnabled || os.Getenv("MTLS_ENABLED") == "true" {
		// Verify required files exist
		if *certFile == "" || *keyFile == "" || *caCertFile == "" {
			log.Fatal("mTLS enabled but missing required certificate files. Please provide -cert, -key, and -ca-cert flags.")
		}
		
		// Check if files exist
		if _, err := os.Stat(*certFile); os.IsNotExist(err) {
			log.Fatalf("Server certificate file not found: %s", *certFile)
		}
		if _, err := os.Stat(*keyFile); os.IsNotExist(err) {
			log.Fatalf("Server key file not found: %s", *keyFile)
		}
		if _, err := os.Stat(*caCertFile); os.IsNotExist(err) {
			log.Fatalf("CA certificate file not found: %s", *caCertFile)
		}
		
		log.Printf("mTLS authentication enabled")
		log.Printf("Server certificate: %s", *certFile)
		log.Printf("Server key: %s", *keyFile)
		log.Printf("CA certificate: %s", *caCertFile)
		
		s = server.NewMTLSServer(*port, *dbPath, *staticDir, *templateDir, *certFile, *keyFile, *caCertFile)
	} else {
		s = server.NewServer(*port, *dbPath, *staticDir, *templateDir)
	}
	
	log.Printf("Server starting on port %s", *port)
	log.Printf("Database: %s", *dbPath)
	log.Printf("Static directory: %s", *staticDir)
	log.Printf("Templates directory: %s", *templateDir)
	
	// Start server in a goroutine
	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	s.Close()
	log.Println("Server stopped")
}
