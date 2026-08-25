package goldgym

import (
	"context"
	"errors"
	"gold-gym-be/internal/entity"
	goldSaleEntity "gold-gym-be/internal/entity/sales"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	// "go.opentelemetry.io/otel/trace"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	WithTransactionThSale(ctx context.Context, fn func(tx *gorm.DB) error) error
	GetLastThSaleCode(ctx context.Context, goldid int, code string) (*string, *string, error)
	InsertThSale(ctx context.Context, tx *gorm.DB, items []goldSaleEntity.ThSale) error
	InsertTdSale(ctx context.Context, tx *gorm.DB, items []goldSaleEntity.TdSale) error
	GetLastTdSaleCode(ctx context.Context, saleid string) (*string, error)
	GetThSale(ctx context.Context, goldid int, custid int, name string, outcode string, page, length int) ([]goldSaleEntity.ThSale, error)
	GetTotalThSale(ctx context.Context, goldid int, custid int, name string, outcode string) (int64, error)
	GetThSaleByID(ctx context.Context, goldid int, custid int, saleid string) (*goldSaleEntity.ThSale, error)
	UpdateSalePaymentYN(ctx context.Context, goldid int, saleid string) (int64, error)
	GetOutletGoldIDByCode(ctx context.Context, outcode string) (int, error)
	GetTdSaleBySaleID(ctx context.Context, saleid string) ([]goldSaleEntity.TdSale, error)
	GetSaleOutletInfo(ctx context.Context, goldid int, outcode string) (*goldSaleEntity.SaleOutletInfo, error)
	GetSaleCustomerToko(ctx context.Context, custid int) (string, error)
	IsPosCustomerOptional(ctx context.Context, goldid int, outcode string) (bool, error)
	GetOutletTypeByCode(ctx context.Context, goldid int, outcode string) (string, error)
	GetPosOutletsForAdmin(ctx context.Context, search string) ([]goldSaleEntity.PosOutlet, error)
	SetPosCustomerOptional(ctx context.Context, goldid int, outcode string, optional bool, addedBy string) error
	MarkBookingsPaid(ctx context.Context, bookingIDs []string, saleID string) (int64, error)
	GetTotalProofBytes(ctx context.Context) (int64, error)
	InsertPaymentProof(ctx context.Context, proof goldSaleEntity.PaymentProof) (goldSaleEntity.PaymentProof, error)
	GetPaymentProofsBySale(ctx context.Context, saleID string) ([]goldSaleEntity.PaymentProof, error)
	GetPaymentProofByID(ctx context.Context, proofID int) (*goldSaleEntity.PaymentProof, error)
	GetSaleReportItems(ctx context.Context, goldid int, outcode, date string) ([]goldSaleEntity.SaleReportItem, error)
	GetSaleDailyTotals(ctx context.Context, goldid int, outcode, start, end string) ([]goldSaleEntity.SaleDailyTotal, error)

	// Visibilitas fitur bukti pembayaran (3 gerbang: global/outlet/user)
	IsPaymentProofGloballyEnabled(ctx context.Context) (bool, error)
	SetPaymentProofGlobal(ctx context.Context, enabled bool, updatedBy string) error
	IsOutletProofEnabled(ctx context.Context, goldid int, outcode string) (bool, error)
	IsUserProofEnabled(ctx context.Context, goldid int) (bool, error)
	GetProofAccessOutlets(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessOutlet, error)
	SetProofOutletEnabled(ctx context.Context, goldid int, outcode string, enabled bool) error
	GetProofAccessUsers(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessUser, error)
	SetProofUserEnabled(ctx context.Context, goldid int, enabled bool) error
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	goldgymsale Data
	tracer      opentracing.Tracer
	// tracer trace.Tracer
	logger jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(goldgymSaleData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	// Assign variable dari parameter ke object
	return Service{
		goldgymsale: goldgymSaleData,
		tracer:      tracer,
		logger:      logger,
	}
}

func (s Service) checkPermission(ctx context.Context, _permissions ...string) error {
	claims := ctx.Value(entity.ContextKey("claims"))
	if claims != nil {
		actions := claims.(entity.ContextValue).Get("permissions").(map[string]interface{})
		for _, action := range actions {
			permissions := action.([]interface{})
			for _, permission := range permissions {
				for _, _permission := range _permissions {
					if permission.(string) == _permission {
						return nil
					}
				}
			}
		}
	}
	return errors.New("401 unauthorized")
}
