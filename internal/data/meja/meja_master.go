package goldgym

import (
	"context"
	"time"

	goldMejaEntity "gold-gym-be/internal/entity/meja"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const dbTimeout = 3 * time.Second

func (d *Data) InsertMeja(ctx context.Context, rows []goldMejaEntity.Meja) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if err := d.db.WithContext(ctx).Debug().Create(&rows).Error; err != nil {
		return errors.Wrap(err, "[DATA][InsertMeja]")
	}
	return nil
}

func (d *Data) GetMejaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldMejaEntity.Meja, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var rows []goldMejaEntity.Meja
	err := d.db.WithContext(ctx).Debug().
		Where("meja_gold_id = ? AND meja_outcode = ?", goldid, outcode).
		Order("meja_name ASC").
		Find(&rows).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetMejaByOutlet]")
	}
	return rows, nil
}

func (d *Data) GetExistingMejaNames(ctx context.Context, outcode string) (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var names []string
	err := d.db.WithContext(ctx).Debug().
		Model(&goldMejaEntity.Meja{}).
		Where("meja_outcode = ?", outcode).
		Pluck("meja_name", &names).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetExistingMejaNames]")
	}

	existing := make(map[string]bool, len(names))
	for _, n := range names {
		existing[n] = true
	}
	return existing, nil
}

// UpdateMejaStatus mengubah status sekumpulan meja dari fromStatus ke
// toStatus, HANYA untuk baris yang statusnya masih fromStatus saat ini
// (conditional update). Dipakai untuk RELEASE (ISI->KOSONG) -- caller cuma
// pernah melepas meja_id yang dia sendiri reservasi, jadi tidak ada risiko
// menyentuh reservasi sesi lain. JANGAN pakai fungsi ini untuk reserve
// (lihat ReserveMeja) karena "sebagian berhasil" di sini tidak aman
// dikompensasi begitu saja -- baris yang gagal match bisa jadi milik sesi
// lain, bukan hasil percobaan kita sendiri.
func (d *Data) UpdateMejaStatus(ctx context.Context, outcode string, mejaIDs []int, fromStatus, toStatus string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if len(mejaIDs) == 0 {
		return 0, nil
	}

	res := d.db.WithContext(ctx).Debug().
		Model(&goldMejaEntity.Meja{}).
		Where("meja_outcode = ? AND meja_id IN ? AND meja_status = ?", outcode, mejaIDs, fromStatus).
		Update("meja_status", toStatus)
	if res.Error != nil {
		return 0, errors.Wrap(res.Error, "[DATA][UpdateMejaStatus]")
	}
	return res.RowsAffected, nil
}

// ReserveMeja mencoba mengubah status SEMUA mejaIDs dari KOSONG ke ISI
// secara atomik: SELECT ... FOR UPDATE dulu (row lock) untuk tahu persis
// mana yang benar-benar masih KOSONG saat ini, baru UPDATE kalau semuanya
// tersedia. Kalau tidak semua tersedia, TIDAK ADA baris yang diubah sama
// sekali (rollback implisit lewat return error di dalam Transaction) --
// jadi tidak pernah butuh "kompensasi" yang berisiko menyentuh reservasi
// sesi lain seperti pada pendekatan conditional-update biasa.
func (d *Data) ReserveMeja(ctx context.Context, outcode string, mejaIDs []int) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if len(mejaIDs) == 0 {
		return 0, nil
	}

	var affected int64
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockedIDs []int
		if err := tx.Debug().Raw(
			"SELECT meja_id FROM meja WHERE meja_outcode = ? AND meja_id IN ? AND meja_status = ? FOR UPDATE",
			outcode, mejaIDs, goldMejaEntity.MejaStatusKosong,
		).Scan(&lockedIDs).Error; err != nil {
			return err
		}

		if len(lockedIDs) != len(mejaIDs) {
			// Tidak semua yang diminta masih KOSONG -- jangan ubah apapun,
			// biarkan caller tahu lewat affected < len(mejaIDs).
			affected = int64(len(lockedIDs))
			return nil
		}

		res := tx.Debug().Model(&goldMejaEntity.Meja{}).
			Where("meja_outcode = ? AND meja_id IN ?", outcode, lockedIDs).
			Update("meja_status", goldMejaEntity.MejaStatusIsi)
		if res.Error != nil {
			return res.Error
		}
		affected = res.RowsAffected
		return nil
	})
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][ReserveMeja]")
	}
	return affected, nil
}
