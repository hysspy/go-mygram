package services

import (
	"errors"
	"mygram/models"
	"mygram/repositories"
)

type SocialMediaService interface {
	GetSocialMedias(userId uint) ([]models.SocialMedia, error)
	GetSocialMediaByID(id uint) (models.SocialMedia, error)
	CreateSocialMedia(socialMedia models.SocialMedia) (*models.SocialMedia, error)
	UpdateSocialMedia(socialMedia models.SocialMedia) (models.SocialMedia, error)
	DeleteSocialMedia(socialMediaId uint, userId uint) error
}

type socialMediaService struct {
	repo repositories.SocialMediaRepository
}

func NewSocialMediaService(repo repositories.SocialMediaRepository) SocialMediaService {
	return &socialMediaService{repo}
}

func (s *socialMediaService) GetSocialMedias(userId uint) ([]models.SocialMedia, error) {
	return s.repo.GetSocialMedias(userId)
}

func (s *socialMediaService) GetSocialMediaByID(id uint) (models.SocialMedia, error) {
	return s.repo.GetSocialMediaByID(id)
}

func (s *socialMediaService) CreateSocialMedia(socialMedia models.SocialMedia) (*models.SocialMedia, error) {
	if err := socialMedia.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateSocialMedia(socialMedia)
}

func (s *socialMediaService) UpdateSocialMedia(socialMedia models.SocialMedia) (models.SocialMedia, error) {
	existingSocialMedia, err := s.repo.GetSocialMediaByID(socialMedia.ID)
	if err != nil {
		return models.SocialMedia{}, err
	}

	if existingSocialMedia.UserID != socialMedia.UserID {
		return models.SocialMedia{}, errors.New("unauthorized")
	}

	return s.repo.UpdateSocialMedia(socialMedia)
}

func (s *socialMediaService) DeleteSocialMedia(socialMediaId uint, userId uint) error {
	socialMedia, err := s.repo.GetSocialMediaByID(socialMediaId)
	if err != nil {
		return err
	}
	if socialMedia.UserID != userId {
		return errors.New("unauthorized")
	}
	return s.repo.DeleteSocialMedia(socialMediaId)
}
