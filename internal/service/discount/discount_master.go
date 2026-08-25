package goldgym

import (
	"context"
	"math"
	"math/rand"
	"strings"
	"time"

	discountEntity "gold-gym-be/internal/entity/discount"

	"github.com/pkg/errors"
)

func canManageDiscount(role string) error {
	if role != "SELLER" && role != "ADMIN" {
		return errors.New("hanya penjual atau admin yang bisa mengelola diskon")
	}
	return nil
}

func validateDiscount(discType string, value float64) error {
	switch discType {
	case "PERCENT":
		if value <= 0 || value > 100 {
			return errors.New("nilai diskon persen harus antara 1-100")
		}
	case "NOMINAL":
		if value <= 0 {
			return errors.New("nilai diskon nominal harus lebih dari 0")
		}
	default:
		return errors.New("tipe diskon harus PERCENT atau NOMINAL")
	}
	return nil
}

func (s Service) GetItemsForOutlet(ctx context.Context, goldid int, outcode, name string) ([]discountEntity.ItemForOutlet, error) {
	items, err := s.discount.GetItemsForOutlet(ctx, goldid, outcode, name)
	return items, errors.Wrap(err, "[Service][GetItemsForOutlet]")
}

func (s Service) GetDiscounts(ctx context.Context, goldid int, outcode, name string, page, length int) ([]discountEntity.Discount, discountEntity.MetadataPaginationDetail, error) {
	var metadata discountEntity.MetadataPaginationDetail
	rows, err := s.discount.GetDiscounts(ctx, goldid, outcode, name, page, length)
	if err != nil {
		return rows, metadata, errors.Wrap(err, "[Service][GetDiscounts]")
	}
	total, err := s.discount.GetTotalDiscounts(ctx, goldid, outcode, name)
	if err != nil {
		return rows, metadata, errors.Wrap(err, "[Service][GetDiscounts][GetTotalDiscounts]")
	}
	metadata = discountEntity.MetadataPaginationDetail{
		Page:      page,
		Limit:     length,
		TotalData: int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(length))),
	}
	return rows, metadata, nil
}

func (s Service) GetActiveDiscountsByOutlet(ctx context.Context, goldid int, outcode string) ([]discountEntity.Discount, error) {
	rows, err := s.discount.GetActiveDiscountsByOutlet(ctx, goldid, outcode)
	return rows, errors.Wrap(err, "[Service][GetActiveDiscountsByOutlet]")
}

func (s Service) InsertDiscount(ctx context.Context, goldid int, role string, items []discountEntity.InsertDiscount) (string, error) {
	if err := canManageDiscount(role); err != nil {
		return "Gagal", err
	}
	if len(items) == 0 {
		return "Data Kosong", nil
	}
	actorName, err := s.discount.GetGoldNameByID(ctx, goldid)
	if err != nil {
		actorName = ""
	}
	for _, it := range items {
		scope := it.DiscountScope
		if scope == "" {
			scope = "ITEM"
		}
		if scope != "ITEM" && scope != "TOTAL" {
			return "Gagal", errors.New("discount_scope harus ITEM atau TOTAL")
		}
		if scope == "TOTAL" && it.DiscountType != "PERCENT" {
			return "Gagal", errors.New("diskon TOTAL hanya boleh bertipe PERCENT")
		}
		if scope == "ITEM" && it.DiscountItemID == 0 {
			return "Gagal", errors.New("discount_item_id wajib diisi untuk diskon per item")
		}
		if err := validateDiscount(it.DiscountType, it.DiscountValue); err != nil {
			return "Gagal", err
		}
		status := it.DiscountStatus
		if status == "" {
			status = "ACTIVE"
		}
		itemID := it.DiscountItemID
		itemName := it.DiscountItemName
		if scope == "TOTAL" {
			// sentinel "tidak spesifik item" -- lihat komentar Discount entity
			itemID = 0
			itemName = ""
		}
		row := discountEntity.Discount{
			DiscountGoldID:    goldid,
			DiscountOutcode:   it.DiscountOutcode,
			DiscountScope:     scope,
			DiscountItemID:    itemID,
			DiscountItemName:  itemName,
			DiscountType:      it.DiscountType,
			DiscountValue:     it.DiscountValue,
			DiscountStatus:    status,
			DiscountCreatedBy: actorName,
			DiscountCreatedAt: time.Now(),
		}
		if _, err := s.discount.InsertDiscountWithLog(ctx, row, actorName, role); err != nil {
			return "Gagal", errors.Wrap(err, "[Service][InsertDiscount]")
		}
	}
	return "Berhasil", nil
}

