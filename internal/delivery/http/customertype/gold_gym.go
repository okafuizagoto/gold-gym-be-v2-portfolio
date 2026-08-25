package goldgym

import (
	"context"
	jaegerLog "gold-gym-be/pkg/log"

	goldCustomerTypeEntity "gold-gym-be/internal/entity/customertype"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcCustomerType interface {
	GetCustomerType(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldCustomerTypeEntity.CustomerType, goldCustomerTypeEntity.MetadataPaginationDetail, error)
	InsertCustomerType(ctx context.Context, id int, items []goldCustomerTypeEntity.CustomerType) (string, error)
	UpdateCustomerType(ctx context.Context, items goldCustomerTypeEntity.CustomerType) (string, error)
	DeleteCustomerType(ctx context.Context, goldid, golditemid int, outcode string) (string, error)
}

type (
	// Handler ...
	Handler struct {
		goldgymSvcCustomerType IgoldgymSvcCustomerType
		tracer                 opentracing.Tracer
		logger                 jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(isim IgoldgymSvcCustomerType, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvcCustomerType: isim,
		tracer:                 tracer,
		logger:                 logger,
	}
}
