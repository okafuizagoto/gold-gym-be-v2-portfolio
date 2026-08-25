package goldgym

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	goldEntity "gold-gym-be/internal/entity/goldgym"

	"github.com/raja/argon2pw"
	"github.com/stretchr/testify/assert"
)

// otpCreatedNow adalah helper untuk mock GetValidationGoldOTP -- OTP yang
// baru saja dibuat, masih berlaku (belum kedaluwarsa 10 menit).
func otpCreatedNow() *time.Time {
	t := time.Now()
	return &t
}

// testPasswordHash adalah argon2 hash dari "testpass123", disiapkan sekali di TestMain.
var testPasswordHash string

func TestMain(m *testing.M) {
	var err error
	testPasswordHash, err = argon2pw.GenerateSaltedHash("testpass123")
	if err != nil {
		log.Fatal("failed to hash test password: ", err)
	}
	os.Exit(m.Run())
}

// =============================================================================
// Service method yang TIDAK ditest di sini:
//   - InsertGoldUser        : memanggil SMTP (sendOTP) langsung di tengah fungsi
//   - UpdateOTP             : memanggil SMTP (sendOTP)
//   - PaymentValidation     : memanggil SMTP (sendOTP)
//   - UpdateOTPSubscription : memanggil SMTP (sendOTP)
//   - UpdatePayment         : logika expiry berbasis time.Now(), butuh mock waktu
// =============================================================================

// --- GetGoldUser ---

