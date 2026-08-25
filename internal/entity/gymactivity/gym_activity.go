package gymactivity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const CollectionName = "gym_activities"

// GymActivity represents a gym member's activity log stored in MongoDB
type GymActivity struct {
	ID           bson.ObjectID `bson:"_id,omitempty"    json:"id"`
	MemberID     int                `bson:"member_id"        json:"member_id"`
	MemberEmail  string             `bson:"member_email"     json:"member_email"`
	ActivityType string             `bson:"activity_type"    json:"activity_type"` // "checkin", "checkout", "exercise"
	Description  string             `bson:"description"      json:"description"`
	Duration     int                `bson:"duration_minutes" json:"duration_minutes"`
	CreatedAt    time.Time          `bson:"created_at"       json:"created_at"`
}

// InsertActivityRequest is the request body for POST /gold-gym/v2/activity
type InsertActivityRequest struct {
	MemberID     int    `json:"member_id"        binding:"required"`
	MemberEmail  string `json:"member_email"     binding:"required"`
	ActivityType string `json:"activity_type"    binding:"required"`
	Description  string `json:"description"`
	Duration     int    `json:"duration_minutes"`
}

// UpdateActivityRequest is the request body for PUT /gold-gym/v2/activity?id=<objectid>
type UpdateActivityRequest struct {
	ActivityType string `json:"activity_type"`
	Description  string `json:"description"`
	Duration     int    `json:"duration_minutes"`
}
