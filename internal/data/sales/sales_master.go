package goldgym

import (
	"context"
	"fmt"
	goldSaleEntity "gold-gym-be/internal/entity/sales"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const dbTimeout = 3 * time.Second
const dbTimeoutInsert = 5 * time.Second

func (d *Data) WithTransactionThSale(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn) //begin
}

// func (d *Data) GetThSale(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldSaleEntity.ThSale, error) {
// 	var (
// 		users []goldSaleEntity.ThSale
// 		err   error
// 	)
// 	offset := (page - 1) * length
// 	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
// 	defer cancel()
// 	if page == 0 && length == 0 {
// 		if name == "" {
// 			err = d.db.WithContext(ctx).Where("sale_gold_id = ? AND sale_outcode = ? AND (? = '' or sale_name like ?)", goldid, outcode, "", "").Find(&users).Error
// 			if err != nil {
// 				return nil, errors.Wrap(err, "[DATA] [GetThSale]")
// 			}
// 		}
// 		if name != "" {
// 			name = "%" + name + "%"
// 			err = d.db.WithContext(ctx).Where("sale_gold_id = ? AND sale_outcode = ? AND (? = '' or sale_name like ?)", goldid, outcode, name, name).Find(&users).Error
// 			if err != nil {
// 				return nil, errors.Wrap(err, "[DATA] [GetThSale]")
// 			}
// 		}

// 	} else {
// 		if name != "" {
// 			name = "%" + name + "%"
// 		}
// 		err = d.db.WithContext(ctx).Where("sale_gold_id = ? AND sale_outcode = ? AND (? = '' or sale_name like ?)", goldid, outcode, name, name).Limit(length).Offset(offset).Find(&users).Error
// 		if err != nil {
// 			return nil, errors.Wrap(err, "[DATA] [GetThSale]")
// 		}
// 	}
// 	return users, err
// }

// func (d *Data) GetTotalThSale(ctx context.Context, goldid int, name string, outcode string) (int64, error) {
// 	var (
// 		total int64
// 		err   error
// 	)
// 	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
// 	defer cancel()
// 	query := d.db.WithContext(ctx).
// 		Model(&goldSaleEntity.ThSale{}).
// 		Where("sale_gold_id = ? AND sale_outcode = ?", goldid, outcode)

// 	if name != "" {
// 		query = query.Where("sale_name LIKE ?", "%"+name+"%")
// 	}

// 	err = query.Debug().Count(&total).Error

// 	return total, err
// }

// GetThSale: custid > 0 berarti scoping milik pembeli (sale_cust_id),
// selain itu scoping tenant kasir (sale_gold_id + outcode).
func (d *Data) GetThSale(ctx context.Context, goldid int, custid int, name string, outcode string, page, length int) ([]goldSaleEntity.ThSale, error) {
	var (
		sales []goldSaleEntity.ThSale
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := d.db.WithContext(ctx).Order("sale_created_at desc")
	if custid > 0 {
		query = query.Where("sale_cust_id = ?", custid)
	} else {
		query = query.Where("sale_gold_id = ? AND sale_outcode = ?", goldid, outcode)
	}

	if name != "" {
		like := "%" + name + "%"
		query = query.Where("(sale_trancnum LIKE ? OR sale_salescustomer LIKE ? OR sale_salesperson LIKE ?)", like, like, like)
	}
	if page > 0 && length > 0 {
		query = query.Limit(length).Offset((page - 1) * length)
	}

	err = query.Find(&sales).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetThSale]")
	}
	return sales, nil
}

func (d *Data) GetTotalThSale(ctx context.Context, goldid int, custid int, name string, outcode string) (int64, error) {
	var (
		total int64
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := d.db.WithContext(ctx).Model(&goldSaleEntity.ThSale{})
	if custid > 0 {
		query = query.Where("sale_cust_id = ?", custid)
	} else {
		query = query.Where("sale_gold_id = ? AND sale_outcode = ?", goldid, outcode)
	}

	if name != "" {
		like := "%" + name + "%"
		query = query.Where("(sale_trancnum LIKE ? OR sale_salescustomer LIKE ? OR sale_salesperson LIKE ?)", like, like, like)
	}

	err = query.Count(&total).Error
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][GetTotalThSale]")
	}
	return total, nil
}

