package goldgym

import (
	"context"
	"strings"

	goldAreaEntity "gold-gym-be/internal/entity/area"

	"github.com/pkg/errors"
)

func (s Service) InsertArea(ctx context.Context, goldid int, outcode string, areas []goldAreaEntity.InsertArea) error {
	if outcode == "" {
		return errors.New("outlet wajib dipilih")
	}
	if len(areas) == 0 {
		return errors.New("data area kosong")
	}

	for i := range areas {
		areas[i].AreaName = strings.TrimSpace(areas[i].AreaName)
		areas[i].AreaType = strings.ToUpper(strings.TrimSpace(areas[i].AreaType))

		if areas[i].AreaName == "" {
			return errors.New("nama area wajib diisi")
		}
		if areas[i].AreaType != goldAreaEntity.AreaTypeIndoor && areas[i].AreaType != goldAreaEntity.AreaTypeOutdoor {
			return errors.New("tipe area harus INDOOR atau OUTDOOR")
		}
	}

	if err := s.goldgymarea.InsertArea(ctx, goldid, outcode, areas); err != nil {
		return errors.Wrap(err, "[Service][InsertArea]")
	}
	return nil
}

func (s Service) GetAreaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldAreaEntity.Area, error) {
	areas, err := s.goldgymarea.GetAreaByOutlet(ctx, goldid, outcode)
	if err != nil {
		return nil, errors.Wrap(err, "[Service][GetAreaByOutlet]")
	}
	return areas, nil
}
