package dto

type TurmaResponse struct {
	TurmaID         uint   `json:"turma_id"`
	HorarioInicio   string `json:"horario_inicio"`
	HorarioFim      string `json:"horario_fim"`
	LimiteInscritos int    `json:"limite_inscritos"`
	DiaSemana       string `json:"dia_semana"`
	Sigla           string `json:"sigla"`
	Total_alunos    int    `json:"total_alunos"`

	Local      LocalResponseDto      `json:"local"`
	Modalidade ModalidadeResponseDto `json:"modalidade"`

	Professor string `json:"professor,omitempty"`
}

type ModalidadeResponseDto struct {
	Nome           string  `json:"nome"`
	ValorAluno     float64 `json:"valor_aluno"`
	ValorProfessor float64 `json:"valor_professor"`
}

type LocalResponseDto struct {
	Nome   string `json:"nome"`
	Campus string `json:"campus"`
}
