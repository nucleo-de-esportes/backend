package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository"
	"github.com/nucleo-de-esportes/backend/internal/services"
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

type ConfiguracaoModalidade struct {
	DiasSemana         []string
	HorarioInicio      string
	HorarioFim         string
	Local              string
	HorariosPermitidos []HorarioPermitido // Pra modalidades com varios horarios como nado livre
}

type HorarioPermitido struct {
	Inicio string
	Fim    string
}

var configuracoesModalidades = map[string]ConfiguracaoModalidade{
	"Clube de corrida": {
		DiasSemana:    []string{"segunda", "quarta"},
		HorarioInicio: "06:30",
		HorarioFim:    "07:30",
		Local:         "Bloco 10 campus asa norte",
	},
	"Voleibol": {
		DiasSemana:    []string{"terça", "quinta"},
		HorarioInicio: "11:30",
		HorarioFim:    "12:30",
		Local:         "Bloco 10 campus asa norte",
	},
	"Defesa pessoal": {
		DiasSemana:    []string{"segunda", "quarta"},
		HorarioInicio: "11:30",
		HorarioFim:    "12:30",
		Local:         "Ginásio bloco 4 campus asa norte",
	},
	"Mobilidade e alongamento": {
		DiasSemana:    []string{"terça", "quinta"},
		HorarioInicio: "11:30",
		HorarioFim:    "12:30",
		Local:         "Bloco 10 campus asa norte",
	},
	"Natação": {
		DiasSemana: []string{"segunda", "quarta"},
		Local:      "Piscina ao lado do bloco 10 campus asa norte",
		HorariosPermitidos: []HorarioPermitido{
			{Inicio: "11:00", Fim: "11:50"},
			{Inicio: "11:50", Fim: "12:40"},
		},
	},
	"Nado livre": {
		DiasSemana: []string{"segunda", "quarta", "sexta"},
		Local:      "Piscina",

		HorariosPermitidos: []HorarioPermitido{}, // Validação diferente
	},
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

	//userType, exists := c.Get("user_type")
	//if !exists {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
	//	return
	//}

	//if userType != model.Admin {
	//	c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem criar turmas."})
	//	return
	//}

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

	// Verificar se a modalidade escolhida tem horários definidos
	config, existe := configuracoesModalidades[modalidade.Nome]
	if !existe {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Modalidade não possui horários configurados",
		})
		return
	}

	if len(config.HorariosPermitidos) == 0 && config.HorarioInicio != "" {

		if newTurma.Horario_Inicio != config.HorarioInicio || newTurma.Horario_Fim != config.HorarioFim {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Horário inválido para esta modalidade",
				"horario_correto": gin.H{
					"inicio": config.HorarioInicio,
					"fim":    config.HorarioFim,
				},
			})
			return
		}
	} else if len(config.HorariosPermitidos) > 0 && modalidade.Nome != "Nado livre" {
		// Verificar se o horário fornecido está na lista de horários permitidos
		horarioValido := false
		for _, h := range config.HorariosPermitidos {
			if newTurma.Horario_Inicio == h.Inicio && newTurma.Horario_Fim == h.Fim {
				horarioValido = true
				break
			}
		}
		if !horarioValido {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":                "Horário inválido para esta modalidade",
				"horarios_disponiveis": config.HorariosPermitidos,
			})
			return
		}
	} else if modalidade.Nome == "Nado livre" {
		// Validacao especial para Nado livre (11h as 20h)
		if !services.ValidarHorarioNadoLivre(newTurma.Horario_Inicio, newTurma.Horario_Fim) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Horário deve estar entre 11:00 e 20:00",
			})
			return
		}
	}

	// Garante que a turma so sera criada se a aula for criada tambem, caso a aula nao seja criada a turma tambem não será
	tx := repository.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	turmaModel := model.Turma{
		Horario_Inicio:  newTurma.Horario_Inicio,
		Horario_Fim:     newTurma.Horario_Fim,
		LimiteInscritos: newTurma.LimiteInscritos,
		Sigla:           newTurma.Sigla,
		Local_Id:        newTurma.Local_Id,
		Modalidade_Id:   newTurma.Modalidade_Id,
		Dia_Semana:      strings.Join(config.DiasSemana, ","),
	}

	if err := tx.Create(&turmaModel).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Criar aulas para TODOS os dias da semana da modalidade
	aulasIds := []uint{}
	for _, diaSemana := range config.DiasSemana {
		proximaDataAula := services.CalcularProximaAula(diaSemana)
		dataHoraInicio := services.CombinarDataHora(proximaDataAula, newTurma.Horario_Inicio)
		dataHoraFim := services.CombinarDataHora(proximaDataAula, newTurma.Horario_Fim)

		aulaModel := model.Aula{
			TurmaID:     turmaModel.Turma_id,
			DataHora:    dataHoraInicio,
			DataHoraFim: dataHoraFim,
			CriadoEm:    time.Now(),
		}

		if err := tx.Create(&aulaModel).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao criar aula para " + diaSemana + ": " + err.Error(),
			})
			return
		}

		aulasIds = append(aulasIds, aulaModel.ID)
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar dados"})
		return
	}

	turmaResponse := ConvertToTurmaResponse(newTurma, local.Nome, modalidade.Nome)
	c.JSON(http.StatusCreated, gin.H{
		"message":       "Turma criada com sucesso",
		"data":          turmaResponse,
		"aulas_criadas": len(aulasIds),
		"aulas_ids":     aulasIds,
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

func GetNextClassById(c *gin.Context) {

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Id da turma não encontrado",
			"details": err.Error(),
		})
		return
	}

	var aula model.Aula

	if err := repository.DB.Where("turma_id = ? AND data_hora > NOW()", turmaId).Order("data_hora ASC").First(&aula).Error; err != nil {

		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Nenhuma aula encontrada para esta turma",
			})
		}
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao tentar encontrar proxima aula",
			"details": err.Error(),
		})
		return
	}

	c.JSON(200, aula)

}
