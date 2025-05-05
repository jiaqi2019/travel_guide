package routes

import (
	"travel_guide/controllers"
	"travel_guide/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// User routes
	userController := controllers.NewUserController(db)
	r.POST("/api/register", userController.CreateUser)
	r.POST("/api/login", userController.Login)
	r.GET("/api/users", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), userController.GetUsers)
	r.PUT("/api/users/:id/status", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), userController.UpdateUserStatus)

	// Guide routes
	guideController := controllers.NewGuideController(db)
	guideRoutes := r.Group("/api/guides")
	{
		guideRoutes.POST("", middleware.AuthMiddleware(), guideController.CreateGuide)
		guideRoutes.GET("", guideController.GetGuides)
		guideRoutes.GET("/:id", guideController.GetGuideDetail)
		guideRoutes.GET("/suggestions", guideController.GetSearchSuggestions)
		guideRoutes.GET("/search", middleware.OptionalAuthMiddleware(db), guideController.SearchGuides)
		guideRoutes.GET("/recommendations", middleware.AuthMiddleware(), guideController.GetUserRecommendations)
	}

	// Tag routes
	tagController := controllers.NewTagController(db)
	tagRoutes := r.Group("/api/tags")
	{
		tagRoutes.GET("", tagController.GetAllTags)
		tagRoutes.GET("/related", tagController.GetRelatedTags)
	}

	// Upload routes
	uploadController := controllers.NewUploadController()
	uploadRoutes := r.Group("/api/upload")
	{
		uploadRoutes.POST("/image", middleware.AuthMiddleware(), uploadController.UploadImage)
	}
}
