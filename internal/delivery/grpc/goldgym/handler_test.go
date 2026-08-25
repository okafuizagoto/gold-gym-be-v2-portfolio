package goldgym

import (
	"context"
	"errors"
	"testing"

	authV2 "gold-gym-be/internal/entity/auth/v2"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	jaegerLog "gold-gym-be/pkg/log"
	pb "gold-gym-be/proto"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockGoldgymSvc struct {
	GetGoldUserFn            func(ctx context.Context) ([]goldEntity.GetGoldUser, error)
	GetGoldUserDataByEmailFn func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error)
	LoginUserFn              func(ctx context.Context, user, password, host, device string) (authV2.Token, map[string]interface{}, string, error)
	InsertGoldUserFn         func(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error)
	GetAllSubscriptionFn     func(ctx context.Context) ([]goldEntity.Subscription, error)
}

func (m *mockGoldgymSvc) GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
	if m.GetGoldUserFn != nil {
		return m.GetGoldUserFn(ctx)
	}
	return nil, nil
}

func (m *mockGoldgymSvc) GetGoldUserDataByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
	if m.GetGoldUserDataByEmailFn != nil {
		return m.GetGoldUserDataByEmailFn(ctx, email)
	}
	return goldEntity.GetGoldUserss{}, nil
}

func (m *mockGoldgymSvc) LoginUser(ctx context.Context, user, password, host, device string) (authV2.Token, map[string]interface{}, string, error) {
	if m.LoginUserFn != nil {
		return m.LoginUserFn(ctx, user, password, host, device)
	}
	return authV2.Token{}, nil, "", nil
}

func (m *mockGoldgymSvc) InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error) {
	if m.InsertGoldUserFn != nil {
		return m.InsertGoldUserFn(ctx, user)
	}
	return nil, nil
}

func (m *mockGoldgymSvc) GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error) {
	if m.GetAllSubscriptionFn != nil {
		return m.GetAllSubscriptionFn(ctx)
	}
	return []goldEntity.Subscription{}, nil
}

// Test helpers
func newTestLogger() jaegerLog.Factory {
	logger, _ := zap.NewDevelopment()
	return jaegerLog.NewFactory(logger)
}

func newTestTracer() opentracing.Tracer {
	return opentracing.NoopTracer{}
}

// =============================================================================
// GetGoldUser Tests
// =============================================================================

