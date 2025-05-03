package services

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email não pode estar vazio")
	}

	emailFormat := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailFormat.MatchString(email) {
		return errors.New("formato de email inválido")
	}

	if !strings.Contains(email, "@sempreceub") {
		return errors.New("email institucional deve ser utilizado")
	}

	return nil
}
