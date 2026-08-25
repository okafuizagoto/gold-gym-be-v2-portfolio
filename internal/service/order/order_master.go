package order

import (
	"context"
	"strconv"
	"strings"

	orderEntity "gold-gym-be/internal/entity/order"
	saleEntity "gold-gym-be/internal/entity/sales"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// GetPublicOutlets daftar outlet (non-THERAPY) yang bisa dipilih pembeli.
func (s Service) GetPublicOutlets(ctx context.Context, name string) ([]orderEntity.PublicOutlet, error) {
	outlets, err := s.order.GetPublicOutlets(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetPublicOutlets]")
	}
	return outlets, nil
}

// GetAllOutletsForAdmin daftar semua outlet (+penanda visible) untuk admin.
func (s Service) GetAllOutletsForAdmin(ctx context.Context, name string) ([]orderEntity.PublicOutlet, error) {
	outlets, err := s.order.GetAllOutletsForAdmin(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetAllOutletsForAdmin]")
	}
	return outlets, nil
}

// AddVisibleOutlet admin menandai outlet boleh dilihat pembeli.
func (s Service) AddVisibleOutlet(ctx context.Context, goldid int, outcode, addedBy string) error {
	if goldid <= 0 || strings.TrimSpace(outcode) == "" {
		return errors.New("gold_id dan outcode wajib diisi")
	}
	outlet, err := s.order.GetOutletByCode(ctx, goldid, outcode)
	if err != nil || outlet == nil {
		return errors.New("outlet tidak ditemukan")
	}
	if outlet.OutletType == "THERAPY" {
		return errors.New("outlet terapi tidak bisa ditampilkan ke pembeli")
	}
	if err := s.order.AddVisibleOutlet(ctx, goldid, outcode, addedBy); err != nil {
		return errors.Wrap(err, "[SERVICE][AddVisibleOutlet]")
	}
	return nil
}

// RemoveVisibleOutlet admin mencabut outlet dari daftar yang dilihat pembeli.
func (s Service) RemoveVisibleOutlet(ctx context.Context, goldid int, outcode string) error {
	if goldid <= 0 || strings.TrimSpace(outcode) == "" {
		return errors.New("gold_id dan outcode wajib diisi")
	}
	if err := s.order.RemoveVisibleOutlet(ctx, goldid, outcode); err != nil {
		return errors.Wrap(err, "[SERVICE][RemoveVisibleOutlet]")
	}
	return nil
}

// GetOutletCatalog daftar barang satu outlet untuk dipesan pembeli.
func (s Service) GetOutletCatalog(ctx context.Context, goldid int, outcode, name string) ([]orderEntity.CatalogItem, error) {
	if goldid <= 0 || strings.TrimSpace(outcode) == "" {
		return nil, errors.New("gold_id dan outcode wajib diisi")
	}
	items, err := s.order.GetOutletCatalog(ctx, goldid, outcode, name)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetOutletCatalog]")
	}
	return items, nil
}

