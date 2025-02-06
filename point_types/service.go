package point_types

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
	CreatePointTypeEvent = "pointtypes.create"
	UpdatePointTypeEvent = "pointtypes.update"
	DeletePointTypeEvent = "pointtypes.delete"
)

type PointTypeService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewPointTypeService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *PointTypeService {
	return &PointTypeService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *PointTypeService) Create(req *models.CreatePointTypeRequest) (*models.PointType, error) {
	item := &models.PointType{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create pointtype", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create pointtype: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreatePointTypeEvent, item)

	return s.GetById(item.Id)
}

func (s *PointTypeService) Update(id uint, req *models.UpdatePointTypeRequest) (*models.PointType, error) {
	item := &models.PointType{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find pointtype for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find pointtype: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update pointtype",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update pointtype: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated pointtype",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated pointtype: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdatePointTypeEvent, result)

	return result, nil
}

func (s *PointTypeService) Delete(id uint) error {
	item := &models.PointType{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find pointtype for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find pointtype: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete pointtype",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete pointtype: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeletePointTypeEvent, item)

	return nil
}

func (s *PointTypeService) GetById(id uint) (*models.PointType, error) {
	item := &models.PointType{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get pointtype",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get pointtype: %w", err)
	}

	return item, nil
}

func (s *PointTypeService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.PointType
	var total int64
	query := s.DB.Model(&models.PointType{})
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
		s.Logger.Error("failed to count pointtypes",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count pointtypes: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.PointType{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get pointtypes",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get pointtypes: %w", err)
	}

	// Convert to response type
	responses := make([]*models.PointTypeListResponse, len(items))
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