func (d *Data) GetThSaleByID(ctx context.Context, goldid int, custid int, saleid string) (*goldSaleEntity.ThSale, error) {
	var (
		sale goldSaleEntity.ThSale
		err  error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := d.db.WithContext(ctx).Where("sale_id = ?", saleid)
	if custid > 0 {
		query = query.Where("sale_cust_id = ?", custid)
	} else {
		query = query.Where("sale_gold_id = ?", goldid)
	}
	err = query.First(&sale).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[DATA][GetThSaleByID]")
	}
	return &sale, nil
}

func (d *Data) GetTdSaleBySaleID(ctx context.Context, saleid string) ([]goldSaleEntity.TdSale, error) {
	var (
		details []goldSaleEntity.TdSale
		err     error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err = d.db.WithContext(ctx).
		Where("sale_id = ?", saleid).
		Order("td_id asc").
		Find(&details).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetTdSaleBySaleID]")
	}
	return details, nil
}

// GetSaleCustomerToko mengambil nama toko pembeli terdaftar (gold_toko) untuk nota.
// Kosong jika pembeli tidak punya toko (daftar mandiri) atau custid tidak ditemukan.
func (d *Data) GetSaleCustomerToko(ctx context.Context, custid int) (string, error) {
	var toko string
	err := d.db.WithContext(ctx).
		Table("data_peserta").
		Select("COALESCE(gold_toko, '')").
		Where("gold_id = ?", custid).
		Limit(1).
		Scan(&toko).Error
	if err != nil {
		return "", errors.Wrap(err, "[DATA][GetSaleCustomerToko]")
	}
	return toko, nil
}

// GetTotalProofBytes menjumlahkan ukuran SEMUA foto bukti pembayaran yang
// pernah diupload (validasi kuota global 10 GB).
func (d *Data) GetTotalProofBytes(ctx context.Context) (int64, error) {
	var total int64
	err := d.db.WithContext(ctx).
		Model(&goldSaleEntity.PaymentProof{}).
		Select("COALESCE(SUM(proof_bytes), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][GetTotalProofBytes]")
	}
	return total, nil
}

// GetPaymentProofsBySale daftar bukti pembayaran milik satu nota.
func (d *Data) GetPaymentProofsBySale(ctx context.Context, saleID string) ([]goldSaleEntity.PaymentProof, error) {
	var proofs []goldSaleEntity.PaymentProof
	err := d.db.WithContext(ctx).
		Where("proof_sale_id = ?", saleID).
		Order("proof_uploaded_at asc").
		Find(&proofs).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetPaymentProofsBySale]")
	}
	return proofs, nil
}

// GetPaymentProofByID metadata satu bukti pembayaran (nil jika tidak ada).
func (d *Data) GetPaymentProofByID(ctx context.Context, proofID int) (*goldSaleEntity.PaymentProof, error) {
	var proof goldSaleEntity.PaymentProof
	err := d.db.WithContext(ctx).Where("proof_id = ?", proofID).First(&proof).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[DATA][GetPaymentProofByID]")
	}
	return &proof, nil
}

// InsertPaymentProof menyimpan metadata foto bukti pembayaran.
func (d *Data) InsertPaymentProof(ctx context.Context, proof goldSaleEntity.PaymentProof) (goldSaleEntity.PaymentProof, error) {
	if err := d.db.WithContext(ctx).Create(&proof).Error; err != nil {
		return proof, errors.Wrap(err, "[DATA][InsertPaymentProof]")
	}
	return proof, nil
}

// MarkBookingsPaid menandai booking terapi UNPAID menjadi PAID dengan sale_id
// nota POS gabungan (booking ikut dibayar bersama barang lain dalam satu nota).
func (d *Data) MarkBookingsPaid(ctx context.Context, bookingIDs []string, saleID string) (int64, error) {
	if len(bookingIDs) == 0 {
		return 0, nil
	}
	res := d.db.WithContext(ctx).
		Table("booking").
		Where("booking_id IN ? AND booking_status = ?", bookingIDs, "UNPAID").
		Updates(map[string]interface{}{
			"booking_status":  "PAID",
			"booking_sale_id": saleID,
		})
	if res.Error != nil {
		return 0, errors.Wrap(res.Error, "[DATA][MarkBookingsPaid]")
	}
	return res.RowsAffected, nil
}

