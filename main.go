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

	flag.Parse()

	// Ensure static directory exists
	if err := os.MkdirAll(*staticDir, 0755); err != nil {
		log.Fatalf("Failed to create static directory: %v", err)
	}

	// Create and start server
	s := server.NewServer(*port, *dbPath, *staticDir)
	
	log.Printf("Server starting on port %s", *port)
	log.Printf("Database: %s", *dbPath)
	log.Printf("Static directory: %s", *staticDir)
	
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
