package controllers

import (
	"mygram/models"
	"mygram/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service services.CommentService
}

func NewCommentController(service services.CommentService) *CommentController {
	return &CommentController{service}
}

func (c *CommentController) GetComments(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	comments, err := c.service.GetComments(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []map[string]interface{}
	for _, comment := range comments {
		commentData := map[string]interface{}{
			"id":         comment.ID,
			"message":    comment.Message,
			"photo_id":   comment.PhotoID,
			"user_id":    comment.UserID,
			"updated_at": comment.UpdatedAt,
			"created_at": comment.CreatedAt,
			"user": map[string]interface{}{
				"id":       comment.User.ID,
				"email":    comment.User.Email,
				"username": comment.User.Username,
			},
			"photo": map[string]interface{}{
				"id":        comment.Photo.ID,
				"title":     comment.Photo.Title,
				"caption":   comment.Photo.Caption,
				"photo_url": comment.Photo.PhotoURL,
				"user_id":   comment.Photo.UserID,
			},
		}
		response = append(response, commentData)
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *CommentController) CreateComment(ctx *gin.Context) {
	var comment models.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.UserID = ctx.MustGet("userID").(uint)
	createdComment, err := c.service.CreateComment(comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map the createdPhoto to the response struct
	CommentResponse := models.CommentResponse{
		ID:        createdComment.ID,
		Message:   createdComment.Message,
		PhotoID:   createdComment.PhotoID,
		UserID:    createdComment.UserID,
		CreatedAt: time.Now(),
	}

	ctx.JSON(http.StatusCreated, CommentResponse)
}

func (cc *CommentController) UpdateComment(c *gin.Context) {
	var comment models.Comment
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.ID = uint(id)
	comment.UserID = userId.(uint)

	updatedComment, err := cc.service.UpdateComment(comment)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var upResponses []models.UpdateCommentResponse
	upResponse := models.UpdateCommentResponse{
		ID:        updatedComment.ID,
		Message:   updatedComment.Message,
		PhotoID:   updatedComment.PhotoID,
		UserID:    updatedComment.UserID,
		CreatedAt: updatedComment.CreatedAt,
		Photo: models.PhotoInfo{
			ID:        updatedComment.Photo.ID,
			Title:     updatedComment.Photo.Title,
			Caption:   updatedComment.Photo.Caption,
			PhotoURL:  updatedComment.Photo.PhotoURL,
			UserID:    updatedComment.Photo.UserID,
			UpdatedAt: updatedComment.Photo.UpdatedAt,
		},
	}
	upResponses = append(upResponses, upResponse)

	c.JSON(http.StatusOK, upResponses)
}

func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentIdStr := ctx.Param("commentId")
	commentId, err := strconv.ParseUint(commentIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := ctx.MustGet("userID").(uint)
	if err := c.service.DeleteComment(uint(commentId), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
