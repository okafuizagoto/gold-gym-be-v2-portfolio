package goldgym

import (
	"encoding/json"
	httpHelper "gold-gym-be/internal/delivery/http"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/response"
	"io/ioutil"
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
// func (h *Handler) UpdateGoldGymGin(w http.ResponseWriter, r *http.Request) {
func (h *Handler) UpdateGoldGymGin(c *gin.Context) {
	var (
		result             interface{}
		metadata           interface{}
		err                error
		resp               response.Response
		types              string
		updategoldsubsuser goldEntity.UpdateSubs
		updatepassword     goldEntity.UpdatePassword
		updatenama         goldEntity.UpdateNama
		updatekartu        goldEntity.UpdateKartu
		logout             goldEntity.Logout
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
	case "updatesubsuser":
		body, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body, &updategoldsubsuser)
		result, err = h.goldgymSvc.UpdateSubscriptionDetail(ctx, updategoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}
	case "updatepassword":
		body, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body, &updatepassword)
		result, err = h.goldgymSvc.UpdateDataPeserta(ctx, updatepassword)
		if err != nil {
			log.Println("err", err)
		}
	case "updatenama":
		body, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body, &updatenama)
		result, err = h.goldgymSvc.UpdateNama(ctx, updatenama)
		if err != nil {
			log.Println("err", err)
		}
	case "updatekartu":
		body, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body, &updatekartu)
		result, err = h.goldgymSvc.UpdateKartu(ctx, updatekartu)
		if err != nil {
			log.Println("err", err)
		}
	case "logout":
		body, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(body, &logout)
		result, err = h.goldgymSvc.Logout(ctx, logout)
		if err != nil {
			log.Println("err", err)
		}
	case "updatevalidationemail":
		// body, _ := ioutil.ReadAll(c.Request.Body)
		// json.Unmarshal(body, &logout)
		// fmt.Println("Result :", logout)
		// fmt.Println("Result2 :", &logout)
		result, err = h.goldgymSvc.UpdateValidationOTP(ctx, c.Request.FormValue("otp"), c.Request.FormValue("email"))
		if err != nil {
			log.Println("err", err)
		}
	case "updateotp":
		// body, _ := ioutil.ReadAll(c.Request.Body)
		// json.Unmarshal(body, &logout)
		// fmt.Println("Result :", logout)
		// fmt.Println("Result2 :", &logout)
		result, err = h.goldgymSvc.UpdateOTP(ctx, c.Request.FormValue("email"))
		if err != nil {
			log.Println("err", err)
		}
	case "updateotpsubscription":
		// body, _ := ioutil.ReadAll(c.Request.Body)
		// json.Unmarshal(body, &logout)
		// fmt.Println("Result :", logout)
		// fmt.Println("Result2 :", &logout)
		// result, _, err = h.goldgymSvc.UpdateOTPSubscription(ctx, c.Request.FormValue("email"))
		result, err = h.goldgymSvc.UpdateOTPSubscription(ctx, c.Request.FormValue("email"))
		if err != nil {
			log.Println("err", err)
		}
	case "updatepaymentsubscription":
		result, err, resp = h.goldgymSvc.UpdatePayment(ctx, c.Request.FormValue("otp"), c.Request.FormValue("email"))
		// 	// case "":
	}

	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())
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

// UpdateGoldGymTokoGin menangani PUT /v2/userdata/toko (wajib token) —
// user login menyimpan nama toko sendiri; saat dia belanja sebagai pembeli,
// nota menampilkan keterangan pembeli dengan nama toko ini.
func (h *Handler) UpdateGoldGymTokoGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("UpdateToko", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	var body goldEntity.UpdateToko
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result, err := h.goldgymSvc.SetToko(ctx, c.GetInt("user"), body.GoldToko)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": result, "message": "nama toko tersimpan"})
}

// UpdateGoldGymBuyerGin menangani PUT /v2/userdata/buyer (wajib token) —
// konfirmasi "daftar sebagai pembeli": set flag gold_buyer_yn = Y sehingga
// menu Mode Pembeli muncul dan menu Daftar Pembeli hilang di aplikasi.
func (h *Handler) UpdateGoldGymBuyerGin(c *gin.Context) {
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("RegisterAsBuyer", ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	result, err := h.goldgymSvc.RegisterAsBuyer(ctx, c.GetInt("user"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": result, "message": "akun terdaftar sebagai pembeli"})
}
