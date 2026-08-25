package goldgym

import (
	"context"
	"fmt"
	goldSaleEntity "gold-gym-be/internal/entity/sales"
	goldStockEntity "gold-gym-be/internal/entity/stock"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// float64 → *decimal.Decimal
func decimalPtrFromFloat(f float64) *decimal.Decimal {
	d := decimal.NewFromFloat(f)
	return &d
}

// string → *decimal.Decimal (hati-hati parsing error)
func decimalPtrFromString(s string) (*decimal.Decimal, error) {
	if s == "" {
		return nil, nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// proofStorageDir direktori penyimpanan foto bukti pembayaran; bisa dioverride
// lewat env PHOTO_STORAGE_DIR (mis. saat backend jalan di pod dengan mount lain).
func proofStorageDir() string {
	if dir := os.Getenv("PHOTO_STORAGE_DIR"); dir != "" {
		return dir
	}
	return "/root/storages/photos"
}

// SavePaymentProof menyimpan foto bukti pembayaran transfer bank:
//   - ukuran per file maksimal 5 MB, hanya file gambar
//   - kuota global SUM(proof_bytes) maksimal 10 GB — jika penuh, upload ditolak
//     dengan pesan "tolong hubungi admin"
//   - file fisik ditulis ke proofStorageDir(), metadata ke tabel payment_proof
func (s Service) SavePaymentProof(ctx context.Context, saleID string, originalName string, mimeType string, content []byte, uploadedBy string) (goldSaleEntity.PaymentProof, error) {
	var proof goldSaleEntity.PaymentProof

	if saleID == "" {
		return proof, errors.New("saleid wajib diisi")
	}
	if len(content) == 0 {
		return proof, errors.New("file bukti pembayaran kosong")
	}
	if len(content) > goldSaleEntity.MaxProofFileBytes {
		return proof, errors.New("ukuran foto maksimal 5 MB")
	}
	if !strings.HasPrefix(strings.ToLower(mimeType), "image/") {
		return proof, errors.New("file harus berupa foto (jpg/png/webp)")
	}

	// validasi kuota global: total seluruh bytes di tabel + file baru
	total, err := s.goldgymsale.GetTotalProofBytes(ctx)
	if err != nil {
		return proof, errors.Wrap(err, "[Service][SavePaymentProof][GetTotalProofBytes]")
	}
	if total+int64(len(content)) > goldSaleEntity.MaxProofTotalBytes {
		return proof, errors.New("penyimpanan bukti pembayaran penuh — tolong hubungi admin")
	}

	dir := proofStorageDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return proof, errors.Wrap(err, "[Service][SavePaymentProof][MkdirAll]")
	}
	ext := strings.ToLower(filepath.Ext(originalName))
	if ext == "" {
		ext = ".jpg"
	}
	filename := uuid.New().String() + ext
	if err := os.WriteFile(filepath.Join(dir, filename), content, 0o644); err != nil {
		return proof, errors.Wrap(err, "[Service][SavePaymentProof][WriteFile]")
	}

	proof = goldSaleEntity.PaymentProof{
		ProofSaleID:       saleID,
		ProofFilename:     filename,
		ProofOriginalName: originalName,
		ProofMime:         mimeType,
		ProofBytes:        int64(len(content)),
		ProofPath:         dir,
		ProofUploadedBy:   uploadedBy,
	}
	saved, err := s.goldgymsale.InsertPaymentProof(ctx, proof)
	if err != nil {
		// metadata gagal disimpan: hapus file supaya tidak jadi sampah tak tercatat
		_ = os.Remove(filepath.Join(dir, filename))
		return proof, errors.Wrap(err, "[Service][SavePaymentProof][InsertPaymentProof]")
	}
	return saved, nil
}

// GetPaymentProofs daftar bukti pembayaran milik satu nota (untuk Sales History).
func (s Service) GetPaymentProofs(ctx context.Context, saleID string) ([]goldSaleEntity.PaymentProof, error) {
	if saleID == "" {
		return nil, errors.New("saleid wajib diisi")
	}
	proofs, err := s.goldgymsale.GetPaymentProofsBySale(ctx, saleID)
	if err != nil {
		return nil, errors.Wrap(err, "[Service][GetPaymentProofs]")
	}
	return proofs, nil
}

// GetPaymentProofFile membaca file foto bukti pembayaran dari disk untuk
// ditampilkan/di-download di aplikasi.
func (s Service) GetPaymentProofFile(ctx context.Context, proofID int) (goldSaleEntity.PaymentProof, []byte, error) {
	var proof goldSaleEntity.PaymentProof
	if proofID <= 0 {
		return proof, nil, errors.New("proofid wajib diisi")
	}
	p, err := s.goldgymsale.GetPaymentProofByID(ctx, proofID)
	if err != nil {
		return proof, nil, errors.Wrap(err, "[Service][GetPaymentProofFile]")
	}
	if p == nil {
		return proof, nil, errors.New("bukti pembayaran tidak ditemukan")
	}
	content, err := os.ReadFile(filepath.Join(p.ProofPath, p.ProofFilename))
	if err != nil {
		return *p, nil, errors.New("file bukti pembayaran tidak ditemukan di server")
	}
	return *p, content, nil
}

// MarkBookingsPaid menandai booking terapi UNPAID ikut lunas lewat nota POS
// gabungan (sale_id nota dipasang di booking supaya slot berubah merah).
func (s Service) MarkBookingsPaid(ctx context.Context, bookingIDs []string, saleID string) (int64, error) {
	if len(bookingIDs) == 0 || saleID == "" {
		return 0, nil
	}
	affected, err := s.goldgymsale.MarkBookingsPaid(ctx, bookingIDs, saleID)
	if err != nil {
		return 0, errors.Wrap(err, "[Service][MarkBookingsPaid]")
	}
	return affected, nil
}

func (s Service) InsertSales(ctx context.Context, goldid int, sales goldSaleEntity.InsertSales) (string, error) {
	var (
		result         string
		header         goldSaleEntity.ThSale
		headerArr      []goldSaleEntity.ThSale
		detail         goldSaleEntity.TdSale
		detailArr      []goldSaleEntity.TdSale
		lastSalesTdRef *string
		saleid         string
		code           string
		err            error
	)
	fmt.Println("MASOK SALES")
	if len(sales.Detail) == 0 {
		return "Data Kosong", nil
	} else {
		code = sales.Header.SaleOutcode
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	formatted := now.Format("0601021504")
	transtime := now.Format("150405")
	dateOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// waktu transaksi manual (khusus ADMIN — handler mengosongkan utk role lain);
	// hanya mempengaruhi sale_transdate & sale_transtime, nomor nota tetap live
	if sales.TransDate != "" {
		if d, errDate := time.ParseInLocation("2006-01-02", sales.TransDate, now.Location()); errDate == nil {
			dateOnly = d
		}
	}
	if sales.TransTime != "" {
		if t, errTime := time.Parse("15:04:05", sales.TransTime); errTime == nil {
			transtime = t.Format("150405")
		} else if t, errTime := time.Parse("15:04", sales.TransTime); errTime == nil {
			transtime = t.Format("150405")
		}
	}

	lastSalesCode, lastSalesID, err := s.goldgymsale.GetLastThSaleCode(ctx, goldid, code)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertSales][GetLastSalesCode]")
	}
	if lastSalesCode == nil || *lastSalesCode == "" {
		emptyCode := ""
		lastSalesCode = &emptyCode
	}
	if *lastSalesCode != "" {
		numberStr := (*lastSalesCode)[len(*lastSalesCode)-4:]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return result, errors.Wrap(err, "[Service][InsertSales][ParseNum][Last!Nil=0]")
		}
		number++
		stringNum := fmt.Sprintf("%04d", number)
		stringCompare := sales.Header.SaleOutcode + formatted + stringNum
		lastSalesCode = &stringCompare
	}
	if *lastSalesCode == "" {
		value := sales.Header.SaleOutcode + formatted + "0001"
		lastSalesCode = &value
	}
	if lastSalesID == nil || *lastSalesID == "" {
		empty := ""
		lastSalesTdRef = &empty
	} else {
		lastSalesTdRef, err = s.goldgymsale.GetLastTdSaleCode(ctx, *lastSalesID)
		if err != nil {
			result = "Gagal"
			return result, errors.Wrap(err, "[Service][InsertSales][GetLastSalesCode]")
		}
	}
	if *lastSalesTdRef != "" {
		numberStr := (*lastSalesTdRef)[len(*lastSalesTdRef)-3:]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return result, errors.Wrap(err, "[Service][InsertSales][ParseNum][Last!Nil=0]")
		}
		number++
		stringNum := fmt.Sprintf("%03d", number)
		// stringCompare := sales.Header.SaleOutcode+formatted + stringNum
		lastSalesTdRef = &stringNum
	}
	if *lastSalesTdRef == "" {
		value := strconv.Itoa(sales.Header.SaleGoldID) + sales.Header.SaleOutcode + formatted + "001"
		lastSalesTdRef = &value
	}
	// pakai sale_id dari handler (dibuat sebelum masuk Kafka) agar FE bisa langsung ambil nota
	saleid = sales.Header.SaleID
	if saleid == "" {
		saleid = uuid.New().String()
	}
	saleTotal, err := decimalPtrFromString(sales.Header.SaleTranstotal)
	if err != nil {
		return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString]")
	}
	saleTotalPayment, err := decimalPtrFromString(sales.Header.SaleTranspayment)
	if err != nil {
		return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString]")
	}
	saleTotalChange, err := decimalPtrFromString(sales.Header.SaleTranschange)
	if err != nil {
		return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString]")
	}
	paymentYN := "N"
	if sales.Header.SalePaymentyn != nil && *sales.Header.SalePaymentyn != "" {
		paymentYN = *sales.Header.SalePaymentyn
	}
	// sale_cust_id merujuk ke data_peserta.gold_id (akun BUYER) — biarkan NULL
	// untuk transaksi walk-in tanpa buyer login; 0 bukan gold_id valid dan
	// akan menabrak FK fk_thsale_buyer.
	custID := sales.Header.SaleCustID
	// toggle tampil customer di nota (khusus THERAPY); default "Y"
	customerShow := "Y"
	if sales.Header.SaleCustomerShow != nil && *sales.Header.SaleCustomerShow == "N" {
		customerShow = "N"
	}
	// Diskon TOTAL (client-trusted, sama pola dengan field harga lain) dan
	// voucher (SUDAH divalidasi+dihitung server-side di handler sebelum
	// masuk Kafka -- lihat insert_gold_gym_sales_gin.go) -- keduanya opsional.
	var totalDiscPercent, totalDiscAmount *decimal.Decimal
	if sales.Header.SaleTotalDiscountPercent != nil {
		totalDiscPercent, err = decimalPtrFromString(*sales.Header.SaleTotalDiscountPercent)
		if err != nil {
			return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString][totalDiscountPercent]")
		}
	}
	if sales.Header.SaleTotalDiscountAmount != nil {
		totalDiscAmount, err = decimalPtrFromString(*sales.Header.SaleTotalDiscountAmount)
		if err != nil {
			return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString][totalDiscountAmount]")
		}
	}
	var voucherPercent, voucherAmount *decimal.Decimal
	if sales.Header.SaleVoucherPercent != nil {
		voucherPercent, err = decimalPtrFromString(*sales.Header.SaleVoucherPercent)
		if err != nil {
			return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString][voucherPercent]")
		}
	}
	if sales.Header.SaleVoucherAmount != nil {
		voucherAmount, err = decimalPtrFromString(*sales.Header.SaleVoucherAmount)
		if err != nil {
			return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString][voucherAmount]")
		}
	}
	header = goldSaleEntity.ThSale{
		SaleID:                   saleid,
		SaleGoldID:               goldid,
		SaleOutcode:              code,
		SaleCustID:               custID,
		SaleTrancnum:             *lastSalesCode,
		SaleTranstime:            &transtime,
		SaleTransdate:            &dateOnly,
		SaleTranstotal:           saleTotal,
		SaleTranspayment:         saleTotalPayment,
		SaleTranschange:          saleTotalChange,
		SaleSalesperson:          sales.Header.SaleSalesperson,
		SaleSalescustomer:        sales.Header.SaleSalescustomer,
		SaleCustomerSource:       sales.Header.SaleCustomerSource,
		SaleCustomerShow:         customerShow,
		SalePaymentyn:            paymentYN,
		SaleCreatedAt:            now,
		SaleUpdatedAt:            nil,
		SaleTotalDiscountPercent: totalDiscPercent,
		SaleTotalDiscountAmount:  totalDiscAmount,
		SaleVoucherCode:          sales.Header.SaleVoucherCode,
		SaleVoucherPercent:       voucherPercent,
		SaleVoucherAmount:        voucherAmount,
	}
	headerArr = append(headerArr, header)

	for _, y := range sales.Detail {
		salePrice, err := decimalPtrFromString(y.SaleSalesprice)
		if err != nil {
			return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString]")
		}
		saleTotalPrice, err := decimalPtrFromString(y.SaleTotalsalesprice)
		if err != nil {
			return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString]")
		}
		// Snapshot diskon dikirim apa adanya dari klien (lihat komentar di
		// TDSaleDetail) -- bukan dihitung ulang di sini.
		var discValue, discAmount *decimal.Decimal
		if y.SaleDiscountValue != nil {
			discValue, err = decimalPtrFromString(*y.SaleDiscountValue)
			if err != nil {
				return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString][discountValue]")
			}
		}
		if y.SaleDiscountAmount != nil {
			discAmount, err = decimalPtrFromString(*y.SaleDiscountAmount)
			if err != nil {
				return result, errors.Wrap(err, "[Service][InsertSales][decimalPtrFromString][discountAmount]")
			}
		}
		detail = goldSaleEntity.TdSale{
			TdID:                  0,
			SaleID:                saleid,
			SaleStockid:           y.SaleStockid,
			SaleStockname:         y.SaleStockname,
			SaleQty:               y.SaleQty,
			SaleSalesprice:        salePrice,
			SaleTotalsalesprice:   saleTotalPrice,
			SalePack:              y.SalePack,
			SaleRef:               lastSalesTdRef,
			SaleCreatedAt:         now,
			SaleLastupdate:        &now,
			SaleDiscountID:        y.SaleDiscountID,
			SaleDiscountType:      y.SaleDiscountType,
			SaleDiscountValue:     discValue,
			SaleDiscountAmount:    discAmount,
			SaleDiscountCreatedAt: y.SaleDiscountCreatedAt,
		}
		detailArr = append(detailArr, detail)
	}
	err = s.goldgymsale.WithTransactionThSale(ctx, func(tx *gorm.DB) error {
		err = s.goldgymsale.InsertThSale(ctx, tx, headerArr)
		if err != nil {
			return errors.Wrap(err, "[Service][InsertSales][InsertThSale]")
		}
		err = s.goldgymsale.InsertTdSale(ctx, tx, detailArr)
		if err != nil {
			return errors.Wrap(err, "[Service][InsertSales][InsertTdSale]")
		}
		// Kurangi stok tiap item nota di transaksi yang sama supaya nota dan
		// pengurangan stok bersifat atomik (dulu stok tidak pernah berkurang).
		if err = s.deductStockOnSale(ctx, tx, header, detailArr); err != nil {
			return errors.Wrap(err, "[Service][InsertSales][deductStockOnSale]")
		}
		return nil
	})
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertSales][WithTransactionThSale]")
	}
	if result != "Gagal" {
		result = "Berhasil"
	}
	return result, err
}

