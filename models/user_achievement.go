package models

import (
	"base/core/app/users"
	"base/core/types"
	"time"

	"gorm.io/gorm"
)

// UserAchievement represents a userachievement entity
type UserAchievement struct {
	Id            uint           `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	UserId        uint           `json:"user_id"`
	User          *users.User    `json:"user,omitempty"`
	AchievementId uint           `json:"achievement_id"`
	Achievement   *Achievement   `json:"achievement,omitempty"`
	Progress      int            `json:"progress"`
	CompletedAt   types.DateTime `json:"completed_at"`
}

// TableName returns the table name for the UserAchievement model
func (item *UserAchievement) TableName() string {
	return "userachievements"
}

// GetId returns the Id of the model
func (item *UserAchievement) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *UserAchievement) GetModelName() string {
	return "userachievement"
}

// UserAchievementListResponse represents the list view response
type UserAchievementListResponse struct {
	Id            uint           `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	UserId        uint           `json:"user_id"`
	AchievementId uint           `json:"achievement_id"`
	Progress      int            `json:"progress"`
	CompletedAt   types.DateTime `json:"completed_at"`
}

// UserAchievementResponse represents the detailed view response
type UserAchievementResponse struct {
	Id            uint           `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty"`
	UserId        uint           `json:"user_id"`
	User          *users.User    `json:"user,omitempty"`
	AchievementId uint           `json:"achievement_id"`
	Achievement   *Achievement   `json:"achievement,omitempty"`
	Progress      int            `json:"progress"`
	CompletedAt   types.DateTime `json:"completed_at"`
}

// CreateUserAchievementRequest represents the request payload for creating a UserAchievement
type CreateUserAchievementRequest struct {
	UserId        uint           `json:"user_id" binding:"required"`
	AchievementId uint           `json:"achievement_id" binding:"required"`
	Progress      int            `json:"progress" binding:"required"`
	CompletedAt   types.DateTime `json:"completed_at" binding:"required"`
}

// UpdateUserAchievementRequest represents the request payload for updating a UserAchievement
type UpdateUserAchievementRequest struct {
	UserId        uint   `json:"user_id,omitempty"`
	AchievementId uint   `json:"achievement_id,omitempty"`
	Progress      string `json:"progress,omitempty"`
	CompletedAt   string `json:"completed_at,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *UserAchievement) ToListResponse() *UserAchievementListResponse {
	if item == nil {
		return nil
	}
	return &UserAchievementListResponse{
		Id:            item.Id,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		UserId:        item.UserId,
		AchievementId: item.AchievementId,
		Progress:      item.Progress,
		CompletedAt:   item.CompletedAt,
	}
}

// ToResponse converts the model to a detailed response
func (item *UserAchievement) ToResponse() *UserAchievementResponse {
	if item == nil {
		return nil
	}
	return &UserAchievementResponse{
		Id:            item.Id,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		DeletedAt:     item.DeletedAt,
		UserId:        item.UserId,
		User:          item.User,
		AchievementId: item.AchievementId,
		Achievement:   item.Achievement,
		Progress:      item.Progress,
		CompletedAt:   item.CompletedAt,
	}
}

// Preload preloads all the model's relationships
func (item *UserAchievement) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("User")
	query = query.Preload("Achievement")
	return query
}
