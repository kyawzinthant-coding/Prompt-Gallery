package models

import (
	"gorm.io/gorm"
	"time"
)

type Prompt struct {
	gorm.Model

	Title       string `gorm:"not null;size:200" json:"title"`
	Description string `gorm:"type:text;not null" json:"description"`

	Language string `gorm:"not null;size:50;index" json:"language"`

	// Difficulty level
	Difficulty DifficultyLevel `gorm:"not null;index" json:"difficulty"`

	// Category/Topic (e.g., "algorithms", "web-development", "data-structures")
	Category string `gorm:"not null;size:100;index" json:"category"`

	// Prompt content and examples
	ProblemStatement string `gorm:"type:text;not null" json:"problem_statement"`

	// Quality control
	IsVerified bool       `gorm:"default:false;index" json:"is_verified"`
	VerifiedBy *uint      `gorm:"index" json:"verified_by,omitempty"` // Foreign key to User (future)
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	// Engagement metrics
	ViewCount      int `gorm:"default:0" json:"view_count"`
	LikeCount      int `gorm:"default:0" json:"like_count"`
	DifficultyVote int `gorm:"default:0" json:"difficulty_vote"` // Average difficulty rating

	Tags string `gorm:"type:text" json:"tags"` // JSON array of tags

	// Author information (for future user system)
	AuthorID    *uint  `gorm:"index" json:"author_id,omitempty"`
	AuthorName  string `gorm:"size:100" json:"author_name,omitempty"`
	AuthorEmail string `gorm:"size:100" json:"author_email,omitempty"`
}

type DifficultyLevel string

const (
	DifficultyBeginner     DifficultyLevel = "beginner"
	DifficultyIntermediate DifficultyLevel = "intermediate"
	DifficultyAdvanced     DifficultyLevel = "advanced"
	DifficultyExpert       DifficultyLevel = "expert"
)

// Valid checks if the difficulty level is valid
func (d DifficultyLevel) Valid() bool {
	switch d {
	case DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced, DifficultyExpert:
		return true
	}
	return false
}

func (Prompt) TableName() string {
	return "prompts"
}

func (p *Prompt) BeforeCreate(tx *gorm.DB) error {
	// Validate difficulty level
	if !p.Difficulty.Valid() {
		p.Difficulty = DifficultyBeginner // Default value
	}
	return nil
}

// PromptFilter represents filtering options for prompts
// Used for search and filtering functionality
type PromptFilter struct {
	Language   string          `json:"language,omitempty"`
	Difficulty DifficultyLevel `json:"difficulty,omitempty"`
	Category   string          `json:"category,omitempty"`
	IsVerified *bool           `json:"is_verified,omitempty"`
	Search     string          `json:"search,omitempty"` // Search in title/description
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
}

// PromptCreateRequest represents the request to create a new prompt
// Similar to DTO (Data Transfer Object) in Java
type PromptCreateRequest struct {
	Title            string          `json:"title" validate:"required,max=200"`
	Description      string          `json:"description" validate:"required"`
	Language         string          `json:"language" validate:"required,max=50"`
	Difficulty       DifficultyLevel `json:"difficulty" validate:"required"`
	Category         string          `json:"category" validate:"required,max=100"`
	ProblemStatement string          `json:"problem_statement" validate:"required"`
	Examples         string          `json:"examples,omitempty"`
	Hints            string          `json:"hints,omitempty"`
	Tags             string          `json:"tags,omitempty"`
	EstimatedTime    int             `json:"estimated_time,omitempty"`
	AuthorName       string          `json:"author_name,omitempty" validate:"max=100"`
	AuthorEmail      string          `json:"author_email,omitempty" validate:"email,max=100"`
}

// ToPrompt converts PromptCreateRequest to Prompt model
func (req *PromptCreateRequest) ToPrompt() *Prompt {
	return &Prompt{
		Title:            req.Title,
		Description:      req.Description,
		Language:         req.Language,
		Difficulty:       req.Difficulty,
		Category:         req.Category,
		ProblemStatement: req.ProblemStatement,

		Tags:        req.Tags,
		AuthorName:  req.AuthorName,
		AuthorEmail: req.AuthorEmail,
	}
}
