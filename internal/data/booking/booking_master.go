package booking

import (
	"context"
	"time"

	bookingEntity "gold-gym-be/internal/entity/booking"

	"gorm.io/gorm"
)

const dbTimeout = 5 * time.Second

// ExpireOverdueBookings menandai booking UNPAID yang jamnya sudah lewat menjadi EXPIRED
// sehingga slotnya kembali available (hijau di frontend).
func (d *Data) ExpireOverdueBookings(ctx context.Context, now time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).
		Model(&bookingEntity.Booking{}).
		Where("booking_status = ? AND booking_start < ?", bookingEntity.StatusUnpaid, now).
		Update("booking_status", bookingEntity.StatusExpired).Error
}

// GetBookingsByDate mengambil booking aktif (PAID/UNPAID) satu outlet pada satu tanggal.
func (d *Data) GetBookingsByDate(ctx context.Context, goldid int, outcode string, date time.Time) ([]bookingEntity.Booking, error) {
	var bookings []bookingEntity.Booking
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("booking_gold_id = ? AND booking_outcode = ? AND booking_date = ? AND booking_status IN ?",
			goldid, outcode, date.Format("2006-01-02"),
			[]string{bookingEntity.StatusPaid, bookingEntity.StatusUnpaid}).
		Order("booking_start asc").
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

// CountActiveInWindow menghitung booking aktif yang menempati jendela 30 menit [winStart, winStart+30m).
// Booking menempati jendela jika booking_start <= winStart < booking_start + duration.
func (d *Data) CountActiveInWindow(ctx context.Context, goldid int, outcode string, winStart time.Time) (int, error) {
	var count int64
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Model(&bookingEntity.Booking{}).
		Where("booking_gold_id = ? AND booking_outcode = ?", goldid, outcode).
		Where("booking_status IN ?", []string{bookingEntity.StatusPaid, bookingEntity.StatusUnpaid}).
		Where("booking_start <= ? AND DATE_ADD(booking_start, INTERVAL booking_duration MINUTE) > ?", winStart, winStart).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (d *Data) InsertBooking(ctx context.Context, b bookingEntity.Booking) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Create(&b).Error
}

func (d *Data) GetBookingByID(ctx context.Context, bookingID string) (*bookingEntity.Booking, error) {
	var b bookingEntity.Booking
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Where("booking_id = ?", bookingID).First(&b).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// UpdateBookingPaid menandai booking lunas + link ke th_sale.
func (d *Data) UpdateBookingPaid(ctx context.Context, bookingID string, saleID string, itemID int, itemName string, price string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).
		Model(&bookingEntity.Booking{}).
		Where("booking_id = ?", bookingID).
		Updates(map[string]interface{}{
			"booking_status":    bookingEntity.StatusPaid,
			"booking_sale_id":   saleID,
			"booking_item_id":   itemID,
			"booking_item_name": itemName,
			"booking_price":     price,
		}).Error
}

func (d *Data) UpdateBookingStatus(ctx context.Context, bookingID string, status string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).
		Model(&bookingEntity.Booking{}).
		Where("booking_id = ?", bookingID).
		Update("booking_status", status).Error
}

// SearchBuyers autocomplete pembeli terdaftar (role BUYER) untuk penjual.
func (d *Data) SearchBuyers(ctx context.Context, name string, limit int) ([]bookingEntity.BuyerInfo, error) {
	var buyers []bookingEntity.BuyerInfo
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	q := d.db.WithContext(ctx).
		Model(&bookingEntity.BuyerInfo{}).
		Where("gold_role = ?", "BUYER")
	if name != "" {
		q = q.Where("(gold_nama LIKE ? OR gold_email LIKE ?)", "%"+name+"%", "%"+name+"%")
	}
	err := q.Order("gold_nama asc").Limit(limit).Find(&buyers).Error
	if err != nil {
		return nil, err
	}
	return buyers, nil
}

