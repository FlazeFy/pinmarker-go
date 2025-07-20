package routes

import (
	"pinmarker/controllers"
	"pinmarker/repositories"
	"pinmarker/services"

	"github.com/gin-gonic/gin"
)

func SetUpDependency(r *gin.Engine) {
	// Setup Repository
	trackRepo := repositories.NewTrackRepository()

	// Setup Service
	trackService := services.NewTrackService(trackRepo)

	// Setup Controller
	trackController := controllers.NewTrackController(trackService)

	// Setup Routes
	SetUpRoutes(r, trackController)

	// Task Scheduler
	SetUpScheduler(trackService)
}
