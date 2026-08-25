package goldgym

import (
	"gold-gym-be/pkg/response"
	"log"
	"net/http"

	// goldSaleEntity "gold-gym-be/internal/entity/sales"

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
func (h *Handler) UpdateGoldGymSaleGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
		types    string
		// updatecustomer goldSaleEntity.UpdateSaleData
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
	case "markpaid":
		// Admin/penjual menandai transaksi BELUM LUNAS menjadi LUNAS
		if c.GetString("role") == "BUYER" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya penjual/admin yang bisa mengubah status pembayaran"})
			return
		}
		saleid := c.Query("saleid")
		if saleid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parameter saleid wajib diisi"})
			return
		}
		userID := c.GetInt("user")
		sale, err := h.goldgymSvcSale.MarkSalePaid(ctx, userID, saleid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    sale,
			"message": "transaksi ditandai lunas",
		})
		return
	case "":
		// case "updatecustomer":
		// userID := c.GetInt("user")
		// body, _ := ioutil.ReadAll(c.Request.Body)
		// json.Unmarshal(body, &updatecustomer)
		// updatecustomer.UpdateData.CustGoldID = userID
		// result, err = h.goldgymSvcSale.UpdateSale(ctx, updatecustomer.UpdateData)
		// if err != nil {
		// 	log.Println("ERR", err)
		// }
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
