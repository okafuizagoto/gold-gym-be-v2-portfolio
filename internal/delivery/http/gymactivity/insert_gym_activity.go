package gymactivity

import (
	"errors"
	"net/http"

	"gold-gym-be/internal/entity"
	gymEntity "gold-gym-be/internal/entity/gymactivity"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// InsertGymActivityGin godoc
// @Summary Insert a gym activity
// @Description Insert a new gym activity log into MongoDB
// @Tags gymactivity
// @Accept  json
// @Produce  json
// @Param   body  body  gymEntity.InsertActivityRequest  true  "Activity payload"
// @Success 201 {object} map[string]string
// @Router /gold-gym/v2/activity [post]
func (h *Handler) InsertGymActivityGin(c *gin.Context) {
	var req gymEntity.InsertActivityRequest

	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)
	span := h.tracer.StartSpan("InsertGymActivity", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info(
		"HTTP request received",
		zap.String("method", c.Request.Method),
		zap.Stringer("url", c.Request.URL),
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertedID, err := h.gymactivitySvc.InsertActivity(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	h.logger.For(ctx).Info(
		"HTTP request done",
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
		zap.String("inserted_id", insertedID),
	)

	c.JSON(http.StatusCreated, gin.H{"inserted_id": insertedID})
}
