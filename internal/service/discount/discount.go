package goldgym

import (
	"context"

	discountEntity "gold-gym-be/internal/entity/discount"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// Data ...
type Data interface {
	GetItemsForOutlet(ctx context.Context, goldid int, outcode, name string) ([]discountEntity.ItemForOutlet, error)
	GetDiscounts(ctx context.Context, goldid int, outcode, name string, page, length int) ([]discountEntity.Discount, error)
	GetTotalDiscounts(ctx context.Context, goldid int, outcode, name string) (int64, error)
	GetActiveDiscountsByOutlet(ctx context.Context, goldid int, outcode string) ([]discountEntity.Discount, error)
	GetDiscountByID(ctx context.Context, goldid int, discountID int) (*discountEntity.Discount, error)
	GetGoldNameByID(ctx context.Context, goldid int) (string, error)
	InsertDiscountWithLog(ctx context.Context, discount discountEntity.Discount, actorName, actorRole string) (int, error)
	UpdateDiscountWithLog(ctx context.Context, existing discountEntity.Discount, newType string, newValue float64, newStatus string, actorName, actorRole string) error
	DeleteDiscountWithLog(ctx context.Context, existing discountEntity.Discount, actorName, actorRole string) error
	GetDiscountHistory(ctx context.Context, discountID int, page, length int) ([]discountEntity.DiscountHistory, error)
	GetTotalDiscountHistory(ctx context.Context, discountID int) (int64, error)
	GetActiveTotalDiscount(ctx context.Context, goldid int, outcode string) (*discountEntity.Discount, error)

	VoucherCodeExists(ctx context.Context, outcode, code string) (bool, error)
	GetVoucherByCode(ctx context.Context, goldid int, outcode, code string) (*discountEntity.Voucher, error)
	InsertVoucher(ctx context.Context, voucher discountEntity.Voucher) error
	GetVouchers(ctx context.Context, goldid int, outcode string, page, length int) ([]discountEntity.Voucher, error)
	GetTotalVouchers(ctx context.Context, goldid int, outcode string) (int64, error)
	GetVoucherByID(ctx context.Context, goldid, voucherID int) (*discountEntity.Voucher, error)
	DeleteVoucherWithLog(ctx context.Context, v discountEntity.Voucher, actorName, actorRole string) error
	RedeemVoucherWithLog(ctx context.Context, goldid int, outcode, code, saleID, actorName, actorRole string) (float64, error)
	GetVoucherHistory(ctx context.Context, goldid int, outcode string, page, length int) ([]discountEntity.VoucherHistory, error)
	GetTotalVoucherHistory(ctx context.Context, goldid int, outcode string) (int64, error)
}

// Service ...
type Service struct {
	discount Data
	tracer   opentracing.Tracer
	logger   jaegerLog.Factory
}

// New ...
func New(discountData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	return Service{
		discount: discountData,
		tracer:   tracer,
		logger:   logger,
	}
}
