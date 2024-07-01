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

func initDB() *gorm.DB {
	// Only load the .env file when running locally
	// Check for a RAILWAY_ENVIRONMENT, if not found, code is running locally
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	// Get the environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// Build the connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, dbPort)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func main() {
	db := initDB()
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