func (d *Data) GetSaleOutletInfo(ctx context.Context, goldid int, outcode string) (*goldSaleEntity.SaleOutletInfo, error) {
	var (
		outlet goldSaleEntity.SaleOutletInfo
		err    error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err = d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet_code, outlet_name, outlet_address, outlet_type").
		Where("outlet_gold_id = ? AND outlet_code = ?", goldid, outcode).
		Limit(1).
		Scan(&outlet).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetSaleOutletInfo]")
	}
	return &outlet, nil
}

// IsPosCustomerOptional true jika outlet (gold_id+outcode) boleh transaksi POS
// tanpa mengisi nama customer. Disimpan di kolom outlet.outlet_customer_optional
// ('Y' = boleh kosong, default; 'N' = wajib isi customer). Default DB 'Y' membuat
// outlet yang belum diatur admin otomatis opsional (boleh dikosongi).
func (d *Data) IsPosCustomerOptional(ctx context.Context, goldid int, outcode string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var flag string
	err := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet_customer_optional").
		Where("outlet_gold_id = ? AND outlet_code = ?", goldid, outcode).
		Limit(1).
		Scan(&flag).Error
	if err != nil {
		return false, errors.Wrap(err, "[DATA][IsPosCustomerOptional]")
	}
	// nilai kosong (outlet tak ditemukan / kolom NULL lawas) diperlakukan opsional
	return !strings.EqualFold(strings.TrimSpace(flag), "N"), nil
}

// GetOutletTypeByCode mengambil tipe outlet (RETAIL/THERAPY) berdasar pemilik+kode.
func (d *Data) GetOutletTypeByCode(ctx context.Context, goldid int, outcode string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var t string
	err := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet_type").
		Where("outlet_gold_id = ? AND outlet_code = ?", goldid, outcode).
		Limit(1).
		Scan(&t).Error
	if err != nil {
		return "", errors.Wrap(err, "[DATA][GetOutletTypeByCode]")
	}
	return t, nil
}

// GetPosOutletsForAdmin daftar outlet RETAIL aktif + alamat + pemilik + penanda
// optional (sudah diberi akses POS-tanpa-customer). search mencocokkan nama
// ATAU alamat outlet (untuk pencarian "pasar" dsb).
func (d *Data) GetPosOutletsForAdmin(ctx context.Context, search string) ([]goldSaleEntity.PosOutlet, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var outlets []goldSaleEntity.PosOutlet
	q := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet.outlet_gold_id, outlet.outlet_code, outlet.outlet_name, outlet.outlet_address, COALESCE(dp.gold_nama, '') AS owner_name, (COALESCE(outlet.outlet_customer_optional, 'Y') <> 'N') AS optional").
		Joins("LEFT JOIN data_peserta dp ON dp.gold_id = outlet.outlet_gold_id").
		Where("outlet.outlet_type = ? AND outlet.outlet_status = ?", "RETAIL", "ACTIVE")
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("outlet.outlet_name LIKE ? OR outlet.outlet_address LIKE ?", like, like)
	}
	if err := q.Order("outlet.outlet_name asc").Scan(&outlets).Error; err != nil {
		return nil, errors.Wrap(err, "[DATA][GetPosOutletsForAdmin]")
	}
	return outlets, nil
}

// SetPosCustomerOptional menyetel apakah outlet boleh POS tanpa customer:
// optional=true => 'Y' (boleh kosong), false => 'N' (wajib isi customer).
// addedBy dipertahankan di signature demi kompatibilitas interface (tidak
// disimpan lagi karena atribut kini melekat pada baris outlet).
func (d *Data) SetPosCustomerOptional(ctx context.Context, goldid int, outcode string, optional bool, addedBy string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	flag := "N"
	if optional {
		flag = "Y"
	}
	return d.db.WithContext(ctx).
		Table("outlet").
		Where("outlet_gold_id = ? AND outlet_code = ?", goldid, outcode).
		Update("outlet_customer_optional", flag).Error
}

// ================= Visibilitas fitur bukti pembayaran (admin) =================

