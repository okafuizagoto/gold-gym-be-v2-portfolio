package booking

import (
	"context"
	"strconv"
	"strings"
	"time"

	bookingEntity "gold-gym-be/internal/entity/booking"
	saleEntity "gold-gym-be/internal/entity/sales"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// normalizeTherapyType menerima input FE ("SOFA", "DRAGON"/"KURSI DRAGON", "KURSI")
// dan mengembalikan konstanta tipe terapi; string kosong jika tidak valid.
// "KURSI DRAGON" dicek sebelum "KURSI" agar tidak salah tangkap.
func normalizeTherapyType(raw string) string {
	switch {
	case raw == bookingEntity.TherapyTypeSofa || raw == "sofa" || raw == "Sofa":
		return bookingEntity.TherapyTypeSofa
	case raw == bookingEntity.TherapyTypeDragon || raw == "dragon" || raw == "Dragon" ||
		raw == "KURSI DRAGON" || raw == "kursi dragon" || raw == "Kursi Dragon":
		return bookingEntity.TherapyTypeDragon
	case raw == bookingEntity.TherapyTypeKursi || raw == "kursi" || raw == "Kursi":
		return bookingEntity.TherapyTypeKursi
	}
	return ""
}

// canOverridePrice hanya penjual/admin yang boleh memakai harga input sendiri.
func canOverridePrice(role string) bool {
	return role == "SELLER" || role == "ADMIN"
}

// GetSlots mengembalikan grid slot 30 menit untuk satu outlet & tanggal
// (jam 05:00-17:00 WIB, kapasitas 6 orang per jendela 30 menit).
// Booking UNPAID yang jamnya lewat otomatis di-EXPIRED-kan dulu sehingga
// slotnya kembali available (hijau).
func (s Service) GetSlots(ctx context.Context, outcode string, dateStr string) ([]bookingEntity.Slot, error) {
	date, err := time.ParseInLocation("2006-01-02", dateStr, bookingEntity.Jakarta)
	if err != nil {
		return nil, errors.New("format tanggal harus YYYY-MM-DD")
	}

	outlet, err := s.booking.GetOutletByCode(ctx, outcode)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetSlots][GetOutletByCode]")
	}
	if outlet == nil {
		return nil, errors.New("outlet tidak ditemukan")
	}
	goldid := outlet.OutletGoldID

	now := time.Now().In(bookingEntity.Jakarta)
	if err := s.booking.ExpireOverdueBookings(ctx, now); err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetSlots][ExpireOverdueBookings]")
	}

	bookings, err := s.booking.GetBookingsByDate(ctx, goldid, outcode, date)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetSlots][GetBookingsByDate]")
	}

	var slots []bookingEntity.Slot
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), bookingEntity.OpenHour, 0, 0, 0, bookingEntity.Jakarta)
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(), bookingEntity.CloseHour, 0, 0, 0, bookingEntity.Jakarta)

	for win := dayStart; win.Before(dayEnd); win = win.Add(30 * time.Minute) {
		slot := bookingEntity.Slot{
			Start:    win.Format("15:04"),
			End:      win.Add(30 * time.Minute).Format("15:04"),
			Capacity: bookingEntity.CapacityPerHalfHour,
			Past:     win.Before(now),
			Bookings: []bookingEntity.SlotBooking{},
		}
		for _, b := range bookings {
			bStart := b.BookingStart.In(bookingEntity.Jakarta)
			bEnd := bStart.Add(time.Duration(b.BookingDuration) * time.Minute)
			// booking menempati jendela jika ada overlap [win, win+30m)
			if !bStart.After(win) && bEnd.After(win) {
				slot.Used++
				if b.BookingStatus == bookingEntity.StatusUnpaid {
					slot.HasUnpaid = true
				}
				ttype := ""
				if b.BookingTherapyType != nil {
					ttype = *b.BookingTherapyType
				}
				priceStr := "0"
				if b.BookingPrice != nil && b.BookingPrice.IsPositive() {
					priceStr = strconv.FormatInt(b.BookingPrice.IntPart(), 10)
				}
				slot.Bookings = append(slot.Bookings, bookingEntity.SlotBooking{
					BookingID:    b.BookingID,
					CustName:     b.BookingCustName,
					RegisteredYN: b.BookingRegisteredYN,
					Status:       b.BookingStatus,
					Duration:     b.BookingDuration,
					Start:        bStart.Format("15:04"),
					TherapyType:  ttype,
					Price:        priceStr,
				})
			}
		}
		slot.Available = slot.Capacity - slot.Used
		if slot.Available < 0 {
			slot.Available = 0
		}
		slot.Full = slot.Used >= slot.Capacity
		slots = append(slots, slot)
	}

	return slots, nil
}

