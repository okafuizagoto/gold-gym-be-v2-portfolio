package selleraccess

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// GetGoldGymSellerAccessGin menangani GET /v2/selleraccess (ADMIN only)
//
//	type=list : daftar outlet + status menu Daftar Pembeli/Mode Pembeli
//	            milik penjualnya (name opsional, cari outlet ATAU penjual)
func (h *Handler) GetGoldGymSellerAccessGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("GetSellerAccess", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	if c.GetString("role") != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengakses"})
		return
	}

	switch c.Query("type") {
	case "list":
		rows, err := h.svc.GetAll(ctx, c.Query("name"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": rows})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}
