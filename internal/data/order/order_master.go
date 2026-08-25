package order

import (
	"context"
	"time"

	orderEntity "gold-gym-be/internal/entity/order"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const dbTimeout = 5 * time.Second

// GetPublicOutlets mengambil outlet non-THERAPY aktif yang SUDAH DIKURASI admin
// (ada di buyer_visible_outlet) untuk dipilih pembeli. Join data_peserta untuk
// nama pemilik/penjual. INNER JOIN kurasi => hanya outlet yang di-approve tampil.
func (d *Data) GetPublicOutlets(ctx context.Context, name string) ([]orderEntity.PublicOutlet, error) {
	var outlets []orderEntity.PublicOutlet
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet.outlet_gold_id, outlet.outlet_code, outlet.outlet_name, outlet.outlet_type, outlet.outlet_address, COALESCE(dp.gold_nama, '') AS owner_name, TRUE AS visible").
		Joins("LEFT JOIN data_peserta dp ON dp.gold_id = outlet.outlet_gold_id").
		Joins("INNER JOIN buyer_visible_outlet bvo ON bvo.bvo_gold_id = outlet.outlet_gold_id AND bvo.bvo_outcode = outlet.outlet_code").
		Where("outlet.outlet_type <> ? AND outlet.outlet_status = ?", "THERAPY", "ACTIVE")

	if name != "" {
		query = query.Where("outlet.outlet_name LIKE ?", "%"+name+"%")
	}

	if err := query.Order("outlet.outlet_name asc").Find(&outlets).Error; err != nil {
		return nil, errors.Wrap(err, "[DATA][GetPublicOutlets]")
	}
	return outlets, nil
}

// GetAllOutletsForAdmin daftar SEMUA outlet non-THERAPY aktif + penanda apakah
// sudah dikurasi (visible) untuk layar pengaturan admin.
func (d *Data) GetAllOutletsForAdmin(ctx context.Context, name string) ([]orderEntity.PublicOutlet, error) {
	var outlets []orderEntity.PublicOutlet
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet.outlet_gold_id, outlet.outlet_code, outlet.outlet_name, outlet.outlet_type, outlet.outlet_address, COALESCE(dp.gold_nama, '') AS owner_name, (bvo.bvo_id IS NOT NULL) AS visible").
		Joins("LEFT JOIN data_peserta dp ON dp.gold_id = outlet.outlet_gold_id").
		Joins("LEFT JOIN buyer_visible_outlet bvo ON bvo.bvo_gold_id = outlet.outlet_gold_id AND bvo.bvo_outcode = outlet.outlet_code").
		Where("outlet.outlet_type <> ? AND outlet.outlet_status = ?", "THERAPY", "ACTIVE")

	if name != "" {
		query = query.Where("outlet.outlet_name LIKE ?", "%"+name+"%")
	}

	if err := query.Order("outlet.outlet_name asc").Find(&outlets).Error; err != nil {
		return nil, errors.Wrap(err, "[DATA][GetAllOutletsForAdmin]")
	}
	return outlets, nil
}

// AddVisibleOutlet menandai satu outlet boleh dilihat pembeli (idempoten).
func (d *Data) AddVisibleOutlet(ctx context.Context, goldid int, outcode, addedBy string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	row := orderEntity.VisibleOutlet{
		BvoGoldID:  goldid,
		BvoOutcode: outcode,
		BvoAddedBy: addedBy,
	}
	// ON DUPLICATE KEY UPDATE ringan supaya aman jika sudah ada (unique gold+outcode)
	return d.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "bvo_gold_id"}, {Name: "bvo_outcode"}},
			DoUpdates: clause.AssignmentColumns([]string{"bvo_added_by"}),
		}).
		Create(&row).Error
}

// RemoveVisibleOutlet mencabut satu outlet dari daftar yang boleh dilihat pembeli.
func (d *Data) RemoveVisibleOutlet(ctx context.Context, goldid int, outcode string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).
		Where("bvo_gold_id = ? AND bvo_outcode = ?", goldid, outcode).
		Delete(&orderEntity.VisibleOutlet{}).Error
}

// GetOutletCatalog mengambil barang (stock) satu outlet milik penjual tertentu,
// lengkap dengan harga & pack. Item brand THERAPY (jasa) dikecualikan.
func (d *Data) GetOutletCatalog(ctx context.Context, goldid int, outcode, name string) ([]orderEntity.CatalogItem, error) {
	var items []orderEntity.CatalogItem
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := d.db.WithContext(ctx).
		Table("stock").
		Select("stock.stock_id, stock.stock_name, stock.stock_pack, stock.stock_qty, COALESCE(items.item_price, 0) AS price, COALESCE(items.item_brand, '') AS brand").
		Joins("LEFT JOIN items ON items.item_id = stock.stock_item_id").
		Where("stock.stock_gold_id = ? AND stock.stock_outcode = ?", goldid, outcode).
		Where("COALESCE(items.item_brand, '') <> ?", "THERAPY")

	if name != "" {
		query = query.Where("stock.stock_name LIKE ?", "%"+name+"%")
	}

	if err := query.Order("stock.stock_name asc").Find(&items).Error; err != nil {
		return nil, errors.Wrap(err, "[DATA][GetOutletCatalog]")
	}
	return items, nil
}