// InsertBooking membuat booking baru. Untuk booking berbayar, mengembalikan juga
// payload insert_sales yang harus dipublish handler ke Kafka (arsitektur sama
// dengan insert sales biasa) — sale_id sudah dibuat di sini dan tersimpan di booking.
// Jika ada booking 30 menit UNPAID bersebelahan dengan nama & tipe sama, pembayaran
// digabung menjadi 1 jam (harga 1 jam, satu nota untuk kedua booking).
func (s Service) InsertBooking(ctx context.Context, userID int, role string, creator string, req bookingEntity.InsertBookingRequest) (bookingEntity.Booking, *saleEntity.InsertSaleData, error) {
	var b bookingEntity.Booking

	if req.Duration != 30 && req.Duration != 60 {
		return b, nil, errors.New("durasi harus 30 atau 60 menit")
	}
	if req.Outcode == "" {
		return b, nil, errors.New("outcode wajib diisi")
	}
	ttype := normalizeTherapyType(req.TherapyType)
	if ttype == "" {
		return b, nil, errors.New("tipe terapi harus SOFA, KURSI DRAGON, atau KURSI")
	}
	if req.CustomPrice < 0 {
		return b, nil, errors.New("harga custom tidak boleh negatif")
	}

	outlet, err := s.booking.GetOutletByCode(ctx, req.Outcode)
	if err != nil {
		return b, nil, errors.Wrap(err, "[SERVICE][InsertBooking][GetOutletByCode]")
	}
	if outlet == nil {
		return b, nil, errors.New("outlet tidak ditemukan")
	}
	if outlet.OutletType != "THERAPY" {
		return b, nil, errors.New("booking hanya untuk outlet bertipe THERAPY")
	}

	date, err := time.ParseInLocation("2006-01-02", req.Date, bookingEntity.Jakarta)
	if err != nil {
		return b, nil, errors.New("format tanggal harus YYYY-MM-DD")
	}
	startClock, err := time.Parse("15:04", req.Start)
	if err != nil {
		return b, nil, errors.New("format jam harus HH:MM")
	}
	if startClock.Minute()%30 != 0 {
		return b, nil, errors.New("jam booking harus kelipatan 30 menit")
	}
	start := time.Date(date.Year(), date.Month(), date.Day(), startClock.Hour(), startClock.Minute(), 0, 0, bookingEntity.Jakarta)

	now := time.Now().In(bookingEntity.Jakarta)
	if start.Before(now) {
		return b, nil, errors.New("tidak bisa booking di waktu yang sudah lewat")
	}
	if startClock.Hour() < bookingEntity.OpenHour ||
		start.Add(time.Duration(req.Duration)*time.Minute).After(time.Date(date.Year(), date.Month(), date.Day(), bookingEntity.CloseHour, 0, 0, 0, bookingEntity.Jakarta)) {
		return b, nil, errors.New("di luar jam operasional booking (05:00 - 17:00 WIB)")
	}

	// pembeli login hanya bisa booking atas nama sendiri
	custID := req.CustID
	custName := strings.TrimSpace(req.CustName)
	registered := "N"
	if role == "BUYER" {
		custID = userID
	}
	if custID < 0 {
		return b, nil, errors.New("cust_id tidak valid")
	}
	if custID > 0 {
		buyer, err := s.booking.GetBuyerByID(ctx, custID)
		if err != nil {
			return b, nil, errors.Wrap(err, "[SERVICE][InsertBooking][GetBuyerByID]")
		}
		if buyer == nil {
			return b, nil, errors.New("pembeli terdaftar tidak ditemukan")
		}
		registered = "Y"
		if custName == "" {
			custName = buyer.GoldNama
		}
	}
	if custName == "" {
		return b, nil, errors.New("nama pembeli wajib diisi")
	}
	// kolom booking_cust_name varchar(100)
	if len(custName) > 100 {
		return b, nil, errors.New("nama pembeli maksimal 100 karakter")
	}

	// lazy-expire dulu supaya slot bekas UNPAID kadaluarsa terhitung kosong
	if err := s.booking.ExpireOverdueBookings(ctx, now); err != nil {
		return b, nil, errors.Wrap(err, "[SERVICE][InsertBooking][ExpireOverdueBookings]")
	}

	// validasi kapasitas: setiap jendela 30 menit yang dipakai maksimal 6 orang (total semua tipe)
	for win := start; win.Before(start.Add(time.Duration(req.Duration) * time.Minute)); win = win.Add(30 * time.Minute) {
		count, err := s.booking.CountActiveInWindow(ctx, outlet.OutletGoldID, req.Outcode, win)
		if err != nil {
			return b, nil, errors.Wrap(err, "[SERVICE][InsertBooking][CountActiveInWindow]")
		}
		if count >= bookingEntity.CapacityPerHalfHour {
			return b, nil, errors.New("slot " + win.Format("15:04") + " sudah penuh")
		}
	}

	status := bookingEntity.StatusUnpaid
	var itemID *int
	var itemName *string
	var price *decimal.Decimal
	var salePayload *saleEntity.InsertSaleData
	var saleID *string
	var adjacent *bookingEntity.Booking

	b = bookingEntity.Booking{
		BookingID:           uuid.New().String(),
		BookingGoldID:       outlet.OutletGoldID,
		BookingOutcode:      req.Outcode,
		BookingDate:         date,
		BookingStart:        start,
		BookingDuration:     req.Duration,
		BookingTherapyType:  &ttype,
		BookingCustName:     custName,
		BookingRegisteredYN: registered,
	}
	if custID > 0 {
		b.BookingCustID = &custID
	}
	b.BookingCreatedBy = &creator
	roleCopy := role
	b.BookingCreatedRole = &roleCopy

	if req.Paid {
		// gabungan 1 jam: booking 30 menit UNPAID bersebelahan, nama & tipe sama
		itemDuration := req.Duration
		saleStart := start
		if req.Duration == 30 {
			adjacent, err = s.booking.FindAdjacentUnpaid(ctx, outlet.OutletGoldID, req.Outcode, custName, ttype, b.BookingID, start)
			if err != nil {
				return b, nil, errors.Wrap(err, "[SERVICE][InsertBooking][FindAdjacentUnpaid]")
			}
			if adjacent != nil {
				itemDuration = 60
				adjStart := adjacent.BookingStart.In(bookingEntity.Jakarta)
				if adjStart.Before(saleStart) {
					saleStart = adjStart
				}
			}
		}

		item, err := s.booking.GetTherapyItemByName(ctx, req.Outcode, bookingEntity.TherapyItemName(ttype, itemDuration))
		if err != nil {
			return b, nil, errors.Wrap(err, "[SERVICE][InsertBooking][GetTherapyItemByName]")
		}
		if item == nil {
			return b, nil, errors.New("item terapi " + bookingEntity.TherapyItemName(ttype, itemDuration) + " belum terdaftar di outlet ini")
		}

		finalPrice := item.ItemPrice
		if req.CustomPrice > 0 && canOverridePrice(role) {
			finalPrice = req.CustomPrice
		}

		status = bookingEntity.StatusPaid
		itemID = &item.ItemID
		itemName = &item.ItemName
		p := decimal.NewFromInt(int64(finalPrice))
		price = &p

		label := item.ItemName
		if adjacent != nil {
			label += " (GABUNGAN 2x30 MENIT)"
		}
		payload := s.buildSalePayload(outlet.OutletGoldID, req.Outcode, custName, creator, label, item.ItemPack,
			s.resolveStockID(ctx, req.Outcode, item.ItemID), finalPrice, custID, saleStart)
		salePayload = &payload
		sid := payload.InsertData.Header.SaleID
		saleID = &sid
	}

	// penjual/admin bisa mengunci harga custom sejak booking dibuat (belum bayar);
	// harga tersimpan ini otomatis dipakai saat pembayaran nanti
	if !req.Paid && req.CustomPrice > 0 && canOverridePrice(role) {
		p := decimal.NewFromInt(int64(req.CustomPrice))
		price = &p
	}

	// kolom booking_price — booking belum bayar disimpan dengan harga 0,
	// bukan NULL (NULL eksplisit pernah ditolak database: "booking_price cannot be null")
	if price == nil {
		zero := decimal.Zero
		price = &zero
	}

	b.BookingStatus = status
	b.BookingItemID = itemID
	b.BookingItemName = itemName
	b.BookingPrice = price
	b.BookingSaleID = saleID
	b.BookingCreatedAt = now

	if err := s.booking.InsertBooking(ctx, b); err != nil {
		return b, nil, errors.Wrap(err, "[SERVICE][InsertBooking][InsertBooking]")
	}

	// booking bersebelahan ikut lunas lewat nota gabungan yang sama (harga di baris gabungan = 0)
	if adjacent != nil && saleID != nil && itemID != nil && itemName != nil {
		if err := s.booking.UpdateBookingPaid(ctx, adjacent.BookingID, *saleID, *itemID, *itemName+" (GABUNGAN)", "0"); err != nil {
			return b, nil, errors.Wrap(err, "[SERVICE][InsertBooking][UpdateAdjacentPaid]")
		}
	}

	return b, salePayload, nil
}

