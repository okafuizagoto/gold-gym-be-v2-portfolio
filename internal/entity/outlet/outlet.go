package outlet

import "time"

type Outlet struct {
	OutletID        string     `gorm:"column:outlet_id" db:"outlet_id" json:"outlet_id"`
	OutletGoldID    int        `gorm:"column:outlet_gold_id" db:"outlet_gold_id" json:"outlet_gold_id"`
	OutletCode      string     `gorm:"column:outlet_code" db:"outlet_code" json:"outlet_code"`
	OutletName      string     `gorm:"column:outlet_name" db:"outlet_name" json:"outlet_name"`
	OutletType      string     `gorm:"column:outlet_type" db:"outlet_type" json:"outlet_type"`
	OutletAddress   string     `gorm:"column:outlet_address" db:"outlet_address" json:"outlet_address"`
	OutletStatus    string     `gorm:"column:outlet_status" db:"outlet_status" json:"outlet_status"`
	OutletCreatedAt time.Time  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	OutletUpdateAt  *time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

type InsertOutlet struct {
	OutletName    string `gorm:"column:outlet_name" db:"outlet_name" json:"outlet_name"`
	OutletAddress string `gorm:"column:outlet_address" db:"outlet_address" json:"outlet_address"`
	OutletStatus  string `gorm:"column:outlet_status" db:"outlet_status" json:"outlet_status"`
}

type UpdateOutlet struct {
	OutletGoldID  int    `gorm:"column:outlet_gold_id" db:"outlet_gold_id" json:"outlet_gold_id"`
	OutletCode    string `gorm:"column:outlet_code" db:"outlet_code" json:"outlet_code"`
	OutletAddress string `gorm:"column:outlet_address" db:"outlet_address" json:"outlet_address"`
	OutletStatus  string `gorm:"column:outlet_status" db:"outlet_status" json:"outlet_status"`
}

type InsertOutletData struct {
	InsertData []InsertOutlet `json:"data"`
}

type UpdateOutletData struct {
	UpdateData UpdateOutlet `json:"data"`
}

type OutletCounter struct {
	Prefix            string `gorm:"column:prefix" db:"prefix" json:"prefix"`
	CounterGoldID     int    `gorm:"column:counter_gold_id" db:"counter_gold_id" json:"counter_gold_id"`
	CounterOutletName string `gorm:"column:counter_outlet_name" db:"counter_outlet_name" json:"counter_outlet_name"`
	Counter           int    `gorm:"column:counter" db:"counter" json:"counter"`
}

type InsertOutletCounterJson struct {
	Prefix          string     `json:"prefix"`
	CounterGoldID   int        `json:"counter_gold_id"`
	OutletCode      string     `json:"outlet_code"`
	OutletName      string     `json:"outlet_name"`
	OutletAddress   string     `json:"outlet_address"`
	OutletStatus    string     `gorm:"column:outlet_status" db:"outlet_status" json:"outlet_status"`
	OutletCreatedAt time.Time  `json:"created_at"`
	OutletUpdateAt  *time.Time `json:"updated_at"`
}

type InsertOutletCounter struct {
	Prefix            string `gorm:"column:prefix" db:"prefix" json:"prefix"`
	CounterGoldID     int    `gorm:"column:counter_gold_id" db:"counter_gold_id" json:"counter_gold_id"`
	CounterOutletName string `gorm:"column:counter_outlet_name" db:"counter_outlet_name" json:"counter_outlet_name"`
}

type OutletDeleteHistory struct {
	HistoryID       int        `gorm:"column:history_id" db:"history_id" json:"history_id"`
	OutletGoldID    int        `gorm:"column:outlet_gold_id" db:"outlet_gold_id" json:"outlet_gold_id"`
	OutletID        string     `gorm:"column:outlet_id" db:"outlet_id" json:"outlet_id"`
	OutletCode      string     `gorm:"column:outlet_code" db:"outlet_code" json:"outlet_code"`
	OutletName      string     `gorm:"column:outlet_name" db:"outlet_name" json:"outlet_name"`
	OutletType      string     `gorm:"column:outlet_type" db:"outlet_type" json:"outlet_type"`
	OutletAddress   string     `gorm:"column:outlet_address" db:"outlet_address" json:"outlet_address"`
	OutletStatus    string     `gorm:"column:outlet_status" db:"outlet_status" json:"outlet_status"`
	OutletCreatedAt *time.Time `gorm:"column:outlet_created_at" db:"outlet_created_at" json:"outlet_created_at"`
	DeletedBy       int        `gorm:"column:deleted_by" db:"deleted_by" json:"deleted_by"`
	DeletedAt       time.Time  `gorm:"column:deleted_at" db:"deleted_at" json:"deleted_at"`
}

func (Outlet) TableName() string {
	return "outlet"
}

func (InsertOutlet) TableName() string {
	return "outlet"
}

func (UpdateOutlet) TableName() string {
	return "outlet"
}

func (OutletCounter) TableName() string {
	return "outlet_counter"
}

func (InsertOutletCounter) TableName() string {
	return "outlet_counter"
}

func (OutletDeleteHistory) TableName() string {
	return "outlet_delete_history"
}

type MetadataPaginationDetail struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}
