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
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
	}
)
