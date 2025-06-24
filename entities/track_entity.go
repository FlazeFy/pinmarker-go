package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	Track struct {
		ID               uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
		BatteryIndicator int       `json:"battery_indicator" gorm:"type:varchar(144);not null"`
		TrackLat         string    `json:"track_lat" gorm:"type:varchar(255);not null"`
		TrackLong        string    `json:"track_long" gorm:"type:varchar(255);not null"`
		TrackType        string    `json:"track_type" gorm:"type:varchar(36);not null"`
		AppsSource       string    `json:"app_source" gorm:"type:varchar(36);not null"`
		CreatedAt        time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		CreatedBy        uuid.UUID `json:"created_by" gorm:"not null"`
	}
	// For Response
	ResponseCreateTrack struct {
		Message string `json:"message" example:"Track created"`
		Status  string `json:"status" example:"success"`
		Data    Track  `json:"data"`
	}
	ResponseCreateTrackMulti struct {
		Message string  `json:"message" example:"Track created"`
		Status  string  `json:"status" example:"success"`
		Data    []Track `json:"data"`
	}
	ResponseGetAllTrack struct {
		Message string  `json:"message" example:"Track fetched"`
		Status  string  `json:"status" example:"success"`
		Data    []Track `json:"data"`
	}
	ResponseDeleteTrackById struct {
		Message string `json:"message" example:"Track permanentally deleted"`
		Status  string `json:"status" example:"success"`
	}
	// For Request
	RequestCreateTrack struct {
		BatteryIndicator int       `json:"battery_indicator" example:"85"`
		TrackLat         string    `json:"track_lat" example:"-6.200000"`
		TrackLong        string    `json:"track_long" example:"106.816666"`
		TrackType        string    `json:"track_type" example:"live"`
		AppsSource       string    `json:"app_source" example:"pinmarker"`
		CreatedBy        uuid.UUID `json:"created_by" example:"123e4567-e89b-12d3-a456-426614174000"`
	}
	RequestCreateTrackMulti []struct {
		BatteryIndicator int       `json:"battery_indicator" example:"85"`
		TrackLat         string    `json:"track_lat" example:"-6.200000"`
		TrackLong        string    `json:"track_long" example:"106.816666"`
		TrackType        string    `json:"track_type" example:"live"`
		AppsSource       string    `json:"app_source" example:"pinmarker"`
		CreatedAt        time.Time `json:"created_at" example:"2025-06-23T11:30:15.913505+07:00"`
		CreatedBy        uuid.UUID `json:"created_by" example:"123e4567-e89b-12d3-a456-426614174000"`
	}
)