// deductStockOnSale mengurangi stok untuk tiap item pada satu nota dan mencatat
// riwayatnya (stock_history status "MINUS"). Dijalankan di dalam transaksi yang
// sama dengan insert nota supaya keduanya atomik. Item jasa THERAPY
// (items.item_brand = "THERAPY") dilewati karena stoknya tidak dibatasi, begitu
// pula detail tanpa stock_id / qty <= 0 (mis. jasa manual).
//
// Catatan: karena insert nota berjalan async lewat Kafka worker, fungsi ini
// TIDAK menolak transaksi saat stok tidak cukup (qty_after bisa negatif) —
// menolak di sini hanya akan membuat nota hilang diam-diam tanpa kasir tahu.
// Pencegahan oversell sebaiknya dilakukan di layer HTTP sebelum publish.
func (s Service) deductStockOnSale(ctx context.Context, tx *gorm.DB, header goldSaleEntity.ThSale, details []goldSaleEntity.TdSale) error {
	// Agregasi qty per stock_id (item yang sama bisa muncul di beberapa baris).
	qtyByID := map[string]int{}
	orderedIDs := []string{}
	for _, d := range details {
		if d.SaleStockid == nil || *d.SaleStockid == "" {
			continue
		}
		if d.SaleQty == nil || *d.SaleQty <= 0 {
			continue
		}
		id := *d.SaleStockid
		if _, seen := qtyByID[id]; !seen {
			orderedIDs = append(orderedIDs, id)
		}
		qtyByID[id] += *d.SaleQty
	}
	if len(orderedIDs) == 0 {
		return nil
	}

	// Ambil kondisi stok terkini + brand (JOIN items) untuk semua stock_id.
	type stockRow struct {
		StockID string `gorm:"column:stock_id"`
		ItemID  int    `gorm:"column:stock_item_id"`
		Name    string `gorm:"column:stock_name"`
		Pack    string `gorm:"column:stock_pack"`
		Qty     int    `gorm:"column:stock_qty"`
		Brand   string `gorm:"column:stock_brand"`
	}
	var rows []stockRow
	if err := tx.WithContext(ctx).
		Table("stock").
		Select("stock.stock_id, stock.stock_item_id, stock.stock_name, stock.stock_pack, stock.stock_qty, COALESCE(items.item_brand, '') AS stock_brand").
		Joins("LEFT JOIN items ON items.item_id = stock.stock_item_id").
		Where("stock.stock_gold_id = ? AND stock.stock_id IN ?", header.SaleGoldID, orderedIDs).
		Scan(&rows).Error; err != nil {
		return errors.Wrap(err, "[Service][deductStockOnSale][loadStock]")
	}
	stockByID := make(map[string]stockRow, len(rows))
	for _, r := range rows {
		stockByID[r.StockID] = r
	}

	// Nomor sh_id berurutan; MAX diambil di dalam tx agar konsisten.
	var maxNo int64
	if err := tx.WithContext(ctx).
		Model(&goldStockEntity.StockHistory{}).
		Select("COALESCE(MAX(CAST(SUBSTRING(sh_id, 4) AS UNSIGNED)), 0)").
		Where("sh_id REGEXP '^STH[0-9]+$'").
		Scan(&maxNo).Error; err != nil {
		return errors.Wrap(err, "[Service][deductStockOnSale][lastHistoryCode]")
	}

	createdBy := ""
	if header.SaleSalesperson != nil {
		createdBy = *header.SaleSalesperson
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	var histories []goldStockEntity.StockHistory

	for _, id := range orderedIDs {
		row, ok := stockByID[id]
		if !ok {
			// stock_id tak punya row stok (mis. jasa) -> lewati.
			continue
		}
		if strings.EqualFold(row.Brand, "THERAPY") {
			// jasa terapi: stok tak dibatasi, jangan dikurangi.
			continue
		}
		qty := qtyByID[id]
		qtyBefore := row.Qty
		qtyAfter := qtyBefore - qty

		if err := tx.WithContext(ctx).
			Table("stock").
			Where("stock_gold_id = ? AND stock_id = ?", header.SaleGoldID, id).
			Updates(map[string]interface{}{
				"stock_qty":         gorm.Expr("stock_qty - ?", qty),
				"stock_qty_update":  gorm.Expr("NOW()"),
				"stock_last_update": gorm.Expr("NOW()"),
			}).Error; err != nil {
			return errors.Wrap(err, "[Service][deductStockOnSale][updateQty]")
		}

		maxNo++
		histories = append(histories, goldStockEntity.StockHistory{
			ID:                    fmt.Sprintf("STH%06d", maxNo),
			StockHistoryStockID:   id,
			StockHistoryGoldID:    header.SaleGoldID,
			StockHistoryCode:      header.SaleOutcode,
			StockHistoryItemID:    row.ItemID,
			StockHistoryName:      row.Name,
			StockHistoryPack:      row.Pack,
			StockHistoryStatus:    "MINUS",
			StockHistoryQtyBefore: qtyBefore,
			StockHistoryQtyChange: qty,
			StockHistoryQtyAfter:  qtyAfter,
			StockHistoryNote:      "STOCK UPDATE - SALES",
			StockHistoryCreatedAt: now,
			StockHistoryCreatedBy: createdBy,
		})
	}

	if len(histories) == 0 {
		return nil
	}
	if err := tx.WithContext(ctx).CreateInBatches(histories, 350).Error; err != nil {
		return errors.Wrap(err, "[Service][deductStockOnSale][insertHistory]")
	}
	return nil
}
