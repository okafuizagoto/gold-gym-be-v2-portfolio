package gymactivity

import (
	"errors"
	"net/http"
	"strconv"

	"gold-gym-be/internal/entity"
	"gold-gym-be/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// GetGymActivityGin godoc
// @Summary Get gym activities
// @Description Get gym activity logs from MongoDB. Pass ?member_id=X to filter by member.
// @Tags gymactivity
// @Accept  json
// @Produce  json
// @Param   member_id  query  int  false  "Member ID (0 or omit = all)"
// @Success 200 {object} response.Response
// @Router /gold-gym/v2/activity [get]
func (h *Handler) GetGymActivityGin(c *gin.Context) {
	var (
		result interface{}
		err    error
		resp   response.Response
	)

	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)
	span := h.tracer.StartSpan("GetGymActivity", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info(
		"HTTP request received",
		zap.String("method", c.Request.Method),
		zap.Stringer("url", c.Request.URL),
	)

	memberID := 0
	if raw := c.Query("member_id"); raw != "" {
		memberID, _ = strconv.Atoi(raw)
	}

	result, err = h.gymactivitySvc.GetActivities(ctx, memberID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
	)

	resp.Data = result
	c.JSON(http.StatusOK, resp)
}
