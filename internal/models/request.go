package models

import (
	"gorm.io/gorm"
)

// PromptRequest represents a request for a custom prompt from users
// This is when users want to request specific prompts to be created by authorized users
// Similar to a ContactForm or RequestForm model in Express.js
type PromptRequest struct {
	gorm.Model

	// Requester information (anyone can submit requests)
	RequesterName  string `gorm:"not null;size:100" json:"requester_name"`
	RequesterEmail string `gorm:"not null;size:100;index" json:"requester_email"`

	// What kind of prompt they want
	RequestedTitle      string          `gorm:"not null;size:200" json:"requested_title"`
	RequestedLanguage   string          `gorm:"not null;size:50" json:"requested_language"`
	RequestedDifficulty DifficultyLevel `gorm:"not null" json:"requested_difficulty"`
	RequestedCategory   string          `gorm:"not null;size:100" json:"requested_category"`

	// Detailed description of what they want
	Description          string `gorm:"type:text;not null" json:"description"`
	SpecificRequirements string `gorm:"type:text" json:"specific_requirements,omitempty"`

	// Additional context
	UseCase         string `gorm:"type:text" json:"use_case,omitempty"`         // What they'll use it for
	PreferredTopics string `gorm:"type:text" json:"preferred_topics,omitempty"` // Specific topics they want covered

	// Request management
	Status   RequestStatus `gorm:"not null;default:'pending';index" json:"status"`
	Priority Priority      `gorm:"not null;default:'normal'" json:"priority"`

	// Admin/Expert assignment (for future)
	AssignedToID *uint  `gorm:"index" json:"assigned_to_id,omitempty"` // Who will create the prompt
	AssignedBy   *uint  `gorm:"index" json:"assigned_by_id,omitempty"` // Who assigned it
	AssignedAt   *int64 `json:"assigned_at,omitempty"`

	// Completion tracking
	CompletedPromptID *uint  `gorm:"index" json:"completed_prompt_id,omitempty"` // Link to created prompt
	CompletedAt       *int64 `json:"completed_at,omitempty"`

	// Communication
	AdminNotes      string `gorm:"type:text" json:"admin_notes,omitempty"`      // Internal notes
	ResponseMessage string `gorm:"type:text" json:"response_message,omitempty"` // Response to requester

	// Flags
	IsUrgent   bool `gorm:"default:false;index" json:"is_urgent"`
	IsRejected bool `gorm:"default:false;index" json:"is_rejected"`

	// Estimated effort (set by admins)
	EstimatedHours int `gorm:"default:0" json:"estimated_hours,omitempty"`
}

// RequestStatus represents the status of a prompt request
// Similar to order status or ticket status in Express.js apps
type RequestStatus string

const (
	StatusPending    RequestStatus = "pending"     // Just submitted
	StatusInReview   RequestStatus = "in_review"   // Being reviewed by admin
	StatusApproved   RequestStatus = "approved"    // Approved, waiting for assignment
	StatusAssigned   RequestStatus = "assigned"    // Assigned to someone
	StatusInProgress RequestStatus = "in_progress" // Being worked on
	StatusCompleted  RequestStatus = "completed"   // Prompt created and published
	StatusRejected   RequestStatus = "rejected"    // Request rejected
	StatusOnHold     RequestStatus = "on_hold"     // Temporarily paused
)

// Valid checks if the request status is valid
func (s RequestStatus) Valid() bool {
	switch s {
	case StatusPending, StatusInReview, StatusApproved, StatusAssigned,
		StatusInProgress, StatusCompleted, StatusRejected, StatusOnHold:
		return true
	}
	return false
}

// Priority represents the priority level of a request
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// Valid checks if the priority is valid
func (p Priority) Valid() bool {
	switch p {
	case PriorityLow, PriorityNormal, PriorityHigh, PriorityUrgent:
		return true
	}
	return false
}

// TableName specifies the table name for GORM
func (PromptRequest) TableName() string {
	return "prompt_requests"
}

// BeforeCreate hook - called before creating a record
func (pr *PromptRequest) BeforeCreate(tx *gorm.DB) error {
	// Set defaults and validate
	if !pr.Status.Valid() {
		pr.Status = StatusPending
	}
	if !pr.Priority.Valid() {
		pr.Priority = PriorityNormal
	}
	if !pr.RequestedDifficulty.Valid() {
		pr.RequestedDifficulty = DifficultyBeginner
	}

	// Set urgent priority if marked as urgent
	if pr.IsUrgent && pr.Priority == PriorityNormal {
		pr.Priority = PriorityHigh
	}

	return nil
}

