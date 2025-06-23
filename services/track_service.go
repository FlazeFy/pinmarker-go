package services

import (
	"pinmarker/entities"
	"pinmarker/repositories"
)

// Track Interface
type TrackService interface {
	Create(track *entities.Track) error
}

// Track Struct
type trackService struct {
	trackRepo repositories.TrackRepository
}

// Track Constructor
func NewTrackService(trackRepo repositories.TrackRepository) TrackService {
	return &trackService{
		trackRepo: trackRepo,
	}
}

func (s *trackService) Create(track *entities.Track) error {
	// Repo : Create Track
	if err := s.trackRepo.Create(track); err != nil {
		return err
	}

	return nil
}
