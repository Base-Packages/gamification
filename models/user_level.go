package models

import (
	"base/core/app/users"
	"base/core/types"
	"time"

	"gorm.io/gorm"
)

// UserLevel represents a userlevel entity
type UserLevel struct {
	Id             uint           `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	UserId         uint           `json:"user_id"`
	User           *users.User    `json:"user,omitempty"`
	CurrentLevelId uint           `json:"current_level_id"`
	CurrentLevel   *Level         `json:"current_level,omitempty"`
	CurrentXp      int            `json:"current_xp"`
	LastLeveledUp  types.DateTime `json:"last_leveled_up"`
}

// TableName returns the table name for the UserLevel model
func (item *UserLevel) TableName() string {
	return "userlevels"
}

// GetId returns the Id of the model
func (item *UserLevel) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *UserLevel) GetModelName() string {
	return "userlevel"
}

// UserLevelListResponse represents the list view response
type UserLevelListResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UserId         uint           `json:"user_id"`
	CurrentLevelId uint           `json:"current_level_id"`
	CurrentXp      int            `json:"current_xp"`
	LastLeveledUp  types.DateTime `json:"last_leveled_up"`
}

// UserLevelResponse represents the detailed view response
type UserLevelResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty"`
	UserId         uint           `json:"user_id"`
	User           *users.User    `json:"user,omitempty"`
	CurrentLevelId uint           `json:"current_level_id"`
	CurrentLevel   *Level         `json:"current_level,omitempty"`
	CurrentXp      int            `json:"current_xp"`
	LastLeveledUp  types.DateTime `json:"last_leveled_up"`
}

// CreateUserLevelRequest represents the request payload for creating a UserLevel
type CreateUserLevelRequest struct {
	UserId         uint           `json:"user_id" binding:"required"`
	CurrentLevelId uint           `json:"current_level_id" binding:"required"`
	CurrentXp      int            `json:"current_xp" binding:"required"`
	LastLeveledUp  types.DateTime `json:"last_leveled_up" binding:"required"`
}

// UpdateUserLevelRequest represents the request payload for updating a UserLevel
type UpdateUserLevelRequest struct {
	UserId         uint   `json:"user_id,omitempty"`
	CurrentLevelId uint   `json:"current_level_id,omitempty"`
	CurrentXp      string `json:"current_xp,omitempty"`
	LastLeveledUp  string `json:"last_leveled_up,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *UserLevel) ToListResponse() *UserLevelListResponse {
	if item == nil {
		return nil
	}
	return &UserLevelListResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		UserId:         item.UserId,
		CurrentLevelId: item.CurrentLevelId,
		CurrentXp:      item.CurrentXp,
		LastLeveledUp:  item.LastLeveledUp,
	}
}

// ToResponse converts the model to a detailed response
func (item *UserLevel) ToResponse() *UserLevelResponse {
	if item == nil {
		return nil
	}
	return &UserLevelResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt,
		UserId:         item.UserId,
		User:           item.User,
		CurrentLevelId: item.CurrentLevelId,
		CurrentLevel:   item.CurrentLevel,
		CurrentXp:      item.CurrentXp,
		LastLeveledUp:  item.LastLeveledUp,
	}
}

// Preload preloads all the model's relationships
func (item *UserLevel) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("User")
	query = query.Preload("CurrentLevel")
	return query
}
