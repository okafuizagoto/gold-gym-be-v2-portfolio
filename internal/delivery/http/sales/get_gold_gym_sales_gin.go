package goldgym

import (
	"errors"
	"gold-gym-be/internal/entity"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) GetGoldGymSaleGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)
	span := h.tracer.StartSpan("GetGoldGym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	// id, _ := strconv.Atoi(c.Query("id"))

	types := c.Query("type")
	// pembeli (BUYER) hanya melihat transaksi miliknya sendiri (sale_cust_id)
	custID := 0
	if c.GetString("role") == "BUYER" {
		custID = c.GetInt("user")
	}
	switch types {
	case "getallsales":
		userID := c.GetInt("user")
		page, _ := strconv.Atoi(c.Query("page"))
		length, _ := strconv.Atoi(c.Query("length"))
		result, metadata, err = h.goldgymSvcSale.GetSales(ctx, userID, custID, c.Query("name"), c.Query("code"), page, length)
	case "getsaledetail":
		userID := c.GetInt("user")
		result, err = h.goldgymSvcSale.GetSaleDetail(ctx, userID, custID, c.Query("saleid"))
	case "salesreport":
		// laporan penjualan penjual (per hari/minggu/bulan), scope outlet sendiri
		if c.GetString("role") == "BUYER" {
			c.JSON(http.StatusForbidden, gin.H{"error": "tidak diizinkan"})
			return
		}
		userID := c.GetInt("user")
		result, err = h.goldgymSvcSale.GetSalesReport(ctx, userID, c.Query("code"), c.Query("mode"), c.Query("date"))
	case "receipt":
		userID := c.GetInt("user")
		pdfBytes, trancnum, errPdf := h.goldgymSvcSale.GetSaleReceiptPDF(ctx, userID, custID, c.Query("saleid"))
		if errPdf != nil {
			if errors.Is(errPdf, entity.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "sale not found"})
				return
			}
			h.logger.For(ctx).Error("failed to generate receipt pdf", zap.Error(errPdf))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate receipt"})
			return
		}
		c.Header("Content-Disposition", `inline; filename="nota-`+trancnum+`.pdf"`)
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
		return
	case "proofs":
		// daftar bukti pembayaran transfer milik satu nota (Sales History)
		saleID := c.Query("saleid")
		if saleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter saleid wajib diisi"})
			return
		}
		proofs, errProofs := h.goldgymSvcSale.GetPaymentProofs(ctx, saleID)
		if errProofs != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errProofs.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": proofs})
		return
	case "proofphoto":
		// unduh/tampilkan file foto bukti pembayaran
		proofID, _ := strconv.Atoi(c.Query("proofid"))
		proof, content, errProof := h.goldgymSvcSale.GetPaymentProofFile(ctx, proofID)
		if errProof != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": errProof.Error()})
			return
		}
		mime := proof.ProofMime
		if mime == "" {
			mime = "image/jpeg"
		}
		c.Header("Content-Disposition", `inline; filename="`+proof.ProofFilename+`"`)
		c.Data(http.StatusOK, mime, content)
		return
	case "customerrequired":
		// dipakai POS FE: apakah outlet ini wajib mengisi nama customer
		outcode := c.Query("code")
		gid, _ := h.goldgymSvcSale.ResolveOutletGoldID(ctx, outcode)
		required := false
		if gid > 0 {
			required, _ = h.goldgymSvcSale.IsPosCustomerRequired(ctx, gid, outcode)
		}
		c.JSON(http.StatusOK, gin.H{"required": required})
		return
	case "posoutlets":
		// ADMIN: daftar outlet RETAIL + alamat + penanda akses "POS tanpa customer".
		// name mencocokkan nama ATAU alamat outlet (untuk pencarian alamat).
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengakses"})
			return
		}
		outlets, errPos := h.goldgymSvcSale.GetPosOutletsForAdmin(ctx, c.Query("name"))
		if errPos != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errPos.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": outlets})
		return
	case "proofvisibility":
		// dipakai FE (POS/belanja pembeli/Sales History): apakah fitur bukti
		// pembayaran tampil untuk user ini di outlet ini (code boleh kosong).
		userID := c.GetInt("user")
		enabled, errProof := h.goldgymSvcSale.IsPaymentProofFeatureEnabled(ctx, userID, c.Query("code"))
		if errProof != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errProof.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"enabled": enabled})
		return
	case "proofglobal":
		// ADMIN: status gerbang global fitur bukti pembayaran.
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengakses"})
			return
		}
		enabled, errGlobal := h.goldgymSvcSale.GetPaymentProofGlobal(ctx)
		if errGlobal != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errGlobal.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"enabled": enabled})
		return
	case "proofoutlets":
		// ADMIN: daftar SEMUA outlet + status aktif/nonaktif fitur bukti pembayaran.
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengakses"})
			return
		}
		outlets, errProofOut := h.goldgymSvcSale.GetProofAccessOutlets(ctx, c.Query("name"))
		if errProofOut != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errProofOut.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": outlets})
		return
	case "proofusers":
		// ADMIN: daftar user (penjual/pembeli) + status aktif/nonaktif fitur bukti pembayaran.
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengakses"})
			return
		}
		users, errProofUsers := h.goldgymSvcSale.GetProofAccessUsers(ctx, c.Query("name"))
		if errProofUsers != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errProofUsers.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": users})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type: " + types})
		return
	}

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, entity.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, entity.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	// }

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", c.Request.Method, c.Request.URL)
	h.logger.For(ctx).Info(
		"HTTP request done",
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
	)
	c.JSON(200, resp)

	return
}
