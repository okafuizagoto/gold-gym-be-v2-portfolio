package goldgym

import (
	"context"

	goldMejaEntity "gold-gym-be/internal/entity/meja"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcMeja interface {
	InsertMejaBulk(ctx context.Context, goldid int, outcode string, rows []goldMejaEntity.InsertMeja) error
	GetMejaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldMejaEntity.Meja, error)
	ReserveMeja(ctx context.Context, outcode string, mejaIDs []int) error
	ReleaseMeja(ctx context.Context, outcode string, mejaIDs []int) error
}

type (
	// Handler ...
	Handler struct {
		goldgymSvcMeja IgoldgymSvcMeja
		tracer         opentracing.Tracer
		logger         jaegerLog.Factory
	}
)

// New for bridging meja handler initialization
func New(isst IgoldgymSvcMeja, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvcMeja: isst,
		tracer:         tracer,
		logger:         logger,
	}
}
