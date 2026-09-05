package goldgym

import (
	"gold-gym-be/pkg/response"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gold-gym-be/internal/config"
	goldSaleEntity "gold-gym-be/internal/entity/sales"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (h *Handler) InsertGoldGymSaleGin(c *gin.Context) {
	var (
		result     interface{}
		metadata   interface{}
		resp       response.Response
		err        error
		insertsale goldSaleEntity.InsertSaleData
		cfg        *config.Config // Configuration object
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	types := c.Query("type")
	switch types {
	case "":
	case "uploadproof":
		// upload foto bukti pembayaran transfer bank (multipart: field "file")
		saleID := c.Query("saleid")
		if saleID == "" {
			saleID = c.PostForm("saleid")
		}
		fileHeader, errFile := c.FormFile("file")
		if errFile != nil {
			c.JSON(400, gin.H{"error": "file bukti pembayaran wajib diupload"})
			return
		}
		if fileHeader.Size > goldSaleEntity.MaxProofFileBytes {
			c.JSON(400, gin.H{"error": "ukuran foto maksimal 5 MB"})
			return
		}
		f, errOpen := fileHeader.Open()
		if errOpen != nil {
			c.JSON(400, gin.H{"error": "file bukti pembayaran tidak bisa dibaca"})
			return
		}
		defer f.Close()
		content, errRead := io.ReadAll(io.LimitReader(f, goldSaleEntity.MaxProofFileBytes+1))
		if errRead != nil {
			c.JSON(400, gin.H{"error": "file bukti pembayaran tidak bisa dibaca"})
			return
		}
		uploadedBy := c.GetString("creator")
		if uploadedBy == "" {
			uploadedBy = strconv.Itoa(c.GetInt("user"))
		}
		proof, errSave := h.goldgymSvcSale.SavePaymentProof(ctx, saleID,
			fileHeader.Filename, fileHeader.Header.Get("Content-Type"), content, uploadedBy,
			c.GetInt("user"), c.GetString("role") == "ADMIN")
		if errSave != nil {
			c.JSON(400, gin.H{"error": errSave.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"data":    proof,
			"message": "bukti pembayaran tersimpan",
		})
		return
	case "insertsales":
		cfg, _ = config.Get()
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&insertsale); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if len(insertsale.InsertData.Detail) == 0 {
			c.JSON(400, gin.H{"error": "detail penjualan kosong"})
			return
		}
		// insert diproses async oleh Kafka consumer — data invalid harus ditolak
		// di sini, kalau tidak transaksi hilang diam-diam di worker
		if strings.TrimSpace(insertsale.InsertData.Header.SaleOutcode) == "" {
			c.JSON(400, gin.H{"error": "sale_outcode wajib diisi"})
			return
		}
		for i, d := range insertsale.InsertData.Detail {
			if d.SaleStockid == nil || strings.TrimSpace(*d.SaleStockid) == "" {
				c.JSON(400, gin.H{"error": "sale_stockid item wajib diisi"})
				return
			}
			if d.SaleQty == nil || *d.SaleQty <= 0 {
				c.JSON(400, gin.H{"error": "qty item ke-" + strconv.Itoa(i+1) + " harus lebih dari 0"})
				return
			}
		}
		insertsale.InsertData.Header.SaleGoldID = userID
		// waktu transaksi manual hanya untuk ADMIN; role lain selalu live
		if c.GetString("role") != "ADMIN" {
			insertsale.InsertData.TransDate = ""
			insertsale.InsertData.TransTime = ""
		}
		// tenant nota = pemilik outlet tempat transaksi. Di-resolve dari outcode
		// untuk SEMUA role supaya penjual/admin yang belanja di outlet lain
		// (mode pembeli) notanya tercatat di outlet tujuan, bukan miliknya.
		outletGoldID, errOutlet := h.goldgymSvcSale.ResolveOutletGoldID(ctx, insertsale.InsertData.Header.SaleOutcode)
		if errOutlet != nil || outletGoldID <= 0 {
			c.JSON(400, gin.H{"error": "outlet tidak ditemukan"})
			return
		}
		insertsale.InsertData.Header.SaleGoldID = outletGoldID
		// Customer WAJIB untuk POS di outlet RETAIL, kecuali outlet sudah diberi
		// akses "POS tanpa customer" oleh admin. THERAPY tidak wajib (dari booking).
		if c.GetString("role") != "BUYER" {
			required, _ := h.goldgymSvcSale.IsPosCustomerRequired(ctx, outletGoldID, insertsale.InsertData.Header.SaleOutcode)
			cust := ""
			if insertsale.InsertData.Header.SaleSalescustomer != nil {
				cust = strings.TrimSpace(*insertsale.InsertData.Header.SaleSalescustomer)
			}
			if required && cust == "" {
				c.JSON(400, gin.H{"error": "nama customer wajib diisi"})
				return
			}
		}
		// pembeli (BUYER): identitas pembeli disimpan di sale_cust_id (gold_id
		// dari data_peserta). Role lain (kasir/SELLER/ADMIN jualan walk-in)
		// tidak boleh percaya nilai dari FE (klien lama kirim 0) — 0 bukan
		// gold_id valid dan akan menabrak FK ke data_peserta, jadi paksa NULL.
		if c.GetString("role") == "BUYER" {
			buyerID := userID
			insertsale.InsertData.Header.SaleCustID = &buyerID
		} else {
			insertsale.InsertData.Header.SaleCustID = nil
		}
		// sale_id dibuat di sini supaya bisa dikembalikan ke client
		// walaupun insert diproses async oleh Kafka consumer
		if insertsale.InsertData.Header.SaleID == "" {
			insertsale.InsertData.Header.SaleID = uuid.New().String()
		}
		// Voucher divalidasi & DIKONSUMSI SEKARANG (sinkron, sebelum antre ke
		// Kafka) -- bukan resource yang bisa dipercaya dari klien seperti
		// diskon per-item/total (lihat komentar SaleVoucherCode). Kalau kode
		// tidak valid/sudah kedaluwarsa/sudah dipakai, tolak SEBELUM nota
		// tercatat, supaya kasir dapat konfirmasi langsung.
		if insertsale.InsertData.Header.SaleVoucherCode != nil &&
			strings.TrimSpace(*insertsale.InsertData.Header.SaleVoucherCode) != "" {
			role := c.GetString("role")
			actorName := strconv.Itoa(userID)
			percent, errRedeem := h.voucherSvc.RedeemVoucher(ctx, outletGoldID,
				insertsale.InsertData.Header.SaleOutcode,
				*insertsale.InsertData.Header.SaleVoucherCode,
				insertsale.InsertData.Header.SaleID, actorName, role)
			if errRedeem != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": errRedeem.Error()})
				return
			}
			amount := "0"
			if total, errTotal := decimal.NewFromString(insertsale.InsertData.Header.SaleTranstotal); errTotal == nil {
				amt := total.Mul(decimal.NewFromFloat(percent)).Div(decimal.NewFromInt(100))
				amount = amt.StringFixed(0)
			}
			percentStr := decimal.NewFromFloat(percent).StringFixed(2)
			insertsale.InsertData.Header.SaleVoucherPercent = &percentStr
			insertsale.InsertData.Header.SaleVoucherAmount = &amount
		}
		// Meja yang direservasi kasir (lihat modul meja) diterjemahkan jadi
		// nama SEBELUM publish, supaya "No Meja" ikut tersimpan di baris
		// th_sale yang sama saat consumer insert (lihat ThSale.SaleMejaNames).
		// Gagal resolve tidak menggagalkan nota -- cukup NULL, ConfirmSaleMeja
		// di bawah tetap jadi sumber cadangan (tabel sale_meja).
		if len(insertsale.InsertData.MejaIDs) > 0 {
			if names, errNames := h.goldgymSvcSale.GetMejaNamesByIDs(ctx,
				insertsale.InsertData.Header.SaleOutcode, insertsale.InsertData.MejaIDs); errNames == nil && len(names) > 0 {
				joined := strings.Join(names, ", ")
				insertsale.InsertData.Header.SaleMejaNames = &joined
			} else if errNames != nil {
				h.logger.For(ctx).Error("failed to resolve meja names", zap.Error(errNames))
			}
		}
		if err := h.kafkaProd.Publish(ctx, cfg.Kafka.Topics.Sales, "insert_sales", insertsale.InsertData); err != nil {
			h.logger.For(ctx).Error("failed to publish to kafka", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue request"})
			return
		}

		resp := gin.H{
			"message": "insert request queued via kafka",
			"status":  202,
			"sale_id": insertsale.InsertData.Header.SaleID,
		}
		// booking terapi yang ikut dibayar lewat nota ini ditandai PAID
		// dengan sale_id yang sama (satu nota gabungan booking + barang)
		if len(insertsale.InsertData.BookingIDs) > 0 {
			if _, errMark := h.goldgymSvcSale.MarkBookingsPaid(ctx, insertsale.InsertData.BookingIDs, insertsale.InsertData.Header.SaleID); errMark != nil {
				h.logger.For(ctx).Error("failed to mark bookings paid", zap.Error(errMark))
				resp["booking_warning"] = "nota tersimpan tapi status booking gagal diperbarui, tandai lunas dari menu booking"
			}
		}
		// meja yang direservasi kasir di picker POS (lihat modul meja)
		// dicatat sinkron di sini, sama filosofi seperti blok booking di
		// atas -- gagal simpan meja tidak menggagalkan nota, cuma warning.
		if len(insertsale.InsertData.MejaIDs) > 0 {
			if _, errMeja := h.goldgymSvcSale.ConfirmSaleMeja(ctx, insertsale.InsertData.MejaIDs,
				insertsale.InsertData.Header.SaleOutcode, insertsale.InsertData.Header.SaleID); errMeja != nil {
				h.logger.For(ctx).Error("failed to confirm sale meja", zap.Error(errMeja))
				resp["meja_warning"] = "nota tersimpan tapi data meja gagal dicatat, catat manual di menu Kelola Meja"
			}
		}
		c.JSON(http.StatusAccepted, resp)
		return

	case "insertsalesnew":
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&insertsale); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		insertsale.InsertData.Header.SaleGoldID = userID
		result, err = h.goldgymSvcSale.InsertSales(ctx, userID, insertsale.InsertData)
		if err != nil {
			log.Println("err", err)
		}
	case "savecustomeraccess":
		// ADMIN: simpan akses "POS tanpa customer" untuk sekumpulan outlet yang
		// sedang tampil di pencarian (per baris: optional true=diberi, false=dicabut).
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengatur"})
			return
		}
		var body goldSaleEntity.SavePosAccessRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		addedBy := c.GetString("creator")
		if addedBy == "" {
			addedBy = strconv.Itoa(c.GetInt("user"))
		}
		if err := h.goldgymSvcSale.SavePosCustomerAccess(ctx, body.Items, addedBy); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "akses POS tersimpan"})
		return
	case "saveproofglobal":
		// ADMIN: nyalakan/matikan fitur bukti pembayaran untuk SEMUA user.
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengatur"})
			return
		}
		var body struct {
			Enabled bool `json:"enabled"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		updatedBy := c.GetString("creator")
		if updatedBy == "" {
			updatedBy = strconv.Itoa(c.GetInt("user"))
		}
		if err := h.goldgymSvcSale.SetPaymentProofGlobal(ctx, body.Enabled, updatedBy); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "pengaturan global tersimpan"})
		return
	case "saveproofoutlets":
		// ADMIN: simpan status aktif/nonaktif fitur bukti pembayaran untuk
		// sekumpulan outlet yang sedang tampil di pencarian.
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengatur"})
			return
		}
		var body struct {
			Items []goldSaleEntity.ProofAccessOutletItem `json:"items"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := h.goldgymSvcSale.SaveProofAccessOutlets(ctx, body.Items); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "akses outlet tersimpan"})
		return
	case "saveproofusers":
		// ADMIN: simpan status aktif/nonaktif fitur bukti pembayaran untuk
		// sekumpulan user yang sedang tampil di pencarian.
		if c.GetString("role") != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "hanya admin yang boleh mengatur"})
			return
		}
		var body struct {
			Items []goldSaleEntity.ProofAccessUserItem `json:"items"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := h.goldgymSvcSale.SaveProofAccessUsers(ctx, body.Items); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "akses user tersimpan"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type: " + types})
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
