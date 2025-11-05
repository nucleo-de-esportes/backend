package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/nucleo-de-esportes/backend/internal/dto"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository"
	"github.com/nucleo-de-esportes/backend/internal/services"
)

type RegisterRequest struct {
	Email     string         `json:"email" binding:"required"`
	Password  string         `json:"password" binding:"required"`
	User_type model.UserType `json:"user_type"`
	Nome      string         `json:"nome" binding:"required"`
}

type RegisterResponse struct {
	Email     string         `json:"email"`
	User_id   string         `json:"user_id"`
	User_type model.UserType `json:"user_type"`
	Nome      string         `json:"nome"`
}

// RegsiterUser godoc
// @Summary Registra um novo usuário
// @Description Cria um novo usuário com email, senha, tipo e nome
// @Tags Usuário
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Dados do novo usuário"
// @Success 201 {object} map[string]interface{} "Usuario cadastrado com sucesso"
// @Failure 400 {object} map[string]interface{} "Credenciais incorretas ou tipo de usuário inválido"
// @Failure 500 {object} map[string]interface{} "Erro ao tentar cadastrar usuario"
// @Router /user/register [post]

func RegisterUser(c *gin.Context) {

	var data RegisterRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais incorretas"})
		return
	}

	if validateEmail := services.ValidateEmail(data.Email); validateEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validateEmail.Error()})
		return

	}

	if strings.HasSuffix(data.Email, "@sempreceub.com") || strings.HasSuffix(data.Email, "@ceub.edu.br") {
		data.User_type = model.Aluno
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email não permitido para registro"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao processar senha",
		})
		return
	}

	newUser := model.User{
		User_id:   datatypes.UUID(uuid.New()),
		User_type: data.User_type,
		Email:     data.Email,
		Nome:      data.Nome,
		Password:  string(hashedPassword),
	}

	insertUser := services.CreateUser(newUser)

	if insertUser != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao inserir usuario na tabela",
			"details": insertUser.Error(),
		})
		return
	}

	userResponse := RegisterResponse{
		User_id:   newUser.User_id.String(),
		Email:     data.Email,
		User_type: data.User_type,
		Nome:      data.Nome,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario cadastrado com sucesso",
		"usuario": userResponse,
	})

}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User_id   string         `json:"user_id"`
	Email     string         `json:"email"`
	Nome      string         `json:"nome"`
	User_type model.UserType `json:"user_type"`
	Message   string         `json:"message"`
	Token     string         `json:"token"`
}

// LoginUser godoc
// @Summary Realiza login do usuário
// @Description Autentica um usuário existente e retorna token JWT e dados do usuário
// @Tags Usuário
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credenciais de login"
// @Success 200 {object} map[string]interface{} "Login realizado com sucesso!"
// @Failure 400 {object} map[string]interface{} "email ou senha incorretos"
// @Failure 401 {object} map[string]interface{} "Falha ao tentar autenticar usuário"
// @Failure 500 {object} map[string]interface{} "Erro ao tentar buscar informações do usuário"
// @Router /user/login [post]
func LoginUser(c *gin.Context) {
	var data LoginRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email ou senha incorretos",
		})
		return
	}

	var user model.User
	if err := repository.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Email ou senha incorretos",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro interno do servidor",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email ou senha incorretos",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject":   user.User_id.String(),
		"exp":       time.Now().Add(2 * time.Hour).Unix(),
		"user_type": user.User_type,
	})

	tokenString, err := token.SignedString([]byte(services.GetSecretKey()))
	if err != nil {
		log.Printf("Erro ao assinar token: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao tentar autenticar usuário",
		})
		return
	}

	response := LoginResponse{
		User_id:   user.User_id.String(),
		Email:     user.Email,
		Nome:      user.Nome,
		User_type: user.User_type,
		Message:   "Login realizado com sucesso!",
		Token:     tokenString,
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 12*3600, "", "", false, true)

	c.JSON(http.StatusOK, response)
}

type InscricaoRequest struct {
	TurmaID int64 `json:"turma_id"`
}

