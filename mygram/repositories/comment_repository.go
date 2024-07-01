package repositories

import (
	"mygram/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	GetComments(userId uint) ([]models.Comment, error)
	GetCommentByID(id uint) (models.Comment, error)
	CreateComment(comment models.Comment) (*models.Comment, error)
	UpdateComment(comment models.Comment) (models.Comment, error)
	DeleteComment(commentId uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db}
}

func (r *commentRepository) GetComments(userId uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Where("user_id = ?", userId).Find(&comments).Error; err != nil {
		return nil, err
	}
	if err := r.db.Preload("User").Preload("Photo").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) GetCommentByID(id uint) (models.Comment, error) {
	var comment models.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

func (r *commentRepository) CreateComment(comment models.Comment) (*models.Comment, error) {
	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) UpdateComment(comment models.Comment) (models.Comment, error) {
	if err := r.db.Save(&comment).Error; err != nil {
		return models.Comment{}, err
	}
	if err := r.db.Preload("Photo").First(&comment, comment.ID).Error; err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

func (r *commentRepository) DeleteComment(commentId uint) error {
	if err := r.db.Delete(&models.Comment{}, commentId).Error; err != nil {
		return err
	}
	return nil
}
