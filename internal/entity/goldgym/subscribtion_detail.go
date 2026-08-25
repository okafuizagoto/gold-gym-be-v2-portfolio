package goldgym

type SubscriptionDetail struct {
	GoldId              int     `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldMenuId          int     `gorm:"column:gold_menuid" db:"gold_menuid" json:"gold_menuid"`
	GoldNamaPaket       string  `gorm:"column:gold_namapaket" db:"gold_namapaket" json:"gold_namapaket"`
	GoldNamaLayanan     string  `gorm:"column:gold_namalayanan" db:"gold_namalayanan" json:"gold_namalayanan"`
	GoldHarga           float64 `gorm:"column:gold_harga" db:"gold_harga" json:"gold_harga"`
	GoldJadwal          string  `gorm:"column:gold_jadwal" db:"gold_jadwal" json:"gold_jadwal"`
	GoldListLatihan     string  `gorm:"column:gold_listlatihan" db:"gold_listlatihan" json:"gold_listlatihan"`
	GoldJumlahpertemuan int     `gorm:"column:gold_jumlahpertemuan" db:"gold_jumlahpertemuan" json:"gold_jumlahpertemuan"`
	GoldDurasi          int     `gorm:"column:gold_durasi" db:"gold_durasi" json:"gold_durasi"`
	GoldStatuslangganan string  `gorm:"column:gold_statuslangganan" db:"gold_statuslangganan" json:"gold_statuslangganan"`
}

type DeleteSubs struct {
	GoldId     int `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldMenuId int `gorm:"column:gold_menuid" db:"gold_menuid" json:"gold_menuid"`
}

type UpdateSubs struct {
	GoldJumlahpertemuan int `gorm:"column:gold_jumlahpertemuan" db:"gold_jumlahpertemuan" json:"gold_jumlahpertemuan"`
	GoldId              int `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldMenuId          int `gorm:"column:gold_menuid" db:"gold_menuid" json:"gold_menuid"`
}

func (SubscriptionDetail) TableName() string {
	return "subscription_detail"
}

func (DeleteSubs) TableName() string {
	return "subscription_detail"
}

func (UpdateSubs) TableName() string {
	return "subscription_detail"
}
