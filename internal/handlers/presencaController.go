package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository"
)

// ConfirmarPresenca godoc
// @Summary Confirma presença de um aluno em uma aula
// @Description O aluno autenticado confirma presença em uma aula específica
// @Tags Presença
// @Accept json
// @Produce json
// @Param id path int true "ID da Aula"
// @Success 200 {object} map[string]interface{} "Presença confirmada com sucesso"
// @Failure 400 {object} map[string]interface{} "Aula inválida"
// @Failure 401 {object} map[string]interface{} "Usuário não autenticado"
// @Failure 404 {object} map[string]interface{} "Aula não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro ao registrar presença"
// @Security BearerAuth
// @Router /aulas/{id}/presenca [put]
func ConfirmarPresenca(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	aulaId := c.Param("id")
	var aula model.Aula
	if err := repository.DB.First(&aula, aulaId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aula não encontrada"})
		return
	}

	uid := uuid.MustParse(userID.(string))

	var presenca model.Presenca
	if err := repository.DB.
		Where("aula_id = ? AND user_id = ?", aula.ID, uid).
		First(&presenca).Error; err == nil {
		presenca.Presente = true
		if err := repository.DB.Save(&presenca).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar presença"})
			return
		}
	} else {
		novaPresenca := model.Presenca{
			AulaID:   aula.ID,
			UserID:   uid,
			Presente: true,
		}
		if err := repository.DB.Create(&novaPresenca).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar presença"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Presença confirmada com sucesso",
		"aula_id": aula.ID,
		"user_id": uid,
	})
}
