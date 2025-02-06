package user_activities

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
	CreateUserActivityEvent = "useractivities.create"
	UpdateUserActivityEvent = "useractivities.update"
	DeleteUserActivityEvent = "useractivities.delete"
)

type UserActivityService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewUserActivityService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *UserActivityService {
	return &UserActivityService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *UserActivityService) Create(req *models.CreateUserActivityRequest) (*models.UserActivity, error) {
	item := &models.UserActivity{
		UserId:         req.UserId,
		ActivityTypeId: req.ActivityTypeId,
		PointsEarned:   req.PointsEarned,
		Metadata:       req.Metadata,
		CompletedAt:    req.CompletedAt,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create useractivity", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create useractivity: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateUserActivityEvent, item)

	return s.GetById(item.Id)
}

func (s *UserActivityService) Update(id uint, req *models.UpdateUserActivityRequest) (*models.UserActivity, error) {
	item := &models.UserActivity{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find useractivity for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find useractivity: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.UserId != 0 {
		updates["user_id"] = req.UserId
	}
	if req.ActivityTypeId != 0 {
		updates["activity_type_id"] = req.ActivityTypeId
	}
	if req.PointsEarned != "" {
		updates["points_earned"] = req.PointsEarned
	}
	if req.Metadata != "" {
		updates["metadata"] = req.Metadata
	}
	if req.CompletedAt != "" {
		updates["completed_at"] = req.CompletedAt
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update useractivity",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update useractivity: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated useractivity",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated useractivity: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateUserActivityEvent, result)

	return result, nil
}

func (s *UserActivityService) Delete(id uint) error {
	item := &models.UserActivity{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find useractivity for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find useractivity: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete useractivity",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete useractivity: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteUserActivityEvent, item)

	return nil
}

func (s *UserActivityService) GetById(id uint) (*models.UserActivity, error) {
	item := &models.UserActivity{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get useractivity",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get useractivity: %w", err)
	}

	return item, nil
}

func (s *UserActivityService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.UserActivity
	var total int64
	query := s.DB.Model(&models.UserActivity{})
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
		s.Logger.Error("failed to count useractivities",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count useractivities: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.UserActivity{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get useractivities",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get useractivities: %w", err)
	}

	// Convert to response type
	responses := make([]*models.UserActivityListResponse, len(items))
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
