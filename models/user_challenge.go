package models

import (
	"base/core/app/users"
	"base/core/types"
	"time"

	"gorm.io/gorm"
)

// UserChallenge represents a userchallenge entity
type UserChallenge struct {
	Id            uint           `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	UserId        uint           `json:"user_id"`
	User          *users.User    `json:"user,omitempty"`
	ChallengeId   uint           `json:"challenge_id"`
	Challenge     *Challenge     `json:"challenge,omitempty"`
	Progress      int            `json:"progress"`
	CompletedAt   types.DateTime `json:"completed_at"`
	RewardClaimed bool           `json:"reward_claimed"`
}

// TableName returns the table name for the UserChallenge model
func (item *UserChallenge) TableName() string {
	return "userchallenges"
}

// GetId returns the Id of the model
func (item *UserChallenge) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *UserChallenge) GetModelName() string {
	return "userchallenge"
}

// UserChallengeListResponse represents the list view response
type UserChallengeListResponse struct {
	Id            uint           `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	UserId        uint           `json:"user_id"`
	ChallengeId   uint           `json:"challenge_id"`
	Progress      int            `json:"progress"`
	CompletedAt   types.DateTime `json:"completed_at"`
	RewardClaimed bool           `json:"reward_claimed"`
}

// UserChallengeResponse represents the detailed view response
type UserChallengeResponse struct {
	Id            uint           `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty"`
	UserId        uint           `json:"user_id"`
	User          *users.User    `json:"user,omitempty"`
	ChallengeId   uint           `json:"challenge_id"`
	Challenge     *Challenge     `json:"challenge,omitempty"`
	Progress      int            `json:"progress"`
	CompletedAt   types.DateTime `json:"completed_at"`
	RewardClaimed bool           `json:"reward_claimed"`
}

// CreateUserChallengeRequest represents the request payload for creating a UserChallenge
type CreateUserChallengeRequest struct {
	UserId        uint           `json:"user_id" binding:"required"`
	ChallengeId   uint           `json:"challenge_id" binding:"required"`
	Progress      int            `json:"progress" binding:"required"`
	CompletedAt   types.DateTime `json:"completed_at" binding:"required"`
	RewardClaimed bool           `json:"reward_claimed" binding:"required"`
}

// UpdateUserChallengeRequest represents the request payload for updating a UserChallenge
type UpdateUserChallengeRequest struct {
	UserId        uint   `json:"user_id,omitempty"`
	ChallengeId   uint   `json:"challenge_id,omitempty"`
	Progress      string `json:"progress,omitempty"`
	CompletedAt   string `json:"completed_at,omitempty"`
	RewardClaimed string `json:"reward_claimed,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *UserChallenge) ToListResponse() *UserChallengeListResponse {
	if item == nil {
		return nil
	}
	return &UserChallengeListResponse{
		Id:            item.Id,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		UserId:        item.UserId,
		ChallengeId:   item.ChallengeId,
		Progress:      item.Progress,
		CompletedAt:   item.CompletedAt,
		RewardClaimed: item.RewardClaimed,
	}
}

// ToResponse converts the model to a detailed response
func (item *UserChallenge) ToResponse() *UserChallengeResponse {
	if item == nil {
		return nil
	}
	return &UserChallengeResponse{
		Id:            item.Id,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		DeletedAt:     item.DeletedAt,
		UserId:        item.UserId,
		User:          item.User,
		ChallengeId:   item.ChallengeId,
		Challenge:     item.Challenge,
		Progress:      item.Progress,
		CompletedAt:   item.CompletedAt,
		RewardClaimed: item.RewardClaimed,
	}
}

// Preload preloads all the model's relationships
func (item *UserChallenge) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("User")
	query = query.Preload("Challenge")
	return query
}