func TestGetGoldUser(t *testing.T) {
	users := []goldEntity.GetGoldUser{
		{GoldId: 1, GoldEmail: "a@test.com", GoldNama: "Budi"},
		{GoldId: 2, GoldEmail: "b@test.com", GoldNama: "Susi"},
	}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    []goldEntity.GetGoldUser
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				GetGoldUserFn: func(_ context.Context) ([]goldEntity.GetGoldUser, error) {
					return users, nil
				},
			},
			want: users,
		},
		{
			name: "empty result",
			repo: &mockRepo{
				GetGoldUserFn: func(_ context.Context) ([]goldEntity.GetGoldUser, error) {
					return []goldEntity.GetGoldUser{}, nil
				},
			},
			want: []goldEntity.GetGoldUser{},
		},
		{
			name: "repo error",
			repo: &mockRepo{
				GetGoldUserFn: func(_ context.Context) ([]goldEntity.GetGoldUser, error) {
					return nil, errors.New("db down")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.GetGoldUser(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- GetGoldUserByEmail ---

func TestGetGoldUserByEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name:  "registered and validated",
			email: "budi@test.com",
			repo: &mockRepo{
				GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
					return goldEntity.GetGoldUserss{GoldId: 1, GoldEmail: "budi@test.com", GoldValidasiYN: "Y"}, nil
				},
			},
			want: "TERDAFTAR",
		},
		{
			name:  "not registered",
			email: "new@test.com",
			repo: &mockRepo{
				GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
					return goldEntity.GetGoldUserss{}, nil
				},
			},
			want: "TIDAK TERDAFTAR",
		},
		{
			name:  "registered but not validated",
			email: "budi@test.com",
			repo: &mockRepo{
				GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
					return goldEntity.GetGoldUserss{GoldId: 1, GoldEmail: "budi@test.com", GoldValidasiYN: "N"}, nil
				},
			},
			want: "BELUM TERVALIDASI",
		},
		{
			name:  "repo error",
			email: "err@test.com",
			repo: &mockRepo{
				GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
					return goldEntity.GetGoldUserss{}, errors.New("timeout")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.GetGoldUserByEmail(context.Background(), tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- GetAllSubscription ---

func TestGetAllSubscription(t *testing.T) {
	subs := []goldEntity.Subscription{
		{GoldNamaPaket: "Basic", GoldHarga: 100},
		{GoldNamaPaket: "Premium", GoldHarga: 250},
	}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    []goldEntity.Subscription
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
					return subs, nil
				},
			},
			want: subs,
		},
		{
			name: "empty",
			repo: &mockRepo{
				GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
					return []goldEntity.Subscription{}, nil
				},
			},
			want: []goldEntity.Subscription{},
		},
		{
			name: "repo error",
			repo: &mockRepo{
				GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.GetAllSubscription(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- LoginUser ---

func TestLoginUser(t *testing.T) {
	mockUser := goldEntity.GetGoldUserss{
		GoldId:       1,
		GoldEmail:    "budi@test.com",
		GoldNama:     "Budi Santoso",
		GoldPassword: testPasswordHash,
	}

	tests := []struct {
		name    string
		user    string
		pass    string
		host    string
		device  string
		repo    *mockRepo
		wantErr bool
	}{
		{
			name: "success",
			user: "budi@test.com",
			pass: "testpass123",
			host: "127.0.0.1",
			repo: &mockRepo{
				GetPasswordByUserFn: func(_ context.Context, _ string) (string, error) {
					return testPasswordHash, nil
				},
				GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
					return mockUser, nil
				},
				UpdateLastLoginFn: func(_ context.Context, u goldEntity.GetGoldUserss) error {
					// verifikasi host dicapture dengan benar
					if u.GoldLastLoginHost != "127.0.0.1" {
						return errors.New("host tidak sesuai")
					}
					return nil
				},
			},
		},
		{
			name: "invalid password",
			user: "budi@test.com",
			pass: "wrongpassword",
			host: "127.0.0.1",
			repo: &mockRepo{
				GetPasswordByUserFn: func(_ context.Context, _ string) (string, error) {
					return testPasswordHash, nil
				},
				GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
					return mockUser, nil
				},
			},
			wantErr: true,
		},
		{
			name: "GetPasswordByUser error",
			user: "budi@test.com",
			pass: "testpass123",
			host: "127.0.0.1",
			repo: &mockRepo{
				GetPasswordByUserFn: func(_ context.Context, _ string) (string, error) {
					return "", errors.New("user not found")
				},
			},
			wantErr: true,
		},
		{
			name: "GetGoldUserByEmail error",
			user: "budi@test.com",
			pass: "testpass123",
			host: "127.0.0.1",
			repo: &mockRepo{
				GetPasswordByUserFn: func(_ context.Context, _ string) (string, error) {
					return testPasswordHash, nil
				},
				GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
					return goldEntity.GetGoldUserss{}, errors.New("db error")
				},
			},
			wantErr: true,
		},
		{
			name: "UpdateLastLogin error",
			user: "budi@test.com",
			pass: "testpass123",
			host: "127.0.0.1",
			repo: &mockRepo{
				GetPasswordByUserFn: func(_ context.Context, _ string) (string, error) {
					return testPasswordHash, nil
				},
				GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
					return mockUser, nil
				},
				UpdateLastLoginFn: func(_ context.Context, _ goldEntity.GetGoldUserss) error {
					return errors.New("update failed")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			token, metadata, _, err := svc.LoginUser(context.Background(), tt.user, tt.pass, tt.host, tt.device)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, token.AccessToken)
			assert.Equal(t, "Bearer", token.TokenType)
			assert.Equal(t, int64(43200), token.ExpiresIn) // 12 jam
			assert.Equal(t, "Budi Santoso", metadata["username"])
		})
	}
}

// --- DeleteSubscriptionHeader ---

