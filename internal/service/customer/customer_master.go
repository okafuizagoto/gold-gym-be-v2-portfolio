package goldgym

import (
	"context"
	"fmt"
	"math"
	"strconv"

	goldCustomerEntity "gold-gym-be/internal/entity/customer"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (s Service) GetCustomer(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldCustomerEntity.Customer, goldCustomerEntity.MetadataPaginationDetail, error) {
	var (
		items          []goldCustomerEntity.Customer
		metadataDetail goldCustomerEntity.MetadataPaginationDetail
		totalPage      int
		err            error
	)
	items, err = s.goldgymcustomer.GetCustomer(ctx, goldid, name, outcode, page, length)
	if err != nil {
		return items, metadataDetail, errors.Wrap(err, "[Service][GetCustomer]")
	}
	totalCustomer, err := s.goldgymcustomer.GetTotalCustomer(ctx, goldid, name, outcode)
	if err != nil {
		return items, metadataDetail, errors.Wrap(err, "[Service][GetCustomer][GetTotalCustomer]")
	}
	if page == 0 && length == 0 {
		totalPage = 0
	} else {
		totalPage = int(math.Ceil(float64(totalCustomer) / float64(length)))
	}
	metadataDetail = goldCustomerEntity.MetadataPaginationDetail{
		Page:      page,
		Limit:     length,
		TotalData: int(totalCustomer),
		TotalPage: totalPage,
	}
	return items, metadataDetail, err
}

func (s Service) InsertCustomer(ctx context.Context, id int, items []goldCustomerEntity.Customer) (string, error) {
	var (
		result string
		// goldUser goldEntity.GetGoldUserss
		number int
		code   string
		err    error
	)

	// goldUser, err = s.goldgymuser.GetGoldUserByEmail(ctx, items[0].CustomerEmail)
	// if err != nil {
	// 	result = "Gagal"
	// 	return result, errors.Wrap(err, "[Service][InsertCustomer][GetGoldUserByEmail]")
	// }
	if len(items) == 0 {
		return "Data Kosong", nil
	} else {
		code = items[0].CustOutcode
	}
	lastCustomerCode, err := s.goldgymcustomer.GetLastCustomerCode(ctx, id, code)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][GetLastCustomerCode]")
	}
	for x := range items {
		if items[x].CustEmail != nil && *items[x].CustEmail == "" {
			items[x].CustEmail = nil
		}
		if items[x].CustPhone != nil && *items[x].CustPhone == "" {
			items[x].CustPhone = nil
		}
		if items[x].CustAddress != nil && *items[x].CustAddress == "" {
			items[x].CustAddress = nil
		}
		items[x].CustGoldID = id
		if *lastCustomerCode != "" {
			if x == 0 {
				numberStr := (*lastCustomerCode)[3:]
				number, err = strconv.Atoi(numberStr)
				if err != nil {
					return result, errors.Wrap(err, "[Service][InsertCustomer][ParseNum][Last!Nil=0]")
				}
				number++
				stringNum := fmt.Sprintf("CST%06d", number)
				lastCustomerCode = &stringNum
			}
			if x > 0 {
				numberStr := (*lastCustomerCode)[3:]
				number, err = strconv.Atoi(numberStr)
				if err != nil {
					return result, errors.Wrap(err, "[Service][InsertCustomer][ParseNum][Last!Nil>0]")
				}
				number++
				stringNum := fmt.Sprintf("CST%06d", number)
				lastCustomerCode = &stringNum
			}
		}
		if *lastCustomerCode == "" {
			if x == 0 {
				value := "CST000001"
				lastCustomerCode = &value
			}
			// if x > 0 {
			// 	numberStr := (*lastCustomerCode)[3:]
			// 	number, err = strconv.Atoi(numberStr)
			// 	if err != nil {
			// 		return result, errors.Wrap(err, "[Service][InsertCustomer][ParseNum][LastNil>0]")
			// 	}
			// 	number++
			// 	stringNum := fmt.Sprintf("CST%06d", number)
			// 	lastCustomerCode = &stringNum
			// }
		}
		items[x].CustCode = *lastCustomerCode
	}
	// if len(items) >= 1 && len(items) < 350 {
	err = s.goldgymcustomer.WithTransactionCustomer(ctx, func(tx *gorm.DB) error {
		err = s.goldgymcustomer.InsertCustomer(ctx, tx, id, items)
		if err != nil {
			return errors.Wrap(err, "[Service][InsertCustomer]")
		}
		return nil
	})
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertCustomer]")
	}
	if result != "Gagal" {
		result = "Berhasil"
	}
	return result, err
}

func (s Service) UpdateCustomer(ctx context.Context, items goldCustomerEntity.Customer) (string, error) {
	var (
		result string
		err    error
	)

	err = s.goldgymcustomer.UpdateCustomer(ctx, items)
	if err != nil {
		result = "Gagal"
		return result, err
	}
	result = "Berhasil"
	return result, err
}

func (s Service) DeleteCustomer(ctx context.Context, goldid, golditemid int, outcode string) (string, error) {
	var (
		result string
		err    error
	)

	err = s.goldgymcustomer.DeleteCustomer(ctx, goldid, golditemid, outcode)
	if err != nil {
		result = "Gagal"
		return result, err
	}
	result = "Berhasil"
	return result, err
}