// IsPaymentProofGloballyEnabled gerbang global (app_settings.payment_proof_enabled).
// Baris belum ada / nilai selain 'N' dianggap aktif (default aman: tidak
// mematikan fitur yang sudah berjalan sebelum migrasi ini dijalankan).
func (d *Data) IsPaymentProofGloballyEnabled(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var value string
	err := d.db.WithContext(ctx).
		Table("app_settings").
		Select("setting_value").
		Where("setting_key = ?", "payment_proof_enabled").
		Limit(1).
		Scan(&value).Error
	if err != nil {
		return false, errors.Wrap(err, "[DATA][IsPaymentProofGloballyEnabled]")
	}
	return !strings.EqualFold(strings.TrimSpace(value), "N"), nil
}

func (d *Data) SetPaymentProofGlobal(ctx context.Context, enabled bool, updatedBy string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	value := "N"
	if enabled {
		value = "Y"
	}
	return d.db.WithContext(ctx).
		Exec(`INSERT INTO app_settings (setting_key, setting_value, updated_by)
              VALUES (?, ?, ?)
              ON DUPLICATE KEY UPDATE setting_value = VALUES(setting_value), updated_by = VALUES(updated_by)`,
			"payment_proof_enabled", value, updatedBy).Error
}

// IsOutletProofEnabled gerbang per outlet (outlet.outlet_proof_enabled).
func (d *Data) IsOutletProofEnabled(ctx context.Context, goldid int, outcode string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var flag string
	err := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet_proof_enabled").
		Where("outlet_gold_id = ? AND outlet_code = ?", goldid, outcode).
		Limit(1).
		Scan(&flag).Error
	if err != nil {
		return false, errors.Wrap(err, "[DATA][IsOutletProofEnabled]")
	}
	// outlet tak ditemukan (mis. pembeli belum pilih outlet) dianggap aktif;
	// gerbang global/user tetap berlaku
	return !strings.EqualFold(strings.TrimSpace(flag), "N"), nil
}

// IsUserProofEnabled gerbang per akun (data_peserta.gold_proof_enabled).
func (d *Data) IsUserProofEnabled(ctx context.Context, goldid int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var flag string
	err := d.db.WithContext(ctx).
		Table("data_peserta").
		Select("gold_proof_enabled").
		Where("gold_id = ?", goldid).
		Limit(1).
		Scan(&flag).Error
	if err != nil {
		return false, errors.Wrap(err, "[DATA][IsUserProofEnabled]")
	}
	return !strings.EqualFold(strings.TrimSpace(flag), "N"), nil
}

// GetProofAccessOutlets daftar SEMUA outlet (RETAIL & THERAPY) + pemilik +
// status aktif/nonaktif fitur bukti pembayaran, untuk layar admin "per outlet".
func (d *Data) GetProofAccessOutlets(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessOutlet, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var outlets []goldSaleEntity.ProofAccessOutlet
	q := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet.outlet_gold_id, outlet.outlet_code, outlet.outlet_name, outlet.outlet_type, COALESCE(dp.gold_nama, '') AS owner_name, (COALESCE(outlet.outlet_proof_enabled, 'Y') <> 'N') AS enabled").
		Joins("LEFT JOIN data_peserta dp ON dp.gold_id = outlet.outlet_gold_id").
		Where("outlet.outlet_status = ?", "ACTIVE")
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("outlet.outlet_name LIKE ? OR outlet.outlet_address LIKE ?", like, like)
	}
	if err := q.Order("outlet.outlet_name asc").Scan(&outlets).Error; err != nil {
		return nil, errors.Wrap(err, "[DATA][GetProofAccessOutlets]")
	}
	return outlets, nil
}

func (d *Data) SetProofOutletEnabled(ctx context.Context, goldid int, outcode string, enabled bool) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	flag := "N"
	if enabled {
		flag = "Y"
	}
	return d.db.WithContext(ctx).
		Table("outlet").
		Where("outlet_gold_id = ? AND outlet_code = ?", goldid, outcode).
		Update("outlet_proof_enabled", flag).Error
}

