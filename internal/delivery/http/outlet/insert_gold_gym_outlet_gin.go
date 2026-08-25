package goldgym

import (
	"fmt"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"

	goldOutletEntity "gold-gym-be/internal/entity/outlet"

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
// func (h *Handler) InsertGoldGymGin(w http.ResponseWriter, r *http.Request) {
func (h *Handler) InsertGoldGymOutletGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error

		resp response.Response

		inseroutlet goldOutletEntity.InsertOutletData
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	types := c.Query("type")
	switch types {
	case "insertoutlet":
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&inseroutlet); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			fmt.Println("error", err.Error())
			return
		}
		err = h.goldgymSvcOutlet.InsertOutlet(ctx, inseroutlet.InsertData, userID)
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
