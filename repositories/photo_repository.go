package repositories

import (
	"mygram/models"

	"gorm.io/gorm"
)

type PhotoRepository interface {
	GetPhotos(userId uint) ([]models.Photo, error)
	GetPhotoByID(photoId uint) (*models.Photo, error)
	CreatePhoto(photo models.Photo) (*models.Photo, error)
	UpdatePhoto(photo models.Photo) (*models.Photo, error)
	SavePhoto(photo *models.Photo) error
	DeletePhoto(photoId uint) error
}

type photoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) PhotoRepository {
	return &photoRepository{db}
}

func (r *photoRepository) GetPhotos(userId uint) ([]models.Photo, error) {
	var photos []models.Photo
	if err := r.db.Where("user_id = ?", userId).Find(&photos).Error; err != nil {
		return nil, err
	}
	if err := r.db.Preload("User").Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (r *photoRepository) GetPhotoByID(photoId uint) (*models.Photo, error) {
	var photo models.Photo
	if err := r.db.First(&photo, photoId).Error; err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *photoRepository) SavePhoto(photo *models.Photo) error {
	return r.db.Save(photo).Error
}

func (r *photoRepository) CreatePhoto(photo models.Photo) (*models.Photo, error) {
	if err := r.db.Create(&photo).Error; err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *photoRepository) UpdatePhoto(photo models.Photo) (*models.Photo, error) {
	if err := r.db.Save(&photo).Error; err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *photoRepository) DeletePhoto(photoId uint) error {
	if err := r.db.Delete(&models.Photo{}, photoId).Error; err != nil {
		return err
	}
	return nil
}
