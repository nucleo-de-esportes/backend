package model

import (
	"time"

	"github.com/google/uuid"
)

type Matricula struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	User_id    uuid.UUID `json:"user_id" gorm:"not null"`
	Turma_id   int64     `json:"turma_id" gorm:"not null"`
	Created_At time.Time `json:"created_at" gorm:"autoCreateTime"`

	User  User  `json:"-" gorm:"foreignKey:User_id"`
	Turma Turma `json:"turma" gorm:"foreignKey:Turma_id"`
}

func (Matricula) TableName() string {
	return "matricula"
}