func TestDeleteSubscriptionHeader(t *testing.T) {
	input := goldEntity.DeleteSubs{GoldId: 1, GoldMenuId: 2}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				DeleteSubscriptionDetailFn: func(_ context.Context, _ goldEntity.DeleteSubs) error {
					return nil
				},
			},
			want: "Berhasil",
		},
		{
			name: "repo error",
			repo: &mockRepo{
				DeleteSubscriptionDetailFn: func(_ context.Context, _ goldEntity.DeleteSubs) error {
					return errors.New("delete failed")
				},
			},
			want:    "Detail - Gagal",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.DeleteSubscriptionHeader(context.Background(), input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- UpdateSubscriptionDetail ---

func TestUpdateSubscriptionDetail(t *testing.T) {
	input := goldEntity.UpdateSubs{GoldId: 1, GoldMenuId: 2, GoldJumlahpertemuan: 10}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				UpdateSubscriptionDetailFn: func(_ context.Context, _ goldEntity.UpdateSubs) error {
					return nil
				},
			},
			want: "Berhasil",
		},
		{
			name: "repo error",
			repo: &mockRepo{
				UpdateSubscriptionDetailFn: func(_ context.Context, _ goldEntity.UpdateSubs) error {
					return errors.New("update failed")
				},
			},
			want:    "Gagal",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.UpdateSubscriptionDetail(context.Background(), input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- UpdateDataPeserta ---

func TestUpdateDataPeserta(t *testing.T) {
	tests := []struct {
		name    string
		input   goldEntity.UpdatePassword
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name: "success - OTP match",
			input: goldEntity.UpdatePassword{
				GoldEmail:    "budi@test.com",
				GoldOTP:      "123456",
				GoldPassword: "newpass",
			},
			repo: &mockRepo{
				GetValidationGoldOTPFn: func(_ context.Context, _ string) (goldEntity.GetValidationGoldOTP, error) {
					return goldEntity.GetValidationGoldOTP{GoldOTP: "123456", GoldOtpCreatedAt: otpCreatedNow()}, nil
				},
				UpdateDataPesertaFn: func(_ context.Context, _ goldEntity.UpdatePassword) error {
					return nil
				},
				UpdateOtpIsNullFn: func(_ context.Context, _ string) error {
					return nil
				},
			},
			want: "Berhasil",
		},
		{
			name: "OTP not found in DB",
			input: goldEntity.UpdatePassword{
				GoldEmail:    "budi@test.com",
				GoldOTP:      "999999",
				GoldPassword: "newpass",
			},
			repo: &mockRepo{
				GetValidationGoldOTPFn: func(_ context.Context, _ string) (goldEntity.GetValidationGoldOTP, error) {
					return goldEntity.GetValidationGoldOTP{}, nil // OTP kosong di DB
				},
			},
			want: "Please Validation OTP First",
		},
		{
			name: "email empty",
			input: goldEntity.UpdatePassword{
				GoldEmail:    "",
				GoldOTP:      "123456",
				GoldPassword: "newpass",
			},
			repo: &mockRepo{
				GetValidationGoldOTPFn: func(_ context.Context, _ string) (goldEntity.GetValidationGoldOTP, error) {
					return goldEntity.GetValidationGoldOTP{GoldOTP: "123456"}, nil
				},
			},
			want: "Please Field the Email",
		},
		{
			name: "OTP mismatch",
			input: goldEntity.UpdatePassword{
				GoldEmail:    "budi@test.com",
				GoldOTP:      "111111",
				GoldPassword: "newpass",
			},
			repo: &mockRepo{
				GetValidationGoldOTPFn: func(_ context.Context, _ string) (goldEntity.GetValidationGoldOTP, error) {
					return goldEntity.GetValidationGoldOTP{GoldOTP: "999999"}, nil
				},
			},
			want: "OTP is incorrect (validation otp)",
		},
		{
			name: "UpdateOtpIsNull error",
			input: goldEntity.UpdatePassword{
				GoldEmail:    "budi@test.com",
				GoldOTP:      "123456",
				GoldPassword: "newpass",
			},
			repo: &mockRepo{
				GetValidationGoldOTPFn: func(_ context.Context, _ string) (goldEntity.GetValidationGoldOTP, error) {
					return goldEntity.GetValidationGoldOTP{GoldOTP: "123456", GoldOtpCreatedAt: otpCreatedNow()}, nil
				},
				UpdateDataPesertaFn: func(_ context.Context, _ goldEntity.UpdatePassword) error {
					return nil
				},
				UpdateOtpIsNullFn: func(_ context.Context, _ string) error {
					return errors.New("update failed")
				},
			},
			want:    "Gagal",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.UpdateDataPeserta(context.Background(), tt.input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- UpdateNama ---

func TestUpdateNama(t *testing.T) {
	input := goldEntity.UpdateNama{GoldNama: "Budi Baru", GoldEmail: "budi@test.com"}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				UpdateNamaFn: func(_ context.Context, _ goldEntity.UpdateNama) error {
					return nil
				},
			},
			want: "Berhasil",
		},
		{
			name: "repo error",
			repo: &mockRepo{
				UpdateNamaFn: func(_ context.Context, _ goldEntity.UpdateNama) error {
					return errors.New("update failed")
				},
			},
			want:    "Gagal",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.UpdateNama(context.Background(), input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- UpdateKartu ---

func TestUpdateKartu(t *testing.T) {
	input := goldEntity.UpdateKartu{
		GoldNomorKartu: "4111111111111111",
		GoldCvv:        "123",
		GoldEmail:      "budi@test.com",
	}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name: "success - values are hashed before repo call",
			repo: &mockRepo{
				UpdateKartuFn: func(_ context.Context, kartu goldEntity.UpdateKartu) error {
					// pastikan plaintext sudah dihash
					if kartu.GoldNomorKartu == "4111111111111111" {
						return errors.New("GoldNomorKartu seharusnya sudah dihash")
					}
					if kartu.GoldCvv == "123" {
						return errors.New("GoldCvv seharusnya sudah dihash")
					}
					return nil
				},
			},
			want: "Berhasil",
		},
		{
			name: "repo error",
			repo: &mockRepo{
				UpdateKartuFn: func(_ context.Context, _ goldEntity.UpdateKartu) error {
					return errors.New("update failed")
				},
			},
			want:    "Gagal",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.UpdateKartu(context.Background(), input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- Logout ---

func TestLogout(t *testing.T) {
	input := goldEntity.Logout{GoldEmail: "budi@test.com"}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				LogoutFn: func(_ context.Context, _ goldEntity.Logout) error {
					return nil
				},
			},
			want: "Berhasil",
		},
		{
			name: "repo error",
			repo: &mockRepo{
				LogoutFn: func(_ context.Context, _ goldEntity.Logout) error {
					return errors.New("logout failed")
				},
			},
			want:    "Gagal",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.Logout(context.Background(), input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- GetSubsWithUser ---

func TestGetSubsWithUser(t *testing.T) {
	data := []goldEntity.GetSubsWithUser{
		{GoldId: 1, GoldEmail: "a@test.com", GoldNama: "Budi"},
	}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    []goldEntity.GetSubsWithUser
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				GetSubsWithUserFn: func(_ context.Context) ([]goldEntity.GetSubsWithUser, error) {
					return data, nil
				},
			},
			want: data,
		},
		{
			name: "empty",
			repo: &mockRepo{
				GetSubsWithUserFn: func(_ context.Context) ([]goldEntity.GetSubsWithUser, error) {
					return []goldEntity.GetSubsWithUser{}, nil
				},
			},
			want: []goldEntity.GetSubsWithUser{},
		},
		{
			name: "repo error",
			repo: &mockRepo{
				GetSubsWithUserFn: func(_ context.Context) ([]goldEntity.GetSubsWithUser, error) {
					return nil, errors.New("query failed")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.GetSubsWithUser(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- UpdateValidationOTP ---

func TestUpdateValidationOTP(t *testing.T) {
	tests := []struct {
		name    string
		otp     string
		email   string
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name:  "OTP match - success",
			otp:   "123456",
			email: "budi@test.com",
			repo: &mockRepo{
				GetValidationGoldOTPFn: func(_ context.Context, _ string) (goldEntity.GetValidationGoldOTP, error) {
					return goldEntity.GetValidationGoldOTP{GoldOTP: "123456"}, nil
				},
				UpdateValidationOTPFn: func(_ context.Context, _ string) error {
					return nil
				},
				UpdateOtpIsNullFn: func(_ context.Context, _ string) error {
					return nil
				},
			},
			want: "Berhasil",
		},
		{
			name:  "OTP mismatch",
			otp:   "111111",
			email: "budi@test.com",
			repo: &mockRepo{
				GetValidationGoldOTPFn: func(_ context.Context, _ string) (goldEntity.GetValidationGoldOTP, error) {
					return goldEntity.GetValidationGoldOTP{GoldOTP: "999999"}, nil
				},
			},
			want: "OTP is incorrect",
			// errors.Wrap(nil, ...) = nil, jadi err tetap nil di path ini
		},
		{
			name:  "GetValidationGoldOTP error",
			otp:   "111111",
			email: "budi@test.com",
			repo: &mockRepo{
				GetValidationGoldOTPFn: func(_ context.Context, _ string) (goldEntity.GetValidationGoldOTP, error) {
					return goldEntity.GetValidationGoldOTP{}, errors.New("db error")
				},
			},
			want:    "OTP is incorrect",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.UpdateValidationOTP(context.Background(), tt.otp, tt.email)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- InsertSubscriptionUser ---

func TestInsertSubscriptionUser(t *testing.T) {
	products := []goldEntity.Subscription{
		{GoldNamaPaket: "Basic", GoldNamaLayanan: "Gym", GoldHarga: 100, GoldJadwal: "Mon", GoldListLatihan: "Push", GoldJumlahpertemuan: 4, GoldDurasi: 30},
		{GoldNamaPaket: "Premium", GoldNamaLayanan: "Gym+", GoldHarga: 200, GoldJadwal: "Tue", GoldListLatihan: "Pull", GoldJumlahpertemuan: 8, GoldDurasi: 60},
	}
	mockUser := goldEntity.GetGoldUserss{GoldId: 5, GoldEmail: "budi@test.com"}

	t.Run("success - single detail MenuId 1", func(t *testing.T) {
		var capturedHeader goldEntity.SubscriptionAll

		svc := newTestService(&mockRepo{
			GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
				return products, nil
			},
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return mockUser, nil
			},
			InsertSubscriptionDetailFn: func(_ context.Context, detail goldEntity.SubscriptionDetail) error {
				assert.Equal(t, 5, detail.GoldId)
				assert.Equal(t, "Basic", detail.GoldNamaPaket)
				assert.Equal(t, "Belum Berlangganan", detail.GoldStatuslangganan)
				return nil
			},
			InsertSubscriptionFn: func(_ context.Context, header goldEntity.SubscriptionAll) error {
				capturedHeader = header
				return nil
			},
		})

		input := goldEntity.InsertSubsAll{
			HeaderData: goldEntity.SubscriptionAll{GoldEmail: "budi@test.com"},
			DetailData: []goldEntity.SubscriptionDetail{{GoldMenuId: 1}},
		}

		got, err := svc.InsertSubscriptionUser(context.Background(), input)
		assert.NoError(t, err)
		assert.Equal(t, "Berhasil", got)
		assert.Equal(t, 5, capturedHeader.GoldId)
	})

	t.Run("success - single detail MenuId > 1", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
				return products, nil
			},
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return mockUser, nil
			},
			InsertSubscriptionDetailFn: func(_ context.Context, detail goldEntity.SubscriptionDetail) error {
				// MenuId 2 → ambil dari header[MenuId-1] = header[1] = "Premium"
				assert.Equal(t, "Premium", detail.GoldNamaPaket)
				return nil
			},
			InsertSubscriptionFn: func(_ context.Context, _ goldEntity.SubscriptionAll) error {
				return nil
			},
		})

		input := goldEntity.InsertSubsAll{
			HeaderData: goldEntity.SubscriptionAll{GoldEmail: "budi@test.com"},
			DetailData: []goldEntity.SubscriptionDetail{{GoldMenuId: 2}},
		}

		got, err := svc.InsertSubscriptionUser(context.Background(), input)
		assert.NoError(t, err)
		assert.Equal(t, "Berhasil", got)
	})

	t.Run("success - multiple details bulk insert", func(t *testing.T) {
		var capturedHeader goldEntity.SubscriptionAll
		var capturedBulk []goldEntity.SubscriptionDetail

		svc := newTestService(&mockRepo{
			GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
				return products, nil
			},
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return mockUser, nil
			},
			BulkInsertSubscriptionDetailFn: func(_ context.Context, details []goldEntity.SubscriptionDetail) error {
				capturedBulk = details
				return nil
			},
			InsertSubscriptionFn: func(_ context.Context, header goldEntity.SubscriptionAll) error {
				capturedHeader = header
				return nil
			},
		})

		input := goldEntity.InsertSubsAll{
			HeaderData: goldEntity.SubscriptionAll{GoldEmail: "budi@test.com"},
			DetailData: []goldEntity.SubscriptionDetail{
				{GoldMenuId: 1},
				{GoldMenuId: 2},
			},
		}

		got, err := svc.InsertSubscriptionUser(context.Background(), input)
		assert.NoError(t, err)
		assert.Equal(t, "Berhasil", got)
		assert.Len(t, capturedBulk, 2)
		assert.Equal(t, float64(300), capturedHeader.GoldTotalharga) // 100 + 200
	})

	t.Run("email not found", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
				return products, nil
			},
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return goldEntity.GetGoldUserss{}, nil // empty = not found
			},
		})

		input := goldEntity.InsertSubsAll{
			HeaderData: goldEntity.SubscriptionAll{GoldEmail: "notfound@test.com"},
			DetailData: []goldEntity.SubscriptionDetail{{GoldMenuId: 1}},
		}

		got, _ := svc.InsertSubscriptionUser(context.Background(), input)
		assert.Equal(t, "Detail - Gagal - Email Tidak Tersedia", got)
	})

	t.Run("GetGoldUserByEmail error", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
				return products, nil
			},
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return goldEntity.GetGoldUserss{}, errors.New("db error")
			},
		})

		input := goldEntity.InsertSubsAll{
			HeaderData: goldEntity.SubscriptionAll{GoldEmail: "budi@test.com"},
			DetailData: []goldEntity.SubscriptionDetail{{GoldMenuId: 1}},
		}

		_, err := svc.InsertSubscriptionUser(context.Background(), input)
		assert.Error(t, err)
	})

	t.Run("InsertSubscriptionDetail repo error", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
				return products, nil
			},
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return mockUser, nil
			},
			InsertSubscriptionDetailFn: func(_ context.Context, _ goldEntity.SubscriptionDetail) error {
				return errors.New("insert detail failed")
			},
		})

		input := goldEntity.InsertSubsAll{
			HeaderData: goldEntity.SubscriptionAll{GoldEmail: "budi@test.com"},
			DetailData: []goldEntity.SubscriptionDetail{{GoldMenuId: 1}},
		}

		got, err := svc.InsertSubscriptionUser(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, "Detail - Gagal", got)
	})

	t.Run("InsertSubscription header error", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetAllSubscriptionFn: func(_ context.Context) ([]goldEntity.Subscription, error) {
				return products, nil
			},
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return mockUser, nil
			},
			InsertSubscriptionDetailFn: func(_ context.Context, _ goldEntity.SubscriptionDetail) error {
				return nil
			},
			InsertSubscriptionFn: func(_ context.Context, _ goldEntity.SubscriptionAll) error {
				return errors.New("insert header failed")
			},
		})

		input := goldEntity.InsertSubsAll{
			HeaderData: goldEntity.SubscriptionAll{GoldEmail: "budi@test.com"},
			DetailData: []goldEntity.SubscriptionDetail{{GoldMenuId: 1}},
		}

		got, err := svc.InsertSubscriptionUser(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, "Header - Gagal", got)
	})
}

