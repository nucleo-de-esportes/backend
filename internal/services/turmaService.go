package services

import (
	"strconv"
	"strings"
	"time"
)

func CalcularProximaAula(diaSemana string) time.Time {
	hoje := time.Now()

	// Mapear dias da semana em português para números
	diasSemana := map[string]time.Weekday{
		"domingo": time.Sunday,
		"segunda": time.Monday,
		"terça":   time.Tuesday,
		"terca":   time.Tuesday,
		"quarta":  time.Wednesday,
		"quinta":  time.Thursday,
		"sexta":   time.Friday,
		"sábado":  time.Saturday,
		"sabado":  time.Saturday,
	}

	diaDesejado, existe := diasSemana[strings.ToLower(diaSemana)]
	if !existe {
		// Se não encontrar, retorna amanhã
		return hoje.AddDate(0, 0, 1)
	}

	diasParaProximo := (int(diaDesejado) - int(hoje.Weekday()) + 7) % 7
	if diasParaProximo == 0 {
		diasParaProximo = 7 // Se for hoje, agendar para a próxima semana
	}

	return hoje.AddDate(0, 0, diasParaProximo)
}

func CombinarDataHora(data time.Time, horario string) time.Time {
	// Assumindo que horario está no formato "HH:MM"
	horarioParts := strings.Split(horario, ":")
	if len(horarioParts) != 2 {
		return data
	}

	hora, _ := strconv.Atoi(horarioParts[0])
	minuto, _ := strconv.Atoi(horarioParts[1])

	return time.Date(data.Year(), data.Month(), data.Day(), hora, minuto, 0, 0, data.Location())
}

func ValidarHorarioNadoLivre(inicio, fim string) bool {

	inicioTime, err1 := time.Parse("15:04", inicio)
	fimTime, err2 := time.Parse("15:04", fim)
	limiteInicio, _ := time.Parse("15:04", "11:00")
	limiteFim, _ := time.Parse("15:04", "20:00")

	if err1 != nil || err2 != nil {
		return false
	}

	// Verificar se está dentro do range 11h-20h
	return !inicioTime.Before(limiteInicio) &&
		!inicioTime.After(limiteFim) &&
		!fimTime.Before(limiteInicio) &&
		!fimTime.After(limiteFim) &&
		inicioTime.Before(fimTime)
}

// ParseDateParam converte um parâmetro de data (formato "2006-01-02" ou timestamp Unix)
// para um intervalo de data (início e fim do dia)
func ParseDateParam(diaParam string) (dataInicio, dataFim time.Time, err error) {
	now := time.Now()

	// Se não foi fornecido parâmetro, usar dia atual
	if diaParam == "" {
		dataInicio = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		dataFim = dataInicio.Add(24 * time.Hour)
		return dataInicio, dataFim, nil
	}

	// Tentar parse como data no formato "2006-01-02"
	parsedTime, err := time.Parse("2006-01-02", diaParam)
	if err == nil {
		dataInicio = time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, parsedTime.Location())
		dataFim = dataInicio.Add(24 * time.Hour)
		return dataInicio, dataFim, nil
	}

	// Tentar parse como timestamp Unix
	timestamp, errInt := strconv.ParseInt(diaParam, 10, 64)
	if errInt != nil {
		return time.Time{}, time.Time{}, err
	}

	parsedTime = time.Unix(timestamp, 0)
	dataInicio = time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, parsedTime.Location())
	dataFim = dataInicio.Add(24 * time.Hour)
	return dataInicio, dataFim, nil
}
