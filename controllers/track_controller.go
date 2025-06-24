package controllers

import (
	"fmt"
	"math"
	"net/http"
	"pinmarker/configs"
	"pinmarker/entities"
	"pinmarker/services"
	"pinmarker/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrackController struct {
	TrackService services.TrackService
}

func NewTrackController(trackService services.TrackService) *TrackController {
	return &TrackController{TrackService: trackService}
}

// @Summary      Create Track
// @Description  Create an track
// @Tags         Track
// @Accept       json
// @Produce      json
// @Param        request  body  entities.RequestCreateTrack  true  "Post Track Request Body"
// @Success      201  {object}  entities.ResponseCreateTrack
// @Failure      404  {object}  entities.ResponseBadRequest
// @Router       /api/v1/tracks [post]
func (tr *TrackController) CreateTrack(c *gin.Context) {
	// Model
	var req entities.Track

	// Validator JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validator Field
	if req.TrackLat == "" {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "track latitude is required")
		return
	}
	if req.TrackLong == "" {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "track longitude is required")
		return
	}
	if req.TrackType == "" {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "track type is required")
		return
	}
	if req.AppsSource == "" {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "app source is required")
		return
	}

	// Validator UUID
	if req.CreatedBy == uuid.Nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "created by is required and must be a valid UUID")
		return
	}

	// Validator : Track Type & Apps Source
	if !utils.ValidatorContains(configs.TrackTypes, req.TrackType) {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "track type is not valid")
		return
	}
	if !utils.ValidatorContains(configs.AppsSources, req.AppsSource) {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "app source is not valid")
		return
	}

	// Service : Create Track
	err := tr.TrackService.CreateTrack(&req)
	if err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.MessageResponseBuild(c, "success", "track", "post", http.StatusCreated, &req, nil)
}

// @Summary      Create Track Multiple
// @Description  Create multiple track
// @Tags         Track
// @Accept       json
// @Produce      json
// @Param        request  body  entities.RequestCreateTrackMulti  true  "Post Track Multiple Request Body"
// @Success      201  {object}  entities.ResponseCreateTrackMulti
// @Failure      404  {object}  entities.ResponseBadRequest
// @Router       /api/v1/tracks/multi [post]
func (tr *TrackController) CreateTrackMulti(c *gin.Context) {
	// Validator JSON
	var req []*entities.Track
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate each item
	for i, track := range req {
		// Validator Field
		if track.TrackLat == "" {
			utils.MessageResponseErrorBuild(c, http.StatusBadRequest, fmt.Sprintf("track latitude is required at index %d", i))
			return
		}
		if track.TrackLong == "" {
			utils.MessageResponseErrorBuild(c, http.StatusBadRequest, fmt.Sprintf("track longitude is required at index %d", i))
			return
		}
		if track.TrackType == "" {
			utils.MessageResponseErrorBuild(c, http.StatusBadRequest, fmt.Sprintf("track type is required at index %d", i))
			return
		}
		if track.AppsSource == "" {
			utils.MessageResponseErrorBuild(c, http.StatusBadRequest, fmt.Sprintf("app source is required at index %d", i))
			return
		}
		if track.CreatedAt.String() == "" {
			utils.MessageResponseErrorBuild(c, http.StatusBadRequest, fmt.Sprintf("created at is required at index %d", i))
			return
		}
		if track.CreatedBy == uuid.Nil {
			utils.MessageResponseErrorBuild(c, http.StatusBadRequest, fmt.Sprintf("created by must be a valid UUID at index %d", i))
			return
		}

		// Validator : Track Type & Apps Source
		if !utils.ValidatorContains(configs.TrackTypes, track.TrackType) {
			utils.MessageResponseErrorBuild(c, http.StatusBadRequest, fmt.Sprintf("track type is not valid at index %d", i))
			return
		}
		if !utils.ValidatorContains(configs.AppsSources, track.AppsSource) {
			utils.MessageResponseErrorBuild(c, http.StatusBadRequest, fmt.Sprintf("app source is not valid at index %d", i))
			return
		}
	}

	// Service : Create Track Multi
	if err := tr.TrackService.CreateTrackMulti(req); err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Response
	utils.MessageResponseBuild(c, "success", "track", "post", http.StatusCreated, req, nil)
}

// @Summary      Get All Track
// @Description  Returns a list of track in pagination format
// @Tags         Track
// @Accept       json
// @Produce      json
// @Success      200  {object}  entities.ResponseGetAllTrack
// @Failure      404  {object}  entities.ResponseNotFound
// @Router       /api/v1/tracks/{app_source}/{created_by} [get]
// @Param        created_by  path  string  true  "created_by must be UUID"
// @Param        app_source  path  string  true  "app_source (such as: pinmarker, mi-fik, myride, or kumande)"
func (tr *TrackController) GetAllTrack(c *gin.Context) {
	// Param
	appsSource := c.Param("app_source")
	createdByRaw := c.Param("created_by")

	// Convert to UUID
	createdBy, err := uuid.Parse(createdByRaw)
	if err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "created by is not valid")
		return
	}

	// Validator : App Source
	if !utils.ValidatorContains(configs.AppsSources, appsSource) {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "app source is not valid")
		return
	}

	// Pagination
	pagination := utils.PaginationBuilder(c)

	// Service : Get All Track
	track, total, err := tr.TrackService.GetAllTrack(pagination, appsSource, createdBy)
	if err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, err.Error())
	}

	// Response
	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	metadata := gin.H{
		"total":       total,
		"page":        pagination.Page,
		"limit":       pagination.Limit,
		"total_pages": totalPages,
	}
	utils.MessageResponseBuild(c, "success", "track", "get", http.StatusOK, track, metadata)
}

// @Summary      Delete Track By ID
// @Description  Delete track by given id
// @Tags         Track
// @Accept       json
// @Produce      json
// @Success      200  {object}  entities.ResponseDeleteTrackById
// @Failure      404  {object}  entities.ResponseNotFound
// @Router       /api/v1/tracks/{app_source}/{created_by}/{track_id} [delete]
// @Param        created_by  path  string  true  "created_by must be UUID"
// @Param        app_source  path  string  true  "app_source (such as: pinmarker, mi-fik, myride, or kumande)"
// @Param        track_id  path  string  true  "track_id must be UUID"
func (tr *TrackController) DeleteTrackById(c *gin.Context) {
	// Param
	appsSource := c.Param("app_source")
	createdByRaw := c.Param("created_by")
	trackIdRaw := c.Param("track_id")

	// Convert to UUID
	createdBy, err := uuid.Parse(createdByRaw)
	if err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "created by is not valid")
		return
	}
	trackID, err := uuid.Parse(trackIdRaw)
	if err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "track id is not valid")
		return
	}

	// Validator : App Source
	if !utils.ValidatorContains(configs.AppsSources, appsSource) {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "app source is not valid")
		return
	}

	// Service : Get All Track
	err = tr.TrackService.DeleteTrackByID(appsSource, createdBy, trackID)
	if err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.MessageResponseBuild(c, "success", "track", "soft delete", http.StatusOK, nil, nil)
}
