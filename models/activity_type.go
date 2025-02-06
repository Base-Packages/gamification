package models

import (
	"time"

	"gorm.io/gorm"
)

// ActivityType represents a activitytype entity
type ActivityType struct {
	Id             uint           `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Category       string         `json:"category"`
	PointsValue    int            `json:"points_value"`
	CooldownPeriod int            `json:"cooldown_period"`
	IsActive       bool           `json:"is_active"`
}

// TableName returns the table name for the ActivityType model
func (item *ActivityType) TableName() string {
	return "activitytypes"
}

// GetId returns the Id of the model
func (item *ActivityType) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *ActivityType) GetModelName() string {
	return "activitytype"
}

// ActivityTypeListResponse represents the list view response
type ActivityTypeListResponse struct {
	Id             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	PointsValue    int       `json:"points_value"`
	CooldownPeriod int       `json:"cooldown_period"`
	IsActive       bool      `json:"is_active"`
}

// ActivityTypeResponse represents the detailed view response
type ActivityTypeResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Category       string         `json:"category"`
	PointsValue    int            `json:"points_value"`
	CooldownPeriod int            `json:"cooldown_period"`
	IsActive       bool           `json:"is_active"`
}

// CreateActivityTypeRequest represents the request payload for creating a ActivityType
type CreateActivityTypeRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"required"`
	Category       string `json:"category" binding:"required"`
	PointsValue    int    `json:"points_value" binding:"required"`
	CooldownPeriod int    `json:"cooldown_period" binding:"required"`
	IsActive       bool   `json:"is_active" binding:"required"`
}

// UpdateActivityTypeRequest represents the request payload for updating a ActivityType
type UpdateActivityTypeRequest struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Category       string `json:"category,omitempty"`
	PointsValue    string `json:"points_value,omitempty"`
	CooldownPeriod string `json:"cooldown_period,omitempty"`
	IsActive       string `json:"is_active,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *ActivityType) ToListResponse() *ActivityTypeListResponse {
	if item == nil {
		return nil
	}
	return &ActivityTypeListResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		Name:           item.Name,
		Description:    item.Description,
		Category:       item.Category,
		PointsValue:    item.PointsValue,
		CooldownPeriod: item.CooldownPeriod,
		IsActive:       item.IsActive,
	}
}

// ToResponse converts the model to a detailed response
func (item *ActivityType) ToResponse() *ActivityTypeResponse {
	if item == nil {
		return nil
	}
	return &ActivityTypeResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt,
		Name:           item.Name,
		Description:    item.Description,
		Category:       item.Category,
		PointsValue:    item.PointsValue,
		CooldownPeriod: item.CooldownPeriod,
		IsActive:       item.IsActive,
	}
}

// Preload preloads all the model's relationships
func (item *ActivityType) Preload(db *gorm.DB) *gorm.DB {
	query := db
	return query
}
