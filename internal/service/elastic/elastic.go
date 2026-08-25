package elastic

import (
	"context"

	elasticEntity "gold-gym-be/internal/entity/elastic"

	"github.com/opentracing/opentracing-go"
	jaegerLog "gold-gym-be/pkg/log"
)

// RepoData defines the data layer interface consumed by this service
type RepoData interface {
	IndexDocument(ctx context.Context, index string, doc elasticEntity.UserDocument) (string, error)
	SearchDocuments(ctx context.Context, index string, query string) ([]elasticEntity.UserDocument, error)
	GetDocumentByID(ctx context.Context, index string, id string) (elasticEntity.UserDocument, error)
	UpdateDocument(ctx context.Context, index string, id string, doc elasticEntity.UserDocument) error
	DeleteDocument(ctx context.Context, index string, id string) error
}

// Service holds the data layer dependency
type Service struct {
	elastic RepoData
	tracer  opentracing.Tracer
	logger  jaegerLog.Factory
}

// New creates a new elastic Service
func New(elasticData RepoData, tracer opentracing.Tracer, logger jaegerLog.Factory) *Service {
	return &Service{
		elastic: elasticData,
		tracer:  tracer,
		logger:  logger,
	}
}
