package goldgym

import "time"

type Item struct {
	ItemsID          int       `gorm:"column:item_id;primaryKey" db:"item_id" json:"item_id"`
	ItemsGoldID      int       `gorm:"column:item_gold_id" db:"item_gold_id" json:"item_gold_id"`
	ItemsCode        string    `gorm:"column:item_code" db:"item_code" json:"item_code"`
	ItemsName        string    `gorm:"column:item_name" db:"item_name" json:"item_name"`
	ItemsType        string    `gorm:"column:item_type" db:"item_type" json:"item_type"`
	ItemsPack        string    `gorm:"column:item_pack" db:"item_pack" json:"item_pack"`
	ItemsPrice       int       `gorm:"column:item_price" db:"item_price" json:"item_price"`
	ItemsBrand       string    `gorm:"column:item_brand" db:"item_brand" json:"item_brand"`
	ItemsDescription string    `gorm:"column:item_description" db:"item_description" json:"item_description"`
	ItemsStatus      string    `gorm:"column:item_status" db:"item_status" json:"item_status"`
	ItemsPhoto       string    `gorm:"column:item_photo" db:"item_photo" json:"item_photo"`
	ItemsPhotoBytes  int       `gorm:"column:item_photo_bytes" db:"item_photo_bytes" json:"item_photo_bytes"`
	ItemsCreatedAt   time.Time `gorm:"column:item_created_at" db:"item_created_at" json:"item_created_at"`
	ItemsUpdatedAt   time.Time `gorm:"column:item_updated_at" db:"item_updated_at" json:"item_updated_at"`
}

// InsertItemsResult respons InsertItems -- ItemID hanya terisi kalau insert
// tepat 1 item ke 1 outlet spesifik (bukan applyAllOutlets), karena hanya
// kasus itu yang punya satu item_id pasti untuk ditempeli foto sesudahnya.
type InsertItemsResult struct {
	Message string `json:"message"`
	ItemID  int    `json:"item_id,omitempty"`
}

type InsertItem struct {
	// ItemsID diisi balik oleh GORM (auto-increment) setelah Create -- perlu
	// primaryKey supaya GORM tahu ini kolom yang ditulis balik ke struct.
	ItemsID          int    `gorm:"column:item_id;primaryKey" json:"item_id,omitempty"`
	ItemsGoldID      int    `gorm:"column:item_gold_id" db:"item_gold_id" json:"item_gold_id"`
	ItemsOutletCode  string `gorm:"column:item_outcode" db:"item_outcode" json:"item_outcode"`
	ItemsCode        string `gorm:"column:item_code" db:"item_code" json:"item_code"`
	ItemsName        string `gorm:"column:item_name" db:"item_name" json:"item_name"`
	ItemsType        string `gorm:"column:item_type" db:"item_type" json:"item_type"`
	ItemsPack        string `gorm:"column:item_pack" db:"item_pack" json:"item_pack"`
	ItemsPrice       int    `gorm:"column:item_price" db:"item_price" json:"item_price"`
	ItemsBrand       string `gorm:"column:item_brand" db:"item_brand" json:"item_brand"`
	ItemsDescription string `gorm:"column:item_description" db:"item_description" json:"item_description"`
	ItemsStatus      string `gorm:"column:item_status" db:"item_status" json:"item_status"`
	ItemsPlace       string `gorm:"column:item_place" db:"item_place" json:"item_place"`
	ItemsEmail       string `gorm:"-" json:"item_email"`
}

type UpdateItems struct {
	ItemsName        string `gorm:"column:item_name" db:"item_name" json:"item_name"`
	ItemsType        string `gorm:"column:item_type" db:"item_type" json:"item_type"`
	ItemsPack        string `gorm:"column:item_pack" db:"item_pack" json:"item_pack"`
	ItemsPrice       int    `gorm:"column:item_price" db:"item_price" json:"item_price"`
	ItemsBrand       string `gorm:"column:item_brand" db:"item_brand" json:"item_brand"`
	ItemsDescription string `gorm:"column:item_description" db:"item_description" json:"item_description"`
	ItemsStatus      string `gorm:"column:item_status" db:"item_status" json:"item_status"`
	ItemsOutletCode  string `gorm:"column:item_outcode" db:"item_outcode" json:"item_outcode"`
	ItemsGoldID      int    `gorm:"-" json:"item_gold_id"`
	ItemsID          int    `gorm:"-" json:"item_id"`
}

type UpdateItemsData struct {
	UpdateData UpdateItems `json:"data"`
}

type MetadataPaginationDetail struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

type InsertItemData struct {
	ItemData        []InsertItem `json:"data"`
	ApplyAllOutlets bool         `json:"apply_all_outlets,omitempty"`
}

// MaxItemPhotoBytes ukuran maksimal satu foto item (2 MB)
const MaxItemPhotoBytes = 2 * 1024 * 1024

func (Item) TableName() string {
	return "items"
}

func (InsertItem) TableName() string {
	return "items"
}

func (UpdateItems) TableName() string {
	return "items"
}