// PayBooking menandai booking UNPAID menjadi PAID dan mengembalikan payload
// insert_sales untuk dipublish handler ke Kafka. Jika ada booking 30 menit UNPAID
// bersebelahan dengan nama & tipe sama, keduanya digabung jadi satu nota 1 jam
// dengan harga 1 jam. customPrice > 0 (khusus SELLER/ADMIN) menggantikan harga default.
func (s Service) PayBooking(ctx context.Context, creator string, role string, bookingID string, itemID int, customPrice int) (bookingEntity.Booking, *saleEntity.InsertSaleData, error) {
	var b bookingEntity.Booking

	if customPrice < 0 {
		return b, nil, errors.New("harga custom tidak boleh negatif")
	}

	existing, err := s.booking.GetBookingByID(ctx, bookingID)
	if err != nil {
		return b, nil, errors.Wrap(err, "[SERVICE][PayBooking][GetBookingByID]")
	}
	if existing == nil {
		return b, nil, errors.New("booking tidak ditemukan")
	}
	if existing.BookingStatus == bookingEntity.StatusPaid {
		return b, nil, errors.New("booking sudah lunas")
	}
	if existing.BookingStatus != bookingEntity.StatusUnpaid {
		return b, nil, errors.New("booking sudah " + existing.BookingStatus + ", tidak bisa dibayar")
	}

	var item *bookingEntity.TherapyItem
	var adjacent *bookingEntity.Booking
	saleStart := existing.BookingStart.In(bookingEntity.Jakarta)

	if existing.BookingTherapyType != nil && *existing.BookingTherapyType != "" {
		ttype := *existing.BookingTherapyType
		itemDuration := existing.BookingDuration
		if existing.BookingDuration == 30 {
			adjacent, err = s.booking.FindAdjacentUnpaid(ctx, existing.BookingGoldID, existing.BookingOutcode,
				existing.BookingCustName, ttype, existing.BookingID, existing.BookingStart)
			if err != nil {
				return b, nil, errors.Wrap(err, "[SERVICE][PayBooking][FindAdjacentUnpaid]")
			}
			if adjacent != nil {
				itemDuration = 60
				adjStart := adjacent.BookingStart.In(bookingEntity.Jakarta)
				if adjStart.Before(saleStart) {
					saleStart = adjStart
				}
			}
		}
		item, err = s.booking.GetTherapyItemByName(ctx, existing.BookingOutcode, bookingEntity.TherapyItemName(ttype, itemDuration))
		if err != nil {
			return b, nil, errors.Wrap(err, "[SERVICE][PayBooking][GetTherapyItemByName]")
		}
		if item == nil {
			return b, nil, errors.New("item terapi " + bookingEntity.TherapyItemName(ttype, itemDuration) + " belum terdaftar di outlet ini")
		}
	} else {
		// booking lama tanpa tipe terapi: pakai item eksplisit
		if itemID <= 0 && existing.BookingItemID != nil {
			itemID = *existing.BookingItemID
		}
		if itemID <= 0 {
			return b, nil, errors.New("item terapi wajib dipilih")
		}
		item, err = s.booking.GetTherapyItem(ctx, itemID)
		if err != nil {
			return b, nil, errors.Wrap(err, "[SERVICE][PayBooking][GetTherapyItem]")
		}
		if item == nil {
			return b, nil, errors.New("item terapi tidak ditemukan")
		}
	}

	finalPrice := item.ItemPrice
	if customPrice > 0 && canOverridePrice(role) {
		finalPrice = customPrice
	} else if existing.BookingPrice != nil && existing.BookingPrice.IsPositive() {
		// harga custom yang dikunci penjual/admin saat booking dibuat
		finalPrice = int(existing.BookingPrice.IntPart())
	}

	custID := 0
	if existing.BookingCustID != nil {
		custID = *existing.BookingCustID
	}
	label := item.ItemName
	if adjacent != nil {
		label += " (GABUNGAN 2x30 MENIT)"
	}
	payload := s.buildSalePayload(existing.BookingGoldID, existing.BookingOutcode, existing.BookingCustName, creator,
		label, item.ItemPack, s.resolveStockID(ctx, existing.BookingOutcode, item.ItemID), finalPrice, custID, saleStart)
	saleID := payload.InsertData.Header.SaleID

	if err := s.booking.UpdateBookingPaid(ctx, bookingID, saleID, item.ItemID, item.ItemName, strconv.Itoa(finalPrice)); err != nil {
		return b, nil, errors.Wrap(err, "[SERVICE][PayBooking][UpdateBookingPaid]")
	}
	if adjacent != nil {
		if err := s.booking.UpdateBookingPaid(ctx, adjacent.BookingID, saleID, item.ItemID, item.ItemName+" (GABUNGAN)", "0"); err != nil {
			return b, nil, errors.Wrap(err, "[SERVICE][PayBooking][UpdateAdjacentPaid]")
		}
	}

	updated, err := s.booking.GetBookingByID(ctx, bookingID)
	if err != nil || updated == nil {
		return b, &payload, nil
	}
	return *updated, &payload, nil
}

