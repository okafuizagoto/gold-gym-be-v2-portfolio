package elastic

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteElasticGin handles DELETE /gold-gym/v2/elastic
// Query params:
//
//	?type=delete&index=<index>&id=<id>
func (h *Handler) DeleteElasticGin(c *gin.Context) {
	ctx := c.Request.Context()
	types := c.Query("type")
	index := c.Query("index")

	if index == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index is required"})
		return
	}

	switch types {
	case "delete":
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required for delete"})
			return
		}

		if err := h.elasticSvc.DeleteUser(ctx, index, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":      id,
				"message": "document deleted successfully",
			},
			"metadata": nil,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'delete'"})
	}
}