// GetProofAccessUsers daftar akun (penjual retail/therapy & pembeli) + status
// aktif/nonaktif fitur bukti pembayaran, untuk layar admin "per user".
// search mencocokkan nama ATAU email.
func (d *Data) GetProofAccessUsers(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessUser, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	var users []goldSaleEntity.ProofAccessUser
	q := d.db.WithContext(ctx).
		Table("data_peserta").
		Select("gold_id, gold_nama, gold_email, gold_role, gold_buyer_yn, (COALESCE(gold_proof_enabled, 'Y') <> 'N') AS enabled").
		Where("gold_role IN ?", []string{"SELLER", "BUYER"})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("gold_nama LIKE ? OR gold_email LIKE ?", like, like)
	}
	if err := q.Order("gold_nama asc").Scan(&users).Error; err != nil {
		return nil, errors.Wrap(err, "[DATA][GetProofAccessUsers]")
	}
	return users, nil
}

func (d *Data) SetProofUserEnabled(ctx context.Context, goldid int, enabled bool) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	flag := "N"
	if enabled {
		flag = "Y"
	}
	return d.db.WithContext(ctx).
		Table("data_peserta").
		Where("gold_id = ?", goldid).
		Update("gold_proof_enabled", flag).Error
}

func (d *Data) GetLastThSaleCode(ctx context.Context, goldid int, code string) (*string, *string, error) {
	var (
		result   goldSaleEntity.ThSaleCounter
		saleid   *string
		trancnum *string
		err      error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Debug().Model(&goldSaleEntity.ThSale{}).Select("sale_id, sale_trancnum").Where("sale_gold_id = ? AND sale_outcode = ?", goldid, code).Order("sale_trancnum desc, sale_created_at desc").Limit(1).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// tidak ada data, tapi bukan error fatal
			return nil, nil, nil
		}
		return nil, nil, err
	}
	saleid = &result.SaleID
	trancnum = &result.SaleTrancnum
	return trancnum, saleid, err
}

