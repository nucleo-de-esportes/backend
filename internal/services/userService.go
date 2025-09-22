package services

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/nucleo-de-esportes/backend/internal/repository"

	"github.com/nucleo-de-esportes/backend/internal/model"
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email não pode estar vazio")
	}

	emailFormat := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailFormat.MatchString(email) {
		return errors.New("formato de email inválido")
	}

	if !strings.Contains(email, "@sempreceub") && !strings.Contains(email, "@ceub") {
		return errors.New("email institucional deve ser utilizado")
	}

	return nil
}

func CreateUser(user model.User) error {

	if err := repository.DB.Create(&user).Error; err != nil {
		return errors.New("erro ao inserir usuario na tabela")
	}

	return nil

}

func GetSecretKey() string {
	return os.Getenv("SECRET_KEY")
}
