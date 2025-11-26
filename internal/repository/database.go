package repository

import (
	"log"

	"github.com/nucleo-de-esportes/backend/internal/config"
	"github.com/nucleo-de-esportes/backend/internal/model"
	"github.com/nucleo-de-esportes/backend/internal/repository/connector"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init inicializa a conexão com o banco e executa as migrações.
func Init(cfg config.DatabaseConfig) {
	conn, err := connector.New(cfg.Driver)
	if err != nil {
		log.Fatal(err)
	}

	database, err := conn.Connect(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco (%s): %v", cfg.Driver, err)
	}

	DB = database
	log.Printf("Conexão feita com sucesso usando %s!", cfg.Driver)

	if err := migrate(); err != nil {
		log.Printf("Erro na migração das tabelas: %v", err)
	} else {
		log.Printf("Tabelas migradas com sucesso!")
	}
}

func migrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Local{},
		&model.Modalidade{},
		&model.Turma{},
		&model.Matricula{},
		&model.Aula{},
		&model.Presenca{},
		&model.Aviso{},
	)
}
