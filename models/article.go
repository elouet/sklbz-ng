package models

import (
	"time"

	"gorm.io/gorm"
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
	gorm.Model        // Adds ID, CreatedAt, UpdatedAt, DeletedAt
	Title            string     `json:"title" gorm:"size:255;not null"`
	Content          string     `json:"content" gorm:"type:text;not null"`
	Author           string     `json:"author" gorm:"size:100;not null"`
	Date             time.Time  `json:"date"`
	Visibility       Visibility `json:"visibility" gorm:"size:20;default:public"`
	Categories       []string   `json:"categories" gorm:"type:text"` // Stored as JSON
	Tags             []string   `json:"tags" gorm:"type:text"`       // Stored as JSON
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
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Author     string    `json:"author"`
	Date       time.Time `json:"date"`
	Visibility Visibility `json:"visibility"`
	Categories []string  `json:"categories"`
	Tags       []string  `json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName sets the table name for Article
func (Article) TableName() string {
	return "articles"
}

// BeforeCreate sets the Date field to current time if not set
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.Date.IsZero() {
		a.Date = time.Now()
	}
	return nil
}

// ArticleRepository handles database operations for articles
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository creates a new ArticleRepository instance
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// AutoMigrate runs database migrations
func (r *ArticleRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&Article{})
}

// Create creates a new article in the database
func (r *ArticleRepository) Create(article *Article) error {
	return r.db.Create(article).Error
}

// Get retrieves an article by ID
func (r *ArticleRepository) Get(id uint) (*Article, error) {
	var article Article
	err := r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetByStringID retrieves an article by string ID (for backward compatibility)
func (r *ArticleRepository) GetByStringID(id string) (*Article, error) {
	var article Article
	err := r.db.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// Update updates an existing article
func (r *ArticleRepository) Update(id uint, updates map[string]interface{}) (*Article, error) {
	var article Article
	err := r.db.Model(&article).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	// Fetch the updated article
	err = r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// Delete removes an article from the database
func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Delete(&Article{}, id).Error
}

// DeleteByStringID removes an article by string ID
func (r *ArticleRepository) DeleteByStringID(id string) error {
	return r.db.Where("id = ?", id).Delete(&Article{}).Error
}

// List returns all articles with optional filtering
func (r *ArticleRepository) List(filters map[string]interface{}) ([]*Article, error) {
	var articles []*Article
	
	query := r.db.Model(&Article{})
	
	// Apply filters
	if author, ok := filters["author"].(string); ok && author != "" {
		query = query.Where("author = ?", author)
	}
	if visibility, ok := filters["visibility"].(string); ok && visibility != "" {
		query = query.Where("visibility = ?", visibility)
	}
	
	// For categories and tags, we need to handle JSON arrays
	// This is a simplified approach - for production, consider using GORM custom types
	if category, ok := filters["category"].(string); ok && category != "" {
		// This is a basic implementation - for better performance, consider using a join table
		query = query.Where("categories LIKE ?", "%"+category+"%")
	}
	if tag, ok := filters["tag"].(string); ok && tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	
	err := query.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	
	return articles, nil
}

// ListAll returns all articles without filtering
func (r *ArticleRepository) ListAll() ([]*Article, error) {
	var articles []*Article
	err := r.db.Find(&articles).Error
	if err != nil {
		return nil, err
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
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}
