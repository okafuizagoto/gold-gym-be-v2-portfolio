package goldgym

import (
	"context"
	"math"
	"strings"
	"time"

	goldOutletEntity "gold-gym-be/internal/entity/outlet"

	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

func GenerateOutletPrefix(name string) string {
	words := strings.Fields(strings.ToLower(name))

	if len(words) >= 3 {
		return string(words[0][0]) + string(words[1][0]) + string(words[2][0])
	}

	if len(words) == 2 {
		w1, w2 := words[0], words[1]

		if len(w1) <= 2 {
			return string(w1[0]) + string(w1[len(w1)-1]) + string(w2[len(w2)-1])
		}

		if len(w2) <= 2 && len(w1) >= 3 {
			return string(w1[0]) + string(w1[2]) + string(w2[0])
		}

		if len(w1) >= 3 {
			return string(w1[0]) + string(w1[2]) + string(w2[0])
		}
	}

	if len(words) == 1 {
		w := words[0]
		if len(w) >= 3 {
			return w[:3]
		}
		return w
	}

	return "xxx"
}

func ToTitleCase(input string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(input))
}

func (s Service) InsertOutlet(ctx context.Context, outletData []goldOutletEntity.InsertOutlet, goldid int) error {
	// outlet.OutletID = uuid.New().String()
	var (
		outlet       goldOutletEntity.Outlet
		counterOne   goldOutletEntity.InsertOutletCounterJson
		counterArr   []goldOutletEntity.InsertOutletCounterJson
		counterDbOne goldOutletEntity.InsertOutletCounter
		counterDbArr []goldOutletEntity.InsertOutletCounter
	)
	outlet.OutletGoldID = goldid
	if len(outletData) == 1 {
		outlet.OutletName = outletData[0].OutletName
		outlet.OutletAddress = outletData[0].OutletAddress
		outlet.OutletStatus = outletData[0].OutletStatus
		err := s.goldgymoutlet.WithTransactionOutlet(ctx, func(tx *gorm.DB) error {

			// 1. prefix
			prefix := GenerateOutletPrefix(outlet.OutletName)
			words := strings.ToUpper(prefix)

			counterOne = goldOutletEntity.InsertOutletCounterJson{
				CounterGoldID:   outlet.OutletGoldID,
				Prefix:          words,
				OutletCode:      "",
				OutletName:      ToTitleCase(outlet.OutletName),
				OutletAddress:   ToTitleCase(outlet.OutletAddress),
				OutletStatus:    strings.ToUpper(outlet.OutletStatus),
				OutletCreatedAt: time.Now(),
				OutletUpdateAt:  nil,
			}
			counterArr = append(counterArr, counterOne)
			counterDbOne = goldOutletEntity.InsertOutletCounter{
				CounterGoldID:     outlet.OutletGoldID,
				Prefix:            words,
				CounterOutletName: ToTitleCase(outlet.OutletName),
			}
			counterDbArr = append(counterDbArr, counterDbOne)

			// 2. sequence
			arrInsert, err := s.goldgymoutlet.GetNextSequence(ctx, tx, counterDbArr, counterArr)
			if err != nil {
				return err //rollback
			}

			// // 3. code
			// outlet.OutletCode = fmt.Sprintf("%s%03d", words, seq)
			// outlet.OutletName = ToTitleCase(outletData[0].OutletName)
			// outlet.OutletAddress = ToTitleCase(outletData[0].OutletAddress)
			// outlet.OutletID = uuid.New().String()
			// outlet.OutletCreatedAt = time.Now()
			// outlet.OutletUpdateAt = nil
			// outletArr = append(outletArr, outlet)

			// 4. insert
			if err := s.goldgymoutlet.InsertOutlet(ctx, tx, arrInsert); err != nil {
				return err //rollback
			}

			return nil
		})

		if err != nil {
			return errors.Wrap(err, "[Service][InsertOutlet]")
		}
	} else {
		err := s.goldgymoutlet.WithTransactionOutlet(ctx, func(tx *gorm.DB) error {
			loc, _ := time.LoadLocation("Asia/Jakarta")
			OutletCreatedAt := time.Now().In(loc)
			for _, y := range outletData {
				// 1. prefix
				prefix := GenerateOutletPrefix(y.OutletName)
				words := strings.ToUpper(prefix)

				counterOne = goldOutletEntity.InsertOutletCounterJson{
					CounterGoldID:   outlet.OutletGoldID,
					Prefix:          words,
					OutletCode:      "",
					OutletName:      ToTitleCase(y.OutletName),
					OutletAddress:   ToTitleCase(y.OutletAddress),
					OutletStatus:    strings.ToUpper(y.OutletStatus),
					OutletCreatedAt: OutletCreatedAt,
					OutletUpdateAt:  nil,
				}
				counterArr = append(counterArr, counterOne)

				counterDbOne = goldOutletEntity.InsertOutletCounter{
					CounterGoldID:     outlet.OutletGoldID,
					Prefix:            words,
					CounterOutletName: ToTitleCase(y.OutletName),
				}
				counterDbArr = append(counterDbArr, counterDbOne)
			}
			// 2. sequence
			arrInsert, err := s.goldgymoutlet.GetNextSequence(ctx, tx, counterDbArr, counterArr)
			if err != nil {
				return err //rollback
			}
			err = s.goldgymoutlet.InsertOutlet(ctx, tx, arrInsert)
			if err != nil {
				return err //rollback
			}
			// for _, y := range outletData {
			// 	outlet = goldOutletEntity.Outlet{
			// 		OutletID:        "",
			// 		OutletGoldID:    goldid,
			// 		OutletCode:      fmt.Sprintf("%s%03d", words, seq),
			// 		OutletName:      ToTitleCase(y.OutletName),
			// 		OutletAddress:   ToTitleCase(y.OutletAddress),
			// 		OutletCreatedAt: time.Now(),
			// 		OutletUpdateAt:  nil,
			// 	}
			// 	outletArr = append(outletArr, outlet)
			// }

			// 	if len(items) >= 1 {
			// 	err = s.goldgymitems.InsertItems(ctx, items)
			// 	if err != nil {
			// 		result = "Gagal"
			// 		return result, errors.Wrap(err, "[Service][InsertItems]")
			// 	}
			// }
			// if len(items) > 1000 {
			// 	limitzI := 500
			// 	totalzI := len(items)
			// 	countzI := int(math.Ceil(float64(totalzI) / float64(limitzI)))
			// 	for i := 0; i < countzI; i++ {
			// 		startzI := limitzI * i
			// 		endzI := limitzI * (i + 1)
			// 		if endzI > totalzI {
			// 			endzI = totalzI
			// 		}
			// 		tempUpdatez := items[startzI:endzI]
			// 		err = s.goldgymitems.InsertItems(ctx, tempUpdatez)
			// 		if err != nil {
			// 			result = "Gagal"
			// 			log.Println(err, "[Service][InsertItems]")
			// 		}
			// 	}
			// }
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "[Service][InsertOutlet]][Bulk]")
		}
	}
	return nil
}

