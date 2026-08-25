package gymactivity

import (
	"context"

	gymEntity "gold-gym-be/internal/entity/gymactivity"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// IgymactivitySvc defines the service interface consumed by this handler
type IgymactivitySvc interface {
	GetActivities(ctx context.Context, memberID int) ([]gymEntity.GymActivity, error)
	InsertActivity(ctx context.Context, req gymEntity.InsertActivityRequest) (string, error)
	UpdateActivity(ctx context.Context, id string, req gymEntity.UpdateActivityRequest) error
	DeleteActivity(ctx context.Context, id string) error
}

// Handler holds dependencies for gymactivity HTTP handlers
type Handler struct {
	gymactivitySvc IgymactivitySvc
	tracer         opentracing.Tracer
	logger         jaegerLog.Factory
}

// New creates a new gymactivity Handler
func New(svc IgymactivitySvc, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		gymactivitySvc: svc,
		tracer:         tracer,
		logger:         logger,
	}
}
