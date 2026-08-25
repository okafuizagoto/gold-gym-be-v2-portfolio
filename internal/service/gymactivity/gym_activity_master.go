package gymactivity

import (
	"context"
	"log"
	"time"

	gymEntity "gold-gym-be/internal/entity/gymactivity"

	"github.com/pkg/errors"
)

// GetActivities retrieves gym activities, optionally filtered by memberID (0 = all)
func (s *Service) GetActivities(ctx context.Context, memberID int) ([]gymEntity.GymActivity, error) {
	log.Println("service GetActivities")

	activities, err := s.repo.GetActivities(ctx, memberID)
	if err != nil {
		return nil, errors.Wrap(err, "[Service][GetActivities]")
	}

	return activities, nil
}

// InsertActivity persists a new gym activity document and returns the inserted document ID
func (s *Service) InsertActivity(ctx context.Context, req gymEntity.InsertActivityRequest) (string, error) {
	log.Println("service InsertActivity")

	activity := gymEntity.GymActivity{
		MemberID:     req.MemberID,
		MemberEmail:  req.MemberEmail,
		ActivityType: req.ActivityType,
		Description:  req.Description,
		Duration:     req.Duration,
		CreatedAt:    time.Now().UTC(),
	}

	insertedID, err := s.repo.InsertActivity(ctx, activity)
	if err != nil {
		return "", errors.Wrap(err, "[Service][InsertActivity]")
	}

	return insertedID, nil
}

// UpdateActivity updates an existing gym activity document by its hex ObjectID
func (s *Service) UpdateActivity(ctx context.Context, id string, req gymEntity.UpdateActivityRequest) error {
	log.Println("service UpdateActivity", id)

	if err := s.repo.UpdateActivity(ctx, id, req); err != nil {
		return errors.Wrap(err, "[Service][UpdateActivity]")
	}

	return nil
}

// DeleteActivity removes a gym activity document by its hex ObjectID
func (s *Service) DeleteActivity(ctx context.Context, id string) error {
	log.Println("service DeleteActivity", id)

	if err := s.repo.DeleteActivity(ctx, id); err != nil {
		return errors.Wrap(err, "[Service][DeleteActivity]")
	}

	return nil
}
