package goldgym

import (
	"fmt"
	goldItemsEntity "gold-gym-be/internal/entity/items"
	"gold-gym-be/pkg/response"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) InsertGoldGymItemsGin(c *gin.Context) {
	var (
		result   interface{}
		metadata interface{}
		resp     response.Response
		err      error

		insertitems goldItemsEntity.InsertItemData
	)
	ctx := c.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", c.Request.Method), zap.Stringer("url", c.Request.URL))

	types := c.Query("type")
	switch types {
	case "insertitems":
		userID := c.GetInt("user")
		if err := c.ShouldBindJSON(&insertitems); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("insertitems", insertitems)
		result, err = h.goldgymSvcItems.InsertItems(ctx, userID, insertitems.ItemData, insertitems.ApplyAllOutlets)
		if err != nil {
			log.Println("err", err)
		}
	case "uploadphoto":
		// upload foto item (multipart: field "file", item_id lewat query/form)
		itemID, _ := strconv.Atoi(c.Query("item_id"))
		if itemID == 0 {
			itemID, _ = strconv.Atoi(c.PostForm("item_id"))
		}
		fileHeader, errFile := c.FormFile("file")
		if errFile != nil {
			c.JSON(400, gin.H{"error": "file foto wajib diupload"})
			return
		}
		if fileHeader.Size > goldItemsEntity.MaxItemPhotoBytes {
			c.JSON(400, gin.H{"error": "ukuran foto maksimal 2 MB"})
			return
		}
		f, errOpen := fileHeader.Open()
		if errOpen != nil {
			c.JSON(400, gin.H{"error": "file foto tidak bisa dibaca"})
			return
		}
		defer f.Close()
		content, errRead := io.ReadAll(io.LimitReader(f, goldItemsEntity.MaxItemPhotoBytes+1))
		if errRead != nil {
			c.JSON(400, gin.H{"error": "file foto tidak bisa dibaca"})
			return
		}
		item, errSave := h.goldgymSvcItems.SaveItemPhoto(ctx, itemID,
			fileHeader.Filename, fileHeader.Header.Get("Content-Type"), content,
			c.GetInt("user"), c.GetString("role") == "ADMIN")
		if errSave != nil {
			c.JSON(400, gin.H{"error": errSave.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"data":    item,
			"message": "foto item tersimpan",
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
