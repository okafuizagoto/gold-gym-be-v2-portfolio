package firebaseclient

import (
	"context"
	"encoding/json"
	"gold-gym-be/internal/config"
	"gold-gym-be/pkg/errors"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

// Client ...
type Client struct {
	App             *firebase.App
	FirestoreClient *firestore.Client
	StorageClient   *storage.Client
}

// NewClient ...
func NewClient(cfg *config.Config, credentials map[string]string) (*Client, error) {
	var (
		c   Client
		ctx context.Context
	)

	ctx = context.Background()
	//
	cb, err := json.Marshal(credentials)
	if err != nil {
		return &c, errors.Wrap(err, "[FIREBASE] Failed to marshal credentials!")
	}
	//
	option := option.WithCredentialsJSON(cb)
	config := &firebase.Config{
		ProjectID:     cfg.Firebase.ProjectID,
		StorageBucket: cfg.Firebase.ProjectID + ".appspot.com",
	}
	//
	c.App, err = firebase.NewApp(ctx, config, option)
	if err != nil {
		return &c, errors.Wrap(err, "[FIREBASE] Failed to initiate firebase client!")
	}
	//
	c.FirestoreClient, err = c.App.Firestore(ctx)
	if err != nil {
		return &c, errors.Wrap(err, "[FIREBASE] Failed to initiate firestore client!")
	}
	//
	c.StorageClient, err = c.App.Storage(ctx)
	if err != nil {
		return &c, errors.Wrap(err, "[FIREBASE] Failed to initiate storage client!")
	}
	return &c, err
}
