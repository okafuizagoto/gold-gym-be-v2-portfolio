package booking

import (
	"net/http"
	"strconv"

	"gold-gym-be/internal/config"
	bookingEntity "gold-gym-be/internal/entity/booking"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// GetGoldGymBookingGin menangani GET /v2/booking
// type=slots   : grid slot per outlet per tanggal (code, date)
// type=buyers  : autocomplete pembeli terdaftar (name)
func (h *Handler) GetGoldGymBookingGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("GetBooking", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "slots":
		outcode := c.Query("code")
		date := c.Query("date")
		if outcode == "" || date == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter code dan date wajib diisi"})
			return
		}
		slots, err := h.bookingSvc.GetSlots(ctx, outcode, date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": slots})
	case "buyers":
		buyers, err := h.bookingSvc.SearchBuyers(ctx, c.Query("name"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": buyers})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// InsertGoldGymBookingGin menangani POST /v2/booking?type=insertbooking
// Booking berbayar → insert_sales dipublish ke Kafka (antrian yang sama dengan POS),
// sale_id dikembalikan supaya FE bisa langsung menarik nota PDF.
func (h *Handler) InsertGoldGymBookingGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("InsertBooking", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "insertbooking":
		var body bookingEntity.InsertBookingData
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID := c.GetInt("user")
		role := c.GetString("role")
		creator := c.GetString("creator")
		if creator == "" {
			creator = strconv.Itoa(userID)
		}

		booking, salePayload, err := h.bookingSvc.InsertBooking(ctx, userID, role, creator, body.InsertData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if salePayload != nil {
			cfg, _ := config.Get()
			if err := h.kafkaProd.Publish(ctx, cfg.Kafka.Topics.Sales, "insert_sales", salePayload.InsertData); err != nil {
				h.logger.For(ctx).Error("failed to publish booking sale to kafka", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "booking tersimpan tapi antrian pembayaran gagal, coba bayar ulang dari daftar booking"})
				return
			}
		}

		resp := gin.H{
			"data":    booking,
			"message": "booking berhasil dibuat",
		}
		if booking.BookingSaleID != nil {
			resp["sale_id"] = *booking.BookingSaleID
		}
		c.JSON(http.StatusCreated, resp)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// UpdateGoldGymBookingGin menangani PUT /v2/booking?type=paybooking&bookingid=...&itemid=...
func (h *Handler) UpdateGoldGymBookingGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("UpdateBooking", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "paybooking":
		bookingID := c.Query("bookingid")
		if bookingID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter bookingid wajib diisi"})
			return
		}
		itemID, _ := strconv.Atoi(c.Query("itemid"))
		customPrice, _ := strconv.Atoi(c.Query("customprice"))
		userID := c.GetInt("user")
		role := c.GetString("role")
		creator := c.GetString("creator")
		if creator == "" {
			creator = strconv.Itoa(userID)
		}

		booking, salePayload, err := h.bookingSvc.PayBooking(ctx, creator, role, bookingID, itemID, customPrice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if salePayload != nil {
			cfg, _ := config.Get()
			if err := h.kafkaProd.Publish(ctx, cfg.Kafka.Topics.Sales, "insert_sales", salePayload.InsertData); err != nil {
				h.logger.For(ctx).Error("failed to publish booking sale to kafka", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "status booking terupdate tapi antrian penjualan gagal"})
				return
			}
		}

		resp := gin.H{
			"data":    booking,
			"message": "booking lunas",
		}
		if booking.BookingSaleID != nil {
			resp["sale_id"] = *booking.BookingSaleID
		}
		c.JSON(http.StatusOK, resp)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}

// DeleteGoldGymBookingGin menangani DELETE /v2/booking?type=removebooking&bookingid=...
// Hanya penjual terapi atau admin yang bisa menghapus booking UNPAID;
// penghapusan dicatat di tabel booking_remove_log.
func (h *Handler) DeleteGoldGymBookingGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("DeleteBooking", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	switch c.Query("type") {
	case "removebooking":
		bookingID := c.Query("bookingid")
		if bookingID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter bookingid wajib diisi"})
			return
		}
		userID := c.GetInt("user")
		role := c.GetString("role")
		remover := c.GetString("creator")
		if remover == "" {
			remover = strconv.Itoa(userID)
		}

		if err := h.bookingSvc.RemoveBooking(ctx, role, remover, bookingID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "booking dihapus"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}
