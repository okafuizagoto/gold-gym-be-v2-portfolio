package selleraccess

import (
	"context"
	"time"

	entity "gold-gym-be/internal/entity/selleraccess"

	"github.com/pkg/errors"
)

const dbTimeout = 5 * time.Second

// GetAll daftar semua outlet aktif (semua tipe, RETAIL & THERAPY) + status
// dua flag ADMIN milik penjual pemiliknya. name (opsional) memfilter
// berdasarkan nama outlet ATAU nama penjual.
func (d *Data) GetAll(ctx context.Context, name string) ([]entity.SellerMenuAccess, error) {
	var rows []entity.SellerMenuAccess
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet.outlet_gold_id, outlet.outlet_code, outlet.outlet_name, outlet.outlet_type, outlet.outlet_address, "+
			"COALESCE(dp.gold_nama, '') AS owner_name, "+
			"(COALESCE(dp.gold_menu_daftar_pembeli, 'Y') = 'Y') AS daftar_pembeli_active, "+
			"(COALESCE(dp.gold_menu_mode_pembeli, 'Y') = 'Y') AS mode_pembeli_active").
		Joins("LEFT JOIN data_peserta dp ON dp.gold_id = outlet.outlet_gold_id").
		Where("outlet.outlet_status = ?", "ACTIVE")

	if name != "" {
		query = query.Where("outlet.outlet_name LIKE ? OR dp.gold_nama LIKE ?", "%"+name+"%", "%"+name+"%")
	}

	if err := query.Order("outlet.outlet_name asc").Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "[DATA][SellerAccess][GetAll]")
	}
	return rows, nil
}

// SetDaftarPembeli aktif/nonaktifkan menu "Daftar Pembeli" untuk 1 akun penjual.
func (d *Data) SetDaftarPembeli(ctx context.Context, goldID int, active bool) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	val := "N"
	if active {
		val = "Y"
	}
	err := d.db.WithContext(ctx).Table("data_peserta").
		Where("gold_id = ?", goldID).
		Update("gold_menu_daftar_pembeli", val).Error
	if err != nil {
		return errors.Wrap(err, "[DATA][SellerAccess][SetDaftarPembeli]")
	}
	return nil
}

// SetModePembeli aktif/nonaktifkan menu "Mode Pembeli" untuk 1 akun penjual.
func (d *Data) SetModePembeli(ctx context.Context, goldID int, active bool) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	val := "N"
	if active {
		val = "Y"
	}
	err := d.db.WithContext(ctx).Table("data_peserta").
		Where("gold_id = ?", goldID).
		Update("gold_menu_mode_pembeli", val).Error
	if err != nil {
		return errors.Wrap(err, "[DATA][SellerAccess][SetModePembeli]")
	}
	return nil
}
