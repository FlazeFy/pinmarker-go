package routes

import (
	"pinmarker/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouteTrack(api *gin.RouterGroup, trackController *controllers.TrackController) {
	track := api.Group("/tracks")
	{
		track.POST("/", trackController.CreateTrack)
		track.POST("/multi", trackController.CreateTrackMulti)
		track.GET("/:app_source/:created_by", trackController.GetAllTrack)
		track.DELETE("/:app_source/:created_by/:track_id", trackController.DeleteTrackById)
	}
}
