package selleraccess

import (
	"net/http"

	entity "gold-gym-be/internal/entity/selleraccess"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// UpdateGoldGymSellerAccessGin menangani PUT /v2/selleraccess (ADMIN only)
//
//	type=daftarpembeli : set status menu "Daftar Pembeli" 1 akun penjual
//	type=modepembeli   : set status menu "Mode Pembeli" 1 akun penjual
//
// body: {"gold_id": 1, "active": true}
func (h *Handler) UpdateGoldGymSellerAccessGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("UpdateSellerAccess", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	if c.GetString("role") != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengatur"})
		return
	}

	var body entity.SetMenuAccessRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.GoldID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gold_id wajib diisi"})
		return
	}

	switch c.Query("type") {
	case "daftarpembeli":
		if err := h.svc.SetDaftarPembeli(ctx, body.GoldID, body.Active); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "status menu Daftar Pembeli diperbarui"})
	case "modepembeli":
		if err := h.svc.SetModePembeli(ctx, body.GoldID, body.Active); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "status menu Mode Pembeli diperbarui"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
	}
}