// --- InsertSubscriptionDetail (service method) ---

func TestInsertSubscriptionDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetOneSubscriptionFn: func(_ context.Context, _ int) (goldEntity.Subscription, error) {
				return goldEntity.Subscription{
					GoldNamaPaket: "Basic", GoldNamaLayanan: "Gym", GoldHarga: 100,
					GoldJadwal: "Mon", GoldListLatihan: "Push", GoldJumlahpertemuan: 4, GoldDurasi: 30,
				}, nil
			},
			GetSubscriptionHeaderFn: func(_ context.Context, _ int) (goldEntity.SubscriptionHeader, error) {
				return goldEntity.SubscriptionHeader{GoldID: 1}, nil // non-empty
			},
			InsertSubscriptionDetailFn: func(_ context.Context, detail goldEntity.SubscriptionDetail) error {
				assert.Equal(t, "Basic", detail.GoldNamaPaket)
				assert.Equal(t, "Belum Berlangganan", detail.GoldStatuslangganan)
				return nil
			},
		})

		input := goldEntity.SubscriptionDetail{GoldId: 1, GoldMenuId: 1}
		got, err, resp := svc.InsertSubscriptionDetail(context.Background(), input)
		assert.Equal(t, "Berhasil", got)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.StatusCode)
	})

	t.Run("subscription header empty", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetOneSubscriptionFn: func(_ context.Context, _ int) (goldEntity.Subscription, error) {
				return goldEntity.Subscription{GoldNamaPaket: "Basic"}, nil
			},
			GetSubscriptionHeaderFn: func(_ context.Context, _ int) (goldEntity.SubscriptionHeader, error) {
				return goldEntity.SubscriptionHeader{}, nil // empty → trigger error
			},
		})

		input := goldEntity.SubscriptionDetail{GoldId: 1, GoldMenuId: 1}
		got, _, resp := svc.InsertSubscriptionDetail(context.Background(), input)
		assert.Equal(t, "Subscription Header Empty", got)
		assert.Equal(t, 501, resp.StatusCode)
		assert.True(t, resp.Error.Status)
	})

	t.Run("InsertSubscriptionDetail repo error", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetOneSubscriptionFn: func(_ context.Context, _ int) (goldEntity.Subscription, error) {
				return goldEntity.Subscription{GoldNamaPaket: "Basic"}, nil
			},
			GetSubscriptionHeaderFn: func(_ context.Context, _ int) (goldEntity.SubscriptionHeader, error) {
				return goldEntity.SubscriptionHeader{GoldID: 1}, nil
			},
			InsertSubscriptionDetailFn: func(_ context.Context, _ goldEntity.SubscriptionDetail) error {
				return errors.New("insert failed")
			},
		})

		input := goldEntity.SubscriptionDetail{GoldId: 1, GoldMenuId: 1}
		got, err, _ := svc.InsertSubscriptionDetail(context.Background(), input)
		assert.Equal(t, "Gagal", got)
		assert.Error(t, err)
	})
}

