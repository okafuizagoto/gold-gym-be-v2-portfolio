package goldgym

import (
	"bytes"
	"errors"
	"gold-gym-be/internal/entity"
	"gold-gym-be/pkg/response"
	"image"
	"image/png"
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
func (h *Handler) GetGoldGymGin(c *gin.Context) {
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
	case "getgoldgym":
		result, err = h.goldgymSvc.GetGoldUser(ctx)
		log.Println("deliverygolduser", result)
	case "getregistrationmode":
		// publik (dipanggil layar Register sebelum login) — mode saat ini
		// untuk pendaftaran mandiri: BOTH / BUYER_ONLY / SELLER_ONLY.
		mode, modeErr := h.goldgymSvc.GetRegistrationMode(ctx)
		if modeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": modeErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"mode": mode}})
		return
	case "golduserbyemail":
		result, err = h.goldgymSvc.GetGoldUserByEmail(ctx, c.Query("email"))
	case "allsubscription":
		result, err = h.goldgymSvc.GetAllSubscription(ctx)
	case "getuserandsubsdetail":
		result, err = h.goldgymSvc.GetSubsWithUser(ctx)
	case "gettotalpayment":
		result, err = h.goldgymSvc.GetSubscriptionHeaderTotalHarga(ctx, c.Query("email"))
	case "getimages":
		id, _ := strconv.Atoi(c.Query("id"))
		result, err = h.goldgymSvc.GetTestingImage(ctx, id)

		// Example image data in []byte (this should be replaced with your actual data source)
		imgData := []byte{ /* your PNG image data here */ }

		// Type assertion
		imgData, ok := result.([]byte)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid image data",
			})
			log.Fatal("The result is not of type []byte")
		}

		// Create a buffer from the image data
		imgBuffer := bytes.NewReader(imgData)

		// Decode the image data to get the image.Image object
		img, _, err := image.Decode(imgBuffer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to decode image",
			})
			return
		}

		// Set the appropriate header for PNG image
		c.Header("Content-Type", "image/png")

		// Encode the image to PNG format and write it to the response
		if err := png.Encode(c.Writer, img); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unable to encode image",
			})
			return
		}
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