func (s Service) UpdateOutlet(ctx context.Context, items goldOutletEntity.UpdateOutlet) error {
	var (
		err error
	)

	err = s.goldgymoutlet.UpdateOutlet(ctx, items)
	if err != nil {
		return err
	}
	return err
}

func (s Service) DeleteOutlet(ctx context.Context, goldid int, code string) error {
	var (
		err error
	)
	err = s.goldgymoutlet.DeleteOutlet(ctx, goldid, code)
	if err != nil {
		return err
	}
	return err
}

func (s Service) GetOutlet(ctx context.Context, goldid int, name, status string, page, length int) ([]goldOutletEntity.Outlet, goldOutletEntity.MetadataPaginationDetail, error) {
	var (
		outlets        []goldOutletEntity.Outlet
		metadataDetail goldOutletEntity.MetadataPaginationDetail
		totalPage      int
		err            error
	)
	outlets, err = s.goldgymoutlet.GetOutlet(ctx, goldid, name, status, page, length)
	if err != nil {
		return outlets, metadataDetail, errors.Wrap(err, "[Service][GetOutlet]")
	}
	totalOutlets, err := s.goldgymoutlet.GetTotalOutlet(ctx, goldid, name)
	if err != nil {
		return outlets, metadataDetail, errors.Wrap(err, "[Service][GetOutlet][GetTotalOutlet]")
	}
	if page == 0 && length == 0 {
		totalPage = 0
	} else {
		totalPage = int(math.Ceil(float64(totalOutlets) / float64(length)))
	}
	metadataDetail = goldOutletEntity.MetadataPaginationDetail{
		Page:      page,
		Limit:     length,
		TotalData: int(totalOutlets),
		TotalPage: totalPage,
	}
	return outlets, metadataDetail, err
}