// --- GetSubscriptionHeaderTotalHarga ---

func TestGetSubscriptionHeaderTotalHarga(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return goldEntity.GetGoldUserss{GoldId: 1}, nil
			},
			GetSubscriptionHeaderTotalHargaFn: func(_ context.Context, id int) (goldEntity.SubscriptionHeaderPayment, error) {
				assert.Equal(t, 1, id) // verifikasi id dari user
				return goldEntity.SubscriptionHeaderPayment{GoldValidasiPayment: "Y"}, nil
			},
		})

		got, err := svc.GetSubscriptionHeaderTotalHarga(context.Background(), "budi@test.com")
		assert.NoError(t, err)
		assert.Equal(t, "Y", got.GoldValidasiPayment)
	})

	t.Run("GetGoldUserByEmail error", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return goldEntity.GetGoldUserss{}, errors.New("not found")
			},
		})

		_, err := svc.GetSubscriptionHeaderTotalHarga(context.Background(), "budi@test.com")
		assert.Error(t, err)
	})

	t.Run("GetSubscriptionHeaderTotalHarga repo error", func(t *testing.T) {
		svc := newTestService(&mockRepo{
			GetGoldUserByEmailFn: func(_ context.Context, _ string) (goldEntity.GetGoldUserss, error) {
				return goldEntity.GetGoldUserss{GoldId: 1}, nil
			},
			GetSubscriptionHeaderTotalHargaFn: func(_ context.Context, _ int) (goldEntity.SubscriptionHeaderPayment, error) {
				return goldEntity.SubscriptionHeaderPayment{}, errors.New("query failed")
			},
		})

		_, err := svc.GetSubscriptionHeaderTotalHarga(context.Background(), "budi@test.com")
		assert.Error(t, err)
	})
}