// PromptRequestCreateRequest represents the public form submission
// This is what comes from the frontend form - POST /api/requests
// Similar to req.body in Express.js contact forms
type PromptRequestCreateRequest struct {
	RequesterName        string          `json:"requester_name" validate:"required,max=100"`
	RequesterEmail       string          `json:"requester_email" validate:"required,email,max=100"`
	RequestedTitle       string          `json:"requested_title" validate:"required,max=200"`
	RequestedLanguage    string          `json:"requested_language" validate:"required,max=50"`
	RequestedDifficulty  DifficultyLevel `json:"requested_difficulty" validate:"required"`
	RequestedCategory    string          `json:"requested_category" validate:"required,max=100"`
	Description          string          `json:"description" validate:"required"`
	SpecificRequirements string          `json:"specific_requirements,omitempty"`
	UseCase              string          `json:"use_case,omitempty"`
	PreferredTopics      string          `json:"preferred_topics,omitempty"`
	IsUrgent             bool            `json:"is_urgent,omitempty"`
}

// PromptRequestUpdateRequest represents admin updates to requests
// This is for PATCH /api/admin/requests/:id - only admins can use this
type PromptRequestUpdateRequest struct {
	Status            *RequestStatus `json:"status,omitempty"`
	Priority          *Priority      `json:"priority,omitempty"`
	AssignedToID      *uint          `json:"assigned_to_id,omitempty"`
	AdminNotes        *string        `json:"admin_notes,omitempty"`
	ResponseMessage   *string        `json:"response_message,omitempty"`
	EstimatedHours    *int           `json:"estimated_hours,omitempty"`
	CompletedPromptID *uint          `json:"completed_prompt_id,omitempty"`
}

// ToPromptRequest converts PromptRequestCreateRequest to PromptRequest model
// Similar to creating a new model instance from request body in Express.js
func (req *PromptRequestCreateRequest) ToPromptRequest() *PromptRequest {
	priority := PriorityNormal
	if req.IsUrgent {
		priority = PriorityHigh
	}

	return &PromptRequest{
		RequesterName:        req.RequesterName,
		RequesterEmail:       req.RequesterEmail,
		RequestedTitle:       req.RequestedTitle,
		RequestedLanguage:    req.RequestedLanguage,
		RequestedDifficulty:  req.RequestedDifficulty,
		RequestedCategory:    req.RequestedCategory,
		Description:          req.Description,
		SpecificRequirements: req.SpecificRequirements,
		UseCase:              req.UseCase,
		PreferredTopics:      req.PreferredTopics,
		IsUrgent:             req.IsUrgent,
		Priority:             priority,
		Status:               StatusPending,
	}
}

// RequestFilter represents filtering options for admin panel
// Similar to query params in Express.js: /api/admin/requests?status=pending&priority=high
type RequestFilter struct {
	Status              RequestStatus   `json:"status,omitempty"`
	Priority            Priority        `json:"priority,omitempty"`
	RequestedLanguage   string          `json:"requested_language,omitempty"`
	RequestedDifficulty DifficultyLevel `json:"requested_difficulty,omitempty"`
	RequestedCategory   string          `json:"requested_category,omitempty"`
	IsUrgent            *bool           `json:"is_urgent,omitempty"`
	IsRejected          *bool           `json:"is_rejected,omitempty"`
	AssignedToID        *uint           `json:"assigned_to_id,omitempty"`
	RequesterEmail      string          `json:"requester_email,omitempty"`
	Search              string          `json:"search,omitempty"` // Search in title/description
}

// PromptRequestResponse represents what we send back to clients
// We might hide some admin fields from public API responses
type PromptRequestResponse struct {
	ID                  uint            `json:"id"`
	RequesterName       string          `json:"requester_name"`
	RequesterEmail      string          `json:"requester_email"`
	RequestedTitle      string          `json:"requested_title"`
	RequestedLanguage   string          `json:"requested_language"`
	RequestedDifficulty DifficultyLevel `json:"requested_difficulty"`
	RequestedCategory   string          `json:"requested_category"`
	Description         string          `json:"description"`
	Status              RequestStatus   `json:"status"`
	Priority            Priority        `json:"priority"`
	CompletedPromptID   *uint           `json:"completed_prompt_id,omitempty"`
	ResponseMessage     string          `json:"response_message,omitempty"`
	CreatedAt           int64           `json:"created_at"`
	UpdatedAt           int64           `json:"updated_at"`
}

// ToResponse converts PromptRequest to PromptRequestResponse
// Similar to selecting what data to send in Express.js responses
func (pr *PromptRequest) ToResponse() *PromptRequestResponse {
	return &PromptRequestResponse{
		ID:                  pr.ID,
		RequesterName:       pr.RequesterName,
		RequesterEmail:      pr.RequesterEmail,
		RequestedTitle:      pr.RequestedTitle,
		RequestedLanguage:   pr.RequestedLanguage,
		RequestedDifficulty: pr.RequestedDifficulty,
		RequestedCategory:   pr.RequestedCategory,
		Description:         pr.Description,
		Status:              pr.Status,
		Priority:            pr.Priority,
		CompletedPromptID:   pr.CompletedPromptID,
		ResponseMessage:     pr.ResponseMessage,
		CreatedAt:           pr.CreatedAt.Unix(),
		UpdatedAt:           pr.UpdatedAt.Unix(),
	}
}
