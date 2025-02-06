package models

import (
	"base/core/storage"
	"time"

	"gorm.io/gorm"
)

// Level represents a level entity
type Level struct {
	Id          uint                `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	DeletedAt   gorm.DeletedAt      `json:"deleted_at,omitempty" gorm:"index"`
	LevelNumber int                 `json:"level_number"`
	XpRequired  int                 `json:"xp_required"`
	Title       string              `json:"title"`
	Rewards     string              `json:"rewards"`
	Icon        *storage.Attachment `json:"icon,omitempty" gorm:"polymorphic:Model"`
}

// TableName returns the table name for the Level model
func (item *Level) TableName() string {
	return "levels"
}

// GetId returns the Id of the model
func (item *Level) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *Level) GetModelName() string {
	return "level"
}

// LevelListResponse represents the list view response
type LevelListResponse struct {
	Id          uint                `json:"id"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	LevelNumber int                 `json:"level_number"`
	XpRequired  int                 `json:"xp_required"`
	Title       string              `json:"title"`
	Rewards     string              `json:"rewards"`
	Icon        *storage.Attachment `json:"icon,omitempty"`
}

// LevelResponse represents the detailed view response
type LevelResponse struct {
	Id          uint                `json:"id"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	DeletedAt   gorm.DeletedAt      `json:"deleted_at,omitempty"`
	LevelNumber int                 `json:"level_number"`
	XpRequired  int                 `json:"xp_required"`
	Title       string              `json:"title"`
	Rewards     string              `json:"rewards"`
	Icon        *storage.Attachment `json:"icon,omitempty"`
}

// CreateLevelRequest represents the request payload for creating a Level
type CreateLevelRequest struct {
	LevelNumber int                 `json:"level_number" binding:"required"`
	XpRequired  int                 `json:"xp_required" binding:"required"`
	Title       string              `json:"title" binding:"required"`
	Rewards     string              `json:"rewards" binding:"required"`
	Icon        *storage.Attachment `json:"icon,omitempty"`
}

// UpdateLevelRequest represents the request payload for updating a Level
type UpdateLevelRequest struct {
	LevelNumber string              `json:"level_number,omitempty"`
	XpRequired  string              `json:"xp_required,omitempty"`
	Title       string              `json:"title,omitempty"`
	Rewards     string              `json:"rewards,omitempty"`
	Icon        *storage.Attachment `json:"icon,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *Level) ToListResponse() *LevelListResponse {
	if item == nil {
		return nil
	}
	return &LevelListResponse{
		Id:          item.Id,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		LevelNumber: item.LevelNumber,
		XpRequired:  item.XpRequired,
		Title:       item.Title,
		Rewards:     item.Rewards,
		Icon:        item.Icon,
	}
}

// ToResponse converts the model to a detailed response
func (item *Level) ToResponse() *LevelResponse {
	if item == nil {
		return nil
	}
	return &LevelResponse{
		Id:          item.Id,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		DeletedAt:   item.DeletedAt,
		LevelNumber: item.LevelNumber,
		XpRequired:  item.XpRequired,
		Title:       item.Title,
		Rewards:     item.Rewards,
		Icon:        item.Icon,
	}
}

// Preload preloads all the model's relationships
func (item *Level) Preload(db *gorm.DB) *gorm.DB {
	query := db
	return query
}
