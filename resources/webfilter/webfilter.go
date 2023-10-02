package webfilter

import (
	"web-filter/models"

	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
)

func RetornoJSON(webFilters []models.WebFilter, c *gin.Context) []WebFilterReturn {
	responseData := make([]WebFilterReturn, len(webFilters))

	for i, wf := range webFilters {
		var data WebFilterReturn
		if err := json.Unmarshal([]byte(wf.Data), &data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Falha ao processar os dados do WebFilter",
			})
			return []WebFilterReturn{}
		}

		data.ID = wf.ID

		responseData[i] = data
	}

	return responseData
}
