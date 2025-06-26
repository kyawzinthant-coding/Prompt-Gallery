package services

import (
	"PromptGallery/internal/models"
	"PromptGallery/internal/repositories"
	"errors"
	"fmt"
)

type PromptService struct {
	promptRepo *repositories.PromptRepository
}

func NewPromptService(promptRepo *repositories.PromptRepository) *PromptService {
	return &PromptService{
		promptRepo: promptRepo,
	}
}

// shape the response data
type PromptResponse struct {
	ID               uint                   `json:"id"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Language         string                 `json:"language"`
	Difficulty       models.DifficultyLevel `json:"difficulty"`
	Category         string                 `json:"category"`
	ProblemStatement string                 `json:"problem_statement"`
	IsVerified       bool                   `json:"is_verified"`
	ViewCount        int                    `json:"view_count"`
	LikeCount        int                    `json:"like_count"`
	Tags             string                 `json:"tags"`
	AuthorName       string                 `json:"author_name,omitempty"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

type PaginationPromptResponse struct {
	Data       []PromptResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

func (s *PromptService) GetAllPrompts(filter models.PromptFilter, page, limit int) (*PaginationPromptResponse, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	if filter.Difficulty != "" && !filter.Difficulty.Valid() {
		return nil, errors.New("invalid difficulty")
	}

	prompts, total, err := s.promptRepo.FindAll(filter, page, limit)
	if err != nil {
		return nil, err
	}

	promptResponses := make([]PromptResponse, len(prompts))
	for i, prompt := range prompts {
		promptResponses[i] = s.transformToResponse(&prompt)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &PaginationPromptResponse{
		Data:       promptResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil

}

func (s *PromptService) GetPromptByID(id uint) (*PromptResponse, error) {
	if id == 0 {
		return nil, errors.New("invalid prompt id")
	}

	prompt, err := s.promptRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find prompt: %w", err)
	}

	go func() {
		_ = s.promptRepo.IncrementViewCount(id)
	}()

	response := s.transformToResponse(prompt)
	return &response, nil
}

func (s *PromptService) CreatePrompt(createReq *models.PromptCreateRequest) (*PromptResponse, error) {
	if err := s.validateCreateRequest(createReq); err != nil {
		return nil, err
	}

	prompt := createReq.ToPrompt()

	if prompt.Difficulty == "" {
		prompt.Difficulty = models.DifficultyBeginner
	}

	createdPrompt, err := s.promptRepo.Create(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt: %w", err)
	}

	response := s.transformToResponse(createdPrompt)
	return &response, nil
}

func (s *PromptService) DeletePrompt(id uint) error {
	if id == 0 {
		return errors.New("invalid prompt id")
	}

	exists, err := s.promptRepo.Exists(id)

	if err != nil {
		return fmt.Errorf("failed to check if prompt exists: %w", err)
	}
	if !exists {
		return errors.New("prompt not found")
	}

	err = s.promptRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete prompt: %w", err)
	}
	return nil
}

func (s *PromptService) GetPopularPrompts(limit int) ([]PromptResponse, error) {
	// Business logic - validate limit
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Call repository
	prompts, err := s.promptRepo.FindPopular(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular prompts: %w", err)
	}

	// Transform results
	responses := make([]PromptResponse, len(prompts))
	for i, prompt := range prompts {
		responses[i] = s.transformToResponse(&prompt)
	}

	return responses, nil
}

func (s *PromptService) validateCreateRequest(req *models.PromptCreateRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if req.Language == "" {
		return errors.New("language is required")
	}
	if req.Category == "" {
		return errors.New("category is required")
	}
	if req.ProblemStatement == "" {
		return errors.New("problem statement is required")
	}
	if req.Difficulty != "" && !req.Difficulty.Valid() {
		return errors.New("invalid difficulty level")
	}

	return nil
}

func (s *PromptService) transformToResponse(prompt *models.Prompt) PromptResponse {
	return PromptResponse{
		ID:               prompt.ID,
		Title:            prompt.Title,
		Description:      prompt.Description,
		Language:         prompt.Language,
		Difficulty:       prompt.Difficulty,
		Category:         prompt.Category,
		ProblemStatement: prompt.ProblemStatement,
		IsVerified:       prompt.IsVerified,
		ViewCount:        prompt.ViewCount,
		LikeCount:        prompt.LikeCount,
		Tags:             prompt.Tags,
		AuthorName:       prompt.AuthorName,
		CreatedAt:        prompt.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        prompt.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
