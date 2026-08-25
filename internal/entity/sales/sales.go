package goldgym

import (
	"time"

	"github.com/shopspring/decimal"
)

type SalesHeader struct {
	SaleID           string  `db:"sale_id" json:"sale_id"`
	SaleTransdate    string  `db:"sale_transdate" json:"sale_transdate"`
	SaleTransTime    string  `db:"sale_transtime" json:"sale_transtime"`
	SaleTranstotal   float64 `db:"sale_transtotal" json:"sale_transtotal"`
	SaleTranspayment float64 `db:"sale_transpayment" json:"sale_transpayment"`
	SaleTranschange  float64 `db:"sale_transchange" json:"sale_transchange"`
	SaleSalesperson  string  `db:"sale_salesperson" json:"sale_salesperson"`
}

type SalesDetail struct {
	SaleID         string `db:"sale_id" json:"sale_id"`
	SaleStockID    string `db:"sale_stockid" json:"sale_stockid"`
	SaleStockcode  string `db:"sale_stockcode" json:"sale_stockcode"`
	SaleStockname  string `db:"sale_stockname" json:"sale_stockname"`
	SaleQty        string `db:"sale_qty" json:"sale_qty"`
	SaleSalesprice string `db:"sale_salesprice" json:"sale_salesprice"`
	SalePack       string `db:"sale_pack" json:"sale_pack"`
	SaleLastupdate string `db:"sale_lastupdate" json:"sale_lastupdate"`
}

type ThSale struct {
	SaleID            string           `gorm:"column:sale_id;primaryKey" db:"sale_id" json:"sale_id"`
	SaleGoldID        int              `gorm:"column:sale_gold_id;not null" db:"sale_gold_id" json:"sale_gold_id"`
	SaleOutcode       string           `gorm:"column:sale_outcode;size:20;not null" db:"sale_outcode" json:"sale_outcode"`
	SaleCustID        *int             `gorm:"column:sale_cust_id" db:"sale_cust_id" json:"sale_cust_id"`
	SaleTrancnum      string           `gorm:"column:sale_trancnum;size:50;not null" db:"sale_trancnum" json:"sale_trancnum"`
	SaleTranstime     *string          `gorm:"column:sale_transtime;size:8" db:"sale_transtime" json:"sale_transtime,omitempty"`
	SaleTransdate     *time.Time       `gorm:"column:sale_transdate" db:"sale_transdate" json:"sale_transdate,omitempty"`
	SaleTranstotal    *decimal.Decimal `gorm:"column:sale_transtotal" db:"sale_transtotal" json:"sale_transtotal,omitempty"`
	SaleTranspayment  *decimal.Decimal `gorm:"column:sale_transpayment" db:"sale_transpayment" json:"sale_transpayment,omitempty"`
	SaleTranschange   *decimal.Decimal `gorm:"column:sale_transchange" db:"sale_transchange" json:"sale_transchange,omitempty"`
	SaleSalesperson   *string          `gorm:"column:sale_salesperson;size:15" db:"sale_salesperson" json:"sale_salesperson,omitempty"`
	SaleSalescustomer *string          `gorm:"column:sale_salescustomer;size:100" db:"sale_salescustomer" json:"sale_salescustomer,omitempty"`
	// sumber nama customer: PESERTA (gold_toko) / CUSTOMER (tabel customer) /
	// MANUAL (ketikan kasir) / BOOKING (nama booking terapi)
	SaleCustomerSource *string `gorm:"column:sale_customer_source;size:10" db:"sale_customer_source" json:"sale_customer_source,omitempty"`
	// Y/N: apakah nama customer dicetak di nota (toggle khusus outlet THERAPY)
	SaleCustomerShow string     `gorm:"column:sale_customer_show;size:1;default:Y" db:"sale_customer_show" json:"sale_customer_show"`
	SalePaymentyn    string     `gorm:"column:sale_paymentyn;size:1;default:N" db:"sale_paymentyn" json:"sale_paymentyn"`
	SaleCreatedAt    time.Time  `gorm:"column:sale_created_at" db:"sale_created_at" json:"sale_created_at"`
	SaleUpdatedAt    *time.Time `gorm:"column:sale_updated_at" db:"sale_updated_at" json:"sale_updated_at,omitempty"`
	// Diskon TOTAL (persen dari total belanja, auto-aktif per outlet) --
	// snapshot dari sisi klien saat checkout, sama filosofinya dengan
	// diskon per-item di td_sale (lihat SaleDiscountID cs).
	SaleTotalDiscountPercent *decimal.Decimal `gorm:"column:sale_total_discount_percent" db:"sale_total_discount_percent" json:"sale_total_discount_percent,omitempty"`
	SaleTotalDiscountAmount  *decimal.Decimal `gorm:"column:sale_total_discount_amount" db:"sale_total_discount_amount" json:"sale_total_discount_amount,omitempty"`
	// Voucher sekali pakai -- BEDA dari diskon total di atas: persen & nominal
	// di sini dihitung & divalidasi SERVER-SIDE (bukan dari klien) saat
	// RedeemVoucher dipanggil, karena voucher adalah resource langka yang
	// harus dikonsumsi atomik (lihat insert_gold_gym_sales_gin.go).
	SaleVoucherCode    *string          `gorm:"column:sale_voucher_code" db:"sale_voucher_code" json:"sale_voucher_code,omitempty"`
	SaleVoucherPercent *decimal.Decimal `gorm:"column:sale_voucher_percent" db:"sale_voucher_percent" json:"sale_voucher_percent,omitempty"`
	SaleVoucherAmount  *decimal.Decimal `gorm:"column:sale_voucher_amount" db:"sale_voucher_amount" json:"sale_voucher_amount,omitempty"`
}