// --- UploadTestingImages ---

func TestUploadTestingImages(t *testing.T) {
	input := goldEntity.Testings{ID: "img1", TestingImages: []byte("fake-image-data")}

	tests := []struct {
		name    string
		repo    *mockRepo
		want    string
		wantErr bool
	}{
		{
			name: "success",
			repo: &mockRepo{
				UploadTestingImagesFn: func(_ context.Context, _ goldEntity.Testings) (string, error) {
					return "Sukses", nil
				},
			},
			want: "Sukses",
		},
		{
			name: "repo error",
			repo: &mockRepo{
				UploadTestingImagesFn: func(_ context.Context, _ goldEntity.Testings) (string, error) {
					return "", errors.New("upload failed")
				},
			},
			want:    "Gagal",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.UploadTestingImages(context.Background(), input)
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// --- GetTestingImage ---

func TestGetTestingImage(t *testing.T) {
	imageData := []byte("fake-image-bytes")

	tests := []struct {
		name    string
		id      int
		repo    *mockRepo
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			repo: &mockRepo{
				GetTestingImagesFn: func(_ context.Context, _ int) ([]byte, error) {
					return imageData, nil
				},
			},
			want: imageData,
		},
		{
			name: "repo error",
			id:   999,
			repo: &mockRepo{
				GetTestingImagesFn: func(_ context.Context, _ int) ([]byte, error) {
					return nil, errors.New("not found")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService(tt.repo)
			got, err := svc.GetTestingImage(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
