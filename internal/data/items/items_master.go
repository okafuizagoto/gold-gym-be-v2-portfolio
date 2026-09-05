package goldgym

import (
	"context"
	"fmt"
	goldItemsEntity "gold-gym-be/internal/entity/items"
	goldQuotaEntity "gold-gym-be/internal/entity/quota"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const dbTimeout = 3 * time.Second
const dbTimeoutInsert = 5 * time.Second

func (d *Data) WithTransactionItems(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn) //begin
}

func (d *Data) GetItems(ctx context.Context, goldid int, name string, outcode string, page, length int) ([]goldItemsEntity.Item, error) {
	var (
		users []goldItemsEntity.Item
		err   error
	)
	offset := (page - 1) * length
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	if page == 0 && length == 0 {
		if name == "" {
			err = d.db.WithContext(ctx).Where("item_gold_id = ? AND item_outcode = ? AND (? = '' or item_name like ?)", goldid, outcode, "", "").Find(&users).Error
			if err != nil {
				return nil, errors.Wrap(err, "[DATA] [GetItems]")
			}
		}
		if name != "" {
			name = "%" + name + "%"
			err = d.db.WithContext(ctx).Where("item_gold_id = ? AND item_outcode = ? AND (? = '' or item_name like ?)", goldid, outcode, name, name).Find(&users).Error
			if err != nil {
				return nil, errors.Wrap(err, "[DATA] [GetItems]")
			}
		}

	} else {
		if name != "" {
			name = "%" + name + "%"
		}
		err = d.db.WithContext(ctx).Where("item_gold_id = ? AND item_outcode = ? AND (? = '' or item_name like ?)", goldid, outcode, name, name).Limit(length).Offset(offset).Find(&users).Error
		if err != nil {
			return nil, errors.Wrap(err, "[DATA] [GetItems]")
		}
	}
	return users, err
}

func (d *Data) GetTotalItems(ctx context.Context, goldid int, name string, outcode string) (int64, error) {
	var (
		total int64
		err   error
	)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	query := d.db.WithContext(ctx).
		Model(&goldItemsEntity.Item{}).
		Where("item_gold_id = ? AND item_outcode = ?", goldid, outcode)

	if name != "" {
		query = query.Where("item_name LIKE ?", "%"+name+"%")
	}

	err = query.Debug().Count(&total).Error

	return total, err
}

// GetLastItemCode mengembalikan kode ITM tertinggi di satu outlet ("" jika
// belum ada). Diambil dari MAX angka kode ber-pola ITM%06d saja:
// kode non-ITM (mis. THR-* hasil seed item terapi) diabaikan supaya tidak
// merusak parser angka, dan tidak difilter item_place supaya item hasil
// seed SQL (place NULL) tetap terhitung — dulu filter place membuat
// generator mulai ulang dari ITM000001 dan kena duplicate key.
func (d *Data) GetLastItemCode(ctx context.Context, goldid int, code string) (*string, error) {
	var maxNo int64
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Debug().Model(&goldItemsEntity.Item{}).
		Select("COALESCE(MAX(CAST(SUBSTRING(item_code, 4) AS UNSIGNED)), 0)").
		Where("item_gold_id = ? AND item_outcode = ?", goldid, code).
		Where("item_code REGEXP '^ITM[0-9]+$'").
		Scan(&maxNo).Error
	if err != nil {
		return nil, err
	}
	result := ""
	if maxNo > 0 {
		result = fmt.Sprintf("ITM%06d", maxNo)
	}
	return &result, nil
}

// GetOutletCodesByGoldID mengembalikan kode semua outlet aktif milik satu
// gold_id -- dipakai mode "Semua Outlet" saat Add Items (fan-out per outlet).
func (d *Data) GetOutletCodesByGoldID(ctx context.Context, goldid int) ([]string, error) {
	var codes []string
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Table("outlet").
		Where("outlet_gold_id = ? AND outlet_status = ?", goldid, "ACTIVE").
		Pluck("outlet_code", &codes).Error
	return codes, err
}

