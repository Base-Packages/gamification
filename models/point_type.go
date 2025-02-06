package models

import (
	"time"

	"gorm.io/gorm"
)

// PointType represents a pointtype entity
type PointType struct {
	Id          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
}

// TableName returns the table name for the PointType model
func (item *PointType) TableName() string {
	return "pointtypes"
}

// GetId returns the Id of the model
func (item *PointType) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *PointType) GetModelName() string {
	return "pointtype"
}

// PointTypeListResponse represents the list view response
type PointTypeListResponse struct {
	Id          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
}

// PointTypeResponse represents the detailed view response
type PointTypeResponse struct {
	Id          uint           `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
}

// CreatePointTypeRequest represents the request payload for creating a PointType
type CreatePointTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Icon        string `json:"icon" binding:"required"`
}

// UpdatePointTypeRequest represents the request payload for updating a PointType
type UpdatePointTypeRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *PointType) ToListResponse() *PointTypeListResponse {
	if item == nil {
		return nil
	}
	return &PointTypeListResponse{
		Id:          item.Id,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		Name:        item.Name,
		Description: item.Description,
		Icon:        item.Icon,
	}
}

// ToResponse converts the model to a detailed response
func (item *PointType) ToResponse() *PointTypeResponse {
	if item == nil {
		return nil
	}
	return &PointTypeResponse{
		Id:          item.Id,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		DeletedAt:   item.DeletedAt,
		Name:        item.Name,
		Description: item.Description,
		Icon:        item.Icon,
	}
}

// Preload preloads all the model's relationships
func (item *PointType) Preload(db *gorm.DB) *gorm.DB {
	query := db
	return query
}