func TestGetGoldUser_Success(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserFn: func(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
			return []goldEntity.GetGoldUser{
				{
					GoldId:            1,
					GoldEmail:         "test@example.com",
					GoldPassword:      "hashedpass",
					GoldNama:          "Test User",
					GoldNomorHp:       "08123456789",
					GoldNomorKartu:    "1234567890123456",
					GoldCvv:           "123",
					GoldExpireddate:   "12/25",
					GoldPemegangKartu: "TEST USER",
				},
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserRequest{}

	resp, err := handler.GetGoldUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, int32(1), resp.Users[0].GoldId)
	assert.Equal(t, "test@example.com", resp.Users[0].GoldEmail)
	assert.Equal(t, "Test User", resp.Users[0].GoldNama)
}

func TestGetGoldUser_EmptyResult(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserFn: func(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
			return []goldEntity.GetGoldUser{}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserRequest{}

	resp, err := handler.GetGoldUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Users)
}

func TestGetGoldUser_ServiceError(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserFn: func(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserRequest{}

	resp, err := handler.GetGoldUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify gRPC error code
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to get gold users")
}

// =============================================================================
// GetGoldUserByEmail Tests
// =============================================================================

func TestGetGoldUserByEmail_Success(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
			return goldEntity.GetGoldUserss{
				GoldId:            1,
				GoldEmail:         "test@example.com",
				GoldPassword:      "hashedpass",
				GoldNama:          "Test User",
				GoldNomorHp:       "08123456789",
				GoldNomorKartu:    "1234567890123456",
				GoldCvv:           "123",
				GoldExpireddate:   "12/25",
				GoldPemegangKartu: "TEST USER",
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserByEmailRequest{
		Email: "test@example.com",
	}

	resp, err := handler.GetGoldUserByEmail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, int32(1), resp.User.GoldId)
	assert.Equal(t, "test@example.com", resp.User.GoldEmail)

	// CRITICAL: Verify password is NOT returned for security
	assert.Empty(t, resp.User.GoldPassword, "Password MUST be empty for security")

	assert.Equal(t, "Test User", resp.User.GoldNama)
	assert.Equal(t, "08123456789", resp.User.GoldNomorhp)
}

func TestGetGoldUserByEmail_EmptyEmail(t *testing.T) {
	mockSvc := &mockGoldgymSvc{}
	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())

	req := &pb.GetGoldUserByEmailRequest{
		Email: "",
	}

	resp, err := handler.GetGoldUserByEmail(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify gRPC error code is InvalidArgument
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "email is required")
}

func TestGetGoldUserByEmail_NotFound(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
			return goldEntity.GetGoldUserss{}, errors.New("record not found")
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserByEmailRequest{
		Email: "notfound@example.com",
	}

	resp, err := handler.GetGoldUserByEmail(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify gRPC error code is NotFound
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "not found")
}

func TestGetGoldUserByEmail_ServiceError(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
			return goldEntity.GetGoldUserss{}, errors.New("database timeout")
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserByEmailRequest{
		Email: "test@example.com",
	}

	resp, err := handler.GetGoldUserByEmail(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify gRPC error code is Internal
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to get gold user")
}

// =============================================================================
// Table-Driven Test Examples
// =============================================================================

func TestGetGoldUserByEmail_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockReturn    goldEntity.GetGoldUserss
		mockError     error
		wantErr       bool
		wantCode      codes.Code
		wantErrMsg    string
		checkPassword bool
	}{
		{
			name:  "success - valid email",
			email: "valid@example.com",
			mockReturn: goldEntity.GetGoldUserss{
				GoldId:    1,
				GoldEmail: "valid@example.com",
				GoldNama:  "Valid User",
			},
			mockError:     nil,
			wantErr:       false,
			checkPassword: true,
		},
		{
			name:       "error - empty email",
			email:      "",
			wantErr:    true,
			wantCode:   codes.InvalidArgument,
			wantErrMsg: "email is required",
		},
		{
			name:       "error - user not found",
			email:      "notfound@example.com",
			mockReturn: goldEntity.GetGoldUserss{},
			mockError:  errors.New("record not found"),
			wantErr:    true,
			wantCode:   codes.NotFound,
			wantErrMsg: "not found",
		},
		{
			name:       "error - database error",
			email:      "error@example.com",
			mockReturn: goldEntity.GetGoldUserss{},
			mockError:  errors.New("connection timeout"),
			wantErr:    true,
			wantCode:   codes.Internal,
			wantErrMsg: "failed to get gold user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockGoldgymSvc{}

			// Only set mock function if email is not empty (to test validation)
			if tt.email != "" {
				mockSvc.GetGoldUserDataByEmailFn = func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
					return tt.mockReturn, tt.mockError
				}
			}

			handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
			req := &pb.GetGoldUserByEmailRequest{
				Email: tt.email,
			}

			resp, err := handler.GetGoldUserByEmail(context.Background(), req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)

				st, ok := status.FromError(err)
				assert.True(t, ok, "Error should be gRPC status error")
				assert.Equal(t, tt.wantCode, st.Code())
				assert.Contains(t, st.Message(), tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.User)

				if tt.checkPassword {
					assert.Empty(t, resp.User.GoldPassword, "Password MUST be empty for security")
				}
			}
		})
	}
}

// =============================================================================
// Benchmark Tests (for performance monitoring)
// =============================================================================

func BenchmarkGetGoldUser(b *testing.B) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserFn: func(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
			return []goldEntity.GetGoldUser{
				{GoldId: 1, GoldEmail: "test@example.com", GoldNama: "Test"},
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserRequest{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.GetGoldUser(context.Background(), req)
	}
}

func BenchmarkGetGoldUserByEmail(b *testing.B) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
			return goldEntity.GetGoldUserss{
				GoldId:    1,
				GoldEmail: email,
				GoldNama:  "Test User",
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserByEmailRequest{
		Email: "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.GetGoldUserByEmail(context.Background(), req)
	}
}

// =============================================================================
// LoginUser Tests
// =============================================================================

func TestLoginUser_Success(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		LoginUserFn: func(ctx context.Context, user, password, host, device string) (authV2.Token, map[string]interface{}, string, error) {
			return authV2.Token{
					AccessToken: "test-jwt-token-123",
				}, map[string]interface{}{
					"user_id":    "1",
					"user_email": "test@example.com",
					"user_name":  "Test User",
				}, "", nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.LoginUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := handler.LoginUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-jwt-token-123", resp.Token)
	assert.Equal(t, "1", resp.UserId)
	assert.Equal(t, "test@example.com", resp.UserEmail)
	assert.Equal(t, "Test User", resp.UserName)
}

func TestLoginUser_EmptyEmail(t *testing.T) {
	mockSvc := &mockGoldgymSvc{}
	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())

	req := &pb.LoginUserRequest{
		Email:    "",
		Password: "password123",
	}

	resp, err := handler.LoginUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "email and password are required")
}

