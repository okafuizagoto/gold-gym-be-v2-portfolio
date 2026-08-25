package goldgym

import (
	"context"

	goldSaleEntity "gold-gym-be/internal/entity/sales"

	"github.com/pkg/errors"
)

// IsPaymentProofFeatureEnabled menentukan apakah fitur bukti pembayaran
// (upload di POS/belanja pembeli, tombol lihat di Sales History) tampil untuk
// user ini. Tiga gerbang independen — global, per outlet (jika outcode
// diisi), per user — SEMUA harus 'Y' (AND logic); admin bisa mematikan lewat
// salah satu gerbang saja. outcode boleh kosong (mis. pembeli belum pilih
// outlet); dalam kasus itu gerbang outlet dilewati.
func (s Service) IsPaymentProofFeatureEnabled(ctx context.Context, userGoldID int, outcode string) (bool, error) {
	global, err := s.goldgymsale.IsPaymentProofGloballyEnabled(ctx)
	if err != nil {
		return false, errors.Wrap(err, "[Service][IsPaymentProofFeatureEnabled][global]")
	}
	if !global {
		return false, nil
	}

	if userGoldID > 0 {
		userEnabled, err := s.goldgymsale.IsUserProofEnabled(ctx, userGoldID)
		if err != nil {
			return false, errors.Wrap(err, "[Service][IsPaymentProofFeatureEnabled][user]")
		}
		if !userEnabled {
			return false, nil
		}
	}

	if outcode != "" {
		ownerGoldID, err := s.ResolveOutletGoldID(ctx, outcode)
		if err != nil {
			return false, errors.Wrap(err, "[Service][IsPaymentProofFeatureEnabled][resolve outlet]")
		}
		if ownerGoldID > 0 {
			outletEnabled, err := s.goldgymsale.IsOutletProofEnabled(ctx, ownerGoldID, outcode)
			if err != nil {
				return false, errors.Wrap(err, "[Service][IsPaymentProofFeatureEnabled][outlet]")
			}
			if !outletEnabled {
				return false, nil
			}
		}
	}

	return true, nil
}

// ---- ADMIN: kelola 3 gerbang ----

func (s Service) GetPaymentProofGlobal(ctx context.Context) (bool, error) {
	enabled, err := s.goldgymsale.IsPaymentProofGloballyEnabled(ctx)
	if err != nil {
		return false, errors.Wrap(err, "[Service][GetPaymentProofGlobal]")
	}
	return enabled, nil
}

func (s Service) SetPaymentProofGlobal(ctx context.Context, enabled bool, updatedBy string) error {
	if err := s.goldgymsale.SetPaymentProofGlobal(ctx, enabled, updatedBy); err != nil {
		return errors.Wrap(err, "[Service][SetPaymentProofGlobal]")
	}
	return nil
}

func (s Service) GetProofAccessOutlets(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessOutlet, error) {
	outlets, err := s.goldgymsale.GetProofAccessOutlets(ctx, search)
	if err != nil {
		return nil, errors.Wrap(err, "[Service][GetProofAccessOutlets]")
	}
	return outlets, nil
}

// SaveProofAccessOutlets menyimpan keputusan admin untuk sekumpulan outlet
// yang sedang tampil di pencarian; outlet di luar daftar ini tidak tersentuh.
func (s Service) SaveProofAccessOutlets(ctx context.Context, items []goldSaleEntity.ProofAccessOutletItem) error {
	for _, it := range items {
		if it.GoldID <= 0 || it.Outcode == "" {
			continue
		}
		if err := s.goldgymsale.SetProofOutletEnabled(ctx, it.GoldID, it.Outcode, it.Enabled); err != nil {
			return errors.Wrap(err, "[Service][SaveProofAccessOutlets]")
		}
	}
	return nil
}

func (s Service) GetProofAccessUsers(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessUser, error) {
	users, err := s.goldgymsale.GetProofAccessUsers(ctx, search)
	if err != nil {
		return nil, errors.Wrap(err, "[Service][GetProofAccessUsers]")
	}
	return users, nil
}

// SaveProofAccessUsers menyimpan keputusan admin untuk sekumpulan user yang
// sedang tampil di pencarian; user di luar daftar ini tidak tersentuh.
func (s Service) SaveProofAccessUsers(ctx context.Context, items []goldSaleEntity.ProofAccessUserItem) error {
	for _, it := range items {
		if it.GoldID <= 0 {
			continue
		}
		if err := s.goldgymsale.SetProofUserEnabled(ctx, it.GoldID, it.Enabled); err != nil {
			return errors.Wrap(err, "[Service][SaveProofAccessUsers]")
		}
	}
	return nil
}
