package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nedpals/supabase-go"
	"github.com/nucleo-de-esportes/backend/model"
)

type Turma struct {
	Horario_Inicio  string `json:"horario_inicio"`
	Horario_Fim     string `json:"horario_fim"`
	LimiteInscritos int64  `json:"limite_inscritos"`
	Dia_Semana      string `json:"dia_semana"`
	Sigla           string `json:"sigla"`
	Local_Id        int64  `json:"local_id"`
	Modalidade_Id   int64  `json:"modalidade_id"`
}

type TurmaResponse struct {
	Horario_Inicio  string `json:"horario_inicio"`
	Horario_Fim     string `json:"horario_fim"`
	LimiteInscritos int64  `json:"limite_inscritos"`
	Dia_Semana      string `json:"dia_semana"`
	Sigla           string `json:"sigla"`
	Local_nome      string `json:"local"`
	Modalidade_nome string `json:"modalidade"`
}

type NomeResponse struct {
	Nome string `json:"nome"`
}

func ConvertToTurmaResponse(turma Turma, localNome string, modalidadeNome string) TurmaResponse {
	var response TurmaResponse

	copier.Copy(&response, &turma)
	response.Local_nome = localNome
	response.Modalidade_nome = modalidadeNome
	return response
}

// @Summary Cria uma nova turma
// @Description Cadastra uma nova turma no sistema, validando local e modalidade
// @Tags Turmas
// @Accept json
// @Produce json
// @Param turma body controller.Turma true "Dados da nova turma"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /turmas [post]
func CreateTurma(c *gin.Context, supabase *supabase.Client) {

	var newTurma Turma

	if err := c.ShouldBindJSON(&newTurma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais incorretas"})
		return
	}

	if newTurma.LimiteInscritos > 30 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Limite de 30 inscritos ultrapassado",
		})
		return
	}

	var validateLocal []model.Local
	var validadeModalidade []model.Modalidade

	localIdString := strconv.FormatInt(newTurma.Local_Id, 10)
	modalideIdString := strconv.FormatInt(newTurma.Modalidade_Id, 10)

	invalidId := supabase.DB.From("local").Select("*").Eq("local_id", localIdString).Execute(&validateLocal)
	invalidId2 := supabase.DB.From("modalidade").Select("*").Eq("modalidade_id", modalideIdString).Execute(&validadeModalidade)

	if invalidId != nil || len(validateLocal) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Local não encontrado"})
		return
	}

	if invalidId2 != nil || len(validadeModalidade) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Local não encontrado"})
		return
	}

	var result []Turma
	err := supabase.DB.From("turma").Insert(newTurma).Execute(&result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var localName []NomeResponse
	var modalidadeName []NomeResponse

	errLocating := supabase.DB.From("local").Select("nome").Eq("local_id", localIdString).Execute(&localName)
	if errLocating != nil || len(localName) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar nome do local"})
		return
	}

	errLocating2 := supabase.DB.From("modalidade").Select("nome").Eq("modalidade_id", modalideIdString).Execute(&modalidadeName)
	if errLocating2 != nil || len(modalidadeName) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar nome da modalidade"})
		return
	}

	turmaResponse := ConvertToTurmaResponse(newTurma, localName[0].Nome, modalidadeName[0].Nome)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Turma criada com sucesso",
		"data":    turmaResponse,
	})

}
// @Summary Deleta uma turma
// @Description Deleta uma turma com base no ID
// @Tags Turmas
// @Produce json
// @Param id path int true "ID da Turma"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /turmas/{id} [delete]
func DeleteTurma(c *gin.Context, supabase *supabase.Client) {

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Turma não encontrada"})
		return
	}

	turmaIdString := strconv.FormatInt(turmaId, 10)

	err = supabase.DB.From("turma").Delete().Eq("turma_id", turmaIdString).Execute(nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao deletar a turma",
			"causa": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Turma deletada com sucesso!",
		"deletedId": turmaId,
	})

}
// @Summary Busca turma por ID
// @Description Retorna uma turma com base no ID
// @Tags Turmas
// @Produce json
// @Param id path int true "ID da Turma"
// @Success 200 {object} []controller.Turma
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /turmas/{id} [get]
func GetTurmaById(c *gin.Context, supabase *supabase.Client) {

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Turma não encontrada"})
		return
	}

	turmaIdString := strconv.FormatInt(turmaId, 10)

	var viewTurma []Turma
	err = supabase.DB.From("turma").Select("*").Eq("turma_id", turmaIdString).Execute(&viewTurma)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar a turma",
			"causa": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, viewTurma)
}
// @Summary Lista todas as turmas
// @Description Retorna uma lista com todas as turmas cadastradas
// @Tags Turmas
// @Produce json
// @Success 200 {array} controller.TurmaResponse
// @Failure 500 {object} map[string]string
// @Router /turmas [get]
func GetAllTurmas(c *gin.Context, supabase *supabase.Client) {

	var turmas []Turma
	err := supabase.DB.From("turma").Select("*").Execute(&turmas)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar turmas",
			"causa": err.Error(),
		})
		return
	}

	var turmasResponse []TurmaResponse

	for _, turma := range turmas {

		localIdString := strconv.FormatInt(turma.Local_Id, 10)
		modalidadeIdString := strconv.FormatInt(turma.Modalidade_Id, 10)

		var localName []NomeResponse
		var modalidadeName []NomeResponse

		errLoc := supabase.DB.From("local").Select("nome").Eq("local_id", localIdString).Execute(&localName)
		errMod := supabase.DB.From("modalidade").Select("nome").Eq("modalidade_id", modalidadeIdString).Execute(&modalidadeName)

		if errLoc != nil || len(localName) == 0 || errMod != nil || len(modalidadeName) == 00 {
			continue
		}

		convertResponse := ConvertToTurmaResponse(turma, localName[0].Nome, modalidadeName[0].Nome)
		turmasResponse = append(turmasResponse, convertResponse)
	}

	c.JSON(http.StatusOK, turmasResponse)
}
// @Summary Atualiza uma turma
// @Description Atualiza os dados de uma turma existente
// @Tags Turmas
// @Accept json
// @Produce json
// @Param id path int true "ID da Turma"
// @Param turma body controller.Turma true "Dados atualizados da turma"
// @Success 200 {object} []controller.Turma
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /turmas/{id} [put]
func UpdateTurma(c *gin.Context, supabase *supabase.Client) {

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Turma não encontrada"})
		return
	}

	turmaIdString := strconv.FormatInt(turmaId, 10)

	var newTurma Turma

	if err := c.ShouldBindJSON(&newTurma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais incorretas"})
		return
	}

	var updatedTurma []Turma
	err2 := supabase.DB.From("turma").Update(newTurma).Eq("turma_id", turmaIdString).Execute(&updatedTurma)

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao tentar atualizar turma",
			"causa": err2.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedTurma)
}
