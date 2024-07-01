package repositories

import (
	"mygram/models"

	"gorm.io/gorm"
)

type SocialMediaRepository interface {
	GetSocialMedias(userId uint) ([]models.SocialMedia, error)
	GetSocialMediaByID(id uint) (models.SocialMedia, error)
	CreateSocialMedia(socialMedia models.SocialMedia) (*models.SocialMedia, error)
	UpdateSocialMedia(socialMedia models.SocialMedia) (models.SocialMedia, error)
	DeleteSocialMedia(socialMediaId uint) error
}

type socialMediaRepository struct {
	db *gorm.DB
}

func NewSocialMediaRepository(db *gorm.DB) SocialMediaRepository {
	return &socialMediaRepository{db}
}

func (r *socialMediaRepository) GetSocialMedias(userId uint) ([]models.SocialMedia, error) {
	var socialMedias []models.SocialMedia
	if err := r.db.Where("user_id = ?", userId).Find(&socialMedias).Error; err != nil {
		return nil, err
	}
	if err := r.db.Preload("User").Find(&socialMedias).Error; err != nil {
		return nil, err
	}
	return socialMedias, nil
}

func (r *socialMediaRepository) GetSocialMediaByID(id uint) (models.SocialMedia, error) {
	var socialMedia models.SocialMedia
	if err := r.db.First(&socialMedia, id).Error; err != nil {
		return models.SocialMedia{}, err
	}
	return socialMedia, nil
}

func (r *socialMediaRepository) CreateSocialMedia(socialMedia models.SocialMedia) (*models.SocialMedia, error) {
	if err := r.db.Create(&socialMedia).Error; err != nil {
		return nil, err
	}
	return &socialMedia, nil
}

func (r *socialMediaRepository) UpdateSocialMedia(socialMedia models.SocialMedia) (models.SocialMedia, error) {
	if err := r.db.Save(&socialMedia).Error; err != nil {
		return models.SocialMedia{}, err
	}
	return socialMedia, nil
}

func (r *socialMediaRepository) DeleteSocialMedia(socialMediaId uint) error {
	if err := r.db.Delete(&models.SocialMedia{}, socialMediaId).Error; err != nil {
		return err
	}
	return nil
}
