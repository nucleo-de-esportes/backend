package repository

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {

	err := godotenv.Load(filepath.Join("../../", "dbVariables.env"))
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env")
	}

	db_name := os.Getenv("DB_NAME")
	db_password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/Sao_Paulo", db_password, db_name)
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
