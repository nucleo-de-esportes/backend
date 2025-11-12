package handlers

import (
	"net/http"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository"

)


type avisosRequest struct{
	Titulo string `json:"titulo" binding:"required"`
	Mensagem string `json:"mensagem" binding:"required"`
	tipo_destinatario string `json:"tipo_destinatario" binding:"required"`
	id_destinatario uint `json:"id_destinatario,omitempty" binding:"required"`

}

type avisosResponse struct{
	Id uint `json:"id_aviso"`
	Titulo string `json:"titulo"`
	Mensagem string `json:"mensagem"`
	Status string `json:"status"`
	Data_envio time.Time `json:"data_envio"`
}

func CreateAviso(c *gin.Context){

	var newAviso avisosRequest

	if err := c.ShouldBindJSON(&newAviso); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais incorretas"})
		return
	}

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if ((userType != model.Admin) && (userType != model.Professor)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores e professores podem acessar essa função."})
		return
	
	
	}

	turmaId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Id da turma não encontrado",
			"details": err.Error(),
		})
		return
	}

	if err := repository.DB.First(&turmaId).Error; err != nil{
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Nenhuma turma encontrada",
			})
		}
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao tentar encontrar turma",
			"details": err.Error(),
		})
		return
		
	}

	if err := repository.DB.Create(&newAviso).Error; err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error": "Erro ao criar novo aviso",
			"deatils": err.Error(),
		})
		return
	}


	avisoResponse := avisosResponse{
		Id: model.Avisos.Id,
		Titulo: newAviso.Titulo,
		Mensagem: newAviso.Mensagem,
		Data_envio: model.Avisos.Data_envio,
		Status: "enviado",
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Aviso criado!",
		"data": avisoResponse,
	})



}