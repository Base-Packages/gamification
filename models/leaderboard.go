package models

import (
	"time"

	"gorm.io/gorm"
)

// Leaderboard represents a leaderboard entity
type Leaderboard struct {
	Id             uint           `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	Period         string         `json:"period"`
	ResetFrequency string         `json:"reset_frequency"`
	IsActive       bool           `json:"is_active"`
}

// TableName returns the table name for the Leaderboard model
func (item *Leaderboard) TableName() string {
	return "leaderboards"
}

// GetId returns the Id of the model
func (item *Leaderboard) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *Leaderboard) GetModelName() string {
	return "leaderboard"
}

// LeaderboardListResponse represents the list view response
type LeaderboardListResponse struct {
	Id             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Period         string    `json:"period"`
	ResetFrequency string    `json:"reset_frequency"`
	IsActive       bool      `json:"is_active"`
}

// LeaderboardResponse represents the detailed view response
type LeaderboardResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty"`
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	Period         string         `json:"period"`
	ResetFrequency string         `json:"reset_frequency"`
	IsActive       bool           `json:"is_active"`
}

// CreateLeaderboardRequest represents the request payload for creating a Leaderboard
type CreateLeaderboardRequest struct {
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required"`
	Period         string `json:"period" binding:"required"`
	ResetFrequency string `json:"reset_frequency" binding:"required"`
	IsActive       bool   `json:"is_active" binding:"required"`
}

// UpdateLeaderboardRequest represents the request payload for updating a Leaderboard
type UpdateLeaderboardRequest struct {
	Name           string `json:"name,omitempty"`
	Type           string `json:"type,omitempty"`
	Period         string `json:"period,omitempty"`
	ResetFrequency string `json:"reset_frequency,omitempty"`
	IsActive       string `json:"is_active,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *Leaderboard) ToListResponse() *LeaderboardListResponse {
	if item == nil {
		return nil
	}
	return &LeaderboardListResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		Name:           item.Name,
		Type:           item.Type,
		Period:         item.Period,
		ResetFrequency: item.ResetFrequency,
		IsActive:       item.IsActive,
	}
}

// ToResponse converts the model to a detailed response
func (item *Leaderboard) ToResponse() *LeaderboardResponse {
	if item == nil {
		return nil
	}
	return &LeaderboardResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt,
		Name:           item.Name,
		Type:           item.Type,
		Period:         item.Period,
		ResetFrequency: item.ResetFrequency,
		IsActive:       item.IsActive,
	}
}

// Preload preloads all the model's relationships
func (item *Leaderboard) Preload(db *gorm.DB) *gorm.DB {
	query := db
	return query
}
