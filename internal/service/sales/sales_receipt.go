package goldgym

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// formatRupiah mengubah decimal menjadi format "12.345" (separator ribuan titik)
func formatRupiah(d *decimal.Decimal) string {
	if d == nil {
		return "0"
	}
	raw := d.StringFixed(0)
	neg := strings.HasPrefix(raw, "-")
	if neg {
		raw = raw[1:]
	}
	var sb strings.Builder
	for i, c := range raw {
		if i > 0 && (len(raw)-i)%3 == 0 {
			sb.WriteByte('.')
		}
		sb.WriteRune(c)
	}
	if neg {
		return "-" + sb.String()
	}
	return sb.String()
}

// GetSaleReceiptPDF membuat nota PDF (ukuran struk 58mm) dari th_sale + td_sale.
// Return: bytes PDF, nomor nota (untuk nama file), error.
func (s Service) GetSaleReceiptPDF(ctx context.Context, goldid int, custid int, saleid string) ([]byte, string, error) {
	sale, err := s.GetSaleDetail(ctx, goldid, custid, saleid)
	if err != nil {
		return nil, "", err
	}

	outletName := sale.Header.SaleOutcode
	outletAddress := ""
	outletType := ""
	// info outlet memakai gold_id di header supaya valid juga saat diakses pembeli
	outlet, err := s.goldgymsale.GetSaleOutletInfo(ctx, sale.Header.SaleGoldID, sale.Header.SaleOutcode)
	if err == nil && outlet != nil {
		outletType = outlet.OutletType
		if outlet.OutletName != "" {
			outletName = outlet.OutletName
			outletAddress = outlet.OutletAddress
		}
	}

	// knob global untuk menaik/menurunkan SEMUA ukuran font sekaligus (0 = default).
	// Dideklarasikan paling atas supaya pageHeight & closure divider/row memakainya.
	fontSize := 10.0

	// pageLn := 0.0

	// Tinggi kertas dihitung dinamis (struk thermal = 1 halaman menerus):
	// bagian tetap ~95mm + tiap item ~11mm, plus tambahan mengikuti fontSize
	// (huruf lebih besar => tiap baris lebih tinggi). Angka disetel dari
	// pengukuran GetY dengan cadangan ~7-10mm supaya tidak pernah pecah ke
	// halaman ke-2 walau font diperbesar/nama item panjang.
	discountLineCount := 0
	for _, d := range sale.Detail {
		if d.SaleDiscountValue != nil && d.SaleDiscountAmount != nil {
			discountLineCount++
		}
	}
	extraTotalLines := 0.0
	if sale.Header.SaleTotalDiscountAmount != nil {
		extraTotalLines++
	}
	if sale.Header.SaleVoucherAmount != nil {
		extraTotalLines++
	}
	pageHeight := 95.0 + 54.0 + float64(len(sale.Detail))*16.0 + float64(discountLineCount)*5.0 + extraTotalLines*6.0
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size:    gofpdf.SizeType{Wd: 125, Ht: pageHeight},
	})
	pdf.SetMargins(5, 2, 0)
	pdf.SetAutoPageBreak(true, 2)
	pdf.AddPage()

	lineWidth := 120.0

	divider := func() {
		pdf.SetFont("Courier", "B", 8+fontSize)
		pdf.CellFormat(lineWidth, 4.5, strings.Repeat("-", 45), "", 1, "C", false, 0, "")
	}

	// Lebar kolom label mengikuti ukuran font baris (8+fontSize) secara dinamis:
	// diukur dari label terpanjang "Customer :" supaya titik dua tetap sejajar dan
	// label tidak terpotong berapa pun fontSize-nya. +1.5mm buffer kecil.
	pdf.SetFont("Courier", "", 8+fontSize)
	labelWidth := pdf.GetStringWidth("Customer :") + 1.5
	rowH := 5.0 + fontSize*0.3 // tinggi baris ikut naik sedikit mengikuti besar font

	// value bisa lebih panjang dari lebar kolom pada ukuran huruf yang besar
	// (nomor nota, nama pelanggan dll) — pakai MultiCell supaya turun baris
	// otomatis, bukan terpotong di tepi kertas.
	// ": " sengaja digabung ke label (CellFormat), bukan ikut dikirim ke
	// MultiCell: MultiCell membungkus teks di spasi terakhir yang ditemukan,
	// jadi kalau value tidak punya spasi sama sekali (mis. nomor nota),
	// ": " ikut terpisah sendirian ke baris atas dan value turun penuh ke
	// baris berikutnya (titik dua jadi tidak sejajar dengan isinya).
	row := func(label, value string) {
		pdf.SetFont("Courier", "", 8+fontSize)
		pdf.CellFormat(labelWidth, rowH, label+" :", "", 0, "L", false, 0, "")
		pdf.MultiCell(lineWidth-labelWidth, rowH, value, "", "L", false)
	}

	// Kop nota — MultiCell supaya nama outlet panjang turun baris, bukan kepotong
	pdf.SetFont("Courier", "", 12+fontSize)
	pdf.MultiCell(lineWidth, 6, outletName, "", "C", false)
	// pageLn = pageLn + 2.0
	pdf.Ln(2)
	if outletAddress != "" {
		pdf.SetFont("Courier", "", 8+fontSize)
		pdf.MultiCell(lineWidth, 4.5, outletAddress, "", "C", false)
	}
	pdf.Ln(2)
	divider()
	pdf.Ln(2)

	transDate := ""
	if sale.Header.SaleTransdate != nil {
		transDate = sale.Header.SaleTransdate.Format("02-01-2006")
	}
	transTime := ""
	if sale.Header.SaleTranstime != nil {
		transTime = strings.TrimSpace(*sale.Header.SaleTranstime)
		// sale_transtime disimpan "HHMMSS" (mis. "143025") → tampilkan
		// jam:menit:detik. Fallback ke nilai apa adanya kalau formatnya lain.
		if t, e := time.Parse("150405", transTime); e == nil {
			transTime = t.Format("15:04:05")
		} else if t, e := time.Parse("1504", transTime); e == nil {
			transTime = t.Format("15:04:05")
		}
	}
	row("No Nota", sale.Header.SaleTrancnum)
	row("Tanggal", strings.TrimSpace(transDate+" "+transTime))
	if sale.Header.SaleSalesperson != nil && *sale.Header.SaleSalesperson != "" {
		row("Kasir", *sale.Header.SaleSalesperson)
	}
	// Customer di bawah Kasir; untuk THERAPY dihormati toggle tampil-customer
	// (sale_customer_show='N' => nama customer tidak dicetak).
	if sale.Header.SaleSalescustomer != nil && *sale.Header.SaleSalescustomer != "" &&
		!(strings.EqualFold(outletType, "THERAPY") && sale.Header.SaleCustomerShow == "N") {
		row("Customer", *sale.Header.SaleSalescustomer)
	}
	// nama toko hanya ada pada pembeli yang didaftarkan penjual (gold_toko)
	if sale.Header.SaleCustID != nil && *sale.Header.SaleCustID > 0 {
		toko, err := s.goldgymsale.GetSaleCustomerToko(ctx, *sale.Header.SaleCustID)
		if err == nil && toko != "" {
			row("Toko Customer", toko)
		}
	}
	// meja (fitur Atur Meja, retail) -- sudah tersimpan langsung di
	// th_sale.sale_meja_names saat insert, tidak perlu query tambahan lagi.
	if sale.Header.SaleMejaNames != nil && *sale.Header.SaleMejaNames != "" {
		row("No Meja", *sale.Header.SaleMejaNames)
	}
	pdf.Ln(2)
	divider()
	pdf.Ln(2)

	// Daftar item
	for x, d := range sale.Detail {
		name := ""
		if d.SaleStockname != nil {
			name = *d.SaleStockname
		}
		// hilangkan anotasi booking "... (Booking 21-07 05:00 - nama)" dari nama
		// item di nota (nama customer-nya tetap tampil di baris Customer/Toko)
		if i := strings.Index(name, " (Booking "); i >= 0 {
			name = strings.TrimSpace(name[:i])
		}
		qty := 0
		if d.SaleQty != nil {
			qty = *d.SaleQty
		}
		if x == 0 {
			pdf.Ln(2)
		}
		if x > 0 {
			// pageLn = pageLn + 4.0
			pdf.Ln(3)
		}
		pdf.SetFont("Courier", "", 8+fontSize)
		pdf.MultiCell(lineWidth, 6, name, "", "L", false)
		qtyPrice := ""
		if d.SaleSalesprice != nil {
			qtyPrice = formatRupiah(d.SaleSalesprice)
		}
		// pageLn = pageLn + 2.0
		pdf.Ln(2)
		pdf.CellFormat(28, 4.5, "  "+strconv.Itoa(qty)+" x "+qtyPrice, "", 0, "L", false, 0, "")
		pdf.CellFormat(lineWidth-28, 4.5, formatRupiah(d.SaleTotalsalesprice), "", 1, "R", false, 0, "")

		if d.SaleDiscountValue != nil && d.SaleDiscountAmount != nil {
			label := "  Diskon"
			if d.SaleDiscountType != nil && *d.SaleDiscountType == "PERCENT" {
				label = "  Diskon (" + d.SaleDiscountValue.StringFixed(0) + "%)"
			}
			pdf.SetFont("Courier", "", 7+fontSize)
			pdf.CellFormat(28, 4.5, label, "", 0, "L", false, 0, "")
			pdf.CellFormat(lineWidth-28, 4.5, "-"+formatRupiah(d.SaleDiscountAmount), "", 1, "R", false, 0, "")
		}
	}
	pdf.Ln(2)
	divider()
	pdf.Ln(2)

	// Total pembayaran
	totalRow := func(label string, value *decimal.Decimal, bold bool) {
		style := ""
		if bold {
			style = "B"
		}
		pdf.SetFont("Courier", style, 10+fontSize)
		pdf.CellFormat(22, 5.5, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(lineWidth-22, 5.5, formatRupiah(value), "", 1, "R", false, 0, "")
	}
	if sale.Header.SaleTotalDiscountAmount != nil {
		label := "Diskon Total"
		if sale.Header.SaleTotalDiscountPercent != nil {
			label = "Diskon Total (" + sale.Header.SaleTotalDiscountPercent.StringFixed(0) + "%)"
		}
		pdf.SetFont("Courier", "", 8+fontSize)
		pdf.CellFormat(lineWidth-22, 4.5, label, "", 0, "L", false, 0, "")
		amt := *sale.Header.SaleTotalDiscountAmount
		pdf.CellFormat(22, 4.5, "-"+formatRupiah(&amt), "", 1, "R", false, 0, "")
		pdf.Ln(1)
	}
	if sale.Header.SaleVoucherAmount != nil {
		code := ""
		if sale.Header.SaleVoucherCode != nil {
			code = *sale.Header.SaleVoucherCode
		}
		label := "Voucher " + code
		if sale.Header.SaleVoucherPercent != nil {
			label += " (" + sale.Header.SaleVoucherPercent.StringFixed(0) + "%)"
		}
		pdf.SetFont("Courier", "", 8+fontSize)
		pdf.CellFormat(lineWidth-22, 4.5, label, "", 0, "L", false, 0, "")
		amt := *sale.Header.SaleVoucherAmount
		pdf.CellFormat(22, 4.5, "-"+formatRupiah(&amt), "", 1, "R", false, 0, "")
		pdf.Ln(1)
	}
	totalRow("TOTAL", sale.Header.SaleTranstotal, true)
	pdf.Ln(2)
	totalRow("BAYAR", sale.Header.SaleTranspayment, false)
	pdf.Ln(2)
	totalRow("KEMBALI", sale.Header.SaleTranschange, false)
	pdf.Ln(6)

	status := "BELUM LUNAS"
	if sale.Header.SalePaymentyn == "Y" {
		status = "LUNAS"
	}
	pdf.SetFont("Courier", "B", 10+fontSize)
	pdf.CellFormat(lineWidth, 5.5, "STATUS: "+status, "", 1, "C", false, 0, "")
	divider()

	pdf.Ln(2)
	pdf.SetFont("Courier", "", 8+fontSize)
	pdf.MultiCell(lineWidth, 4.5, "Terima kasih atas kunjungan Anda", "", "C", false)
	pdf.Ln(6)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", errors.Wrap(err, "[Service][GetSaleReceiptPDF][Output]")
	}

	return buf.Bytes(), sale.Header.SaleTrancnum, nil
}
