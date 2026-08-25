package discount

import (
	"net/http"
	"strconv"

	discountEntity "gold-gym-be/internal/entity/discount"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func atoiDefault(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// GetGoldGymDiscountGin menangani GET /v2/discount
//
//	type=getitemsforoutlet : sumber item-picker saat tambah diskon (code, name)
//	type=getdiscounts       : daftar diskon outlet (code, name, page, length)
//	type=getactivebyoutlet  : semua diskon ACTIVE di outlet (code) -- dipakai POS
//	type=gethistory         : riwayat satu diskon (discountid, page, length)
func (h *Handler) GetGoldGymDiscountGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("GetDiscount", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	userID := c.GetInt("user")

	switch c.Query("type") {
	case "getitemsforoutlet":
		items, err := h.discountSvc.GetItemsForOutlet(ctx, userID, c.Query("code"), c.Query("name"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": items})
	case "getdiscounts":
		page := atoiDefault(c.Query("page"), 1)
		length := atoiDefault(c.Query("length"), 10)
		rows, metadata, err := h.discountSvc.GetDiscounts(ctx, userID, c.Query("code"), c.Query("name"), page, length)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": rows, "metadata": metadata})
	case "getactivebyoutlet":
		rows, err := h.discountSvc.GetActiveDiscountsByOutlet(ctx, userID, c.Query("code"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": rows})
	case "gethistory":
		discountID, _ := strconv.Atoi(c.Query("discountid"))
		page := atoiDefault(c.Query("page"), 1)
		length := atoiDefault(c.Query("length"), 10)
		rows, metadata, err := h.discountSvc.GetDiscountHistory(ctx, discountID, page, length)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": rows, "metadata": metadata})
	case "generatevouchercode":
		code, err := h.discountSvc.GenerateVoucherCode(ctx, c.Query("code"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"voucher_code": code}})
	case "getvouchers":
		page := atoiDefault(c.Query("page"), 1)
		length := atoiDefault(c.Query("length"), 10)
		rows, metadata, err := h.discountSvc.GetVouchers(ctx, userID, c.Query("code"), page, length)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": rows, "metadata": metadata})
	case "checkvoucher":
		// Pratinjau (TIDAK mengonsumsi) -- dipakai POS untuk tampilkan
		// potongan sebelum kasir input jumlah bayar. Konsumsi sungguhan
		// terjadi saat insertsales (modul sales), dengan goldid outlet
		// (bukan c.GetInt("user")) supaya konsisten kalau kasir belanja
		// lintas outlet -- di sini masih pakai goldid akun sendiri karena
		// ini murni pratinjau untuk kasir yang jualan di outlet miliknya.
		v, err := h.discountSvc.CheckVoucher(ctx, userID, c.Query("code"), c.Query("voucher"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": v})
	case "getvoucherhistory":
		page := atoiDefault(c.Query("page"), 1)
		length := atoiDefault(c.Query("length"), 10)
		rows, metadata, err := h.discountSvc.GetVoucherHistory(ctx, userID, c.Query("code"), page, length)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": rows, "metadata": metadata})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// InsertGoldGymDiscountGin menangani POST /v2/discount?type=insertdiscount
func (h *Handler) InsertGoldGymDiscountGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("InsertDiscount", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "insertdiscount":
		var body discountEntity.InsertDiscountData
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := c.GetInt("user")
		role := c.GetString("role")
		result, err := h.discountSvc.InsertDiscount(ctx, userID, role, body.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": result})
	case "insertvoucher":
		var body struct {
			Data discountEntity.InsertVoucher `json:"data"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := c.GetInt("user")
		role := c.GetString("role")
		code, err := h.discountSvc.InsertVoucher(ctx, userID, role, body.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Berhasil", "data": gin.H{"voucher_code": code}})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// UpdateGoldGymDiscountGin menangani PUT /v2/discount?type=updatediscount
func (h *Handler) UpdateGoldGymDiscountGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("UpdateDiscount", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "updatediscount":
		var body discountEntity.UpdateDiscountData
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := c.GetInt("user")
		role := c.GetString("role")
		result, err := h.discountSvc.UpdateDiscount(ctx, userID, role, body.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": result})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// DeleteGoldGymDiscountGin menangani DELETE /v2/discount?type=deletediscount&discountid=...
func (h *Handler) DeleteGoldGymDiscountGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("DeleteDiscount", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "deletediscount":
		discountID, err := strconv.Atoi(c.Query("discountid"))
		if err != nil || discountID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter discountid wajib diisi"})
			return
		}
		userID := c.GetInt("user")
		role := c.GetString("role")
		result, err := h.discountSvc.DeleteDiscount(ctx, userID, role, discountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": result})
	case "deletevoucher":
		voucherID, err := strconv.Atoi(c.Query("voucherid"))
		if err != nil || voucherID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter voucherid wajib diisi"})
			return
		}
		userID := c.GetInt("user")
		role := c.GetString("role")
		result, err := h.discountSvc.DeleteVoucher(ctx, userID, role, voucherID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": result})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}
