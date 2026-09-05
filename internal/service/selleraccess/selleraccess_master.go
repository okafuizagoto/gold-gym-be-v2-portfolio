package selleraccess

import (
	"context"

	entity "gold-gym-be/internal/entity/selleraccess"
)

func (s Service) GetAll(ctx context.Context, name string) ([]entity.SellerMenuAccess, error) {
	return s.data.GetAll(ctx, name)
}

func (s Service) SetDaftarPembeli(ctx context.Context, goldID int, active bool) error {
	return s.data.SetDaftarPembeli(ctx, goldID, active)
}

func (s Service) SetModePembeli(ctx context.Context, goldID int, active bool) error {
	return s.data.SetModePembeli(ctx, goldID, active)
}
