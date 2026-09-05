package storage

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// GetGoldGymStorageGin menangani GET /v2/storage -- ringkasan pemakaian +
// daftar foto (item + bukti pembayaran) milik user yang sedang login. Menu
// Storage TIDAK tersedia untuk role ADMIN.
func (h *Handler) GetGoldGymStorageGin(c *gin.Context) {
	ctx := c.Request.Context()
	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("GetStorage", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	if c.GetString("role") == "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin tidak memiliki kuota storage"})
		return
	}

	summary, err := h.storageSvc.GetSummary(ctx, c.GetInt("user"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// DeleteGoldGymStorageGin menangani DELETE /v2/storage?type=<item_photo|payment_proof>&id=<id>.
func (h *Handler) DeleteGoldGymStorageGin(c *gin.Context) {
	ctx := c.Request.Context()
	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("DeleteStorage", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	role := c.GetString("role")
	if role == "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin tidak memiliki kuota storage"})
		return
	}

	sourceType := c.Query("type")
	sourceID, _ := strconv.Atoi(c.Query("id"))
	userID := c.GetInt("user")
	deletedBy := c.GetString("creator")
	if deletedBy == "" {
		deletedBy = strconv.Itoa(userID)
	}

	if err := h.storageSvc.DeleteEntry(ctx, sourceType, sourceID, userID, false, deletedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "foto berhasil dihapus"})
}
