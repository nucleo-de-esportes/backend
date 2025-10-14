package model

import (
	"time"
)

type Aula struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TurmaID     int64     `gorm:"not null" json:"turma_id"`
	DataHora    time.Time `gorm:"not null" json:"data_hora"`
	DataHoraFim time.Time `gorm:"not null" json:"data_hora_fim"`
	CriadoEm    time.Time `gorm:"autoCreateTime" json:"criado_em"`

	Turma     Turma      `gorm:"foreignKey:Turma_id;constraint:OnDelete:CASCADE" json:"turma"`
	Presencas []Presenca `json:"presencas"`
}

func (Aula) TableName() string {
	return "aula"
}
