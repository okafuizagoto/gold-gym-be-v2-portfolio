package goldgym

import (
	"context"
	"errors"
	"time"

	discountEntity "gold-gym-be/internal/entity/discount"

	"gorm.io/gorm"
)

const dbTimeout = 3 * time.Second

var (
	errVoucherNotFound = errors.New("kode voucher tidak ditemukan atau sudah tidak valid")
	errVoucherExpired  = errors.New("kode voucher sudah kedaluwarsa")
)

// GetItemsForOutlet = sumber item-picker saat Tambah Diskon (item aktif milik
// outlet ini saja).
func (d *Data) GetItemsForOutlet(ctx context.Context, goldid int, outcode, name string) ([]discountEntity.ItemForOutlet, error) {
	var items []discountEntity.ItemForOutlet
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	q := d.db.WithContext(ctx).Table("items").
		Where("item_gold_id = ? AND item_outcode = ? AND item_status = ?", goldid, outcode, "ACTIVE")
	if name != "" {
		q = q.Where("item_name LIKE ?", "%"+name+"%")
	}
	err := q.Order("item_name asc").Find(&items).Error
	return items, err
}

func (d *Data) GetDiscounts(ctx context.Context, goldid int, outcode, name string, page, length int) ([]discountEntity.Discount, error) {
	var rows []discountEntity.Discount
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	offset := (page - 1) * length
	q := d.db.WithContext(ctx).
		Where("discount_gold_id = ? AND discount_outcode = ?", goldid, outcode)
	if name != "" {
		q = q.Where("discount_item_name LIKE ?", "%"+name+"%")
	}
	err := q.Order("discount_created_at desc").Limit(length).Offset(offset).Find(&rows).Error
	return rows, err
}

func (d *Data) GetTotalDiscounts(ctx context.Context, goldid int, outcode, name string) (int64, error) {
	var total int64
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	q := d.db.WithContext(ctx).Model(&discountEntity.Discount{}).
		Where("discount_gold_id = ? AND discount_outcode = ?", goldid, outcode)
	if name != "" {
		q = q.Where("discount_item_name LIKE ?", "%"+name+"%")
	}
	err := q.Count(&total).Error
	return total, err
}

// GetActiveDiscountsByOutlet = semua diskon ACTIVE di outlet ini -- dipakai
// POS (FE) untuk auto-apply diskon saat item masuk keranjang.
func (d *Data) GetActiveDiscountsByOutlet(ctx context.Context, goldid int, outcode string) ([]discountEntity.Discount, error) {
	var rows []discountEntity.Discount
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("discount_gold_id = ? AND discount_outcode = ? AND discount_status = ?", goldid, outcode, "ACTIVE").
		Find(&rows).Error
	return rows, err
}

// GetActiveTotalDiscount = diskon TOTAL (persen dari total keranjang) yang
// sedang ACTIVE untuk outlet ini -- diambil yang paling baru dibuat kalau
// ada lebih dari satu (tidak ada penegakan "cuma boleh 1" di level DB).
func (d *Data) GetActiveTotalDiscount(ctx context.Context, goldid int, outcode string) (*discountEntity.Discount, error) {
	var row discountEntity.Discount
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("discount_gold_id = ? AND discount_outcode = ? AND discount_scope = ? AND discount_status = ?",
			goldid, outcode, "TOTAL", "ACTIVE").
		Order("discount_created_at desc").
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (d *Data) GetDiscountByID(ctx context.Context, goldid int, discountID int) (*discountEntity.Discount, error) {
	var row discountEntity.Discount
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("discount_id = ? AND discount_gold_id = ?", discountID, goldid).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// GetGoldNameByID = nama aktor untuk dicatat di discount_history (tidak ada
// nama terisi di context request, cuma gold_id).
func (d *Data) GetGoldNameByID(ctx context.Context, goldid int) (string, error) {
	var name string
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Table("data_peserta").
		Select("gold_nama").
		Where("gold_id = ?", goldid).
		Scan(&name).Error
	return name, err
}

func newHistoryRow(discount discountEntity.Discount, action, actorName, actorRole string) discountEntity.DiscountHistory {
	return discountEntity.DiscountHistory{
		HistoryDiscountID:        discount.DiscountID,
		HistoryAction:            action,
		HistoryGoldID:            discount.DiscountGoldID,
		HistoryActorName:         actorName,
		HistoryActorRole:         actorRole,
		HistoryOutcode:           discount.DiscountOutcode,
		HistoryItemID:            discount.DiscountItemID,
		HistoryItemName:          discount.DiscountItemName,
		HistoryDiscountType:      discount.DiscountType,
		HistoryDiscountValue:     discount.DiscountValue,
		HistoryDiscountStatus:    discount.DiscountStatus,
		HistoryDiscountCreatedAt: &discount.DiscountCreatedAt,
	}
}

// InsertDiscountWithLog insert diskon baru + catat 'INSERT' di history, dalam
// satu transaksi -- pola sama seperti RemoveBookingWithLog (booking module).
func (d *Data) InsertDiscountWithLog(ctx context.Context, discount discountEntity.Discount, actorName, actorRole string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&discount).Error; err != nil {
			return err
		}
		hist := newHistoryRow(discount, "INSERT", actorName, actorRole)
		return tx.Create(&hist).Error
	})
	return discount.DiscountID, err
}

