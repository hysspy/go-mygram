package controllers

import (
	"mygram/models"
	"mygram/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SocialMediaController struct {
	service services.SocialMediaService
}

func NewSocialMediaController(service services.SocialMediaService) *SocialMediaController {
	return &SocialMediaController{service}
}

func (c *SocialMediaController) GetSocialMedias(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	socialMedias, err := c.service.GetSocialMedias(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var smResponses []models.GetSMResponse
	for _, socialm := range socialMedias {
		smResponse := models.GetSMResponse{
			ID:             socialm.ID,
			Name:           socialm.Name,
			SocialMediaURL: socialm.SocialMediaURL,
			UserID:         socialm.UserID,
			CreatedAt:      socialm.CreatedAt,
			UpdatedAt:      socialm.UpdatedAt,
			User: models.UserInfox{
				ID:              socialm.User.ID,
				Email:           socialm.User.Email,
				Username:        socialm.User.Username,
				ProfileImageURL: socialm.User.ProfileImageURL,
			},
		}
		smResponses = append(smResponses, smResponse)
	}

	ctx.JSON(http.StatusOK, smResponses)
}

func (c *SocialMediaController) CreateSocialMedia(ctx *gin.Context) {
	var socialMedia models.SocialMedia
	if err := ctx.ShouldBindJSON(&socialMedia); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	socialMedia.UserID = ctx.MustGet("userID").(uint)
	createdSocialMedia, err := c.service.CreateSocialMedia(socialMedia)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map the createdSM to the response struct
	SMResponse := models.SMResponse{
		ID:             createdSocialMedia.ID,
		Name:           createdSocialMedia.Name,
		SocialMediaURL: createdSocialMedia.SocialMediaURL,
		UserID:         createdSocialMedia.UserID,
		CreatedAt:      time.Now(),
	}

	ctx.JSON(http.StatusCreated, SMResponse)
}

func (sc *SocialMediaController) UpdateSocialMedia(c *gin.Context) {
	var socialMedia models.SocialMedia
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("socialMediaId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid social media ID"})
		return
	}

	if err := c.ShouldBindJSON(&socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	socialMedia.ID = uint(id)
	socialMedia.UserID = userId.(uint)

	updatedSocialMedia, err := sc.service.UpdateSocialMedia(socialMedia)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Map the updatedSM to the response struct
	SMUpdateResponse := models.SMUpdateResponse{
		ID:             updatedSocialMedia.ID,
		Name:           updatedSocialMedia.Name,
		SocialMediaURL: updatedSocialMedia.SocialMediaURL,
		UserID:         updatedSocialMedia.UserID,
		UpdatedAt:      time.Now(),
	}

	c.JSON(http.StatusOK, SMUpdateResponse)
}

func (c *SocialMediaController) DeleteSocialMedia(ctx *gin.Context) {
	socialMediaIdStr := ctx.Param("socialMediaId")
	socialMediaId, err := strconv.ParseUint(socialMediaIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid social media ID"})
		return
	}

	userID := ctx.MustGet("userID").(uint)
	if err := c.service.DeleteSocialMedia(uint(socialMediaId), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Social media deleted successfully"})
}
