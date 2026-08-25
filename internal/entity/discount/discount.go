package goldgym

import "time"

// Discount = satu diskon aktif -- scope 'ITEM' (per produk, discount_item_id
// terisi) atau 'TOTAL' (persen dari total penjualan, auto-aktif per outlet;
// discount_item_id=0/discount_item_name="" jadi sentinel "tidak spesifik item").
type Discount struct {
	DiscountID        int        `gorm:"column:discount_id;primaryKey" db:"discount_id" json:"discount_id"`
	DiscountGoldID    int        `gorm:"column:discount_gold_id" db:"discount_gold_id" json:"discount_gold_id"`
	DiscountOutcode   string     `gorm:"column:discount_outcode" db:"discount_outcode" json:"discount_outcode"`
	DiscountScope     string     `gorm:"column:discount_scope" db:"discount_scope" json:"discount_scope"` // ITEM | TOTAL
	DiscountItemID    int        `gorm:"column:discount_item_id" db:"discount_item_id" json:"discount_item_id"`
	DiscountItemName  string     `gorm:"column:discount_item_name" db:"discount_item_name" json:"discount_item_name"`
	DiscountType      string     `gorm:"column:discount_type" db:"discount_type" json:"discount_type"` // PERCENT | NOMINAL
	DiscountValue     float64    `gorm:"column:discount_value" db:"discount_value" json:"discount_value"`
	DiscountStatus    string     `gorm:"column:discount_status" db:"discount_status" json:"discount_status"` // ACTIVE | NONACTIVE
	DiscountCreatedBy string     `gorm:"column:discount_created_by" db:"discount_created_by" json:"discount_created_by"`
	DiscountCreatedAt time.Time  `gorm:"column:discount_created_at" db:"discount_created_at" json:"discount_created_at"`
	DiscountUpdatedAt *time.Time `gorm:"column:discount_updated_at" db:"discount_updated_at" json:"discount_updated_at,omitempty"`
}

// InsertDiscount = payload satu baris insert diskon dari FE.
type InsertDiscount struct {
	DiscountOutcode  string  `json:"discount_outcode"`
	DiscountScope    string  `json:"discount_scope"` // ITEM | TOTAL, kosong = ITEM
	DiscountItemID   int     `json:"discount_item_id"`
	DiscountItemName string  `json:"discount_item_name"`
	DiscountType     string  `json:"discount_type"`
	DiscountValue    float64 `json:"discount_value"`
	DiscountStatus   string  `json:"discount_status"`
}

type InsertDiscountData struct {
	Data []InsertDiscount `json:"data"`
}

// UpdateDiscount = payload update satu diskon (identitas via DiscountID).
type UpdateDiscount struct {
	DiscountID     int     `json:"discount_id"`
	DiscountType   string  `json:"discount_type"`
	DiscountValue  float64 `json:"discount_value"`
	DiscountStatus string  `json:"discount_status"`
}

type UpdateDiscountData struct {
	Data UpdateDiscount `json:"data"`
}

// DiscountHistory = satu baris jejak audit (INSERT/UPDATE/DELETE) diskon.
type DiscountHistory struct {
	HistoryID                int        `gorm:"column:history_id;primaryKey" db:"history_id" json:"history_id"`
	HistoryDiscountID        int        `gorm:"column:history_discount_id" db:"history_discount_id" json:"history_discount_id"`
	HistoryAction            string     `gorm:"column:history_action" db:"history_action" json:"history_action"`
	HistoryGoldID            int        `gorm:"column:history_gold_id" db:"history_gold_id" json:"history_gold_id"`
	HistoryActorName         string     `gorm:"column:history_actor_name" db:"history_actor_name" json:"history_actor_name"`
	HistoryActorRole         string     `gorm:"column:history_actor_role" db:"history_actor_role" json:"history_actor_role"`
	HistoryOutcode           string     `gorm:"column:history_outcode" db:"history_outcode" json:"history_outcode"`
	HistoryItemID            int        `gorm:"column:history_item_id" db:"history_item_id" json:"history_item_id"`
	HistoryItemName          string     `gorm:"column:history_item_name" db:"history_item_name" json:"history_item_name"`
	HistoryDiscountType      string     `gorm:"column:history_discount_type" db:"history_discount_type" json:"history_discount_type"`
	HistoryDiscountValue     float64    `gorm:"column:history_discount_value" db:"history_discount_value" json:"history_discount_value"`
	HistoryDiscountStatus    string     `gorm:"column:history_discount_status" db:"history_discount_status" json:"history_discount_status"`
	HistoryDiscountCreatedAt *time.Time `gorm:"column:history_discount_created_at" db:"history_discount_created_at" json:"history_discount_created_at,omitempty"`
	HistoryChangedAt         time.Time  `gorm:"column:history_changed_at" db:"history_changed_at" json:"history_changed_at"`
}

