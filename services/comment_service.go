package services

import (
	"errors"
	"mygram/models"
	"mygram/repositories"
)

type CommentService interface {
	GetComments(userId uint) ([]models.Comment, error)
	GetCommentByID(commentId uint) (models.Comment, error)
	CreateComment(comment models.Comment) (*models.Comment, error)
	UpdateComment(models.Comment) (models.Comment, error)
	DeleteComment(commentId uint, userId uint) error
}

type commentService struct {
	repo repositories.CommentRepository
}

func NewCommentService(repo repositories.CommentRepository) CommentService {
	return &commentService{repo}
}

func (s *commentService) GetComments(userId uint) ([]models.Comment, error) {
	return s.repo.GetComments(userId)
}

func (s *commentService) GetCommentByID(id uint) (models.Comment, error) {
	return s.repo.GetCommentByID(id)
}

func (s *commentService) CreateComment(comment models.Comment) (*models.Comment, error) {
	if err := comment.Validate(); err != nil {
		return nil, err
	}
	return s.repo.CreateComment(comment)
}

func (s *commentService) UpdateComment(comment models.Comment) (models.Comment, error) {
	existingComment, err := s.repo.GetCommentByID(comment.ID)
	if err != nil {
		return models.Comment{}, err
	}

	if existingComment.UserID != comment.UserID {
		return models.Comment{}, errors.New("unauthorized")
	}

	return s.repo.UpdateComment(comment)
}

func (s *commentService) DeleteComment(commentId uint, userId uint) error {
	comment, err := s.repo.GetCommentByID(commentId)
	if err != nil {
		return err
	}
	if comment.UserID != userId {
		return errors.New("unauthorized")
	}
	return s.repo.DeleteComment(commentId)
}
