package goldgym

import (
	"context"
	"fmt"
	goldItemsEntity "gold-gym-be/internal/entity/items"
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (s Service) GetItems(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldItemsEntity.Item, goldItemsEntity.MetadataPaginationDetail, error) {
	var (
		items          []goldItemsEntity.Item
		metadataDetail goldItemsEntity.MetadataPaginationDetail
		totalPage      int
		err            error
	)
	items, err = s.goldgymitems.GetItems(ctx, goldid, name, outcode, page, length)
	if err != nil {
		return items, metadataDetail, errors.Wrap(err, "[Service][GetItems]")
	}
	totalItems, err := s.goldgymitems.GetTotalItems(ctx, goldid, name, outcode)
	if err != nil {
		return items, metadataDetail, errors.Wrap(err, "[Service][GetItems][GetTotalItems]")
	}
	if page == 0 && length == 0 {
		totalPage = 0
	} else {
		totalPage = int(math.Ceil(float64(totalItems) / float64(length)))
	}
	metadataDetail = goldItemsEntity.MetadataPaginationDetail{
		Page:      page,
		Limit:     length,
		TotalData: int(totalItems),
		TotalPage: totalPage,
	}
	return items, metadataDetail, err
}

// InsertItems adalah entrypoint publik: kalau applyAllOutlets=false, perilaku
// persis seperti sebelumnya (satu outlet, dari items[0].ItemsOutletCode).
// Kalau true ("Semua Outlet" di FE), item yang sama di-fan-out ke semua
// outlet aktif milik gold_id ini -- setiap outlet dapat penomoran item_code
// sendiri (bukan baris yang dibagi), lewat insertItemsForOutlet yang sama.
func (s Service) InsertItems(ctx context.Context, id int, items []goldItemsEntity.InsertItem, applyAllOutlets bool) (string, error) {
	if len(items) == 0 {
		return "Data Kosong", nil
	}
	if !applyAllOutlets {
		return s.insertItemsForOutlet(ctx, id, items, items[0].ItemsOutletCode)
	}

	codes, err := s.goldgymitems.GetOutletCodesByGoldID(ctx, id)
	if err != nil {
		return "Gagal", errors.Wrap(err, "[Service][InsertItems][GetOutletCodesByGoldID]")
	}
	if len(codes) == 0 {
		return "Gagal", errors.New("tidak ada outlet aktif untuk akun ini")
	}
	for _, outcode := range codes {
		batch := make([]goldItemsEntity.InsertItem, len(items))
		copy(batch, items)
		if _, err := s.insertItemsForOutlet(ctx, id, batch, outcode); err != nil {
			return "Gagal", errors.Wrap(err, "[Service][InsertItems][ApplyAllOutlets][outcode="+outcode+"]")
		}
	}
	return "Berhasil", nil
}

// insertItemsForOutlet adalah body InsertItems yang lama, sekarang menerima
// outcode secara eksplisit (bukan cuma dari items[0].ItemsOutletCode) supaya
// bisa dipanggil berulang untuk mode "Semua Outlet".
func (s Service) insertItemsForOutlet(ctx context.Context, id int, items []goldItemsEntity.InsertItem, code string) (string, error) {
	var (
		result string
		number int
		err    error
	)
	for x := range items {
		items[x].ItemsOutletCode = code
	}
	lastItemCode, err := s.goldgymitems.GetLastItemCode(ctx, id, code)
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][GetLastItemCode]")
	}
	for x := range items {
		items[x].ItemsType = "STOCK"
		items[x].ItemsGoldID = id
		if *lastItemCode != "" {
			if x == 0 {
				numberStr := (*lastItemCode)[3:]
				number, err = strconv.Atoi(numberStr)
				if err != nil {
					return result, errors.Wrap(err, "[Service][InsertItems][ParseNum][Last!Nil=0]")
				}
				number++
				stringNum := fmt.Sprintf("ITM%06d", number)
				lastItemCode = &stringNum
			}
			if x > 0 {
				numberStr := (*lastItemCode)[3:]
				number, err = strconv.Atoi(numberStr)
				if err != nil {
					return result, errors.Wrap(err, "[Service][InsertItems][ParseNum][Last!Nil>0]")
				}
				number++
				stringNum := fmt.Sprintf("ITM%06d", number)
				lastItemCode = &stringNum
			}
		}
		if *lastItemCode == "" {
			if x == 0 {
				value := "ITM000001"
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
		items[x].ItemsCode = *lastItemCode
		items[x].ItemsPlace = "TOKO"
	}
	// if len(items) >= 1 && len(items) < 350 {
	err = s.goldgymitems.WithTransactionItems(ctx, func(tx *gorm.DB) error {
		err = s.goldgymitems.InsertItems(ctx, tx, id, items)
		if err != nil {
			return errors.Wrap(err, "[Service][InsertItems]")
		}
		return nil
	})
	if err != nil {
		result = "Gagal"
		return result, errors.Wrap(err, "[Service][InsertItems]")
	}
	// item brand THERAPY (jasa) otomatis dibuatkan baris stock agar langsung
	// tampil di menu insert sales tanpa perlu Add Stock
	for _, y := range items {
		if strings.EqualFold(strings.TrimSpace(y.ItemsBrand), "THERAPY") {
			if err := s.goldgymitems.EnsureTherapyStock(ctx, id, code); err != nil {
				return "Gagal", errors.Wrap(err, "[Service][InsertItems][EnsureTherapyStock]")
			}
			break
		}
	}
	// }
	// if len(items) > 350 {
	// 	// limitzI := 500
	// 	// totalzI := len(items)
	// 	// countzI := int(math.Ceil(float64(totalzI) / float64(limitzI)))
	// 	// for i := 0; i < countzI; i++ {
	// 	// 	startzI := limitzI * i
	// 	// 	endzI := limitzI * (i + 1)
	// 	// 	if endzI > totalzI {
	// 	// 		endzI = totalzI
	// 	// 	}
	// 	// 	tempUpdatez := items[startzI:endzI]
	// 	err = s.goldgymitems.InsertItems(ctx, items)
	// 	// if err != nil {
	// 	// 	result = "Gagal"
	// 	// 	log.Println(err, "[Service][InsertItems]")
	// 	// }
	// 	// }
	// }
	if result != "Gagal" {
		result = "Berhasil"
	}
	return result, err
}

func (s Service) UpdateItems(ctx context.Context, items goldItemsEntity.UpdateItems) (string, error) {
	var (
		result string
		err    error
	)

	// pemilik item bisa berbeda dari user login (admin mengedit item outlet
	// milik akun lain) — pakai gold_id pemilik untuk EnsureTherapyStock
	if ownerGoldID, errOwner := s.goldgymitems.GetItemGoldID(ctx, items.ItemsID, items.ItemsOutletCode); errOwner == nil && ownerGoldID > 0 {
		items.ItemsGoldID = ownerGoldID
	}

	err = s.goldgymitems.UpdateItems(ctx, items)
	if err != nil {
		result = "Gagal"
		return result, err
	}
	// merek diubah menjadi THERAPY → buatkan baris stock otomatis
	if strings.EqualFold(strings.TrimSpace(items.ItemsBrand), "THERAPY") {
		if err := s.goldgymitems.EnsureTherapyStock(ctx, items.ItemsGoldID, items.ItemsOutletCode); err != nil {
			return "Gagal", errors.Wrap(err, "[Service][UpdateItems][EnsureTherapyStock]")
		}
	}
	result = "Berhasil"
	return result, err
}

func (s Service) DeleteItems(ctx context.Context, goldid, golditemid int, outcode string) (string, error) {
	var (
		result string
		err    error
	)

	err = s.goldgymitems.DeleteItems(ctx, goldid, golditemid, outcode)
	if err != nil {
		result = "Gagal"
		return result, err
	}
	result = "Berhasil"
	return result, err
}