// UpdateDiscountWithLog update diskon existing + catat 'UPDATE' (snapshot
// keadaan SETELAH diubah), dalam satu transaksi.
func (d *Data) UpdateDiscountWithLog(ctx context.Context, existing discountEntity.Discount, newType string, newValue float64, newStatus string, actorName, actorRole string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&discountEntity.Discount{}).
			Where("discount_id = ? AND discount_gold_id = ?", existing.DiscountID, existing.DiscountGoldID).
			Updates(map[string]interface{}{
				"discount_type":   newType,
				"discount_value":  newValue,
				"discount_status": newStatus,
			}).Error; err != nil {
			return err
		}
		updated := existing
		updated.DiscountType = newType
		updated.DiscountValue = newValue
		updated.DiscountStatus = newStatus
		hist := newHistoryRow(updated, "UPDATE", actorName, actorRole)
		return tx.Create(&hist).Error
	})
}

// DeleteDiscountWithLog catat 'DELETE' (snapshot keadaan SEBELUM dihapus)
// lalu hapus baris diskon, dalam satu transaksi -- log dulu baru hapus, sama
// urutannya dengan RemoveBookingWithLog.
func (d *Data) DeleteDiscountWithLog(ctx context.Context, existing discountEntity.Discount, actorName, actorRole string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		hist := newHistoryRow(existing, "DELETE", actorName, actorRole)
		if err := tx.Create(&hist).Error; err != nil {
			return err
		}
		return tx.Where("discount_id = ?", existing.DiscountID).Delete(&discountEntity.Discount{}).Error
	})
}

func (d *Data) GetDiscountHistory(ctx context.Context, discountID int, page, length int) ([]discountEntity.DiscountHistory, error) {
	var rows []discountEntity.DiscountHistory
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	offset := (page - 1) * length
	err := d.db.WithContext(ctx).
		Where("history_discount_id = ?", discountID).
		Order("history_changed_at desc").Limit(length).Offset(offset).Find(&rows).Error
	return rows, err
}

func (d *Data) GetTotalDiscountHistory(ctx context.Context, discountID int) (int64, error) {
	var total int64
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Model(&discountEntity.DiscountHistory{}).
		Where("history_discount_id = ?", discountID).Count(&total).Error
	return total, err
}

// ---------------------------------------------------------------------------
// Voucher

func (d *Data) VoucherCodeExists(ctx context.Context, outcode, code string) (bool, error) {
	var count int64
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Model(&discountEntity.Voucher{}).
		Where("voucher_outcode = ? AND voucher_code = ?", outcode, code).
		Count(&count).Error
	return count > 0, err
}

func (d *Data) InsertVoucher(ctx context.Context, voucher discountEntity.Voucher) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Create(&voucher).Error
}

func (d *Data) GetVouchers(ctx context.Context, goldid int, outcode string, page, length int) ([]discountEntity.Voucher, error) {
	var rows []discountEntity.Voucher
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	offset := (page - 1) * length
	err := d.db.WithContext(ctx).
		Where("voucher_gold_id = ? AND voucher_outcode = ?", goldid, outcode).
		Order("voucher_created_at desc").Limit(length).Offset(offset).Find(&rows).Error
	return rows, err
}

func (d *Data) GetTotalVouchers(ctx context.Context, goldid int, outcode string) (int64, error) {
	var total int64
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Model(&discountEntity.Voucher{}).
		Where("voucher_gold_id = ? AND voucher_outcode = ?", goldid, outcode).Count(&total).Error
	return total, err
}

