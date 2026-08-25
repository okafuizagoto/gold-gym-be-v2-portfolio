package order

import (
	"context"

	orderEntity "gold-gym-be/internal/entity/order"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// Data kontrak akses data pesanan pembeli.
type Data interface {
	GetPublicOutlets(ctx context.Context, name string) ([]orderEntity.PublicOutlet, error)
	GetAllOutletsForAdmin(ctx context.Context, name string) ([]orderEntity.PublicOutlet, error)
	AddVisibleOutlet(ctx context.Context, goldid int, outcode, addedBy string) error
	RemoveVisibleOutlet(ctx context.Context, goldid int, outcode string) error
	GetOutletCatalog(ctx context.Context, goldid int, outcode, name string) ([]orderEntity.CatalogItem, error)
	GetOutletByCode(ctx context.Context, goldid int, outcode string) (*orderEntity.PublicOutlet, error)
	GetBuyerName(ctx context.Context, buyerID int) (string, error)
	InsertOrder(ctx context.Context, header orderEntity.Order, details []orderEntity.OrderDetail) error
	GetOrdersByBuyer(ctx context.Context, buyerID int) ([]orderEntity.Order, error)
	GetOrdersBySeller(ctx context.Context, goldid int, status string) ([]orderEntity.Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*orderEntity.Order, error)
	GetOrderDetails(ctx context.Context, orderID string) ([]orderEntity.OrderDetail, error)
	UpdateOrderStatus(ctx context.Context, orderID, status string, reason *string) error
	FinishOrder(ctx context.Context, orderID, saleID string) error
}

// Service ...
type Service struct {
	order  Data
	tracer opentracing.Tracer
	logger jaegerLog.Factory
}

// New ...
func New(orderData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	return Service{
		order:  orderData,
		tracer: tracer,
		logger: logger,
	}
}
