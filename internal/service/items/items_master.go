package goldgym

import (
	"context"
	"fmt"
	goldItemsEntity "gold-gym-be/internal/entity/items"
	goldQuotaEntity "gold-gym-be/internal/entity/quota"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
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
func (s Service) InsertItems(ctx context.Context, id int, items []goldItemsEntity.InsertItem, applyAllOutlets bool) (goldItemsEntity.InsertItemsResult, error) {
	if len(items) == 0 {
		return goldItemsEntity.InsertItemsResult{Message: "Data Kosong"}, nil
	}
	if !applyAllOutlets {
		msg, err := s.insertItemsForOutlet(ctx, id, items, items[0].ItemsOutletCode)
		result := goldItemsEntity.InsertItemsResult{Message: msg}
		// item_id cuma pasti & tunggal kalau tepat 1 item ke 1 outlet -- dipakai
		// FE untuk menempel foto item tepat setelah item dibuat.
		if err == nil && len(items) == 1 {
			result.ItemID = items[0].ItemsID
		}
		return result, err
	}

	codes, err := s.goldgymitems.GetOutletCodesByGoldID(ctx, id)
	if err != nil {
		return goldItemsEntity.InsertItemsResult{Message: "Gagal"}, errors.Wrap(err, "[Service][InsertItems][GetOutletCodesByGoldID]")
	}
	if len(codes) == 0 {
		return goldItemsEntity.InsertItemsResult{Message: "Gagal"}, errors.New("tidak ada outlet aktif untuk akun ini")
	}
	for _, outcode := range codes {
		batch := make([]goldItemsEntity.InsertItem, len(items))
		copy(batch, items)
		if _, err := s.insertItemsForOutlet(ctx, id, batch, outcode); err != nil {
			return goldItemsEntity.InsertItemsResult{Message: "Gagal"}, errors.Wrap(err, "[Service][InsertItems][ApplyAllOutlets][outcode="+outcode+"]")
		}
	}
	return goldItemsEntity.InsertItemsResult{Message: "Berhasil"}, nil
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

// itemPhotoStorageDir direktori penyimpanan foto item -- pakai env var yang SAMA
// dengan foto bukti pembayaran (PHOTO_STORAGE_DIR) supaya satu folder dipakai
// bersama; nama file dibedakan lewat prefix "item_" (lihat SaveItemPhoto).
func itemPhotoStorageDir() string {
	if dir := os.Getenv("PHOTO_STORAGE_DIR"); dir != "" {
		return dir
	}
	return "/root/storages/photos"
}

// bytesToKB pembulatan ke ATAS ke KB terdekat -- dipakai konsisten di semua
// operasi kuota (naik saat upload, turun saat hapus) supaya tidak drift.
func bytesToKB(b int) int {
	if b <= 0 {
		return 0
	}
	return (b + 1023) / 1024
}

// SaveItemPhoto menyimpan foto item: ukuran per file maksimal 2 MB, hanya file
// gambar. File fisik ditulis ke itemPhotoStorageDir(), nama file+ukuran ke
// items.item_photo/item_photo_bytes. Foto lama (kalau ada) dihapus setelah
// foto baru berhasil tersimpan. Kuota storage 30MB/user berlaku untuk semua
// role KECUALI admin (isAdmin=true melewati pengecekan & tidak menambah counter).
func (s Service) SaveItemPhoto(ctx context.Context, itemID int, originalName string, mimeType string, content []byte, goldID int, isAdmin bool) (goldItemsEntity.Item, error) {
	var item goldItemsEntity.Item

	if itemID <= 0 {
		return item, errors.New("item_id wajib diisi")
	}
	if len(content) == 0 {
		return item, errors.New("file foto kosong")
	}
	if len(content) > goldItemsEntity.MaxItemPhotoBytes {
		return item, errors.New("ukuran foto maksimal 2 MB")
	}
	if !strings.HasPrefix(strings.ToLower(mimeType), "image/") {
		return item, errors.New("file harus berupa foto (jpg/png/webp)")
	}

	existing, err := s.goldgymitems.GetItemByID(ctx, itemID)
	if err != nil {
		return item, errors.Wrap(err, "[Service][SaveItemPhoto][GetItemByID]")
	}

	// kuota: hanya selisih (foto baru - foto lama yg diganti) yang perlu
	// muat di sisa kuota, bukan ukuran penuh foto baru
	newKB := bytesToKB(len(content))
	oldKB := bytesToKB(existing.ItemsPhotoBytes)
	deltaKB := newKB - oldKB
	if !isAdmin && deltaKB > 0 {
		usedKB, errQuota := s.goldgymitems.GetUserStorageUsedKB(ctx, goldID)
		if errQuota != nil {
			return item, errors.Wrap(errQuota, "[Service][SaveItemPhoto][GetUserStorageUsedKB]")
		}
		if usedKB+deltaKB > goldQuotaEntity.MaxUserStorageKB {
			return item, errors.New("Kapasitas penyimpanan foto Anda sudah penuh (30 MB). Hapus beberapa foto di menu Storage sebelum mengunggah yang baru.")
		}
	}

	dir := itemPhotoStorageDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return item, errors.Wrap(err, "[Service][SaveItemPhoto][MkdirAll]")
	}
	ext := strings.ToLower(filepath.Ext(originalName))
	if ext == "" {
		ext = ".jpg"
	}
	filename := "item_" + uuid.New().String() + ext
	if err := os.WriteFile(filepath.Join(dir, filename), content, 0o644); err != nil {
		return item, errors.Wrap(err, "[Service][SaveItemPhoto][WriteFile]")
	}

	if err := s.goldgymitems.UpdateItemPhoto(ctx, itemID, filename, len(content)); err != nil {
		// metadata gagal disimpan: hapus file supaya tidak jadi sampah tak tercatat
		_ = os.Remove(filepath.Join(dir, filename))
		return item, errors.Wrap(err, "[Service][SaveItemPhoto][UpdateItemPhoto]")
	}
	// foto lama diganti -- hapus supaya tidak jadi sampah
	if existing.ItemsPhoto != "" && existing.ItemsPhoto != filename {
		_ = os.Remove(filepath.Join(dir, existing.ItemsPhoto))
	}
	if !isAdmin && deltaKB != 0 {
		_ = s.goldgymitems.AddUserStorageUsedKB(ctx, goldID, deltaKB)
	}

	existing.ItemsPhoto = filename
	existing.ItemsPhotoBytes = len(content)
	return existing, nil
}

// DeleteItemPhoto mengosongkan foto item (aksi hapus di menu Storage):
// verifikasi kepemilikan, hapus file fisik, kembalikan kuota, catat history.
func (s Service) DeleteItemPhoto(ctx context.Context, itemID int, goldID int, deletedBy string) error {
	item, err := s.goldgymitems.GetItemByID(ctx, itemID)
	if err != nil {
		return errors.Wrap(err, "[Service][DeleteItemPhoto][GetItemByID]")
	}
	if item.ItemsGoldID != goldID {
		return errors.New("item bukan milik Anda")
	}
	if item.ItemsPhoto == "" {
		return errors.New("item tidak punya foto")
	}
	if err := s.goldgymitems.ClearItemPhoto(ctx, itemID); err != nil {
		return errors.Wrap(err, "[Service][DeleteItemPhoto][ClearItemPhoto]")
	}
	_ = os.Remove(filepath.Join(itemPhotoStorageDir(), item.ItemsPhoto))
	if deltaKB := bytesToKB(item.ItemsPhotoBytes); deltaKB > 0 {
		_ = s.goldgymitems.AddUserStorageUsedKB(ctx, goldID, -deltaKB)
	}
	_ = s.goldgymitems.InsertStorageDeleteHistory(ctx, goldQuotaEntity.StorageDeleteHistory{
		GoldID:           goldID,
		SourceType:       goldQuotaEntity.SourceTypeItemPhoto,
		SourceID:         itemID,
		OriginalFilename: item.ItemsPhoto,
		FileBytes:        item.ItemsPhotoBytes,
		ContextLabel:     item.ItemsName,
		DeletedBy:        deletedBy,
	})
	return nil
}

// GetItemPhotoFile membaca file foto item dari disk untuk ditampilkan di FE.
func (s Service) GetItemPhotoFile(ctx context.Context, itemID int) (string, []byte, error) {
	item, err := s.goldgymitems.GetItemByID(ctx, itemID)
	if err != nil {
		return "", nil, errors.Wrap(err, "[Service][GetItemPhotoFile][GetItemByID]")
	}
	if item.ItemsPhoto == "" {
		return "", nil, errors.New("item tidak punya foto")
	}
	content, err := os.ReadFile(filepath.Join(itemPhotoStorageDir(), item.ItemsPhoto))
	if err != nil {
		return "", nil, errors.Wrap(err, "[Service][GetItemPhotoFile][ReadFile]")
	}
	return item.ItemsPhoto, content, nil
}
