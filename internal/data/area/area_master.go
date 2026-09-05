package goldgym

import (
	"context"
	"time"

	goldAreaEntity "gold-gym-be/internal/entity/area"

	"github.com/pkg/errors"
)

const dbTimeout = 3 * time.Second

func (d *Data) InsertArea(ctx context.Context, goldid int, outcode string, areas []goldAreaEntity.InsertArea) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows := make([]goldAreaEntity.Area, 0, len(areas))
	for _, a := range areas {
		rows = append(rows, goldAreaEntity.Area{
			AreaGoldID:   goldid,
			AreaOutcode:  outcode,
			AreaName:     a.AreaName,
			AreaType:     a.AreaType,
			AreaCreateAt: time.Now(),
		})
	}

	if err := d.db.WithContext(ctx).Debug().Create(&rows).Error; err != nil {
		return errors.Wrap(err, "[DATA][InsertArea]")
	}
	return nil
}

func (d *Data) GetAreaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldAreaEntity.Area, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var areas []goldAreaEntity.Area
	err := d.db.WithContext(ctx).Debug().
		Where("area_gold_id = ? AND area_outcode = ?", goldid, outcode).
		Order("area_id ASC").
		Find(&areas).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetAreaByOutlet]")
	}
	return areas, nil
}
