package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title" gorm:"not null" validate:"required"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url" gorm:"not null" validate:"required,url"`
	User      User      `json:"user" gorm:"foreignKey:UserID" validate:"-"`
	Comments  []Comment `json:"comments"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type PhotoResponse struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null" validate:"required"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url" gorm:"not null" validate:"required,url"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type GetPhotoResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      UserInfo  `json:"user"`
}

type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UpdatePhotoInput struct {
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption" binding:"required"`
	PhotoURL string `json:"photo_url" binding:"required,url"`
}

type UpdatePhotoResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UserID    uint      `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Photo) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
