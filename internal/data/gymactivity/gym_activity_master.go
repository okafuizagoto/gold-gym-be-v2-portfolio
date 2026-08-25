package gymactivity

import (
	"context"
	"time"

	gymEntity "gold-gym-be/internal/entity/gymactivity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GetActivities fetches gym activity documents. Pass memberID=0 to get all.
func (r *Repository) GetActivities(ctx context.Context, memberID int) ([]gymEntity.GymActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := r.db.Collection(gymEntity.CollectionName)

	filter := bson.D{}
	if memberID > 0 {
		filter = bson.D{{Key: "member_id", Value: memberID}}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(100)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activities []gymEntity.GymActivity
	if err := cursor.All(ctx, &activities); err != nil {
		return nil, err
	}

	return activities, nil
}

// InsertActivity inserts a new GymActivity document and returns the inserted ID as string
func (r *Repository) InsertActivity(ctx context.Context, activity gymEntity.GymActivity) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	coll := r.db.Collection(gymEntity.CollectionName)

	result, err := coll.InsertOne(ctx, activity)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		return oid.Hex(), nil
	}

	return "", nil
}

// UpdateActivity updates fields of an existing GymActivity document by its hex ID.
// Returns mongo.ErrNoDocuments (wrapped) when no document matches.
func (r *Repository) UpdateActivity(ctx context.Context, id string, req gymEntity.UpdateActivityRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "activity_type", Value: req.ActivityType},
		{Key: "description", Value: req.Description},
		{Key: "duration_minutes", Value: req.Duration},
	}}}

	coll := r.db.Collection(gymEntity.CollectionName)
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errNotFound
	}

	return nil
}

// DeleteActivity removes a GymActivity document by its hex ID.
// Returns errNotFound when no document matches.
func (r *Repository) DeleteActivity(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coll := r.db.Collection(gymEntity.CollectionName)
	result, err := coll.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errNotFound
	}

	return nil
}
