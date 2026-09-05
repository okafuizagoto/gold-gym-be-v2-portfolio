package goldgym

import (
	"context"

	goldMejaEntity "gold-gym-be/internal/entity/meja"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	InsertMeja(ctx context.Context, rows []goldMejaEntity.Meja) error
	GetMejaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldMejaEntity.Meja, error)
	GetExistingMejaNames(ctx context.Context, outcode string) (map[string]bool, error)
	UpdateMejaStatus(ctx context.Context, outcode string, mejaIDs []int, fromStatus, toStatus string) (int64, error)
	ReserveMeja(ctx context.Context, outcode string, mejaIDs []int) (int64, error)
}

// Service ...
type Service struct {
	goldgymmeja Data
	tracer      opentracing.Tracer
	logger      jaegerLog.Factory
}

// New ...
func New(goldgymMejaData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	return Service{
		goldgymmeja: goldgymMejaData,
		tracer:      tracer,
		logger:      logger,
	}
}
