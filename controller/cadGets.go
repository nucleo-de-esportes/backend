package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
	"github.com/nucleo-de-esportes/backend/model"
)

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