package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nucleo-de-esportes/backend/internal/dto"
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

type TurmaComAlunosResponse struct {
	Turma  TurmaInfoResponse   `json:"turma"`
	Alunos []AlunoInfoResponse `json:"alunos"`
}

type TurmaInfoResponse struct {
	IdTurma    int64  `json:"id_turma"`
	Modalidade string `json:"modalidade"`
	Sigla      string `json:"sigla"`
	QtdAulas   int64  `json:"qtd_aulas"`
}

type AlunoInfoResponse struct {
	IdAluno   string `json:"id_aluno"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Presencas int64  `json:"presencas"`
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
	HorariosPermitidos []HorarioPermitido 
}

type HorarioPermitido struct {
	Inicio string
	Fim    string
}

var configuracoesModalidades = map[string]ConfiguracaoModalidade{
	"Clube de Corrida": {
		DiasSemana:    []string{"Segunda", "Quarta"},
		HorarioInicio: "06:30",
		HorarioFim:    "07:30",
		Local:         "Bloco 10 campus asa norte",
	},
	"Voleibol": {
		DiasSemana:    []string{"Terça", "Quinta"},
		HorarioInicio: "11:30",
		HorarioFim:    "12:30",
		Local:         "Bloco 10 campus asa norte",
	},
	"Defesa Pessoal": {
		DiasSemana:    []string{"Segunda", "Quarta"},
		HorarioInicio: "11:30",
		HorarioFim:    "12:30",
		Local:         "Ginásio bloco 4 campus asa norte",
	},
	"Mobilidade e Alongamento": {
		DiasSemana:    []string{"Terça", "Quinta"},
		HorarioInicio: "11:30",
		HorarioFim:    "12:30",
		Local:         "Bloco 10 campus asa norte",
	},
	"Natacao": {
		DiasSemana: []string{"segunda", "quarta"},
		Local:      "Piscina ao lado do bloco 10 campus asa norte",
		HorariosPermitidos: []HorarioPermitido{
			{Inicio: "11:00", Fim: "11:50"},
			{Inicio: "11:50", Fim: "12:40"},
		},
	},
	"Nado Livre": {
		DiasSemana: []string{"segunda", "quarta", "sexta"},
		Local:      "Piscina",

		HorariosPermitidos: []HorarioPermitido{}, 
	},
}

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

	var config ConfiguracaoModalidade
	var existe bool

	config, existe = configuracoesModalidades[modalidade.Nome]

	if !existe {
		for chave, valor := range configuracoesModalidades {
			if strings.EqualFold(strings.TrimSpace(chave), strings.TrimSpace(modalidade.Nome)) {
				config = valor
				existe = true
				break
			}
		}
	}

	if !existe {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Modalidade '%s' não possui configuração correspondente no sistema.", modalidade.Nome),
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
	} else if len(config.HorariosPermitidos) > 0 && !strings.EqualFold(modalidade.Nome, "Nado livre") {
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
	} else if strings.EqualFold(modalidade.Nome, "Nado livre") {
		if !services.ValidarHorarioNadoLivre(newTurma.Horario_Inicio, newTurma.Horario_Fim) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Horário deve estar entre 11:00 e 20:00",
			})
			return
		}
	}

	var count int64
	repository.DB.Model(&model.Turma{}).Where("modalidade_id = ?", newTurma.Modalidade_Id).Count(&count)

	prefix := ""
	words := strings.Split(modalidade.Nome, " ")
	for _, word := range words {
		if len(word) > 0 {
			prefix += string(word[0])
		}
	}
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}
	prefix = strings.ToUpper(prefix)
	newTurma.Sigla = fmt.Sprintf("%s-%03d", prefix, count+1)

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

func GetTurmaById(c *gin.Context) {

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != model.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem realizar esta ação."})
		return
	}

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Turma não encontrada"})
		return
	}

	var turma model.Turma
	if err := repository.DB.
		Preload("Professor").
		Where("turma_id = ?", turmaId).
		First(&turma).Error; err != nil {
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

	var response dto.TurmaResponse
	response.TurmaID = uint(turma.Turma_id)
	response.HorarioInicio = turma.Horario_Inicio
	response.HorarioFim = turma.Horario_Fim
	response.LimiteInscritos = int(turma.LimiteInscritos)
	response.DiaSemana = turma.Dia_Semana
	response.Sigla = turma.Sigla
	response.Local = dto.LocalResponseDto{Nome: local.Nome, Campus: local.Campus}
	response.Modalidade = dto.ModalidadeResponseDto{Nome: modalidade.Nome, ValorAluno: modalidade.Valor_aluno, ValorProfessor: modalidade.Valor_professor}
	var total int64
	repository.DB.Model(&model.Matricula{}).
		Where("turma_id = ?", turma.Turma_id).
		Count(&total)
	response.Total_alunos = int(total)

	professorName := ""
	if turma.Professor != nil {
		professorName = turma.Professor.Nome
	}

	response.Professor = professorName

	c.JSON(http.StatusOK, response)
}

func GetAllTurmas(c *gin.Context) {

	var turmas []model.Turma
	if err := repository.DB.
		Preload("Professor").
		Find(&turmas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar turmas",
			"causa": err.Error(),
		})
		return
	}

	var turmasResponse []dto.TurmaResponse

	for _, turma := range turmas {

		var local model.Local
		if err := repository.DB.Where("local_id = ?", turma.Local_Id).First(&local).Error; err != nil {
			continue
		}

		var modalidade model.Modalidade
		if err := repository.DB.Where("modalidade_id = ?", turma.Modalidade_Id).First(&modalidade).Error; err != nil {
			continue
		}

		var total int64
		repository.DB.Model(&model.Matricula{}).
			Where("turma_id = ?", turma.Turma_id).
			Count(&total)

		professorName := ""
		if turma.Professor != nil {
			professorName = turma.Professor.Nome
		}

		convertResponse := dto.TurmaResponse{
			TurmaID:         uint(turma.Turma_id),
			HorarioInicio:   turma.Horario_Inicio,
			HorarioFim:      turma.Horario_Fim,
			LimiteInscritos: int(turma.LimiteInscritos),
			Sigla:           turma.Sigla,
			Local:           dto.LocalResponseDto{Nome: local.Nome, Campus: local.Campus},
			Modalidade:      dto.ModalidadeResponseDto{Nome: modalidade.Nome, ValorAluno: modalidade.Valor_aluno, ValorProfessor: modalidade.Valor_professor},
			DiaSemana:       turma.Dia_Semana,
			Total_alunos:    int(total),
			Professor:       professorName,
		}
		turmasResponse = append(turmasResponse, convertResponse)
	}

	c.JSON(http.StatusOK, turmasResponse)
}

func UpdateTurma(c *gin.Context) {

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != model.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem realizar esta ação."})
		return
	}

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
			return 
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao tentar encontrar proxima aula",
			"details": err.Error(),
		})
		return 
	}

	c.JSON(200, aula)

}

func GetAlunosByTurmaId(c *gin.Context) {

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de turma inválido",
			"causa": err.Error(),
		})
		return
	}

	var turma model.Turma
	if err := repository.DB.
		Preload("Modalidade").
		Preload("Matriculas.User").
		Where("turma_id = ?", turmaId).
		First(&turma).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Turma não encontrada",
		})
		return
	}

	var qtdAulas int64
	if err := repository.DB.Model(&model.Aula{}).
		Where("turma_id = ? AND data_hora < NOW()", turmaId).
		Count(&qtdAulas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao contar aulas",
			"causa": err.Error(),
		})
		return
	}

	turmaInfo := TurmaInfoResponse{
		IdTurma:    turma.Turma_id,
		Modalidade: turma.Modalidade.Nome,
		Sigla:      turma.Sigla,
		QtdAulas:   qtdAulas,
	}

	alunos := make([]AlunoInfoResponse, 0, len(turma.Matriculas))

	for _, matricula := range turma.Matriculas {
		var presencas int64
		if err := repository.DB.
			Table("presencas").
			Joins("INNER JOIN aula ON presencas.aula_id = aula.id").
			Where("aula.turma_id = ? AND presencas.user_id = ? AND presencas.presente = true AND aula.data_hora < NOW()",
				turmaId, matricula.User_id).
			Count(&presencas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao contar presenças",
				"causa": err.Error(),
			})
			return
		}

		aluno := AlunoInfoResponse{
			IdAluno:   matricula.User.User_id.String(),
			Nome:      matricula.User.Nome,
			Email:     matricula.User.Email,
			Presencas: presencas,
		}
		alunos = append(alunos, aluno)
	}

	response := TurmaComAlunosResponse{
		Turma:  turmaInfo,
		Alunos: alunos,
	}

	c.JSON(http.StatusOK, response)
}

type AulaProfessorResponse struct {
	IdAula     uint   `json:"id_aula"`
	IdTurma    int64  `json:"id_turma"`
	Modalidade string `json:"modalidade"`
	Sigla      string `json:"sigla"`
	Local      string `json:"local"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
}

