package goldgym

import (
	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/response"
	"log"

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
// func (h *Handler) DeleteGoldGymGin(w http.ResponseWriter, r *http.Request) {
func (h *Handler) DeleteGoldGymGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
		// types              string
		deletegoldsubsuser goldEntity.DeleteSubs
	)
	// defer resp.RenderJSON(w, r)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	// ctx := r.Context()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	// types = r.FormValue("type")
	types := c.Query("type")
	switch types {
	case "deletesubsuser":
		// body, _ := ioutil.ReadAll(r.Body)
		// json.Unmarshal(body, &deletegoldsubsuser)
		if err := c.ShouldBindJSON(&deletegoldsubsuser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err = h.goldgymSvc.DeleteSubscriptionHeader(ctx, deletegoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}
	}

	if err != nil {
		// resp = httpHelper.ParseErrorCode(err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		//
		log.Printf("[ERROR] %s %s - %v\n", c.Request.Method, c.Request.URL, err)
		h.logger.For(ctx).Error("HTTP request error", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL), zap.Error(err))
		return
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", c.Request.Method, c.Request.URL)
	h.logger.For(ctx).Info("HTTP request done", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	c.JSON(200, resp)
	return
}
