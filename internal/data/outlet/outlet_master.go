package goldgym

import (
	"context"
	"fmt"
	"strings"
	"time"

	goldOutletEntity "gold-gym-be/internal/entity/outlet"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const dbTimeout = 3 * time.Second

func (d *Data) GetNextSequence(
	ctx context.Context,
	tx *gorm.DB,
	reqs []goldOutletEntity.InsertOutletCounter,
	outletDataJson []goldOutletEntity.InsertOutletCounterJson,
) ([]goldOutletEntity.Outlet, error) {
	var (
		resultOne goldOutletEntity.Outlet
		resultArr []goldOutletEntity.Outlet
	)

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if len(reqs) == 0 {
		return nil, nil
	}

	// ========================
	// 1. BULK UPSERT
	// ========================
	// 	valueStrings := []string{}
	// 	valueArgs := []interface{}{}
	// 	fmt.Println("test - reqs", reqs)
	// 	for _, r := range reqs {
	// 		valueStrings = append(valueStrings, "(?, ?, ?, 1)")
	// 		valueArgs = append(valueArgs, r.Prefix, r.CounterGoldID, r.CounterOutletName)
	// 	}

	// 	query := fmt.Sprintf(`
	// 		INSERT INTO outlet_counter (prefix, counter_gold_id, counter_outlet_name, counter)
	// 		VALUES %s
	// 		ON DUPLICATE KEY UPDATE counter = counter + 1
	// 	`, strings.Join(valueStrings, ","))
	// fmt.Println("query", query)
	// 	if err := tx.WithContext(ctx).Exec(query, valueArgs...).Error; err != nil {
	// 		return nil, err
	// 	}

	type CounterResult struct {
		Prefix            string
		CounterGoldID     int
		CounterOutletName string
		Counter           int
	}

	counterMap := make(map[string][]int)

	keyOccurrence := make(map[string]int)

	results := make([]CounterResult, 0, len(reqs))

	for _, r := range reqs {
		var next int

		err := tx.Raw(`
		SELECT COALESCE(MAX(counter), 0) + 1
		FROM outlet_counter
		WHERE prefix = ? AND counter_gold_id = ?
		FOR UPDATE
	`, r.Prefix, r.CounterGoldID).Scan(&next).Error
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%s_%d", r.Prefix, r.CounterGoldID)

		offset := keyOccurrence[key]
		// counterMap[key] = next + offset
		counter := next + offset

		keyOccurrence[key]++

		results = append(results, CounterResult{
			Prefix:            r.Prefix,
			CounterGoldID:     r.CounterGoldID,
			CounterOutletName: r.CounterOutletName,
			Counter:           counter,
		})
	}

	// ========================
	// 2. BULK SELECT
	// ========================
	// prefixes := []string{}
	// for _, r := range reqs {
	// 	prefixes = append(prefixes, r.Prefix)
	// }
	// fmt.Println("prefixes", prefixes)
	// var results []struct {
	// 	Prefix            string
	// 	CounterOutletName string
	// 	Counter           int
	// }

	// err := tx.WithContext(ctx).
	// 	Raw(`SELECT prefix, counter_outlet_name, counter FROM outlet_counter WHERE prefix IN ?`, prefixes).
	// 	Scan(&results).Error
	// if err != nil {
	// 	return nil, err
	// }

	valueStrings := []string{}
	valueArgs := []interface{}{}

	// for _, r := range reqs {
	for _, r := range results {
		// key := fmt.Sprintf("%s_%d", r.Prefix, r.CounterGoldID)
		// counter := counterMap[key]
		// key := fmt.Sprintf("%s_%d", r.Prefix, r.CounterGoldID)
		counterMap[r.CounterOutletName] = append(counterMap[r.CounterOutletName], r.Counter)
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			r.Prefix,
			r.CounterGoldID,
			r.CounterOutletName,
			r.Counter,
		)
	}

	query := fmt.Sprintf(`
	INSERT INTO outlet_counter (prefix, counter_gold_id, counter_outlet_name, counter)
	VALUES %s
`, strings.Join(valueStrings, ","))

	if err := tx.Debug().Exec(query, valueArgs...).Error; err != nil {
		return nil, err
	}

	// ========================
	// 3. MAP RESULT
	// ========================
	// fmt.Println("test", results)
	// resMap := make(map[string]goldOutletEntity.InsertOutletCounterPrefix)
	// for _, r := range results {
	// 	prefixCounter = goldOutletEntity.InsertOutletCounterPrefix{
	// 		Prefix:  r.Prefix,
	// 		Counter: r.Counter,
	// 	}
	// 	resMap[r.CounterOutletName] = prefixCounter
	// }
	// for _, s := range outletDataJson {
	// 	resultOne = goldOutletEntity.Outlet{
	// 		OutletID:        uuid.New().String(),
	// 		OutletGoldID:    s.CounterGoldID,
	// 		OutletCode:      fmt.Sprintf("%s%03d", resMap[s.OutletName].Prefix, resMap[s.OutletName].Counter),
	// 		OutletName:      s.OutletName,
	// 		OutletAddress:   s.OutletAddress,
	// 		OutletCreatedAt: time.Now(),
	// 		OutletUpdateAt:  nil,
	// 	}
	// 	fmt.Println("test Arr", resMap[s.OutletName].Prefix, resMap[s.OutletName].Counter)
	// 	resultArr = append(resultArr, resultOne)
	// }

	counterIndex := make(map[string]int)

	for _, s := range outletDataJson {
		// key := fmt.Sprintf("%s_%d", s.Prefix, s.CounterGoldID)

		counters := counterMap[s.OutletName] // []int, misal [1, 2, 3]
		idx := counterIndex[s.OutletName]    // ambil index saat ini
		counter := counters[idx]             // ambil counter ke-idx
		counterIndex[s.OutletName]++
		// for _ y := counter {
		resultOne = goldOutletEntity.Outlet{
			OutletID:        uuid.New().String(),
			OutletGoldID:    s.CounterGoldID,
			OutletCode:      fmt.Sprintf("%s%03d", s.Prefix, counter),
			OutletName:      s.OutletName,
			OutletAddress:   s.OutletAddress,
			OutletStatus:    s.OutletStatus,
			OutletCreatedAt: time.Now(),
			OutletUpdateAt:  nil,
		}

		resultArr = append(resultArr, resultOne)
	}

	return resultArr, nil
}

