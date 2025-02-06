package models

import (
	"time"

	"gorm.io/gorm"
)

// AchievementCriteria represents a achievementcriteria entity
type AchievementCriteria struct {
	Id             uint           `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	AchievementId  uint           `json:"achievement_id"`
	Achievement    *Achievement   `json:"achievement,omitempty"`
	ActivityTypeId uint           `json:"activity_type_id"`
	ActivityType   *ActivityType  `json:"activity_type,omitempty"`
	RequiredCount  int            `json:"required_count"`
	TimeFrame      int            `json:"time_frame"`
}

// TableName returns the table name for the AchievementCriteria model
func (item *AchievementCriteria) TableName() string {
	return "achievementcriteria"
}

// GetId returns the Id of the model
func (item *AchievementCriteria) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *AchievementCriteria) GetModelName() string {
	return "achievementcriteria"
}

// AchievementCriteriaListResponse represents the list view response
type AchievementCriteriaListResponse struct {
	Id             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	AchievementId  uint      `json:"achievement_id"`
	ActivityTypeId uint      `json:"activity_type_id"`
	RequiredCount  int       `json:"required_count"`
	TimeFrame      int       `json:"time_frame"`
}

// AchievementCriteriaResponse represents the detailed view response
type AchievementCriteriaResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty"`
	AchievementId  uint           `json:"achievement_id"`
	Achievement    *Achievement   `json:"achievement,omitempty"`
	ActivityTypeId uint           `json:"activity_type_id"`
	ActivityType   *ActivityType  `json:"activity_type,omitempty"`
	RequiredCount  int            `json:"required_count"`
	TimeFrame      int            `json:"time_frame"`
}

// CreateAchievementCriteriaRequest represents the request payload for creating a AchievementCriteria
type CreateAchievementCriteriaRequest struct {
	AchievementId  uint `json:"achievement_id" binding:"required"`
	ActivityTypeId uint `json:"activity_type_id" binding:"required"`
	RequiredCount  int  `json:"required_count" binding:"required"`
	TimeFrame      int  `json:"time_frame" binding:"required"`
}

// UpdateAchievementCriteriaRequest represents the request payload for updating a AchievementCriteria
type UpdateAchievementCriteriaRequest struct {
	AchievementId  uint   `json:"achievement_id,omitempty"`
	ActivityTypeId uint   `json:"activity_type_id,omitempty"`
	RequiredCount  string `json:"required_count,omitempty"`
	TimeFrame      string `json:"time_frame,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *AchievementCriteria) ToListResponse() *AchievementCriteriaListResponse {
	if item == nil {
		return nil
	}
	return &AchievementCriteriaListResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		AchievementId:  item.AchievementId,
		ActivityTypeId: item.ActivityTypeId,
		RequiredCount:  item.RequiredCount,
		TimeFrame:      item.TimeFrame,
	}
}

// ToResponse converts the model to a detailed response
func (item *AchievementCriteria) ToResponse() *AchievementCriteriaResponse {
	if item == nil {
		return nil
	}
	return &AchievementCriteriaResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt,
		AchievementId:  item.AchievementId,
		Achievement:    item.Achievement,
		ActivityTypeId: item.ActivityTypeId,
		ActivityType:   item.ActivityType,
		RequiredCount:  item.RequiredCount,
		TimeFrame:      item.TimeFrame,
	}
}

// Preload preloads all the model's relationships
func (item *AchievementCriteria) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Achievement")
	query = query.Preload("ActivityType")
	return query
}
