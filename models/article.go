package models

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Visibility represents the visibility level of an article
type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
	VisibilityDraft   Visibility = "draft"
)

// Article represents a blog article with metadata
type Article struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Author     string      `json:"author"`
	Date       time.Time   `json:"date"`
	Visibility Visibility  `json:"visibility"`
	Categories []string    `json:"categories"`
	Tags       []string    `json:"tags"`
	UpdatedAt  time.Time   `json:"updated_at"`
	CreatedAt  time.Time   `json:"created_at"`
}

// ArticleRequest represents the request payload for creating/updating an article
type ArticleRequest struct {
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Author     string     `json:"author"`
	Visibility Visibility `json:"visibility"`
	Categories []string   `json:"categories"`
	Tags       []string   `json:"tags"`
}

// ArticleResponse represents the response payload for an article
type ArticleResponse struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Author     string     `json:"author"`
	Date       time.Time  `json:"date"`
	Visibility Visibility `json:"visibility"`
	Categories []string   `json:"categories"`
	Tags       []string   `json:"tags"`
	UpdatedAt  time.Time  `json:"updated_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ArticleStorage handles JSON file storage for articles
type ArticleStorage struct {
	BasePath string
}

// NewArticleStorage creates a new ArticleStorage instance
func NewArticleStorage(basePath string) *ArticleStorage {
	return &ArticleStorage{BasePath: basePath}
}

// ensureDirectory ensures the storage directory exists
func (s *ArticleStorage) ensureDirectory() error {
	if err := os.MkdirAll(s.BasePath, 0755); err != nil {
		return err
	}
	return nil
}

// getFilePath returns the full path for an article JSON file
func (s *ArticleStorage) getFilePath(id string) string {
	return filepath.Join(s.BasePath, id+".json")
}

// Save saves an article to disk as JSON
func (s *ArticleStorage) Save(article *Article) error {
	if err := s.ensureDirectory(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(article, "", "  ")
	if err != nil {
		return err
	}

	filePath := s.getFilePath(article.ID)
	return os.WriteFile(filePath, data, 0644)
}

// Get retrieves an article from disk by ID
func (s *ArticleStorage) Get(id string) (*Article, error) {
	filePath := s.getFilePath(id)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("article not found")
		}
		return nil, err
	}

	var article Article
	if err := json.Unmarshal(data, &article); err != nil {
		return nil, err
	}

	return &article, nil
}

// Delete removes an article from disk
func (s *ArticleStorage) Delete(id string) error {
	filePath := s.getFilePath(id)
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.New("article not found")
	}

	return os.Remove(filePath)
}

// List returns all article IDs
func (s *ArticleStorage) List() ([]string, error) {
	files, err := os.ReadDir(s.BasePath)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			id := file.Name()[:len(file.Name())-5] // Remove .json extension
			ids = append(ids, id)
		}
	}

	return ids, nil
}

// ListAll returns all articles
func (s *ArticleStorage) ListAll() ([]*Article, error) {
	ids, err := s.List()
	if err != nil {
		return nil, err
	}

	var articles []*Article
	for _, id := range ids {
		article, err := s.Get(id)
		if err != nil {
			continue // Skip articles that can't be read
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// ConvertToResponse converts an Article to ArticleResponse
func (a *Article) ToResponse() *ArticleResponse {
	return &ArticleResponse{
		ID:         a.ID,
		Title:      a.Title,
		Content:    a.Content,
		Author:     a.Author,
		Date:       a.Date,
		Visibility: a.Visibility,
		Categories: a.Categories,
		Tags:       a.Tags,
		UpdatedAt:  a.UpdatedAt,
		CreatedAt:  a.CreatedAt,
	}
}
