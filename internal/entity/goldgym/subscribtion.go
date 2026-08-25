package goldgym

import (
	"time"

	"gopkg.in/guregu/null.v3/zero"
)

type SubscriptionAll struct {
	GoldEmail           string    `json:"gold_email"`
	GoldId              int       `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldTotalharga      float64   `gorm:"column:gold_totalharga" db:"gold_totalharga" json:"gold_totalharga"`
	GoldValidasiPayment string    `gorm:"column:gold_validasipayment" db:"gold_validasipayment" json:"gold_validasipayment"`
	GoldOTP             string    `gorm:"column:gold_otp" db:"gold_otp" json:"gold_otp"`
	GoldLastupdate      time.Time `gorm:"column:gold_lastupdate" db:"gold_lastupdate" json:"gold_lastupdate"`
	// GoldMenuId int `gorm:"column:gold_menuid" json:"gold_menuid"`
	// GoldNamaPaket   string  `gorm:"column:gold_namapaket" json:"gold_namapaket"`
	// GoldNamaLayanan string  `gorm:"column:gold_namalayanan" json:"gold_namalayanan"`
	// GoldHarga       float64 `gorm:"column:gold_harga" json:"gold_harga"`
}

type SubscriptionHeader struct {
	GoldID              int         `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldTotalharga      zero.Float  `gorm:"column:gold_totalharga" db:"gold_totalharga" json:"gold_totalharga"`
	GoldValidasiPayment string      `gorm:"column:gold_validasipayment" db:"gold_validasipayment" json:"gold_validasipayment"`
	GoldOTP             zero.String `gorm:"column:gold_otp" db:"gold_otp" json:"gold_otp"`
	GoldLastupdate      zero.String `gorm:"column:gold_lastupdate" db:"gold_lastupdate" json:"gold_lastupdate"`
}

type SubscriptionHeaderPayment struct {
	// GoldID              int         `gorm:"column:gold_id" json:"gold_id"`
	GoldTotalharga      zero.Float `gorm:"column:gold_totalharga" db:"gold_totalharga" json:"gold_totalharga"`
	GoldValidasiPayment string     `gorm:"column:gold_validasipayment" db:"gold_validasipayment" json:"gold_validasipayment"`
	// GoldOTP             zero.String `gorm:"column:gold_otp" json:"gold_otp"`
	// GoldLastupdate      zero.String `gorm:"column:gold_lastupdate" json:"gold_lastupdate"`
}

type UpdatePayment struct {
	GoldID int `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
}

func (SubscriptionAll) TableName() string {
	return "subscription"
}

func (SubscriptionHeader) TableName() string {
	return "subscription"
}

func (SubscriptionHeaderPayment) TableName() string {
	return "subscription"
}
