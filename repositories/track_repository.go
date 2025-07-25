package repositories

import (
	"context"
	"fmt"
	"pinmarker/configs"
	"pinmarker/entities"
	"pinmarker/utils"
	"sort"
	"strings"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/google/uuid"
)

// Track Interface
type TrackRepository interface {
	Create(track *entities.Track) error
	CreateBatch(tracks []*entities.Track) error
	FindAll(pagination utils.Pagination, appsSource string, createdBy uuid.UUID) ([]*entities.Track, int, error)
	DeleteByID(appsSource string, createdBy uuid.UUID, trackID uuid.UUID) error
	FindAppsUserTotal() ([]*entities.AppCount, error)
	DeleteAllTracksByDaysCreated(days int) (int64, error)
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

func (r *trackRepository) CreateBatch(tracks []*entities.Track) error {
	// Prepare multi-path data
	updates := make(map[string]interface{})

	for _, track := range tracks {
		// Default Field
		track.ID = uuid.New()

		// Converter : Struct To Map
		data, err := utils.ConverterStructToMap(track)
		if err != nil {
			continue
		}

		// Doc Name
		path := fmt.Sprintf("%s/%s/user_%s/%s", configs.TrackDoc, track.AppsSource, track.CreatedBy.String(), track.ID.String())
		updates[path] = data
	}

	// Query
	ref := r.firebaseClient.NewRef("/")
	if err := ref.Update(r.firebaseCtx, updates); err != nil {
		return fmt.Errorf("failed to batch insert to Firebase: %w", err)
	}

	return nil
}

func (r *trackRepository) FindAll(pagination utils.Pagination, appsSource string, createdBy uuid.UUID) ([]*entities.Track, int, error) {
	// Doc Name
	docName := fmt.Sprintf("%s/%s/user_%s", configs.TrackDoc, appsSource, createdBy.String())
	ref := r.firebaseClient.NewRef(docName)

	// Query
	var result map[string]map[string]interface{}
	if err := ref.Get(r.firebaseCtx, &result); err != nil {
		return nil, 0, fmt.Errorf("failed to read from Firebase: %w", err)
	}

	// Converter : Map To Struct
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

	// Sort Descending
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].CreatedAt.After(tracks[j].CreatedAt)
	})

	// Pagination
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

func (r *trackRepository) DeleteByID(appsSource string, createdBy uuid.UUID, trackID uuid.UUID) error {
	// Doc Name
	docName := fmt.Sprintf("%s/%s/user_%s/%s", configs.TrackDoc, appsSource, createdBy.String(), trackID.String())
	ref := r.firebaseClient.NewRef(docName)

	// Check Existence
	var existing map[string]interface{}
	err := ref.Get(r.firebaseCtx, &existing)
	if err != nil {
		return fmt.Errorf("failed to read before delete: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("Track not found")
	}

	// Query
	if err := ref.Delete(r.firebaseCtx); err != nil {
		return fmt.Errorf("failed to delete from Firebase: %w", err)
	}

	return nil
}

func (r *trackRepository) FindAppsUserTotal() ([]*entities.AppCount, error) {
	// Doc Name
	ref := r.firebaseClient.NewRef(configs.TrackDoc)

	// Query
	var raw map[string]map[string]interface{}
	if err := ref.Get(r.firebaseCtx, &raw); err != nil {
		return nil, fmt.Errorf("failed to read apps from Firebase: %w", err)
	}
	appCounts := make([]*entities.AppCount, 0)

	// Count each app's track
	for appName, users := range raw {
		appCounts = append(appCounts, &entities.AppCount{
			AppName: appName,
			Total:   len(users),
		})
	}

	return appCounts, nil
}

func (r *trackRepository) DeleteAllTracksByDaysCreated(days int) (int64, error) {
	var deletedCount int64

	// Doc Name
	rootRef := r.firebaseClient.NewRef(configs.TrackDoc)

	// Fetch All Tracks Data
	var allApps map[string]map[string]map[string]interface{}
	err := rootRef.Get(r.firebaseCtx, &allApps)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch all tracks: %w", err)
	}

	// Cutoff Time
	cutoff := time.Now().AddDate(0, 0, -days)

	// All Apps
	for appName, users := range allApps {
		for userKey, tracks := range users {
			if !strings.HasPrefix(userKey, "user_") {
				continue
			}

			for trackID, trackDataRaw := range tracks {
				trackData, ok := trackDataRaw.(map[string]interface{})
				if !ok {
					continue
				}

				createdAtStr, ok := trackData["created_at"].(string)
				if !ok {
					continue
				}

				createdAt, err := time.Parse(time.RFC3339Nano, createdAtStr)
				if err != nil {
					continue
				}

				if createdAt.Before(cutoff) {
					path := fmt.Sprintf("%s/%s/%s/%s", configs.TrackDoc, appName, userKey, trackID)
					ref := r.firebaseClient.NewRef(path)

					// Query
					if err := ref.Delete(r.firebaseCtx); err != nil {
						return deletedCount, fmt.Errorf("failed to delete track %s: %w", path, err)
					}

					deletedCount++
				}
			}
		}
	}

	return deletedCount, nil
}
