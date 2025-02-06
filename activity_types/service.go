package activity_types

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
	CreateActivityTypeEvent = "activitytypes.create"
	UpdateActivityTypeEvent = "activitytypes.update"
	DeleteActivityTypeEvent = "activitytypes.delete"
)

type ActivityTypeService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewActivityTypeService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *ActivityTypeService {
	return &ActivityTypeService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *ActivityTypeService) Create(req *models.CreateActivityTypeRequest) (*models.ActivityType, error) {
	item := &models.ActivityType{
		Name:           req.Name,
		Description:    req.Description,
		Category:       req.Category,
		PointsValue:    req.PointsValue,
		CooldownPeriod: req.CooldownPeriod,
		IsActive:       req.IsActive,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create activitytype", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create activitytype: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateActivityTypeEvent, item)

	return s.GetById(item.Id)
}

func (s *ActivityTypeService) Update(id uint, req *models.UpdateActivityTypeRequest) (*models.ActivityType, error) {
	item := &models.ActivityType{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find activitytype for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find activitytype: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.PointsValue != "" {
		updates["points_value"] = req.PointsValue
	}
	if req.CooldownPeriod != "" {
		updates["cooldown_period"] = req.CooldownPeriod
	}
	if req.IsActive != "" {
		updates["is_active"] = req.IsActive
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update activitytype",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update activitytype: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated activitytype",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated activitytype: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateActivityTypeEvent, result)

	return result, nil
}

func (s *ActivityTypeService) Delete(id uint) error {
	item := &models.ActivityType{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find activitytype for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find activitytype: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete activitytype",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete activitytype: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteActivityTypeEvent, item)

	return nil
}

func (s *ActivityTypeService) GetById(id uint) (*models.ActivityType, error) {
	item := &models.ActivityType{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get activitytype",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get activitytype: %w", err)
	}

	return item, nil
}

func (s *ActivityTypeService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.ActivityType
	var total int64
	query := s.DB.Model(&models.ActivityType{})
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
		s.Logger.Error("failed to count activitytypes",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count activitytypes: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.ActivityType{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get activitytypes",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get activitytypes: %w", err)
	}

	// Convert to response type
	responses := make([]*models.ActivityTypeListResponse, len(items))
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
