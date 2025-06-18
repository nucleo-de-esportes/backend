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
	Turma_id        int64  `json:"turma_id"`
	Horario_Inicio  string `json:"horario_inicio"`
	Horario_Fim     string `json:"horario_fim"`
	LimiteInscritos int64  `json:"limite_inscritos"`
	Dia_Semana      string `json:"dia_semana"`
	Sigla           string `json:"sigla"`
	Local_nome      string `json:"local"`
	Modalidade_nome string `json:"modalidade"`
}

type TurmaGet struct {
	Turma
	Turma_id int64 `json:"turma_id"`
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

// CreateTurma godoc
// @Summary Cria uma nova turma
// @Description Cria uma nova turma com dados como horário, limite de inscritos, dia da semana, local e modalidade.
// @Tags Turmas
// @Accept json
// @Produce json
// @Param turma body controller.Turma true "Dados da nova turma"
// @Success 201 {object} map[string]interface{} "Turma criada com sucesso"
// @Failure 400 {object} map[string]interface{} "Credenciais incorretas | Limite de 30 inscritos ultrapassado | Local não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro ao buscar nome do local | Erro ao buscar nome da modalidade | Erro interno"
// @Router /turma [post]
func CreateTurma(c *gin.Context, supabase *supabase.Client) {

	var newTurma Turma

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem criar turmas."})
		return
	}

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

// DeleteTurma godoc
// @Summary Deleta uma turma
// @Description Deleta uma turma pelo ID
// @Tags Turmas
// @Param id path int true "ID da Turma"
// @Produce json
// @Success 200 {object} map[string]interface{} "Turma deletada com sucesso!"
// @Failure 400 {object} map[string]interface{} "Turma não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro ao deletar a turma"
// @Router /turma/{id} [delete]
func DeleteTurma(c *gin.Context, supabase *supabase.Client) {

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem criar turmas."})
		return
	}

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

// GetTurmaById godoc
// @Summary Busca turma por ID
// @Description Retorna os dados completos de uma turma com base no ID
// @Tags Turmas
// @Param id path int true "ID da Turma"
// @Produce json
// @Success 200 {object} TurmaResponse
// @Failure 400 {object} map[string]interface{} "Turma não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro ao buscar a turma | Erro ao tentar localizar local ou modalidade"
// @Router /turma/{id} [get]
func GetTurmaById(c *gin.Context, supabase *supabase.Client) {

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem criar turmas."})
		return
	}

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

	localIdString := strconv.FormatInt(viewTurma[0].Local_Id, 10)
	modalidadeIdString := strconv.FormatInt(viewTurma[0].Modalidade_Id, 10)

	var localName []NomeResponse
	var modalidadeName []NomeResponse

	errLoc := supabase.DB.From("local").Select("nome").Eq("local_id", localIdString).Execute(&localName)
	errMod := supabase.DB.From("modalidade").Select("nome").Eq("modalidade_id", modalidadeIdString).Execute(&modalidadeName)

	if errLoc != nil || len(localName) == 0 || errMod != nil || len(modalidadeName) == 00 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao tentar localizar local ou modalidade"})
		return
	}

	var response TurmaResponse
	copier.Copy(&response, &viewTurma[0])
	response.Local_nome = localName[0].Nome
	response.Modalidade_nome = modalidadeName[0].Nome
	response.Turma_id = turmaId

	c.JSON(http.StatusOK, response)
}

// GetAllTurmas godoc
// @Summary Lista todas as turmas
// @Description Retorna uma lista com todas as turmas cadastradas
// @Tags Turmas
// @Produce json
// @Success 200 {array} TurmaResponse
// @Failure 500 {object} map[string]interface{} "Erro ao buscar turmas"
// @Router /turma [get]
func GetAllTurmas(c *gin.Context, supabase *supabase.Client) {

	var turmas []TurmaGet
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

		convertResponse := TurmaResponse{
			Turma_id:        turma.Turma_id,
			Horario_Inicio:  turma.Horario_Inicio,
			Horario_Fim:     turma.Horario_Fim,
			LimiteInscritos: turma.LimiteInscritos,
			Dia_Semana:      turma.Dia_Semana,
			Sigla:           turma.Sigla,
			Local_nome:      localName[0].Nome,
			Modalidade_nome: modalidadeName[0].Nome,
		}
		turmasResponse = append(turmasResponse, convertResponse)

	}

	c.JSON(http.StatusOK, turmasResponse)
}

// UpdateTurma godoc
// @Summary Atualiza uma turma
// @Description Atualiza os dados de uma turma existente com base no ID
// @Tags Turmas
// @Accept json
// @Produce json
// @Param id path int true "ID da Turma"
// @Param turma body controller.Turma true "Dados atualizados da turma"
// @Success 200 {object} []Turma "Turma atualizada"
// @Failure 400 {object} map[string]interface{} "Credenciais incorretas | Turma não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro ao tentar atualizar turma"
// @Router /turma/{id} [put]
func UpdateTurma(c *gin.Context, supabase *supabase.Client) {

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem criar turmas."})
		return
	}

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
