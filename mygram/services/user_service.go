package services

import (
	"errors"
	"mygram/models"
	"mygram/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user models.User) (*models.User, error)
	Login(email, password string) (string, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user models.User) (*models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(user models.User) (*models.User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	return s.repo.CreateUser(user)
}

func (s *userService) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) UpdateUser(user models.User) (*models.User, error) {
	updatedUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}
