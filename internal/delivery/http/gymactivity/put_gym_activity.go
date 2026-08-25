package gymactivity

import (
	"net/http"

	gymEntity "gold-gym-be/internal/entity/gymactivity"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// UpdateGymActivityGin godoc
// @Summary Update a gym activity
// @Description Update an existing gym activity document in MongoDB by ObjectID
// @Tags gymactivity
// @Accept  json
// @Produce  json
// @Param   id    query  string                         true  "Document ObjectID (hex)"
// @Param   body  body   gymEntity.UpdateActivityRequest true  "Fields to update"
// @Success 200 {object} map[string]string
// @Router /gold-gym/v2/activity [put]
func (h *Handler) UpdateGymActivityGin(c *gin.Context) {
	var req gymEntity.UpdateActivityRequest

	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)
	span := h.tracer.StartSpan("UpdateGymActivity", ext.RPCServerOption(spanCtx))
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.gymactivitySvc.UpdateActivity(ctx, id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	h.logger.For(ctx).Info(
		"HTTP request done",
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
		zap.String("updated_id", id),
	)

	c.JSON(http.StatusOK, gin.H{"message": "activity updated", "id": id})
}
