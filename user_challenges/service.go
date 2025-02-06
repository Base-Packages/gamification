package user_challenges

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
	CreateUserChallengeEvent = "userchallenges.create"
	UpdateUserChallengeEvent = "userchallenges.update"
	DeleteUserChallengeEvent = "userchallenges.delete"
)

type UserChallengeService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewUserChallengeService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *UserChallengeService {
	return &UserChallengeService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *UserChallengeService) Create(req *models.CreateUserChallengeRequest) (*models.UserChallenge, error) {
	item := &models.UserChallenge{
		UserId:        req.UserId,
		ChallengeId:   req.ChallengeId,
		Progress:      req.Progress,
		CompletedAt:   req.CompletedAt,
		RewardClaimed: req.RewardClaimed,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create userchallenge", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create userchallenge: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateUserChallengeEvent, item)

	return s.GetById(item.Id)
}

func (s *UserChallengeService) Update(id uint, req *models.UpdateUserChallengeRequest) (*models.UserChallenge, error) {
	item := &models.UserChallenge{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find userchallenge for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find userchallenge: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.UserId != 0 {
		updates["user_id"] = req.UserId
	}
	if req.ChallengeId != 0 {
		updates["challenge_id"] = req.ChallengeId
	}
	if req.Progress != "" {
		updates["progress"] = req.Progress
	}
	if req.CompletedAt != "" {
		updates["completed_at"] = req.CompletedAt
	}
	if req.RewardClaimed != "" {
		updates["reward_claimed"] = req.RewardClaimed
	}

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update userchallenge",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update userchallenge: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated userchallenge",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated userchallenge: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateUserChallengeEvent, result)

	return result, nil
}

func (s *UserChallengeService) Delete(id uint) error {
	item := &models.UserChallenge{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find userchallenge for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find userchallenge: %w", err)
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete userchallenge",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete userchallenge: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteUserChallengeEvent, item)

	return nil
}

func (s *UserChallengeService) GetById(id uint) (*models.UserChallenge, error) {
	item := &models.UserChallenge{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get userchallenge",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get userchallenge: %w", err)
	}

	return item, nil
}

func (s *UserChallengeService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.UserChallenge
	var total int64
	query := s.DB.Model(&models.UserChallenge{})
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
		s.Logger.Error("failed to count userchallenges",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count userchallenges: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.UserChallenge{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get userchallenges",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get userchallenges: %w", err)
	}

	// Convert to response type
	responses := make([]*models.UserChallengeListResponse, len(items))
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
