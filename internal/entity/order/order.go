package order

import (
	"time"

	"github.com/shopspring/decimal"
)

// Status pesanan pembeli
const (
	StatusWaiting = "WAITING" // menunggu konfirmasi penjual
	StatusProcess = "PROCESS" // penjual konfirmasi, sedang diproses
	StatusFinish  = "FINISH"  // selesai -> nota dibuat (masuk th_sale)
	StatusReject  = "REJECT"  // penjual menolak (dengan alasan)
)

// Tipe pembayaran
const (
	PayTunai    = "TUNAI"
	PayTransfer = "TRANSFER"
)

// Order header pesanan pembeli (tabel buyer_order)
type Order struct {
	OrderID           string           `gorm:"column:order_id;primaryKey" json:"order_id"`
	OrderBuyerID      int              `gorm:"column:order_buyer_id" json:"order_buyer_id"`
	OrderBuyerName    string           `gorm:"column:order_buyer_name" json:"order_buyer_name"`
	OrderGoldID       int              `gorm:"column:order_gold_id" json:"order_gold_id"`
	OrderOutcode      string           `gorm:"column:order_outcode" json:"order_outcode"`
	OrderOutletName   string           `gorm:"column:order_outlet_name" json:"order_outlet_name"`
	OrderTotal        decimal.Decimal  `gorm:"column:order_total" json:"order_total"`
	OrderPayType      string           `gorm:"column:order_pay_type" json:"order_pay_type"`
	OrderPaidYN       string           `gorm:"column:order_paid_yn" json:"order_paid_yn"`
	OrderStatus       string           `gorm:"column:order_status" json:"order_status"`
	OrderRejectReason *string          `gorm:"column:order_reject_reason" json:"order_reject_reason"`
	OrderSaleID       *string          `gorm:"column:order_sale_id" json:"order_sale_id"`
	OrderCreatedAt    time.Time        `gorm:"column:order_created_at" json:"order_created_at"`
	OrderUpdatedAt    *time.Time       `gorm:"column:order_updated_at" json:"order_updated_at"`
}

func (Order) TableName() string { return "buyer_order" }

// OrderDetail baris item pesanan (tabel buyer_order_detail)
type OrderDetail struct {
	OdID        int             `gorm:"column:od_id;primaryKey;autoIncrement" json:"od_id"`
	OdOrderID   string          `gorm:"column:od_order_id" json:"od_order_id"`
	OdStockID   string          `gorm:"column:od_stock_id" json:"od_stock_id"`
	OdStockName string          `gorm:"column:od_stock_name" json:"od_stock_name"`
	OdQty       int             `gorm:"column:od_qty" json:"od_qty"`
	OdPrice     decimal.Decimal `gorm:"column:od_price" json:"od_price"`
	OdTotal     decimal.Decimal `gorm:"column:od_total" json:"od_total"`
	OdPack      *string         `gorm:"column:od_pack" json:"od_pack"`
}

func (OrderDetail) TableName() string { return "buyer_order_detail" }

// OrderWithDetail response detail pesanan (header + item)
type OrderWithDetail struct {
	Header Order         `json:"header"`
	Detail []OrderDetail `json:"detail"`
}

// PublicOutlet outlet yang bisa dipilih pembeli (semua outlet RETAIL aktif,
// lintas penjual). GoldNama = nama pemilik/penjual (join data_peserta).
type PublicOutlet struct {
	OutletGoldID  int    `gorm:"column:outlet_gold_id" json:"outlet_gold_id"`
	OutletCode    string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName    string `gorm:"column:outlet_name" json:"outlet_name"`
	OutletType    string `gorm:"column:outlet_type" json:"outlet_type"`
	OutletAddress string `gorm:"column:outlet_address" json:"outlet_address"`
	OwnerName     string `gorm:"column:owner_name" json:"owner_name"`
	// Visible hanya diisi pada daftar admin: true jika outlet sudah dikurasi
	// (boleh dilihat pembeli). Pada daftar pembeli selalu true.
	Visible bool `gorm:"column:visible" json:"visible"`
}

// VisibleOutlet baris kurasi outlet-untuk-pembeli (tabel buyer_visible_outlet).
type VisibleOutlet struct {
	BvoID      int    `gorm:"column:bvo_id;primaryKey;autoIncrement" json:"bvo_id"`
	BvoGoldID  int    `gorm:"column:bvo_gold_id" json:"bvo_gold_id"`
	BvoOutcode string `gorm:"column:bvo_outcode" json:"bvo_outcode"`
	BvoAddedBy string `gorm:"column:bvo_added_by" json:"bvo_added_by"`
}

func (VisibleOutlet) TableName() string { return "buyer_visible_outlet" }

// CatalogItem barang yang bisa dipesan dari satu outlet (join stock+items)
type CatalogItem struct {
	StockID   string `gorm:"column:stock_id" json:"stock_id"`
	StockName string `gorm:"column:stock_name" json:"stock_name"`
	StockPack string `gorm:"column:stock_pack" json:"stock_pack"`
	StockQty  int    `gorm:"column:stock_qty" json:"stock_qty"`
	Price     int    `gorm:"column:price" json:"price"`
	Brand     string `gorm:"column:brand" json:"brand"`
}

// InsertOrderLine satu baris item dari FE saat pembeli checkout
type InsertOrderLine struct {
	StockID   string `json:"stock_id"`
	StockName string `json:"stock_name"`
	Qty       int    `json:"qty"`
	Price     int    `json:"price"`
	Pack      string `json:"pack"`
}

// InsertOrderRequest body checkout pembeli
type InsertOrderRequest struct {
	GoldID     int               `json:"gold_id"`     // pemilik outlet (dari daftar outlet)
	Outcode    string            `json:"outcode"`     // kode outlet
	OutletName string            `json:"outlet_name"` // nama outlet (tampilan)
	PayType    string            `json:"pay_type"`    // TUNAI / TRANSFER
	Lines      []InsertOrderLine `json:"lines"`
}

type InsertOrderData struct {
	InsertData InsertOrderRequest `json:"data"`
}