type ThSaleCounter struct {
	SaleID       string `gorm:"column:sale_id;primaryKey" db:"sale_id" json:"sale_id"`
	SaleTrancnum string `gorm:"column:sale_trancnum;size:50;not null" db:"sale_trancnum" json:"sale_trancnum"`
}

type TdSale struct {
	TdID                int              `gorm:"column:td_id;primaryKey;autoIncrement" db:"td_id" json:"td_id"`
	SaleID              string           `gorm:"column:sale_id;not null" db:"sale_id" json:"sale_id"`
	SaleStockid         *string          `gorm:"column:sale_stockid" db:"sale_stockid" json:"sale_stockid,omitempty"`
	SaleStockname       *string          `gorm:"column:sale_stockname" db:"sale_stockname" json:"sale_stockname,omitempty"`
	SaleQty             *int             `gorm:"column:sale_qty" db:"sale_qty" json:"sale_qty,omitempty"`
	SaleSalesprice      *decimal.Decimal `gorm:"column:sale_salesprice" db:"sale_salesprice" json:"sale_salesprice,omitempty"`
	SaleTotalsalesprice *decimal.Decimal `gorm:"column:sale_totalsalesprice" db:"sale_totalsalesprice" json:"sale_totalsalesprice,omitempty"`
	SalePack            *string          `gorm:"column:sale_pack;size:10" db:"sale_pack" json:"sale_pack,omitempty"`
	SaleRef             *string          `gorm:"column:sale_ref;size:50;index:idx_td_sale_ref" db:"sale_ref" json:"sale_ref,omitempty"`
	SaleCreatedAt       time.Time        `gorm:"column:sale_created_at" db:"sale_created_at" json:"sale_created_at"`
	SaleLastupdate      *time.Time       `gorm:"column:sale_lastupdate" db:"sale_lastupdate" json:"sale_lastupdate,omitempty"`
	// Snapshot diskon yang aktif saat baris ini masuk nota (NULL = tanpa
	// diskon). Disalin apa adanya dari sisi klien (POS) saat checkout, BUKAN
	// dihitung ulang di sini -- lihat komentar di TDSaleDetail.
	SaleDiscountID        *int             `gorm:"column:sale_discount_id" db:"sale_discount_id" json:"sale_discount_id,omitempty"`
	SaleDiscountType      *string          `gorm:"column:sale_discount_type" db:"sale_discount_type" json:"sale_discount_type,omitempty"`
	SaleDiscountValue     *decimal.Decimal `gorm:"column:sale_discount_value" db:"sale_discount_value" json:"sale_discount_value,omitempty"`
	SaleDiscountAmount    *decimal.Decimal `gorm:"column:sale_discount_amount" db:"sale_discount_amount" json:"sale_discount_amount,omitempty"`
	SaleDiscountCreatedAt *time.Time       `gorm:"column:sale_discount_created_at" db:"sale_discount_created_at" json:"sale_discount_created_at,omitempty"`
}

type InsertSales struct {
	Header THSaleDetail   `json:"header"`
	Detail []TDSaleDetail `json:"detail"`
	// BookingIDs booking terapi UNPAID yang ikut dibayar lewat nota ini
	// (digabung barang lain jadi satu nota di POS) — handler menandai
	// booking tersebut PAID dengan sale_id nota ini. Diabaikan consumer Kafka.
	BookingIDs []string `json:"booking_ids,omitempty"`
	// TransDate ("2006-01-02") & TransTime ("15:04"/"15:04:05"): waktu transaksi
	// manual, KHUSUS role ADMIN (handler mengosongkan untuk role lain).
	// Kosong = waktu sekarang (live).
	TransDate string `json:"trans_date,omitempty"`
	TransTime string `json:"trans_time,omitempty"`
}

