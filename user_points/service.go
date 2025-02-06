package user_points

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
	CreateUserPointEvent = "userpoints.create"
	UpdateUserPointEvent = "userpoints.update"
	DeleteUserPointEvent = "userpoints.delete"
)

type UserPointService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewUserPointService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *UserPointService {
	return &UserPointService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *UserPointService) Create(req *models.CreateUserPointRequest) (*models.UserPoint, error) {
	item := &models.UserPoint{
		UserId:         req.UserId,
		PointTypeId:    req.PointTypeId,
		CurrentBalance: req.CurrentBalance,
		LifetimeEarned: req.LifetimeEarned,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create userpoint", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create userpoint: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateUserPointEvent, item)

	return s.GetById(item.Id)
}

func (s *UserPointService) Update(id uint, req *models.UpdateUserPointRequest) (*models.UserPoint, error) {
	item := &models.UserPoint{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find userpoint for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find userpoint: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.UserId != 0 {
		updates["user_id"] = req.UserId
	}
	if req.PointTypeId != 0 {
		updates["point_type_id"] = req.PointTypeId
	}
	if req.CurrentBalance != "" {
		updates["current_balance"] = req.CurrentBalance
	}
	if req.LifetimeEarned != "" {
		updates["lifetime_earned"] = req.LifetimeEarned
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update userpoint",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update userpoint: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated userpoint",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated userpoint: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateUserPointEvent, result)

	return result, nil
}

func (s *UserPointService) Delete(id uint) error {
	item := &models.UserPoint{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find userpoint for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find userpoint: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete userpoint",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete userpoint: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteUserPointEvent, item)

	return nil
}

func (s *UserPointService) GetById(id uint) (*models.UserPoint, error) {
	item := &models.UserPoint{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get userpoint",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get userpoint: %w", err)
	}

	return item, nil
}

func (s *UserPointService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.UserPoint
	var total int64
	query := s.DB.Model(&models.UserPoint{})
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
		s.Logger.Error("failed to count userpoints",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count userpoints: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.UserPoint{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get userpoints",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get userpoints: %w", err)
	}

	// Convert to response type
	responses := make([]*models.UserPointListResponse, len(items))
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
