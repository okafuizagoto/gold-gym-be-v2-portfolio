package goldgym

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gold-gym-be/internal/entity/auth/v2"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	jaegerLog "gold-gym-be/pkg/log"

	"go.uber.org/zap"

	"github.com/opentracing/opentracing-go"
)

type mockService struct {
	users []goldEntity.GetGoldUser
	err   error
}

// Implement all methods from IgoldgymSvc interface
func (m *mockService) GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
	return m.users, m.err
}

func (m *mockService) InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error) {
	return "success", m.err
}

func (m *mockService) GetGoldUserByEmail(ctx context.Context, email string) (string, error) {
	return "TERDAFTAR", m.err
}

func (m *mockService) LoginUser(ctx context.Context, _user, _password string, _host, device string) (auth.Token, map[string]interface{}, string, error) {
	return auth.Token{}, nil, "", m.err
}

func (m *mockService) GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error) {
	return []goldEntity.Subscription{}, m.err
}

func (m *mockService) InsertSubscriptionUser(ctx context.Context, subs goldEntity.InsertSubsAll) (string, error) {
	return "success", m.err
}

func (m *mockService) DeleteSubscriptionHeader(ctx context.Context, subs goldEntity.DeleteSubs) (string, error) {
	return "success", m.err
}

func (m *mockService) UpdateSubscriptionDetail(ctx context.Context, subs goldEntity.UpdateSubs) (string, error) {
	return "success", m.err
}

func (m *mockService) UpdateDataPeserta(ctx context.Context, subs goldEntity.UpdatePassword) (string, error) {
	return "success", m.err
}

func (m *mockService) UpdateNama(ctx context.Context, subs goldEntity.UpdateNama) (string, error) {
	return "success", m.err
}

func (m *mockService) UpdateKartu(ctx context.Context, subs goldEntity.UpdateKartu) (string, error) {
	return "success", m.err
}

func (m *mockService) Logout(ctx context.Context, subs goldEntity.Logout) (string, error) {
	return "success", m.err
}

func (m *mockService) GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error) {
	return []goldEntity.GetSubsWithUser{}, m.err
}

func (m *mockService) UpdateValidationOTP(ctx context.Context, otp string, email string) (string, error) {
	return "success", m.err
}

func (m *mockService) UpdateOTP(ctx context.Context, email string) (string, error) {
	return "success", m.err
}

func (m *mockService) InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) (string, error, response.Response) {
	return "success", m.err, response.Response{}
}

func (m *mockService) UpdateOTPSubscription(ctx context.Context, email string) (string, error) {
	return "success", m.err
}

func (m *mockService) UpdatePayment(ctx context.Context, otp string, email string) (string, error, response.Response) {
	return "success", m.err, response.Response{}
}

func (m *mockService) GetSubscriptionHeaderTotalHarga(ctx context.Context, email string) (goldEntity.SubscriptionHeaderPayment, error) {
	return goldEntity.SubscriptionHeaderPayment{}, m.err
}

func (m *mockService) UploadTestingImages(ctx context.Context, testing goldEntity.Testings) (string, error) {
	return "success", m.err
}

func (m *mockService) GetTestingImage(ctx context.Context, id int) ([]byte, error) {
	return []byte{}, m.err
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/gold-gym/v2/userdata", h.GetGoldGymGin)
	return r
}

func newTestLogger() jaegerLog.Factory {
	logger, _ := zap.NewDevelopment()
	return jaegerLog.NewFactory(logger)
}

func newTestTracer() opentracing.Tracer {
	return opentracing.NoopTracer{}
}

func TestGetGoldGym_Success(t *testing.T) {
	svc := &mockService{
		users: []goldEntity.GetGoldUser{
			{GoldId: 1, GoldNama: "Budi"},
		},
	}

	h := New(svc, nil, newTestTracer(), newTestLogger())
	r := setupRouter(h)

	req, _ := http.NewRequest("GET", "/gold-gym/v2/userdata?type=getgoldgym", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Budi")
}

func TestGetGoldGym_Error(t *testing.T) {
	svc := &mockService{
		err: errors.New("db down"),
	}

	h := New(svc, nil, newTestTracer(), newTestLogger())
	r := setupRouter(h)

	req, _ := http.NewRequest("GET", "/gold-gym/v2/userdata?type=getgoldgym", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func (m *mockService) RegisterBuyer(ctx context.Context, req goldEntity.RegisterBuyerRequest) (goldEntity.BuyerRow, error) {
	return goldEntity.BuyerRow{GoldNama: req.GoldNama, GoldEmail: req.GoldEmail, GoldRole: "BUYER"}, m.err
}