type InsertSaleData struct {
	InsertData InsertSales `json:"data"`
}

// PaymentProof metadata foto bukti pembayaran transfer bank (tabel
// payment_proof); file fisiknya disimpan backend di direktori
// PHOTO_STORAGE_DIR (default /root/storages/photos).
type PaymentProof struct {
	ProofID           int       `gorm:"column:proof_id;primaryKey;autoIncrement" json:"proof_id"`
	ProofSaleID       string    `gorm:"column:proof_sale_id" json:"proof_sale_id"`
	ProofFilename     string    `gorm:"column:proof_filename" json:"proof_filename"`
	ProofOriginalName string    `gorm:"column:proof_original_name" json:"proof_original_name"`
	ProofMime         string    `gorm:"column:proof_mime" json:"proof_mime"`
	ProofBytes        int64     `gorm:"column:proof_bytes" json:"proof_bytes"`
	ProofPath         string    `gorm:"column:proof_path" json:"proof_path"`
	ProofUploadedBy   string    `gorm:"column:proof_uploaded_by" json:"proof_uploaded_by"`
	ProofUploadedAt   time.Time `gorm:"column:proof_uploaded_at;autoCreateTime" json:"proof_uploaded_at"`
}

func (PaymentProof) TableName() string {
	return "payment_proof"
}

const (
	// MaxProofFileBytes ukuran maksimal satu foto bukti pembayaran (5 MB)
	MaxProofFileBytes = 5 * 1024 * 1024
	// MaxProofTotalBytes kuota total seluruh foto bukti pembayaran (10 GB) —
	// jika terlampaui, upload ditolak dengan pesan hubungi admin
	MaxProofTotalBytes = 10000000000
)

type THSaleDetail struct {
	SaleID            string     `gorm:"column:sale_id;primaryKey" db:"sale_id" json:"sale_id"`
	SaleGoldID        int        `gorm:"column:sale_gold_id;not null" db:"sale_gold_id" json:"sale_gold_id"`
	SaleOutcode       string     `gorm:"column:sale_outcode;size:20;not null" db:"sale_outcode" json:"sale_outcode"`
	SaleCustID        *int       `gorm:"column:sale_cust_id" db:"sale_cust_id" json:"sale_cust_id"`
	SaleTranstime     *string    `gorm:"column:sale_transtime;size:8" db:"sale_transtime" json:"sale_transtime,omitempty"`
	SaleTransdate     *time.Time `gorm:"column:sale_transdate" db:"sale_transdate" json:"sale_transdate,omitempty"`
	SaleTranstotal    string     `gorm:"column:sale_transtotal" db:"sale_transtotal" json:"sale_transtotal,omitempty"`
	SaleTranspayment  string     `gorm:"column:sale_transpayment" db:"sale_transpayment" json:"sale_transpayment,omitempty"`
	SaleTranschange   string     `gorm:"column:sale_transchange" db:"sale_transchange" json:"sale_transchange,omitempty"`
	SaleSalesperson   *string    `gorm:"column:sale_salesperson;size:15" db:"sale_salesperson" json:"sale_salesperson,omitempty"`
	SaleSalescustomer *string    `gorm:"column:sale_salescustomer;size:100" db:"sale_salescustomer" json:"sale_salescustomer,omitempty"`
	// sumber nama customer (PESERTA/CUSTOMER/MANUAL/BOOKING) — dari FE
	SaleCustomerSource *string `gorm:"column:sale_customer_source;size:10" db:"sale_customer_source" json:"sale_customer_source,omitempty"`
	// toggle tampil customer di nota (Y/N) — khusus THERAPY; default Y
	SaleCustomerShow *string `gorm:"column:sale_customer_show;size:1" db:"sale_customer_show" json:"sale_customer_show,omitempty"`
	SalePaymentyn    *string `gorm:"column:sale_paymentyn;size:1;default:N" db:"sale_paymentyn" json:"sale_paymentyn,omitempty"`
	// Diskon TOTAL (persen dari total belanja, auto-aktif per outlet) --
	// dikirim apa adanya dari klien (dihitung FE dari GetActiveDiscountsByOutlet
	// scope=TOTAL), sama trust boundary-nya dengan sale_transtotal dkk.
	SaleTotalDiscountPercent *string `json:"sale_total_discount_percent,omitempty"`
	SaleTotalDiscountAmount  *string `json:"sale_total_discount_amount,omitempty"`
	// Kode voucher yang diketik kasir di POS -- HANYA kode-nya yang dipercaya
	// dari klien; persen & nominal potongan SELALU dihitung ulang server-side
	// saat redeem (lihat RedeemVoucher), tidak pernah dari input klien.
	SaleVoucherCode *string `json:"sale_voucher_code,omitempty"`
	// Diisi HANDLER (bukan klien) setelah RedeemVoucher berhasil -- dua
	// field ini yang benar-benar disimpan ke th_sale, bukan input klien.
	SaleVoucherPercent *string `json:"sale_voucher_percent,omitempty"`
	SaleVoucherAmount  *string `json:"sale_voucher_amount,omitempty"`
}

