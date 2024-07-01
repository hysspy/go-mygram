package services

import (
	"errors"
	"mygram/models"
	"mygram/repositories"
)

type PhotoService interface {
	GetPhotos(userId uint) ([]models.Photo, error)
	GetPhotoByID(photoId uint) (*models.Photo, error)
	CreatePhoto(photo models.Photo) (*models.Photo, error)
	IsUserAuthorized(userID uint, photoID uint) (bool, error)
	UpdatePhoto(photoID uint, input models.UpdatePhotoInput) (models.Photo, error)
	DeletePhoto(photoId uint, userId uint) error
}

type photoService struct {
	repo repositories.PhotoRepository
}

func NewPhotoService(repo repositories.PhotoRepository) PhotoService {
	return &photoService{repo}
}

func (s *photoService) GetPhotos(userId uint) ([]models.Photo, error) {
	return s.repo.GetPhotos(userId)
}

func (s *photoService) GetPhotoByID(photoId uint) (*models.Photo, error) {
	return s.repo.GetPhotoByID(photoId)
}

func (s *photoService) CreatePhoto(photo models.Photo) (*models.Photo, error) {
	createdPhoto, err := s.repo.CreatePhoto(photo)
	if err != nil {
		return nil, err
	}
	return createdPhoto, nil
}

func (s *photoService) IsUserAuthorized(userID uint, photoID uint) (bool, error) {
	photo, err := s.repo.GetPhotoByID(photoID)
	if err != nil {
		return false, err
	}

	return photo.UserID == userID, nil
}

func (s *photoService) UpdatePhoto(photoID uint, input models.UpdatePhotoInput) (models.Photo, error) {
	photo, err := s.repo.GetPhotoByID(photoID)
	if err != nil {
		return models.Photo{}, err
	}

	photo.Title = input.Title
	photo.Caption = input.Caption
	photo.PhotoURL = input.PhotoURL

	err = s.repo.SavePhoto(photo)
	if err != nil {
		return models.Photo{}, err
	}

	return *photo, nil
}

func (s *photoService) DeletePhoto(photoId uint, userId uint) error {
	photo, err := s.repo.GetPhotoByID(photoId)
	if err != nil {
		return err
	}
	if photo.UserID != userId {
		return errors.New("unauthorized")
	}
	return s.repo.DeletePhoto(photoId)
}
