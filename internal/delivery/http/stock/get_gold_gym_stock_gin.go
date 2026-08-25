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

// Getgoldgym godoc
// @Summary Get entries of all goldgyms
// @Description Get entries of all goldgyms
// @Tags goldgym
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200

// @Router /v1/profiles [get]
// func (h *Handler) GetGoldGym(w http.ResponseWriter, r *http.Request) {
func (h *Handler) GetGoldGymStockGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
		// types    string
	)
	// defer resp.RenderJSON(w, r)

	// spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	// span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))

	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)
	span := h.tracer.StartSpan("GetGoldGym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	// ctx := r.Context()
	ctx = opentracing.ContextWithSpan(ctx, span)
	// h.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	// Your code here
	// types = r.FormValue("type")
	types := c.Query("type")
	switch types {
	// stock -----------------------------------------------------------------------------------------------
	case "getonestock":
		result, err = h.goldgymSvcStock.GetOneStockProduct(ctx, c.Query("stockcode"), c.Query("stockname"), c.Query("stockid"))
		// stock -----------------------------------------------------------------------------------------------
		log.Printf("testDelivery %+v", result)
	case "getallstock":
		// result, err = h.goldgymSvcStock.GetAllStockHeader(ctx)
		userID := c.GetInt("user")
		limit, _ := strconv.Atoi(c.Query("page"))
		offset, _ := strconv.Atoi(c.Query("length"))
		result, metadata, err = h.goldgymSvcStock.GetStock(ctx, userID, c.Query("name"), c.Query("code"), limit, offset)
		// stock -----------------------------------------------------------------------------------------------
		// log.Printf("testDelivery %+v", result)
	case "getallstockredis":
		result, err = h.goldgymSvcStock.GetAllStockHeaderToRedis(ctx)
	case "getstockbyname":
		userID := c.GetInt("user")
		result, err = h.goldgymSvcStock.GetStockByName(ctx, userID, c.Query("name"))

	}

	// if err != nil {
	// resp = httpHelper.ParseErrorCode(err.Error())

	// // log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
	// log.Printf("[ERROR] %s %s - %v\n", c.Request.Method, c.Request.URL, err)
	// // h.logger.For(ctx).Error("HTTP request error", zap.String("method", r.Method), zap.Stringer("url", r.URL), zap.Error(err))
	// h.logger.For(ctx).Error(
	// 	"HTTP request error",
	// 	zap.String("method", c.Request.Method),
	// 	zap.String("url", c.Request.URL.String()),
	// 	zap.Error(err),
	// )
	// c.JSON(resp.StatusCode, resp)

	// c.JSON(http.StatusInternalServerError, gin.H{
	// 	"error": err.Error(),
	// })
	// return

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
	// log.Printf("[INFO] %s %s\n", r.Method, r.URL)
	log.Printf("[INFO] %s %s\n", c.Request.Method, c.Request.URL)
	h.logger.For(ctx).Info(
		"HTTP request done",
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
	)
	// h.logger.For(ctx).Info("HTTP request done", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	c.JSON(200, resp)

	return
}
