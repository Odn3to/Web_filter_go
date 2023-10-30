package controllers

import (
    "web-filter/models"
	"web-filter/database"
    "web-filter/resources/webfilter"
    "web-filter/resources/squid"

	"github.com/gin-gonic/gin"
	"net/http"
    "bytes"
	"encoding/json"
	"fmt"
    "os/exec"
)

// @Summary new filter - WebFilter
// @Description Cria filtro no WebFilter - Squid
// @ID newFilter
// @Param   Requisição     body    webfilter.WebFilterRequest     true        "Especificação do Filtro"
// @Success 200 {object} webfilter.Response
// @Router /webfilter/new [post]
func CriaWebFilter(c *gin.Context) {
	var webFilter models.WebFilter
	var request webfilter.WebFilterRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid credential",
        })
        return
    }

	webFilter.Data = fmt.Sprintf(`{"nome": "%s", "url": "%s"}`, request.Nome, request.URL)

	database.DB.Create(&webFilter)

	response := gin.H{
        "message": "WebFilter Salvo com sucesso!",
        "data":    webFilter,
    }

	c.JSON(http.StatusOK, response)
}

// @Summary Get - WebFilters
// @Description Busca os WebFilters - Squid
// @ID getFilters
// @Param   searchValue     path     string     false     "Valor da pesquisa"     default
// @Success 200 {object} webfilter.Response
// @Router /webfilter/search/:searchValue [get]
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

// @Summary Edit - WebFilters
// @Description Edita os WebFilters - Squid
// @ID editFilters
// @Param   id     path     string     false     "id WebFilter"     default
// @Success 200 {object} webfilter.Response
// @Router /webfilter/edit/:id [put]
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

// @Summary Delete - WebFilters
// @Description Deleta os WebFilters - Squid
// @ID deleteFilters
// @Param   id     path     string     false     "id WebFilter"     default
// @Success 200 {object} webfilter.Response
// @Router /webfilter/delete/:id [delete]
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

// @Summary Apply - WebFilters
// @Description Aplica as configurações do WebFilters - Squid
// @ID applyFilters
// @Param   searchValue     path     string     false     "Valor da pesquisa"     default
// @Success 200 {object} webfilter.Response
// @Router /webfilter/apply [get]
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

// @Summary Get status SQUID
// @Description Pega o status do Squid
// @ID getStatusSquid
// @Success 200 {object} webfilter.ResponseSquid
// @Router /webfilter/status [get]
func GetStatusSquid(c *gin.Context){
    cmd := exec.Command("systemctl", "is-active", "squid")
    output, err := cmd.CombinedOutput()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erro ao pegar status do squid",
        })
        return
    }

    text := "NÃO ATIVADO"
    classText := "alert-red"

    if string(output) == "active\n" {
        text = "ATIVADO"
        classText = "alert-green"
    }

    response := gin.H{
        "text":  text,
        "class": classText,
    }

    c.JSON(http.StatusOK, response)
}

func TokenValidationMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	tokenData := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	jsonData, err := json.Marshal(tokenData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating JSON body",
		})
        c.Abort()
		return
	}

	body := bytes.NewReader(jsonData)

	resp, err := http.Post("http://172.23.58.10/auth/login/validador", "application/json", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error verifying token",
		})
        c.Abort()
		return
	}
	defer resp.Body.Close()

	// Se o token não for válido, retorne um erro
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid token",
		})
        c.Abort()
		return
	}

	c.Next()
}