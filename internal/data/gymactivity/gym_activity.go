package gymactivity

import (
	"errors"

	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// errNotFound is returned when a document is not found in MongoDB
var errNotFound = errors.New("activity not found")

// Repository holds a MongoDB database handle and observability dependencies
type Repository struct {
	db     *mongo.Database
	tracer opentracing.Tracer
	logger jaegerLog.Factory
}

// New creates a new gymactivity Repository
func New(db *mongo.Database, tracer opentracing.Tracer, logger jaegerLog.Factory) *Repository {
	return &Repository{
		db:     db,
		tracer: tracer,
		logger: logger,
	}
}
