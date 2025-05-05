package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UserType int

const (
	Professor UserType = iota
	Aluno
	Admin
)

func (u UserType) ConvertToString() string {
	types := map[UserType]string{
		Professor: "professor",
		Aluno:     "aluno",
		Admin:     "admin",
	}
	return types[u]
}

func ConvertToType(user_type string) (UserType, error) {

	user := strings.ToLower(user_type)
	types := map[string]UserType{
		"professor": Professor,
		"aluno":     Aluno,
		"admin":     Admin,
	}
	value, exists := types[user]
	if !exists {
		return 0, fmt.Errorf("tipo de usuário inválido: %s", user)
	}

	return value, nil

}

type User struct {
	User_id   uuid.UUID `json:"user_id"`
	User_type UserType  `json:"user_type"`
	Email     string    `json:"email"`
	Nome      string    `json:"nome"`
}
