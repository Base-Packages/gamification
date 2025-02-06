package models

import (
	"base/core/storage"
	"time"

	"gorm.io/gorm"
)

// Achievement represents a achievement entity
type Achievement struct {
	Id              uint                `json:"id" gorm:"primarykey"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	DeletedAt       gorm.DeletedAt      `json:"deleted_at,omitempty" gorm:"index"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Icon            *storage.Attachment `json:"icon,omitempty" gorm:"polymorphic:Model"`
	Category        string              `json:"category"`
	DifficultyLevel int                 `json:"difficulty_level"`
	IsHidden        bool                `json:"is_hidden"`
	IsActive        bool                `json:"is_active"`
}

// TableName returns the table name for the Achievement model
func (item *Achievement) TableName() string {
	return "achievements"
}

// GetId returns the Id of the model
func (item *Achievement) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *Achievement) GetModelName() string {
	return "achievement"
}

// AchievementListResponse represents the list view response
type AchievementListResponse struct {
	Id              uint                `json:"id"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Icon            *storage.Attachment `json:"icon,omitempty"`
	Category        string              `json:"category"`
	DifficultyLevel int                 `json:"difficulty_level"`
	IsHidden        bool                `json:"is_hidden"`
	IsActive        bool                `json:"is_active"`
}

// AchievementResponse represents the detailed view response
type AchievementResponse struct {
	Id              uint                `json:"id"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	DeletedAt       gorm.DeletedAt      `json:"deleted_at,omitempty"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Icon            *storage.Attachment `json:"icon,omitempty"`
	Category        string              `json:"category"`
	DifficultyLevel int                 `json:"difficulty_level"`
	IsHidden        bool                `json:"is_hidden"`
	IsActive        bool                `json:"is_active"`
}

// CreateAchievementRequest represents the request payload for creating a Achievement
type CreateAchievementRequest struct {
	Name            string              `json:"name" binding:"required"`
	Description     string              `json:"description" binding:"required"`
	Icon            *storage.Attachment `json:"icon,omitempty"`
	Category        string              `json:"category" binding:"required"`
	DifficultyLevel int                 `json:"difficulty_level" binding:"required"`
	IsHidden        bool                `json:"is_hidden" binding:"required"`
	IsActive        bool                `json:"is_active" binding:"required"`
}

// UpdateAchievementRequest represents the request payload for updating a Achievement
type UpdateAchievementRequest struct {
	Name            string              `json:"name,omitempty"`
	Description     string              `json:"description,omitempty"`
	Icon            *storage.Attachment `json:"icon,omitempty"`
	Category        string              `json:"category,omitempty"`
	DifficultyLevel string              `json:"difficulty_level,omitempty"`
	IsHidden        string              `json:"is_hidden,omitempty"`
	IsActive        string              `json:"is_active,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *Achievement) ToListResponse() *AchievementListResponse {
	if item == nil {
		return nil
	}
	return &AchievementListResponse{
		Id:              item.Id,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
		Name:            item.Name,
		Description:     item.Description,
		Icon:            item.Icon,
		Category:        item.Category,
		DifficultyLevel: item.DifficultyLevel,
		IsHidden:        item.IsHidden,
		IsActive:        item.IsActive,
	}
}

// ToResponse converts the model to a detailed response
func (item *Achievement) ToResponse() *AchievementResponse {
	if item == nil {
		return nil
	}
	return &AchievementResponse{
		Id:              item.Id,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
		DeletedAt:       item.DeletedAt,
		Name:            item.Name,
		Description:     item.Description,
		Icon:            item.Icon,
		Category:        item.Category,
		DifficultyLevel: item.DifficultyLevel,
		IsHidden:        item.IsHidden,
		IsActive:        item.IsActive,
	}
}

// Preload preloads all the model's relationships
func (item *Achievement) Preload(db *gorm.DB) *gorm.DB {
	query := db
	return query
}
