package models

import (
	"base/core/types"
	"time"

	"gorm.io/gorm"
)

// Challenge represents a challenge entity
type Challenge struct {
	Id          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	StartDate   types.DateTime `json:"start_date"`
	EndDate     types.DateTime `json:"end_date"`
	RewardType  string         `json:"reward_type"`
	RewardValue string         `json:"reward_value"`
	IsActive    bool           `json:"is_active"`
}

// TableName returns the table name for the Challenge model
func (item *Challenge) TableName() string {
	return "challenges"
}

// GetId returns the Id of the model
func (item *Challenge) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *Challenge) GetModelName() string {
	return "challenge"
}

// ChallengeListResponse represents the list view response
type ChallengeListResponse struct {
	Id          uint           `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	StartDate   types.DateTime `json:"start_date"`
	EndDate     types.DateTime `json:"end_date"`
	RewardType  string         `json:"reward_type"`
	RewardValue string         `json:"reward_value"`
	IsActive    bool           `json:"is_active"`
}

// ChallengeResponse represents the detailed view response
type ChallengeResponse struct {
	Id          uint           `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	StartDate   types.DateTime `json:"start_date"`
	EndDate     types.DateTime `json:"end_date"`
	RewardType  string         `json:"reward_type"`
	RewardValue string         `json:"reward_value"`
	IsActive    bool           `json:"is_active"`
}

// CreateChallengeRequest represents the request payload for creating a Challenge
type CreateChallengeRequest struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description" binding:"required"`
	StartDate   types.DateTime `json:"start_date" binding:"required"`
	EndDate     types.DateTime `json:"end_date" binding:"required"`
	RewardType  string         `json:"reward_type" binding:"required"`
	RewardValue string         `json:"reward_value" binding:"required"`
	IsActive    bool           `json:"is_active" binding:"required"`
}

// UpdateChallengeRequest represents the request payload for updating a Challenge
type UpdateChallengeRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	RewardType  string `json:"reward_type,omitempty"`
	RewardValue string `json:"reward_value,omitempty"`
	IsActive    string `json:"is_active,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *Challenge) ToListResponse() *ChallengeListResponse {
	if item == nil {
		return nil
	}
	return &ChallengeListResponse{
		Id:          item.Id,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		Name:        item.Name,
		Description: item.Description,
		StartDate:   item.StartDate,
		EndDate:     item.EndDate,
		RewardType:  item.RewardType,
		RewardValue: item.RewardValue,
		IsActive:    item.IsActive,
	}
}

// ToResponse converts the model to a detailed response
func (item *Challenge) ToResponse() *ChallengeResponse {
	if item == nil {
		return nil
	}
	return &ChallengeResponse{
		Id:          item.Id,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		DeletedAt:   item.DeletedAt,
		Name:        item.Name,
		Description: item.Description,
		StartDate:   item.StartDate,
		EndDate:     item.EndDate,
		RewardType:  item.RewardType,
		RewardValue: item.RewardValue,
		IsActive:    item.IsActive,
	}
}

// Preload preloads all the model's relationships
func (item *Challenge) Preload(db *gorm.DB) *gorm.DB {
	query := db
	return query
}
