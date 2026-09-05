package goldgym

import (
	"context"
	"database/sql"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

const dbTimeout = 3 * time.Second
const dbTimeoutInsert = 5 * time.Second

func (d *Data) GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error) {
	var (
		users []goldEntity.GetGoldUser
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, err
}

func (d *Data) GetGoldUserByID(ctx context.Context, id string) (goldEntity.GetGoldUserss, error) {
	var (
		user goldEntity.GetGoldUserss
		err  error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Where("gold_id = ?", id).First(&user).Error
	if err != nil {
		return goldEntity.GetGoldUserss{}, err
	}

	return user, err
}

func (d *Data) GetGoldUserByEmail(ctx context.Context, email string) (goldEntity.GetGoldUserss, error) {
	var (
		user goldEntity.GetGoldUserss
		err  error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Where("gold_email = ?", email).First(&user).Error
	if err != nil {
		return goldEntity.GetGoldUserss{}, err
	}

	return user, err
}

func (d *Data) GetGoldUserByEmailLogin(ctx context.Context, email string, password string) (goldEntity.LoginUser, error) {
	var (
		user goldEntity.LoginUser
		err  error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Where("gold_email = ? AND gold_password = ?", email, password).First(&user).Error
	if err != nil {
		return goldEntity.LoginUser{}, err
	}
	return user, err
}

func (d *Data) InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeoutInsert)
	defer cancel()
	tx := d.db.WithContext(ctx)
	// gold_expireddate datetime nullable, tapi GoldExpireddate string kosong
	// ("") bukan NULL yang valid -> ditolak sql_mode strict. Kolom kartu ini
	// legacy & opsional (registrasi modern tidak mengisi kartu), jadi
	// di-omit dari INSERT saat kosong supaya DB pakai NULL.
	if strings.TrimSpace(user.GoldExpireddate) == "" {
		tx = tx.Omit("gold_expireddate")
	}
	err := tx.Create(&user).Error
	if err != nil {
		return "Gagal", err
	}
	return "Sukses", nil
}

func (d *Data) GetGoldToken(ctx context.Context) (goldEntity.LoginToken, error) {
	var (
		user goldEntity.LoginToken
		err  error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).First(&user).Error
	if err != nil {
		return goldEntity.LoginToken{}, err
	}
	return user, err
}

func (d *Data) UpdateGoldToken(ctx context.Context, user goldEntity.LoginTokenDataPeserta) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.LoginToken{}).Where("gold_email = ?", user.GoldEmail).Update("gold_token", user.GoldToken).Error
}

func (d *Data) GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error) {
	var (
		users []goldEntity.Subscription
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return []goldEntity.Subscription{}, err
	}
	return users, err
}

func (d *Data) InsertSubscription(ctx context.Context, user goldEntity.SubscriptionAll) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *Data) InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Create(user).Error

}

func (d *Data) DeleteSubscriptionDetail(ctx context.Context, user goldEntity.DeleteSubs) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Where("gold_id = ? AND gold_menuid = ?", user.GoldId, user.GoldMenuId).Delete(&goldEntity.SubscriptionDetail{}).Error
}

func (d *Data) UpdateSubscriptionDetail(ctx context.Context, user goldEntity.UpdateSubs) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.SubscriptionDetail{}).Where("gold_id = ? AND gold_menuid = ?", user.GoldId, user.GoldMenuId).Update("gold_jumlahpertemuan", user.GoldJumlahpertemuan).Error
}

func (d *Data) UpdateDataPeserta(ctx context.Context, user goldEntity.UpdatePassword) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.GetGoldUser{}).Where("gold_email = ? and gold_otp = ?", user.GoldEmail, user.GoldOTP).Update("gold_password", user.GoldPassword).Error
}

func (d *Data) UpdateNama(ctx context.Context, user goldEntity.UpdateNama) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.GetGoldUser{}).Where("gold_email = ?", user.GoldEmail).Update("gold_nama", user.GoldNama).Error
}

// UpdateGoldBuyerYN menandai akun sudah mendaftar sebagai pembeli
// (flag menu Mode Pembeli di aplikasi).
func (d *Data) UpdateGoldBuyerYN(ctx context.Context, goldid int) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.GetGoldUser{}).Where("gold_id = ?", goldid).Update("gold_buyer_yn", "Y").Error
}

