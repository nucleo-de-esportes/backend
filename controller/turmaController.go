package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
)

type Turma struct {
	Horario         time.Time `json:"horario"`
	LimiteInscritos int64     `json:"limite_inscritos"`
	Dia_Semana      string    `json:"dia_semana"`
	Sigla           string    `json:"sigla"`
	Local_Id        int64     `json:"local_id"`
	Modalidade_Id   int64     `json:"modalidade_id"`
}

func CreateTurma(c *gin.Context, supabase *supabase.Client) {

	var newTurma Turma

	if err := c.ShouldBindJSON(&newTurma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais incorretas"})
		return
	}

	var result []Turma
	err := supabase.DB.From("turma").Insert(newTurma).Execute(&result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Turma criada com sucesso",
		"data":    result,
	})

}
