package controllers

import (
	"mygram/models"
	"mygram/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PhotoController struct {
	service services.PhotoService
}

func NewPhotoController(service services.PhotoService) *PhotoController {
	return &PhotoController{service}
}

func (c *PhotoController) GetPhotos(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	photos, err := c.service.GetPhotos(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var photoResponses []models.GetPhotoResponse
	for _, photo := range photos {
		photoResponse := models.GetPhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoURL:  photo.PhotoURL,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: models.UserInfo{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		}
		photoResponses = append(photoResponses, photoResponse)
	}
	ctx.JSON(http.StatusOK, photoResponses)
}

func (c *PhotoController) CreatePhoto(ctx *gin.Context) {
	var photo models.Photo
	if err := ctx.ShouldBindJSON(&photo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perform validation
	validate := validator.New()
	if err := validate.Struct(photo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo.UserID = ctx.MustGet("userID").(uint)
	createdPhoto, err := c.service.CreatePhoto(photo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map the createdPhoto to the response struct
	photoResponse := models.PhotoResponse{
		ID:        createdPhoto.ID,
		Title:     createdPhoto.Title,
		Caption:   createdPhoto.Caption,
		PhotoURL:  createdPhoto.PhotoURL,
		UserID:    createdPhoto.UserID,
		CreatedAt: time.Now(),
	}

	ctx.JSON(http.StatusCreated, photoResponse)
}

func (c *PhotoController) UpdatePhoto(ctx *gin.Context) {
	var input models.UpdatePhotoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photoID, err := strconv.Atoi(ctx.Param("photoId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	userID, _ := ctx.Get("userID")
	authorized, err := c.service.IsUserAuthorized(userID.(uint), uint(photoID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !authorized {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	photo, err := c.service.UpdatePhoto(uint(photoID), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map the UpdatedPhoto to the response struct
	updatephotoResponse := models.UpdatePhotoResponse{
		ID:        photo.ID,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoURL:  photo.PhotoURL,
		UserID:    photo.UserID,
		UpdatedAt: time.Now(),
	}

	ctx.JSON(http.StatusOK, updatephotoResponse)
}

func (c *PhotoController) DeletePhoto(ctx *gin.Context) {
	photoIdStr := ctx.Param("photoId")
	photoId, err := strconv.ParseUint(photoIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	userID := ctx.MustGet("userID").(uint)
	if err := c.service.DeletePhoto(uint(photoId), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