func (d *Data) InsertItems(ctx context.Context, tx *gorm.DB, id int, items []goldItemsEntity.InsertItem) error {
	var (
		err error
	)
	for _, y := range items {
		y.ItemsGoldID = id
	}
	if len(items) == 0 {
		return nil
	}

	if len(items) == 1 {
		// pointer wajib supaya GORM menulis balik item_id auto-increment ke
		// items[0] (dipakai InsertItems service untuk upload foto sesudahnya)
		err = tx.WithContext(ctx).Create(&items[0]).Error
	}
	if len(items) > 1 && len(items) <= 350 {
		err = tx.WithContext(ctx).CreateInBatches(&items, 350).Error
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
				var brand interface{}
				if v.ItemsBrand != "" {
					brand = v.ItemsBrand
				} else {
					brand = nil
				}

				var desc interface{}
				if v.ItemsDescription != "" {
					desc = v.ItemsDescription
				} else {
					desc = nil
				}

				valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
				valueArgs = append(valueArgs,
					v.ItemsGoldID,
					v.ItemsOutletCode,
					v.ItemsCode,
					v.ItemsName,
					v.ItemsType,
					v.ItemsPack,
					v.ItemsPrice,
					brand,
					desc,
					v.ItemsStatus,
					"TOKO",
				)
			}

			query := fmt.Sprintf(qInsertItems, strings.Join(valueStrings, ","))
			fmt.Println("query", query)
			fmt.Println("valueArgs", valueArgs)

			if err := tx.WithContext(ctx).Debug().Exec(query, valueArgs...).Error; err != nil {
				return errors.Wrap(err, "[DATA][BulkInsertItems]")
			}
		}
	}

	return err
}

func (d *Data) UpdateItems(ctx context.Context, items goldItemsEntity.UpdateItems) error {

	updates := map[string]interface{}{}

	if items.ItemsName != "" {
		updates["item_name"] = items.ItemsName
	}

	if items.ItemsType != "" {
		updates["item_type"] = items.ItemsType
	}

	if items.ItemsPack != "" {
		updates["item_pack"] = items.ItemsPack
	}

	if items.ItemsPrice != 0 {
		updates["item_price"] = items.ItemsPrice
	}

	if items.ItemsBrand != "" {
		updates["item_brand"] = items.ItemsBrand
	}

	if items.ItemsDescription != "" {
		updates["item_description"] = items.ItemsDescription
	}

	if items.ItemsStatus != "" {
		updates["item_status"] = items.ItemsStatus
	}

	// item_id sudah primary key; filter item_gold_id dari token membuat update
	// jadi 0 baris (tapi tetap "sukses") saat admin mengedit item outlet milik
	// akun lain — cukup jaga dengan outcode outlet yang sedang dibuka
	return d.db.WithContext(ctx).Debug().Model(&goldItemsEntity.UpdateItems{}).Where("item_outcode = ? AND item_id = ?", items.ItemsOutletCode, items.ItemsID).Updates(updates).Error
}

// GetItemGoldID mengembalikan pemilik (item_gold_id) sebuah item — dipakai agar
// EnsureTherapyStock menarget tenant pemilik item, bukan user yang sedang login.
func (d *Data) GetItemGoldID(ctx context.Context, itemID int, outcode string) (int, error) {
	var goldID int
	err := d.db.WithContext(ctx).
		Table("items").
		Select("item_gold_id").
		Where("item_id = ? AND item_outcode = ?", itemID, outcode).
		Scan(&goldID).Error
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][GetItemGoldID]")
	}
	return goldID, nil
}

func (d *Data) DeleteItems(ctx context.Context, goldid, golditemid int, outcode string) error {
	return d.db.WithContext(ctx).Debug().Where("item_gold_id = ? AND item_outcode = ? AND item_id = ?", goldid, outcode, golditemid).Delete(&goldItemsEntity.UpdateItems{}).Error
}

// GetItemByID mengembalikan satu baris item lengkap (dipakai upload/serving foto:
// cek item ada, ambil nama file foto lama untuk dihapus, ambil nama file untuk serve).
func (d *Data) GetItemByID(ctx context.Context, itemID int) (goldItemsEntity.Item, error) {
	var item goldItemsEntity.Item
	err := d.db.WithContext(ctx).Where("item_id = ?", itemID).First(&item).Error
	if err != nil {
		return item, errors.Wrap(err, "[DATA][GetItemByID]")
	}
	return item, nil
}

