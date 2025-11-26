package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository"
)

type AvisoRequest struct {
	Titulo        string               `json:"titulo" binding:"required"`
	Mensagem      string               `json:"mensagem" binding:"required"`
	Destinatarios []model.Destinatario `json:"destinatarios" binding:"required"`
}

// Funcao esta funcionando, porem 'Destinatario' esta como string no model, deve ser provavelmente um enum, so deve enviar aviso para turma visto que o id vindo da url e de
// uma turma, necessita de ajustes futuros

func CreateAviso(c *gin.Context) {

	var req AvisoRequest

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != model.Admin && userType != model.Professor {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem acessar essa função."})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Credenciais incorretas", "details": err.Error()})
		return
	}

	turmaID, err := strconv.Atoi(c.Param("turma_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Id da turma inválido"})
		return
	}

	var turma model.Turma
	if err := repository.DB.First(&turma, turmaID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Turma não encontrada"})
		return
	}

	aviso := model.Aviso{
		Titulo:    req.Titulo,
		Mensagem:  req.Mensagem,
		Status:    "enviado",
		DataEnvio: time.Now(),
	}

	if err := repository.DB.Create(&aviso).Error; err != nil {
		c.JSON(500, gin.H{"error": "Erro ao salvar aviso"})
		return
	}

	for i := range req.Destinatarios {
		req.Destinatarios[i].DestinoID = uint(turmaID)
		req.Destinatarios[i].AvisoID = aviso.Id
		repository.DB.Create(&req.Destinatarios[i])
	}

	c.JSON(http.StatusCreated, gin.H{
		"id_aviso":   aviso.Id,
		"titulo":     aviso.Titulo,
		"mensagem":   aviso.Mensagem,
		"data_envio": aviso.DataEnvio,
		"status":     aviso.Status,
	})
}
