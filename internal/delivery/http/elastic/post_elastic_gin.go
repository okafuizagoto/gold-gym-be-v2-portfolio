package elastic

import (
	"net/http"

	elasticEntity "gold-gym-be/internal/entity/elastic"

	"github.com/gin-gonic/gin"
)

// PostElasticGin handles POST /gold-gym/v2/elastic
// Query params:
//
//	?type=index&index=<index>
//
// Body: JSON UserDocument
func (h *Handler) PostElasticGin(c *gin.Context) {
	ctx := c.Request.Context()
	types := c.Query("type")
	index := c.Query("index")

	if index == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index is required"})
		return
	}

	switch types {
	case "index":
		var doc elasticEntity.UserDocument
		if err := c.ShouldBindJSON(&doc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		docID, err := h.elasticSvc.IndexUser(ctx, index, doc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":      docID,
				"message": "document indexed successfully",
			},
			"metadata": nil,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'index'"})
	}
}