type AulasProfessorResponse struct {
	Aulas []AulaProfessorResponse `json:"aulas"`
}

func GetAulasByProfessor(c *gin.Context) {
	professorId := c.Param("id")
	if professorId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID do professor é obrigatório",
		})
		return
	}

	var professor model.User
	if err := repository.DB.Where("user_id = ? AND user_type = ?", professorId, model.Professor).First(&professor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Professor não encontrado",
		})
		return
	}

	diaParam := c.Query("dia")
	dataInicio, dataFim, err := services.ParseDateParam(diaParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de data inválido. Use 2006-01-02 ou timestamp Unix",
			"causa": err.Error(),
		})
		return
	}

	var turmas []model.Turma
	if err := repository.DB.Where("professor_id = ?", professorId).Find(&turmas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar turmas do professor",
			"causa": err.Error(),
		})
		return
	}

	turmaIds := make([]int64, len(turmas))
	for i, turma := range turmas {
		turmaIds[i] = turma.Turma_id
	}

	if len(turmaIds) == 0 {
		c.JSON(http.StatusOK, AulasProfessorResponse{Aulas: []AulaProfessorResponse{}})
		return
	}

	var aulas []model.Aula
	if err := repository.DB.
		Preload("Turma.Modalidade").
		Preload("Turma.Local").
		Where("turma_id IN ? AND data_hora >= ? AND data_hora < ?", turmaIds, dataInicio, dataFim).
		Order("data_hora ASC").
		Find(&aulas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar aulas",
			"causa": err.Error(),
		})
		return
	}

	aulasResponse := make([]AulaProfessorResponse, 0, len(aulas))
	for _, aula := range aulas {
		aulaResp := AulaProfessorResponse{
			IdAula:     aula.ID,
			IdTurma:    aula.TurmaID,
			Modalidade: aula.Turma.Modalidade.Nome,
			Sigla:      aula.Turma.Sigla,
			Local:      aula.Turma.Local.Nome,
			HoraInicio: aula.DataHora.Format("15:04"),
			HoraFim:    aula.DataHoraFim.Format("15:04"),
		}
		aulasResponse = append(aulasResponse, aulaResp)
	}

	response := AulasProfessorResponse{
		Aulas: aulasResponse,
	}

	c.JSON(http.StatusOK, response)
}