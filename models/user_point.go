package models

import (
	"base/core/app/users"
	"time"

	"gorm.io/gorm"
)

// UserPoint represents a userpoint entity
type UserPoint struct {
	Id             uint           `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	UserId         uint           `json:"user_id"`
	User           *users.User    `json:"user,omitempty"`
	PointTypeId    uint           `json:"point_type_id"`
	PointType      *PointType     `json:"point_type,omitempty"`
	CurrentBalance int            `json:"current_balance"`
	LifetimeEarned int            `json:"lifetime_earned"`
}

// TableName returns the table name for the UserPoint model
func (item *UserPoint) TableName() string {
	return "userpoints"
}

// GetId returns the Id of the model
func (item *UserPoint) GetId() uint {
	return item.Id
}

// GetModelName returns the model name
func (item *UserPoint) GetModelName() string {
	return "userpoint"
}

// UserPointListResponse represents the list view response
type UserPointListResponse struct {
	Id             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserId         uint      `json:"user_id"`
	PointTypeId    uint      `json:"point_type_id"`
	CurrentBalance int       `json:"current_balance"`
	LifetimeEarned int       `json:"lifetime_earned"`
}

// UserPointResponse represents the detailed view response
type UserPointResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty"`
	UserId         uint           `json:"user_id"`
	User           *users.User    `json:"user,omitempty"`
	PointTypeId    uint           `json:"point_type_id"`
	PointType      *PointType     `json:"point_type,omitempty"`
	CurrentBalance int            `json:"current_balance"`
	LifetimeEarned int            `json:"lifetime_earned"`
}

// CreateUserPointRequest represents the request payload for creating a UserPoint
type CreateUserPointRequest struct {
	UserId         uint `json:"user_id" binding:"required"`
	PointTypeId    uint `json:"point_type_id" binding:"required"`
	CurrentBalance int  `json:"current_balance" binding:"required"`
	LifetimeEarned int  `json:"lifetime_earned" binding:"required"`
}

// UpdateUserPointRequest represents the request payload for updating a UserPoint
type UpdateUserPointRequest struct {
	UserId         uint   `json:"user_id,omitempty"`
	PointTypeId    uint   `json:"point_type_id,omitempty"`
	CurrentBalance string `json:"current_balance,omitempty"`
	LifetimeEarned string `json:"lifetime_earned,omitempty"`
}

// ToListResponse converts the model to a list response
func (item *UserPoint) ToListResponse() *UserPointListResponse {
	if item == nil {
		return nil
	}
	return &UserPointListResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		UserId:         item.UserId,
		PointTypeId:    item.PointTypeId,
		CurrentBalance: item.CurrentBalance,
		LifetimeEarned: item.LifetimeEarned,
	}
}

// ToResponse converts the model to a detailed response
func (item *UserPoint) ToResponse() *UserPointResponse {
	if item == nil {
		return nil
	}
	return &UserPointResponse{
		Id:             item.Id,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt,
		UserId:         item.UserId,
		User:           item.User,
		PointTypeId:    item.PointTypeId,
		PointType:      item.PointType,
		CurrentBalance: item.CurrentBalance,
		LifetimeEarned: item.LifetimeEarned,
	}
}

// Preload preloads all the model's relationships
func (item *UserPoint) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("User")
	query = query.Preload("PointType")
	return query
}