// InscreverAluno godoc
// @Summary Realiza a inscrição de um aluno em uma turma
// @Description Inscreve o usuário autenticado em uma turma específica
// @Tags Usuário
// @Accept json
// @Produce json
// @Param   inscricao    body     InscricaoRequest true "ID da Turma"
// @Success 201 {object} map[string]interface{} "Inscrição realizada com sucesso!"
// @Failure 401 {object} map[string]interface{} "Token ausente ou inválido"
// @Failure 404 {object} map[string]interface{} "Turma não encontrada"
// @Security BearerAuth
// @Router /user/inscricao [post]
func InscreverAluno(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	var request InscricaoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var turma model.Turma
	if err := repository.DB.First(&turma, request.TurmaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turma não encontrada"})
		return
	}

	var alreadyExists model.Matricula
	if err := repository.DB.Where("user_id = ? AND turma_id = ?", userID, request.TurmaID).
		First(&alreadyExists).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Você já está inscrito nesta turma"})
		return
	}

	inscricao := model.Matricula{
		User_id:    datatypes.UUID(uuid.MustParse(userID.(string))),
		Turma_id:   request.TurmaID,
		Created_At: time.Now(),
	}

	if err := repository.DB.Create(&inscricao).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar inscrição"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Inscrição realizada com sucesso!"})
}