// UpdateItemPhoto menyimpan nama file + ukuran foto item yang baru diupload.
func (d *Data) UpdateItemPhoto(ctx context.Context, itemID int, filename string, bytes int) error {
	return d.db.WithContext(ctx).Model(&goldItemsEntity.Item{}).
		Where("item_id = ?", itemID).
		Updates(map[string]interface{}{"item_photo": filename, "item_photo_bytes": bytes}).Error
}

// ClearItemPhoto mengosongkan referensi foto item (dipanggil saat user hapus
// foto lewat menu Storage) -- file fisiknya dihapus terpisah oleh service.
func (d *Data) ClearItemPhoto(ctx context.Context, itemID int) error {
	return d.db.WithContext(ctx).Model(&goldItemsEntity.Item{}).
		Where("item_id = ?", itemID).
		Updates(map[string]interface{}{"item_photo": nil, "item_photo_bytes": nil}).Error
}

// GetUserStorageUsedKB kuota storage foto yang sudah terpakai satu user (KB).
func (d *Data) GetUserStorageUsedKB(ctx context.Context, goldID int) (int, error) {
	var usedKB int
	err := d.db.WithContext(ctx).Table("data_peserta").
		Select("gold_storage_used_kb").Where("gold_id = ?", goldID).
		Scan(&usedKB).Error
	if err != nil {
		return 0, errors.Wrap(err, "[DATA][GetUserStorageUsedKB]")
	}
	return usedKB, nil
}

// AddUserStorageUsedKB menambah (atau mengurangi, deltaKB negatif) kuota
// terpakai user -- tidak pernah turun di bawah 0.
func (d *Data) AddUserStorageUsedKB(ctx context.Context, goldID int, deltaKB int) error {
	return d.db.WithContext(ctx).Table("data_peserta").
		Where("gold_id = ?", goldID).
		Update("gold_storage_used_kb", gorm.Expr("GREATEST(0, gold_storage_used_kb + ?)", deltaKB)).Error
}

// InsertStorageDeleteHistory mencatat foto yang dihapus dari menu Storage --
// audit saja, belum ada tampilan untuk tabel ini.
func (d *Data) InsertStorageDeleteHistory(ctx context.Context, h goldQuotaEntity.StorageDeleteHistory) error {
	return d.db.WithContext(ctx).Create(&h).Error
}

// EnsureTherapyStock membuat baris stock otomatis untuk item brand THERAPY yang belum
// punya stok, sehingga item jasa langsung tampil di menu insert sales tanpa Add Stock.
// Qty diisi 0 dan tidak divalidasi saat penjualan (jasa, tanpa batas).
func (d *Data) EnsureTherapyStock(ctx context.Context, goldid int, outcode string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeoutInsert)
	defer cancel()

	query := `
INSERT INTO stock (stock_id, stock_gold_id, stock_outcode, stock_item_id, stock_name, stock_pack, stock_qty, stock_created_at, stock_update_by)
SELECT CONCAT('STK', LPAD(base.maxnum + ROW_NUMBER() OVER (ORDER BY i.item_id), 6, '0')),
       i.item_gold_id, i.item_outcode, i.item_id, i.item_name, i.item_pack, 0, NOW(), 'SYSTEM'
FROM items i
CROSS JOIN (
    SELECT COALESCE(MAX(CAST(SUBSTRING(stock_id, 4) AS UNSIGNED)), 0) AS maxnum
    FROM stock WHERE stock_id LIKE 'STK%'
) base
LEFT JOIN stock s ON s.stock_item_id = i.item_id AND s.stock_outcode = i.item_outcode AND s.stock_gold_id = i.item_gold_id
WHERE i.item_gold_id = ? AND i.item_outcode = ? AND UPPER(i.item_brand) = 'THERAPY' AND s.stock_id IS NULL`

	err := d.db.WithContext(ctx).Exec(query, goldid, outcode).Error
	if err != nil {
		return errors.Wrap(err, "[DATA][EnsureTherapyStock]")
	}
	return nil
}
