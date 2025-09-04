package model

type Modalidade struct {
	Modalidade_id   int64   `json:"modalidade_id" gorm:"primaryKey;autoIncrement"`
	Nome            string  `json:"nome" gorm:"not null"`
	Valor_aluno     float64 `json:"valor_aluno" gorm:"not null"`
	Valor_professor float64 `json:"valor_professor" gorm:"not null"`
}

func (Modalidade) TableName() string {
	return "modalidade"
}
