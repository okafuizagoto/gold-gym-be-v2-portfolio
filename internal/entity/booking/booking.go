package booking

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	StatusPaid      = "PAID"
	StatusUnpaid    = "UNPAID"
	StatusExpired   = "EXPIRED"
	StatusCancelled = "CANCELLED"

	// Kapasitas maksimal per jendela 30 menit (total gabungan semua tipe terapi)
	CapacityPerHalfHour = 6

	// Jam operasional slot booking (05:00 - 17:00 WIB)
	OpenHour  = 5
	CloseHour = 17

	// Tipe terapi
	TherapyTypeSofa   = "SOFA"
	TherapyTypeDragon = "DRAGON"
	TherapyTypeKursi  = "KURSI"
)

// Jakarta adalah zona waktu tetap UTC+7 (WIB). Semua logika jam booking memakai
// zona ini supaya konsisten di server mana pun; DSN database juga loc=Asia/Jakarta.
var Jakarta = time.FixedZone("WIB", 7*60*60)

// TherapyItemName memetakan tipe terapi + durasi ke nama item jasa di tabel items
// (di-seed oleh migrasi 20260717_booking_therapy_v2.sql & 20260717_therapy_type_kursi.sql):
// SOFA 30 MENIT 15.000, SOFA 1 JAM 25.000,
// KURSI DRAGON 30 MENIT 25.000, KURSI DRAGON 1 JAM 35.000,
// KURSI 30 MENIT 10.000, KURSI 1 JAM 20.000.
func TherapyItemName(therapyType string, duration int) string {
	prefix := "SOFA"
	switch therapyType {
	case TherapyTypeDragon:
		prefix = "KURSI DRAGON"
	case TherapyTypeKursi:
		prefix = "KURSI"
	}
	if duration == 60 {
		return prefix + " 1 JAM"
	}
	return prefix + " 30 MENIT"
}

type Booking struct {
	BookingID           string           `gorm:"column:booking_id;primaryKey" db:"booking_id" json:"booking_id"`
	BookingGoldID       int              `gorm:"column:booking_gold_id" db:"booking_gold_id" json:"booking_gold_id"`
	BookingOutcode      string           `gorm:"column:booking_outcode" db:"booking_outcode" json:"booking_outcode"`
	BookingDate         time.Time        `gorm:"column:booking_date" db:"booking_date" json:"booking_date"`
	BookingStart        time.Time        `gorm:"column:booking_start" db:"booking_start" json:"booking_start"`
	BookingDuration     int              `gorm:"column:booking_duration" db:"booking_duration" json:"booking_duration"`
	BookingTherapyType  *string          `gorm:"column:booking_therapy_type" db:"booking_therapy_type" json:"booking_therapy_type"`
	BookingCustID       *int             `gorm:"column:booking_cust_id" db:"booking_cust_id" json:"booking_cust_id"`
	BookingCustName     string           `gorm:"column:booking_cust_name" db:"booking_cust_name" json:"booking_cust_name"`
	BookingRegisteredYN string           `gorm:"column:booking_registered_yn" db:"booking_registered_yn" json:"booking_registered_yn"`
	BookingStatus       string           `gorm:"column:booking_status" db:"booking_status" json:"booking_status"`
	BookingItemID       *int             `gorm:"column:booking_item_id" db:"booking_item_id" json:"booking_item_id"`
	BookingItemName     *string          `gorm:"column:booking_item_name" db:"booking_item_name" json:"booking_item_name"`
	BookingPrice        *decimal.Decimal `gorm:"column:booking_price" db:"booking_price" json:"booking_price"`
	BookingSaleID       *string          `gorm:"column:booking_sale_id" db:"booking_sale_id" json:"booking_sale_id"`
	BookingCreatedBy    *string          `gorm:"column:booking_created_by" db:"booking_created_by" json:"booking_created_by"`
	BookingCreatedRole  *string          `gorm:"column:booking_created_role" db:"booking_created_role" json:"booking_created_role"`
	BookingCreatedAt    time.Time        `gorm:"column:booking_created_at" db:"booking_created_at" json:"booking_created_at"`
	BookingUpdatedAt    *time.Time       `gorm:"column:booking_updated_at" db:"booking_updated_at" json:"booking_updated_at"`
}

func (Booking) TableName() string {
	return "booking"
}

