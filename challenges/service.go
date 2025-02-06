package challenges

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
	CreateChallengeEvent = "challenges.create"
	UpdateChallengeEvent = "challenges.update"
	DeleteChallengeEvent = "challenges.delete"
)

type ChallengeService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewChallengeService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *ChallengeService {
	return &ChallengeService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *ChallengeService) Create(req *models.CreateChallengeRequest) (*models.Challenge, error) {
	item := &models.Challenge{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		RewardType:  req.RewardType,
		RewardValue: req.RewardValue,
		IsActive:    req.IsActive,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create challenge", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create challenge: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateChallengeEvent, item)

	return s.GetById(item.Id)
}

func (s *ChallengeService) Update(id uint, req *models.UpdateChallengeRequest) (*models.Challenge, error) {
	item := &models.Challenge{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find challenge for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find challenge: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.StartDate != "" {
		updates["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		updates["end_date"] = req.EndDate
	}
	if req.RewardType != "" {
		updates["reward_type"] = req.RewardType
	}
	if req.RewardValue != "" {
		updates["reward_value"] = req.RewardValue
	}
	if req.IsActive != "" {
		updates["is_active"] = req.IsActive
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update challenge",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update challenge: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated challenge",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated challenge: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateChallengeEvent, result)

	return result, nil
}

func (s *ChallengeService) Delete(id uint) error {
	item := &models.Challenge{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find challenge for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find challenge: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete challenge",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete challenge: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteChallengeEvent, item)

	return nil
}

func (s *ChallengeService) GetById(id uint) (*models.Challenge, error) {
	item := &models.Challenge{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get challenge",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}

	return item, nil
}

func (s *ChallengeService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.Challenge
	var total int64
	query := s.DB.Model(&models.Challenge{})
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
		s.Logger.Error("failed to count challenges",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count challenges: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.Challenge{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get challenges",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get challenges: %w", err)
	}

	// Convert to response type
	responses := make([]*models.ChallengeListResponse, len(items))
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
