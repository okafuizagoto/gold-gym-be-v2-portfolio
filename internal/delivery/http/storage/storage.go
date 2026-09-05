package storage

import (
	"context"
	goldStorageEntity "gold-gym-be/internal/entity/storage"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcStorage interface {
	GetSummary(ctx context.Context, goldID int) (goldStorageEntity.StorageSummary, error)
	DeleteEntry(ctx context.Context, sourceType string, sourceID, goldID int, isAdmin bool, deletedBy string) error
}

type Handler struct {
	storageSvc IgoldgymSvcStorage
	tracer     opentracing.Tracer
	logger     jaegerLog.Factory
}

// New for bridging storage handler initialization
func New(is IgoldgymSvcStorage, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		storageSvc: is,
		tracer:     tracer,
		logger:     logger,
	}
}
