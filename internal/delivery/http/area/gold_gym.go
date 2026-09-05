package goldgym

import (
	"context"

	goldAreaEntity "gold-gym-be/internal/entity/area"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcArea interface {
	InsertArea(ctx context.Context, goldid int, outcode string, areas []goldAreaEntity.InsertArea) error
	GetAreaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldAreaEntity.Area, error)
}

type (
	// Handler ...
	Handler struct {
		goldgymSvcArea IgoldgymSvcArea
		tracer         opentracing.Tracer
		logger         jaegerLog.Factory
	}
)

// New for bridging area handler initialization
func New(isst IgoldgymSvcArea, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvcArea: isst,
		tracer:         tracer,
		logger:         logger,
	}
}