type TDSaleDetail struct {
	SaleStockid         *string `gorm:"column:sale_stockid" db:"sale_stockid" json:"sale_stockid,omitempty"`
	SaleStockname       *string `gorm:"column:sale_stockname" db:"sale_stockname" json:"sale_stockname,omitempty"`
	SaleQty             *int    `gorm:"column:sale_qty" db:"sale_qty" json:"sale_qty,omitempty"`
	SaleSalesprice      string  `gorm:"column:sale_salesprice" db:"sale_salesprice" json:"sale_salesprice,omitempty"`
	SaleTotalsalesprice string  `gorm:"column:sale_totalsalesprice" db:"sale_totalsalesprice" json:"sale_totalsalesprice,omitempty"`
	SalePack            *string `gorm:"column:sale_pack;size:10" db:"sale_pack" json:"sale_pack,omitempty"`
	// Snapshot diskon dari sisi klien (POS) saat item dengan diskon aktif
	// masuk keranjang -- backend TIDAK query ulang tabel discount saat
	// insert (InsertSales diproses async via Kafka, bisa lama setelah
	// checkout; query ulang saat itu akan ambil kondisi diskon "sekarang",
	// bukan "saat transaksi terjadi"). Konsisten dengan field harga lain di
	// struct ini yang juga tidak dihitung ulang server-side.
	SaleDiscountID        *int       `json:"sale_discount_id,omitempty"`
	SaleDiscountType      *string    `json:"sale_discount_type,omitempty"`
	SaleDiscountValue     *string    `json:"sale_discount_value,omitempty"`
	SaleDiscountAmount    *string    `json:"sale_discount_amount,omitempty"`
	SaleDiscountCreatedAt *time.Time `json:"sale_discount_created_at,omitempty"`
}

type MetadataPaginationDetail struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

// SaleWithDetail dipakai untuk response detail transaksi & sumber data nota PDF
type SaleWithDetail struct {
	Header ThSale   `json:"header"`
	Detail []TdSale `json:"detail"`
}

// SaleOutletInfo untuk kop nota
type SaleOutletInfo struct {
	OutletCode    string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string `gorm:"column:outlet_name" json:"outlet_name"`
	OutletAddress string `gorm:"column:outlet_address" json:"outlet_address"`
	OutletType    string `gorm:"column:outlet_type" json:"outlet_type"`
}

// PosOutlet baris untuk layar admin "akses POS tanpa customer": outlet RETAIL +
// alamat + pemilik + penanda apakah boleh POS tanpa customer (Optional true =
// boleh kosong; ini default). Disimpan di kolom outlet.outlet_customer_optional.
type PosOutlet struct {
	OutletGoldID  int    `gorm:"column:outlet_gold_id" json:"outlet_gold_id"`
	OutletCode    string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string `gorm:"column:outlet_name" json:"outlet_name"`
	OutletAddress string `gorm:"column:outlet_address" json:"outlet_address"`
	OwnerName     string `gorm:"column:owner_name" json:"owner_name"`
	Optional      bool   `gorm:"column:optional" json:"optional"`
}

// PosAccessItem satu baris keputusan admin saat menyimpan akses (per outlet).
type PosAccessItem struct {
	GoldID   int    `json:"gold_id"`
	Outcode  string `json:"outcode"`
	Optional bool   `json:"optional"`
}

// SavePosAccessRequest body simpan akses POS-tanpa-customer (hanya baris yang
// sedang tampil di hasil pencarian admin — outlet lain tidak tersentuh).
type SavePosAccessRequest struct {
	Items []PosAccessItem `json:"items"`
}

// ================= Visibilitas fitur bukti pembayaran (admin) =================
// Tiga gerbang independen (global / per outlet / per user); fitur upload
// bukti transfer (POS & belanja pembeli) dan tombol lihat bukti (Sales
// History) hanya tampil jika ketiganya aktif ('Y').