// RemoveBooking menghapus booking yang BELUM dibayar (oleh penjual terapi atau admin)
// dan mencatat jejaknya ke tabel booking_remove_log.
func (s Service) RemoveBooking(ctx context.Context, role string, remover string, bookingID string) error {
	if role != "SELLER" && role != "ADMIN" {
		return errors.New("hanya penjual atau admin yang bisa menghapus booking")
	}

	existing, err := s.booking.GetBookingByID(ctx, bookingID)
	if err != nil {
		return errors.Wrap(err, "[SERVICE][RemoveBooking][GetBookingByID]")
	}
	if existing == nil {
		return errors.New("booking tidak ditemukan")
	}
	if existing.BookingStatus != bookingEntity.StatusUnpaid {
		return errors.New("hanya booking yang belum dibayar yang bisa dihapus")
	}

	if err := s.booking.RemoveBookingWithLog(ctx, *existing, remover, role); err != nil {
		return errors.Wrap(err, "[SERVICE][RemoveBooking][RemoveBookingWithLog]")
	}
	return nil
}

// SearchBuyers meneruskan pencarian pembeli terdaftar (untuk autocomplete penjual).
func (s Service) SearchBuyers(ctx context.Context, name string) ([]bookingEntity.BuyerInfo, error) {
	buyers, err := s.booking.SearchBuyers(ctx, name, 20)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][SearchBuyers]")
	}
	return buyers, nil
}