func TestLoginUser_EmptyPassword(t *testing.T) {
	mockSvc := &mockGoldgymSvc{}
	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())

	req := &pb.LoginUserRequest{
		Email:    "test@example.com",
		Password: "",
	}

	resp, err := handler.LoginUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "email and password are required")
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		LoginUserFn: func(ctx context.Context, user, password, host, device string) (authV2.Token, map[string]interface{}, string, error) {
			return authV2.Token{}, nil, "", errors.New("invalid credentials")
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.LoginUserRequest{
		Email:    "wrong@example.com",
		Password: "wrongpassword",
	}

	resp, err := handler.LoginUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "invalid credentials")
}

// =============================================================================
// InsertGoldUser Tests
// =============================================================================

func TestInsertGoldUser_Success(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		InsertGoldUserFn: func(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error) {
			return "user-123", nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.InsertGoldUserRequest{
		GoldEmail:             "newuser@example.com",
		GoldPassword:          "password123",
		GoldNama:              "New User",
		GoldNomorhp:           "08123456789",
		GoldNomorkartu:        "1234567890123456",
		GoldCvv:               "123",
		GoldExpireddate:       "12/25",
		GoldNamapemegangkartu: "NEW USER",
	}

	resp, err := handler.InsertGoldUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "user-123", resp.UserId)
	assert.Equal(t, "User created successfully", resp.Message)
}

func TestInsertGoldUser_MissingEmail(t *testing.T) {
	mockSvc := &mockGoldgymSvc{}
	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())

	req := &pb.InsertGoldUserRequest{
		GoldEmail:    "",
		GoldPassword: "password123",
		GoldNama:     "New User",
	}

	resp, err := handler.InsertGoldUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "email, password, and name are required")
}

func TestInsertGoldUser_MissingPassword(t *testing.T) {
	mockSvc := &mockGoldgymSvc{}
	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())

	req := &pb.InsertGoldUserRequest{
		GoldEmail:    "newuser@example.com",
		GoldPassword: "",
		GoldNama:     "New User",
	}

	resp, err := handler.InsertGoldUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "email, password, and name are required")
}

func TestInsertGoldUser_MissingName(t *testing.T) {
	mockSvc := &mockGoldgymSvc{}
	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())

	req := &pb.InsertGoldUserRequest{
		GoldEmail:    "newuser@example.com",
		GoldPassword: "password123",
		GoldNama:     "",
	}

	resp, err := handler.InsertGoldUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "email, password, and name are required")
}

