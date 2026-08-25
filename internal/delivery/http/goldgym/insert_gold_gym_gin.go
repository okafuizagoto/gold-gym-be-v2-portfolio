package goldgym

import (
	"fmt"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/response"
	"io/ioutil"
	"log"
	"net/http"

	"gold-gym-be/own-pkg/crypto"

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
func (h *Handler) InsertGoldGymGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error

		resp response.Response
		// types          string
		insertgoldloginuser      goldEntity.LogUser
		insertgoldsubsuser       goldEntity.InsertSubsAll
		insertgoldsubsuserdetail goldEntity.SubscriptionDetail
		refreshToken             string
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
	case "registerbuyer":
		var registerBuyer goldEntity.RegisterBuyerRequest
		if err := c.ShouldBindJSON(&registerBuyer); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		buyer, err := h.goldgymSvc.RegisterBuyer(ctx, registerBuyer)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{
			"data":    buyer,
			"message": "registrasi pembeli berhasil",
		})
		return
	case "loginuser":
		// body, _ := ioutil.ReadAll(r.Body)
		// json.Unmarshal(body, &insertgoldloginuser)
		// fmt.Println("Result :", insertgoldloginuser)
		// fmt.Println("Result2 :", &insertgoldloginuser)
		if err := c.ShouldBindJSON(&insertgoldloginuser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		host := c.Request.Host
		device := c.GetHeader("User-Agent")
		fmt.Println("MASOK-DELIVERY")
		result, metadata, refreshToken, err = h.goldgymSvc.LoginUser(ctx, insertgoldloginuser.GoldEmail, insertgoldloginuser.GoldPassword, host, device)
		if err != nil {
			log.Println("err", err)
		}
		c.SetCookie(
			"refresh_token",
			refreshToken,
			604800,
			"/",
			"",
			true,
			true,
		)
		fmt.Println("result :", result)
		fmt.Println("metadata :", metadata)
		fmt.Println("err :", err)
		fmt.Println("Result :", insertgoldloginuser)
	case "insertsubsuser":
		// body, _ := ioutil.ReadAll(c.Request.Body)
		// json.Unmarshal(body, &insertgoldsubsuser)
		if err := c.ShouldBindJSON(&insertgoldsubsuser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err = h.goldgymSvc.InsertSubscriptionUser(ctx, insertgoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}
	case "insertsubsuserdetail":
		// body, _ := ioutil.ReadAll(c.Request.Body)
		// json.Unmarshal(body, &insertgoldsubsuserdetail)
		if err := c.ShouldBindJSON(&insertgoldsubsuser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err, resp = h.goldgymSvc.InsertSubscriptionDetail(ctx, insertgoldsubsuserdetail)
		if err != nil {
			log.Println("err", err)
		}
	case "uploadimages":

		// err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		// if err != nil {
		// 	http.Error(w, "Unable to parse form", http.StatusBadRequest)
		// 	return
		// }
		// -- body, _ := ioutil.ReadAll(r.Body)
		// Retrieve the file from the form field named "image"
		// file, _, err := r.FormFile("file")
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			// http.Error(w, "Unable to get file", http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get file"})
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			// http.Error(w, "Unable to read file", http.StatusInternalServerError)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get file"})
			return
		}

		// fmt.Println("fileBytes", fileBytes)

		// fmt.Println("files2", file)

		test := goldEntity.Testings{
			ID:            c.Request.FormValue("id"),
			TestingImages: fileBytes,
		}

		result, err = h.goldgymSvc.UploadTestingImages(ctx, test)
		if err != nil {
			log.Println("err", err)
		}

	case "testapi":
		fmt.Println("MASOK HANDLER")
		serviceCodeAuth := "401"
		// Get values from headers (matching PHP format)
		timestamp := c.GetHeader("X-TIMESTAMP")
		clientKey := c.GetHeader("X-CLIENT-KEY")
		privateKey := c.GetHeader("Private_Key")

		// Validate required headers
		if timestamp == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"ResponseCode":    serviceCodeAuth + "02",
				"ResponseMessage": "Missing X-TIMESTAMP header",
			})
			return
		}

		if clientKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"ResponseCode":    serviceCodeAuth + "02",
				"ResponseMessage": "Missing X-CLIENT-KEY header",
			})
			return
		}

		if privateKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"ResponseCode":    serviceCodeAuth + "02",
				"ResponseMessage": "Missing Private_Key header",
			})
			return
		}

		// Generate signature: {clientKey}|{timestamp} signed with RSA
		payload := clientKey + "|" + timestamp
		signature, err := crypto.RSASign(privateKey, payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"ResponseCode":    serviceCodeAuth + "02",
				"ResponseMessage": "Failed to generate signature: " + err.Error(),
			})
			return
		}

		// Return signature only (matching PHP format)
		c.JSON(http.StatusOK, gin.H{
			"signature": signature,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'type' tidak dikenali atau tidak diisi"})
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
