package goldgym

// "strings"

// "gold-gym-be/internal/entity/auth/v2"

// "github.com/dgrijalva/jwt-go"

// "go.opentelemetry.io/otel/attribute"
// "go.opentelemetry.io/otel/trace"

// package goldgym

import (
	"context"
	"fmt"
	firebaseEntity "gold-gym-be/internal/entity/firebase"
	goldStockEntity "gold-gym-be/internal/entity/stock"
	"gold-gym-be/pkg/errors"
	"log"
	"math"
	"strconv"
	"time"

	"gorm.io/gorm"
	// "os"
	// "strings"
	// "gold-gym-be/internal/entity/auth/v2"
	// "github.com/dgrijalva/jwt-go"
	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/trace"
)

func (s Service) GetOneStockProduct(ctx context.Context, stockcode string, stockname string, stockid string) (goldStockEntity.GetOneStock, error) {
	log.Println("service GetGoldUser object")

	users, err := s.goldgymstock.GetOneStockProduct(ctx, stockcode, stockname, stockid)
	log.Println("servicegolduser", users)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetGoldUser]")
	}
	// if len(users) = 0 {}
	log.Printf("testService %+v", users)
	return users, nil
}

// func (s Service) InsertStockSales(ctx context.Context, stock goldStockEntity.InsertStockData) (string, error) {
// 	var (
// 		result string
// 		users  goldStockEntity.GetOneStock
// 		err    error
// 	)
// 	log.Println("service GetGoldUser object")

// 	log.Println("stock.StockData.StockCode", stock.StockData.StockID)

// 	users, err = s.goldgymstock.GetOneStockProduct(ctx, stock.StockData.StockID, "", "")
// 	log.Println("servicegolduser", users)
// 	if err != nil {
// 		return result, errors.Wrap(err, "[Service][GetGoldUser]")
// 	}

// 	usersNull := goldStockEntity.GetOneStock{}

// 	log.Printf("users %+v", users)

// 	if users == usersNull {
// 		lastData, err := s.goldgymstock.GetLastStock(ctx)
// 		if err != nil {
// 			result = "Gagal Insert"
// 			return result, errors.Wrap(err, "[Service][GetLastStock]")
// 		}
// 		if lastData == usersNull {
// 			stock.StockData.StockID = "1"
// 			result, err = s.goldgymstock.InsertStockSales(ctx, stock.StockData)
// 			log.Println("servicegoldusers", result)
// 			if err != nil {
// 				result = "Gagal Insert"
// 				return result, errors.Wrap(err, "[Service][InsertStockSales]")
// 			}
// 			result = "Berhasil Insert"
// 		}

// 		if lastData != usersNull {
// 			id, _ := strconv.Atoi(lastData.StockID)
// 			stockid := id + 1
// 			stock.StockData.StockID = strconv.Itoa(stockid)
// 			result, err = s.goldgymstock.InsertStockSales(ctx, stock.StockData)
// 			log.Println("servicegoldusersssss", result)
// 			if err != nil {
// 				result = "Gagal Insert"
// 				return result, errors.Wrap(err, "[Service][InsertStockSales]")
// 			}
// 			result = "Berhasil Insert"
// 		}
// 		// result, err = s.goldgymstock.InsertStockSales(ctx, stock.StockData)
// 		// log.Println("servicegolduser", result)
// 		// if err != nil {
// 		// 	result = "Gagal Insert"
// 		// 	return result, errors.Wrap(err, "[Service][InsertStockSales]")
// 		// }
// 		// result = "Berhasil Insert"
// 	}
// 	if users != usersNull {
// 		_, err = s.goldgymstock.AddStockByDate(ctx, stock.StockData)
// 		if err != nil {
// 			result = "Gagal Insert - Update"
// 			return result, errors.Wrap(err, "[Service][AddStockByDate]")
// 		}

// 		stockUpdated := users.StockQTY + stock.StockData.StockQTY
// 		result, err = s.goldgymstock.UpdateStockQty(ctx, stockUpdated, stock.StockData.StockID)
// 		if err != nil {
// 			result = "Gagal Update"
// 			return result, errors.Wrap(err, "[Service][UpdateStockQty]")
// 		}
// 		result = "Stock Updated"
// 	}
// 	return result, nil
// }

func (s Service) GetAllStockHeader(ctx context.Context) ([]goldStockEntity.GetOneStock, error) {
	// log.Println("service GetGoldUser object")

	users, err := s.goldgymstock.GetAllStockHeader(ctx)
	// log.Println("servicegolduser", users)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetAllStockHeader]")
	}
	// if len(users) = 0 {}
	// log.Printf("testService %+v", users)
	return users, nil
}

