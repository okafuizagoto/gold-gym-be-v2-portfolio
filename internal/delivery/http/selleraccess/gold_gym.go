package selleraccess

import (
	"context"

	entity "gold-gym-be/internal/entity/selleraccess"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcSellerAccess interface {
	GetAll(ctx context.Context, name string) ([]entity.SellerMenuAccess, error)
	SetDaftarPembeli(ctx context.Context, goldID int, active bool) error
	SetModePembeli(ctx context.Context, goldID int, active bool) error
}

type Handler struct {
	svc    IgoldgymSvcSellerAccess
	tracer opentracing.Tracer
	logger jaegerLog.Factory
}

// New bridging selleraccess handler initialization
func New(is IgoldgymSvcSellerAccess, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		svc:    is,
		tracer: tracer,
		logger: logger,
	}
}