// UpdateGoldToko menyimpan nama toko user login (identitas pembeli di nota).
func (d *Data) UpdateGoldToko(ctx context.Context, goldid int, toko string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.GetGoldUser{}).Where("gold_id = ?", goldid).Update("gold_toko", toko).Error
}

func (d *Data) UpdateKartu(ctx context.Context, user goldEntity.UpdateKartu) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.GetGoldUser{}).Where("gold_email = ?", user.GoldEmail).Updates(map[string]interface{}{
		"gold_nomorkartu": user.GoldNomorKartu,
		"gold_cvv":        user.GoldCvv,
	}).Error
}

func (d *Data) Logout(ctx context.Context, user goldEntity.Logout) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.GetGoldUser{}).Where("gold_email = ?", user.GoldEmail).Update("gold_nama", sql.NullString{String: "", Valid: false}).Error
}

func (d *Data) GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error) {
	var (
		user  goldEntity.GetSubsWithUser
		users []goldEntity.GetSubsWithUser
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	rows, err := (*d.stmt)[getSubsWithUser].QueryxContext(ctx)
	if err != nil {
		return users, errors.Wrap(err, "[DATA] [GetGoldUser]")
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.StructScan(&user); err != nil {
			return users, errors.Wrap(err, "[DATA] [GetGoldUser]")
		}
		users = append(users, user)
	}
	return users, err
}

func (d *Data) GetValidationGoldOTP(ctx context.Context, otp string) (goldEntity.GetValidationGoldOTP, error) {
	var (
		users goldEntity.GetValidationGoldOTP
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Where("gold_otp = ?", otp).First(&users).Error
	if err != nil {
		return goldEntity.GetValidationGoldOTP{}, err
	}
	return users, err
}

func (d *Data) UpdateValidationOTP(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.UpdateValidationOTP{}).Where("gold_email = ? AND gold_otp IS NOT NULL", email).Update("gold_validasiyn", "Y").Error
}

func (d *Data) UpdateOtpIsNull(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.UpdateValidationOTP{}).Where("gold_email = ?", email).Updates(map[string]interface{}{
		"gold_otp":            nil,
		"gold_otp_created_at": nil,
	}).Error
}

func (d *Data) UpdateOTP(ctx context.Context, otp string, email string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.UpdateValidationOTP{}).Where("gold_email = ?", email).Updates(map[string]interface{}{
		"gold_otp":            otp,
		"gold_otp_created_at": gorm.Expr("NOW()"),
	}).Error
}

func (d *Data) GetOneSubscription(ctx context.Context, menuid int) (goldEntity.Subscription, error) {
	var (
		users goldEntity.Subscription
		err   error
	)

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Where("gold_menuid = ?", menuid).First(&users).Error
	if err != nil {
		return goldEntity.Subscription{}, err
	}

	return users, err
}

func (d *Data) UpdateOTPSubscription(ctx context.Context, otp string, id int) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.SubscriptionAll{}).Where("gold_id = ?", id).Updates(map[string]interface{}{
		"gold_otp":        otp,
		"gold_lastupdate": gorm.Expr("NOW()"),
	}).Error
}

func (d *Data) BulkInsertSubscriptionDetail(ctx context.Context, user []goldEntity.SubscriptionDetail) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	for _, v := range user {
		query, args, err := sqlx.In(qInsertSubscriptionDetail,
			v.GoldId, v.GoldMenuId, v.GoldNamaPaket, v.GoldNamaLayanan, v.GoldHarga, v.GoldJadwal, v.GoldListLatihan, v.GoldJumlahpertemuan, v.GoldDurasi, v.GoldStatuslangganan)
		if err != nil {
			return errors.Wrap(err, "[DATA][BulkInsertSubscriptionDetail]")
		}
		_, err = d.dbr.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.Wrap(err, "[DATA][BulkInsertSubscriptionDetail]")
		}
	}
	return nil
}