func (s Service) GetActiveTotalDiscount(ctx context.Context, goldid int, outcode string) (*discountEntity.Discount, error) {
	row, err := s.discount.GetActiveTotalDiscount(ctx, goldid, outcode)
	if err != nil {
		return nil, err // termasuk gorm.ErrRecordNotFound kalau tidak ada -- wajar, bukan error fatal
	}
	return row, nil
}

func (s Service) UpdateDiscount(ctx context.Context, goldid int, role string, req discountEntity.UpdateDiscount) (string, error) {
	if err := canManageDiscount(role); err != nil {
		return "Gagal", err
	}
	if err := validateDiscount(req.DiscountType, req.DiscountValue); err != nil {
		return "Gagal", err
	}
	existing, err := s.discount.GetDiscountByID(ctx, goldid, req.DiscountID)
	if err != nil {
		return "Gagal", errors.Wrap(err, "[Service][UpdateDiscount][GetDiscountByID]")
	}
	actorName, err := s.discount.GetGoldNameByID(ctx, goldid)
	if err != nil {
		actorName = ""
	}
	status := req.DiscountStatus
	if status == "" {
		status = "ACTIVE"
	}
	if err := s.discount.UpdateDiscountWithLog(ctx, *existing, req.DiscountType, req.DiscountValue, status, actorName, role); err != nil {
		return "Gagal", errors.Wrap(err, "[Service][UpdateDiscount]")
	}
	return "Berhasil", nil
}

func (s Service) DeleteDiscount(ctx context.Context, goldid int, role string, discountID int) (string, error) {
	if err := canManageDiscount(role); err != nil {
		return "Gagal", err
	}
	existing, err := s.discount.GetDiscountByID(ctx, goldid, discountID)
	if err != nil {
		return "Gagal", errors.Wrap(err, "[Service][DeleteDiscount][GetDiscountByID]")
	}
	actorName, err := s.discount.GetGoldNameByID(ctx, goldid)
	if err != nil {
		actorName = ""
	}
	if err := s.discount.DeleteDiscountWithLog(ctx, *existing, actorName, role); err != nil {
		return "Gagal", errors.Wrap(err, "[Service][DeleteDiscount]")
	}
	return "Berhasil", nil
}

func (s Service) GetDiscountHistory(ctx context.Context, discountID int, page, length int) ([]discountEntity.DiscountHistory, discountEntity.MetadataPaginationDetail, error) {
	var metadata discountEntity.MetadataPaginationDetail
	rows, err := s.discount.GetDiscountHistory(ctx, discountID, page, length)
	if err != nil {
		return rows, metadata, errors.Wrap(err, "[Service][GetDiscountHistory]")
	}
	total, err := s.discount.GetTotalDiscountHistory(ctx, discountID)
	if err != nil {
		return rows, metadata, errors.Wrap(err, "[Service][GetDiscountHistory][GetTotalDiscountHistory]")
	}
	metadata = discountEntity.MetadataPaginationDetail{
		Page:      page,
		Limit:     length,
		TotalData: int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(length))),
	}
	return rows, metadata, nil
}

// ---------------------------------------------------------------------------
// Voucher

const voucherCodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomVoucherCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = voucherCodeChars[rand.Intn(len(voucherCodeChars))]
	}
	return string(b)
}

