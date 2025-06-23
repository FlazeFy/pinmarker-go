package routes

import (
	"pinmarker/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine,
	trackController *controllers.TrackController) {

	// V1 Endpoint
	api := r.Group("/api/v1")

	// Routes Endpoint
	SetUpRouteTrack(api, trackController)
}