func (d *Data) GetBuyerByID(ctx context.Context, goldid int) (*bookingEntity.BuyerInfo, error) {
	var buyer bookingEntity.BuyerInfo
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("gold_id = ? AND gold_role = ?", goldid, "BUYER").
		First(&buyer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &buyer, nil
}

// GetTherapyItem mengambil item terapi (harga booking) dari tabel items.
func (d *Data) GetTherapyItem(ctx context.Context, itemID int) (*bookingEntity.TherapyItem, error) {
	var item bookingEntity.TherapyItem
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Where("item_id = ?", itemID).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetTherapyItemByName mengambil item jasa terapi berdasarkan nama (per outlet),
// dipakai untuk memetakan tipe terapi + durasi → item & harga.
func (d *Data) GetTherapyItemByName(ctx context.Context, outcode string, name string) (*bookingEntity.TherapyItem, error) {
	var item bookingEntity.TherapyItem
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("item_outcode = ? AND UPPER(item_name) = ? AND UPPER(item_brand) = 'THERAPY'", outcode, name).
		First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// FindAdjacentUnpaid mencari booking 30 menit UNPAID milik nama & tipe terapi yang sama
// tepat 30 menit sebelum/sesudah start (untuk digabung jadi 1 jam saat pembayaran).
func (d *Data) FindAdjacentUnpaid(ctx context.Context, goldid int, outcode string, custName string, therapyType string, excludeID string, start time.Time) (*bookingEntity.Booking, error) {
	var b bookingEntity.Booking
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("booking_gold_id = ? AND booking_outcode = ?", goldid, outcode).
		Where("booking_status = ? AND booking_duration = 30", bookingEntity.StatusUnpaid).
		Where("booking_cust_name = ? AND booking_therapy_type = ?", custName, therapyType).
		Where("booking_id <> ?", excludeID).
		Where("booking_start IN ?", []time.Time{start.Add(-30 * time.Minute), start.Add(30 * time.Minute)}).
		Order("booking_start asc").
		First(&b).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// RemoveBookingWithLog menghapus booking dan mencatatnya ke booking_remove_log
// dalam satu transaksi (jejak audit penghapusan oleh penjual/admin).
func (d *Data) RemoveBookingWithLog(ctx context.Context, b bookingEntity.Booking, removedBy string, removedRole string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		date := b.BookingDate
		start := b.BookingStart
		logRow := bookingEntity.RemoveLog{
			LogBookingID:    b.BookingID,
			LogGoldID:       b.BookingGoldID,
			LogOutcode:      b.BookingOutcode,
			LogBookingDate:  &date,
			LogBookingStart: &start,
			LogDuration:     b.BookingDuration,
			LogTherapyType:  b.BookingTherapyType,
			LogCustID:       b.BookingCustID,
			LogCustName:     b.BookingCustName,
			LogStatus:       b.BookingStatus,
			LogRemovedBy:    removedBy,
			LogRemovedRole:  removedRole,
		}
		if err := tx.Create(&logRow).Error; err != nil {
			return err
		}
		return tx.Where("booking_id = ?", b.BookingID).Delete(&bookingEntity.Booking{}).Error
	})
}

// GetStockIDByItem mencari stock_id item di suatu outlet, untuk menjaga relasi
// td_sale.sale_stockid → stock.stock_id → items.item_id pada penjualan booking.
func (d *Data) GetStockIDByItem(ctx context.Context, outcode string, itemID int) (string, error) {
	var stockID string
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Table("stock").
		Select("stock_id").
		Where("stock_outcode = ? AND stock_item_id = ?", outcode, itemID).
		Limit(1).
		Scan(&stockID).Error
	if err != nil {
		return "", err
	}
	return stockID, nil
}

// GetOutletByCode mengambil outlet (gold_id tenant + tipe) berdasarkan kode outlet.
func (d *Data) GetOutletByCode(ctx context.Context, outcode string) (*bookingEntity.OutletInfo, error) {
	var outlet bookingEntity.OutletInfo
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Where("outlet_code = ?", outcode).First(&outlet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &outlet, nil
}
