package elastic

import (
	"context"

	elasticEntity "gold-gym-be/internal/entity/elastic"

	"github.com/opentracing/opentracing-go"
	jaegerLog "gold-gym-be/pkg/log"
)

// IelasticSvc defines the service layer methods used by this handler
type IelasticSvc interface {
	IndexUser(ctx context.Context, index string, doc elasticEntity.UserDocument) (string, error)
	SearchUsers(ctx context.Context, index string, query string) ([]elasticEntity.UserDocument, error)
	GetUserByID(ctx context.Context, index string, id string) (elasticEntity.UserDocument, error)
	UpdateUser(ctx context.Context, index string, id string, doc elasticEntity.UserDocument) error
	DeleteUser(ctx context.Context, index string, id string) error
}

// Handler holds the elastic service dependency
type Handler struct {
	elasticSvc IelasticSvc
	tracer     opentracing.Tracer
	logger     jaegerLog.Factory
}

// New creates a new elastic Handler
func New(svc IelasticSvc, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		elasticSvc: svc,
		tracer:     tracer,
		logger:     logger,
	}
}