func (d *Data) GetLastTdSaleCode(ctx context.Context, saleid string) (*string, error) {
	var (
		result string
		err    error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Debug().Model(&goldSaleEntity.TdSale{}).Select("sale_ref").Where("sale_id = ?", saleid).Order("sale_created_at desc").Limit(1).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (d *Data) InsertThSale(ctx context.Context, tx *gorm.DB, items []goldSaleEntity.ThSale) error {
	var (
		err error
	)
	if len(items) == 0 {
		return nil
	}

	if len(items) == 1 {
		err = tx.WithContext(ctx).Create(&items[0]).Error
	}
	if len(items) > 1 && len(items) <= 350 {
		err = tx.WithContext(ctx).CreateInBatches(items, 350).Error
	}

	if len(items) > 350 {
		batchSize := 500

		for i := 0; i < len(items); i += batchSize {
			end := i + batchSize
			if end > len(items) {
				end = len(items)
			}

			batch := items[i:end]

			valueStrings := []string{}
			valueArgs := []interface{}{}

			for _, v := range batch {

				valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)")
				valueArgs = append(valueArgs,
					v.SaleID,
					v.SaleGoldID,
					v.SaleOutcode,
					v.SaleTrancnum,
					v.SaleTranstime,
					v.SaleTransdate,
					v.SaleTranstotal,
					v.SaleTranspayment,
					v.SaleTranschange,
					v.SaleSalesperson,
					v.SaleSalescustomer,
					v.SalePaymentyn,
					v.SaleUpdatedAt,
				)
			}

			query := fmt.Sprintf(qInsertThSale, strings.Join(valueStrings, ","))

			if err := tx.WithContext(ctx).Debug().Exec(query, valueArgs...).Error; err != nil {
				return errors.Wrap(err, "[DATA][BulkInsertThSale]")
			}
		}
	}

	return err
}

// func (d *Data) UpdateThSale(ctx context.Context, cust goldSaleEntity.ThSale) error {

// 	updates := map[string]interface{}{}

// 	if cust.CustName != "" {
// 		updates["sale_name"] = cust.CustName
// 	}

// 	if cust.CustPhone != nil {
// 		updates["sale_phone"] = cust.CustPhone
// 	}

// 	if cust.CustAddress != nil {
// 		updates["sale_address"] = cust.CustAddress
// 	}

// 	if cust.CustEmail != nil {
// 		updates["sale_email"] = cust.CustEmail
// 	}

// 	if cust.CustStatus != "" {
// 		updates["sale_status"] = cust.CustStatus
// 	}
// 	if cust.CustOutletName != "" {
// 		updates["sale_outlet_name"] = cust.CustOutletName
// 	}

// 	if len(updates) == 0 {
// 		return nil
// 	}

// 	updates["sale_updated_at"] = time.Now()

// 	return d.db.WithContext(ctx).Debug().
// 		Model(&goldSaleEntity.ThSale{}).
// 		Where("sale_gold_id = ? AND sale_outcode = ? AND sale_id = ?",
// 			cust.CustGoldID,
// 			cust.CustOutcode,
// 			cust.CustID,
// 		).
// 		Updates(updates).Error
// }

// func (d *Data) DeleteThSale(ctx context.Context, goldid, goldcustomerid int, outcode string) error {
// 	return d.db.WithContext(ctx).Debug().Where("sale_gold_id = ? AND sale_outcode = ? AND sale_id = ?", goldid, outcode, goldcustomerid).Delete(&goldSaleEntity.ThSale{}).Error
// }

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

// func (d *Data) WithTransactionTdSale(ctx context.Context, fn func(tx *gorm.DB) error) error {
// 	return d.db.WithContext(ctx).Transaction(fn) //begin
// }

// func (d *Data) GetTdSale(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldSaleEntity.TdSale, error) {
// 	var (
// 		users []goldSaleEntity.TdSale
// 		err   error
// 	)
// 	offset := (page - 1) * length
// 	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
// 	defer cancel()
// 	if page == 0 && length == 0 {
// 		if name == "" {
// 			err = d.db.WithContext(ctx).Where("sale_gold_id = ? AND sale_outcode = ? AND (? = '' or sale_name like ?)", goldid, outcode, "", "").Find(&users).Error
// 			if err != nil {
// 				return nil, errors.Wrap(err, "[DATA] [GetTdSale]")
// 			}
// 		}
// 		if name != "" {
// 			name = "%" + name + "%"
// 			err = d.db.WithContext(ctx).Where("sale_gold_id = ? AND sale_outcode = ? AND (? = '' or sale_name like ?)", goldid, outcode, name, name).Find(&users).Error
// 			if err != nil {
// 				return nil, errors.Wrap(err, "[DATA] [GetTdSale]")
// 			}
// 		}

// 	} else {
// 		if name != "" {
// 			name = "%" + name + "%"
// 		}
// 		err = d.db.WithContext(ctx).Where("sale_gold_id = ? AND sale_outcode = ? AND (? = '' or sale_name like ?)", goldid, outcode, name, name).Limit(length).Offset(offset).Find(&users).Error
// 		if err != nil {
// 			return nil, errors.Wrap(err, "[DATA] [GetTdSale]")
// 		}
// 	}
// 	return users, err
// }

// func (d *Data) GetTotalTdSale(ctx context.Context, goldid int, name string, outcode string) (int64, error) {
// 	var (
// 		total int64
// 		err   error
// 	)
// 	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
// 	defer cancel()
// 	query := d.db.WithContext(ctx).
// 		Model(&goldSaleEntity.TdSale{}).
// 		Where("sale_gold_id = ? AND sale_outcode = ?", goldid, outcode)

// 	if name != "" {
// 		query = query.Where("sale_name LIKE ?", "%"+name+"%")
// 	}

// 	err = query.Debug().Count(&total).Error

// 	return total, err
// }

// func (d *Data) GetLastTdSaleCode(ctx context.Context, goldid int, code string) (*string, error) {
// 	var (
// 		result string
// 		err    error
// 	)
// 	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
// 	defer cancel()
// 	err = d.db.WithContext(ctx).Debug().Model(&goldSaleEntity.TdSale{}).Select("sale_trancnum").Where("sale_gold_id = ? AND sale_outcode = ?", goldid, code).Order("sale_trancnum desc, sale_created_at desc").Limit(1).Scan(&result).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &result, err
// }

func (d *Data) InsertTdSale(ctx context.Context, tx *gorm.DB, items []goldSaleEntity.TdSale) error {
	var (
		err error
	)
	if len(items) == 0 {
		return nil
	}

	if len(items) == 1 {
		err = tx.WithContext(ctx).Create(&items[0]).Error
	}
	if len(items) > 1 && len(items) <= 350 {
		err = tx.WithContext(ctx).CreateInBatches(items, 350).Error
	}

	if len(items) > 350 {
		batchSize := 500

		for i := 0; i < len(items); i += batchSize {
			end := i + batchSize
			if end > len(items) {
				end = len(items)
			}

			batch := items[i:end]

			valueStrings := []string{}
			valueArgs := []interface{}{}

			for _, v := range batch {

				valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, NOW())")
				valueArgs = append(valueArgs,
					v.SaleID,
					v.SaleStockid,
					v.SaleStockname,
					v.SaleQty,
					v.SaleSalesprice,
					v.SaleTotalsalesprice,
					v.SalePack,
					v.SaleRef,
				)
			}

			query := fmt.Sprintf(qInsertTdSale, strings.Join(valueStrings, ","))

			if err := tx.WithContext(ctx).Debug().Exec(query, valueArgs...).Error; err != nil {
				return errors.Wrap(err, "[DATA][BulkInsertTdSale]")
			}
		}
	}

	return err
}

