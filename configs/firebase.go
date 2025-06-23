package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebaseApp() {
	configFileName := os.Getenv("FIREBASE_CONFIG_FILENAME")

	opt := option.WithCredentialsFile("configs/" + configFileName)

	conf := &firebase.Config{
		DatabaseURL: os.Getenv("FIREBASE_DB_URL"),
	}

	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		log.Fatalf("Firebase init error: %v\n", err)
	}
	FirebaseApp = app
}

func FirebaseDB() (*db.Client, context.Context, error) {
	ctx := context.Background()
	client, err := FirebaseApp.Database(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Firebase DB: %w", err)
	}

	return client, ctx, nil
}
