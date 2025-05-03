package model

import (
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

type User struct {
	User_id   uuid.UUID `json:"user_id"`
	User_type UserType  `json:"user_type"`
	Email     string    `json:"email"`
}
