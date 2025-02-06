package leaderboards

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
	CreateLeaderboardEvent = "leaderboards.create"
	UpdateLeaderboardEvent = "leaderboards.update"
	DeleteLeaderboardEvent = "leaderboards.delete"
)

type LeaderboardService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewLeaderboardService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *LeaderboardService {
	return &LeaderboardService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *LeaderboardService) Create(req *models.CreateLeaderboardRequest) (*models.Leaderboard, error) {
	item := &models.Leaderboard{
		Name:           req.Name,
		Type:           req.Type,
		Period:         req.Period,
		ResetFrequency: req.ResetFrequency,
		IsActive:       req.IsActive,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create leaderboard", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create leaderboard: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateLeaderboardEvent, item)

	return s.GetById(item.Id)
}

func (s *LeaderboardService) Update(id uint, req *models.UpdateLeaderboardRequest) (*models.Leaderboard, error) {
	item := &models.Leaderboard{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find leaderboard for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find leaderboard: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Period != "" {
		updates["period"] = req.Period
	}
	if req.ResetFrequency != "" {
		updates["reset_frequency"] = req.ResetFrequency
	}
	if req.IsActive != "" {
		updates["is_active"] = req.IsActive
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update leaderboard",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update leaderboard: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated leaderboard",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated leaderboard: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateLeaderboardEvent, result)

	return result, nil
}

func (s *LeaderboardService) Delete(id uint) error {
	item := &models.Leaderboard{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find leaderboard for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find leaderboard: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete leaderboard",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete leaderboard: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteLeaderboardEvent, item)

	return nil
}

func (s *LeaderboardService) GetById(id uint) (*models.Leaderboard, error) {
	item := &models.Leaderboard{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get leaderboard",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	return item, nil
}

func (s *LeaderboardService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.Leaderboard
	var total int64
	query := s.DB.Model(&models.Leaderboard{})
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
		s.Logger.Error("failed to count leaderboards",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count leaderboards: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.Leaderboard{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get leaderboards",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get leaderboards: %w", err)
	}

	// Convert to response type
	responses := make([]*models.LeaderboardListResponse, len(items))
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