// GetVoucherByCode = pratinjau voucher TANPA mengonsumsinya -- dipakai FE
// untuk tampilkan potongan sebelum kasir input jumlah bayar. Validasi
// kedaluwarsa final tetap di RedeemVoucherWithLog saat benar-benar dipakai.
func (d *Data) GetVoucherByCode(ctx context.Context, goldid int, outcode, code string) (*discountEntity.Voucher, error) {
	var row discountEntity.Voucher
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("voucher_gold_id = ? AND voucher_outcode = ? AND voucher_code = ?", goldid, outcode, code).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (d *Data) GetVoucherByID(ctx context.Context, goldid, voucherID int) (*discountEntity.Voucher, error) {
	var row discountEntity.Voucher
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("voucher_id = ? AND voucher_gold_id = ?", voucherID, goldid).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func newVoucherHistoryRow(v discountEntity.Voucher, status string, saleID *string, actorName, actorRole string) discountEntity.VoucherHistory {
	return discountEntity.VoucherHistory{
		HistoryVoucherCode:      v.VoucherCode,
		HistoryOutcode:          v.VoucherOutcode,
		HistoryGoldID:           v.VoucherGoldID,
		HistoryPercent:          v.VoucherPercent,
		HistoryStatus:           status,
		HistorySaleID:           saleID,
		HistoryActorName:        actorName,
		HistoryActorRole:        actorRole,
		HistoryVoucherCreatedAt: &v.VoucherCreatedAt,
	}
}

// DeleteVoucherWithLog = penghapusan manual oleh penjual (belum sempat
// dipakai) -- log 'DELETED' dulu baru hapus baris, sama pola dengan diskon.
func (d *Data) DeleteVoucherWithLog(ctx context.Context, v discountEntity.Voucher, actorName, actorRole string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		hist := newVoucherHistoryRow(v, "DELETED", nil, actorName, actorRole)
		if err := tx.Create(&hist).Error; err != nil {
			return err
		}
		return tx.Where("voucher_id = ?", v.VoucherID).Delete(&discountEntity.Voucher{}).Error
	})
}

// RedeemVoucherWithLog memvalidasi & mengonsumsi voucher secara ATOMIK --
// dipanggil dari modul sales SEBELUM insert sale di-antre ke Kafka (bukan di
// worker async), supaya kasir dapat konfirmasi langsung kalau kode sudah
// tidak valid, dan supaya penggunaan ganda (concurrent) tercegah:
// DELETE...WHERE dulu, cek RowsAffected -- kalau 0 berarti sudah dipakai
// request lain lebih dulu (race), baris history JUGA tidak ditulis kalau
// delete gagal (tidak ada apa pun yang perlu di-log untuk percobaan gagal).
// Return: persen diskon, error (termasuk error "voucher tidak ditemukan/
// sudah tidak valid" dan "voucher sudah kedaluwarsa").
func (d *Data) RedeemVoucherWithLog(ctx context.Context, goldid int, outcode, code, saleID, actorName, actorRole string) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var percent float64
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var v discountEntity.Voucher
		if err := tx.Where("voucher_gold_id = ? AND voucher_outcode = ? AND voucher_code = ?", goldid, outcode, code).
			Take(&v).Error; err != nil {
			return errVoucherNotFound
		}
		if v.VoucherExpiredAt != nil && v.VoucherExpiredAt.Before(time.Now()) {
			// kedaluwarsa: log EXPIRED, hapus baris, tolak pemakaian
			hist := newVoucherHistoryRow(v, "EXPIRED", nil, actorName, actorRole)
			tx.Create(&hist)
			tx.Where("voucher_id = ?", v.VoucherID).Delete(&discountEntity.Voucher{})
			return errVoucherExpired
		}
		result := tx.Where("voucher_id = ?", v.VoucherID).Delete(&discountEntity.Voucher{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errVoucherNotFound // sudah dikonsumsi request lain (race)
		}
		sid := saleID
		hist := newVoucherHistoryRow(v, "USED", &sid, actorName, actorRole)
		if err := tx.Create(&hist).Error; err != nil {
			return err
		}
		percent = v.VoucherPercent
		return nil
	})
	return percent, err
}

func (d *Data) GetVoucherHistory(ctx context.Context, goldid int, outcode string, page, length int) ([]discountEntity.VoucherHistory, error) {
	var rows []discountEntity.VoucherHistory
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	offset := (page - 1) * length
	err := d.db.WithContext(ctx).
		Where("history_gold_id = ? AND history_outcode = ?", goldid, outcode).
		Order("history_changed_at desc").Limit(length).Offset(offset).Find(&rows).Error
	return rows, err
}

func (d *Data) GetTotalVoucherHistory(ctx context.Context, goldid int, outcode string) (int64, error) {
	var total int64
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Model(&discountEntity.VoucherHistory{}).
		Where("history_gold_id = ? AND history_outcode = ?", goldid, outcode).Count(&total).Error
	return total, err
}
