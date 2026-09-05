package goldgym

import (
	"context"
	"gold-gym-be/pkg/kafka"
	jaegerLog "gold-gym-be/pkg/log"

	goldSaleEntity "gold-gym-be/internal/entity/sales"

	"github.com/opentracing/opentracing-go"
)

// IVoucherRedeemer divalidasi & dikonsumsi SEBELUM insert sales diantre ke
// Kafka (bukan di worker async) -- lihat komentar di
// insert_gold_gym_sales_gin.go dan RedeemVoucherWithLog (modul discount).
type IVoucherRedeemer interface {
	RedeemVoucher(ctx context.Context, goldid int, outcode, code, saleID, actorName, actorRole string) (float64, error)
}

type IgoldgymSvcSale interface {
	InsertSales(ctx context.Context, goldid int, sales goldSaleEntity.InsertSales) (string, error)
	GetSales(ctx context.Context, goldid int, custid int, name string, outcode string, page, length int) ([]goldSaleEntity.ThSale, goldSaleEntity.MetadataPaginationDetail, error)
	GetSaleDetail(ctx context.Context, goldid int, custid int, saleid string) (goldSaleEntity.SaleWithDetail, error)
	GetSaleReceiptPDF(ctx context.Context, goldid int, custid int, saleid string) ([]byte, string, error)
	MarkSalePaid(ctx context.Context, goldid int, saleid string) (goldSaleEntity.ThSale, error)
	ResolveOutletGoldID(ctx context.Context, outcode string) (int, error)
	MarkBookingsPaid(ctx context.Context, bookingIDs []string, saleID string) (int64, error)
	ConfirmSaleMeja(ctx context.Context, mejaIDs []int, outcode, saleID string) (int64, error)
	GetMejaNamesByIDs(ctx context.Context, outcode string, mejaIDs []int) ([]string, error)
	SavePaymentProof(ctx context.Context, saleID string, originalName string, mimeType string, content []byte, uploadedBy string, goldID int, isAdmin bool) (goldSaleEntity.PaymentProof, error)
	GetPaymentProofs(ctx context.Context, saleID string) ([]goldSaleEntity.PaymentProof, error)
	GetPaymentProofFile(ctx context.Context, proofID int) (goldSaleEntity.PaymentProof, []byte, error)
	IsPosCustomerRequired(ctx context.Context, goldid int, outcode string) (bool, error)
	GetPosOutletsForAdmin(ctx context.Context, search string) ([]goldSaleEntity.PosOutlet, error)
	SavePosCustomerAccess(ctx context.Context, items []goldSaleEntity.PosAccessItem, addedBy string) error
	GetSalesReport(ctx context.Context, goldid int, outcode, mode, date string) (interface{}, error)

	// Visibilitas fitur bukti pembayaran (3 gerbang: global/outlet/user)
	IsPaymentProofFeatureEnabled(ctx context.Context, userGoldID int, outcode string) (bool, error)
	GetPaymentProofGlobal(ctx context.Context) (bool, error)
	SetPaymentProofGlobal(ctx context.Context, enabled bool, updatedBy string) error
	GetProofAccessOutlets(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessOutlet, error)
	SaveProofAccessOutlets(ctx context.Context, items []goldSaleEntity.ProofAccessOutletItem) error
	GetProofAccessUsers(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessUser, error)
	SaveProofAccessUsers(ctx context.Context, items []goldSaleEntity.ProofAccessUserItem) error
}

type (
	// Handler ...
	Handler struct {
		goldgymSvcSale IgoldgymSvcSale
		voucherSvc     IVoucherRedeemer
		kafkaProd      *kafka.Producer
		tracer         opentracing.Tracer
		logger         jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(isim IgoldgymSvcSale, voucherSvc IVoucherRedeemer, kp *kafka.Producer, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvcSale: isim,
		voucherSvc:     voucherSvc,
		kafkaProd:      kp,
		tracer:         tracer,
		logger:         logger,
	}
}
