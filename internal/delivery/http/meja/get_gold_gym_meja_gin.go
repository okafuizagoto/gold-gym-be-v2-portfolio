package goldgym

import (
	"gold-gym-be/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) GetGoldGymMejaGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("GetGoldGymMeja", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	types := c.Query("type")
	switch types {
	case "listmeja":
		userID := c.GetInt("user")
		result, err = h.goldgymSvcMeja.GetMejaByOutlet(ctx, userID, c.Query("outcode"))
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Printf("[ERROR] %s %s - %v\n", c.Request.Method, c.Request.URL, err)
		h.logger.For(ctx).Error("HTTP request error", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL), zap.Error(err))
		return
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", c.Request.Method, c.Request.URL)
	h.logger.For(ctx).Info("HTTP request done", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	c.JSON(200, resp)
	return
}
