package repositories

import (
	"PromptGallery/internal/models"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type PromptRepository struct {
	db *gorm.DB
}

func NewPromptRepository(db *gorm.DB) *PromptRepository {
	return &PromptRepository{
		db: db,
	}
}

func (r *PromptRepository) FindAll(filter models.PromptFilter, page, limit int) ([]models.Prompt, int64, error) {

	var prompts []models.Prompt
	var total int64

	query := r.db.Model(&models.Prompt{})

	query = r.applyFilters(query, filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// offset pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&prompts).Error; err != nil {
		return nil, 0, err
	}

	return prompts, total, nil
}

func (r *PromptRepository) FindByID(id uint) (*models.Prompt, error) {

	var prompt models.Prompt

	if err := r.db.First(&prompt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("prompt not found")
		}
		return nil, err
	}

	return &prompt, nil
}

func (r *PromptRepository) Create(prompt *models.Prompt) (*models.Prompt, error) {
	if err := r.db.Create(prompt).Error; err != nil {
		return nil, err
	}
	return prompt, nil
}

func (r *PromptRepository) Update(prompt *models.Prompt) (*models.Prompt, error) {
	if err := r.db.Save(prompt).Error; err != nil {
		return nil, err
	}
	return prompt, nil
}

func (r *PromptRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Prompt{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("prompt not found")
	}
	return nil
}

func (r *PromptRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Prompt{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *PromptRepository) FindByLanguage(language string, limit int) ([]models.Prompt, error) {
	var prompts []models.Prompt

	err := r.db.Where("language = ?", language).
		Order("created_at DESC").
		Limit(limit).
		Find(&prompts).Error

	return prompts, err
}

func (r *PromptRepository) FindPopular(limit int) ([]models.Prompt, error) {
	var prompts []models.Prompt

	err := r.db.Order("view_count DESC").
		Limit(limit).
		Find(&prompts).Error

	return prompts, err
}

func (r *PromptRepository) FindByDifficulty(difficulty models.DifficultyLevel, limit int) ([]models.Prompt, error) {
	var prompts []models.Prompt

	err := r.db.Where("difficulty = ?", difficulty).
		Order("created_at DESC").
		Limit(limit).
		Find(&prompts).Error

	return prompts, err
}

func (r *PromptRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Prompt{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *PromptRepository) applyFilters(query *gorm.DB, filter models.PromptFilter) *gorm.DB {
	if filter.Language != "" {
		query = query.Where("language = ?", filter.Language)
	}
	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if filter.IsVerified != nil {
		query = query.Where("is_verified = ?", *filter.IsVerified)
	}

	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(problem_statement) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	return query
}
