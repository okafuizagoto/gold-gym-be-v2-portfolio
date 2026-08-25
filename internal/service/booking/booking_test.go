package booking

import (
	"context"
	"strings"
	"testing"
	"time"

	bookingEntity "gold-gym-be/internal/entity/booking"

	"github.com/shopspring/decimal"
)

type fakeBookingData struct {
	countInWindow int
	inserted      *bookingEntity.Booking
	byID          *bookingEntity.Booking
	adjacent      *bookingEntity.Booking
	itemPrices    map[string]int
	paidCalls     []string
}

func (f *fakeBookingData) ExpireOverdueBookings(ctx context.Context, now time.Time) error {
	return nil
}

func (f *fakeBookingData) GetBookingsByDate(ctx context.Context, goldid int, outcode string, date time.Time) ([]bookingEntity.Booking, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, bookingEntity.Jakarta)
	return []bookingEntity.Booking{
		{
			BookingID:           "b1",
			BookingStart:        start,
			BookingDuration:     60,
			BookingCustName:     "Budi",
			BookingRegisteredYN: "Y",
			BookingStatus:       bookingEntity.StatusUnpaid,
		},
	}, nil
}

func (f *fakeBookingData) CountActiveInWindow(ctx context.Context, goldid int, outcode string, winStart time.Time) (int, error) {
	return f.countInWindow, nil
}

func (f *fakeBookingData) InsertBooking(ctx context.Context, b bookingEntity.Booking) error {
	f.inserted = &b
	return nil
}

func (f *fakeBookingData) GetBookingByID(ctx context.Context, bookingID string) (*bookingEntity.Booking, error) {
	return f.byID, nil
}

func (f *fakeBookingData) UpdateBookingPaid(ctx context.Context, bookingID string, saleID string, itemID int, itemName string, price string) error {
	f.paidCalls = append(f.paidCalls, bookingID+"|"+price)
	return nil
}

func (f *fakeBookingData) UpdateBookingStatus(ctx context.Context, bookingID string, status string) error {
	return nil
}

func (f *fakeBookingData) SearchBuyers(ctx context.Context, name string, limit int) ([]bookingEntity.BuyerInfo, error) {
	return nil, nil
}

func (f *fakeBookingData) GetBuyerByID(ctx context.Context, goldid int) (*bookingEntity.BuyerInfo, error) {
	return &bookingEntity.BuyerInfo{GoldId: goldid, GoldNama: "Budi"}, nil
}

func (f *fakeBookingData) GetTherapyItem(ctx context.Context, itemID int) (*bookingEntity.TherapyItem, error) {
	return &bookingEntity.TherapyItem{ItemID: itemID, ItemName: "Terapi 1 Jam", ItemPrice: 100000}, nil
}

func (f *fakeBookingData) GetTherapyItemByName(ctx context.Context, outcode string, name string) (*bookingEntity.TherapyItem, error) {
	price := 100000
	if p, ok := f.itemPrices[name]; ok {
		price = p
	}
	return &bookingEntity.TherapyItem{ItemID: 7, ItemName: name, ItemPrice: price}, nil
}

func (f *fakeBookingData) FindAdjacentUnpaid(ctx context.Context, goldid int, outcode string, custName string, therapyType string, excludeID string, start time.Time) (*bookingEntity.Booking, error) {
	return f.adjacent, nil
}

func (f *fakeBookingData) RemoveBookingWithLog(ctx context.Context, b bookingEntity.Booking, removedBy string, removedRole string) error {
	return nil
}

func (f *fakeBookingData) GetOutletByCode(ctx context.Context, outcode string) (*bookingEntity.OutletInfo, error) {
	return &bookingEntity.OutletInfo{OutletGoldID: 4, OutletCode: outcode, OutletName: "Testwow", OutletType: "THERAPY"}, nil
}

func (f *fakeBookingData) GetStockIDByItem(ctx context.Context, outcode string, itemID int) (string, error) {
	return "STK000009", nil
}

func newTestService(f *fakeBookingData) Service {
	return Service{booking: f}
}

func futureReq(paid bool) bookingEntity.InsertBookingRequest {
	tomorrow := time.Now().Add(24 * time.Hour)
	return bookingEntity.InsertBookingRequest{
		Outcode:  "TES002",
		Date:     tomorrow.Format("2006-01-02"),
		Start:    "10:00",
		Duration: 60,
		CustName: "Siti",
		Paid:     paid,
		ItemID:   7,

		TherapyType: bookingEntity.TherapyTypeSofa,
	}
}

