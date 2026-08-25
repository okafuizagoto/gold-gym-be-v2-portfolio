package goldgym

import (
	"context"
	"fmt"
	"time"

	"gold-gym-be/internal/entity"
	goldSaleEntity "gold-gym-be/internal/entity/sales"

	"github.com/pkg/errors"
)

var monthNamesID = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

var monthShortID = []string{
	"", "Jan", "Feb", "Mar", "Apr", "Mei", "Jun",
	"Jul", "Agu", "Sep", "Okt", "Nov", "Des",
}

// weekNoOf mengembalikan nomor blok minggu (1-based) dari tanggal-dalam-bulan,
// dengan aturan blok 7 hari dari tgl 1: 1–7=minggu 1, 8–14=minggu 2, dst.
func weekNoOf(day int) int { return (day-1)/7 + 1 }

// GetSalesReport membangun laporan penjualan sesuai mode (day/week/month),
// di-scope ke outlet penjual (goldid + outcode). Semua total dijumlahkan dari
// sale_transtotal (day) / subtotal item (day grand total).
func (s Service) GetSalesReport(ctx context.Context, goldid int, outcode, mode, date string) (interface{}, error) {
	if outcode == "" {
		return nil, errors.Wrap(entity.ErrInvalid, "outcode wajib diisi")
	}
	switch mode {
	case "day":
		return s.reportDay(ctx, goldid, outcode, date)
	case "week":
		return s.reportWeek(ctx, goldid, outcode, date)
	case "month":
		return s.reportMonth(ctx, goldid, outcode, date)
	default:
		return nil, errors.Wrap(entity.ErrInvalid, "mode harus day/week/month")
	}
}

func (s Service) reportDay(ctx context.Context, goldid int, outcode, date string) (goldSaleEntity.SaleReportDay, error) {
	var res goldSaleEntity.SaleReportDay
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return res, errors.Wrap(entity.ErrInvalid, "tanggal harus format YYYY-MM-DD")
	}
	dateStr := d.Format("2006-01-02")

	items, err := s.goldgymsale.GetSaleReportItems(ctx, goldid, outcode, dateStr)
	if err != nil {
		return res, errors.Wrap(err, "[Service][reportDay]")
	}

	var grand float64
	notas := map[string]struct{}{}
	for _, it := range items {
		grand += it.Subtotal
		notas[it.SaleID] = struct{}{}
	}

	res = goldSaleEntity.SaleReportDay{
		Mode:       "day",
		Date:       dateStr,
		Items:      items,
		Count:      len(notas),
		GrandTotal: grand,
	}
	return res, nil
}

func (s Service) reportWeek(ctx context.Context, goldid int, outcode, date string) (goldSaleEntity.SaleReportWeek, error) {
	var res goldSaleEntity.SaleReportWeek
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return res, errors.Wrap(entity.ErrInvalid, "tanggal harus format YYYY-MM-DD")
	}

	year, month := d.Year(), d.Month()
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	wk := weekNoOf(d.Day())
	startDay := (wk-1)*7 + 1
	endDay := wk * 7
	if endDay > daysInMonth {
		endDay = daysInMonth
	}
	start := time.Date(year, month, startDay, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, month, endDay, 0, 0, 0, 0, time.UTC)
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")

	totals, err := s.goldgymsale.GetSaleDailyTotals(ctx, goldid, outcode, startStr, endStr)
	if err != nil {
		return res, errors.Wrap(err, "[Service][reportWeek]")
	}
	byDate := map[string]goldSaleEntity.SaleDailyTotal{}
	for _, t := range totals {
		byDate[t.Date] = t
	}

	// lengkapi semua hari di blok minggu meski tidak ada transaksi (total 0)
	var days []goldSaleEntity.SaleDailyTotal
	var grand float64
	for day := startDay; day <= endDay; day++ {
		ds := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		if t, ok := byDate[ds]; ok {
			days = append(days, t)
			grand += t.Total
		} else {
			days = append(days, goldSaleEntity.SaleDailyTotal{Date: ds, Total: 0, Count: 0})
		}
	}

	res = goldSaleEntity.SaleReportWeek{
		Mode:       "week",
		Label:      fmt.Sprintf("%d–%d %s %d", startDay, endDay, monthShortID[int(month)], year),
		RangeStart: startStr,
		RangeEnd:   endStr,
		Days:       days,
		GrandTotal: grand,
	}
	return res, nil
}

func (s Service) reportMonth(ctx context.Context, goldid int, outcode, date string) (goldSaleEntity.SaleReportMonth, error) {
	var res goldSaleEntity.SaleReportMonth
	// terima YYYY-MM maupun YYYY-MM-DD
	layout := "2006-01-02"
	if len(date) == 7 {
		layout = "2006-01"
	}
	d, err := time.Parse(layout, date)
	if err != nil {
		return res, errors.Wrap(entity.ErrInvalid, "bulan harus format YYYY-MM")
	}

	year, month := d.Year(), d.Month()
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	firstStr := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	lastStr := time.Date(year, month, daysInMonth, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	totals, err := s.goldgymsale.GetSaleDailyTotals(ctx, goldid, outcode, firstStr, lastStr)
	if err != nil {
		return res, errors.Wrap(err, "[Service][reportMonth]")
	}

	weekCount := weekNoOf(daysInMonth) // jumlah blok minggu di bulan ini
	weeks := make([]goldSaleEntity.SaleWeeklyTotal, weekCount)
	for i := 0; i < weekCount; i++ {
		wk := i + 1
		startDay := (wk-1)*7 + 1
		endDay := wk * 7
		if endDay > daysInMonth {
			endDay = daysInMonth
		}
		weeks[i] = goldSaleEntity.SaleWeeklyTotal{
			WeekNo:     wk,
			Label:      fmt.Sprintf("%d–%d", startDay, endDay),
			RangeStart: time.Date(year, month, startDay, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
			RangeEnd:   time.Date(year, month, endDay, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
		}
	}

	var grand float64
	for _, t := range totals {
		td, perr := time.Parse("2006-01-02", t.Date)
		if perr != nil {
			continue
		}
		idx := weekNoOf(td.Day()) - 1
		if idx < 0 || idx >= len(weeks) {
			continue
		}
		weeks[idx].Total += t.Total
		weeks[idx].Count += t.Count
		grand += t.Total
	}

	res = goldSaleEntity.SaleReportMonth{
		Mode:       "month",
		Month:      d.Format("2006-01"),
		Label:      fmt.Sprintf("%s %d", monthNamesID[int(month)], year),
		Weeks:      weeks,
		GrandTotal: grand,
	}
	return res, nil
}