func (d *Data) InsertOutletCounters(ctx context.Context, tx *gorm.DB, outlet goldOutletEntity.Outlet) error {

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	return tx.WithContext(ctx).Create(outlet).Error
}

func (d *Data) InsertOutlet(ctx context.Context, tx *gorm.DB, outlet []goldOutletEntity.Outlet) error {
	var (
		err error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// return tx.WithContext(ctx).Create(outlet).Error

	if len(outlet) == 1 {
		err = tx.WithContext(ctx).Create(outlet[0]).Error
	}
	if len(outlet) > 1 && len(outlet) <= 350 {
		err = tx.WithContext(ctx).CreateInBatches(outlet, 350).Error
	}
	// return err
	if len(outlet) > 350 {

		if len(outlet) == 0 {
			return nil
		}

		batchSize := 500

		for i := 0; i < len(outlet); i += batchSize {
			end := i + batchSize
			if end > len(outlet) {
				end = len(outlet)
			}

			batch := outlet[i:end]

			valueStrings := []string{}
			valueArgs := []interface{}{}

			for _, v := range batch {
				valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?)")
				valueArgs = append(valueArgs,
					v.OutletID,
					v.OutletGoldID,
					v.OutletCode,
					v.OutletName,
					v.OutletAddress,
					v.OutletStatus,
					v.OutletCreatedAt,
					v.OutletUpdateAt,
				)
			}

			query := fmt.Sprintf(qInsertOutlet, strings.Join(valueStrings, ","))

			if err := tx.WithContext(ctx).Debug().Exec(query, valueArgs...).Error; err != nil {
				return errors.Wrap(err, "[DATA][BulkInsertItems]")
			}
		}
		err = nil
	}

	return err
}

