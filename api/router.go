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
		protected.POST("/api/v1/add/url", controllers.CreateURL)
		protected.GET("/api/v1/get/url", controllers.GetURLs)
		protected.DELETE("/api/v1/delete/url/:id", controllers.DeleteURL)
		protected.PUT("/api/v1/update/url/:id", controllers.UpdateURL)
		protected.GET("/api/v1/get/service-urls/:service", controllers.GetURLsByService)

		// API target management routes
		protected.POST("/api/v1/add/api", controllers.CreateAPITarget)
		protected.GET("/api/v1/get/api", controllers.GetAPITargets)
		protected.PUT("/api/v1/update/api/:id", controllers.UpdateAPITarget)
		protected.DELETE("/api/v1/delete/api/:id", controllers.DeleteAPITarget)
		protected.GET("/api/v1/get/service-apis/:service", controllers.GetAPITargetsByService)

		// Discovery route
		protected.POST("/api/v1/url-discovery/service/:service", controllers.StartSpiderScan)
		protected.POST("/api/v1/url-scan/service/:service", controllers.StartActiveScan)
		protected.GET("/api/v1/url-report/service/:service", controllers.GetAlerts)
		protected.GET("/api/v1/check/service/:service", controllers.CheckUrl)
		protected.POST("/api/v1/import-api/service/:service", controllers.ImportApiTarget)
		protected.POST("/api/v1/api-scan/service/:service", controllers.StartApiScan)
		protected.GET("/api/v1/api-report/service/:service", controllers.GetApiAlerts)

	}

	return r
}