// func (d *Data) UpdateTdSale(ctx context.Context, cust goldSaleEntity.TdSale) error {

// 	updates := map[string]interface{}{}

// 	if cust.CustName != "" {
// 		updates["sale_name"] = cust.CustName
// 	}

// 	if cust.CustPhone != nil {
// 		updates["sale_phone"] = cust.CustPhone
// 	}

// 	if cust.CustAddress != nil {
// 		updates["sale_address"] = cust.CustAddress
// 	}

// 	if cust.CustEmail != nil {
// 		updates["sale_email"] = cust.CustEmail
// 	}

// 	if cust.CustStatus != "" {
// 		updates["sale_status"] = cust.CustStatus
// 	}
// 	if cust.CustOutletName != "" {
// 		updates["sale_outlet_name"] = cust.CustOutletName
// 	}

// 	if len(updates) == 0 {
// 		return nil
// 	}

// 	updates["sale_updated_at"] = time.Now()

// 	return d.db.WithContext(ctx).Debug().
// 		Model(&goldSaleEntity.TdSale{}).
// 		Where("sale_gold_id = ? AND sale_outcode = ? AND sale_id = ?",
// 			cust.CustGoldID,
// 			cust.CustOutcode,
// 			cust.CustID,
// 		).
// 		Updates(updates).Error
// }

// func (d *Data) DeleteTdSale(ctx context.Context, goldid, goldcustomerid int, outcode string) error {
// 	return d.db.WithContext(ctx).Debug().Where("sale_gold_id = ? AND sale_outcode = ? AND sale_id = ?", goldid, outcode, goldcustomerid).Delete(&goldSaleEntity.TdSale{}).Error
// }

// UpdateSalePaymentYN menandai transaksi BELUM LUNAS menjadi LUNAS.
// Pembayaran dianggap pas: sale_transpayment = sale_transtotal, kembalian 0.
func (d *Data) UpdateSalePaymentYN(ctx context.Context, goldid int, saleid string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	res := d.db.WithContext(ctx).
		Model(&goldSaleEntity.ThSale{}).
		Where("sale_id = ? AND sale_gold_id = ? AND sale_paymentyn = ?", saleid, goldid, "N").
		Updates(map[string]interface{}{
			"sale_paymentyn":    "Y",
			"sale_transpayment": gorm.Expr("sale_transtotal"),
			"sale_transchange":  0,
		})
	if res.Error != nil {
		return 0, errors.Wrap(res.Error, "[DATA][UpdateSalePaymentYN]")
	}
	return res.RowsAffected, nil
}

// GetOutletGoldIDByCode mengambil gold_id tenant pemilik outlet — dipakai saat
// pembeli (BUYER) membuat transaksi supaya sale_gold_id tetap milik outlet.
func (d *Data) GetOutletGoldIDByCode(ctx context.Context, outcode string) (int, error) {
	var goldID int
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet_gold_id").
		Where("outlet_code = ?", outcode).
		Limit(1).
		Scan(&goldID).Error
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][GetOutletGoldIDByCode]")
	}
	if goldID == 0 {
		return 0, errors.New("outlet tidak ditemukan")
	}
	return goldID, nil
}
