package goldgym

import (
	"context"
	"errors"
	"gold-gym-be/internal/entity"
	jaegerLog "gold-gym-be/pkg/log"
	"time"

	goldEntity "gold-gym-be/internal/entity/goldgym"

	"github.com/opentracing/opentracing-go"
	// "go.opentelemetry.io/otel/trace"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type RepoData interface {
	GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error)
	InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (string, error)
	GetGoldUserByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error)
	GetGoldUserByEmailLogin(ctx context.Context, email string, password string) (goldEntity.LoginUser, error)
	GetGoldToken(ctx context.Context) (goldEntity.LoginToken, error)
	UpdateGoldToken(ctx context.Context, user goldEntity.LoginTokenDataPeserta) error
	GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error)
	InsertSubscription(ctx context.Context, user goldEntity.SubscriptionAll) error
	InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) error
	// DeleteSubscriptionHeader(ctx context.Context, user goldEntity.DeleteSubsHeader) error
	DeleteSubscriptionDetail(ctx context.Context, user goldEntity.DeleteSubs) error
	UpdateSubscriptionDetail(ctx context.Context, user goldEntity.UpdateSubs) error
	UpdateDataPeserta(ctx context.Context, user goldEntity.UpdatePassword) error
	UpdateNama(ctx context.Context, user goldEntity.UpdateNama) error
	UpdateGoldToko(ctx context.Context, goldid int, toko string) error
	UpdateGoldBuyerYN(ctx context.Context, goldid int) error
	UpdateKartu(ctx context.Context, user goldEntity.UpdateKartu) error
	Logout(ctx context.Context, user goldEntity.Logout) error
	GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error)
	GetValidationGoldOTP(ctx context.Context, otp string) (goldEntity.GetValidationGoldOTP, error)
	// UpdateValidationOTP(ctx context.Context, user goldEntity.UpdateValidationOTP) error
	UpdateValidationOTP(ctx context.Context, email string) error
	UpdateOtpIsNull(ctx context.Context, email string) error
	UpdateOTP(ctx context.Context, otp string, email string) error
	GetOneSubscription(ctx context.Context, menuid int) (goldEntity.Subscription, error)
	BulkInsertSubscriptionDetail(ctx context.Context, user []goldEntity.SubscriptionDetail) error
	UpdateOTPSubscription(ctx context.Context, otp string, id int) error
	GetSubscriptionHeader(ctx context.Context, id int) (goldEntity.SubscriptionHeader, error)
	UpdateValidasiPaymentHeader(ctx context.Context, updatePayment goldEntity.UpdatePayment) error
	UpdateValidasiPaymentDetail(ctx context.Context, updatePayment goldEntity.UpdatePayment) error
	GetSubscriptionHeaderTotalHarga(ctx context.Context, id int) (goldEntity.SubscriptionHeaderPayment, error)
	GetPasswordByUser(ctx context.Context, _user string) (string, error)
	UpdateLastLogin(ctx context.Context, _user goldEntity.GetGoldUserss) error

	//testings
	UploadTestingImages(ctx context.Context, testing goldEntity.Testings) (string, error)
	GetTestingImages(ctx context.Context, id int) ([]byte, error)

	GetGoldUserByID(ctx context.Context, id string) (goldEntity.GetGoldUserss, error)

	InsertBuyer(ctx context.Context, buyer goldEntity.BuyerRow) (goldEntity.BuyerRow, error)
	CountGoldUsers(ctx context.Context) (int64, error)
}

type RedisData interface {
	AddToRedis(ctx context.Context, data interface{}, key string, time time.Duration) (err error)
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	goldgym RepoData
	redis   RedisData
	tracer  opentracing.Tracer
	// tracer trace.Tracer
	logger jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(goldgymData RepoData, redisData RedisData, tracer opentracing.Tracer, logger jaegerLog.Factory) *Service {
	// Assign variable dari parameter ke object
	return &Service{
		goldgym: goldgymData,
		redis:   redisData,
		tracer:  tracer,
		logger:  logger,
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
