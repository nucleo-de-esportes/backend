package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository"
)

// GetAllLocais godoc
// @Summary Lista todos os locais
// @Description Retorna uma lista com todos os locais cadastrados
// @Tags Cadastro
// @Success 200 {array} model.Local
// @Failure 500 {object} map[string]interface{} "Erro ao buscar locais"
// @Router /cad/local [get]
func GetAllLocais(c *gin.Context) {
	var locais []model.Local

	if err := repository.DB.Find(&locais).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar locais",
			"causa": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, locais)

}

// GetAllModalidades godoc
// @Summary Lista todas as modalidades
// @Description Retorna uma lista com todas as modalidades cadastradas
// @Tags Cadastro
// @Produce json
// @Success 200 {array} model.Modalidade
// @Failure 500 {object} map[string]interface{} "Erro ao buscar modalidades "
// @Router /cad/mod [get]
func GetAllModalidades(c *gin.Context) {
	var modalidades []model.Modalidade

	if err := repository.DB.Find(&modalidades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar modalidades",
			"causa": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, modalidades)
}
