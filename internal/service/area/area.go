package goldgym

import (
	"context"

	goldAreaEntity "gold-gym-be/internal/entity/area"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	InsertArea(ctx context.Context, goldid int, outcode string, areas []goldAreaEntity.InsertArea) error
	GetAreaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldAreaEntity.Area, error)
}

// Service ...
type Service struct {
	goldgymarea Data
	tracer      opentracing.Tracer
	logger      jaegerLog.Factory
}

// New ...
func New(goldgymAreaData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	return Service{
		goldgymarea: goldgymAreaData,
		tracer:      tracer,
		logger:      logger,
	}
}
