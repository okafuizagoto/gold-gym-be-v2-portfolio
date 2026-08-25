package goldgym

import (
	"context"
	"fmt"
	goldCustomerEntity "gold-gym-be/internal/entity/customer"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const dbTimeout = 3 * time.Second
const dbTimeoutInsert = 5 * time.Second

func (d *Data) WithTransactionCustomer(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn) //begin
}

func (d *Data) GetCustomer(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldCustomerEntity.Customer, error) {
	var (
		users []goldCustomerEntity.Customer
		err   error
	)
	offset := (page - 1) * length
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	if page == 0 && length == 0 {
		if name == "" {
			err = d.db.WithContext(ctx).Where("cust_gold_id = ? AND cust_outcode = ? AND (? = '' or cust_name like ?)", goldid, outcode, "", "").Find(&users).Error
			if err != nil {
				return nil, errors.Wrap(err, "[DATA] [GetCustomer]")
			}
		}
		if name != "" {
			name = "%" + name + "%"
			err = d.db.WithContext(ctx).Where("cust_gold_id = ? AND cust_outcode = ? AND (? = '' or cust_name like ?)", goldid, outcode, name, name).Find(&users).Error
			if err != nil {
				return nil, errors.Wrap(err, "[DATA] [GetCustomer]")
			}
		}

	} else {
		if name != "" {
			name = "%" + name + "%"
		}
		err = d.db.WithContext(ctx).Where("cust_gold_id = ? AND cust_outcode = ? AND (? = '' or cust_name like ?)", goldid, outcode, name, name).Limit(length).Offset(offset).Find(&users).Error
		if err != nil {
			return nil, errors.Wrap(err, "[DATA] [GetCustomer]")
		}
	}
	return users, err
}

func (d *Data) GetTotalCustomer(ctx context.Context, goldid int, name string, outcode string) (int64, error) {
	var (
		total int64
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	query := d.db.WithContext(ctx).
		Model(&goldCustomerEntity.Customer{}).
		Where("cust_gold_id = ? AND cust_outcode = ?", goldid, outcode)

	if name != "" {
		query = query.Where("cust_name LIKE ?", "%"+name+"%")
	}

	err = query.Debug().Count(&total).Error

	return total, err
}

func (d *Data) GetLastCustomerCode(ctx context.Context, goldid int, code string) (*string, error) {
	var (
		result string
		err    error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err = d.db.WithContext(ctx).Debug().Model(&goldCustomerEntity.Customer{}).Select("cust_code").Where("cust_gold_id = ? AND cust_outcode = ?", goldid, code).Order("cust_code desc, cust_created_at desc").Limit(1).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (d *Data) InsertCustomer(ctx context.Context, tx *gorm.DB, id int, items []goldCustomerEntity.Customer) error {
	var (
		err error
	)
	for _, y := range items {
		y.CustGoldID = id
	}
	if len(items) == 0 {
		return nil
	}

	if len(items) == 1 {
		err = tx.WithContext(ctx).Create(&items[0]).Error
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

				// valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?)")
				valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?)")
				valueArgs = append(valueArgs,
					v.CustID,
					v.CustGoldID,
					v.CustOutcode,
					v.CustCode,
					v.CustName,
					v.CustOutletName,
					v.CustPhone,
					v.CustAddress,
					v.CustEmail,
					v.CustStatus,
					v.CustUpdatedAt,
					v.CustCreatedBy,
				)
			}

			query := fmt.Sprintf(qInsertCustomer, strings.Join(valueStrings, ","))

			if err := tx.WithContext(ctx).Debug().Exec(query, valueArgs...).Error; err != nil {
				return errors.Wrap(err, "[DATA][BulkInsertCustomer]")
			}
		}
	}

	return err
}

func (d *Data) UpdateCustomer(ctx context.Context, cust goldCustomerEntity.Customer) error {

	updates := map[string]interface{}{}

	if cust.CustName != "" {
		updates["cust_name"] = cust.CustName
	}

	if cust.CustPhone != nil {
		updates["cust_phone"] = cust.CustPhone
	}

	if cust.CustAddress != nil {
		updates["cust_address"] = cust.CustAddress
	}

	if cust.CustEmail != nil {
		updates["cust_email"] = cust.CustEmail
	}

	if cust.CustStatus != "" {
		updates["cust_status"] = cust.CustStatus
	}
	if cust.CustOutletName != "" {
		updates["cust_outlet_name"] = cust.CustOutletName
	}

	if len(updates) == 0 {
		return nil
	}

	updates["cust_updated_at"] = time.Now()

	return d.db.WithContext(ctx).Debug().
		Model(&goldCustomerEntity.Customer{}).
		Where("cust_gold_id = ? AND cust_outcode = ? AND cust_id = ?",
			cust.CustGoldID,
			cust.CustOutcode,
			cust.CustID,
		).
		Updates(updates).Error
}

func (d *Data) DeleteCustomer(ctx context.Context, goldid, goldcustomerid int, outcode string) error {
	return d.db.WithContext(ctx).Debug().Where("cust_gold_id = ? AND cust_outcode = ? AND cust_id = ?", goldid, outcode, goldcustomerid).Delete(&goldCustomerEntity.Customer{}).Error
}
