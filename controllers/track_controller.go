package controllers

import (
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

func (tr *TrackController) Create(c *gin.Context) {
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

	// Validator : Track Type
	if !utils.ValidatorContains(configs.TrackTypes, req.TrackType) {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "track type is not valid")
		return
	}
	if !utils.ValidatorContains(configs.AppsSources, req.AppsSource) {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "app source is not valid")
		return
	}

	// Service : Create Track
	err := tr.TrackService.Create(&req)
	if err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.MessageResponseBuild(c, "success", "track", "post", http.StatusCreated, &req, nil)
}

func (tr *TrackController) GetAllTrack(c *gin.Context) {
	// Param
	appsSource := c.Param("app_source")
	createdByRaw := c.Param("created_by")

	createdBy, err := uuid.Parse(createdByRaw)
	if err != nil {
		utils.MessageResponseErrorBuild(c, http.StatusBadRequest, "created by is not valid")
		return
	}

	// Validator : Track Type
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
