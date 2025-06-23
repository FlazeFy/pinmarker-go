package routes

import (
	"pinmarker/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouteTrack(api *gin.RouterGroup, trackController *controllers.TrackController) {
	track := api.Group("/tracks")
	{
		track.POST("/", trackController.Create)
	}
}
