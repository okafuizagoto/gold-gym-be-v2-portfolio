package selleraccess

import (
	"context"

	entity "gold-gym-be/internal/entity/selleraccess"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// Data kontrak akses data akses-menu penjual.
type Data interface {
	GetAll(ctx context.Context, name string) ([]entity.SellerMenuAccess, error)
	SetDaftarPembeli(ctx context.Context, goldID int, active bool) error
	SetModePembeli(ctx context.Context, goldID int, active bool) error
}

// Service ...
type Service struct {
	data   Data
	tracer opentracing.Tracer
	logger jaegerLog.Factory
}

// New ...
func New(data Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	return Service{
		data:   data,
		tracer: tracer,
		logger: logger,
	}
}
