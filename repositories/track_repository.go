package repositories

import (
	"context"
	"fmt"
	"pinmarker/configs"
	"pinmarker/entities"
	"pinmarker/utils"
	"sort"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/google/uuid"
)

// Track Interface
type TrackRepository interface {
	Create(track *entities.Track) error
	FindAll(pagination utils.Pagination, appsSource string, createdBy uuid.UUID) ([]*entities.Track, int, error)
}

// Track Struct
type trackRepository struct {
	firebaseClient *db.Client
	firebaseCtx    context.Context
}

// Track Constructor
func NewTrackRepository() TrackRepository {
	client, ctx, err := configs.FirebaseDB()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize Firebase DB: %v", err))
	}

	return &trackRepository{
		firebaseClient: client,
		firebaseCtx:    ctx,
	}
}

func (r *trackRepository) Create(track *entities.Track) error {
	// Default Field
	track.ID = uuid.New()
	track.CreatedAt = time.Now()

	// Converter : Struct To Map
	data, err := utils.ConverterStructToMap(track)
	if err != nil {
		return err
	}

	// Doc Name
	docName := configs.TrackDoc + "/" + track.AppsSource + "/user_" + track.CreatedBy.String()

	// Query
	ref := r.firebaseClient.NewRef(docName).Child(track.ID.String())
	if err := ref.Set(r.firebaseCtx, data); err != nil {
		return fmt.Errorf("failed to save to Firebase: %w", err)
	}

	return nil
}

func (r *trackRepository) FindAll(pagination utils.Pagination, appsSource string, createdBy uuid.UUID) ([]*entities.Track, int, error) {
	// Firebase path
	docName := fmt.Sprintf("%s/%s/user_%s", configs.TrackDoc, appsSource, createdBy.String())
	ref := r.firebaseClient.NewRef(docName)

	// Fetch from Firebase
	var result map[string]map[string]interface{}
	if err := ref.Get(r.firebaseCtx, &result); err != nil {
		return nil, 0, fmt.Errorf("failed to read from Firebase: %w", err)
	}

	// Convert map to []*Track
	tracks := make([]*entities.Track, 0)
	for _, item := range result {
		var track entities.Track
		if err := utils.ConverterMapToStruct(item, &track); err != nil {
			continue
		}
		tracks = append(tracks, &track)
	}

	// Total before pagination
	total := len(tracks)

	// Sort DESC by CreatedAt
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].CreatedAt.After(tracks[j].CreatedAt)
	})

	// Manual pagination
	start := (pagination.Page - 1) * pagination.Limit
	end := start + pagination.Limit

	if start > total {
		return []*entities.Track{}, total, nil
	}
	if end > total {
		end = total
	}

	return tracks[start:end], total, nil
}
