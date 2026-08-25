package goldgym

import (
	"errors"
	"gold-gym-be/internal/entity"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) GetGoldGymItemsGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)
	span := h.tracer.StartSpan("GetGoldGym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	// id, _ := strconv.Atoi(c.Query("id"))

	types := c.Query("type")
	switch types {
	case "getallitems":
		userID := c.GetInt("user")
		limit, _ := strconv.Atoi(c.Query("page"))
		offset, _ := strconv.Atoi(c.Query("length"))
		result, metadata, err = h.goldgymSvcItems.GetItems(ctx, userID, c.Query("name"), c.Query("code"), limit, offset)
	}

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, entity.ErrInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, entity.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	// }

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", c.Request.Method, c.Request.URL)
	h.logger.For(ctx).Info(
		"HTTP request done",
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
	)
	c.JSON(200, resp)

	return
}
