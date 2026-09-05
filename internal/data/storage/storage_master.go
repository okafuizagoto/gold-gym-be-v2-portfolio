package goldgym

import (
	"context"
	goldQuotaEntity "gold-gym-be/internal/entity/quota"
	goldStorageEntity "gold-gym-be/internal/entity/storage"
	"time"

	"github.com/pkg/errors"
)

// bytesToKB pembulatan ke ATAS ke KB terdekat -- sama dengan helper serupa di
// modul items/sales, dipakai untuk tampilan ukuran per entry di menu Storage.
func bytesToKB(b int) int {
	if b <= 0 {
		return 0
	}
	return (b + 1023) / 1024
}

// ListItemPhotos foto item milik satu user (menu Storage).
func (d *Data) ListItemPhotos(ctx context.Context, goldID int) ([]goldStorageEntity.StorageEntry, error) {
	type row struct {
		ItemID     int
		ItemName   string
		PhotoBytes int
		UpdatedAt  time.Time
	}
	var rows []row
	err := d.db.WithContext(ctx).Table("items").
		Select("item_id as item_id, item_name as item_name, item_photo_bytes as photo_bytes, item_updated_at as updated_at").
		Where("item_gold_id = ? AND item_photo IS NOT NULL AND item_photo != ''", goldID).
		Scan(&rows).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][ListItemPhotos]")
	}
	entries := make([]goldStorageEntity.StorageEntry, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, goldStorageEntity.StorageEntry{
			SourceType:  goldQuotaEntity.SourceTypeItemPhoto,
			SourceID:    r.ItemID,
			Label:       "Foto Katalog Item",
			ContextText: r.ItemName,
			SizeKB:      bytesToKB(r.PhotoBytes),
			UploadedAt:  r.UpdatedAt,
		})
	}
	return entries, nil
}

// ListPaymentProofs bukti pembayaran milik satu user (menu Storage) -- JOIN
// th_sale hanya untuk nomor nota (sale_trancnum), biar labelnya informatif.
func (d *Data) ListPaymentProofs(ctx context.Context, goldID int) ([]goldStorageEntity.StorageEntry, error) {
	type row struct {
		ProofID    int
		ProofBytes int64
		Trancnum   string
		UploadedAt time.Time
	}
	var rows []row
	err := d.db.WithContext(ctx).Table("payment_proof AS pp").
		Select("pp.proof_id as proof_id, pp.proof_bytes as proof_bytes, ts.sale_trancnum as trancnum, pp.proof_uploaded_at as uploaded_at").
		// COLLATE eksplisit -- payment_proof.proof_sale_id (utf8mb4_uca1400_ai_ci,
		// kolom lebih baru, dibuat setelah default collation server berubah)
		// vs th_sale.sale_id (utf8mb4_unicode_ci, tabel lama) beda collation,
		// tanpa ini MariaDB menolak perbandingan dengan Error 1267 "Illegal
		// mix of collations".
		Joins("JOIN th_sale ts ON ts.sale_id = pp.proof_sale_id COLLATE utf8mb4_unicode_ci").
		Where("pp.proof_gold_id = ?", goldID).
		Scan(&rows).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][ListPaymentProofs]")
	}
	entries := make([]goldStorageEntity.StorageEntry, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, goldStorageEntity.StorageEntry{
			SourceType:  goldQuotaEntity.SourceTypePaymentProof,
			SourceID:    r.ProofID,
			Label:       "Bukti Pembayaran Transaksi POS",
			ContextText: "Nota " + r.Trancnum,
			SizeKB:      bytesToKB(int(r.ProofBytes)),
			UploadedAt:  r.UploadedAt,
		})
	}
	return entries, nil
}

// GetUserStorageUsedKB kuota storage foto yang sudah terpakai satu user (KB).
func (d *Data) GetUserStorageUsedKB(ctx context.Context, goldID int) (int, error) {
	var usedKB int
	err := d.db.WithContext(ctx).Table("data_peserta").
		Select("gold_storage_used_kb").Where("gold_id = ?", goldID).
		Scan(&usedKB).Error
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][GetUserStorageUsedKB]")
	}
	return usedKB, nil
}
