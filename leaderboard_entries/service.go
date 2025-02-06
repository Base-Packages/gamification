package leaderboard_entries

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
	CreateLeaderboardEntryEvent = "leaderboardentries.create"
	UpdateLeaderboardEntryEvent = "leaderboardentries.update"
	DeleteLeaderboardEntryEvent = "leaderboardentries.delete"
)

type LeaderboardEntryService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewLeaderboardEntryService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *LeaderboardEntryService {
	return &LeaderboardEntryService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *LeaderboardEntryService) Create(req *models.CreateLeaderboardEntryRequest) (*models.LeaderboardEntry, error) {
	item := &models.LeaderboardEntry{
		LeaderboardId: req.LeaderboardId,
		UserId:        req.UserId,
		Score:         req.Score,
		Rank:          req.Rank,
		PeriodStart:   req.PeriodStart,
		PeriodEnd:     req.PeriodEnd,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create leaderboardentry", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create leaderboardentry: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateLeaderboardEntryEvent, item)

	return s.GetById(item.Id)
}

func (s *LeaderboardEntryService) Update(id uint, req *models.UpdateLeaderboardEntryRequest) (*models.LeaderboardEntry, error) {
	item := &models.LeaderboardEntry{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find leaderboardentry for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find leaderboardentry: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.LeaderboardId != 0 {
		updates["leaderboard_id"] = req.LeaderboardId
	}
	if req.UserId != 0 {
		updates["user_id"] = req.UserId
	}
	if req.Score != "" {
		updates["score"] = req.Score
	}
	if req.Rank != "" {
		updates["rank"] = req.Rank
	}
	if req.PeriodStart != "" {
		updates["period_start"] = req.PeriodStart
	}
	if req.PeriodEnd != "" {
		updates["period_end"] = req.PeriodEnd
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update leaderboardentry",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update leaderboardentry: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated leaderboardentry",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated leaderboardentry: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateLeaderboardEntryEvent, result)

	return result, nil
}

func (s *LeaderboardEntryService) Delete(id uint) error {
	item := &models.LeaderboardEntry{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find leaderboardentry for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find leaderboardentry: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete leaderboardentry",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete leaderboardentry: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteLeaderboardEntryEvent, item)

	return nil
}

func (s *LeaderboardEntryService) GetById(id uint) (*models.LeaderboardEntry, error) {
	item := &models.LeaderboardEntry{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get leaderboardentry",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get leaderboardentry: %w", err)
	}

	return item, nil
}

func (s *LeaderboardEntryService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.LeaderboardEntry
	var total int64
	query := s.DB.Model(&models.LeaderboardEntry{})
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
		s.Logger.Error("failed to count leaderboardentries",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count leaderboardentries: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.LeaderboardEntry{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get leaderboardentries",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get leaderboardentries: %w", err)
	}

	// Convert to response type
	responses := make([]*models.LeaderboardEntryListResponse, len(items))
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
