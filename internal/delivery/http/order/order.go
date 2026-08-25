package order

import (
	"context"

	orderEntity "gold-gym-be/internal/entity/order"
	saleEntity "gold-gym-be/internal/entity/sales"
	"gold-gym-be/pkg/kafka"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcOrder interface {
	GetPublicOutlets(ctx context.Context, name string) ([]orderEntity.PublicOutlet, error)
	GetAllOutletsForAdmin(ctx context.Context, name string) ([]orderEntity.PublicOutlet, error)
	AddVisibleOutlet(ctx context.Context, goldid int, outcode, addedBy string) error
	RemoveVisibleOutlet(ctx context.Context, goldid int, outcode string) error
	GetOutletCatalog(ctx context.Context, goldid int, outcode, name string) ([]orderEntity.CatalogItem, error)
	InsertOrder(ctx context.Context, buyerID int, req orderEntity.InsertOrderRequest) (orderEntity.Order, error)
	GetBuyerOrders(ctx context.Context, buyerID int) ([]orderEntity.Order, error)
	GetSellerOrders(ctx context.Context, goldid int, status string) ([]orderEntity.Order, error)
	GetOrderDetail(ctx context.Context, orderID string, requesterID int) (orderEntity.OrderWithDetail, error)
	ConfirmOrder(ctx context.Context, orderID string, sellerGoldID int) (orderEntity.Order, error)
	RejectOrder(ctx context.Context, orderID string, sellerGoldID int, reason string) (orderEntity.Order, error)
	FinishOrder(ctx context.Context, orderID string, sellerGoldID int, creator string) (orderEntity.Order, *saleEntity.InsertSaleData, error)
}

type Handler struct {
	orderSvc  IgoldgymSvcOrder
	kafkaProd *kafka.Producer
	tracer    opentracing.Tracer
	logger    jaegerLog.Factory
}

// New bridging order handler initialization
func New(is IgoldgymSvcOrder, kp *kafka.Producer, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		orderSvc:  is,
		kafkaProd: kp,
		tracer:    tracer,
		logger:    logger,
	}
}
