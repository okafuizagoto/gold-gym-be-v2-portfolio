package gymactivity

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// DeleteGymActivityGin godoc
// @Summary Delete a gym activity
// @Description Delete a gym activity document from MongoDB by ObjectID
// @Tags gymactivity
// @Produce  json
// @Param   id  query  string  true  "Document ObjectID (hex)"
// @Success 200 {object} map[string]string
// @Router /gold-gym/v2/activity [delete]
func (h *Handler) DeleteGymActivityGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)
	span := h.tracer.StartSpan("DeleteGymActivity", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info(
		"HTTP request received",
		zap.String("method", c.Request.Method),
		zap.Stringer("url", c.Request.URL),
	)

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'id' is required"})
		return
	}

	if err := h.gymactivitySvc.DeleteActivity(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.logger.For(ctx).Info(
		"HTTP request done",
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
		zap.String("deleted_id", id),
	)

	c.JSON(http.StatusOK, gin.H{"message": "activity deleted", "id": id})
}
