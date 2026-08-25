package goldgym

import (
	"gold-gym-be/pkg/response"
	"log"
	"net/http"

	goldCustomerTypeEntity "gold-gym-be/internal/entity/customertype"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) InsertGoldGymCustomerTypeGin(c *gin.Context) {
	var (
		result         interface{}
		metadata       interface{}
		resp           response.Response
		err            error
		insertcustomer goldCustomerTypeEntity.InsertCustomerTypeData
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	types := c.Query("type")
	switch types {
	case "insertcustomertype":
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&insertcustomer); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err = h.goldgymSvcCustomerType.InsertCustomerType(ctx, userID, insertcustomer.CustomerData)
		if err != nil {
			log.Println("err", err)
		}
	}

	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		resp.StatusCode = 500
		resp.Error.Code = 500
		log.Printf("[ERROR] %s %s - %s\n", c.Request.Method, c.Request.URL, err.Error())
		c.JSON(resp.StatusCode, resp)
		return
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", c.Request.Method, c.Request.URL)
	h.logger.For(ctx).Info("HTTP request done", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	c.JSON(200, resp)
	return
}
