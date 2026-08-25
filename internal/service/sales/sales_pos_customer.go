package goldgym

import (
	"context"
	"strings"

	goldSaleEntity "gold-gym-be/internal/entity/sales"

	"github.com/pkg/errors"
)

// IsPosCustomerRequired menentukan apakah transaksi POS di outlet ini WAJIB
// mengisi nama customer. Wajib bila outlet RETAIL DAN belum diberi akses
// "POS tanpa customer" oleh admin. THERAPY tidak pernah wajib (customer dari booking).
func (s Service) IsPosCustomerRequired(ctx context.Context, goldid int, outcode string) (bool, error) {
	outletType, err := s.goldgymsale.GetOutletTypeByCode(ctx, goldid, outcode)
	if err != nil {
		return false, errors.Wrap(err, "[Service][IsPosCustomerRequired][type]")
	}
	if strings.ToUpper(outletType) != "RETAIL" {
		return false, nil
	}
	optional, err := s.goldgymsale.IsPosCustomerOptional(ctx, goldid, outcode)
	if err != nil {
		return false, errors.Wrap(err, "[Service][IsPosCustomerRequired][optional]")
	}
	return !optional, nil
}

// GetPosOutletsForAdmin daftar outlet RETAIL untuk layar admin (cari alamat/nama).
func (s Service) GetPosOutletsForAdmin(ctx context.Context, search string) ([]goldSaleEntity.PosOutlet, error) {
	outlets, err := s.goldgymsale.GetPosOutletsForAdmin(ctx, search)
	if err != nil {
		return nil, errors.Wrap(err, "[Service][GetPosOutletsForAdmin]")
	}
	return outlets, nil
}

// SavePosCustomerAccess menyimpan keputusan admin untuk sekumpulan outlet yang
// sedang tampil di pencarian: optional=true diberi akses, false dicabut.
// Outlet di luar daftar ini tidak tersentuh.
func (s Service) SavePosCustomerAccess(ctx context.Context, items []goldSaleEntity.PosAccessItem, addedBy string) error {
	for _, it := range items {
		if it.GoldID <= 0 || strings.TrimSpace(it.Outcode) == "" {
			continue
		}
		if err := s.goldgymsale.SetPosCustomerOptional(ctx, it.GoldID, it.Outcode, it.Optional, addedBy); err != nil {
			return errors.Wrap(err, "[Service][SavePosCustomerAccess]")
		}
	}
	return nil
}
