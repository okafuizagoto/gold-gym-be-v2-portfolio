package goldgym

import (
	"encoding/json"
	"gold-gym-be/pkg/response"
	"io/ioutil"
	"log"
	"net/http"

	goldItemsEntity "gold-gym-be/internal/entity/items"

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
// func (h *Handler) UpdateGoldGymGin(w http.ResponseWriter, r *http.Request) {
func (h *Handler) UpdateGoldGymItemsGin(c *gin.Context) {
	var (
		result      interface{}
		metadata    interface{}
		err         error
		resp        response.Response
		types       string
		updateitems goldItemsEntity.UpdateItemsData
	)
	// defer resp.RenderJSON(w, r)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	// ctx := r.Context()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	// Your code here
	types = c.Request.FormValue("type")
	switch types {
	// case "":
	case "updateitems":
		userID := c.GetInt("user")
		body, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body, &updateitems)
		updateitems.UpdateData.ItemsGoldID = userID
		result, err = h.goldgymSvcItems.UpdateItems(ctx, updateitems.UpdateData)
		if err != nil {
			log.Println("ERR", err)
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
