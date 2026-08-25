package goldgym

import (
	"context"
	"gold-gym-be/pkg/kafka"
	jaegerLog "gold-gym-be/pkg/log"

	goldCustomerEntity "gold-gym-be/internal/entity/customer"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcCustomer interface {
	GetCustomer(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldCustomerEntity.Customer, goldCustomerEntity.MetadataPaginationDetail, error)
	InsertCustomer(ctx context.Context, id int, items []goldCustomerEntity.Customer) (string, error)
	UpdateCustomer(ctx context.Context, items goldCustomerEntity.Customer) (string, error)
	DeleteCustomer(ctx context.Context, goldid, golditemid int, outcode string) (string, error)
}

type (
	// Handler ...
	Handler struct {
		goldgymSvcCustomer IgoldgymSvcCustomer
		kafkaProd          *kafka.Producer
		tracer             opentracing.Tracer
		logger             jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(isim IgoldgymSvcCustomer, kp *kafka.Producer, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvcCustomer: isim,
		kafkaProd:          kp,
		tracer:             tracer,
		logger:             logger,
	}
}
