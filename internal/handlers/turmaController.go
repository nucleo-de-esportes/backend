package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository"
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
func CreateTurma(c *gin.Context) {

	var newTurma Turma

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != model.Admin {
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

	var local model.Local
	if err := repository.DB.Where("local_id = ?", newTurma.Local_Id).First(&local).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Local não encontrado"})
		return
	}

	var modalidade model.Modalidade
	if err := repository.DB.Where("modalidade_id = ?", newTurma.Modalidade_Id).First(&modalidade).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Modalidade não encontrada"})
		return
	}

	turmaModel := model.Turma{
		Horario_Inicio:  newTurma.Horario_Inicio,
		Horario_Fim:     newTurma.Horario_Fim,
		LimiteInscritos: newTurma.LimiteInscritos,
		Dia_Semana:      newTurma.Dia_Semana,
		Sigla:           newTurma.Sigla,
		Local_Id:        newTurma.Local_Id,
		Modalidade_Id:   newTurma.Modalidade_Id,
	}

	if err := repository.DB.Create(&turmaModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	turmaResponse := ConvertToTurmaResponse(newTurma, local.Nome, modalidade.Nome)

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
func DeleteTurma(c *gin.Context) {

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem deletar turmas."})
		return
	}

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Turma não encontrada"})
		return
	}

	if err := repository.DB.Where("turma_id = ?", turmaId).Delete(&model.Turma{}).Error; err != nil {
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
func GetTurmaById(c *gin.Context) {

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Turma não encontrada"})
		return
	}

	var turma model.Turma
	if err := repository.DB.Where("turma_id = ?", turmaId).First(&turma).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar a turma",
			"causa": err.Error(),
		})
		return
	}

	var local model.Local
	if err := repository.DB.Where("local_id = ?", turma.Local_Id).First(&local).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao tentar localizar local"})
		return
	}

	var modalidade model.Modalidade
	if err := repository.DB.Where("modalidade_id = ?", turma.Modalidade_Id).First(&modalidade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao tentar localizar modalidade"})
		return
	}

	var response TurmaResponse
	copier.Copy(&response, &turma)
	response.Local_nome = local.Nome
	response.Modalidade_nome = modalidade.Nome
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
func GetAllTurmas(c *gin.Context) {

	var turmas []model.Turma
	if err := repository.DB.Find(&turmas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar turmas",
			"causa": err.Error(),
		})
		return
	}

	var turmasResponse []TurmaResponse

	for _, turma := range turmas {

		var local model.Local
		if err := repository.DB.Where("local_id = ?", turma.Local_Id).First(&local).Error; err != nil {
			continue
		}

		var modalidade model.Modalidade
		if err := repository.DB.Where("modalidade_id = ?", turma.Modalidade_Id).First(&modalidade).Error; err != nil {
			continue
		}

		convertResponse := TurmaResponse{
			Turma_id:        turma.Turma_id,
			Horario_Inicio:  turma.Horario_Inicio,
			Horario_Fim:     turma.Horario_Fim,
			LimiteInscritos: turma.LimiteInscritos,
			Dia_Semana:      turma.Dia_Semana,
			Sigla:           turma.Sigla,
			Local_nome:      local.Nome,
			Modalidade_nome: modalidade.Nome,
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
func UpdateTurma(c *gin.Context) {

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Turma não encontrada"})
		return
	}

	var newTurma Turma

	if err := c.ShouldBindJSON(&newTurma); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais incorretas"})
		return
	}

	if err := repository.DB.Model(&model.Turma{}).Where("turma_id = ?", turmaId).Updates(map[string]interface{}{
		"horario_inicio":   newTurma.Horario_Inicio,
		"horario_fim":      newTurma.Horario_Fim,
		"limite_inscritos": newTurma.LimiteInscritos,
		"dia_semana":       newTurma.Dia_Semana,
		"sigla":            newTurma.Sigla,
		"local_id":         newTurma.Local_Id,
		"modalidade_id":    newTurma.Modalidade_Id,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao tentar atualizar turma",
			"causa": err.Error(),
		})
		return
	}

	var updatedTurma model.Turma
	if err := repository.DB.Where("turma_id = ?", turmaId).First(&updatedTurma).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar turma atualizada",
		})
		return
	}

	c.JSON(http.StatusOK, updatedTurma)
}

func GetNextClassById(c *gin.Context){

	turmaId, err := strv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Id da turma não encontrado",
			"details": err.Error(),

		})
		return
	}

	var aula model.Aula

	if err := repository.DB.Where("turma_id = ? AND data_hora > NOW()",turmaId).Order("data_hora ASC").First(&aula).Error; if err == "record not found"{
		c.JSON(http.StatusOK, gin.H{
			"message": "Nenhuma aula encontrada"
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao tentar encontrar proxima aula"
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, aula)

}
