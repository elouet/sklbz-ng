package config

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database configuration
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB(dbPath string) *gorm.DB {
	var err error

	// Check if database file exists, create directory if needed
	dir := os.Getenv("DB_DIR")
	if dir == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Open SQLite database
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Connected to database: %s", dbPath)
	return DB
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}

// CloseDB closes the database connection
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting database connection: %v", err)
		return
	}
	sqlDB.Close()
}
