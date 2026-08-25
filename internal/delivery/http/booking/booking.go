package booking

import (
	"context"

	bookingEntity "gold-gym-be/internal/entity/booking"
	saleEntity "gold-gym-be/internal/entity/sales"
	"gold-gym-be/pkg/kafka"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcBooking interface {
	GetSlots(ctx context.Context, outcode string, dateStr string) ([]bookingEntity.Slot, error)
	InsertBooking(ctx context.Context, userID int, role string, creator string, req bookingEntity.InsertBookingRequest) (bookingEntity.Booking, *saleEntity.InsertSaleData, error)
	PayBooking(ctx context.Context, creator string, role string, bookingID string, itemID int, customPrice int) (bookingEntity.Booking, *saleEntity.InsertSaleData, error)
	RemoveBooking(ctx context.Context, role string, remover string, bookingID string) error
	SearchBuyers(ctx context.Context, name string) ([]bookingEntity.BuyerInfo, error)
}

type Handler struct {
	bookingSvc IgoldgymSvcBooking
	kafkaProd  *kafka.Producer
	tracer     opentracing.Tracer
	logger     jaegerLog.Factory
}

// New for bridging booking handler initialization
func New(is IgoldgymSvcBooking, kp *kafka.Producer, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		bookingSvc: is,
		kafkaProd:  kp,
		tracer:     tracer,
		logger:     logger,
	}
}