// InsertBookingRequest body dari frontend
type InsertBookingRequest struct {
	Outcode  string `json:"outcode"`
	Date     string `json:"date"`  // "2006-01-02"
	Start    string `json:"start"` // "15:04" — harus kelipatan 30 menit
	Duration int    `json:"duration"`
	CustID   int    `json:"cust_id"`   // gold_id pembeli terdaftar; 0 = belum terdaftar
	CustName string `json:"cust_name"` // wajib jika cust_id 0
	Paid     bool   `json:"paid"`
	ItemID   int    `json:"item_id"` // legacy: item terapi eksplisit (tidak dipakai jika therapy_type diisi)
	PayType  string `json:"pay_type"`

	// TherapyType SOFA atau DRAGON — menentukan item & harga (per durasi)
	TherapyType string `json:"therapy_type"`
	// CustomPrice harga input manual (khusus role SELLER/ADMIN); 0 = harga default
	CustomPrice int `json:"custom_price"`
}

type InsertBookingData struct {
	InsertData InsertBookingRequest `json:"data"`
}

// SlotBooking ringkasan satu booking dalam slot (untuk tampilan grid FE)
type SlotBooking struct {
	BookingID    string `json:"booking_id"`
	CustName     string `json:"cust_name"`
	RegisteredYN string `json:"registered_yn"`
	Status       string `json:"status"`
	Duration     int    `json:"duration"`
	Start        string `json:"start"`
	TherapyType  string `json:"therapy_type"`
	// Price harga tersimpan di booking ("0" jika belum ada) — untuk booking UNPAID
	// bisa berisi harga custom yang dikunci penjual saat booking dibuat
	Price string `json:"price"`
}

// Slot satu jendela 30 menit
type Slot struct {
	Start     string        `json:"start"` // "08:00"
	End       string        `json:"end"`   // "08:30"
	Capacity  int           `json:"capacity"`
	Used      int           `json:"used"`
	Available int           `json:"available"`
	HasUnpaid bool          `json:"has_unpaid"`
	Full      bool          `json:"full"`
	Past      bool          `json:"past"`
	Bookings  []SlotBooking `json:"bookings"`
}

// BuyerInfo hasil pencarian pembeli terdaftar (autocomplete penjual)
type BuyerInfo struct {
	GoldId      int    `gorm:"column:gold_id" json:"gold_id"`
	GoldNama    string `gorm:"column:gold_nama" json:"gold_nama"`
	GoldToko    string `gorm:"column:gold_toko" json:"gold_toko"`
	GoldEmail   string `gorm:"column:gold_email" json:"gold_email"`
	GoldNomorHp string `gorm:"column:gold_nomorhp" json:"gold_nomorhp"`
}

func (BuyerInfo) TableName() string {
	return "data_peserta"
}

// TherapyItem item jasa terapi dari tabel items
type TherapyItem struct {
	ItemID    int    `gorm:"column:item_id" json:"item_id"`
	ItemName  string `gorm:"column:item_name" json:"item_name"`
	ItemPrice int    `gorm:"column:item_price" json:"item_price"`
	ItemPack  string `gorm:"column:item_pack" json:"item_pack"`
}

func (TherapyItem) TableName() string {
	return "items"
}

// RemoveLog jejak penghapusan booking oleh penjual/admin (tabel booking_remove_log)
type RemoveLog struct {
	LogID           int        `gorm:"column:log_id;primaryKey;autoIncrement" json:"log_id"`
	LogBookingID    string     `gorm:"column:log_booking_id" json:"log_booking_id"`
	LogGoldID       int        `gorm:"column:log_gold_id" json:"log_gold_id"`
	LogOutcode      string     `gorm:"column:log_outcode" json:"log_outcode"`
	LogBookingDate  *time.Time `gorm:"column:log_booking_date" json:"log_booking_date"`
	LogBookingStart *time.Time `gorm:"column:log_booking_start" json:"log_booking_start"`
	LogDuration     int        `gorm:"column:log_duration" json:"log_duration"`
	LogTherapyType  *string    `gorm:"column:log_therapy_type" json:"log_therapy_type"`
	LogCustID       *int       `gorm:"column:log_cust_id" json:"log_cust_id"`
	LogCustName     string     `gorm:"column:log_cust_name" json:"log_cust_name"`
	LogStatus       string     `gorm:"column:log_status" json:"log_status"`
	LogRemovedBy    string     `gorm:"column:log_removed_by" json:"log_removed_by"`
	LogRemovedRole  string     `gorm:"column:log_removed_role" json:"log_removed_role"`
	LogRemovedAt    time.Time  `gorm:"column:log_removed_at;autoCreateTime" json:"log_removed_at"`
}

func (RemoveLog) TableName() string {
	return "booking_remove_log"
}

// OutletInfo info outlet untuk validasi booking (tipe THERAPY) dan tenant gold_id
type OutletInfo struct {
	OutletGoldID int    `gorm:"column:outlet_gold_id" json:"outlet_gold_id"`
	OutletCode   string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName   string `gorm:"column:outlet_name" json:"outlet_name"`
	OutletType   string `gorm:"column:outlet_type" json:"outlet_type"`
}

func (OutletInfo) TableName() string {
	return "outlet"
}
