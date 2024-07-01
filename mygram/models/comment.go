package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Comment struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	UserID    uint       `json:"user_id"`
	PhotoID   uint       `json:"photo_id"`
	Message   string     `json:"message" gorm:"not null" validate:"required"`
	User      User       `json:"user" gorm:"foreignKey:UserID" validate:"-"`
	Photo     Photo      `json:"photo" validate:"-" gorm:"foreignKey:PhotoID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index"`
}

type CommentResponse struct {
	ID        uint      `gorm:"primaryKey"`
	Message   string    `json:"message" gorm:"not null" validate:"required"`
	PhotoID   uint      `json:"photo_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateCommentResponse struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message" gorm:"not null" validate:"required"`
	PhotoID   uint      `json:"photo_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Photo     PhotoInfo `json:"photo"`
}

type PhotoInfo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UserID    uint      `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Comment) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