// generateUniqueVoucherCode coba beberapa kali sampai dapat kode yang belum
// dipakai di outlet ini (tabrakan sangat jarang -- 36^8 kombinasi).
func (s Service) generateUniqueVoucherCode(ctx context.Context, outcode string) (string, error) {
	for i := 0; i < 10; i++ {
		code := randomVoucherCode(8)
		exists, err := s.discount.VoucherCodeExists(ctx, outcode, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", errors.New("gagal membuat kode voucher unik, coba lagi")
}

// GenerateVoucherCode dipakai FE untuk tombol "Generate" -- kode ditampilkan
// dulu ke penjual (masih bisa diedit) sebelum benar-benar di-insert.
func (s Service) GenerateVoucherCode(ctx context.Context, outcode string) (string, error) {
	code, err := s.generateUniqueVoucherCode(ctx, outcode)
	return code, errors.Wrap(err, "[Service][GenerateVoucherCode]")
}

func (s Service) InsertVoucher(ctx context.Context, goldid int, role string, req discountEntity.InsertVoucher) (string, error) {
	if err := canManageDiscount(role); err != nil {
		return "Gagal", err
	}
	if req.VoucherPercent <= 0 || req.VoucherPercent > 100 {
		return "Gagal", errors.New("persentase voucher harus antara 1-100")
	}
	code := strings.ToUpper(strings.TrimSpace(req.VoucherCode))
	if code == "" {
		generated, err := s.generateUniqueVoucherCode(ctx, req.VoucherOutcode)
		if err != nil {
			return "Gagal", err
		}
		code = generated
	} else {
		exists, err := s.discount.VoucherCodeExists(ctx, req.VoucherOutcode, code)
		if err != nil {
			return "Gagal", errors.Wrap(err, "[Service][InsertVoucher][VoucherCodeExists]")
		}
		if exists {
			return "Gagal", errors.New("kode voucher sudah dipakai, pilih kode lain")
		}
	}
	actorName, err := s.discount.GetGoldNameByID(ctx, goldid)
	if err != nil {
		actorName = ""
	}
	voucher := discountEntity.Voucher{
		VoucherGoldID:    goldid,
		VoucherOutcode:   req.VoucherOutcode,
		VoucherCode:      code,
		VoucherPercent:   req.VoucherPercent,
		VoucherExpiredAt: req.VoucherExpired,
		VoucherCreatedBy: actorName,
		VoucherCreatedAt: time.Now(),
	}
	if err := s.discount.InsertVoucher(ctx, voucher); err != nil {
		return "Gagal", errors.Wrap(err, "[Service][InsertVoucher]")
	}
	return code, nil
}

func (s Service) GetVouchers(ctx context.Context, goldid int, outcode string, page, length int) ([]discountEntity.Voucher, discountEntity.MetadataPaginationDetail, error) {
	var metadata discountEntity.MetadataPaginationDetail
	rows, err := s.discount.GetVouchers(ctx, goldid, outcode, page, length)
	if err != nil {
		return rows, metadata, errors.Wrap(err, "[Service][GetVouchers]")
	}
	total, err := s.discount.GetTotalVouchers(ctx, goldid, outcode)
	if err != nil {
		return rows, metadata, errors.Wrap(err, "[Service][GetVouchers][GetTotalVouchers]")
	}
	metadata = discountEntity.MetadataPaginationDetail{
		Page: page, Limit: length, TotalData: int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(length))),
	}
	return rows, metadata, nil
}

func (s Service) DeleteVoucher(ctx context.Context, goldid int, role string, voucherID int) (string, error) {
	if err := canManageDiscount(role); err != nil {
		return "Gagal", err
	}
	existing, err := s.discount.GetVoucherByID(ctx, goldid, voucherID)
	if err != nil {
		return "Gagal", errors.Wrap(err, "[Service][DeleteVoucher][GetVoucherByID]")
	}
	actorName, err := s.discount.GetGoldNameByID(ctx, goldid)
	if err != nil {
		actorName = ""
	}
	if err := s.discount.DeleteVoucherWithLog(ctx, *existing, actorName, role); err != nil {
		return "Gagal", errors.Wrap(err, "[Service][DeleteVoucher]")
	}
	return "Berhasil", nil
}

// CheckVoucher = pratinjau voucher (tidak mengonsumsi) -- dipakai POS untuk
// tampilkan potongan sebelum kasir input jumlah bayar.
func (s Service) CheckVoucher(ctx context.Context, goldid int, outcode, code string) (*discountEntity.Voucher, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	v, err := s.discount.GetVoucherByCode(ctx, goldid, outcode, code)
	if err != nil {
		return nil, errors.New("kode voucher tidak ditemukan atau sudah tidak valid")
	}
	if v.VoucherExpiredAt != nil && v.VoucherExpiredAt.Before(time.Now()) {
		return nil, errors.New("kode voucher sudah kedaluwarsa")
	}
	return v, nil
}

// RedeemVoucher dipanggil dari modul sales (bukan lewat HTTP endpoint
// tersendiri) saat kasir insert sales dengan kode voucher terisi.
func (s Service) RedeemVoucher(ctx context.Context, goldid int, outcode, code, saleID, actorName, actorRole string) (float64, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	percent, err := s.discount.RedeemVoucherWithLog(ctx, goldid, outcode, code, saleID, actorName, actorRole)
	return percent, err
}

func (s Service) GetVoucherHistory(ctx context.Context, goldid int, outcode string, page, length int) ([]discountEntity.VoucherHistory, discountEntity.MetadataPaginationDetail, error) {
	var metadata discountEntity.MetadataPaginationDetail
	rows, err := s.discount.GetVoucherHistory(ctx, goldid, outcode, page, length)
	if err != nil {
		return rows, metadata, errors.Wrap(err, "[Service][GetVoucherHistory]")
	}
	total, err := s.discount.GetTotalVoucherHistory(ctx, goldid, outcode)
	if err != nil {
		return rows, metadata, errors.Wrap(err, "[Service][GetVoucherHistory][GetTotalVoucherHistory]")
	}
	metadata = discountEntity.MetadataPaginationDetail{
		Page: page, Limit: length, TotalData: int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(length))),
	}
	return rows, metadata, nil
}
