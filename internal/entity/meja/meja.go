package meja

import "time"

type Meja struct {
	MejaID       int        `gorm:"column:meja_id;primaryKey" db:"meja_id" json:"meja_id"`
	MejaGoldID   int        `gorm:"column:meja_gold_id" db:"meja_gold_id" json:"meja_gold_id"`
	MejaOutcode  string     `gorm:"column:meja_outcode" db:"meja_outcode" json:"meja_outcode"`
	MejaAreaID   int        `gorm:"column:meja_area_id" db:"meja_area_id" json:"meja_area_id"`
	MejaName     string     `gorm:"column:meja_name" db:"meja_name" json:"meja_name"`
	MejaCapacity int        `gorm:"column:meja_capacity" db:"meja_capacity" json:"meja_capacity"`
	MejaStatus   string     `gorm:"column:meja_status" db:"meja_status" json:"meja_status"`
	MejaCreateAt time.Time  `gorm:"column:created_at" db:"created_at" json:"created_at"`
	MejaUpdateAt *time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at"`
}

func (Meja) TableName() string {
	return "meja"
}

const (
	MejaStatusKosong = "KOSONG"
	MejaStatusIsi    = "ISI"
)

type InsertMeja struct {
	MejaName     string `json:"meja_name"`
	MejaCapacity int    `json:"meja_capacity"`
	MejaAreaID   int    `json:"meja_area_id"`
}

type InsertMejaData struct {
	Outcode    string       `json:"outcode"`
	InsertData []InsertMeja `json:"data"`
}

type MejaStatusRequest struct {
	Outcode string `json:"outcode"`
	MejaIDs []int  `json:"meja_ids"`
}

// SaleMeja adalah snapshot meja yang dipakai pada satu nota -- lihat
// tabel sale_meja. Ditulis oleh modul sales (ConfirmSaleMeja), dibaca
// oleh modul sales juga saat generate PDF nota (GetSaleMejaNames).
type SaleMeja struct {
	ID        int       `gorm:"column:id;primaryKey" db:"id" json:"id"`
	SaleID    string    `gorm:"column:sale_id" db:"sale_id" json:"sale_id"`
	MejaID    int       `gorm:"column:meja_id" db:"meja_id" json:"meja_id"`
	MejaName  string    `gorm:"column:meja_name" db:"meja_name" json:"meja_name"`
	CreatedAt time.Time `gorm:"column:created_at" db:"created_at" json:"created_at"`
}

func (SaleMeja) TableName() string {
	return "sale_meja"
}
