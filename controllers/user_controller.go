package controllers

import (
	"mygram/models"
	"mygram/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := c.service.Register(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (c *UserController) Login(ctx *gin.Context) {
	var user struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(user.Email, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	idParam := ctx.Param("userId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ProfileImageURL = user.GenerateProfileImageURL()

	ctx.JSON(http.StatusOK, gin.H{
		"id":                user.ID,
		"username":          user.Username,
		"profile_image_url": user.ProfileImageURL,
	})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userIdStr := ctx.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userIDFromToken := ctx.MustGet("userID").(uint)
	if userIDFromToken != uint(userId) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this user"})
		return
	}

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uint(userId)

	updatedUser, err := c.service.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to response struct
	userResponse := models.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Username:  updatedUser.Username,
		Age:       updatedUser.Age,
		UpdatedAt: &updatedUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userIdStr := ctx.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userIDFromToken := ctx.MustGet("userID").(uint)
	if userIDFromToken != uint(userId) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete this user"})
		return
	}

	if err := c.service.DeleteUser(uint(userId)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
