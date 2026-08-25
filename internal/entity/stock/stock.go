package goldgym

import "time"

type GetOneStock struct {
	StockID         string `gorm:"column:stock_id" db:"stock_id" json:"stock_id"`
	StockGoldID     int    `gorm:"column:stock_gold_id" db:"stock_gold_id" json:"stock_gold_id"`
	StockOutletCode string `gorm:"column:stock_outcode" db:"stock_outcode" json:"stock_outcode"`
	StockItemID     int    `gorm:"column:stock_item_id" db:"stock_item_id" json:"stock_item_id"`
	StockName       string `gorm:"column:stock_name" db:"stock_name" json:"stock_name"`
	StockPack       string `gorm:"column:stock_pack" db:"stock_pack" json:"stock_pack"`
	StockQTY        int    `gorm:"column:stock_qty" db:"stock_qty" json:"stock_qty"`
	StocCreatedAt   string `gorm:"column:stock_created_at" db:"stock_created_at" json:"stock_created_at"`
	StocLastUpdate  string `gorm:"column:stock_last_update" db:"stock_last_update" json:"stock_last_update"`
	StocQtyUpdate   string `gorm:"column:stock_qty_update" db:"stock_qty_update" json:"stock_qty_update"`
	StockUpdateBy   string `gorm:"column:stock_update_by" db:"stock_update_by" json:"stock_update_by"`
	// harga diambil dari items.item_price via join (read-only, bukan kolom tabel stock)
	StockPrice int `gorm:"column:stock_price;->" db:"stock_price" json:"stock_price"`
	// merek diambil dari items.item_brand via join; brand THERAPY = jasa, stok tak dibatasi
	StockBrand string `gorm:"column:stock_brand;->" db:"stock_brand" json:"stock_brand"`
}

type InsertStock struct {
	StockID         string `gorm:"column:stock_id" db:"stock_id" json:"stock_id"`
	StockGoldID     int    `gorm:"column:stock_gold_id" db:"stock_gold_id" json:"stock_gold_id"`
	StockOutletCode string `gorm:"column:stock_outcode" db:"stock_outcode" json:"stock_outcode"`
	StockItemID     int    `gorm:"column:stock_item_id" db:"stock_item_id" json:"stock_item_id"`
	StockName       string `gorm:"column:stock_name" db:"stock_name" json:"stock_name"`
	StockPack       string `gorm:"column:stock_pack" db:"stock_pack" json:"stock_pack"`
	StockQTY        int    `gorm:"column:stock_qty" db:"stock_qty" json:"stock_qty"`
	StocCreatedAt   string `gorm:"column:stock_created_at" db:"stock_created_at" json:"stock_created_at"`
	StocLastUpdate  string `gorm:"column:stock_last_update" db:"stock_last_update" json:"stock_last_update"`
	StocQtyUpdate   string `gorm:"column:stock_qty_update" db:"stock_qty_update" json:"stock_qty_update"`
	StockUpdateBy   string `gorm:"column:stock_update_by" db:"stock_update_by" json:"stock_update_by"`
}

type MetadataPaginationDetail struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

type InsertStockData struct {
	StockData []InsertStock `json:"data"`
}

type StockHistory struct {
	ID                    string    `gorm:"column:sh_id" db:"sh_id" json:"sh_id"`
	StockHistoryStockID   string    `gorm:"column:sh_stock_id" db:"sh_stock_id" json:"sh_stock_id"`
	StockHistoryGoldID    int       `gorm:"column:sh_gold_id" db:"sh_gold_id" json:"sh_gold_id"`
	StockHistoryCode      string    `gorm:"column:sh_outcode" db:"sh_outcode" json:"sh_outcode"`
	StockHistoryItemID    int       `gorm:"column:sh_item_id" db:"sh_item_id" json:"sh_item_id"`
	StockHistoryName      string    `gorm:"column:sh_name" db:"sh_name" json:"sh_name"`
	StockHistoryPack      string    `gorm:"column:sh_pack" db:"sh_pack" json:"sh_pack"`
	StockHistoryStatus    string    `gorm:"column:sh_status" db:"sh_status" json:"sh_status"`
	StockHistoryQtyBefore int       `gorm:"column:sh_qty_before" db:"sh_qty_before" json:"sh_qty_before"`
	StockHistoryQtyChange int       `gorm:"column:sh_qty_change" db:"sh_qty_change" json:"sh_qty_change"`
	StockHistoryQtyAfter  int       `gorm:"column:sh_qty_after" db:"sh_qty_after" json:"sh_qty_after"`
	StockHistoryNote      string    `gorm:"column:sh_note" db:"sh_note" json:"sh_note"`
	StockHistoryCreatedAt time.Time `gorm:"column:sh_created_at" db:"sh_created_at" json:"sh_created_at"`
	StockHistoryCreatedBy string    `gorm:"column:sh_created_by" db:"sh_created_by" json:"sh_created_by"`
}

func (GetOneStock) TableName() string {
	return "stock"
}

func (InsertStock) TableName() string {
	return "stock"
}

func (StockHistory) TableName() string {
	return "stock_history"
}
