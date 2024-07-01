package main

import (
	"fmt"
	"log"
	"mygram/controllers"
	"mygram/middlewares"
	"mygram/models"
	"mygram/repositories"
	"mygram/services"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	currentTime := time.Now()
	fmt.Println(strings.Repeat(" ", 100))
	fmt.Println(strings.Repeat("$", 100))
	fmt.Printf("Connected to PostgreSQL database at %s\n", currentTime.Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("$", 100))
	fmt.Println(strings.Repeat(" ", 100))

	// Automigrate models
	db.AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})
	db.AutoMigrate(&models.User{}, &models.Photo{})

	// Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	photoRepo := repositories.NewPhotoRepository(db)
	photoService := services.NewPhotoService(photoRepo)
	photoController := controllers.NewPhotoController(photoService)

	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)
	commentController := controllers.NewCommentController(commentService)

	socialMediaRepo := repositories.NewSocialMediaRepository(db)
	socialMediaService := services.NewSocialMediaService(socialMediaRepo)
	socialMediaController := controllers.NewSocialMediaController(socialMediaService)

	router := gin.Default()
	router.GET("/", salamsilaturahmi)
	users := router.Group("/users")
	users.POST("/login", userController.Login)
	users.POST("/register", userController.Register)

	auth := router.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		auth.GET("/photos", photoController.GetPhotos)
		auth.POST("/photos", photoController.CreatePhoto)
		auth.PUT("/photos/:photoId", photoController.UpdatePhoto)
		auth.DELETE("/photos/:photoId", photoController.DeletePhoto)

		auth.GET("/comments", commentController.GetComments)
		auth.POST("/comments", commentController.CreateComment)
		auth.PUT("/comments/:commentId", commentController.UpdateComment)
		auth.DELETE("/comments/:commentId", commentController.DeleteComment)

		auth.GET("/socialmedias", socialMediaController.GetSocialMedias)
		auth.POST("/socialmedias", socialMediaController.CreateSocialMedia)
		auth.PUT("/socialmedias/:socialMediaId", socialMediaController.UpdateSocialMedia)
		auth.DELETE("/socialmedias/:socialMediaId", socialMediaController.DeleteSocialMedia)

		auth.GET("/users/:userId", userController.GetUser)
		auth.PUT("/users/:userId", userController.UpdateUser)
		auth.DELETE("/users/:userId", userController.DeleteUser)
	}

	router.Run()
}

func salamsilaturahmi(c *gin.Context) {
	c.String(http.StatusOK, "MyGram Rest API / Final Project | UK | Kartu Prakerja Herwin Yudha Setyawan | Basic Golang untuk BackEnd Developer")
}
