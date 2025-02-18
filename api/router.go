package api

import (
	"go_back/config"
	"go_back/controllers"
	"go_back/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/signup", controllers.Register)
	r.POST("/login", controllers.Login)

	// Protected routes
	protected := r.Group("/", middleware.AuthMiddleware(cfg))
	{
		protected.POST("/urls", controllers.CreateURL)
		protected.GET("/urls", controllers.GetURLs)
		protected.DELETE("/urls/:id", controllers.DeleteURL)
		protected.PUT("/urls/:id", controllers.UpdateURL)

		protected.GET("/urls/service/:service", controllers.GetURLsByService)

		// Discovery route
		//protected.POST("/discover/:service", controllers.DiscoverURLs)
		//protected.GET("/spider/service/:service", controllers.StartSpiderScan)
	}

	return r
}