func AtribuirProfessor(c *gin.Context) {
	var input struct {
		Turma_id     int64     `json:"turma_id"`
		Professor_id uuid.UUID `json:"professor_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if err := repository.DB.Model(&model.Turma{}).
		Where("turma_id = ?", input.Turma_id).
		Update("professor_id", input.Professor_id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar turma"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Professor atribuído à turma com sucesso"})
}

type InscricaoComTurma struct {
	Turma TurmaResponseUser `json:"turma"`
}

type TurmaResponseUser struct {
	Turma_id        int64        `json:"turma_id"`
	Horario_Inicio  string       `json:"horario_inicio"`
	Horario_Fim     string       `json:"horario_fim"`
	LimiteInscritos int64        `json:"limite_inscritos"`
	Dia_Semana      string       `json:"dia_semana"`
	Sigla           string       `json:"sigla"`
	Local           NomeResponse `json:"local"`
	Modalidade      NomeResponse `json:"modalidade"`
}

type UserResponse struct {
	User_id   string         `json:"user_id"`
	User_type model.UserType `json:"user_type"`
	Email     string         `json:"email"`
	Nome      string         `json:"nome"`

	Matriculas      []model.Matricula `json:"matriculas,omitempty"`
	TurmasProfessor []model.Turma     `json:"turmas_professor,omitempty"`
}

func GetTurmasByUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var matriculas []model.Matricula
	if err := repository.DB.
		Preload("Turma.Local").
		Preload("Turma.Modalidade").
		Where("user_id = ?", uuid.MustParse(userId.(string))).
		Find(&matriculas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar turmas", "detalhe": err.Error()})
		return
	}

	var turmasInscritas []dto.TurmaResponse
	var turmasMinistradas []dto.TurmaResponse

	for _, m := range matriculas {
		t := m.Turma

		var total int64
		repository.DB.Model(&model.Matricula{}).
			Where("turma_id = ?", t.Turma_id).
			Count(&total)

		turmasInscritas = append(turmasInscritas, dto.TurmaResponse{
			TurmaID:         uint(t.Turma_id),
			HorarioInicio:   t.Horario_Inicio,
			HorarioFim:      t.Horario_Fim,
			LimiteInscritos: int(t.LimiteInscritos),
			DiaSemana:       t.Dia_Semana,
			Sigla:           t.Sigla,
			Total_alunos:    int(total),
			Local: dto.LocalResponseDto{
				Nome:   t.Local.Nome,
				Campus: t.Local.Campus,
			},
			Modalidade: dto.ModalidadeResponseDto{
				Nome:           t.Modalidade.Nome,
				ValorAluno:     t.Modalidade.Valor_aluno,
				ValorProfessor: t.Modalidade.Valor_professor,
			},
		})
	}

	var turmasProfessor []model.Turma
	if err := repository.DB.
		Preload("Local").
		Preload("Modalidade").
		Where("professor_id = ?", userId).
		Find(&turmasProfessor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar turmas (professor)", "detalhe": err.Error()})
		return
	}

	for _, t := range turmasProfessor {
		var total int64
		repository.DB.Model(&model.Matricula{}).
			Where("turma_id = ?", t.Turma_id).
			Count(&total)

		turmasMinistradas = append(turmasMinistradas, dto.TurmaResponse{
			TurmaID:         uint(t.Turma_id),
			HorarioInicio:   t.Horario_Inicio,
			HorarioFim:      t.Horario_Fim,
			LimiteInscritos: int(t.LimiteInscritos),
			DiaSemana:       t.Dia_Semana,
			Sigla:           t.Sigla,
			Total_alunos:    int(total),
			Local: dto.LocalResponseDto{
				Nome:   t.Local.Nome,
				Campus: t.Local.Campus,
			},
			Modalidade: dto.ModalidadeResponseDto{
				Nome:           t.Modalidade.Nome,
				ValorAluno:     t.Modalidade.Valor_aluno,
				ValorProfessor: t.Modalidade.Valor_professor,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"turmas ministradas": turmasMinistradas,
		"turmas inscritas": turmasInscritas})
}

func GetUsers(c *gin.Context) {

	var users []model.User

	if err := repository.DB.
		Preload("Matriculas.Turma.Local").
		Preload("Matriculas.Turma.Modalidade").
		Preload("TurmasProfessor.Local").
		Preload("TurmasProfessor.Modalidade").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar usuários",
			"causa": err.Error(),
		})
		return
	}

	var result []UserResponse

	for _, user := range users {
		userResp := UserResponse{
			User_id:         user.User_id.String(),
			User_type:       user.User_type,
			Email:           user.Email,
			Nome:            user.Nome,
			Matriculas:      user.Matriculas,
			TurmasProfessor: user.TurmasProfessor,
		}

		switch user.User_type {
		case model.Aluno:
			userResp.Matriculas = user.Matriculas

		case model.Professor:
			userResp.TurmasProfessor = user.TurmasProfessor

		}
		result = append(result, userResp)

	}
	c.JSON(200, result)
}

func GetUserById(c *gin.Context) {

	var user model.User

	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Usuário não encontrado",
			"causa": err.Error()})
		return
	}

	if err := repository.DB.First(&user, "user_id = ?", userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar usuário",
			"causa": err.Error(),
		})
		return
	}

	userResp := UserResponse{
		User_id:         user.User_id.String(),
		User_type:       user.User_type,
		Email:           user.Email,
		Nome:            user.Nome,
		Matriculas:      user.Matriculas,
		TurmasProfessor: user.TurmasProfessor,
	}

	switch user.User_type {
	case model.Aluno:
		userResp.Matriculas = user.Matriculas

	case model.Professor:
		userResp.TurmasProfessor = user.TurmasProfessor

	}

	c.JSON(200, userResp)
}

func DeleteUserById(c *gin.Context) {

	var user model.User

	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Usuário não encontrado",
			"causa": err.Error()})
		return
	}

	if err := repository.DB.First(&user, "user_id = ?", userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar usuário",
			"causa": err.Error(),
		})
		return
	}

	if err := repository.DB.Delete(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuario nao encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro interno do servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario deletado com sucesso",
	})

}

func DeleteUserTurma(c *gin.Context) {

	userId, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Usuário não encontrado",
			"causa": err.Error()})
		return
	}

	turmaId := c.Param("turma_id")

	// Verifica se a matricula existe antes de deletar
	var matricula model.Matricula
	if err := repository.DB.Where("user_id = ? AND turma_id = ?", userId, turmaId).First(&matricula).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Matricula nao encontrada para este usuario e turma",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro interno do servidor",
		})
		return
	}

	if err := repository.DB.Where("user_id = ? AND turma_id = ?", userId, turmaId).Delete(&model.Matricula{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao remover usuário da turma",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuário removido da turma com sucesso!",
	})

}

type UpdateUserRequest struct {
	UserType model.UserType `json:"user_type"`
}

func UpdateUser(c *gin.Context) {

	userType, exists := c.Get("user_type")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tipo de usuário não encontrado"})
		return
	}

	if userType != model.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permissão negada. Apenas administradores podem acessar essa função."})
		return
	}

	userId := c.Param("id")

	var userResponse UserResponse
	var user model.User

	if err := repository.DB.First(&user, "user_id = ?", userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar usuário",
			"causa": err.Error(),
		})
		return
	}

	var newUserType UpdateUserRequest

	if err := c.ShouldBindJSON(&newUserType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais incorretas"})
		return
	}

	if err := repository.DB.Model(&user).Where("user_id = ?", userId).Update("user_type", newUserType.UserType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao tentar atualizar informações do usuário",
			"details": err.Error(),
		})
		return
	}
	userResponse.User_type = user.User_type
	userResponse.User_id = userId
	userResponse.Email = user.Email
	userResponse.Nome = user.Nome
	userResponse.Matriculas = user.Matriculas
	userResponse.TurmasProfessor = user.TurmasProfessor

	c.JSON(http.StatusOK, userResponse)
}
