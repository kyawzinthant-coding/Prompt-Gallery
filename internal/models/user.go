package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

// User represents authorized users who can create/edit/delete prompts
// Similar to User model in Express.js apps with roles
type User struct {
	gorm.Model

	// Basic information
	Name     string `gorm:"not null;size:100" json:"name"`
	Email    string `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Username string `gorm:"uniqueIndex;not null;size:50" json:"username"`

	// Authentication (for future - could use JWT, OAuth, etc.)
	PasswordHash string `gorm:"not null" json:"-"` // Never send in JSON responses

	// Authorization
	Role     UserRole `gorm:"not null;default:'contributor'" json:"role"`
	IsActive bool     `gorm:"default:true;index" json:"is_active"`

	// Profile information
	Bio      string `gorm:"type:text" json:"bio,omitempty"`
	Website  string `gorm:"size:200" json:"website,omitempty"`
	Avatar   string `gorm:"size:500" json:"avatar,omitempty"` // URL to avatar image
	Location string `gorm:"size:100" json:"location,omitempty"`

	// Specialties - what they're good at
	Specialties string `gorm:"type:text" json:"specialties,omitempty"` // JSON array of languages/topics

	// Social Links
	GithubUsername  string `gorm:"size:100" json:"github_username,omitempty"`
	TwitterUsername string `gorm:"size:100" json:"twitter_username,omitempty"`
	LinkedinProfile string `gorm:"size:200" json:"linkedin_profile,omitempty"`

	// Statistics
	PromptsCreated  int `gorm:"default:0" json:"prompts_created"`
	PromptsVerified int `gorm:"default:0" json:"prompts_verified"` // How many they've verified
	RequestsHandled int `gorm:"default:0" json:"requests_handled"` // How many requests they've completed
}

// UserRole represents different user roles
// Similar to role-based access control in Express.js apps
type UserRole string

const (
	RoleContributor UserRole = "contributor" // Can create and edit their own prompts
	RoleModerator   UserRole = "moderator"   // Can verify prompts, manage requests
	RoleAdmin       UserRole = "admin"       // Full access to everything
	RoleSuperAdmin  UserRole = "super_admin" // System administration
)

// Valid checks if the user role is valid
func (r UserRole) Valid() bool {
	switch r {
	case RoleContributor, RoleModerator, RoleAdmin, RoleSuperAdmin:
		return true
	}
	return false
}

// CanCreatePrompts checks if user can create prompts
func (r UserRole) CanCreatePrompts() bool {
	return r == RoleContributor || r == RoleModerator || r == RoleAdmin || r == RoleSuperAdmin
}

// CanVerifyPrompts checks if user can verify prompts
func (r UserRole) CanVerifyPrompts() bool {
	return r == RoleModerator || r == RoleAdmin || r == RoleSuperAdmin
}

// CanManageRequests checks if user can manage prompt requests
func (r UserRole) CanManageRequests() bool {
	return r == RoleModerator || r == RoleAdmin || r == RoleSuperAdmin
}

// CanManageUsers checks if user can manage other users
func (r UserRole) CanManageUsers() bool {
	return r == RoleAdmin || r == RoleSuperAdmin
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if !u.Role.Valid() {
		u.Role = RoleContributor
	}
	return nil
}

// UserCreateRequest represents request to create a new user (admin only)
// Similar to user creation in Express.js admin panels
type UserCreateRequest struct {
	Name     string   `json:"name" validate:"required,max=100"`
	Email    string   `json:"email" validate:"required,email,max=100"`
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role,omitempty"`
	Bio      string   `json:"bio,omitempty"`
	Website  string   `json:"website,omitempty" validate:"omitempty,url"`
}

// UserUpdateRequest represents updates to user profile
// This is for PATCH /api/users/:id or /api/profile
type UserUpdateRequest struct {
	Name            *string  `json:"name,omitempty" validate:"omitempty,max=100"`
	Email           *string  `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Username        *string  `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Bio             *string  `json:"bio,omitempty"`
	Website         *string  `json:"website,omitempty" validate:"omitempty,url"`
	Avatar          *string  `json:"avatar,omitempty" validate:"omitempty,url"`
	Location        *string  `json:"location,omitempty" validate:"omitempty,max=100"`
	Specialties     []string `json:"specialties,omitempty"`
	GithubUsername  *string  `json:"github_username,omitempty" validate:"omitempty,max=100"`
	TwitterUsername *string  `json:"twitter_username,omitempty" validate:"omitempty,max=100"`
	LinkedinProfile *string  `json:"linkedin_profile,omitempty" validate:"omitempty,url"`
}

// UserAdminUpdateRequest for admin-only updates (role, status, etc.)
type UserAdminUpdateRequest struct {
	Role     *UserRole `json:"role,omitempty"`
	IsActive *bool     `json:"is_active,omitempty"`
}

// UserResponse represents what we send back to clients
// Excludes sensitive information like password hash
type UserResponse struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Username        string   `json:"username"`
	Role            UserRole `json:"role"`
	IsActive        bool     `json:"is_active"`
	Bio             string   `json:"bio,omitempty"`
	Website         string   `json:"website,omitempty"`
	Avatar          string   `json:"avatar,omitempty"`
	Location        string   `json:"location,omitempty"`
	Specialties     []string `json:"specialties"`
	PromptsCreated  int      `json:"prompts_created"`
	PromptsVerified int      `json:"prompts_verified"`
	RequestsHandled int      `json:"requests_handled"`
	GithubUsername  string   `json:"github_username,omitempty"`
	TwitterUsername string   `json:"twitter_username,omitempty"`
	LinkedinProfile string   `json:"linkedin_profile,omitempty"`
	CreatedAt       int64    `json:"created_at"`
	UpdatedAt       int64    `json:"updated_at"`
}

// ToResponse converts User to UserResponse
// Similar to user serialization in Express.js - exclude sensitive data
func (u *User) ToResponse() *UserResponse {
	// Convert specialties string to slice
	var specialties []string
	if u.Specialties != "" {
		// Parse JSON string to slice
		if err := json.Unmarshal([]byte(u.Specialties), &specialties); err != nil {
			// If parsing fails, return empty slice
			specialties = []string{}
		}
	}

	return &UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Username:        u.Username,
		Role:            u.Role,
		IsActive:        u.IsActive,
		Bio:             u.Bio,
		Website:         u.Website,
		Avatar:          u.Avatar,
		Location:        u.Location,
		Specialties:     specialties,
		PromptsCreated:  u.PromptsCreated,
		PromptsVerified: u.PromptsVerified,
		RequestsHandled: u.RequestsHandled,
		GithubUsername:  u.GithubUsername,
		TwitterUsername: u.TwitterUsername,
		LinkedinProfile: u.LinkedinProfile,
		CreatedAt:       u.CreatedAt.Unix(),
		UpdatedAt:       u.UpdatedAt.Unix(),
	}
}

// SetSpecialties converts a slice of strings to JSON and sets it
func (u *User) SetSpecialties(specialties []string) error {
	data, err := json.Marshal(specialties)
	if err != nil {
		return err
	}
	u.Specialties = string(data)
	return nil
}

// GetSpecialties returns specialties as a slice
func (u *User) GetSpecialties() []string {
	var specialties []string
	if u.Specialties != "" {
		json.Unmarshal([]byte(u.Specialties), &specialties)
	}
	return specialties
}

// ToCreateRequest converts UserCreateRequest to User model
func (req *UserCreateRequest) ToUser() *User {
	return &User{
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
		Role:     req.Role,
		Bio:      req.Bio,
		Website:  req.Website,
		// PasswordHash will be set separately after hashing
	}
}
