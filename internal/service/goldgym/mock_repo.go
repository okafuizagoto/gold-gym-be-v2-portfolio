package goldgym

import (
	"context"

	goldEntity "gold-gym-be/internal/entity/goldgym"
)

// mockRepo mengimplementasikan interface RepoData untuk unit test.
//
// Pola penggunaan:
//   - Setiap method di RepoData punya satu field fungsi (Fn).
//   - Jika field Fn nilal, method kembalikan zero value.
//   - Di test, set hanya field Fn yang dibutuhkan.
//
// Cara menambah method baru (setelah menambah method baru di interface RepoData):
//  1. Tambah field fungsi baru mengikuti pola di bawah, misal: NamaMethodFn func(...) (...).
//  2. Tambah implementasi method-nya mengikuti pola yang sudah ada (4 baris).
type mockRepo struct {
	GetGoldUserFn                     func(ctx context.Context) ([]goldEntity.GetGoldUser, error)
	InsertGoldUserFn                  func(ctx context.Context, user goldEntity.GetGoldUsers) (string, error)
	GetGoldUserByEmailFn              func(ctx context.Context, email string) (goldEntity.GetGoldUserss, error)
	GetGoldUserByEmailLoginFn         func(ctx context.Context, email string, password string) (goldEntity.LoginUser, error)
	GetGoldTokenFn                    func(ctx context.Context) (goldEntity.LoginToken, error)
	UpdateGoldTokenFn                 func(ctx context.Context, user goldEntity.LoginTokenDataPeserta) error
	GetAllSubscriptionFn              func(ctx context.Context) ([]goldEntity.Subscription, error)
	InsertSubscriptionFn              func(ctx context.Context, user goldEntity.SubscriptionAll) error
	InsertSubscriptionDetailFn        func(ctx context.Context, user goldEntity.SubscriptionDetail) error
	DeleteSubscriptionDetailFn        func(ctx context.Context, user goldEntity.DeleteSubs) error
	UpdateSubscriptionDetailFn        func(ctx context.Context, user goldEntity.UpdateSubs) error
	UpdateDataPesertaFn               func(ctx context.Context, user goldEntity.UpdatePassword) error
	UpdateNamaFn                      func(ctx context.Context, user goldEntity.UpdateNama) error
	UpdateKartuFn                     func(ctx context.Context, user goldEntity.UpdateKartu) error
	LogoutFn                          func(ctx context.Context, user goldEntity.Logout) error
	GetSubsWithUserFn                 func(ctx context.Context) ([]goldEntity.GetSubsWithUser, error)
	GetValidationGoldOTPFn            func(ctx context.Context, otp string) (goldEntity.GetValidationGoldOTP, error)
	UpdateValidationOTPFn             func(ctx context.Context, email string) error
	UpdateOtpIsNullFn                 func(ctx context.Context, email string) error
	UpdateOTPFn                       func(ctx context.Context, otp string, email string) error
	GetOneSubscriptionFn              func(ctx context.Context, menuid int) (goldEntity.Subscription, error)
	BulkInsertSubscriptionDetailFn    func(ctx context.Context, user []goldEntity.SubscriptionDetail) error
	UpdateOTPSubscriptionFn           func(ctx context.Context, otp string, id int) error
	GetSubscriptionHeaderFn           func(ctx context.Context, id int) (goldEntity.SubscriptionHeader, error)
	UpdateValidasiPaymentHeaderFn     func(ctx context.Context, updatePayment goldEntity.UpdatePayment) error
	UpdateValidasiPaymentDetailFn     func(ctx context.Context, updatePayment goldEntity.UpdatePayment) error
	GetSubscriptionHeaderTotalHargaFn func(ctx context.Context, id int) (goldEntity.SubscriptionHeaderPayment, error)
	GetPasswordByUserFn               func(ctx context.Context, _user string) (string, error)
	UpdateLastLoginFn                 func(ctx context.Context, _user goldEntity.GetGoldUserss) error
	UploadTestingImagesFn             func(ctx context.Context, testing goldEntity.Testings) (string, error)
	GetTestingImagesFn                func(ctx context.Context, id int) ([]byte, error)
	GetGoldUserByIDFn                 func(ctx context.Context, id string) (goldEntity.GetGoldUserss, error)
	CountGoldUsersFn                  func(ctx context.Context) (int64, error)
	GetRegistrationModeFn             func(ctx context.Context) (string, error)
	SetRegistrationModeFn             func(ctx context.Context, mode string, updatedBy string) error
}