func TestInsertBookingFullSlotRejected(t *testing.T) {
	f := &fakeBookingData{countInWindow: bookingEntity.CapacityPerHalfHour}
	svc := newTestService(f)

	_, _, err := svc.InsertBooking(context.Background(), 10, "SELLER", "kasir", futureReq(false))
	if err == nil || !strings.Contains(err.Error(), "penuh") {
		t.Fatalf("expected slot penuh error, got %v", err)
	}
	if f.inserted != nil {
		t.Fatal("booking should not be inserted when slot full")
	}
}

// Regresi: kolom booking_price NOT NULL — booking belum bayar harus tersimpan
// dengan harga 0, bukan nil (database menolak NULL eksplisit).
func TestInsertBookingUnpaidPriceZeroNotNull(t *testing.T) {
	f := &fakeBookingData{countInWindow: 0}
	svc := newTestService(f)

	b, sale, err := svc.InsertBooking(context.Background(), 10, "SELLER", "kasir", futureReq(false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sale != nil {
		t.Fatal("unpaid booking must not build sale payload")
	}
	if b.BookingPrice == nil {
		t.Fatal("booking_price must not be nil (kolom NOT NULL)")
	}
	if !b.BookingPrice.IsZero() {
		t.Fatalf("unpaid booking price must be 0, got %s", b.BookingPrice)
	}
	if f.inserted == nil || f.inserted.BookingPrice == nil {
		t.Fatal("inserted row must carry non-nil booking_price")
	}
}

func TestInsertBookingPaidBuildsSalePayload(t *testing.T) {
	f := &fakeBookingData{countInWindow: 0}
	svc := newTestService(f)

	b, sale, err := svc.InsertBooking(context.Background(), 10, "SELLER", "kasir", futureReq(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.BookingStatus != bookingEntity.StatusPaid {
		t.Fatalf("expected PAID, got %s", b.BookingStatus)
	}
	if sale == nil {
		t.Fatal("expected sale payload for paid booking")
	}
	if sale.InsertData.Header.SaleID == "" || b.BookingSaleID == nil || *b.BookingSaleID != sale.InsertData.Header.SaleID {
		t.Fatal("booking_sale_id must match generated sale_id")
	}
	if sale.InsertData.Header.SaleGoldID != 4 {
		t.Fatalf("sale_gold_id must be outlet tenant id, got %d", sale.InsertData.Header.SaleGoldID)
	}
	if sale.InsertData.Header.SaleTranstotal != "100000" {
		t.Fatalf("unexpected total %s", sale.InsertData.Header.SaleTranstotal)
	}
}

func TestInsertBookingBuyerBooksForSelf(t *testing.T) {
	f := &fakeBookingData{countInWindow: 0}
	svc := newTestService(f)

	req := futureReq(false)
	req.CustID = 0
	req.CustName = ""

	b, _, err := svc.InsertBooking(context.Background(), 42, "BUYER", "buyer@mail.com", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.BookingCustID == nil || *b.BookingCustID != 42 {
		t.Fatal("buyer booking must reference own gold_id")
	}
	if b.BookingRegisteredYN != "Y" {
		t.Fatal("buyer booking must be marked registered")
	}
}

func TestPayBookingMergesAdjacentHalfHours(t *testing.T) {
	ttype := bookingEntity.TherapyTypeSofa
	start := time.Date(2026, 7, 20, 6, 0, 0, 0, bookingEntity.Jakarta)
	f := &fakeBookingData{
		byID: &bookingEntity.Booking{
			BookingID:          "b-main",
			BookingGoldID:      4,
			BookingOutcode:     "TES002",
			BookingStart:       start,
			BookingDuration:    30,
			BookingTherapyType: &ttype,
			BookingCustName:    "Siti",
			BookingStatus:      bookingEntity.StatusUnpaid,
		},
		adjacent: &bookingEntity.Booking{
			BookingID:       "b-adj",
			BookingStart:    start.Add(30 * time.Minute),
			BookingDuration: 30,
			BookingStatus:   bookingEntity.StatusUnpaid,
		},
		itemPrices: map[string]int{"SOFA 1 JAM": 25000, "SOFA 30 MENIT": 15000},
	}
	svc := newTestService(f)

	_, sale, err := svc.PayBooking(context.Background(), "kasir", "SELLER", "b-main", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sale == nil {
		t.Fatal("expected sale payload")
	}
	// 2x30 menit bersebelahan → satu nota 1 jam dengan harga 1 jam (sofa 25.000)
	if sale.InsertData.Header.SaleTranstotal != "25000" {
		t.Fatalf("expected merged 1-hour price 25000, got %s", sale.InsertData.Header.SaleTranstotal)
	}
	if len(sale.InsertData.Detail) != 1 || sale.InsertData.Detail[0].SaleStockname == nil ||
		!strings.Contains(*sale.InsertData.Detail[0].SaleStockname, "GABUNGAN") {
		t.Fatal("merged sale line must be labeled GABUNGAN")
	}
	// kedua booking harus ditandai lunas dengan nota yang sama
	if len(f.paidCalls) != 2 {
		t.Fatalf("expected both bookings marked paid, got %v", f.paidCalls)
	}
}

func TestPayBookingCustomPriceSellerOnly(t *testing.T) {
	ttype := bookingEntity.TherapyTypeDragon
	start := time.Date(2026, 7, 20, 7, 0, 0, 0, bookingEntity.Jakarta)
	makeFake := func() *fakeBookingData {
		return &fakeBookingData{
			byID: &bookingEntity.Booking{
				BookingID:          "b1",
				BookingGoldID:      4,
				BookingOutcode:     "TES002",
				BookingStart:       start,
				BookingDuration:    60,
				BookingTherapyType: &ttype,
				BookingCustName:    "Budi",
				BookingStatus:      bookingEntity.StatusUnpaid,
			},
			itemPrices: map[string]int{"KURSI DRAGON 1 JAM": 35000},
		}
	}

	// penjual boleh pakai harga custom
	_, sale, err := newTestService(makeFake()).PayBooking(context.Background(), "kasir", "SELLER", "b1", 0, 50000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sale.InsertData.Header.SaleTranstotal != "50000" {
		t.Fatalf("seller custom price ignored, got %s", sale.InsertData.Header.SaleTranstotal)
	}

	// pembeli tidak boleh — tetap harga default
	_, sale, err = newTestService(makeFake()).PayBooking(context.Background(), "buyer", "BUYER", "b1", 0, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sale.InsertData.Header.SaleTranstotal != "35000" {
		t.Fatalf("buyer must not override price, got %s", sale.InsertData.Header.SaleTranstotal)
	}
}

// Harga custom yang dikunci penjual saat membuat booking UNPAID harus dipakai
// saat pembayaran walau parameter customprice tidak dikirim lagi.
func TestPayBookingUsesStoredCustomPrice(t *testing.T) {
	ttype := bookingEntity.TherapyTypeSofa
	stored := decimal.NewFromInt(42000)
	start := time.Date(2026, 7, 20, 7, 0, 0, 0, bookingEntity.Jakarta)
	f := &fakeBookingData{
		byID: &bookingEntity.Booking{
			BookingID:          "b1",
			BookingGoldID:      4,
			BookingOutcode:     "TES002",
			BookingStart:       start,
			BookingDuration:    60,
			BookingTherapyType: &ttype,
			BookingCustName:    "Budi",
			BookingStatus:      bookingEntity.StatusUnpaid,
			BookingPrice:       &stored,
		},
		itemPrices: map[string]int{"SOFA 1 JAM": 25000},
	}

	_, sale, err := newTestService(f).PayBooking(context.Background(), "kasir", "SELLER", "b1", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sale.InsertData.Header.SaleTranstotal != "42000" {
		t.Fatalf("stored custom price ignored, got %s", sale.InsertData.Header.SaleTranstotal)
	}
}

func TestGetSlotsMarksUnpaidWindows(t *testing.T) {
	f := &fakeBookingData{}
	svc := newTestService(f)

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	slots, err := svc.GetSlots(context.Background(), "TES002", tomorrow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(slots) != (bookingEntity.CloseHour-bookingEntity.OpenHour)*2 {
		t.Fatalf("expected %d slots, got %d", (bookingEntity.CloseHour-bookingEntity.OpenHour)*2, len(slots))
	}
	// booking 09:00 durasi 60 menit menempati jendela 09:00 dan 09:30
	for _, s := range slots {
		switch s.Start {
		case "09:00", "09:30":
			if s.Used != 1 || !s.HasUnpaid {
				t.Fatalf("slot %s: expected used=1 hasUnpaid=true, got used=%d hasUnpaid=%v", s.Start, s.Used, s.HasUnpaid)
			}
		default:
			if s.Used != 0 {
				t.Fatalf("slot %s should be empty", s.Start)
			}
		}
	}
}