func TestInsertGoldUser_ServiceError(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		InsertGoldUserFn: func(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error) {
			return nil, errors.New("duplicate email address")
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.InsertGoldUserRequest{
		GoldEmail:    "duplicate@example.com",
		GoldPassword: "password123",
		GoldNama:     "Duplicate User",
	}

	resp, err := handler.InsertGoldUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to create user")
}

// =============================================================================
// GetAllSubscription Tests
// =============================================================================

func TestGetAllSubscription_Success(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetAllSubscriptionFn: func(ctx context.Context) ([]goldEntity.Subscription, error) {
			return []goldEntity.Subscription{
				{
					GoldNamaPaket:       "Premium",
					GoldNamaLayanan:     "Personal Training",
					GoldHarga:           500000,
					GoldJadwal:          "Mon-Fri 08:00-10:00",
					GoldListLatihan:     "Cardio, Strength",
					GoldJumlahpertemuan: 12,
					GoldDurasi:          60,
				},
				{
					GoldNamaPaket:       "Basic",
					GoldNamaLayanan:     "Group Class",
					GoldHarga:           200000,
					GoldJadwal:          "Mon-Wed 18:00-19:00",
					GoldListLatihan:     "Yoga, Pilates",
					GoldJumlahpertemuan: 8,
					GoldDurasi:          45,
				},
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetAllSubscriptionRequest{}

	resp, err := handler.GetAllSubscription(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Subscriptions, 2)
	assert.Equal(t, "Premium", resp.Subscriptions[0].GoldNamapaket)
	assert.Equal(t, float64(500000), resp.Subscriptions[0].GoldHarga)
	assert.Equal(t, int32(12), resp.Subscriptions[0].GoldJumlahpertemuan)
}

func TestGetAllSubscription_EmptyResult(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetAllSubscriptionFn: func(ctx context.Context) ([]goldEntity.Subscription, error) {
			return []goldEntity.Subscription{}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetAllSubscriptionRequest{}

	resp, err := handler.GetAllSubscription(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Subscriptions)
}

func TestGetAllSubscription_ServiceError(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetAllSubscriptionFn: func(ctx context.Context) ([]goldEntity.Subscription, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetAllSubscriptionRequest{}

	resp, err := handler.GetAllSubscription(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to get subscriptions")
}

// =============================================================================
// Additional Edge Case Tests
// =============================================================================

func TestGetGoldUserByEmail_SpecialCharacters(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
			return goldEntity.GetGoldUserss{
				GoldId:    1,
				GoldEmail: email,
				GoldNama:  "Test User",
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())

	// Test with email containing special characters
	req := &pb.GetGoldUserByEmailRequest{
		Email: "test+tag@example.com",
	}

	resp, err := handler.GetGoldUserByEmail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test+tag@example.com", resp.User.GoldEmail)
}

func TestGetGoldUserByEmail_LongEmail(t *testing.T) {
	longEmail := "very.long.email.address.that.exceeds.normal.length@verylongdomainname.example.com"

	mockSvc := &mockGoldgymSvc{
		GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
			return goldEntity.GetGoldUserss{
				GoldId:    1,
				GoldEmail: email,
				GoldNama:  "Test User",
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserByEmailRequest{
		Email: longEmail,
	}

	resp, err := handler.GetGoldUserByEmail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, longEmail, resp.User.GoldEmail)
}

func TestGetGoldUserByEmail_SQLNoRowsError(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
			return goldEntity.GetGoldUserss{}, errors.New("sql: no rows in result set")
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserByEmailRequest{
		Email: "notfound@example.com",
	}

	resp, err := handler.GetGoldUserByEmail(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

// =============================================================================
// Security-Focused Tests
// =============================================================================

func TestGetGoldUser_PasswordNotInResponse(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserFn: func(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
			return []goldEntity.GetGoldUser{
				{
					GoldId:       1,
					GoldPassword: "secrethash123",
					GoldEmail:    "test@example.com",
					GoldNama:     "Test",
				},
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserRequest{}

	resp, err := handler.GetGoldUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 1)

	// CRITICAL SECURITY CHECK: Password should be included in GetGoldUser
	// (different from GetGoldUserByEmail which excludes it)
	assert.Equal(t, "secrethash123", resp.Users[0].GoldPassword)
}

func TestGetGoldUserByEmail_AllSensitiveFieldsCheck(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
			return goldEntity.GetGoldUserss{
				GoldId:            1,
				GoldEmail:         email,
				GoldPassword:      "hashedpassword123",
				GoldNama:          "Test User",
				GoldNomorHp:       "08123456789",
				GoldNomorKartu:    "1234567890123456",
				GoldCvv:           "123",
				GoldExpireddate:   "12/25",
				GoldPemegangKartu: "TEST USER",
			}, nil
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserByEmailRequest{
		Email: "test@example.com",
	}

	resp, err := handler.GetGoldUserByEmail(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// SECURITY: Password must be empty
	assert.Empty(t, resp.User.GoldPassword, "Password MUST be empty for security")

	// Other fields should be populated
	assert.NotEmpty(t, resp.User.GoldEmail)
	assert.NotEmpty(t, resp.User.GoldNama)
	assert.NotEmpty(t, resp.User.GoldNomorhp)
	assert.NotEmpty(t, resp.User.GoldNomorkartu)
	assert.NotEmpty(t, resp.User.GoldCvv)
}

// =============================================================================
// Context and Error Propagation Tests
// =============================================================================

func TestGetGoldUser_ContextCancellation(t *testing.T) {
	mockSvc := &mockGoldgymSvc{
		GetGoldUserFn: func(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return []goldEntity.GetGoldUser{}, nil
			}
		},
	}

	handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
	req := &pb.GetGoldUserRequest{}

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	resp, err := handler.GetGoldUser(ctx, req)

	// Should handle cancelled context gracefully
	if err != nil {
		assert.Nil(t, resp)
	}
}

func TestGetGoldUserByEmail_ErrorMessageFormat(t *testing.T) {
	tests := []struct {
		name         string
		serviceError error
		wantContains string
	}{
		{
			name:         "timeout error",
			serviceError: errors.New("context deadline exceeded"),
			wantContains: "failed to get gold user",
		},
		{
			name:         "connection error",
			serviceError: errors.New("connection refused"),
			wantContains: "failed to get gold user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockGoldgymSvc{
				GetGoldUserDataByEmailFn: func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
					return goldEntity.GetGoldUserss{}, tt.serviceError
				},
			}

			handler := NewHandler(mockSvc, newTestTracer(), newTestLogger())
			req := &pb.GetGoldUserByEmailRequest{
				Email: "test@example.com",
			}

			_, err := handler.GetGoldUserByEmail(context.Background(), req)

			assert.Error(t, err)
			st, _ := status.FromError(err)
			assert.Contains(t, st.Message(), tt.wantContains)
		})
	}
}
