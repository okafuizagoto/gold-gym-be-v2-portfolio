package goldgym

import (
	"gold-gym-be/pkg/response"
	"log"

	goldMejaEntity "gold-gym-be/internal/entity/meja"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) InsertGoldGymMejaGin(c *gin.Context) {
	var (
		result     interface{}
		metadata   interface{}
		err        error
		resp       response.Response
		insertmeja goldMejaEntity.InsertMejaData
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("InsertGoldGymMeja", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	types := c.Query("type")
	switch types {
	case "insertmeja":
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&insertmeja); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = h.goldgymSvcMeja.InsertMejaBulk(ctx, userID, insertmeja.Outcode, insertmeja.InsertData)
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
