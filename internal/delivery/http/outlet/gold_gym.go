package goldgym

import (
	"context"
	goldOutletEntity "gold-gym-be/internal/entity/outlet"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcOutlet interface {
	InsertOutlet(ctx context.Context, outletData []goldOutletEntity.InsertOutlet, goldid int) error
	UpdateOutlet(ctx context.Context, items goldOutletEntity.UpdateOutlet) error
	DeleteOutlet(ctx context.Context, goldid int, code string) error
	GetOutlet(ctx context.Context, goldid int, name, status string, page, length int) ([]goldOutletEntity.Outlet, goldOutletEntity.MetadataPaginationDetail, error)
}

type (
	// Handler ...
	Handler struct {
		goldgymSvcOutlet IgoldgymSvcOutlet
		tracer           opentracing.Tracer
		logger           jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(isst IgoldgymSvcOutlet, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvcOutlet: isst,
		tracer:           tracer,
		logger:           logger,
	}
}
