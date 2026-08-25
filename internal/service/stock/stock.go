package goldgym

import (
	"context"
	"errors"
	"gold-gym-be/internal/entity"
	firebaseEntity "gold-gym-be/internal/entity/firebase"
	goldStockEntity "gold-gym-be/internal/entity/stock"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	// "go.opentelemetry.io/otel/trace"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
	GetOneStockProduct(ctx context.Context, stockcode string, stockname string, stockid string) (goldStockEntity.GetOneStock, error)
	GetLastStock(ctx context.Context) (goldStockEntity.GetOneStock, error)
	InsertStockSales(ctx context.Context, stock goldStockEntity.InsertStock) (string, error)
	AddStockByDate(ctx context.Context, stock goldStockEntity.InsertStock) (string, error)
	GetStockByID(ctx context.Context, stockcode string) ([]goldStockEntity.GetOneStock, error)
	UpdateStockQty(ctx context.Context, stockQty int, stockcode string) (string, error)
	GetAllStockHeader(ctx context.Context) ([]goldStockEntity.GetOneStock, error)
	GetAllStockHeaderToRedis(ctx context.Context) ([]goldStockEntity.GetOneStock, error) // GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error)
	GetFromFirebase(ctx context.Context, userID string) (*firebaseEntity.User, error)
	CreateUser(ctx context.Context, user firebaseEntity.User) (string, error)
	GetStockByName(ctx context.Context, goldid int, stockname string) ([]goldStockEntity.GetOneStock, error)
	// -----------------------------------------------------------------------------------------------------------
	GetLastStockCode(ctx context.Context, goldid int, code string) (*string, error)
	WithTransactionStock(ctx context.Context, fn func(tx *gorm.DB) error) error
	InsertStockNew(ctx context.Context, tx *gorm.DB, id int, stock []goldStockEntity.InsertStock) error
	GetLastHistoryStockCode(ctx context.Context, goldid int, code string) (*string, error)
	InsertStockHistoryNew(ctx context.Context, tx *gorm.DB, id int, stockHistory []goldStockEntity.StockHistory) error
	GetStockExist(ctx context.Context, goldid int, outcode string, itemid []int) (map[int]goldStockEntity.GetOneStock, error)
	GetStock(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldStockEntity.GetOneStock, error)
	GetTotalStock(ctx context.Context, goldid int, name string, outcode string) (int64, error)
	// InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (string, error)
	// GetGoldUserByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error)
	// GetGoldUserByEmailLogin(ctx context.Context, email string, password string) (goldEntity.LoginUser, error)
	// GetGoldToken(ctx context.Context) (goldEntity.LoginToken, error)
	// UpdateGoldToken(ctx context.Context, user goldEntity.LoginTokenDataPeserta) error
	// GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error)
	// InsertSubscription(ctx context.Context, user goldEntity.SubscriptionAll) error
	// InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) error
	// // DeleteSubscriptionHeader(ctx context.Context, user goldEntity.DeleteSubsHeader) error
	// DeleteSubscriptionDetail(ctx context.Context, user goldEntity.DeleteSubs) error
	// UpdateSubscriptionDetail(ctx context.Context, user goldEntity.UpdateSubs) error
	// UpdateDataPeserta(ctx context.Context, user goldEntity.UpdatePassword) error
	// UpdateNama(ctx context.Context, user goldEntity.UpdateNama) error
	// UpdateKartu(ctx context.Context, user goldEntity.UpdateKartu) error
	// Logout(ctx context.Context, user goldEntity.Logout) error
	// GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error)
	// GetValidationGoldOTP(ctx context.Context, otp string) (goldEntity.GetValidationGoldOTP, error)
	// // UpdateValidationOTP(ctx context.Context, user goldEntity.UpdateValidationOTP) error
	// UpdateValidationOTP(ctx context.Context, email string) error
	// UpdateOtpIsNull(ctx context.Context, email string) error
	// UpdateOTP(ctx context.Context, otp string, email string) error
	// GetOneSubscription(ctx context.Context, menuid int) (goldEntity.Subscription, error)
	// BulkInsertSubscriptionDetail(ctx context.Context, user []goldEntity.SubscriptionDetail) error
	// UpdateOTPSubscription(ctx context.Context, otp string, id int) error
	// GetSubscriptionHeader(ctx context.Context, id int) (goldEntity.SubscriptionHeader, error)
	// UpdateValidasiPaymentHeader(ctx context.Context, updatePayment goldEntity.UpdatePayment) error
	// UpdateValidasiPaymentDetail(ctx context.Context, updatePayment goldEntity.UpdatePayment) error
	// GetSubscriptionHeaderTotalHarga(ctx context.Context, id int) (goldEntity.SubscriptionHeaderPayment, error)
	// GetPasswordByUser(ctx context.Context, _user string) (string, error)
	// UpdateLastLogin(ctx context.Context, _user goldEntity.GetGoldUserss) error

}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	goldgymstock Data
	tracer       opentracing.Tracer
	// tracer trace.Tracer
	logger jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(goldgymStockData Data, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	// Assign variable dari parameter ke object
	return Service{
		goldgymstock: goldgymStockData,
		tracer:       tracer,
		logger:       logger,
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
