package achievement_criteria

import (
	"fmt"
	"math"

	"base/core/emitter"
	"base/core/logger"
	"base/core/packages/gamification/models"
	"base/core/storage"
	"base/core/types"

	"gorm.io/gorm"
)

const (
	CreateAchievementCriteriaEvent = "achievementcriteria.create"
	UpdateAchievementCriteriaEvent = "achievementcriteria.update"
	DeleteAchievementCriteriaEvent = "achievementcriteria.delete"
)

type AchievementCriteriaService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewAchievementCriteriaService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *AchievementCriteriaService {
	return &AchievementCriteriaService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *AchievementCriteriaService) Create(req *models.CreateAchievementCriteriaRequest) (*models.AchievementCriteria, error) {
	item := &models.AchievementCriteria{
		AchievementId:  req.AchievementId,
		ActivityTypeId: req.ActivityTypeId,
		RequiredCount:  req.RequiredCount,
		TimeFrame:      req.TimeFrame,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create achievementcriteria", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create achievementcriteria: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateAchievementCriteriaEvent, item)

	return s.GetById(item.Id)
}

func (s *AchievementCriteriaService) Update(id uint, req *models.UpdateAchievementCriteriaRequest) (*models.AchievementCriteria, error) {
	item := &models.AchievementCriteria{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find achievementcriteria for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find achievementcriteria: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.AchievementId != 0 {
		updates["achievement_id"] = req.AchievementId
	}
	if req.ActivityTypeId != 0 {
		updates["activity_type_id"] = req.ActivityTypeId
	}
	if req.RequiredCount != "" {
		updates["required_count"] = req.RequiredCount
	}
	if req.TimeFrame != "" {
		updates["time_frame"] = req.TimeFrame
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update achievementcriteria",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update achievementcriteria: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated achievementcriteria",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated achievementcriteria: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateAchievementCriteriaEvent, result)

	return result, nil
}

func (s *AchievementCriteriaService) Delete(id uint) error {
	item := &models.AchievementCriteria{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find achievementcriteria for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find achievementcriteria: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete achievementcriteria",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete achievementcriteria: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteAchievementCriteriaEvent, item)

	return nil
}

func (s *AchievementCriteriaService) GetById(id uint) (*models.AchievementCriteria, error) {
	item := &models.AchievementCriteria{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get achievementcriteria",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get achievementcriteria: %w", err)
	}

	return item, nil
}

func (s *AchievementCriteriaService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.AchievementCriteria
	var total int64
	query := s.DB.Model(&models.AchievementCriteria{})
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
		s.Logger.Error("failed to count achievementcriteria",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count achievementcriteria: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.AchievementCriteria{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get achievementcriteria",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get achievementcriteria: %w", err)
	}

	// Convert to response type
	responses := make([]*models.AchievementCriteriaListResponse, len(items))
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
