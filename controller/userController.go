package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nedpals/supabase-go"
	supa "github.com/nedpals/supabase-go"
	"github.com/nucleo-de-esportes/backend/model"
	"github.com/nucleo-de-esportes/backend/services"
)

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	User_type string
	Nome      string `json:"nome" binding:"required"`
}

type RegisterResponse struct {
	Email     string    `json:"email"`
	User_id   uuid.UUID `json:"user_id"`
	User_type string    `json:"user_type"`
	Nome      string    `json:"nome"`
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
func RegsiterUser(c *gin.Context, supabase *supa.Client) {

	var data RegisterRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credenciais incorretas"})
		return
	}

	if validateEmail := services.ValidateEmail(data.Email); validateEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validateEmail.Error()})
		return

	}

	if strings.HasSuffix(data.Email, "@sempreceub.com") {
		data.User_type = "aluno"
	} else if strings.HasSuffix(data.Email, "@ceub.com") {
		data.User_type = "professor"
	}
	var user_type model.UserType

	user_type, err := model.ConvertToType(data.User_type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	user, err := supabase.Auth.SignUp(c, supa.UserCredentials{
		Email:    data.Email,
		Password: data.Password,
		Data: map[string]interface{}{
			"user_type": data.User_type,
			"nome":      data.Nome,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao tentar cadastrar usuario",
			"details": err.Error(),
		})
		return
	}

	newUser := model.User{
		User_id:   uuid.MustParse(user.ID),
		User_type: user_type,
		Email:     data.Email,
		Nome:      data.Nome,
	}

	userResponse := RegisterResponse{
		User_id:   uuid.MustParse(user.ID),
		Email:     data.Email,
		User_type: data.User_type,
		Nome:      data.Nome,
	}

	//var result []model.User
	insertUser := supabase.DB.From("usuario").Insert(newUser).Execute(nil)

	if insertUser != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao inserir usuario na tabela",
			"details": insertUser.Error(),
		})
		return
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
	User_id uuid.UUID `json:"user_id"`
	Email   string    `json:"email"`
	Nome    string    `json:"nome"`
	Token   string    `json:"token"`
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
func LoginUser(c *gin.Context, supabase *supa.Client) {

	var data LoginRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email ou senha incorretos"})
		return
	}

	login, err := supabase.Auth.SignIn(c, supa.UserCredentials{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Falha ao tentar autenticar usuário",
			"details": err.Error(),
		})
		return

	}

	var userData []model.User

	err = supabase.DB.From("usuario").Select("*").Eq("user_id", login.User.ID).Execute(&userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao tentar buscar informações do usuário",
			"details": err.Error()})
		return
	}

	response := LoginResponse{
		User_id: uuid.MustParse(login.User.ID),
		Email:   userData[0].Email,
		Nome:    userData[0].Nome,
		Token:   login.AccessToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login realizado com sucesso!",
		"usuario": response,
	})
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
func InscreverAluno(c *gin.Context, supabase *supabase.Client) {

	userID := c.GetString("user_id")

	var turmaId InscricaoRequest
	if err := c.ShouldBindJSON(&turmaId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	turmaIdString := strconv.FormatInt(turmaId.TurmaID, 10)

	var turma []map[string]interface{}
	err := supabase.DB.From("turma").Select("*").Eq("turma_id", turmaIdString).Execute(&turma)
	if err != nil || len(turma) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Turma não encontrada"})
		return
	}

	var inscricoes []map[string]interface{}
	err = supabase.DB.From("inscricao").Select("*").
		Eq("user_id", userID).Eq("turma_id", turmaIdString).
		Execute(&inscricoes)

	if err == nil && len(inscricoes) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Você já está inscrito nesta turma"})
		return
	}

	data := map[string]interface{}{
		"user_id":        userID,
		"turma_id":       turmaId.TurmaID,
		"created_at": time.Now(),
	}

	//var result map[string]interface{}
	err = supabase.DB.From("inscricao").Insert(data).Execute(nil)
	if err != nil {
    // Adicione esta linha para ver o erro detalhado no terminal onde o servidor está rodando
    	log.Printf("ERRO NO BANCO DE DADOS AO INSCREVER: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar inscrição"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Inscrição realizada com sucesso!"})
}
