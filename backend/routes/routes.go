package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"wira-dashboard/handlers"
	"wira-dashboard/middleware"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// Create handlers
	rankingHandler := handlers.NewHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// API routes group
	api := r.Group("/api")
	{
		// Public auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile routes
			user := protected.Group("/user")
			{
				user.GET("/profile", authHandler.GetProfile)
				user.POST("/change-password", authHandler.ChangePassword)
				user.GET("/activities", authHandler.GetUserActivities)
			}

			// 2FA routes
			twoFA := protected.Group("/2fa")
			{
				twoFA.GET("/status", authHandler.Get2FAStatus)
				twoFA.POST("/setup", authHandler.Setup2FA)
				twoFA.POST("/enable", authHandler.Enable2FA)
				twoFA.POST("/disable", authHandler.Disable2FA)
			}

			// Rankings endpoints (now protected)
			rankings := protected.Group("/rankings")
			{
				rankings.GET("", rankingHandler.GetRankings)
				rankings.GET("/search", rankingHandler.SearchRankings)
				rankings.GET("/stats", rankingHandler.GetClassStats)
			}
		}

		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"message": "WIRA Ranking API is running",
			})
		})
	}
}
