package elastic

import (
	"net/http"

	elasticEntity "gold-gym-be/internal/entity/elastic"

	"github.com/gin-gonic/gin"
)

// GetElasticGin handles GET /gold-gym/v2/elastic
// Query params:
//
//	?type=search&index=<index>&query=<query>
//	?type=getbyid&index=<index>&id=<id>
func (h *Handler) GetElasticGin(c *gin.Context) {
	ctx := c.Request.Context()
	types := c.Query("type")
	index := c.Query("index")

	if index == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index is required"})
		return
	}

	switch types {
	case "search":
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query is required for search"})
			return
		}

		docs, err := h.elasticSvc.SearchUsers(ctx, index, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": elasticEntity.SearchResult{
				Total: len(docs),
				Hits:  docs,
			},
			"metadata": nil,
		})

	case "getbyid":
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required for getbyid"})
			return
		}

		doc, err := h.elasticSvc.GetUserByID(ctx, index, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":     doc,
			"metadata": nil,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'search' or 'getbyid'"})
	}
}