func (s Service) GetAllStockHeaderToRedis(ctx context.Context) ([]goldStockEntity.GetOneStock, error) {
	// log.Println("service GetGoldUser object")

	users, err := s.goldgymstock.GetAllStockHeaderToRedis(ctx)
	// log.Println("servicegolduser", users)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetAllStockHeaderToRedis]")
	}
	// if len(users) = 0 {}
	// log.Printf("testService %+v", users)
	return users, nil
}

func (s Service) GetFromFirebase(ctx context.Context, userID string) (*firebaseEntity.User, error) {
	// log.Println("service GetGoldUser object")

	users, err := s.goldgymstock.GetFromFirebase(ctx, userID)
	// log.Println("servicegolduser", users)
	if err != nil {
		return users, errors.Wrap(err, "[Service][GetFromFirebase]")
	}
	// if len(users) = 0 {}
	// log.Printf("testService %+v", users)
	return users, nil
}

func (s Service) CreateUser(ctx context.Context, user firebaseEntity.User) (string, error) {
	var (
		result string
		err    error
	)

	result, err = s.goldgymstock.CreateUser(ctx, user)
	log.Println("servicegoldusers", result)
	if err != nil {
		result = "Gagal Insert"
		return result, errors.Wrap(err, "[Service][CreateUser]")
	}

	return result, nil
}

func (s Service) GetStockByName(ctx context.Context, goldid int, stockname string) ([]goldStockEntity.GetOneStock, error) {
	var (
		result []goldStockEntity.GetOneStock
		err    error
	)
	result, err = s.goldgymstock.GetStockByName(ctx, goldid, stockname)
	if err != nil {
		return result, errors.Wrap(err, "[Service][GetStockByName]")
	}
	return result, err
}

