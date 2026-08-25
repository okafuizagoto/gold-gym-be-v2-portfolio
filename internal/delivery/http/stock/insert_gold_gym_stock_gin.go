package goldgym

import (
	"gold-gym-be/internal/entity/firebase"
	goldStockEntity "gold-gym-be/internal/entity/stock"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"

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
func (h *Handler) InsertGoldGymStockGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error

		resp response.Response
		// types          string
		// insertgoldsubsuser goldEntity.InsertSubsAll
		insertUserFirebase firebase.User
		insertstock        goldStockEntity.InsertStockData
		// header                   http.Header
		// testings                 goldEntity.Testings
	)
	// defer resp.RenderJSON(w, r)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	// ctx := r.Context()
	ctx = opentracing.ContextWithSpan(ctx, span)
	// h.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	// Your code here
	// types = r.FormValue("type")
	types := c.Query("type")
	switch types {
	case "insertuserfirebase":
		// body, _ := ioutil.ReadAll(c.Request.Body)
		// json.Unmarshal(body, &insertUserFirebase)
		if err := c.ShouldBindJSON(&insertUserFirebase); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err = h.goldgymSvcStock.CreateUser(ctx, insertUserFirebase)
	// case "insertstock":
	// 	// body, _ := ioutil.ReadAll(c.Request.Body)
	// 	// json.Unmarshal(body, &insertstock)
	// 	if err := c.ShouldBindJSON(&insertgoldsubsuser); err != nil {
	// 		c.JSON(400, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	result, err = h.goldgymSvcStock.InsertStockSales(ctx, insertstock)
	// 	if err != nil {
	// 		log.Println("err", err)
	// 	}
	case "insertstocknew":
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&insertstock); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err = h.goldgymSvcStock.InsertStockSalesNew(ctx, userID, insertstock.StockData)
		if err != nil {
			log.Println("err", err)
		}
		// case "":
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
