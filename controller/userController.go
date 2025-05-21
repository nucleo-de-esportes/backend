package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	supa "github.com/nedpals/supabase-go"
	"github.com/nucleo-de-esportes/backend/model"
	"github.com/nucleo-de-esportes/backend/services"
)

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	User_type string `json:"user_type" binding:"required"`
	Nome      string `json:"nome" binding:"required"`
}

type RegisterResponse struct {
	Email     string    `json:"email"`
	User_id   uuid.UUID `json:"user_id"`
	User_type string    `json:"user_type"`
	Nome      string    `json:"nome"`
}

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

	var result []model.User
	insertUser := supabase.DB.From("usuario").Insert(newUser).Execute(&result)

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
	User_id   uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	User_type string    `json:"user_type"`
	Nome      string    `json:"nome"`
	Token     string    `json:"token"`
}

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
		User_id:   uuid.MustParse(login.User.ID),
		Email:     userData[0].Email,
		User_type: userData[0].User_type.ConvertToString(),
		Nome:      userData[0].Nome,
		Token:     login.AccessToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login realizado com sucesso!",
		"usuario": response,
	})

}
