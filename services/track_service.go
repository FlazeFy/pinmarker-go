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
	GetAppsUserTotal() ([]*entities.AppCount, error)
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
	return s.trackRepo.Create(track)
}

func (s *trackService) CreateTrackMulti(track []*entities.Track) error {
	return s.trackRepo.CreateBatch(track)
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
	return s.trackRepo.DeleteByID(appsSource, createdBy, trackID)
}

func (s *trackService) GetAppsUserTotal() ([]*entities.AppCount, error) {
	return s.trackRepo.FindAppsUserTotal()
}
