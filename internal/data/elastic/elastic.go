package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/opentracing/opentracing-go"

	jaegerLog "gold-gym-be/pkg/log"
)

// Repository holds ES client and dependencies
type Repository struct {
	client *elasticsearch.Client
	tracer opentracing.Tracer
	logger jaegerLog.Factory
}

// New creates a new Elasticsearch Repository
func New(client *elasticsearch.Client, tracer opentracing.Tracer, logger jaegerLog.Factory) *Repository {
	return &Repository{
		client: client,
		tracer: tracer,
		logger: logger,
	}
}
