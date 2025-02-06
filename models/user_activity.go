package models

import (
	"base/core/app/users"
	"base/core/types"
	"time"

	"gorm.io/gorm"
)

// UserActivity represents a useractivity entity
type UserActivity struct {
	Id             uint           `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	UserId         uint           `json:"user_id"`
	User           *users.User    `json:"user,omitempty"`
	ActivityTypeId uint           `json:"activity_type_id"`
	ActivityType   *ActivityType  `json:"activity_type,omitempty"`
	PointsEarned   int            `json:"points_earned"`
	Metadata       string         `json:"metadata"`
	CompletedAt    types.DateTime `json:"completed_at"`
}

// TableName returns the table name for the UserActivity model
func (item *UserActivity) TableName() string {
	return "useractivities"
}

// GetId returns the Id of the model
func (item *UserActivity) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *UserActivity) GetModelName() string {
	return "useractivity"
}

// UserActivityListResponse represents the list view response
type UserActivityListResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UserId         uint           `json:"user_id"`
	ActivityTypeId uint           `json:"activity_type_id"`
	PointsEarned   int            `json:"points_earned"`
	Metadata       string         `json:"metadata"`
	CompletedAt    types.DateTime `json:"completed_at"`
}

// UserActivityResponse represents the detailed view response
type UserActivityResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty"`
	UserId         uint           `json:"user_id"`
	User           *users.User    `json:"user,omitempty"`
	ActivityTypeId uint           `json:"activity_type_id"`
	ActivityType   *ActivityType  `json:"activity_type,omitempty"`
	PointsEarned   int            `json:"points_earned"`
	Metadata       string         `json:"metadata"`
	CompletedAt    types.DateTime `json:"completed_at"`
}

// CreateUserActivityRequest represents the request payload for creating a UserActivity
type CreateUserActivityRequest struct {
	UserId         uint           `json:"user_id" binding:"required"`
	ActivityTypeId uint           `json:"activity_type_id" binding:"required"`
	PointsEarned   int            `json:"points_earned" binding:"required"`
	Metadata       string         `json:"metadata" binding:"required"`
	CompletedAt    types.DateTime `json:"completed_at" binding:"required"`
}

// UpdateUserActivityRequest represents the request payload for updating a UserActivity
type UpdateUserActivityRequest struct {
	UserId         uint   `json:"user_id,omitempty"`
	ActivityTypeId uint   `json:"activity_type_id,omitempty"`
	PointsEarned   string `json:"points_earned,omitempty"`
	Metadata       string `json:"metadata,omitempty"`
	CompletedAt    string `json:"completed_at,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *UserActivity) ToListResponse() *UserActivityListResponse {
	if item == nil {
		return nil
	}
	return &UserActivityListResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		UserId:         item.UserId,
		ActivityTypeId: item.ActivityTypeId,
		PointsEarned:   item.PointsEarned,
		Metadata:       item.Metadata,
		CompletedAt:    item.CompletedAt,
	}
}

// ToResponse converts the model to a detailed response
func (item *UserActivity) ToResponse() *UserActivityResponse {
	if item == nil {
		return nil
	}
	return &UserActivityResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt,
		UserId:         item.UserId,
		User:           item.User,
		ActivityTypeId: item.ActivityTypeId,
		ActivityType:   item.ActivityType,
		PointsEarned:   item.PointsEarned,
		Metadata:       item.Metadata,
		CompletedAt:    item.CompletedAt,
	}
}

// Preload preloads all the model's relationships
func (item *UserActivity) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("User")
	query = query.Preload("ActivityType")
	return query
}
