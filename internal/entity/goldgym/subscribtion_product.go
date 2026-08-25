package goldgym

type Subscription struct {
	// GoldMenuId          int     `gorm:"column:gold_menuid" json:"gold_menuid"`
	GoldNamaPaket       string  `gorm:"column:gold_namapaket" db:"gold_namapaket" json:"gold_namapaket"`
	GoldNamaLayanan     string  `gorm:"column:gold_namalayanan" db:"gold_namalayanan" json:"gold_namalayanan"`
	GoldHarga           float64 `gorm:"column:gold_harga" db:"gold_harga" json:"gold_harga"`
	GoldJadwal          string  `gorm:"column:gold_jadwal" db:"gold_jadwal" json:"gold_jadwal"`
	GoldListLatihan     string  `gorm:"column:gold_listlatihan" db:"gold_listlatihan" json:"gold_listlatihan"`
	GoldJumlahpertemuan int     `gorm:"column:gold_jumlahpertemuan" db:"gold_jumlahpertemuan" json:"gold_jumlahpertemuan"`
	GoldDurasi          int     `gorm:"column:gold_durasi" db:"gold_durasi" json:"gold_durasi"`
}

func (Subscription) TableName() string {
	return "subscription_product"
}
