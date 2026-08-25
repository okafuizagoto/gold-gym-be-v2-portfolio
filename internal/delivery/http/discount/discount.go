package discount

import (
	"context"

	discountEntity "gold-gym-be/internal/entity/discount"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvcDiscount interface {
	GetItemsForOutlet(ctx context.Context, goldid int, outcode, name string) ([]discountEntity.ItemForOutlet, error)
	GetDiscounts(ctx context.Context, goldid int, outcode, name string, page, length int) ([]discountEntity.Discount, discountEntity.MetadataPaginationDetail, error)
	GetActiveDiscountsByOutlet(ctx context.Context, goldid int, outcode string) ([]discountEntity.Discount, error)
	InsertDiscount(ctx context.Context, goldid int, role string, items []discountEntity.InsertDiscount) (string, error)
	UpdateDiscount(ctx context.Context, goldid int, role string, req discountEntity.UpdateDiscount) (string, error)
	DeleteDiscount(ctx context.Context, goldid int, role string, discountID int) (string, error)
	GetDiscountHistory(ctx context.Context, discountID int, page, length int) ([]discountEntity.DiscountHistory, discountEntity.MetadataPaginationDetail, error)

	GenerateVoucherCode(ctx context.Context, outcode string) (string, error)
	InsertVoucher(ctx context.Context, goldid int, role string, req discountEntity.InsertVoucher) (string, error)
	GetVouchers(ctx context.Context, goldid int, outcode string, page, length int) ([]discountEntity.Voucher, discountEntity.MetadataPaginationDetail, error)
	DeleteVoucher(ctx context.Context, goldid int, role string, voucherID int) (string, error)
	GetVoucherHistory(ctx context.Context, goldid int, outcode string, page, length int) ([]discountEntity.VoucherHistory, discountEntity.MetadataPaginationDetail, error)
	CheckVoucher(ctx context.Context, goldid int, outcode, code string) (*discountEntity.Voucher, error)
}

type Handler struct {
	discountSvc IgoldgymSvcDiscount
	tracer      opentracing.Tracer
	logger      jaegerLog.Factory
}

// New for bridging discount handler initialization
func New(is IgoldgymSvcDiscount, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		discountSvc: is,
		tracer:      tracer,
		logger:      logger,
	}
}
