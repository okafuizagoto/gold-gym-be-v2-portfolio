package booking

import (
	"context"
	"time"

	bookingEntity "gold-gym-be/internal/entity/booking"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// Data ...
type Data interface {
	ExpireOverdueBookings(ctx context.Context, now time.Time) error
	GetBookingsByDate(ctx context.Context, goldid int, outcode string, date time.Time) ([]bookingEntity.Booking, error)
	CountActiveInWindow(ctx context.Context, goldid int, outcode string, winStart time.Time) (int, error)
	InsertBooking(ctx context.Context, b bookingEntity.Booking) error
	GetBookingByID(ctx context.Context, bookingID string) (*bookingEntity.Booking, error)
	UpdateBookingPaid(ctx context.Context, bookingID string, saleID string, itemID int, itemName string, price string) error
	UpdateBookingStatus(ctx context.Context, bookingID string, status string) error
	SearchBuyers(ctx context.Context, name string, limit int) ([]bookingEntity.BuyerInfo, error)
	GetBuyerByID(ctx context.Context, goldid int) (*bookingEntity.BuyerInfo, error)
	GetTherapyItem(ctx context.Context, itemID int) (*bookingEntity.TherapyItem, error)
	GetTherapyItemByName(ctx context.Context, outcode string, name string) (*bookingEntity.TherapyItem, error)
	FindAdjacentUnpaid(ctx context.Context, goldid int, outcode string, custName string, therapyType string, excludeID string, start time.Time) (*bookingEntity.Booking, error)
	RemoveBookingWithLog(ctx context.Context, b bookingEntity.Booking, removedBy string, removedRole string) error
	GetOutletByCode(ctx context.Context, outcode string) (*bookingEntity.OutletInfo, error)
	GetStockIDByItem(ctx context.Context, outcode string, itemID int) (string, error)
}

// Service ...
type Service struct {
	booking Data
	tracer  opentracing.Tracer
	logger  jaegerLog.Factory
}

// New ...
func New(bookingData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	return Service{
		booking: bookingData,
		tracer:  tracer,
		logger:  logger,
	}
}