func (d Data) GetSubscriptionHeader(ctx context.Context, id int) (goldEntity.SubscriptionHeader, error) {
	var (
		users goldEntity.SubscriptionHeader
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Where("gold_id = ?", id).First(&users).Error
	if err != nil {
		return goldEntity.SubscriptionHeader{}, err
	}
	return users, err
}

func (d Data) UpdateValidasiPaymentHeader(ctx context.Context, updatePayment goldEntity.UpdatePayment) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.SubscriptionAll{}).Where("gold_id = ?", updatePayment.GoldID).Updates(map[string]interface{}{
		"gold_validasipayment": "Y",
		"gold_lastupdate":      gorm.Expr("NOW()"),
	}).Error

}

func (d Data) UpdateValidasiPaymentDetail(ctx context.Context, updatePayment goldEntity.UpdatePayment) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.SubscriptionDetail{}).Where("gold_id = ?", updatePayment.GoldID).Updates(map[string]interface{}{
		"gold_startdate":       gorm.Expr("NOW()"),
		"gold_enddate":         gorm.Expr("DATE_ADD(NOW(), INTERVAL 30 DAY)"),
		"gold_statuslangganan": "Berlangganan",
	}).Error
}

func (d Data) GetSubscriptionHeaderTotalHarga(ctx context.Context, id int) (goldEntity.SubscriptionHeaderPayment, error) {
	var (
		users goldEntity.SubscriptionHeaderPayment
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Where("gold_id = ?", id).First(&users).Error
	if err != nil {
		return goldEntity.SubscriptionHeaderPayment{}, err
	}
	return users, err
}

func (d Data) GetPasswordByUser(ctx context.Context, _user string) (string, error) {
	var (
		users goldEntity.GetGoldUser
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Where("gold_email = ?", _user).First(&users).Error
	if err != nil {
		return "", err
	}
	return users.GoldPassword, nil
}

func (d Data) UpdateLastLogin(ctx context.Context, _user goldEntity.GetGoldUserss) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Model(&goldEntity.GetGoldUser{}).Where("gold_email = ?", _user.GoldEmail).Updates(map[string]interface{}{
		"gold_last_login":      gorm.Expr("NOW()"),
		"gold_last_login_host": _user.GoldLastLoginHost,
	}).Error
}

func (d Data) UploadTestingImages(ctx context.Context, testing goldEntity.Testings) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Create(&testing).Error
	if err != nil {
		return "Gagal", err
	}

	return "Sukses", err

}

func (d Data) GetTestingImages(ctx context.Context, id int) ([]byte, error) {
	var image []byte
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Where("id = ? ", id).Scan(&image).Error
	if err != nil {
		return []byte{}, err
	}

	return image, nil
}

func (d *Data) InsertBuyer(ctx context.Context, buyer goldEntity.BuyerRow) (goldEntity.BuyerRow, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeoutInsert)
	defer cancel()
	err := d.db.WithContext(ctx).Create(&buyer).Error
	if err != nil {
		return buyer, err
	}
	return buyer, nil
}

// CountGoldUsers total seluruh akun terdaftar (semua role), dipakai untuk
// menegakkan batas MaxRegisteredUsers saat registrasi.
func (d *Data) CountGoldUsers(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var total int64
	err := d.db.WithContext(ctx).Model(&goldEntity.GetGoldUser{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetRegistrationMode mode pendaftaran mandiri saat ini (app_settings.registration_mode):
// BOTH (default, bisa pilih pembeli/penjual), BUYER_ONLY, atau SELLER_ONLY.
// Baris belum ada / nilai kosong dianggap BOTH (default aman: perilaku
// existing tidak berubah sampai admin eksplisit mengubah).
func (d *Data) GetRegistrationMode(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var value string
	err := d.db.WithContext(ctx).
		Table("app_settings").
		Select("setting_value").
		Where("setting_key = ?", "registration_mode").
		Limit(1).
		Scan(&value).Error
	if err != nil {
		return "", err
	}
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return "BOTH", nil
	}
	return value, nil
}

func (d *Data) SetRegistrationMode(ctx context.Context, mode string, updatedBy string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).
		Exec(`INSERT INTO app_settings (setting_key, setting_value, updated_by)
              VALUES (?, ?, ?)
              ON DUPLICATE KEY UPDATE setting_value = VALUES(setting_value), updated_by = VALUES(updated_by)`,
			"registration_mode", mode, updatedBy).Error
}
