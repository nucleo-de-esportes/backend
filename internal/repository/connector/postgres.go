package connector

import (
	"fmt"

	"github.com/nucleo-de-esportes/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConnector struct{}

func (p *PostgresConnector) Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}
