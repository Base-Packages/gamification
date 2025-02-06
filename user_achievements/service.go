package user_achievements

import (
	"fmt"
	"math"

	"base/core/emitter"
	"base/core/logger"
	"base/core/storage"
	"base/core/types"
	"base/packages/gamification/models"

	"gorm.io/gorm"
)

const (
	CreateUserAchievementEvent = "userachievements.create"
	UpdateUserAchievementEvent = "userachievements.update"
	DeleteUserAchievementEvent = "userachievements.delete"
)

type UserAchievementService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewUserAchievementService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *UserAchievementService {
	return &UserAchievementService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *UserAchievementService) Create(req *models.CreateUserAchievementRequest) (*models.UserAchievement, error) {
	item := &models.UserAchievement{
		UserId:        req.UserId,
		AchievementId: req.AchievementId,
		Progress:      req.Progress,
		CompletedAt:   req.CompletedAt,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create userachievement", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create userachievement: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateUserAchievementEvent, item)

	return s.GetById(item.Id)
}

func (s *UserAchievementService) Update(id uint, req *models.UpdateUserAchievementRequest) (*models.UserAchievement, error) {
	item := &models.UserAchievement{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find userachievement for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find userachievement: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.UserId != 0 {
		updates["user_id"] = req.UserId
	}
	if req.AchievementId != 0 {
		updates["achievement_id"] = req.AchievementId
	}
	if req.Progress != "" {
		updates["progress"] = req.Progress
	}
	if req.CompletedAt != "" {
		updates["completed_at"] = req.CompletedAt
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update userachievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update userachievement: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated userachievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated userachievement: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateUserAchievementEvent, result)

	return result, nil
}

func (s *UserAchievementService) Delete(id uint) error {
	item := &models.UserAchievement{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find userachievement for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find userachievement: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete userachievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete userachievement: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteUserAchievementEvent, item)

	return nil
}

func (s *UserAchievementService) GetById(id uint) (*models.UserAchievement, error) {
	item := &models.UserAchievement{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get userachievement",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get userachievement: %w", err)
	}

	return item, nil
}

func (s *UserAchievementService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.UserAchievement
	var total int64
	query := s.DB.Model(&models.UserAchievement{})
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
		s.Logger.Error("failed to count userachievements",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count userachievements: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.UserAchievement{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get userachievements",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get userachievements: %w", err)
	}

	// Convert to response type
	responses := make([]*models.UserAchievementListResponse, len(items))
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
