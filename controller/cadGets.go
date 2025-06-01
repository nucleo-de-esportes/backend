package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
	"github.com/nucleo-de-esportes/backend/model"
)

// GetAllLocais godoc
// @Summary Lista todos os locais
// @Description Retorna uma lista com todos os locais cadastrados
// @Tags Cadastro
// @Success 200 {array} model.Local
// @Failure 500 {object} map[string]interface{} "Erro ao buscar locais"
// @Router /cad/local [get]
func GetAllLocais(c *gin.Context, supabase *supabase.Client) {
	var locais []model.Local

	err := supabase.DB.From("local").Select("*").Execute(&locais)
	if err != nil {
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
func GetAllModalidades(c *gin.Context, supabase *supabase.Client) {
	var modalidades []model.Modalidade

	err := supabase.DB.From("modalidade").Select("*").Execute(&modalidades)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar modalidades",
			"causa": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, modalidades)
}