// Voucher = kode voucher sekali pakai, diskon persen dari TOTAL keranjang.
// Dikonsumsi (dihapus dari tabel ini) saat dipakai di POS -- lihat
// discount_voucher_history untuk jejaknya.
type Voucher struct {
	VoucherID        int        `gorm:"column:voucher_id;primaryKey" db:"voucher_id" json:"voucher_id"`
	VoucherGoldID    int        `gorm:"column:voucher_gold_id" db:"voucher_gold_id" json:"voucher_gold_id"`
	VoucherOutcode   string     `gorm:"column:voucher_outcode" db:"voucher_outcode" json:"voucher_outcode"`
	VoucherCode      string     `gorm:"column:voucher_code" db:"voucher_code" json:"voucher_code"`
	VoucherPercent   float64    `gorm:"column:voucher_percent" db:"voucher_percent" json:"voucher_percent"`
	VoucherExpiredAt *time.Time `gorm:"column:voucher_expired_at" db:"voucher_expired_at" json:"voucher_expired_at,omitempty"`
	VoucherCreatedBy string     `gorm:"column:voucher_created_by" db:"voucher_created_by" json:"voucher_created_by"`
	VoucherCreatedAt time.Time  `gorm:"column:voucher_created_at" db:"voucher_created_at" json:"voucher_created_at"`
}

type InsertVoucher struct {
	VoucherOutcode string     `json:"voucher_outcode"`
	VoucherCode    string     `json:"voucher_code"` // kosong = auto-generate
	VoucherPercent float64    `json:"voucher_percent"`
	VoucherExpired *time.Time `json:"voucher_expired_at,omitempty"`
}

// VoucherHistory = jejak audit voucher: dibuat (tidak dicatat di sini, cukup
// tabel voucher itu sendiri), dipakai (USED, terisi history_sale_id), atau
// dihapus manual oleh penjual (DELETED) sebelum sempat dipakai.
type VoucherHistory struct {
	HistoryID               int        `gorm:"column:history_id;primaryKey" db:"history_id" json:"history_id"`
	HistoryVoucherCode      string     `gorm:"column:history_voucher_code" db:"history_voucher_code" json:"history_voucher_code"`
	HistoryOutcode          string     `gorm:"column:history_outcode" db:"history_outcode" json:"history_outcode"`
	HistoryGoldID           int        `gorm:"column:history_gold_id" db:"history_gold_id" json:"history_gold_id"`
	HistoryPercent          float64    `gorm:"column:history_percent" db:"history_percent" json:"history_percent"`
	HistoryStatus           string     `gorm:"column:history_status" db:"history_status" json:"history_status"` // USED | EXPIRED | DELETED
	HistorySaleID           *string    `gorm:"column:history_sale_id" db:"history_sale_id" json:"history_sale_id,omitempty"`
	HistoryActorName        string     `gorm:"column:history_actor_name" db:"history_actor_name" json:"history_actor_name"`
	HistoryActorRole        string     `gorm:"column:history_actor_role" db:"history_actor_role" json:"history_actor_role"`
	HistoryVoucherCreatedAt *time.Time `gorm:"column:history_voucher_created_at" db:"history_voucher_created_at" json:"history_voucher_created_at,omitempty"`
	HistoryChangedAt        time.Time  `gorm:"column:history_changed_at" db:"history_changed_at" json:"history_changed_at"`
}

// ItemForOutlet = sumber item-picker saat menambah diskon (item yang sudah
// ada di outlet ini, dari tabel items -- bukan tabel discount).
type ItemForOutlet struct {
	ItemID    int    `gorm:"column:item_id" db:"item_id" json:"item_id"`
	ItemName  string `gorm:"column:item_name" db:"item_name" json:"item_name"`
	ItemPrice int    `gorm:"column:item_price" db:"item_price" json:"item_price"`
}

type MetadataPaginationDetail struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func (Discount) TableName() string        { return "discount" }
func (DiscountHistory) TableName() string { return "discount_history" }
func (ItemForOutlet) TableName() string   { return "items" }
func (Voucher) TableName() string         { return "discount_voucher" }
func (VoucherHistory) TableName() string  { return "discount_voucher_history" }
