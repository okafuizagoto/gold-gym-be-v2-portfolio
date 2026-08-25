package gymactivity

import (
	"context"

	gymEntity "gold-gym-be/internal/entity/gymactivity"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// RepoData defines the data layer contract consumed by this service
type RepoData interface {
	GetActivities(ctx context.Context, memberID int) ([]gymEntity.GymActivity, error)
	InsertActivity(ctx context.Context, activity gymEntity.GymActivity) (string, error)
	UpdateActivity(ctx context.Context, id string, req gymEntity.UpdateActivityRequest) error
	DeleteActivity(ctx context.Context, id string) error
}

// Service holds business logic for gymactivity domain
type Service struct {
	repo   RepoData
	tracer opentracing.Tracer
	logger jaegerLog.Factory
}

// New creates a new gymactivity Service
func New(repo RepoData, tracer opentracing.Tracer, logger jaegerLog.Factory) *Service {
	return &Service{
		repo:   repo,
		tracer: tracer,
		logger: logger,
	}
}
