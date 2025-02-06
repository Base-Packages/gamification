package models

import (
	"base/core/app/users"
	"base/core/types"
	"time"

	"gorm.io/gorm"
)

// LeaderboardEntry represents a leaderboardentry entity
type LeaderboardEntry struct {
	Id            uint           `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	LeaderboardId uint           `json:"leaderboard_id"`
	Leaderboard   *Leaderboard   `json:"leaderboard,omitempty"`
	UserId        uint           `json:"user_id"`
	User          *users.User    `json:"user,omitempty"`
	Score         int            `json:"score"`
	Rank          int            `json:"rank"`
	PeriodStart   types.DateTime `json:"period_start"`
	PeriodEnd     types.DateTime `json:"period_end"`
}

// TableName returns the table name for the LeaderboardEntry model
func (item *LeaderboardEntry) TableName() string {
	return "leaderboardentries"
}

// GetId returns the Id of the model
func (item *LeaderboardEntry) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *LeaderboardEntry) GetModelName() string {
	return "leaderboardentry"
}

// LeaderboardEntryListResponse represents the list view response
type LeaderboardEntryListResponse struct {
	Id            uint           `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	LeaderboardId uint           `json:"leaderboard_id"`
	UserId        uint           `json:"user_id"`
	Score         int            `json:"score"`
	Rank          int            `json:"rank"`
	PeriodStart   types.DateTime `json:"period_start"`
	PeriodEnd     types.DateTime `json:"period_end"`
}

// LeaderboardEntryResponse represents the detailed view response
type LeaderboardEntryResponse struct {
	Id            uint           `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty"`
	LeaderboardId uint           `json:"leaderboard_id"`
	Leaderboard   *Leaderboard   `json:"leaderboard,omitempty"`
	UserId        uint           `json:"user_id"`
	User          *users.User    `json:"user,omitempty"`
	Score         int            `json:"score"`
	Rank          int            `json:"rank"`
	PeriodStart   types.DateTime `json:"period_start"`
	PeriodEnd     types.DateTime `json:"period_end"`
}

// CreateLeaderboardEntryRequest represents the request payload for creating a LeaderboardEntry
type CreateLeaderboardEntryRequest struct {
	LeaderboardId uint           `json:"leaderboard_id" binding:"required"`
	UserId        uint           `json:"user_id" binding:"required"`
	Score         int            `json:"score" binding:"required"`
	Rank          int            `json:"rank" binding:"required"`
	PeriodStart   types.DateTime `json:"period_start" binding:"required"`
	PeriodEnd     types.DateTime `json:"period_end" binding:"required"`
}

// UpdateLeaderboardEntryRequest represents the request payload for updating a LeaderboardEntry
type UpdateLeaderboardEntryRequest struct {
	LeaderboardId uint   `json:"leaderboard_id,omitempty"`
	UserId        uint   `json:"user_id,omitempty"`
	Score         string `json:"score,omitempty"`
	Rank          string `json:"rank,omitempty"`
	PeriodStart   string `json:"period_start,omitempty"`
	PeriodEnd     string `json:"period_end,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *LeaderboardEntry) ToListResponse() *LeaderboardEntryListResponse {
	if item == nil {
		return nil
	}
	return &LeaderboardEntryListResponse{
		Id:            item.Id,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		LeaderboardId: item.LeaderboardId,
		UserId:        item.UserId,
		Score:         item.Score,
		Rank:          item.Rank,
		PeriodStart:   item.PeriodStart,
		PeriodEnd:     item.PeriodEnd,
	}
}

// ToResponse converts the model to a detailed response
func (item *LeaderboardEntry) ToResponse() *LeaderboardEntryResponse {
	if item == nil {
		return nil
	}
	return &LeaderboardEntryResponse{
		Id:            item.Id,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		DeletedAt:     item.DeletedAt,
		LeaderboardId: item.LeaderboardId,
		Leaderboard:   item.Leaderboard,
		UserId:        item.UserId,
		User:          item.User,
		Score:         item.Score,
		Rank:          item.Rank,
		PeriodStart:   item.PeriodStart,
		PeriodEnd:     item.PeriodEnd,
	}
}

// Preload preloads all the model's relationships
func (item *LeaderboardEntry) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Leaderboard")
	query = query.Preload("User")
	return query
}