// ProofAccessOutlet baris layar admin "per outlet": semua outlet (RETAIL &
// THERAPY, tidak seperti PosOutlet yang RETAIL saja) + penanda aktif/nonaktif.
type ProofAccessOutlet struct {
	OutletGoldID int    `gorm:"column:outlet_gold_id" json:"outlet_gold_id"`
	OutletCode   string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName   string `gorm:"column:outlet_name" json:"outlet_name"`
	OutletType   string `gorm:"column:outlet_type" json:"outlet_type"`
	OwnerName    string `gorm:"column:owner_name" json:"owner_name"`
	Enabled      bool   `gorm:"column:enabled" json:"enabled"`
}

// ProofAccessOutletItem satu baris keputusan admin saat menyimpan (per outlet).
type ProofAccessOutletItem struct {
	GoldID  int    `json:"gold_id"`
	Outcode string `json:"outcode"`
	Enabled bool   `json:"enabled"`
}

// ProofAccessUser baris layar admin "per user": penjual (retail/therapy)
// maupun pembeli + penanda aktif/nonaktif.
type ProofAccessUser struct {
	GoldID  int    `gorm:"column:gold_id" json:"gold_id"`
	Name    string `gorm:"column:gold_nama" json:"name"`
	Email   string `gorm:"column:gold_email" json:"email"`
	Role    string `gorm:"column:gold_role" json:"role"`
	BuyerYN string `gorm:"column:gold_buyer_yn" json:"buyer_yn"`
	Enabled bool   `gorm:"column:enabled" json:"enabled"`
}

// ProofAccessUserItem satu baris keputusan admin saat menyimpan (per user).
type ProofAccessUserItem struct {
	GoldID  int  `json:"gold_id"`
	Enabled bool `json:"enabled"`
}

func (ThSaleCounter) TableName() string {
	return "th_sale"
}

func (ThSale) TableName() string {
	return "th_sale"
}

func (TdSale) TableName() string {
	return "td_sale"
}

// ================= Laporan Penjualan (report) =================

// SaleReportItem satu baris item pada laporan per-hari. Beberapa baris dengan
// customer sama ditampilkan tanpa garis pemisah di kolom customer (FE).
type SaleReportItem struct {
	SaleID      string  `json:"sale_id"`
	Trancnum    string  `json:"trancnum"`
	Customer    string  `json:"customer"`
	Salesperson string  `json:"salesperson"`
	TransTime   string  `json:"trans_time"`
	ItemName    string  `json:"item_name"`
	Qty         int     `json:"qty"`
	Price       float64 `json:"price"`    // harga satuan
	Subtotal    float64 `json:"subtotal"` // total harga per item (qty x harga)
	Remaining   int     `json:"remaining"`
}

// SaleDailyTotal total penjualan satu hari (dipakai laporan per-minggu).
type SaleDailyTotal struct {
	Date  string  `json:"date"` // YYYY-MM-DD
	Total float64 `json:"total"`
	Count int     `json:"count"` // jumlah nota
}

// SaleWeeklyTotal total penjualan satu blok minggu (dipakai laporan per-bulan).
type SaleWeeklyTotal struct {
	WeekNo     int     `json:"week_no"`
	Label      string  `json:"label"` // mis. "1–7"
	RangeStart string  `json:"range_start"`
	RangeEnd   string  `json:"range_end"`
	Total      float64 `json:"total"`
	Count      int     `json:"count"`
}

// SaleReportDay respons laporan per-hari.
type SaleReportDay struct {
	Mode       string           `json:"mode"`
	Date       string           `json:"date"`
	Items      []SaleReportItem `json:"items"`
	Count      int              `json:"count"` // jumlah nota unik
	GrandTotal float64          `json:"grand_total"`
}

// SaleReportWeek respons laporan per-minggu (list total per hari).
type SaleReportWeek struct {
	Mode       string           `json:"mode"`
	Label      string           `json:"label"`
	RangeStart string           `json:"range_start"`
	RangeEnd   string           `json:"range_end"`
	Days       []SaleDailyTotal `json:"days"`
	GrandTotal float64          `json:"grand_total"`
}

// SaleReportMonth respons laporan per-bulan (list total per minggu).
type SaleReportMonth struct {
	Mode       string            `json:"mode"`
	Month      string            `json:"month"` // YYYY-MM
	Label      string            `json:"label"`
	Weeks      []SaleWeeklyTotal `json:"weeks"`
	GrandTotal float64           `json:"grand_total"`
}
