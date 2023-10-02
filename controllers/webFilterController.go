package controllers

import (
    "web-filter/models"
	"web-filter/database"
    "web-filter/resources/webfilter"
    "web-filter/resources/squid"

	"github.com/gin-gonic/gin"
	"net/http"

	"fmt"
)

func CriaWebFilter(c *gin.Context) {
	var webFilter models.WebFilter
	var request webfilter.WebFilterRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid credential",
        })
        return
    }
	//fmt.Printf("Resqyest: %+v\n", request)

	webFilter.Data = fmt.Sprintf(`{"nome": "%s", "url": "%s"}`, request.Nome, request.URL)

	database.DB.Create(&webFilter)

	response := gin.H{
        "message": "WebFilter Salvo com sucesso!",
        "data":    webFilter,
    }

	c.JSON(http.StatusOK, response)
}

func PesquisaWebFilter(c *gin.Context) {
    value := c.Params.ByName("searchValue")
	var webFilters []models.WebFilter

    // Recupera todos os registros da tabela WebFilter
    if value != "" {
        database.DB.Where("data->>'nome' LIKE ? or data->>'url' LIKE ?", "%"+value+"%", "%"+value+"%").Find(&webFilters)
    } else {
        database.DB.Find(&webFilters)
    }

    responseData := webfilter.RetornoJSON(webFilters, c)

    response := gin.H{
        "message": "WebFilters recuperados com sucesso!",
        "data": responseData,
    }

    c.JSON(http.StatusOK, response)
}

func EditarWebFilter(c *gin.Context) {
    id := c.Params.ByName("id")
    
    // Verifique se o registro com o ID especificado existe
    var webFilter models.WebFilter
    if err := database.DB.First(&webFilter, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Registro não encontrado",
        })
        return
    }

    var request webfilter.WebFilterRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Credencial inválida",
        })
        return
    }

    // Atualize os campos do registro com os valores da solicitação
    webFilter.Data = fmt.Sprintf(`{"nome": "%s", "url": "%s"}`, request.Nome, request.URL)

    // Salve as alterações no registro
    if err := database.DB.Save(&webFilter).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Falha ao atualizar o registro",
        })
        return
    }

    response := gin.H{
        "message": "WebFilter atualizado com sucesso!",
        "data":    webFilter,
    }

    c.JSON(http.StatusOK, response)
}

func DeleteWebFilter(c *gin.Context) {
    id := c.Params.ByName("id")

    // Verifique se o registro com o ID especificado existe
    var webFilter models.WebFilter
    if err := database.DB.First(&webFilter, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Registro não encontrado",
        })
        return
    }

    // Exclua o registro
    if err := database.DB.Delete(&webFilter).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Falha ao excluir o registro",
        })
        return
    }

    response := gin.H{
        "message": "WebFilter excluído com sucesso!",
        "data":    id,
    }

    c.JSON(http.StatusOK, response)
}

func ApplyWebFilter(c *gin.Context) {
    var webFilters []models.WebFilter

    // Recupera todos os registros da tabela WebFilter
    database.DB.Find(&webFilters)

    //configura o squid
    squid.ConfiguradorSquid(webFilters)

    response := gin.H{
        "message": "WebFilter Apply com sucesso!",
    }

    c.JSON(http.StatusOK, response)
}