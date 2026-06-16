package main

import (
	"flag"
	"log"
	"os"

	"sklbz-ng/server"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8080", "Port to listen on")
	dataDir := flag.String("data", "./data", "Directory to store article JSON files")
	staticDir := flag.String("static", "./static", "Directory to serve static files from")

	flag.Parse()

	// Create and start server
	s := server.NewServer(*port, *dataDir, *staticDir)
	
	// Ensure data directory exists
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Ensure static directory exists
	if err := os.MkdirAll(*staticDir, 0755); err != nil {
		log.Fatalf("Failed to create static directory: %v", err)
	}

	log.Printf("Server starting on port %s", *port)
	log.Printf("Data directory: %s", *dataDir)
	log.Printf("Static directory: %s", *staticDir)
	
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
