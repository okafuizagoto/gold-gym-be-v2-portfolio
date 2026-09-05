package goldgym

import (
	"bytes"
	"context"
	"testing"
	"time"

	goldQuotaEntity "gold-gym-be/internal/entity/quota"
	goldSaleEntity "gold-gym-be/internal/entity/sales"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type fakeSalesData struct{}

func (fakeSalesData) WithTransactionThSale(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return nil
}
func (fakeSalesData) GetLastThSaleCode(ctx context.Context, goldid int, code string) (*string, *string, error) {
	return nil, nil, nil
}
func (fakeSalesData) InsertThSale(ctx context.Context, tx *gorm.DB, items []goldSaleEntity.ThSale) error {
	return nil
}
func (fakeSalesData) InsertTdSale(ctx context.Context, tx *gorm.DB, items []goldSaleEntity.TdSale) error {
	return nil
}
func (fakeSalesData) GetLastTdSaleCode(ctx context.Context, saleid string) (*string, error) {
	return nil, nil
}
func (fakeSalesData) GetThSale(ctx context.Context, goldid int, custid int, name string, outcode string, page, length int) ([]goldSaleEntity.ThSale, error) {
	return nil, nil
}
func (fakeSalesData) GetTotalThSale(ctx context.Context, goldid int, custid int, name string, outcode string) (int64, error) {
	return 0, nil
}
func (fakeSalesData) GetThSaleByID(ctx context.Context, goldid int, custid int, saleid string) (*goldSaleEntity.ThSale, error) {
	now := time.Now()
	transtime := "141530"
	person := "KASIR01"
	customer := "Budi"
	total := decimal.NewFromInt(45000)
	payment := decimal.NewFromInt(50000)
	change := decimal.NewFromInt(5000)
	return &goldSaleEntity.ThSale{
		SaleID:            saleid,
		SaleGoldID:        goldid,
		SaleOutcode:       "OUT01",
		SaleTrancnum:      "OUT0125071514150001",
		SaleTranstime:     &transtime,
		SaleTransdate:     &now,
		SaleTranstotal:    &total,
		SaleTranspayment:  &payment,
		SaleTranschange:   &change,
		SaleSalesperson:   &person,
		SaleSalescustomer: &customer,
		SalePaymentyn:     "Y",
		SaleCreatedAt:     now,
	}, nil
}
func (fakeSalesData) GetTdSaleBySaleID(ctx context.Context, saleid string) ([]goldSaleEntity.TdSale, error) {
	name := "Protein Bar"
	qty := 3
	price := decimal.NewFromInt(15000)
	totalPrice := decimal.NewFromInt(45000)
	return []goldSaleEntity.TdSale{
		{
			SaleID:              saleid,
			SaleStockname:       &name,
			SaleQty:             &qty,
			SaleSalesprice:      &price,
			SaleTotalsalesprice: &totalPrice,
		},
	}, nil
}
func (fakeSalesData) GetSaleOutletInfo(ctx context.Context, goldid int, outcode string) (*goldSaleEntity.SaleOutletInfo, error) {
	return &goldSaleEntity.SaleOutletInfo{
		OutletCode:    outcode,
		OutletName:    "Gold Gym Outlet 1",
		OutletAddress: "Jl. Contoh No. 1, Jakarta",
	}, nil
}

func (fakeSalesData) GetSaleCustomerToko(ctx context.Context, custid int) (string, error) {
	return "", nil
}

func (fakeSalesData) ConfirmSaleMeja(ctx context.Context, mejaIDs []int, outcode, saleID string) (int64, error) {
	return int64(len(mejaIDs)), nil
}

func (fakeSalesData) GetMejaNamesByIDs(ctx context.Context, outcode string, mejaIDs []int) ([]string, error) {
	return nil, nil
}

func (fakeSalesData) IsPosCustomerOptional(ctx context.Context, goldid int, outcode string) (bool, error) {
	return false, nil
}

func (fakeSalesData) GetOutletTypeByCode(ctx context.Context, goldid int, outcode string) (string, error) {
	return "RETAIL", nil
}

func (fakeSalesData) GetPosOutletsForAdmin(ctx context.Context, search string) ([]goldSaleEntity.PosOutlet, error) {
	return nil, nil
}

func (fakeSalesData) SetPosCustomerOptional(ctx context.Context, goldid int, outcode string, optional bool, addedBy string) error {
	return nil
}

func (fakeSalesData) MarkBookingsPaid(ctx context.Context, bookingIDs []string, saleID string) (int64, error) {
	return int64(len(bookingIDs)), nil
}

func (fakeSalesData) GetTotalProofBytes(ctx context.Context) (int64, error) {
	return 0, nil
}

func (fakeSalesData) InsertPaymentProof(ctx context.Context, proof goldSaleEntity.PaymentProof) (goldSaleEntity.PaymentProof, error) {
	return proof, nil
}

func (fakeSalesData) GetPaymentProofsBySale(ctx context.Context, saleID string) ([]goldSaleEntity.PaymentProof, error) {
	return nil, nil
}

func (fakeSalesData) GetPaymentProofByID(ctx context.Context, proofID int) (*goldSaleEntity.PaymentProof, error) {
	return nil, nil
}
func (fakeSalesData) DeletePaymentProof(ctx context.Context, proofID int) error {
	return nil
}
func (fakeSalesData) GetUserStorageUsedKB(ctx context.Context, goldID int) (int, error) {
	return 0, nil
}
func (fakeSalesData) AddUserStorageUsedKB(ctx context.Context, goldID int, deltaKB int) error {
	return nil
}
func (fakeSalesData) InsertStorageDeleteHistory(ctx context.Context, h goldQuotaEntity.StorageDeleteHistory) error {
	return nil
}

func (fakeSalesData) GetSaleReportItems(ctx context.Context, goldid int, outcode, date string) ([]goldSaleEntity.SaleReportItem, error) {
	return nil, nil
}

func (fakeSalesData) GetSaleDailyTotals(ctx context.Context, goldid int, outcode, start, end string) ([]goldSaleEntity.SaleDailyTotal, error) {
	return nil, nil
}

func (fakeSalesData) IsPaymentProofGloballyEnabled(ctx context.Context) (bool, error) {
	return true, nil
}

func (fakeSalesData) SetPaymentProofGlobal(ctx context.Context, enabled bool, updatedBy string) error {
	return nil
}

func (fakeSalesData) IsOutletProofEnabled(ctx context.Context, goldid int, outcode string) (bool, error) {
	return true, nil
}

func (fakeSalesData) IsUserProofEnabled(ctx context.Context, goldid int) (bool, error) {
	return true, nil
}

func (fakeSalesData) GetProofAccessOutlets(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessOutlet, error) {
	return nil, nil
}

func (fakeSalesData) SetProofOutletEnabled(ctx context.Context, goldid int, outcode string, enabled bool) error {
	return nil
}

func (fakeSalesData) GetProofAccessUsers(ctx context.Context, search string) ([]goldSaleEntity.ProofAccessUser, error) {
	return nil, nil
}

func (fakeSalesData) SetProofUserEnabled(ctx context.Context, goldid int, enabled bool) error {
	return nil
}

func TestGetSaleReceiptPDF(t *testing.T) {
	svc := New(fakeSalesData{}, opentracing.NoopTracer{}, jaegerLog.Factory{})

	pdfBytes, trancnum, err := svc.GetSaleReceiptPDF(context.Background(), 1, 0, "abc-123")
	if err != nil {
		t.Fatalf("GetSaleReceiptPDF error: %v", err)
	}
	if trancnum != "OUT0125071514150001" {
		t.Fatalf("unexpected trancnum: %s", trancnum)
	}
	if !bytes.HasPrefix(pdfBytes, []byte("%PDF")) {
		t.Fatalf("output is not a valid PDF, first bytes: %q", pdfBytes[:10])
	}
	if len(pdfBytes) < 500 {
		t.Fatalf("PDF too small: %d bytes", len(pdfBytes))
	}
}

func TestFormatRupiah(t *testing.T) {
	cases := map[int64]string{
		0:        "0",
		999:      "999",
		1000:     "1.000",
		45000:    "45.000",
		1234567:  "1.234.567",
		-1500000: "-1.500.000",
	}
	for input, want := range cases {
		d := decimal.NewFromInt(input)
		if got := formatRupiah(&d); got != want {
			t.Errorf("formatRupiah(%d) = %s, want %s", input, got, want)
		}
	}
	if got := formatRupiah(nil); got != "0" {
		t.Errorf("formatRupiah(nil) = %s, want 0", got)
	}
}

func (fakeSalesData) UpdateSalePaymentYN(ctx context.Context, goldid int, saleid string) (int64, error) {
	return 1, nil
}

func (fakeSalesData) GetOutletGoldIDByCode(ctx context.Context, outcode string) (int, error) {
	return 1, nil
}
