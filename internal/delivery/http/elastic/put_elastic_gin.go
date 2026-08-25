package elastic

import (
	"net/http"

	elasticEntity "gold-gym-be/internal/entity/elastic"

	"github.com/gin-gonic/gin"
)

// PutElasticGin handles PUT /gold-gym/v2/elastic
// Query params:
//
//	?type=update&index=<index>&id=<id>
//
// Body: JSON partial UserDocument (hanya field yang ingin diubah)
func (h *Handler) PutElasticGin(c *gin.Context) {
	ctx := c.Request.Context()
	types := c.Query("type")
	index := c.Query("index")

	if index == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index is required"})
		return
	}

	switch types {
	case "update":
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required for update"})
			return
		}

		var doc elasticEntity.UserDocument
		if err := c.ShouldBindJSON(&doc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := h.elasticSvc.UpdateUser(ctx, index, id, doc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":      id,
				"message": "document updated successfully",
			},
			"metadata": nil,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'update'"})
	}
}