func (d *Data) WithTransactionOutlet(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn) //begin
}

func (d *Data) UpdateOutlet(ctx context.Context, updateoutlet goldOutletEntity.UpdateOutlet) error {
	updates := map[string]interface{}{}
	if updateoutlet.OutletAddress != "" {
		updates["outlet_address"] = updateoutlet.OutletAddress
		updates["outlet_status"] = updateoutlet.OutletStatus
	}

	return d.db.WithContext(ctx).Debug().Model(&goldOutletEntity.UpdateOutlet{}).Where("outlet_gold_id = ? AND outlet_code = ?", updateoutlet.OutletGoldID, updateoutlet.OutletCode).Updates(updates).Error
}

func (d *Data) DeleteOutlet(ctx context.Context, goldid int, code string) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var outlet goldOutletEntity.Outlet
		if err := tx.Debug().Where("outlet_gold_id = ? AND outlet_code = ?", goldid, code).First(&outlet).Error; err != nil {
			return err
		}

		history := goldOutletEntity.OutletDeleteHistory{
			OutletGoldID:    outlet.OutletGoldID,
			OutletID:        outlet.OutletID,
			OutletCode:      outlet.OutletCode,
			OutletName:      outlet.OutletName,
			OutletType:      outlet.OutletType,
			OutletAddress:   outlet.OutletAddress,
			OutletStatus:    outlet.OutletStatus,
			OutletCreatedAt: &outlet.OutletCreatedAt,
			DeletedBy:       goldid,
		}
		if err := tx.Debug().Create(&history).Error; err != nil {
			return err
		}

		return tx.Debug().Where("outlet_gold_id = ? AND outlet_code = ?", goldid, code).Delete(&goldOutletEntity.Outlet{}).Error
	})
}

func (d *Data) GetOutlet(ctx context.Context, goldid int, name, status string, page, length int) ([]goldOutletEntity.Outlet, error) {
	var (
		outlet []goldOutletEntity.Outlet
		err    error
	)
	offset := (page - 1) * length
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	if page == 0 && length == 0 {
		if name == "" {
			err = d.db.WithContext(ctx).Debug().Where("outlet_gold_id = ? AND (? = '' or outlet_name like ?) AND (? = '' or outlet_status = ?)", goldid, "", "", status, status).Find(&outlet).Error
			if err != nil {
				return nil, errors.Wrap(err, "[DATA] [GetOutlet]")
			}
		}
		if name != "" {
			name = "%" + name + "%"
			err = d.db.WithContext(ctx).Debug().Where("outlet_gold_id = ? AND (? = '' or outlet_name like ?)", goldid, name, name).Find(&outlet).Error
			if err != nil {
				return nil, errors.Wrap(err, "[DATA] [GetOutlet]")
			}
		}

	} else {
		if name != "" {
			name = "%" + name + "%"
		}
		err = d.db.WithContext(ctx).Debug().Where("outlet_gold_id = ? AND (? = '' or outlet_name like ?)", goldid, name, name).Limit(length).Offset(offset).Find(&outlet).Error
		if err != nil {
			return nil, errors.Wrap(err, "[DATA] [GetOutlet]")
		}
	}
	return outlet, err
}

func (d *Data) GetTotalOutlet(ctx context.Context, goldid int, name string) (int64, error) {
	var (
		total int64
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	query := d.db.WithContext(ctx).
		Model(&goldOutletEntity.Outlet{}).
		Where("outlet_gold_id = ?", goldid)

	if name != "" {
		query = query.Where("outlet_name LIKE ?", "%"+name+"%")
	}

	err = query.Count(&total).Error

	return total, err
}
