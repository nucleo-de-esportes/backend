package model

import (
	"github.com/google/uuid"
)

type Turma struct {
	Turma_id        int64      `json:"turma_id" gorm:"primaryKey;autoIncrement"`
	Horario_Inicio  string     `json:"horario_inicio"`
	Horario_Fim     string     `json:"horario_fim"`
	LimiteInscritos int64      `json:"limite_inscritos"`
	Dia_Semana      string     `json:"dia_semana"`
	Sigla           string     `json:"sigla"`
	Local_Id        int64      `json:"local_id"`
	Modalidade_Id   int64      `json:"modalidade_id"`
	Professor_id    *uuid.UUID `json:"professor_id" gorm:"type:uuid"`

	Local      Local       `json:"local" gorm:"foreignKey:Local_Id"`
	Modalidade Modalidade  `json:"modalidade" gorm:"foreignKey:Modalidade_Id"`
	Matriculas []Matricula `json:"matriculas" gorm:"foreignKey:Turma_id"`
	Professor  *User       `json:"professor" gorm:"foreignKey:Professor_id"`
}

func (Turma) TableName() string {
	return "turma"
}