// GetOutletByCode mengambil satu outlet berdasarkan gold_id pemilik + kode.
// outlet_code tidak unik global, jadi wajib disertai gold_id pemilik.
func (d *Data) GetOutletByCode(ctx context.Context, goldid int, outcode string) (*orderEntity.PublicOutlet, error) {
	var outlet orderEntity.PublicOutlet
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := d.db.WithContext(ctx).
		Table("outlet").
		Select("outlet.outlet_gold_id, outlet.outlet_code, outlet.outlet_name, outlet.outlet_type, outlet.outlet_address, COALESCE(dp.gold_nama, '') AS owner_name").
		Joins("LEFT JOIN data_peserta dp ON dp.gold_id = outlet.outlet_gold_id").
		Where("outlet.outlet_gold_id = ? AND outlet.outlet_code = ?", goldid, outcode).
		Take(&outlet).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetOutletByCode]")
	}
	return &outlet, nil
}

// GetBuyerName mengambil nama pembeli dari data_peserta.
func (d *Data) GetBuyerName(ctx context.Context, buyerID int) (string, error) {
	var name string
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Table("data_peserta").
		Select("gold_nama").
		Where("gold_id = ?", buyerID).
		Take(&name).Error
	if err != nil {
		return "", errors.Wrap(err, "[DATA][GetBuyerName]")
	}
	return name, nil
}

// InsertOrder menyimpan header + detail pesanan dalam satu transaksi.
func (d *Data) InsertOrder(ctx context.Context, header orderEntity.Order, details []orderEntity.OrderDetail) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&header).Error; err != nil {
			return errors.Wrap(err, "[DATA][InsertOrder][header]")
		}
		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return errors.Wrap(err, "[DATA][InsertOrder][detail]")
			}
		}
		return nil
	})
}

// GetOrdersByBuyer daftar pesanan milik satu pembeli (dashboard pembeli).
func (d *Data) GetOrdersByBuyer(ctx context.Context, buyerID int) ([]orderEntity.Order, error) {
	var orders []orderEntity.Order
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("order_buyer_id = ?", buyerID).
		Order("order_created_at desc").
		Find(&orders).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetOrdersByBuyer]")
	}
	return orders, nil
}

// GetOrdersBySeller daftar pesanan masuk untuk satu penjual (menu order penjual).
// status kosong = semua status.
func (d *Data) GetOrdersBySeller(ctx context.Context, goldid int, status string) ([]orderEntity.Order, error) {
	var orders []orderEntity.Order
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	query := d.db.WithContext(ctx).Where("order_gold_id = ?", goldid)
	if status != "" {
		query = query.Where("order_status = ?", status)
	}
	err := query.Order("order_created_at desc").Find(&orders).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetOrdersBySeller]")
	}
	return orders, nil
}

// GetOrderByID mengambil satu pesanan (header saja).
func (d *Data) GetOrderByID(ctx context.Context, orderID string) (*orderEntity.Order, error) {
	var o orderEntity.Order
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).Where("order_id = ?", orderID).Take(&o).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetOrderByID]")
	}
	return &o, nil
}

// GetOrderDetails mengambil item satu pesanan.
func (d *Data) GetOrderDetails(ctx context.Context, orderID string) ([]orderEntity.OrderDetail, error) {
	var details []orderEntity.OrderDetail
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	err := d.db.WithContext(ctx).
		Where("od_order_id = ?", orderID).
		Order("od_id asc").
		Find(&details).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetOrderDetails]")
	}
	return details, nil
}

// UpdateOrderStatus mengubah status (opsional alasan penolakan).
func (d *Data) UpdateOrderStatus(ctx context.Context, orderID, status string, reason *string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	updates := map[string]interface{}{"order_status": status}
	if reason != nil {
		updates["order_reject_reason"] = *reason
	}
	return d.db.WithContext(ctx).
		Model(&orderEntity.Order{}).
		Where("order_id = ?", orderID).
		Updates(updates).Error
}

// FinishOrder menandai pesanan FINISH dan menautkan sale_id nota yang dibuat.
func (d *Data) FinishOrder(ctx context.Context, orderID, saleID string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return d.db.WithContext(ctx).
		Model(&orderEntity.Order{}).
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{
			"order_status":  orderEntity.StatusFinish,
			"order_sale_id": saleID,
		}).Error
}