// InsertOrder membuat pesanan baru dari pembeli. Status awal WAITING.
// TRANSFER dianggap lunas (paid_yn=Y, pembeli upload bukti); TUNAI belum lunas.
func (s Service) InsertOrder(ctx context.Context, buyerID int, req orderEntity.InsertOrderRequest) (orderEntity.Order, error) {
	var empty orderEntity.Order

	if buyerID <= 0 {
		return empty, errors.New("identitas pembeli tidak valid")
	}
	payType := strings.ToUpper(strings.TrimSpace(req.PayType))
	if payType != orderEntity.PayTunai && payType != orderEntity.PayTransfer {
		return empty, errors.New("tipe pembayaran harus TUNAI atau TRANSFER")
	}
	if req.GoldID <= 0 || strings.TrimSpace(req.Outcode) == "" {
		return empty, errors.New("outlet tujuan wajib dipilih")
	}
	if len(req.Lines) == 0 {
		return empty, errors.New("keranjang belanja kosong")
	}

	// validasi outlet tujuan: harus ada dan bukan THERAPY
	outlet, err := s.order.GetOutletByCode(ctx, req.GoldID, req.Outcode)
	if err != nil || outlet == nil {
		return empty, errors.New("outlet tujuan tidak ditemukan")
	}
	if outlet.OutletType == "THERAPY" {
		return empty, errors.New("outlet terapi tidak menerima pesanan barang")
	}

	outletName := req.OutletName
	if strings.TrimSpace(outletName) == "" {
		outletName = outlet.OutletName
	}

	buyerName, err := s.order.GetBuyerName(ctx, buyerID)
	if err != nil {
		buyerName = ""
	}

	orderID := uuid.New().String()
	total := decimal.Zero
	details := make([]orderEntity.OrderDetail, 0, len(req.Lines))
	for i, l := range req.Lines {
		if strings.TrimSpace(l.StockID) == "" {
			return empty, errors.New("stock_id item ke-" + strconv.Itoa(i+1) + " kosong")
		}
		if l.Qty <= 0 {
			return empty, errors.New("qty item ke-" + strconv.Itoa(i+1) + " harus lebih dari 0")
		}
		price := decimal.NewFromInt(int64(l.Price))
		lineTotal := price.Mul(decimal.NewFromInt(int64(l.Qty)))
		total = total.Add(lineTotal)
		var pack *string
		if strings.TrimSpace(l.Pack) != "" {
			p := l.Pack
			pack = &p
		}
		details = append(details, orderEntity.OrderDetail{
			OdOrderID:   orderID,
			OdStockID:   l.StockID,
			OdStockName: l.StockName,
			OdQty:       l.Qty,
			OdPrice:     price,
			OdTotal:     lineTotal,
			OdPack:      pack,
		})
	}

	paidYN := "N"
	if payType == orderEntity.PayTransfer {
		paidYN = "Y"
	}

	header := orderEntity.Order{
		OrderID:         orderID,
		OrderBuyerID:    buyerID,
		OrderBuyerName:  buyerName,
		OrderGoldID:     req.GoldID,
		OrderOutcode:    req.Outcode,
		OrderOutletName: outletName,
		OrderTotal:      total,
		OrderPayType:    payType,
		OrderPaidYN:     paidYN,
		OrderStatus:     orderEntity.StatusWaiting,
	}

	if err := s.order.InsertOrder(ctx, header, details); err != nil {
		return empty, errors.Wrap(err, "[SERVICE][InsertOrder]")
	}
	return header, nil
}

// GetBuyerOrders daftar pesanan milik pembeli (dashboard).
func (s Service) GetBuyerOrders(ctx context.Context, buyerID int) ([]orderEntity.Order, error) {
	orders, err := s.order.GetOrdersByBuyer(ctx, buyerID)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetBuyerOrders]")
	}
	return orders, nil
}

// GetSellerOrders daftar pesanan masuk untuk penjual.
func (s Service) GetSellerOrders(ctx context.Context, goldid int, status string) ([]orderEntity.Order, error) {
	orders, err := s.order.GetOrdersBySeller(ctx, goldid, status)
	if err != nil {
		return nil, errors.Wrap(err, "[SERVICE][GetSellerOrders]")
	}
	return orders, nil
}

// GetOrderDetail header + item satu pesanan. Otorisasi: pemesan sendiri, atau
// penjual pemilik outlet pesanan tersebut.
func (s Service) GetOrderDetail(ctx context.Context, orderID string, requesterID int) (orderEntity.OrderWithDetail, error) {
	var out orderEntity.OrderWithDetail
	header, err := s.order.GetOrderByID(ctx, orderID)
	if err != nil || header == nil {
		return out, errors.New("pesanan tidak ditemukan")
	}
	if header.OrderBuyerID != requesterID && header.OrderGoldID != requesterID {
		return out, errors.New("tidak berhak mengakses pesanan ini")
	}
	details, err := s.order.GetOrderDetails(ctx, orderID)
	if err != nil {
		return out, errors.Wrap(err, "[SERVICE][GetOrderDetail]")
	}
	out.Header = *header
	out.Detail = details
	return out, nil
}

// ConfirmOrder penjual mengkonfirmasi pesanan WAITING -> PROCESS.
func (s Service) ConfirmOrder(ctx context.Context, orderID string, sellerGoldID int) (orderEntity.Order, error) {
	header, err := s.guardSellerAction(ctx, orderID, sellerGoldID, orderEntity.StatusWaiting)
	if err != nil {
		return orderEntity.Order{}, err
	}
	if err := s.order.UpdateOrderStatus(ctx, orderID, orderEntity.StatusProcess, nil); err != nil {
		return orderEntity.Order{}, errors.Wrap(err, "[SERVICE][ConfirmOrder]")
	}
	header.OrderStatus = orderEntity.StatusProcess
	return *header, nil
}