func (m *mockRepo) GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
	if m.GetGoldUserFn != nil {
		return m.GetGoldUserFn(ctx)
	}
	return nil, nil
}

func (m *mockRepo) InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (string, error) {
	if m.InsertGoldUserFn != nil {
		return m.InsertGoldUserFn(ctx, user)
	}
	return "", nil
}

func (m *mockRepo) GetGoldUserByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
	if m.GetGoldUserByEmailFn != nil {
		return m.GetGoldUserByEmailFn(ctx, email)
	}
	return goldEntity.GetGoldUserss{}, nil
}

func (m *mockRepo) GetGoldUserByEmailLogin(ctx context.Context, email string, password string) (goldEntity.LoginUser, error) {
	if m.GetGoldUserByEmailLoginFn != nil {
		return m.GetGoldUserByEmailLoginFn(ctx, email, password)
	}
	return goldEntity.LoginUser{}, nil
}

func (m *mockRepo) GetGoldToken(ctx context.Context) (goldEntity.LoginToken, error) {
	if m.GetGoldTokenFn != nil {
		return m.GetGoldTokenFn(ctx)
	}
	return goldEntity.LoginToken{}, nil
}

func (m *mockRepo) UpdateGoldToken(ctx context.Context, user goldEntity.LoginTokenDataPeserta) error {
	if m.UpdateGoldTokenFn != nil {
		return m.UpdateGoldTokenFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error) {
	if m.GetAllSubscriptionFn != nil {
		return m.GetAllSubscriptionFn(ctx)
	}
	return nil, nil
}

func (m *mockRepo) InsertSubscription(ctx context.Context, user goldEntity.SubscriptionAll) error {
	if m.InsertSubscriptionFn != nil {
		return m.InsertSubscriptionFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) error {
	if m.InsertSubscriptionDetailFn != nil {
		return m.InsertSubscriptionDetailFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) DeleteSubscriptionDetail(ctx context.Context, user goldEntity.DeleteSubs) error {
	if m.DeleteSubscriptionDetailFn != nil {
		return m.DeleteSubscriptionDetailFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) UpdateSubscriptionDetail(ctx context.Context, user goldEntity.UpdateSubs) error {
	if m.UpdateSubscriptionDetailFn != nil {
		return m.UpdateSubscriptionDetailFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) UpdateDataPeserta(ctx context.Context, user goldEntity.UpdatePassword) error {
	if m.UpdateDataPesertaFn != nil {
		return m.UpdateDataPesertaFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) UpdateNama(ctx context.Context, user goldEntity.UpdateNama) error {
	if m.UpdateNamaFn != nil {
		return m.UpdateNamaFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) UpdateGoldToko(ctx context.Context, goldid int, toko string) error {
	return nil
}

func (m *mockRepo) UpdateGoldBuyerYN(ctx context.Context, goldid int) error {
	return nil
}

func (m *mockRepo) UpdateKartu(ctx context.Context, user goldEntity.UpdateKartu) error {
	if m.UpdateKartuFn != nil {
		return m.UpdateKartuFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) Logout(ctx context.Context, user goldEntity.Logout) error {
	if m.LogoutFn != nil {
		return m.LogoutFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error) {
	if m.GetSubsWithUserFn != nil {
		return m.GetSubsWithUserFn(ctx)
	}
	return nil, nil
}

func (m *mockRepo) GetValidationGoldOTP(ctx context.Context, otp string) (goldEntity.GetValidationGoldOTP, error) {
	if m.GetValidationGoldOTPFn != nil {
		return m.GetValidationGoldOTPFn(ctx, otp)
	}
	return goldEntity.GetValidationGoldOTP{}, nil
}

func (m *mockRepo) UpdateValidationOTP(ctx context.Context, email string) error {
	if m.UpdateValidationOTPFn != nil {
		return m.UpdateValidationOTPFn(ctx, email)
	}
	return nil
}

func (m *mockRepo) UpdateOtpIsNull(ctx context.Context, email string) error {
	if m.UpdateOtpIsNullFn != nil {
		return m.UpdateOtpIsNullFn(ctx, email)
	}
	return nil
}

func (m *mockRepo) UpdateOTP(ctx context.Context, otp string, email string) error {
	if m.UpdateOTPFn != nil {
		return m.UpdateOTPFn(ctx, otp, email)
	}
	return nil
}

func (m *mockRepo) GetOneSubscription(ctx context.Context, menuid int) (goldEntity.Subscription, error) {
	if m.GetOneSubscriptionFn != nil {
		return m.GetOneSubscriptionFn(ctx, menuid)
	}
	return goldEntity.Subscription{}, nil
}

func (m *mockRepo) BulkInsertSubscriptionDetail(ctx context.Context, user []goldEntity.SubscriptionDetail) error {
	if m.BulkInsertSubscriptionDetailFn != nil {
		return m.BulkInsertSubscriptionDetailFn(ctx, user)
	}
	return nil
}

func (m *mockRepo) UpdateOTPSubscription(ctx context.Context, otp string, id int) error {
	if m.UpdateOTPSubscriptionFn != nil {
		return m.UpdateOTPSubscriptionFn(ctx, otp, id)
	}
	return nil
}

func (m *mockRepo) GetSubscriptionHeader(ctx context.Context, id int) (goldEntity.SubscriptionHeader, error) {
	if m.GetSubscriptionHeaderFn != nil {
		return m.GetSubscriptionHeaderFn(ctx, id)
	}
	return goldEntity.SubscriptionHeader{}, nil
}

func (m *mockRepo) UpdateValidasiPaymentHeader(ctx context.Context, updatePayment goldEntity.UpdatePayment) error {
	if m.UpdateValidasiPaymentHeaderFn != nil {
		return m.UpdateValidasiPaymentHeaderFn(ctx, updatePayment)
	}
	return nil
}

func (m *mockRepo) UpdateValidasiPaymentDetail(ctx context.Context, updatePayment goldEntity.UpdatePayment) error {
	if m.UpdateValidasiPaymentDetailFn != nil {
		return m.UpdateValidasiPaymentDetailFn(ctx, updatePayment)
	}
	return nil
}

func (m *mockRepo) GetSubscriptionHeaderTotalHarga(ctx context.Context, id int) (goldEntity.SubscriptionHeaderPayment, error) {
	if m.GetSubscriptionHeaderTotalHargaFn != nil {
		return m.GetSubscriptionHeaderTotalHargaFn(ctx, id)
	}
	return goldEntity.SubscriptionHeaderPayment{}, nil
}

func (m *mockRepo) GetPasswordByUser(ctx context.Context, _user string) (string, error) {
	if m.GetPasswordByUserFn != nil {
		return m.GetPasswordByUserFn(ctx, _user)
	}
	return "", nil
}

func (m *mockRepo) UpdateLastLogin(ctx context.Context, _user goldEntity.GetGoldUserss) error {
	if m.UpdateLastLoginFn != nil {
		return m.UpdateLastLoginFn(ctx, _user)
	}
	return nil
}

func (m *mockRepo) UploadTestingImages(ctx context.Context, testing goldEntity.Testings) (string, error) {
	if m.UploadTestingImagesFn != nil {
		return m.UploadTestingImagesFn(ctx, testing)
	}
	return "", nil
}

func (m *mockRepo) GetTestingImages(ctx context.Context, id int) ([]byte, error) {
	if m.GetTestingImagesFn != nil {
		return m.GetTestingImagesFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepo) GetGoldUserByID(ctx context.Context, id string) (goldEntity.GetGoldUserss, error) {
	if m.GetGoldUserByIDFn != nil {
		return m.GetGoldUserByIDFn(ctx, id)
	}
	return goldEntity.GetGoldUserss{}, nil
}

func (m *mockRepo) InsertBuyer(ctx context.Context, buyer goldEntity.BuyerRow) (goldEntity.BuyerRow, error) {
	return buyer, nil
}

func (m *mockRepo) CountGoldUsers(ctx context.Context) (int64, error) {
	if m.CountGoldUsersFn != nil {
		return m.CountGoldUsersFn(ctx)
	}
	return 0, nil
}

func (m *mockRepo) GetRegistrationMode(ctx context.Context) (string, error) {
	if m.GetRegistrationModeFn != nil {
		return m.GetRegistrationModeFn(ctx)
	}
	return "BOTH", nil
}

func (m *mockRepo) SetRegistrationMode(ctx context.Context, mode string, updatedBy string) error {
	if m.SetRegistrationModeFn != nil {
		return m.SetRegistrationModeFn(ctx, mode, updatedBy)
	}
	return nil
}
