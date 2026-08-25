package goldgym

import (
	"fmt"
	"net/http"

	"gold-gym-be/internal/config"
	goldEntity "gold-gym-be/internal/entity/goldgym"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// InsertGoldGymKafkaGin menerima HTTP request lalu mem-publish payload ke Kafka.
// Response: 202 Accepted (async — proses dilanjutkan oleh Kafka consumer worker).
//
// Endpoint: POST /gold-gym/v2/userdata/kafka?type=insertuser
func (h *Handler) InsertGoldGymKafkaGin(c *gin.Context) {
	var (
		cfg *config.Config // Configuration object
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("InsertGoldGymKafka", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info("Kafka HTTP request received",
		zap.String("method", c.Request.Method),
		zap.Stringer("url", c.Request.URL),
	)

	types := c.Query("type")
	switch types {
	case "insertuser":
		var user goldEntity.GetGoldUsers
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.kafkaProd.Publish(ctx, cfg.Kafka.Topics.HTTPInsert, "insert_user", user); err != nil {
			fmt.Println("MASOK3")
			h.logger.For(ctx).Error("failed to publish to kafka", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue request"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message": "insert request queued via kafka",
			"status":  202,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type: " + types})
	}
}
