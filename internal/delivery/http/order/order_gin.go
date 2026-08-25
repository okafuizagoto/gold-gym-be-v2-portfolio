package order

import (
	"net/http"
	"strconv"

	"gold-gym-be/internal/config"
	orderEntity "gold-gym-be/internal/entity/order"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// GetGoldGymOrderGin menangani GET /v2/order
//
//	type=outlets      : daftar outlet (non-THERAPY) untuk pembeli pilih (name opsional)
//	type=catalog      : barang satu outlet (goldid, code, name opsional)
//	type=buyerorders  : daftar pesanan milik pembeli login (dashboard)
//	type=sellerorders : daftar pesanan masuk untuk penjual login (status opsional)
//	type=orderdetail  : detail satu pesanan (orderid)
func (h *Handler) GetGoldGymOrderGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("GetOrder", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "outlets":
		outlets, err := h.orderSvc.GetPublicOutlets(ctx, c.Query("name"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": outlets})
	case "alloutlets":
		// khusus ADMIN: semua outlet + penanda visible untuk pengaturan kurasi
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengakses"})
			return
		}
		outlets, err := h.orderSvc.GetAllOutletsForAdmin(ctx, c.Query("name"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": outlets})
	case "catalog":
		goldid, _ := strconv.Atoi(c.Query("goldid"))
		outcode := c.Query("code")
		items, err := h.orderSvc.GetOutletCatalog(ctx, goldid, outcode, c.Query("name"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": items})
	case "buyerorders":
		buyerID := c.GetInt("user")
		orders, err := h.orderSvc.GetBuyerOrders(ctx, buyerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": orders})
	case "sellerorders":
		goldid := c.GetInt("user")
		orders, err := h.orderSvc.GetSellerOrders(ctx, goldid, c.Query("status"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": orders})
	case "orderdetail":
		orderID := c.Query("orderid")
		if orderID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter orderid wajib diisi"})
			return
		}
		detail, err := h.orderSvc.GetOrderDetail(ctx, orderID, c.GetInt("user"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": detail})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// InsertGoldGymOrderGin menangani POST /v2/order?type=insertorder
func (h *Handler) InsertGoldGymOrderGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("InsertOrder", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "insertorder":
		var body orderEntity.InsertOrderData
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		buyerID := c.GetInt("user")
		order, err := h.orderSvc.InsertOrder(ctx, buyerID, body.InsertData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"data":    order,
			"message": "pesanan berhasil dibuat",
		})
	case "addvisible":
		// khusus ADMIN: tambahkan outlet ke daftar yang boleh dilihat pembeli
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengatur"})
			return
		}
		var body struct {
			GoldID  int    `json:"gold_id"`
			Outcode string `json:"outcode"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		addedBy := c.GetString("creator")
		if addedBy == "" {
			addedBy = strconv.Itoa(c.GetInt("user"))
		}
		if err := h.orderSvc.AddVisibleOutlet(ctx, body.GoldID, body.Outcode, addedBy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "outlet ditampilkan ke pembeli"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// DeleteGoldGymOrderGin menangani DELETE /v2/order?type=removevisible (ADMIN):
// mencabut satu outlet dari daftar yang boleh dilihat pembeli.
func (h *Handler) DeleteGoldGymOrderGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("DeleteOrder", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "removevisible":
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengatur"})
			return
		}
		goldid, _ := strconv.Atoi(c.Query("goldid"))
		outcode := c.Query("code")
		if err := h.orderSvc.RemoveVisibleOutlet(ctx, goldid, outcode); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "outlet disembunyikan dari pembeli"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// UpdateGoldGymOrderGin menangani PUT /v2/order
//
//	type=confirm : penjual konfirmasi (WAITING -> PROCESS)
//	type=reject  : penjual tolak (WAITING -> REJECT), body {"reason": "..."}
//	type=finish  : penjual selesai proses (PROCESS -> FINISH), nota dipublish ke Kafka
func (h *Handler) UpdateGoldGymOrderGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("UpdateOrder", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	orderID := c.Query("orderid")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter orderid wajib diisi"})
		return
	}
	sellerGoldID := c.GetInt("user")

	switch c.Query("type") {
	case "confirm":
		order, err := h.orderSvc.ConfirmOrder(ctx, orderID, sellerGoldID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": order, "message": "pesanan dikonfirmasi"})
	case "reject":
		var body struct {
			Reason string `json:"reason"`
		}
		_ = c.ShouldBindJSON(&body)
		if body.Reason == "" {
			body.Reason = c.Query("reason")
		}
		order, err := h.orderSvc.RejectOrder(ctx, orderID, sellerGoldID, body.Reason)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": order, "message": "pesanan ditolak"})
	case "finish":
		creator := c.GetString("creator")
		if creator == "" {
			creator = strconv.Itoa(sellerGoldID)
		}
		order, salePayload, err := h.orderSvc.FinishOrder(ctx, orderID, sellerGoldID, creator)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if salePayload != nil {
			cfg, _ := config.Get()
			if err := h.kafkaProd.Publish(ctx, cfg.Kafka.Topics.Sales, "insert_sales", salePayload.InsertData); err != nil {
				h.logger.For(ctx).Error("failed to publish order sale to kafka", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "pesanan selesai tapi antrian nota gagal, coba lagi"})
				return
			}
		}
		resp := gin.H{"data": order, "message": "pesanan selesai"}
		if order.OrderSaleID != nil {
			resp["sale_id"] = *order.OrderSaleID
		}
		c.JSON(http.StatusOK, resp)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}
