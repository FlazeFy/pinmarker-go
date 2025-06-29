package services

import (
	"errors"
	"pinmarker/entities"
	"pinmarker/repositories"
	"pinmarker/utils"

	"github.com/google/uuid"
)

// Track Interface
type TrackService interface {
	CreateTrack(track *entities.Track) error
	CreateTrackMulti(track []*entities.Track) error
	GetAllTrack(pagination utils.Pagination, appsSource string, createdBy uuid.UUID) ([]*entities.Track, int, error)
	DeleteTrackByID(appsSource string, createdBy uuid.UUID, trackID uuid.UUID) error
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

func (s *trackService) CreateTrack(track *entities.Track) error {
	// Repo : Create Track
	if err := s.trackRepo.Create(track); err != nil {
		return err
	}

	return nil
}

func (s *trackService) CreateTrackMulti(track []*entities.Track) error {
	// Repo : Create Track
	if err := s.trackRepo.CreateBatch(track); err != nil {
		return err
	}

	return nil
}

func (s *trackService) GetAllTrack(pagination utils.Pagination, appsSource string, createdBy uuid.UUID) ([]*entities.Track, int, error) {
	// Repo : Get All Track
	track, total, err := s.trackRepo.FindAll(pagination, appsSource, createdBy)
	if err != nil {
		return nil, 0, err
	}
	if track == nil {
		return nil, 0, errors.New("tracks not found")
	}

	return track, total, nil
}

func (s *trackService) DeleteTrackByID(appsSource string, createdBy uuid.UUID, trackID uuid.UUID) error {
	// Repo : Delete Track
	if err := s.trackRepo.DeleteByID(appsSource, createdBy, trackID); err != nil {
		return err
	}

	return nil
}
