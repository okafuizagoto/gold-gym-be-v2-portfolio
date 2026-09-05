package goldgym

import (
	"context"
	goldQuotaEntity "gold-gym-be/internal/entity/quota"
	goldStorageEntity "gold-gym-be/internal/entity/storage"
	"sort"

	"github.com/pkg/errors"
)

// GetSummary ringkasan pemakaian storage + daftar foto (item + bukti
// pembayaran digabung, terbaru dulu) milik satu user -- untuk menu Storage.
func (s Service) GetSummary(ctx context.Context, goldID int) (goldStorageEntity.StorageSummary, error) {
	var summary goldStorageEntity.StorageSummary

	usedKB, err := s.storage.GetUserStorageUsedKB(ctx, goldID)
	if err != nil {
		return summary, errors.Wrap(err, "[Service][GetSummary][GetUserStorageUsedKB]")
	}
	itemEntries, err := s.storage.ListItemPhotos(ctx, goldID)
	if err != nil {
		return summary, errors.Wrap(err, "[Service][GetSummary][ListItemPhotos]")
	}
	proofEntries, err := s.storage.ListPaymentProofs(ctx, goldID)
	if err != nil {
		return summary, errors.Wrap(err, "[Service][GetSummary][ListPaymentProofs]")
	}

	entries := append(itemEntries, proofEntries...)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].UploadedAt.After(entries[j].UploadedAt)
	})

	summary.UsedKB = usedKB
	summary.QuotaKB = goldQuotaEntity.MaxUserStorageKB
	summary.Entries = entries
	return summary, nil
}

// DeleteEntry menghapus satu foto (aksi hapus di menu Storage). ADMIN tidak
// punya kuota/menu Storage -- ditolak lebih dulu di sini sebagai
// defense-in-depth (FE juga sudah menyembunyikan menunya).
func (s Service) DeleteEntry(ctx context.Context, sourceType string, sourceID, goldID int, isAdmin bool, deletedBy string) error {
	if isAdmin {
		return errors.New("admin tidak memiliki kuota storage")
	}
	switch sourceType {
	case goldQuotaEntity.SourceTypeItemPhoto:
		return s.itemsSvc.DeleteItemPhoto(ctx, sourceID, goldID, deletedBy)
	case goldQuotaEntity.SourceTypePaymentProof:
		return s.salesSvc.DeletePaymentProofPhoto(ctx, sourceID, goldID, deletedBy)
	default:
		return errors.New("jenis storage tidak dikenal")
	}
}
