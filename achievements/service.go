package achievements

import (
	"fmt"
	"math"
	"mime/multipart"

	"base/core/emitter"
	"base/core/logger"
	"base/core/storage"
	"base/core/types"
	"base/packages/gamification/models"

	"gorm.io/gorm"
)

const (
	CreateAchievementEvent = "achievements.create"
	UpdateAchievementEvent = "achievements.update"
	DeleteAchievementEvent = "achievements.delete"
)

type AchievementService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewAchievementService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *AchievementService {
	return &AchievementService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *AchievementService) Create(req *models.CreateAchievementRequest) (*models.Achievement, error) {
	item := &models.Achievement{
		Name:        req.Name,
		Description: req.Description,
		// Icon attachment is handled via separate endpoint
		Category:        req.Category,
		DifficultyLevel: req.DifficultyLevel,
		IsHidden:        req.IsHidden,
		IsActive:        req.IsActive,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create achievement", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create achievement: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateAchievementEvent, item)

	return s.GetById(item.Id)
}

func (s *AchievementService) Update(id uint, req *models.UpdateAchievementRequest) (*models.Achievement, error) {
	item := &models.Achievement{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find achievement for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find achievement: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	// Icon attachment is handled via separate endpoint
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.DifficultyLevel != "" {
		updates["difficulty_level"] = req.DifficultyLevel
	}
	if req.IsHidden != "" {
		updates["is_hidden"] = req.IsHidden
	}
	if req.IsActive != "" {
		updates["is_active"] = req.IsActive
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update achievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update achievement: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated achievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated achievement: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateAchievementEvent, result)

	return result, nil
}

func (s *AchievementService) Delete(id uint) error {
	item := &models.Achievement{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find achievement for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find achievement: %w", err)
	}

	// Delete file attachments if any
	if item.Icon != nil {
		if err := s.Storage.Delete(item.Icon); err != nil {
			s.Logger.Error("failed to delete icon",
				logger.String("error", err.Error()),
				logger.Int("id", int(id)))
			return fmt.Errorf("failed to delete icon: %w", err)
		}
	}

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete achievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete achievement: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteAchievementEvent, item)

	return nil
}

func (s *AchievementService) GetById(id uint) (*models.Achievement, error) {
	item := &models.Achievement{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get achievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get achievement: %w", err)
	}

	return item, nil
}

func (s *AchievementService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.Achievement
	var total int64
	query := s.DB.Model(&models.Achievement{})
	// Set default values if nil
	defaultPage := 1
	defaultLimit := 10
	if page == nil {
		page = &defaultPage
	}
	if limit == nil {
		limit = &defaultLimit
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		s.Logger.Error("failed to count achievements",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count achievements: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.Achievement{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get achievements",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get achievements: %w", err)
	}

	// Convert to response type
	responses := make([]*models.AchievementListResponse, len(items))
	for i, item := range items {
		responses[i] = item.ToListResponse()
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(*limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &types.PaginatedResponse{
		Data: responses,
		Pagination: types.Pagination{
			Total:      int(total),
			Page:       *page,
			PageSize:   *limit,
			TotalPages: totalPages,
		},
	}, nil
}

// UploadIcon uploads a file for the Achievement's Icon field
func (s *AchievementService) UploadIcon(id uint, file *multipart.FileHeader) (*models.Achievement, error) {
	item := &models.Achievement{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find achievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find achievement: %w", err)
	}

	// Delete existing file if any
	if item.Icon != nil {
		if err := s.Storage.Delete(item.Icon); err != nil {
			s.Logger.Error("failed to delete existing icon",
				logger.String("error", err.Error()),
				logger.Int("id", int(id)))
			return nil, fmt.Errorf("failed to delete existing icon: %w", err)
		}
	}

	// Attach new file
	attachment, err := s.Storage.Attach(item, "icon", file)
	if err != nil {
		s.Logger.Error("failed to attach icon",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to attach icon: %w", err)
	}

	// Update the model with the new attachment
	if err := s.DB.Model(item).Association("Icon").Replace(attachment); err != nil {
		s.Logger.Error("failed to associate icon",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to associate icon: %w", err)
	}

	return s.GetById(id)
}

// RemoveIcon removes the file from the Achievement's Icon field
func (s *AchievementService) RemoveIcon(id uint) (*models.Achievement, error) {
	item := &models.Achievement{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find achievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find achievement: %w", err)
	}

	if item.Icon == nil {
		return item, nil
	}

	if err := s.Storage.Delete(item.Icon); err != nil {
		s.Logger.Error("failed to delete icon",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to delete icon: %w", err)
	}

	// Clear the association
	if err := s.DB.Model(item).Association("Icon").Clear(); err != nil {
		s.Logger.Error("failed to clear icon association",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to clear icon association: %w", err)
	}

	return s.GetById(id)
}
