package repository

import (
	"fmt"
	"log"

	"github.com/nucleo-de-esportes/backend/internal/config"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg config.DatabaseConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Erro ao tentar conectar com o banco de dados: ", err)
	}

	DB = database
	log.Printf("Conexao feita com sucesso!")

	err = DB.AutoMigrate(
		&model.User{},
		&model.Local{},
		&model.Modalidade{},
		&model.Turma{},
		&model.Matricula{},
		&model.Aula{},
		&model.Presenca{},
	)
	if err != nil {
		log.Printf("Erro na migração das tabelas: %v", err)
	} else {
		log.Printf("Tabelas migradas com sucesso!")
	}
}
