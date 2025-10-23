package model

import (
	"gorm.io/datatypes"
)

type UserType string

const (
	Professor UserType = "professor"
	Aluno     UserType = "aluno"
	Admin     UserType = "admin"
)

type User struct {
	User_id   datatypes.UUID `json:"user_id" gorm:"primaryKey;default:gen_random_uuid()"`
	User_type UserType       `json:"user_type"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Nome      string         `json:"nome" gorm:"not null"`
	Password  string         `json:"-" gorm:"not null"`

	Matriculas      []Matricula `json:"matriculas" gorm:"foreignKey:User_id"`
	TurmasProfessor []Turma     `json:"turmas_professor" gorm:"foreignKey:Professor_id"`
}

func (User) TableName() string {
	return "usuario"
}