func (s Service) InsertStockSalesNew(ctx context.Context, goldid int, stock []goldStockEntity.InsertStock) (string, error) {
	var (
		result string
		// goldUser goldEntity.GetGoldUserss
		stockHistory    goldStockEntity.StockHistory
		stockHistoryArr []goldStockEntity.StockHistory
		number          int
		numberHistory   int
		code            string
		itemid          []int
		err             error
	)
	if len(stock) == 0 {
		return "Data Kosong", nil
	} else {
		code = stock[0].StockOutletCode
	}
	for _, h := range stock {
		itemid = append(itemid, h.StockItemID)
	}
	lastItemCode, err := s.goldgymstock.GetLastStockCode(ctx, goldid, code)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertStockSalesNew][GetLastStockCode]")
	}
	lastStockHistoryCode, err := s.goldgymstock.GetLastHistoryStockCode(ctx, goldid, code)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertStockSalesNew][GetLastStockCode]")
	}
	fmt.Println("testlast", lastStockHistoryCode)
	stockExist, err := s.goldgymstock.GetStockExist(ctx, goldid, code, itemid)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertStockSalesNew][GetStockExist]")
	}
	for x, y := range stock {
		stock[x].StockGoldID = goldid
		if *lastItemCode != "" {
			_, ok := stockExist[y.StockItemID]
			if len(stockExist) == 0 || !ok {
				if x == 0 {
					numberStr := (*lastItemCode)[3:]
					number, err = strconv.Atoi(numberStr)
					if err != nil {
						return result, errors.Wrap(err, "[Service][InsertStockSalesNew][ParseNum][Last!Nil=0]")
					}
					number++
					stringNum := fmt.Sprintf("STK%06d", number)
					lastItemCode = &stringNum
				}
				if x > 0 {
					numberStr := (*lastItemCode)[3:]
					number, err = strconv.Atoi(numberStr)
					if err != nil {
						return result, errors.Wrap(err, "[Service][InsertStockSalesNew][ParseNum][Last!Nil>0]")
					}
					number++
					stringNum := fmt.Sprintf("STK%06d", number)
					lastItemCode = &stringNum
				}
			} else {
				value := stockExist[y.StockItemID].StockID
				lastItemCode = &value
			}
		}
		if *lastItemCode == "" {
			if x == 0 {
				value := "STK000001"
				lastItemCode = &value
			}
			// if x > 0 {
			// 	numberStr := (*lastItemCode)[3:]
			// 	number, err = strconv.Atoi(numberStr)
			// 	if err != nil {
			// 		return result, errors.Wrap(err, "[Service][InsertItems][ParseNum][LastNil>0]")
			// 	}
			// 	number++
			// 	stringNum := fmt.Sprintf("ITM%06d", number)
			// 	lastItemCode = &stringNum
			// }
		}
		// ---------------------------------------------------------------------------
		if *lastStockHistoryCode != "" {
			if x == 0 {
				numberStrHistory := (*lastStockHistoryCode)[3:]
				numberHistory, err = strconv.Atoi(numberStrHistory)
				if err != nil {
					return result, errors.Wrap(err, "[Service][InsertStockSalesNew][ParseNumHistory][Last!Nil=0]")
				}
				numberHistory++
				stringNumHistory := fmt.Sprintf("STH%06d", numberHistory)
				lastStockHistoryCode = &stringNumHistory
			}
			if x > 0 {
				numberStrHistory := (*lastStockHistoryCode)[3:]
				numberHistory, err = strconv.Atoi(numberStrHistory)
				if err != nil {
					return result, errors.Wrap(err, "[Service][InsertStockSalesNew][ParseNumHistory][Last!Nil>0]")
				}
				numberHistory = numberHistory + x
				stringNumHistory := fmt.Sprintf("STH%06d", numberHistory)
				lastStockHistoryCode = &stringNumHistory
			}
		}
		if *lastStockHistoryCode == "" {
			if x == 0 {
				value := "STH000001"
				lastStockHistoryCode = &value
			}
			// if x > 0 {
			// 	numberStr := (*lastItemCode)[3:]
			// 	number, err = strconv.Atoi(numberStr)
			// 	if err != nil {
			// 		return result, errors.Wrap(err, "[Service][InsertItems][ParseNum][LastNil>0]")
			// 	}
			// 	number++
			// 	stringNum := fmt.Sprintf("ITM%06d", number)
			// 	lastItemCode = &stringNum
			// }
		}
		stock[x].StockID = *lastItemCode
		_, ok := stockExist[y.StockItemID]
		if len(stockExist) == 0 || !ok {
			stockHistory = goldStockEntity.StockHistory{
				ID:                    *lastStockHistoryCode,
				StockHistoryStockID:   *lastItemCode,
				StockHistoryGoldID:    goldid,
				StockHistoryCode:      y.StockOutletCode,
				StockHistoryItemID:    y.StockItemID,
				StockHistoryStatus:    "PLUS",
				StockHistoryQtyBefore: 0,
				StockHistoryQtyChange: 0,
				StockHistoryQtyAfter:  y.StockQTY,
				StockHistoryNote:      "STOCK AWAL - INSERT",
				StockHistoryCreatedAt: time.Now(),
				StockHistoryCreatedBy: y.StockUpdateBy,
			}
		} else {
			stockHistory = goldStockEntity.StockHistory{
				ID:                    *lastStockHistoryCode,
				StockHistoryStockID:   stockExist[y.StockItemID].StockID,
				StockHistoryGoldID:    goldid,
				StockHistoryCode:      y.StockOutletCode,
				StockHistoryItemID:    y.StockItemID,
				StockHistoryName:      y.StockName,
				StockHistoryPack:      y.StockPack,
				StockHistoryStatus:    "PLUS",
				StockHistoryQtyBefore: stockExist[y.StockItemID].StockQTY,
				StockHistoryQtyChange: y.StockQTY,
				StockHistoryQtyAfter:  stockExist[y.StockItemID].StockQTY + y.StockQTY,
				StockHistoryNote:      "STOCK UPDATE - UPDATE",
				StockHistoryCreatedAt: time.Now(),
				StockHistoryCreatedBy: y.StockUpdateBy,
			}
		}
		stockHistoryArr = append(stockHistoryArr, stockHistory)

	}
	// if len(items) >= 1 && len(items) < 350 {
	err = s.goldgymstock.WithTransactionStock(ctx, func(tx *gorm.DB) error {
		err = s.goldgymstock.InsertStockNew(ctx, tx, goldid, stock)
		if err != nil {
			return errors.Wrap(err, "[Service][InsertStockSalesNew][InsertStockNew]")
		}
		err = s.goldgymstock.InsertStockHistoryNew(ctx, tx, goldid, stockHistoryArr)
		if err != nil {
			return errors.Wrap(err, "[Service][InsertStockSalesNew][InsertStockHistoryNew]")
		}

		return nil
	})
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertStockSalesNew][WithTransactionStock]")
	}

	if result != "Gagal" {
		result = "Berhasil"
	}
	return result, err
}

func (s Service) GetStock(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldStockEntity.GetOneStock, goldStockEntity.MetadataPaginationDetail, error) {
	var (
		items          []goldStockEntity.GetOneStock
		metadataDetail goldStockEntity.MetadataPaginationDetail
		totalPage      int
		err            error
	)
	items, err = s.goldgymstock.GetStock(ctx, goldid, name, outcode, page, length)
	if err != nil {
		return items, metadataDetail, errors.Wrap(err, "[Service][GetItems]")
	}
	totalItems, err := s.goldgymstock.GetTotalStock(ctx, goldid, name, outcode)
	if err != nil {
		return items, metadataDetail, errors.Wrap(err, "[Service][GetItems][GetTotalItems]")
	}
	if page == 0 && length == 0 {
		totalPage = 0
	} else {
		totalPage = int(math.Ceil(float64(totalItems) / float64(length)))
	}
	metadataDetail = goldStockEntity.MetadataPaginationDetail{
		Page:      page,
		Limit:     length,
		TotalData: int(totalItems),
		TotalPage: totalPage,
	}
	return items, metadataDetail, err
}

