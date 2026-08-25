package goldgym

import (
	"gold-gym-be/internal/config"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"

	goldCustomerEntity "gold-gym-be/internal/entity/customer"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) InsertGoldGymCustomerGin(c *gin.Context) {
	var (
		result         interface{}
		metadata       interface{}
		resp           response.Response
		err            error
		insertcustomer goldCustomerEntity.InsertCustomerData
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	types := c.Query("type")
	switch types {
	case "insertcustomer":
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&insertcustomer); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err = h.goldgymSvcCustomer.InsertCustomer(ctx, userID, insertcustomer.CustomerData)
		if err != nil {
			log.Println("err", err)
		}
	case "bulkinsertcustomer":
		// insert massal: dipublish ke Kafka (async) supaya request tidak menunggu
		// proses insert banyak baris. Worker memanggil InsertCustomer yang sama.
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&insertcustomer); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if len(insertcustomer.CustomerData) == 0 {
			c.JSON(400, gin.H{"error": "data customer kosong"})
			return
		}
		cfg, _ := config.Get()
		bulk := goldCustomerEntity.BulkInsertCustomer{GoldID: userID, Items: insertcustomer.CustomerData}
		if err := h.kafkaProd.Publish(ctx, cfg.Kafka.Topics.Sales, "insert_customers_bulk", bulk); err != nil {
			h.logger.For(ctx).Error("failed to publish bulk customer to kafka", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengantre insert massal"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message": "insert massal customer diantrekan",
			"status":  202,
			"count":   len(insertcustomer.CustomerData),
		})
		return
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
