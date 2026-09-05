package selleraccess

// SellerMenuAccess satu baris outlet + status 2 flag ADMIN milik penjual
// pemiliknya (gold_menu_daftar_pembeli/gold_menu_mode_pembeli di
// data_peserta). Flag ini melekat di akun penjual (gold_id), bukan per
// outlet — kalau 1 penjual punya beberapa outlet, semua baris outletnya
// akan menampilkan status yang sama.
type SellerMenuAccess struct {
	OutletGoldID        int    `gorm:"column:outlet_gold_id" json:"outlet_gold_id"`
	OutletCode          string `gorm:"column:outlet_code" json:"outlet_code"`
	OutletName          string `gorm:"column:outlet_name" json:"outlet_name"`
	OutletType          string `gorm:"column:outlet_type" json:"outlet_type"`
	OutletAddress       string `gorm:"column:outlet_address" json:"outlet_address"`
	OwnerName           string `gorm:"column:owner_name" json:"owner_name"`
	DaftarPembeliActive bool   `gorm:"column:daftar_pembeli_active" json:"daftar_pembeli_active"`
	ModePembeliActive   bool   `gorm:"column:mode_pembeli_active" json:"mode_pembeli_active"`
}

// SetMenuAccessRequest body PUT untuk 1 akun penjual (gold_id).
type SetMenuAccessRequest struct {
	GoldID int  `json:"gold_id"`
	Active bool `json:"active"`
}
