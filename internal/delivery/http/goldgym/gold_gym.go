package goldgym

import (
	"context"
	"gold-gym-be/pkg/kafka"
	jaegerLog "gold-gym-be/pkg/log"
	"gold-gym-be/pkg/response"

	"gold-gym-be/internal/entity/auth/v2"
	goldEntity "gold-gym-be/internal/entity/goldgym"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvc interface {
	GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error)
	InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error)
	GetGoldUserByEmail(ctx context.Context, email string) (string, error)
	// LoginUser(ctx context.Context, user goldEntity.LogUser) (interface{}, goldEntity.LoginUser, error)
	LoginUser(ctx context.Context, _user, _password string, _host, device string) (auth.Token, map[string]interface{}, string, error)
	GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error)
	InsertSubscriptionUser(ctx context.Context, subs goldEntity.InsertSubsAll) (string, error)
	DeleteSubscriptionHeader(ctx context.Context, subs goldEntity.DeleteSubs) (string, error)
	UpdateSubscriptionDetail(ctx context.Context, subs goldEntity.UpdateSubs) (string, error)
	UpdateDataPeserta(ctx context.Context, subs goldEntity.UpdatePassword) (string, error)
	UpdateNama(ctx context.Context, subs goldEntity.UpdateNama) (string, error)
	SetToko(ctx context.Context, goldid int, toko string) (string, error)
	RegisterAsBuyer(ctx context.Context, goldid int) (string, error)
	UpdateKartu(ctx context.Context, subs goldEntity.UpdateKartu) (string, error)
	Logout(ctx context.Context, subs goldEntity.Logout) (string, error)
	GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error)
	UpdateValidationOTP(ctx context.Context, otp string, email string) (string, error)
	UpdateOTP(ctx context.Context, email string) (string, error)
	InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) (string, error, response.Response)
	// UpdateOTPSubscription(ctx context.Context, email string) (string, time.Time, error)
	UpdateOTPSubscription(ctx context.Context, email string) (string, error)
	UpdatePayment(ctx context.Context, otp string, email string) (string, error, response.Response)
	GetSubscriptionHeaderTotalHarga(ctx context.Context, email string) (goldEntity.SubscriptionHeaderPayment, error)

	UploadTestingImages(ctx context.Context, testing goldEntity.Testings) (string, error)
	GetTestingImage(ctx context.Context, id int) ([]byte, error)

	RegisterBuyer(ctx context.Context, req goldEntity.RegisterBuyerRequest) (goldEntity.BuyerRow, error)

	GetRegistrationMode(ctx context.Context) (string, error)
	SetRegistrationMode(ctx context.Context, mode string, updatedBy string) error
}

type (
	// Handler ...
	Handler struct {
		goldgymSvc IgoldgymSvc
		kafkaProd  *kafka.Producer
		tracer     opentracing.Tracer
		logger     jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(is IgoldgymSvc, kp *kafka.Producer, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvc: is,
		kafkaProd:  kp,
		tracer:     tracer,
		logger:     logger,
	}
}