// RejectOrder penjual menolak pesanan WAITING -> REJECT (dengan alasan).
func (s Service) RejectOrder(ctx context.Context, orderID string, sellerGoldID int, reason string) (orderEntity.Order, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return orderEntity.Order{}, errors.New("alasan penolakan wajib diisi")
	}
	header, err := s.guardSellerAction(ctx, orderID, sellerGoldID, orderEntity.StatusWaiting)
	if err != nil {
		return orderEntity.Order{}, err
	}
	if err := s.order.UpdateOrderStatus(ctx, orderID, orderEntity.StatusReject, &reason); err != nil {
		return orderEntity.Order{}, errors.Wrap(err, "[SERVICE][RejectOrder]")
	}
	header.OrderStatus = orderEntity.StatusReject
	header.OrderRejectReason = &reason
	return *header, nil
}

// FinishOrder penjual menyelesaikan pesanan PROCESS -> FINISH. Nota dibuat:
// mengembalikan payload insert_sales untuk dipublish ke Kafka oleh handler.
func (s Service) FinishOrder(ctx context.Context, orderID string, sellerGoldID int, creator string) (orderEntity.Order, *saleEntity.InsertSaleData, error) {
	header, err := s.guardSellerAction(ctx, orderID, sellerGoldID, orderEntity.StatusProcess)
	if err != nil {
		return orderEntity.Order{}, nil, err
	}
	details, err := s.order.GetOrderDetails(ctx, orderID)
	if err != nil {
		return orderEntity.Order{}, nil, errors.Wrap(err, "[SERVICE][FinishOrder][details]")
	}

	saleID := uuid.New().String()
	payload := s.buildSalePayload(saleID, *header, details, creator)

	if err := s.order.FinishOrder(ctx, orderID, saleID); err != nil {
		return orderEntity.Order{}, nil, errors.Wrap(err, "[SERVICE][FinishOrder]")
	}
	header.OrderStatus = orderEntity.StatusFinish
	header.OrderSaleID = &saleID
	return *header, &payload, nil
}

// guardSellerAction memvalidasi bahwa pesanan ada, milik outlet penjual ini, dan
// berada pada status yang diharapkan sebelum aksi diproses.
func (s Service) guardSellerAction(ctx context.Context, orderID string, sellerGoldID int, expected string) (*orderEntity.Order, error) {
	header, err := s.order.GetOrderByID(ctx, orderID)
	if err != nil || header == nil {
		return nil, errors.New("pesanan tidak ditemukan")
	}
	if header.OrderGoldID != sellerGoldID {
		return nil, errors.New("pesanan ini bukan milik outlet Anda")
	}
	if header.OrderStatus != expected {
		return nil, errors.New("status pesanan sudah berubah, muat ulang daftar pesanan")
	}
	return header, nil
}

// buildSalePayload menyusun payload insert_sales dari pesanan yang selesai.
// custID = pembeli (nota membawa identitas pembeli). paymentyn: TRANSFER lunas (Y),
// TUNAI belum lunas (N) -> penjual tandai lunas di Sales History.
func (s Service) buildSalePayload(saleID string, o orderEntity.Order, details []orderEntity.OrderDetail, creator string) saleEntity.InsertSaleData {
	total := o.OrderTotal.String()
	paymentYN := "N"
	if o.OrderPayType == orderEntity.PayTransfer {
		paymentYN = "Y"
	}
	custID := o.OrderBuyerID
	custName := o.OrderBuyerName
	salesperson := creator
	if len(salesperson) > 15 {
		salesperson = salesperson[:15]
	}

	detail := make([]saleEntity.TDSaleDetail, 0, len(details))
	for _, d := range details {
		stockID := d.OdStockID
		stockName := d.OdStockName
		qty := d.OdQty
		pack := d.OdPack
		detail = append(detail, saleEntity.TDSaleDetail{
			SaleStockid:         &stockID,
			SaleStockname:       &stockName,
			SaleQty:             &qty,
			SaleSalesprice:      d.OdPrice.String(),
			SaleTotalsalesprice: d.OdTotal.String(),
			SalePack:            pack,
		})
	}

	return saleEntity.InsertSaleData{
		InsertData: saleEntity.InsertSales{
			Header: saleEntity.THSaleDetail{
				SaleID:            saleID,
				SaleGoldID:        o.OrderGoldID,
				SaleOutcode:       o.OrderOutcode,
				SaleCustID:        &custID,
				SaleTranstotal:    total,
				SaleTranspayment:  total,
				SaleTranschange:   "0",
				SaleSalesperson:   &salesperson,
				SaleSalescustomer: &custName,
				SalePaymentyn:     &paymentYN,
			},
			Detail: detail,
		},
	}
}
