package goldgym

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	goldCustomerTypeEntity "gold-gym-be/internal/entity/customertype"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (s Service) GetCustomerType(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldCustomerTypeEntity.CustomerType, goldCustomerTypeEntity.MetadataPaginationDetail, error) {
	var (
		items          []goldCustomerTypeEntity.CustomerType
		metadataDetail goldCustomerTypeEntity.MetadataPaginationDetail
		totalPage      int
		err            error
	)
	items, err = s.goldgymcustomer.GetCustomerType(ctx, goldid, name, outcode, page, length)
	if err != nil {
		return items, metadataDetail, errors.Wrap(err, "[Service][GetCustomerType]")
	}
	totalCustomerType, err := s.goldgymcustomer.GetTotalCustomerType(ctx, goldid, name, outcode)
	if err != nil {
		return items, metadataDetail, errors.Wrap(err, "[Service][GetCustomerType][GetTotalCustomerType]")
	}
	if page == 0 && length == 0 {
		totalPage = 0
	} else {
		totalPage = int(math.Ceil(float64(totalCustomerType) / float64(length)))
	}
	metadataDetail = goldCustomerTypeEntity.MetadataPaginationDetail{
		Page:      page,
		Limit:     length,
		TotalData: int(totalCustomerType),
		TotalPage: totalPage,
	}
	return items, metadataDetail, err
}

func (s Service) InsertCustomerType(ctx context.Context, id int, items []goldCustomerTypeEntity.CustomerType) (string, error) {
	var (
		result string
		// goldUser goldEntity.GetGoldUserss
		number int
		code   string
		err    error
	)

	// goldUser, err = s.goldgymuser.GetGoldUserByEmail(ctx, items[0].CustomerTypeEmail)
	// if err != nil {
	// 	result = "Gagal"
	// 	return result, errors.Wrap(err, "[Service][InsertCustomerType][GetGoldUserByEmail]")
	// }
	if len(items) == 0 {
		return "Data Kosong", nil
	} else {
		code = items[0].CtOutcode
	}
	lastCustomerTypeCode, err := s.goldgymcustomer.GetLastCustomerTypeCode(ctx, id, code)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][GetLastCustomerTypeCode]")
	}
	for x := range items {
		items[x].CtName = strings.ToLower(items[x].CtName)
		items[x].CtCategory = strings.ToLower(items[x].CtCategory)
		items[x].CtGoldID = id
		if *lastCustomerTypeCode != "" {
			if x == 0 {
				numberStr := (*lastCustomerTypeCode)[3:]
				number, err = strconv.Atoi(numberStr)
				if err != nil {
					return result, errors.Wrap(err, "[Service][InsertCustomerType][ParseNum][Last!Nil=0]")
				}
				number++
				stringNum := fmt.Sprintf("CT%06d", number)
				lastCustomerTypeCode = &stringNum
			}
			if x > 0 {
				numberStr := (*lastCustomerTypeCode)[3:]
				number, err = strconv.Atoi(numberStr)
				if err != nil {
					return result, errors.Wrap(err, "[Service][InsertCustomerType][ParseNum][Last!Nil>0]")
				}
				number++
				stringNum := fmt.Sprintf("CT%06d", number)
				lastCustomerTypeCode = &stringNum
			}
		}
		if *lastCustomerTypeCode == "" {
			if x == 0 {
				value := "CT000001"
				lastCustomerTypeCode = &value
			}
			// if x > 0 {
			// 	numberStr := (*lastCustomerTypeCode)[3:]
			// 	number, err = strconv.Atoi(numberStr)
			// 	if err != nil {
			// 		return result, errors.Wrap(err, "[Service][InsertCustomerType][ParseNum][LastNil>0]")
			// 	}
			// 	number++
			// 	stringNum := fmt.Sprintf("CT%06d", number)
			// 	lastCustomerTypeCode = &stringNum
			// }
		}
		items[x].CtCode = *lastCustomerTypeCode
	}
	// if len(items) >= 1 && len(items) < 350 {
	err = s.goldgymcustomer.WithTransactionCustomerType(ctx, func(tx *gorm.DB) error {
		err = s.goldgymcustomer.InsertCustomerType(ctx, tx, id, items)
		if err != nil {
			return errors.Wrap(err, "[Service][InsertCustomerType]")
		}
		return nil
	})
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertCustomerType]")
	}
	if result != "Gagal" {
		result = "Berhasil"
	}
	return result, err
}

func (s Service) UpdateCustomerType(ctx context.Context, items goldCustomerTypeEntity.CustomerType) (string, error) {
	var (
		result string
		err    error
	)

	err = s.goldgymcustomer.UpdateCustomerType(ctx, items)
	if err != nil {
		result = "Gagal"
		return result, err
	}
	result = "Berhasil"
	return result, err
}

func (s Service) DeleteCustomerType(ctx context.Context, goldid, golditemid int, outcode string) (string, error) {
	var (
		result string
		err    error
	)

	err = s.goldgymcustomer.DeleteCustomerType(ctx, goldid, golditemid, outcode)
	if err != nil {
		result = "Gagal"
		return result, err
	}
	result = "Berhasil"
	return result, err
}
