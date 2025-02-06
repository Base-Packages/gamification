package levels

import (
	"fmt"
	"math"
	"mime/multipart"

	"base/core/emitter"
	"base/core/logger"
	"base/core/storage"
	"base/core/types"
	"base/packages/gamification/models"

	"gorm.io/gorm"
)

const (
	CreateLevelEvent = "levels.create"
	UpdateLevelEvent = "levels.update"
	DeleteLevelEvent = "levels.delete"
)

type LevelService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewLevelService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *LevelService {
	return &LevelService{
		DB:      db,
		Emitter: emitter,
		Storage: storage,
		Logger:  logger,
	}
}

func (s *LevelService) Create(req *models.CreateLevelRequest) (*models.Level, error) {
	item := &models.Level{
		LevelNumber: req.LevelNumber,
		XpRequired:  req.XpRequired,
		Title:       req.Title,
		Rewards:     req.Rewards,
		// Icon attachment is handled via separate endpoint
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create level", logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create level: %w", err)
	}

	// Emit create event
	s.Emitter.Emit(CreateLevelEvent, item)

	return s.GetById(item.Id)
}

func (s *LevelService) Update(id uint, req *models.UpdateLevelRequest) (*models.Level, error) {
	item := &models.Level{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find level for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find level: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.LevelNumber != "" {
		updates["level_number"] = req.LevelNumber
	}
	if req.XpRequired != "" {
		updates["xp_required"] = req.XpRequired
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Rewards != "" {
		updates["rewards"] = req.Rewards
	}
	// Icon attachment is handled via separate endpoint

	if err := s.DB.Model(item).Updates(updates).Error; err != nil {
		s.Logger.Error("failed to update level",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to update level: %w", err)
	}

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated level",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get updated level: %w", err)
	}

	// Emit update event
	s.Emitter.Emit(UpdateLevelEvent, result)

	return result, nil
}

func (s *LevelService) Delete(id uint) error {
	item := &models.Level{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find level for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to find level: %w", err)
	}

	// Delete file attachments if any
	if item.Icon != nil {
		if err := s.Storage.Delete(item.Icon); err != nil {
			s.Logger.Error("failed to delete icon",
				logger.String("error", err.Error()),
				logger.Int("id", int(id)))
			return fmt.Errorf("failed to delete icon: %w", err)
		}
	}

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete level",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return fmt.Errorf("failed to delete level: %w", err)
	}

	// Emit delete event
	s.Emitter.Emit(DeleteLevelEvent, item)

	return nil
}

func (s *LevelService) GetById(id uint) (*models.Level, error) {
	item := &models.Level{}

	query := item.Preload(s.DB)

	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get level",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to get level: %w", err)
	}

	return item, nil
}

func (s *LevelService) GetAll(page *int, limit *int) (*types.PaginatedResponse, error) {
	var items []*models.Level
	var total int64
	query := s.DB.Model(&models.Level{})
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
		s.Logger.Error("failed to count levels",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to count levels: %w", err)
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Preload relationships
	query = (&models.Level{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get levels",
			logger.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get levels: %w", err)
	}

	// Convert to response type
	responses := make([]*models.LevelListResponse, len(items))
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

// UploadIcon uploads a file for the Level's Icon field
func (s *LevelService) UploadIcon(id uint, file *multipart.FileHeader) (*models.Level, error) {
	item := &models.Level{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find level",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find level: %w", err)
	}

	// Delete existing file if any
	if item.Icon != nil {
		if err := s.Storage.Delete(item.Icon); err != nil {
			s.Logger.Error("failed to delete existing icon",
				logger.String("error", err.Error()),
				logger.Int("id", int(id)))
			return nil, fmt.Errorf("failed to delete existing icon: %w", err)
		}
	}

	// Attach new file
	attachment, err := s.Storage.Attach(item, "icon", file)
	if err != nil {
		s.Logger.Error("failed to attach icon",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to attach icon: %w", err)
	}

	// Update the model with the new attachment
	if err := s.DB.Model(item).Association("Icon").Replace(attachment); err != nil {
		s.Logger.Error("failed to associate icon",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to associate icon: %w", err)
	}

	return s.GetById(id)
}

// RemoveIcon removes the file from the Level's Icon field
func (s *LevelService) RemoveIcon(id uint) (*models.Level, error) {
	item := &models.Level{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find level",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to find level: %w", err)
	}

	if item.Icon == nil {
		return item, nil
	}

	if err := s.Storage.Delete(item.Icon); err != nil {
		s.Logger.Error("failed to delete icon",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to delete icon: %w", err)
	}

	// Clear the association
	if err := s.DB.Model(item).Association("Icon").Clear(); err != nil {
		s.Logger.Error("failed to clear icon association",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, fmt.Errorf("failed to clear icon association: %w", err)
	}

	return s.GetById(id)
}