// resolveStockID mencari stock_id milik item terapi (relasi td_sale → stock → items).
// Item brand THERAPY otomatis punya baris stock; jika belum ada, fallback ke item_id.
func (s Service) resolveStockID(ctx context.Context, outcode string, itemID int) string {
	stockID, err := s.booking.GetStockIDByItem(ctx, outcode, itemID)
	if err != nil || stockID == "" {
		return strconv.Itoa(itemID)
	}
	return stockID
}

// buildSalePayload menyusun payload insert_sales (dipublish ke Kafka oleh handler)
// untuk booking berbayar. sale_id dibuat di sini supaya FE bisa langsung ambil nota PDF.
// custID > 0 = pembeli terdaftar (nota bisa menampilkan nama toko pembeli).
func (s Service) buildSalePayload(goldid int, outcode, custName, creator, itemLabel, itemPack, stockID string, price int, custID int, start time.Time) saleEntity.InsertSaleData {
	priceStr := strconv.Itoa(price)
	paymentYN := "Y"
	qty := 1
	stockName := itemLabel + " (Booking " + start.Format("02-01 15:04") + ")"
	pack := itemPack
	if pack == "" {
		pack = "SESI"
	}
	// kolom sale_salesperson maksimal 15 karakter
	salesperson := creator
	if len(salesperson) > 15 {
		salesperson = salesperson[:15]
	}

	return saleEntity.InsertSaleData{
		InsertData: saleEntity.InsertSales{
			Header: saleEntity.THSaleDetail{
				SaleID:            uuid.New().String(),
				SaleGoldID:        goldid,
				SaleOutcode:       outcode,
				SaleCustID:        &custID,
				SaleTranstotal:    priceStr,
				SaleTranspayment:  priceStr,
				SaleTranschange:   "0",
				SaleSalesperson:   &salesperson,
				SaleSalescustomer: &custName,
				SalePaymentyn:     &paymentYN,
			},
			Detail: []saleEntity.TDSaleDetail{
				{
					SaleStockid:         &stockID,
					SaleStockname:       &stockName,
					SaleQty:             &qty,
					SaleSalesprice:      priceStr,
					SaleTotalsalesprice: priceStr,
					SalePack:            &pack,
				},
			},
		},
	}
}
