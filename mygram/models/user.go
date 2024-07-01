package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	Username        string `json:"username" validate:"required" gorm:"unique;not null"`
	Email           string `json:"email" validate:"required,email" gorm:"unique;not null"`
	Password        string `json:"password" validate:"required,min=6"`
	Age             uint   `json:"age" gorm:"not null" validate:"required,gt=8"`
	SocialMediaURL  string `gorm:"-"`
	ProfileImageURL string `json:"profile_image_url"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time    `gorm:"index"`
	SocialMedias    []SocialMedia `json:"-"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Age       uint       `json:"age"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) GenerateProfileImageURL() string {
	return fmt.Sprintf("%s/profile/%s", u.SocialMediaURL, u.Username)
}
