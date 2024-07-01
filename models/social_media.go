package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type SocialMedia struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	UserID         uint       `json:"user_id" gorm:"not null"`
	Name           string     `json:"name" gorm:"not null" validate:"required"`
	SocialMediaURL string     `json:"social_media_url" gorm:"not null" validate:"required,url"`
	User           User       `json:"user" gorm:"foreignkey:UserID" validate:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index"`
}

type SMResponse struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Name           string    `json:"name" gorm:"not null" validate:"required"`
	SocialMediaURL string    `json:"social_media_url" gorm:"not null" validate:"required,url"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
}

type SMUpdateResponse struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Name           string    `json:"name" gorm:"not null" validate:"required"`
	SocialMediaURL string    `json:"social_media_url" gorm:"not null" validate:"required,url"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type GetSMResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name" gorm:"not null" validate:"required"`
	SocialMediaURL string    `json:"social_media_url" gorm:"not null" validate:"required,url"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           UserInfox `json:"user"`
}

type UserInfox struct {
	ID              uint   `json:"id"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
}

func (sm *SocialMedia) Validate() error {
	validate := validator.New()
	return validate.Struct(sm)
}
