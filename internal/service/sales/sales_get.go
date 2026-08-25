package goldgym

import (
	"context"
	"math"

	"gold-gym-be/internal/entity"
	goldSaleEntity "gold-gym-be/internal/entity/sales"

	"github.com/pkg/errors"
)

func (s Service) GetSales(ctx context.Context, goldid int, custid int, name string, outcode string, page, length int) ([]goldSaleEntity.ThSale, goldSaleEntity.MetadataPaginationDetail, error) {
	var (
		sales          []goldSaleEntity.ThSale
		metadataDetail goldSaleEntity.MetadataPaginationDetail
		err            error
	)

	sales, err = s.goldgymsale.GetThSale(ctx, goldid, custid, name, outcode, page, length)
	if err != nil {
		return sales, metadataDetail, errors.Wrap(err, "[Service][GetSales]")
	}

	total, err := s.goldgymsale.GetTotalThSale(ctx, goldid, custid, name, outcode)
	if err != nil {
		return sales, metadataDetail, errors.Wrap(err, "[Service][GetSales][GetTotalThSale]")
	}

	totalPage := 1
	if length > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(length)))
	}
	metadataDetail = goldSaleEntity.MetadataPaginationDetail{
		Page:      page,
		Limit:     length,
		TotalData: int(total),
		TotalPage: totalPage,
	}

	return sales, metadataDetail, nil
}

func (s Service) GetSaleDetail(ctx context.Context, goldid int, custid int, saleid string) (goldSaleEntity.SaleWithDetail, error) {
	var result goldSaleEntity.SaleWithDetail

	header, err := s.goldgymsale.GetThSaleByID(ctx, goldid, custid, saleid)
	if err != nil {
		return result, errors.Wrap(err, "[Service][GetSaleDetail][GetThSaleByID]")
	}
	if header == nil {
		return result, entity.ErrNotFound
	}

	details, err := s.goldgymsale.GetTdSaleBySaleID(ctx, saleid)
	if err != nil {
		return result, errors.Wrap(err, "[Service][GetSaleDetail][GetTdSaleBySaleID]")
	}

	result.Header = *header
	result.Detail = details
	return result, nil
}

// MarkSalePaid mengubah transaksi BELUM LUNAS menjadi LUNAS (hanya penjual/admin).
func (s Service) MarkSalePaid(ctx context.Context, goldid int, saleid string) (goldSaleEntity.ThSale, error) {
	var result goldSaleEntity.ThSale

	affected, err := s.goldgymsale.UpdateSalePaymentYN(ctx, goldid, saleid)
	if err != nil {
		return result, errors.Wrap(err, "[Service][MarkSalePaid]")
	}
	if affected == 0 {
		// tidak ada baris berubah: transaksi tidak ada atau sudah lunas
		header, err := s.goldgymsale.GetThSaleByID(ctx, goldid, 0, saleid)
		if err != nil {
			return result, errors.Wrap(err, "[Service][MarkSalePaid][GetThSaleByID]")
		}
		if header == nil {
			return result, entity.ErrNotFound
		}
		return *header, errors.New("transaksi sudah berstatus lunas")
	}

	header, err := s.goldgymsale.GetThSaleByID(ctx, goldid, 0, saleid)
	if err != nil {
		return result, errors.Wrap(err, "[Service][MarkSalePaid][GetThSaleByID]")
	}
	if header != nil {
		result = *header
	}
	return result, nil
}

// ResolveOutletGoldID dipakai handler saat pembeli (BUYER) membuat transaksi.
func (s Service) ResolveOutletGoldID(ctx context.Context, outcode string) (int, error) {
	return s.goldgymsale.GetOutletGoldIDByCode(ctx, outcode)
}
