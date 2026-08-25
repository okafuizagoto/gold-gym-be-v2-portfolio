package goldgym

import (
	"context"

	goldSaleEntity "gold-gym-be/internal/entity/sales"

	"github.com/pkg/errors"
)

// GetSaleReportItems mengambil baris item penjualan satu hari untuk outlet
// (goldid+outcode). Sisa stok (remaining) diambil live dari tabel stock.
// Diurutkan per nama customer supaya baris customer sama bersebelahan (FE
// menggabungkan sel customer yang sama tanpa garis pemisah).
func (d *Data) GetSaleReportItems(ctx context.Context, goldid int, outcode, date string) ([]goldSaleEntity.SaleReportItem, error) {
	var rows []goldSaleEntity.SaleReportItem
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := d.db.WithContext(ctx).
		Table("td_sale AS td").
		Select(`th.sale_id AS sale_id,
			th.sale_trancnum AS trancnum,
			COALESCE(th.sale_salescustomer, '') AS customer,
			COALESCE(th.sale_salesperson, '') AS salesperson,
			COALESCE(NULLIF(th.sale_transtime, ''), DATE_FORMAT(th.sale_created_at, '%H:%i')) AS trans_time,
			COALESCE(td.sale_stockname, '') AS item_name,
			COALESCE(td.sale_qty, 0) AS qty,
			COALESCE(td.sale_salesprice, 0) AS price,
			COALESCE(td.sale_totalsalesprice, 0) AS subtotal,
			COALESCE(s.stock_qty, 0) AS remaining`).
		Joins("JOIN th_sale AS th ON th.sale_id = td.sale_id").
		Joins("LEFT JOIN stock AS s ON s.stock_id = td.sale_stockid AND s.stock_gold_id = th.sale_gold_id").
		Where("th.sale_gold_id = ? AND th.sale_outcode = ?", goldid, outcode).
		Where("DATE(COALESCE(th.sale_transdate, th.sale_created_at)) = ?", date).
		Order("customer ASC, th.sale_created_at ASC, td.td_id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetSaleReportItems]")
	}
	return rows, nil
}

// GetSaleDailyTotals menjumlahkan total penjualan per hari (jumlah nota + total
// rupiah) untuk rentang tanggal [start, end] pada satu outlet. Hari tanpa
// transaksi tidak muncul di sini (dilengkapi di layer service).
func (d *Data) GetSaleDailyTotals(ctx context.Context, goldid int, outcode, start, end string) ([]goldSaleEntity.SaleDailyTotal, error) {
	var rows []goldSaleEntity.SaleDailyTotal
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := d.db.WithContext(ctx).
		Table("th_sale AS th").
		Select(`DATE_FORMAT(COALESCE(th.sale_transdate, th.sale_created_at), '%Y-%m-%d') AS date,
			COALESCE(SUM(th.sale_transtotal), 0) AS total,
			COUNT(*) AS count`).
		Where("th.sale_gold_id = ? AND th.sale_outcode = ?", goldid, outcode).
		Where("DATE(COALESCE(th.sale_transdate, th.sale_created_at)) BETWEEN ? AND ?", start, end).
		Group("DATE(COALESCE(th.sale_transdate, th.sale_created_at))").
		Order("date ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, errors.Wrap(err, "[DATA][GetSaleDailyTotals]")
	}
	return rows, nil
}
