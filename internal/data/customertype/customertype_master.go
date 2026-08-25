package goldgym

import (
	"context"
	"fmt"
	goldCustomerTypeEntity "gold-gym-be/internal/entity/customertype"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const dbTimeout = 3 * time.Second
const dbTimeoutInsert = 5 * time.Second

func (d *Data) WithTransactionCustomerType(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn) //begin
}

func (d *Data) GetCustomerType(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldCustomerTypeEntity.CustomerType, error) {
	var (
		users []goldCustomerTypeEntity.CustomerType
		err   error
	)
	offset := (page - 1) * length
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	if page == 0 && length == 0 {
		if name == "" {
			err = d.db.WithContext(ctx).Where("ct_gold_id = ? AND ct_outcode = ? AND (? = '' or ct_name like ?)", goldid, outcode, "", "").Find(&users).Error
			if err != nil {
				return nil, errors.Wrap(err, "[DATA] [GetCustomerType]")
			}
		}
		if name != "" {
			name = "%" + name + "%"
			err = d.db.WithContext(ctx).Where("ct_gold_id = ? AND ct_outcode = ? AND (? = '' or ct_name like ?)", goldid, outcode, name, name).Find(&users).Error
			if err != nil {
				return nil, errors.Wrap(err, "[DATA] [GetCustomerType]")
			}
		}

	} else {
		if name != "" {
			name = "%" + name + "%"
		}
		err = d.db.WithContext(ctx).Where("ct_gold_id = ? AND ct_outcode = ? AND (? = '' or ct_name like ?)", goldid, outcode, name, name).Limit(length).Offset(offset).Find(&users).Error
		if err != nil {
			return nil, errors.Wrap(err, "[DATA] [GetCustomerType]")
		}
	}
	return users, err
}

func (d *Data) GetTotalCustomerType(ctx context.Context, goldid int, name string, outcode string) (int64, error) {
	var (
		total int64
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	query := d.db.WithContext(ctx).
		Model(&goldCustomerTypeEntity.CustomerType{}).
		Where("ct_gold_id = ? AND ct_outcode = ?", goldid, outcode)

	if name != "" {
		query = query.Where("ct_name LIKE ?", "%"+name+"%")
	}

	err = query.Debug().Count(&total).Error

	return total, err
}

func (d *Data) GetLastCustomerTypeCode(ctx context.Context, goldid int, code string) (*string, error) {
	var (
		result string
		err    error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Debug().Model(&goldCustomerTypeEntity.CustomerType{}).Select("ct_code").Where("ct_gold_id = ? AND ct_outcode = ?", goldid, code).Order("ct_code desc, ct_created_at desc").Limit(1).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (d *Data) InsertCustomerType(ctx context.Context, tx *gorm.DB, id int, items []goldCustomerTypeEntity.CustomerType) error {
	var (
		err error
	)
	for _, y := range items {
		y.CtGoldID = id
	}
	if len(items) == 0 {
		return nil
	}

	if len(items) == 1 {
		err = tx.WithContext(ctx).Create(items[0]).Error
	}
	if len(items) > 1 && len(items) <= 350 {
		err = tx.WithContext(ctx).CreateInBatches(items, 350).Error
	}

	if len(items) > 350 {
		batchSize := 500

		for i := 0; i < len(items); i += batchSize {
			end := i + batchSize
			if end > len(items) {
				end = len(items)
			}

			batch := items[i:end]

			valueStrings := []string{}
			valueArgs := []interface{}{}

			for _, v := range batch {
				var description interface{}
				if v.CtDescription != "" {
					description = v.CtDescription
				} else {
					description = nil
				}

				valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?)")
				valueArgs = append(valueArgs,
					v.CtGoldID,
					v.CtOutcode,
					v.CtCode,
					v.CtName,
					v.CtCategory,
					description,
					v.CtStatus,
					v.CtCreatedAt,
					v.CtUpdatedAt,
				)
			}

			query := fmt.Sprintf(qInsertCustomerType, strings.Join(valueStrings, ","))

			if err := tx.WithContext(ctx).Debug().Exec(query, valueArgs...).Error; err != nil {
				return errors.Wrap(err, "[DATA][BulkInsertCustomerType]")
			}
		}
	}

	return err
}

func (d *Data) UpdateCustomerType(ctx context.Context, cust goldCustomerTypeEntity.CustomerType) error {

	updates := map[string]interface{}{}

	if cust.CtName != "" {
		updates["ct_name"] = cust.CtName
	}

	if cust.CtCategory != "" {
		updates["ct_category"] = cust.CtCategory
	}

	if cust.CtDescription != "" {
		updates["ct_description"] = cust.CtDescription
	}

	if cust.CtStatus != "" {
		updates["ct_status"] = cust.CtStatus
	}

	updates["ct_updated_at"] = time.Now()

	if len(updates) == 0 {
		return nil
	}

	updates["ct_updated_at"] = time.Now()

	return d.db.WithContext(ctx).Debug().
		Model(&goldCustomerTypeEntity.CustomerType{}).
		Where("ct_gold_id = ? AND ct_outcode = ? AND ct_id = ?",
			cust.CtGoldID,
			cust.CtCode,
			cust.CtID,
		).
		Updates(updates).Error
}

func (d *Data) DeleteCustomerType(ctx context.Context, goldid, goldcustomerid int, outcode string) error {
	return d.db.WithContext(ctx).Debug().Where("ct_gold_id = ? AND ct_outcode = ? AND ct_id = ?", goldid, outcode, goldcustomerid).Delete(&goldCustomerTypeEntity.CustomerType{}).Error
}
