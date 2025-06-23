package repositories

import (
	"context"
	"fmt"
	"pinmarker/configs"
	"pinmarker/entities"
	"pinmarker/utils"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/google/uuid"
)

// Track Interface
type TrackRepository interface {
	Create(track *entities.Track) error
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
