package model

type Modalidade struct {
	Modalidade_id int64   `json:"modalidade_id" `
	Nome          string  `json:"nome" `
	Valor         float64 `json:"valor" `
}
